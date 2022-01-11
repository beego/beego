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

package prometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/server/web/context"
)

func TestFilterChain(t *testing.T) {
	filter := (&FilterChainBuilder{}).FilterChain(func(ctx *context.Context) {
		// do nothing
		ctx.Input.SetData("invocation", true)
	})

	ctx := context.NewContext()
	r, _ := http.NewRequest("GET", "/prometheus/user", nil)
	w := httptest.NewRecorder()
	ctx.Reset(w, r)
	ctx.Input.SetData("RouterPattern", "my-route")
	filter(ctx)
	assert.True(t, ctx.Input.GetData("invocation").(bool))
	time.Sleep(1 * time.Second)
}

func TestFilterChainBuilder_report(t *testing.T) {
	ctx := context.NewContext()
	r, _ := http.NewRequest("GET", "/prometheus/user", nil)
	w := httptest.NewRecorder()
	ctx.Reset(w, r)
	fb := &FilterChainBuilder{}
	// without router info
	report(time.Second, ctx, fb.buildVec())

	ctx.Input.SetData("RouterPattern", "my-route")
	report(time.Second, ctx, fb.buildVec())
}
