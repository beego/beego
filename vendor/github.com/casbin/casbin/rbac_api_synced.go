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

// GetRolesForUser gets the roles that a user has.
func (e *SyncedEnforcer) GetRolesForUser(name string) []string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetRolesForUser(name)
}

// GetUsersForRole gets the users that has a role.
func (e *SyncedEnforcer) GetUsersForRole(name string) []string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetUsersForRole(name)
}

// HasRoleForUser determines whether a user has a role.
func (e *SyncedEnforcer) HasRoleForUser(name string, role string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.HasRoleForUser(name, role)
}

// AddRoleForUser adds a role for a user.
// Returns false if the user already has the role (aka not affected).
func (e *SyncedEnforcer) AddRoleForUser(user string, role string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.AddRoleForUser(user, role)
}

// DeleteRoleForUser deletes a role for a user.
// Returns false if the user does not have the role (aka not affected).
func (e *SyncedEnforcer) DeleteRoleForUser(user string, role string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeleteRoleForUser(user, role)
}

// DeleteRolesForUser deletes all roles for a user.
// Returns false if the user does not have any roles (aka not affected).
func (e *SyncedEnforcer) DeleteRolesForUser(user string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeleteRolesForUser(user)
}

// DeleteUser deletes a user.
// Returns false if the user does not exist (aka not affected).
func (e *SyncedEnforcer) DeleteUser(user string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeleteUser(user)
}

// DeleteRole deletes a role.
func (e *SyncedEnforcer) DeleteRole(role string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.Enforcer.DeleteRole(role)
}

// DeletePermission deletes a permission.
// Returns false if the permission does not exist (aka not affected).
func (e *SyncedEnforcer) DeletePermission(permission ...string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeletePermission(permission...)
}

// AddPermissionForUser adds a permission for a user or role.
// Returns false if the user or role already has the permission (aka not affected).
func (e *SyncedEnforcer) AddPermissionForUser(user string, permission ...string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.AddPermissionForUser(user, permission...)
}

// DeletePermissionForUser deletes a permission for a user or role.
// Returns false if the user or role does not have the permission (aka not affected).
func (e *SyncedEnforcer) DeletePermissionForUser(user string, permission ...string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeletePermissionForUser(user, permission...)
}

// DeletePermissionsForUser deletes permissions for a user or role.
// Returns false if the user or role does not have any permissions (aka not affected).
func (e *SyncedEnforcer) DeletePermissionsForUser(user string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.DeletePermissionsForUser(user)
}

// GetPermissionsForUser gets permissions for a user or role.
func (e *SyncedEnforcer) GetPermissionsForUser(user string) [][]string {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.GetPermissionsForUser(user)
}

// HasPermissionForUser determines whether a user has a permission.
func (e *SyncedEnforcer) HasPermissionForUser(user string, permission ...string) bool {
	e.m.Lock()
	defer e.m.Unlock()
	return e.Enforcer.HasPermissionForUser(user, permission...)
}
