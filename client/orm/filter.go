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
)

// FilterChain is used to build a Filter
// don't forget to call next(...) inside your Filter
type FilterChain func(next Filter) Filter

// Filter behavior is a little big strange.
// it's only be called when users call methods of Ormer
// return value is an array. it's a little bit hard to understand,
// for example, the Ormer's Read method only return error
// so the filter processing this method should return an array whose first element is error
// and, Ormer's ReadOrCreateWithCtx return three values, so the Filter's result should contain three values
type Filter func(ctx context.Context, inv *Invocation) []interface{}

var globalFilterChains = make([]FilterChain, 0, 4)

// AddGlobalFilterChain adds a new FilterChain
// All orm instances built after this invocation will use this filterChain,
// but instances built before this invocation will not be affected
func AddGlobalFilterChain(filterChain ...FilterChain) {
	globalFilterChains = append(globalFilterChains, filterChain...)
}
