// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casbin

import (
	"errors"
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/casbin/casbin/effect"
	"github.com/casbin/casbin/log"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/persist/file-adapter"
	"github.com/casbin/casbin/rbac"
	"github.com/casbin/casbin/rbac/default-role-manager"
	"github.com/casbin/casbin/util"
)

// Enforcer is the main interface for authorization enforcement and policy management.
type Enforcer struct {
	modelPath string
	model     model.Model
	fm        model.FunctionMap
	eft       effect.Effector

	adapter persist.Adapter
	watcher persist.Watcher
	rm      rbac.RoleManager

	enabled            bool
	autoSave           bool
	autoBuildRoleLinks bool
}

// NewEnforcer creates an enforcer via file or DB.
// File:
// e := casbin.NewEnforcer("path/to/basic_model.conf", "path/to/basic_policy.csv")
// MySQL DB:
// a := mysqladapter.NewDBAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
// e := casbin.NewEnforcer("path/to/basic_model.conf", a)
func NewEnforcer(params ...interface{}) *Enforcer {
	e := &Enforcer{}

	parsedParamLen := 0
	if len(params) >= 1 {
		enableLog, ok := params[len(params)-1].(bool)
		if ok {
			e.EnableLog(enableLog)

			parsedParamLen++
		}
	}

	if len(params)-parsedParamLen == 2 {
		switch p0 := params[0].(type) {
		case string:
			switch p1 := params[1].(type) {
			case string:
				e.InitWithFile(p0, p1)
			default:
				e.InitWithAdapter(p0, p1.(persist.Adapter))
			}
		default:
			switch params[1].(type) {
			case string:
				panic("Invalid parameters for enforcer.")
			default:
				e.InitWithModelAndAdapter(p0.(model.Model), params[1].(persist.Adapter))
			}
		}
	} else if len(params)-parsedParamLen == 1 {
		switch p0 := params[0].(type) {
		case string:
			e.InitWithFile(p0, "")
		default:
			e.InitWithModelAndAdapter(p0.(model.Model), nil)
		}
	} else if len(params)-parsedParamLen == 0 {
		e.InitWithFile("", "")
	} else {
		panic("Invalid parameters for enforcer.")
	}

	return e
}

// InitWithFile initializes an enforcer with a model file and a policy file.
func (e *Enforcer) InitWithFile(modelPath string, policyPath string) {
	a := fileadapter.NewAdapter(policyPath)
	e.InitWithAdapter(modelPath, a)
}

// InitWithAdapter initializes an enforcer with a database adapter.
func (e *Enforcer) InitWithAdapter(modelPath string, adapter persist.Adapter) {
	m := NewModel(modelPath, "")
	e.InitWithModelAndAdapter(m, adapter)

	e.modelPath = modelPath
}

// InitWithModelAndAdapter initializes an enforcer with a model and a database adapter.
func (e *Enforcer) InitWithModelAndAdapter(m model.Model, adapter persist.Adapter) {
	e.adapter = adapter

	e.model = m
	e.model.PrintModel()
	e.fm = model.LoadFunctionMap()

	e.initialize()

	if e.adapter != nil {
		// error intentionally ignored
		e.LoadPolicy()
	}
}

func (e *Enforcer) initialize() {
	e.rm = defaultrolemanager.NewRoleManager(10)
	e.eft = effect.NewDefaultEffector()
	e.watcher = nil

	e.enabled = true
	e.autoSave = true
	e.autoBuildRoleLinks = true
}

// NewModel creates a model.
func NewModel(text ...string) model.Model {
	m := make(model.Model)

	if len(text) == 2 {
		if text[0] != "" {
			m.LoadModel(text[0])
		}
	} else if len(text) == 1 {
		m.LoadModelFromText(text[0])
	} else if len(text) != 0 {
		panic("Invalid parameters for model.")
	}

	return m
}

// LoadModel reloads the model from the model CONF file.
// Because the policy is attached to a model, so the policy is invalidated and needs to be reloaded by calling LoadPolicy().
func (e *Enforcer) LoadModel() {
	e.model = NewModel()
	e.model.LoadModel(e.modelPath)
	e.model.PrintModel()
	e.fm = model.LoadFunctionMap()
}

// GetModel gets the current model.
func (e *Enforcer) GetModel() model.Model {
	return e.model
}

// SetModel sets the current model.
func (e *Enforcer) SetModel(m model.Model) {
	e.model = m
	e.fm = model.LoadFunctionMap()
}

// GetAdapter gets the current adapter.
func (e *Enforcer) GetAdapter() persist.Adapter {
	return e.adapter
}

// SetAdapter sets the current adapter.
func (e *Enforcer) SetAdapter(adapter persist.Adapter) {
	e.adapter = adapter
}

// SetWatcher sets the current watcher.
func (e *Enforcer) SetWatcher(watcher persist.Watcher) {
	e.watcher = watcher
	// error intentionally ignored
	watcher.SetUpdateCallback(func(string) { e.LoadPolicy() })
}

// SetRoleManager sets the current role manager.
func (e *Enforcer) SetRoleManager(rm rbac.RoleManager) {
	e.rm = rm
}

// SetEffector sets the current effector.
func (e *Enforcer) SetEffector(eft effect.Effector) {
	e.eft = eft
}

// ClearPolicy clears all policy.
func (e *Enforcer) ClearPolicy() {
	e.model.ClearPolicy()
}

// LoadPolicy reloads the policy from file/database.
func (e *Enforcer) LoadPolicy() error {
	e.model.ClearPolicy()
	if err := e.adapter.LoadPolicy(e.model); err != nil {
		return err
	}

	e.model.PrintPolicy()
	if e.autoBuildRoleLinks {
		e.BuildRoleLinks()
	}
	return nil
}

// LoadFilteredPolicy reloads a filtered policy from file/database.
func (e *Enforcer) LoadFilteredPolicy(filter interface{}) error {
	e.model.ClearPolicy()

	var filteredAdapter persist.FilteredAdapter

	// Attempt to cast the Adapter as a FilteredAdapter
	switch adapter := e.adapter.(type) {
	case persist.FilteredAdapter:
		filteredAdapter = adapter
	default:
		return errors.New("filtered policies are not supported by this adapter")
	}
	if err := filteredAdapter.LoadFilteredPolicy(e.model, filter); err != nil {
		return err
	}

	e.model.PrintPolicy()
	if e.autoBuildRoleLinks {
		e.BuildRoleLinks()
	}
	return nil
}

// IsFiltered returns true if the loaded policy has been filtered.
func (e *Enforcer) IsFiltered() bool {
	filteredAdapter, ok := e.adapter.(persist.FilteredAdapter)
	if !ok {
		return false
	}
	return filteredAdapter.IsFiltered()
}

// SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.
func (e *Enforcer) SavePolicy() error {
	if e.IsFiltered() {
		return errors.New("cannot save a filtered policy")
	}
	if err := e.adapter.SavePolicy(e.model); err != nil {
		return err
	}
	if e.watcher != nil {
		return e.watcher.Update()
	}
	return nil
}

// EnableEnforce changes the enforcing state of Casbin, when Casbin is disabled, all access will be allowed by the Enforce() function.
func (e *Enforcer) EnableEnforce(enable bool) {
	e.enabled = enable
}

// EnableLog changes whether Casbin will log messages to the Logger.
func (e *Enforcer) EnableLog(enable bool) {
	log.GetLogger().EnableLog(enable)
}

// EnableAutoSave controls whether to save a policy rule automatically to the adapter when it is added or removed.
func (e *Enforcer) EnableAutoSave(autoSave bool) {
	e.autoSave = autoSave
}

// EnableAutoBuildRoleLinks controls whether to rebuild the role inheritance relations when a role is added or deleted.
func (e *Enforcer) EnableAutoBuildRoleLinks(autoBuildRoleLinks bool) {
	e.autoBuildRoleLinks = autoBuildRoleLinks
}

// BuildRoleLinks manually rebuild the role inheritance relations.
func (e *Enforcer) BuildRoleLinks() {
	// error intentionally ignored
	e.rm.Clear()
	e.model.BuildRoleLinks(e.rm)
}

// Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).
func (e *Enforcer) Enforce(rvals ...interface{}) bool {
	if !e.enabled {
		return true
	}

	functions := make(map[string]govaluate.ExpressionFunction)
	for key, function := range e.fm {
		functions[key] = function
	}
	if _, ok := e.model["g"]; ok {
		for key, ast := range e.model["g"] {
			rm := ast.RM
			functions[key] = util.GenerateGFunction(rm)
		}
	}

	expString := e.model["m"]["m"].Value
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expString, functions)
	if err != nil {
		panic(err)
	}

	var policyEffects []effect.Effect
	var matcherResults []float64
	if policyLen := len(e.model["p"]["p"].Policy); policyLen != 0 {
		policyEffects = make([]effect.Effect, policyLen)
		matcherResults = make([]float64, policyLen)

		for i, pvals := range e.model["p"]["p"].Policy {
			// log.LogPrint("Policy Rule: ", pvals)

			parameters := make(map[string]interface{}, 8)
			for j, token := range e.model["r"]["r"].Tokens {
				parameters[token] = rvals[j]
			}
			for j, token := range e.model["p"]["p"].Tokens {
				parameters[token] = pvals[j]
			}

			result, err := expression.Evaluate(parameters)
			// log.LogPrint("Result: ", result)

			if err != nil {
				policyEffects[i] = effect.Indeterminate
				panic(err)
			}

			switch result := result.(type) {
			case bool:
				if !result {
					policyEffects[i] = effect.Indeterminate
					continue
				}
			case float64:
				if result == 0 {
					policyEffects[i] = effect.Indeterminate
					continue
				} else {
					matcherResults[i] = result
				}
			default:
				panic(errors.New("matcher result should be bool, int or float"))
			}

			if eft, ok := parameters["p_eft"]; ok {
				if eft == "allow" {
					policyEffects[i] = effect.Allow
				} else if eft == "deny" {
					policyEffects[i] = effect.Deny
				} else {
					policyEffects[i] = effect.Indeterminate
				}
			} else {
				policyEffects[i] = effect.Allow
			}

			if e.model["e"]["e"].Value == "priority(p_eft) || deny" {
				break
			}

		}
	} else {
		policyEffects = make([]effect.Effect, 1)
		matcherResults = make([]float64, 1)

		parameters := make(map[string]interface{}, 8)
		for j, token := range e.model["r"]["r"].Tokens {
			parameters[token] = rvals[j]
		}
		for _, token := range e.model["p"]["p"].Tokens {
			parameters[token] = ""
		}

		result, err := expression.Evaluate(parameters)
		// log.LogPrint("Result: ", result)

		if err != nil {
			policyEffects[0] = effect.Indeterminate
			panic(err)
		}

		if result.(bool) {
			policyEffects[0] = effect.Allow
		} else {
			policyEffects[0] = effect.Indeterminate
		}
	}

	// log.LogPrint("Rule Results: ", policyEffects)

	result, err := e.eft.MergeEffects(e.model["e"]["e"].Value, policyEffects, matcherResults)
	if err != nil {
		panic(err)
	}

	// Log request.
	if log.GetLogger().IsEnabled() {
		reqStr := "Request: "
		for i, rval := range rvals {
			if i != len(rvals)-1 {
				reqStr += fmt.Sprintf("%v, ", rval)
			} else {
				reqStr += fmt.Sprintf("%v", rval)
			}
		}
		reqStr += fmt.Sprintf(" ---> %t", result)
		log.LogPrint(reqStr)
	}

	return result
}
