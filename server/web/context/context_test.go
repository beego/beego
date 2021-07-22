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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/session"
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

func TestContext_Session(t *testing.T) {
	c := NewContext()
	if store, err := c.Session(); store != nil || err == nil {
		t.FailNow()
	}
}

func TestContext_Session1(t *testing.T) {
	c := Context{}
	if store, err := c.Session(); store != nil || err == nil {
		t.FailNow()
	}
}

func TestContext_Session2(t *testing.T) {
	c := NewContext()
	c.Input.CruSession = &session.MemSessionStore{}

	if store, err := c.Session(); store == nil || err != nil {
		t.FailNow()
	}
}

func TestSetCookie(t *testing.T) {
	type cookie struct {
		Name     string
		Value    string
		MaxAge   int64
		Path     string
		Domain   string
		Secure   bool
		HttpOnly bool
		SameSite string
	}
	type testItem struct {
		item cookie
		want string
	}
	cases := []struct {
		request string
		valueGp []testItem
	}{
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "Strict"}, "name=value; Max-Age=0; Path=/; SameSite=Strict"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "Lax"}, "name=value; Max-Age=0; Path=/; SameSite=Lax"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "None"}, "name=value; Max-Age=0; Path=/; SameSite=None"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, ""}, "name=value; Max-Age=0; Path=/"}}},
	}
	for _, c := range cases {
		r, _ := http.NewRequest("GET", c.request, nil)
		output := NewOutput()
		output.Context = NewContext()
		output.Context.Reset(httptest.NewRecorder(), r)
		for _, item := range c.valueGp {
			params := item.item
			var others = []interface{}{params.MaxAge, params.Path, params.Domain, params.Secure, params.HttpOnly, params.SameSite}
			output.Context.SetCookie(params.Name, params.Value, others...)
			got := output.Context.ResponseWriter.Header().Get("Set-Cookie")
			if got != item.want {
				t.Fatalf("SetCookie error,should be:\n%v \ngot:\n%v", item.want, got)
			}
		}
	}
}
