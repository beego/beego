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

const (
	notImplemented = "not implemented"
)

// addPolicy adds a rule to the current policy.
func (e *Enforcer) addPolicy(sec string, ptype string, rule []string) bool {
	ruleAdded := e.model.AddPolicy(sec, ptype, rule)
	if !ruleAdded {
		return ruleAdded
	}

	if e.adapter != nil && e.autoSave {
		if err := e.adapter.AddPolicy(sec, ptype, rule); err != nil {
			if err.Error() != notImplemented {
				panic(err)
			}
		}
		if e.watcher != nil {
			// error intentionally ignored
			e.watcher.Update()
		}
	}

	return ruleAdded
}

// removePolicy removes a rule from the current policy.
func (e *Enforcer) removePolicy(sec string, ptype string, rule []string) bool {
	ruleRemoved := e.model.RemovePolicy(sec, ptype, rule)
	if !ruleRemoved {
		return ruleRemoved
	}

	if e.adapter != nil && e.autoSave {
		if err := e.adapter.RemovePolicy(sec, ptype, rule); err != nil {
			if err.Error() != notImplemented {
				panic(err)
			}
		}
		if e.watcher != nil {
			// error intentionally ignored
			e.watcher.Update()
		}
	}

	return ruleRemoved
}

// removeFilteredPolicy removes rules based on field filters from the current policy.
func (e *Enforcer) removeFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) bool {
	ruleRemoved := e.model.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
	if !ruleRemoved {
		return ruleRemoved
	}

	if e.adapter != nil && e.autoSave {
		if err := e.adapter.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...); err != nil {
			if err.Error() != notImplemented {
				panic(err)
			}
		}
		if e.watcher != nil {
			// error intentionally ignored
			e.watcher.Update()
		}
	}

	return ruleRemoved
}
