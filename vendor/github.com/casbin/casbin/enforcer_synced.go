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
	"log"
	"sync"
	"time"

	"github.com/casbin/casbin/persist"
)

// SyncedEnforcer wraps Enforcer and provides synchronized access
type SyncedEnforcer struct {
	*Enforcer
	m        sync.RWMutex
	autoLoad bool
}

// NewSyncedEnforcer creates a synchronized enforcer via file or DB.
func NewSyncedEnforcer(params ...interface{}) *SyncedEnforcer {
	e := &SyncedEnforcer{}
	e.Enforcer = NewEnforcer(params...)
	e.autoLoad = false
	return e
}

// StartAutoLoadPolicy starts a go routine that will every specified duration call LoadPolicy
func (e *SyncedEnforcer) StartAutoLoadPolicy(d time.Duration) {
	e.autoLoad = true
	go func() {
		n := 1
		log.Print("Start automatically load policy")
		for {
			if !e.autoLoad {
				log.Print("Stop automatically load policy")
				break
			}

			// error intentionally ignored
			e.LoadPolicy()
			// Uncomment this line to see when the policy is loaded.
			// log.Print("Load policy for time: ", n)
			n++
			time.Sleep(d)
		}
	}()
}

// StopAutoLoadPolicy causes the go routine to exit.
func (e *SyncedEnforcer) StopAutoLoadPolicy() {
	e.autoLoad = false
}

// SetWatcher sets the current watcher.
func (e *SyncedEnforcer) SetWatcher(watcher persist.Watcher) {
	e.watcher = watcher
	// error intentionally ignored
	watcher.SetUpdateCallback(func(string) { e.LoadPolicy() })
}

// ClearPolicy clears all policy.
func (e *SyncedEnforcer) ClearPolicy() {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.ClearPolicy()
}

// LoadPolicy reloads the policy from file/database.
func (e *SyncedEnforcer) LoadPolicy() error {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.LoadPolicy()
}

// SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.
func (e *SyncedEnforcer) SavePolicy() error {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.SavePolicy()
}

// BuildRoleLinks manually rebuild the role inheritance relations.
func (e *SyncedEnforcer) BuildRoleLinks() {
	e.m.RLock()
	defer e.m.RUnlock()
	e.Enforcer.BuildRoleLinks()
}

// Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).
func (e *SyncedEnforcer) Enforce(rvals ...interface{}) bool {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.Enforce(rvals...)
}

// GetAllSubjects gets the list of subjects that show up in the current policy.
func (e *SyncedEnforcer) GetAllSubjects() []string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetAllSubjects()
}

// GetAllObjects gets the list of objects that show up in the current policy.
func (e *SyncedEnforcer) GetAllObjects() []string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetAllObjects()
}

// GetAllActions gets the list of actions that show up in the current policy.
func (e *SyncedEnforcer) GetAllActions() []string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetAllActions()
}

// GetAllRoles gets the list of roles that show up in the current policy.
func (e *SyncedEnforcer) GetAllRoles() []string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetAllRoles()
}

// GetPolicy gets all the authorization rules in the policy.
func (e *SyncedEnforcer) GetPolicy() [][]string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetPolicy()
}

// GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.
func (e *SyncedEnforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
}

// GetGroupingPolicy gets all the role inheritance rules in the policy.
func (e *SyncedEnforcer) GetGroupingPolicy() [][]string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetGroupingPolicy()
}

// GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.
func (e *SyncedEnforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) [][]string {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.GetFilteredGroupingPolicy(fieldIndex, fieldValues...)
}

// HasPolicy determines whether an authorization rule exists.
func (e *SyncedEnforcer) HasPolicy(params ...interface{}) bool {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.HasPolicy(params...)
}

// AddPolicy adds an authorization rule to the current policy.
// If the rule already exists, the function returns false and the rule will not be added.
// Otherwise the function returns true by adding the new rule.
func (e *SyncedEnforcer) AddPolicy(params ...interface{}) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.AddPolicy(params...)
}

// RemovePolicy removes an authorization rule from the current policy.
func (e *SyncedEnforcer) RemovePolicy(params ...interface{}) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RemovePolicy(params...)
}

// RemoveFilteredPolicy removes an authorization rule from the current policy, field filters can be specified.
func (e *SyncedEnforcer) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
}

// HasGroupingPolicy determines whether a role inheritance rule exists.
func (e *SyncedEnforcer) HasGroupingPolicy(params ...interface{}) bool {
	e.m.RLock()
	defer e.m.RUnlock()
	return e.Enforcer.HasGroupingPolicy(params...)
}

// AddGroupingPolicy adds a role inheritance rule to the current policy.
// If the rule already exists, the function returns false and the rule will not be added.
// Otherwise the function returns true by adding the new rule.
func (e *SyncedEnforcer) AddGroupingPolicy(params ...interface{}) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.AddGroupingPolicy(params...)
}

// RemoveGroupingPolicy removes a role inheritance rule from the current policy.
func (e *SyncedEnforcer) RemoveGroupingPolicy(params ...interface{}) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RemoveGroupingPolicy(params...)
}

// RemoveFilteredGroupingPolicy removes a role inheritance rule from the current policy, field filters can be specified.
func (e *SyncedEnforcer) RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.RemoveFilteredGroupingPolicy(fieldIndex, fieldValues...)
}
