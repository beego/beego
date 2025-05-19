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
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpServerWithCfg(t *testing.T) {
	BConfig.AppName = "Before"
	svr := NewHttpServerWithCfg(BConfig)
	svr.Cfg.AppName = "hello"
	assert.Equal(t, "hello", BConfig.AppName)
}

func TestServerCtrlGet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	CtrlGet("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlGet can't run")
	}
}

func TestServerCtrlPost(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/user", nil)
	w := httptest.NewRecorder()

	CtrlPost("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlPost can't run")
	}
}

func TestServerCtrlHead(t *testing.T) {
	r, _ := http.NewRequest(http.MethodHead, "/user", nil)
	w := httptest.NewRecorder()

	CtrlHead("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlHead can't run")
	}
}

func TestServerCtrlPut(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPut, "/user", nil)
	w := httptest.NewRecorder()

	CtrlPut("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlPut can't run")
	}
}

func TestServerCtrlPatch(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPatch, "/user", nil)
	w := httptest.NewRecorder()

	CtrlPatch("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlPatch can't run")
	}
}

func TestServerCtrlDelete(t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, "/user", nil)
	w := httptest.NewRecorder()

	CtrlDelete("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlDelete can't run")
	}
}

func TestServerCtrlAny(t *testing.T) {
	CtrlAny("/user", ExampleController.Ping)

	for method := range HTTPMETHOD {
		r, _ := http.NewRequest(method, "/user", nil)
		w := httptest.NewRecorder()
		BeeApp.Handlers.ServeHTTP(w, r)
		if w.Body.String() != exampleBody {
			t.Errorf("TestServerCtrlAny can't run")
		}
	}
}

// ExampleInsertFilter is an example of how to use InsertFilter
func ExampleInsertFilter() {

	app := NewHttpServerWithCfg(newBConfig())
	app.Cfg.CopyRequestBody = true
	path := "/api/hello"
	app.Get(path, func(ctx *context.Context) {
		s := "hello world"
		fmt.Println(s)
		_ = ctx.Resp(s)
	})

	app.InsertFilter("*", BeforeStatic, func(ctx *context.Context) {
		fmt.Println("BeforeStatic filter process")
	})

	app.InsertFilter("*", BeforeRouter, func(ctx *context.Context) {
		fmt.Println("BeforeRouter filter process")
	})

	app.InsertFilter("*", BeforeExec, func(ctx *context.Context) {
		fmt.Println("BeforeExec filter process")
	})

	// need to set the WithReturnOnOutput false
	app.InsertFilter("*", AfterExec, func(ctx *context.Context) {
		fmt.Println("AfterExec filter process")
	}, WithReturnOnOutput(false))

	// need to set the WithReturnOnOutput false
	app.InsertFilter("*", FinishRouter, func(ctx *context.Context) {
		fmt.Println("FinishRouter filter process")
	}, WithReturnOnOutput(false))

	reader := strings.NewReader("")
	req := httptest.NewRequest("GET", path, reader)
	req.Header.Set("Accept", "*/*")

	w := httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, req)

	// Output:
	// BeforeStatic filter process
	// BeforeRouter filter process
	// BeforeExec filter process
	// hello world
	// AfterExec filter process
	// FinishRouter filter process
}
