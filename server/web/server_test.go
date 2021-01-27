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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpServerWithCfg(t *testing.T) {

	BConfig.AppName = "Before"
	svr := NewHttpServerWithCfg(BConfig)
	svr.Cfg.AppName = "hello"
	assert.Equal(t, "hello", BConfig.AppName)

}

func TestServerRouterGet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	RouterGet("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerRouterGet can't run")
	}
}

func TestServerRouterPost(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/user", nil)
	w := httptest.NewRecorder()

	RouterPost("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerRouterPost can't run")
	}
}

func TestServerRouterHead(t *testing.T) {
	r, _ := http.NewRequest(http.MethodHead, "/user", nil)
	w := httptest.NewRecorder()

	RouterHead("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerRouterHead can't run")
	}
}

func TestServerRouterPut(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPut, "/user", nil)
	w := httptest.NewRecorder()

	RouterPut("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerRouterPut can't run")
	}
}

func TestServerRouterPatch(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPatch, "/user", nil)
	w := httptest.NewRecorder()

	RouterPatch("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerRouterPatch can't run")
	}
}

func TestServerRouterDelete(t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, "/user", nil)
	w := httptest.NewRecorder()

	RouterDelete("/user", ExampleController.Ping)
	BeeApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerRouterDelete can't run")
	}
}

func TestServerRouterAny(t *testing.T) {
	RouterAny("/user", ExampleController.Ping)

	for method := range HTTPMETHOD {
		r, _ := http.NewRequest(method, "/user", nil)
		w := httptest.NewRecorder()
		BeeApp.Handlers.ServeHTTP(w, r)
		if w.Body.String() != exampleBody {
			t.Errorf("TestServerRouterAny can't run")
		}
	}
}
