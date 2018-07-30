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

// GetAllSubjects gets the list of subjects that show up in the current policy.
func (e *Enforcer) GetAllSubjects() []string {
	return e.GetAllNamedSubjects("p")
}

// GetAllNamedSubjects gets the list of subjects that show up in the current named policy.
func (e *Enforcer) GetAllNamedSubjects(ptype string) []string {
	return e.model.GetValuesForFieldInPolicy("p", ptype, 0)
}

// GetAllObjects gets the list of objects that show up in the current policy.
func (e *Enforcer) GetAllObjects() []string {
	return e.GetAllNamedObjects("p")
}

// GetAllNamedObjects gets the list of objects that show up in the current named policy.
func (e *Enforcer) GetAllNamedObjects(ptype string) []string {
	return e.model.GetValuesForFieldInPolicy("p", ptype, 1)
}

// GetAllActions gets the list of actions that show up in the current policy.
func (e *Enforcer) GetAllActions() []string {
	return e.GetAllNamedActions("p")
}

// GetAllNamedActions gets the list of actions that show up in the current named policy.
func (e *Enforcer) GetAllNamedActions(ptype string) []string {
	return e.model.GetValuesForFieldInPolicy("p", ptype, 2)
}

// GetAllRoles gets the list of roles that show up in the current policy.
func (e *Enforcer) GetAllRoles() []string {
	return e.GetAllNamedRoles("g")
}

// GetAllNamedRoles gets the list of roles that show up in the current named policy.
func (e *Enforcer) GetAllNamedRoles(ptype string) []string {
	return e.model.GetValuesForFieldInPolicy("g", ptype, 1)
}

// GetPolicy gets all the authorization rules in the policy.
func (e *Enforcer) GetPolicy() [][]string {
	return e.GetNamedPolicy("p")
}

// GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.
func (e *Enforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string {
	return e.GetFilteredNamedPolicy("p", fieldIndex, fieldValues...)
}

// GetNamedPolicy gets all the authorization rules in the named policy.
func (e *Enforcer) GetNamedPolicy(ptype string) [][]string {
	return e.model.GetPolicy("p", ptype)
}

// GetFilteredNamedPolicy gets all the authorization rules in the named policy, field filters can be specified.
func (e *Enforcer) GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) [][]string {
	return e.model.GetFilteredPolicy("p", ptype, fieldIndex, fieldValues...)
}

// GetGroupingPolicy gets all the role inheritance rules in the policy.
func (e *Enforcer) GetGroupingPolicy() [][]string {
	return e.GetNamedGroupingPolicy("g")
}

// GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.
func (e *Enforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) [][]string {
	return e.GetFilteredNamedGroupingPolicy("g", fieldIndex, fieldValues...)
}

// GetNamedGroupingPolicy gets all the role inheritance rules in the policy.
func (e *Enforcer) GetNamedGroupingPolicy(ptype string) [][]string {
	return e.model.GetPolicy("g", ptype)
}

// GetFilteredNamedGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.
func (e *Enforcer) GetFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) [][]string {
	return e.model.GetFilteredPolicy("g", ptype, fieldIndex, fieldValues...)
}

// HasPolicy determines whether an authorization rule exists.
func (e *Enforcer) HasPolicy(params ...interface{}) bool {
	return e.HasNamedPolicy("p", params...)
}

// HasNamedPolicy determines whether a named authorization rule exists.
func (e *Enforcer) HasNamedPolicy(ptype string, params ...interface{}) bool {
	if strSlice, ok := params[0].([]string); len(params) == 1 && ok {
		return e.model.HasPolicy("p", ptype, strSlice)
	}

	policy := make([]string, 0)
	for _, param := range params {
		policy = append(policy, param.(string))
	}

	return e.model.HasPolicy("p", ptype, policy)
}

// AddPolicy adds an authorization rule to the current policy.
// If the rule already exists, the function returns false and the rule will not be added.
// Otherwise the function returns true by adding the new rule.
func (e *Enforcer) AddPolicy(params ...interface{}) bool {
	return e.AddNamedPolicy("p", params...)
}

// AddNamedPolicy adds an authorization rule to the current named policy.
// If the rule already exists, the function returns false and the rule will not be added.
// Otherwise the function returns true by adding the new rule.
func (e *Enforcer) AddNamedPolicy(ptype string, params ...interface{}) bool {
	var ruleAdded bool
	if strSlice, ok := params[0].([]string); len(params) == 1 && ok {
		ruleAdded = e.addPolicy("p", ptype, strSlice)
	} else {
		policy := make([]string, 0)
		for _, param := range params {
			policy = append(policy, param.(string))
		}

		ruleAdded = e.addPolicy("p", ptype, policy)
	}

	return ruleAdded
}

// RemovePolicy removes an authorization rule from the current policy.
func (e *Enforcer) RemovePolicy(params ...interface{}) bool {
	return e.RemoveNamedPolicy("p", params...)
}

// RemoveFilteredPolicy removes an authorization rule from the current policy, field filters can be specified.
func (e *Enforcer) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) bool {
	return e.RemoveFilteredNamedPolicy("p", fieldIndex, fieldValues...)
}

// RemoveNamedPolicy removes an authorization rule from the current named policy.
func (e *Enforcer) RemoveNamedPolicy(ptype string, params ...interface{}) bool {
	var ruleRemoved bool
	if strSlice, ok := params[0].([]string); len(params) == 1 && ok {
		ruleRemoved = e.removePolicy("p", ptype, strSlice)
	} else {
		policy := make([]string, 0)
		for _, param := range params {
			policy = append(policy, param.(string))
		}

		ruleRemoved = e.removePolicy("p", ptype, policy)
	}

	return ruleRemoved
}

// RemoveFilteredNamedPolicy removes an authorization rule from the current named policy, field filters can be specified.
func (e *Enforcer) RemoveFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) bool {
	return e.removeFilteredPolicy("p", ptype, fieldIndex, fieldValues...)
}

// HasGroupingPolicy determines whether a role inheritance rule exists.
func (e *Enforcer) HasGroupingPolicy(params ...interface{}) bool {
	return e.HasNamedGroupingPolicy("g", params...)
}

// HasNamedGroupingPolicy determines whether a named role inheritance rule exists.
func (e *Enforcer) HasNamedGroupingPolicy(ptype string, params ...interface{}) bool {
	if strSlice, ok := params[0].([]string); len(params) == 1 && ok {
		return e.model.HasPolicy("g", ptype, strSlice)
	}

	policy := make([]string, 0)
	for _, param := range params {
		policy = append(policy, param.(string))
	}

	return e.model.HasPolicy("g", ptype, policy)
}

// AddGroupingPolicy adds a role inheritance rule to the current policy.
// If the rule already exists, the function returns false and the rule will not be added.
// Otherwise the function returns true by adding the new rule.
func (e *Enforcer) AddGroupingPolicy(params ...interface{}) bool {
	return e.AddNamedGroupingPolicy("g", params...)
}

// AddNamedGroupingPolicy adds a named role inheritance rule to the current policy.
// If the rule already exists, the function returns false and the rule will not be added.
// Otherwise the function returns true by adding the new rule.
func (e *Enforcer) AddNamedGroupingPolicy(ptype string, params ...interface{}) bool {
	var ruleAdded bool
	if strSlice, ok := params[0].([]string); len(params) == 1 && ok {
		ruleAdded = e.addPolicy("g", ptype, strSlice)
	} else {
		policy := make([]string, 0)
		for _, param := range params {
			policy = append(policy, param.(string))
		}

		ruleAdded = e.addPolicy("g", ptype, policy)
	}

	if e.autoBuildRoleLinks {
		e.BuildRoleLinks()
	}
	return ruleAdded
}

// RemoveGroupingPolicy removes a role inheritance rule from the current policy.
func (e *Enforcer) RemoveGroupingPolicy(params ...interface{}) bool {
	return e.RemoveNamedGroupingPolicy("g", params...)
}

// RemoveFilteredGroupingPolicy removes a role inheritance rule from the current policy, field filters can be specified.
func (e *Enforcer) RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) bool {
	return e.RemoveFilteredNamedGroupingPolicy("g", fieldIndex, fieldValues...)
}

// RemoveNamedGroupingPolicy removes a role inheritance rule from the current named policy.
func (e *Enforcer) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) bool {
	var ruleRemoved bool
	if strSlice, ok := params[0].([]string); len(params) == 1 && ok {
		ruleRemoved = e.removePolicy("g", ptype, strSlice)
	} else {
		policy := make([]string, 0)
		for _, param := range params {
			policy = append(policy, param.(string))
		}

		ruleRemoved = e.removePolicy("g", ptype, policy)
	}

	if e.autoBuildRoleLinks {
		e.BuildRoleLinks()
	}
	return ruleRemoved
}

// RemoveFilteredNamedGroupingPolicy removes a role inheritance rule from the current named policy, field filters can be specified.
func (e *Enforcer) RemoveFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) bool {
	ruleRemoved := e.removeFilteredPolicy("g", ptype, fieldIndex, fieldValues...)

	if e.autoBuildRoleLinks {
		e.BuildRoleLinks()
	}
	return ruleRemoved
}

// AddFunction adds a customized function.
func (e *Enforcer) AddFunction(name string, function func(args ...interface{}) (interface{}, error)) {
	e.fm.AddFunction(name, function)
}
