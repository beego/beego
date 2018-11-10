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

import (
	"github.com/casbin/casbin/model"
)

// FilteredAdapter is the interface for Casbin adapters supporting filtered policies.
type FilteredAdapter interface {
	Adapter

	// LoadFilteredPolicy loads only policy rules that match the filter.
	LoadFilteredPolicy(model model.Model, filter interface{}) error
	// IsFiltered returns true if the loaded policy has been filtered.
	IsFiltered() bool
}
