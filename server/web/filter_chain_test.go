// Copyright 2020
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

package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func TestControllerRegisterInsertFilterChain(t *testing.T) {

	InsertFilterChain("/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", "filter-chain")
			next(ctx)
		}
	})

	ns := NewNamespace("/chain")

	ns.Get("/*", func(ctx *context.Context) {
		ctx.Output.Body([]byte("hello"))
	})

	r, _ := http.NewRequest("GET", "/chain/user", nil)
	w := httptest.NewRecorder()

	BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, "filter-chain", w.Header().Get("filter"))
}

func TestFilterChainRouter(t *testing.T) {


	app := NewHttpSever()

	const filterNonMatch = "filter-chain-non-match"
	app.InsertFilterChain("/app/nonMatch/before/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", filterNonMatch)
			next(ctx)
		}
	})

	const filterAll = "filter-chain-all"
	app.InsertFilterChain("/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", filterAll)
			next(ctx)
		}
	})

	app.InsertFilterChain("/app/nonMatch/after/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", filterNonMatch)
			next(ctx)
		}
	})

	app.InsertFilterChain("/app/match/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("match", "yes")
			next(ctx)
		}
	})

	r, _ := http.NewRequest("GET", "/app/match", nil)
	w := httptest.NewRecorder()

	app.Handlers.ServeHTTP(w, r)
	assert.Equal(t, filterAll, w.Header().Get("filter"))
	assert.Equal(t, "yes", w.Header().Get("match"))

	r, _ = http.NewRequest("GET", "/app/match1", nil)
	w = httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, r)
	assert.Equal(t, filterAll, w.Header().Get("filter"))
	assert.NotEqual(t, "yes", w.Header().Get("match"))

	r, _ = http.NewRequest("GET", "/app/nonMatch", nil)
	w = httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, r)
	assert.Equal(t, filterAll, w.Header().Get("filter"))
	assert.NotEqual(t, "yes", w.Header().Get("match"))
}
