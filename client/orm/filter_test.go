// Copyright 2020 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGlobalFilterChain(t *testing.T) {
	AddGlobalFilterChain(func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) []interface{} {
			return next(ctx, inv)
		}
	})
	assert.Equal(t, 1, len(globalFilterChains))
	globalFilterChains = nil
}
