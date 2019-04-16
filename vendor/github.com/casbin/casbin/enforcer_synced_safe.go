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

// NewSyncedEnforcerSafe creates a synchronized enforcer via file or DB.
func NewSyncedEnforcerSafe(params ...interface{}) (enforcer *SyncedEnforcer, err error) {
	e := &SyncedEnforcer{}
	e.Enforcer, err = NewEnforcerSafe(params...)

	if err != nil {
		return nil, err
	}

	e.autoLoad = false
	return e, nil
}

// LoadModelSafe calls LoadModel in a safe way, returns error instead of causing panic.
func (e *SyncedEnforcer) LoadModelSafe() (err error) {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.LoadModelSafe()
}

// EnforceSafe calls Enforce in a safe way, returns error instead of causing panic.
func (e *SyncedEnforcer) EnforceSafe(rvals ...interface{}) (result bool, err error) {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.EnforceSafe(rvals...)
}

// AddPolicySafe calls AddPolicy in a safe way, returns error instead of causing panic.
func (e *SyncedEnforcer) AddPolicySafe(params ...interface{}) (result bool, err error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.AddPolicySafe(params...)
}

// RemovePolicySafe calls RemovePolicy in a safe way, returns error instead of causing panic.
func (e *SyncedEnforcer) RemovePolicySafe(params ...interface{}) (result bool, err error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RemovePolicySafe(params...)
}

// RemoveFilteredPolicySafe calls RemoveFilteredPolicy in a safe way, returns error instead of causing panic.
func (e *SyncedEnforcer) RemoveFilteredPolicySafe(fieldIndex int, fieldValues ...string) (result bool, err error) {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RemoveFilteredPolicySafe(fieldIndex, fieldValues...)
}
