// Copyright 2016 beego Author. All Rights Reserved.
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

package context

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/session"
	webSession "github.com/beego/beego/v2/server/web/session"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXsrfReset_01(t *testing.T) {
	r := &http.Request{}
	c := NewContext()
	c.Request = r
	c.ResponseWriter = &Response{}
	c.ResponseWriter.reset(httptest.NewRecorder())
	c.Output.Reset(c)
	c.Input.Reset(c)
	c.XSRFToken("key", 16)
	if c._xsrfToken == "" {
		t.FailNow()
	}
	token := c._xsrfToken
	c.Reset(&Response{ResponseWriter: httptest.NewRecorder()}, r)
	if c._xsrfToken != "" {
		t.FailNow()
	}
	c.XSRFToken("key", 16)
	if c._xsrfToken == "" {
		t.FailNow()
	}
	if token == c._xsrfToken {
		t.FailNow()
	}
}

func testRequest(t *testing.T, handler *web.ControllerRegister, path string, method string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != code {
		t.Errorf("%s, %s: %d, supposed to be %d", path, method, w.Code, code)
	}
}

func TestContext_Session(t *testing.T) {
	handler := web.NewControllerRegister()

	handler.InsertFilterChain(
		"*",
		session.Session(
			webSession.ProviderMemory,
			webSession.CfgCookieName(`go_session_id`),
			webSession.CfgSetCookie(true),
			webSession.CfgGcLifeTime(3600),
			webSession.CfgMaxLifeTime(3600),
			webSession.CfgSecure(false),
			webSession.CfgCookieLifeTime(3600),
		),
	)
	handler.InsertFilterChain(
		"*",
		func(next web.FilterFunc) web.FilterFunc {
			return func(ctx *Context) {
				if _, err := ctx.Session(); err == nil {
					t.Error()
				}

			}
		},
	)
	handler.Any("*", func(ctx *Context) {
		ctx.Output.SetStatus(200)
	})

	testRequest(t, handler, "/dataset1/resource1", "GET", 200)
}