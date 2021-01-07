// Copyright 2014 beego Author. All Rights Reserved.
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

package web

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

const exampleBody = "hello world"

type ExampleController struct {
	Controller
}

func (m ExampleController) Ping() {
	m.Ctx.Output.Body([]byte(exampleBody))
}

func (m ExampleController) ping() {
	m.Ctx.Output.Body([]byte(exampleBody))
}

func TestNamespaceGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Get("/user", func(ctx *context.Context) {
		ctx.Output.Body([]byte("v1_user"))
	})
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "v1_user" {
		t.Errorf("TestNamespaceGet can't run, get the response is " + w.Body.String())
	}
}

func TestNamespacePost(t *testing.T) {
	r, _ := http.NewRequest("POST", "/v1/user/123", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Post("/user/:id", func(ctx *context.Context) {
		ctx.Output.Body([]byte(ctx.Input.Param(":id")))
	})
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "123" {
		t.Errorf("TestNamespacePost can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNest(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/admin/order", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Namespace(
		NewNamespace("/admin").
			Get("/order", func(ctx *context.Context) {
				ctx.Output.Body([]byte("order"))
			}),
	)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "order" {
		t.Errorf("TestNamespaceNest can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNestParam(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/admin/order/123", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Namespace(
		NewNamespace("/admin").
			Get("/order/:id", func(ctx *context.Context) {
				ctx.Output.Body([]byte(ctx.Input.Param(":id")))
			}),
	)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "123" {
		t.Errorf("TestNamespaceNestParam can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouter(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/api/list", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Router("/api/list", &TestController{}, "*:List")
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("TestNamespaceRouter can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceAutoFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/test/list", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.AutoRouter(&TestController{})
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("user define func can't run")
	}
}

func TestNamespaceFilter(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/user/123", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Filter("before", func(ctx *context.Context) {
		ctx.Output.Body([]byte("this is Filter"))
	}).
		Get("/user/:id", func(ctx *context.Context) {
			ctx.Output.Body([]byte(ctx.Input.Param(":id")))
		})
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "this is Filter" {
		t.Errorf("TestNamespaceFilter can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCond(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v2/test/list", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v2")
	ns.Cond(func(ctx *context.Context) bool {
		return ctx.Input.Domain() == "beego.me"
	}).
		AutoRouter(&TestController{})
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != 405 {
		t.Errorf("TestNamespaceCond can't run get the result " + strconv.Itoa(w.Code))
	}
}

func TestNamespaceInside(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/shop/order/123", nil)
	w := httptest.NewRecorder()
	ns := NewNamespace("/v3",
		NSAutoRouter(&TestController{}),
		NSNamespace("/shop",
			NSGet("/order/:id", func(ctx *context.Context) {
				ctx.Output.Body([]byte(ctx.Input.Param(":id")))
			}),
		),
	)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "123" {
		t.Errorf("TestNamespaceInside can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterGet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.RouterGet("/user", ExampleController.Ping)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterGet can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterPost(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.RouterPost("/user", ExampleController.Ping)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterPost can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterDelete(t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.RouterDelete("/user", ExampleController.Ping)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterDelete can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterPut(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPut, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.RouterPut("/user", ExampleController.Ping)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterPut can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterHead(t *testing.T) {
	r, _ := http.NewRequest(http.MethodHead, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.RouterHead("/user", ExampleController.Ping)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterHead can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterOptions(t *testing.T) {
	r, _ := http.NewRequest(http.MethodOptions, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.RouterOptions("/user", ExampleController.Ping)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterOptions can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterPatch(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPatch, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.RouterPatch("/user", ExampleController.Ping)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterPatch can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouterAny(t *testing.T) {
	ns := NewNamespace("/v1")
	ns.RouterAny("/user", ExampleController.Ping)
	AddNamespace(ns)

	for method, _ := range HTTPMETHOD {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(method, "/v1/user", nil)
		BeeApp.Handlers.ServeHTTP(w, r)
		if w.Body.String() != exampleBody {
			t.Errorf("TestNamespaceRouterAny can't run, get the response is " + w.Body.String())
		}
	}
}

func TestNamespaceNSRouterGet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	NSRouterGet("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterGet can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSRouterPost(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	NSRouterPost("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterPost can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSRouterDelete(t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	NSRouterDelete("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterDelete can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSRouterPut(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPut, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	NSRouterPut("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterPut can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSRouterHead(t *testing.T) {
	r, _ := http.NewRequest(http.MethodHead, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	NSRouterHead("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterHead can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSRouterOptions(t *testing.T) {
	r, _ := http.NewRequest(http.MethodOptions, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	NSRouterOptions("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterOptions can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSRouterPatch(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPatch, "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	NSRouterPatch("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceRouterPatch can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSRouterAny(t *testing.T) {
	ns := NewNamespace("/v1")
	NSRouterAny("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)

	for method, _ := range HTTPMETHOD {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(method, "/v1/user", nil)
		BeeApp.Handlers.ServeHTTP(w, r)
		if w.Body.String() != exampleBody {
			t.Errorf("TestNamespaceRouterAny can't run, get the response is " + w.Body.String())
		}
	}
}
