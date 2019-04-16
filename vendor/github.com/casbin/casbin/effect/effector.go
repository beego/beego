// Copyright 2018 The casbin Authors. All Rights Reserved.
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

package effect

// Effect is the result for a policy rule.
type Effect int

// Values for policy effect.
const (
	Allow Effect = iota
	Indeterminate
	Deny
)

// Effector is the interface for Casbin effectors.
type Effector interface {
	// MergeEffects merges all matching results collected by the enforcer into a single decision.
	MergeEffects(expr string, effects []Effect, results []float64) (bool, error)
}
