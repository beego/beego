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

package rbac

// RoleManager provides interface to define the operations for managing roles.
type RoleManager interface {
	// Clear clears all stored data and resets the role manager to the initial state.
	Clear() error
	// AddLink adds the inheritance link between two roles. role: name1 and role: name2.
	// domain is a prefix to the roles (can be used for other purposes).
	AddLink(name1 string, name2 string, domain ...string) error
	// DeleteLink deletes the inheritance link between two roles. role: name1 and role: name2.
	// domain is a prefix to the roles (can be used for other purposes).
	DeleteLink(name1 string, name2 string, domain ...string) error
	// HasLink determines whether a link exists between two roles. role: name1 inherits role: name2.
	// domain is a prefix to the roles (can be used for other purposes).
	HasLink(name1 string, name2 string, domain ...string) (bool, error)
	// GetRoles gets the roles that a user inherits.
	// domain is a prefix to the roles (can be used for other purposes).
	GetRoles(name string, domain ...string) ([]string, error)
	// GetUsers gets the users that inherits a role.
	// domain is a prefix to the users (can be used for other purposes).
	GetUsers(name string, domain ...string) ([]string, error)
	// PrintRoles prints all the roles to log.
	PrintRoles() error
}
