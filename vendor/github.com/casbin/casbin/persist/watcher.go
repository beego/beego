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

package persist

// Watcher is the interface for Casbin watchers.
type Watcher interface {
	// SetUpdateCallback sets the callback function that the watcher will call
	// when the policy in DB has been changed by other instances.
	// A classic callback is Enforcer.LoadPolicy().
	SetUpdateCallback(func(string)) error
	// Update calls the update callback of other instances to synchronize their policy.
	// It is usually called after changing the policy in DB, like Enforcer.SavePolicy(),
	// Enforcer.AddPolicy(), Enforcer.RemovePolicy(), etc.
	Update() error
}
