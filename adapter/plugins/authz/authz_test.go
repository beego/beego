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

package authz

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin"

	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/context"
	"github.com/beego/beego/v2/adapter/plugins/auth"
)

const (
	authCfg = "authz_model.conf"
	authCsv = "authz_policy.csv"
)

func testRequest(t *testing.T, handler *beego.ControllerRegister, user string, path string, method string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	r.SetBasicAuth(user, "123")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", user, path, method, w.Code, code)
	}
}

func TestBasic(t *testing.T) {
	handler := beego.NewControllerRegister()

	_ = handler.InsertFilter("*", beego.BeforeRouter, auth.Basic("alice", "123"))

	_ = handler.InsertFilter("*", beego.BeforeRouter, NewAuthorizer(casbin.NewEnforcer(authCfg, authCsv)))

	handler.Any("*", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
	})

	const d1r1 = "/dataset1/resource1"
	testRequest(t, handler, "alice", d1r1, "GET", 200)
	testRequest(t, handler, "alice", d1r1, "POST", 200)
	const d1r2 = "/dataset1/resource2"
	testRequest(t, handler, "alice", d1r2, "GET", 200)
	testRequest(t, handler, "alice", d1r2, "POST", 403)
}

func TestPathWildcard(t *testing.T) {
	handler := beego.NewControllerRegister()

	_ = handler.InsertFilter("*", beego.BeforeRouter, auth.Basic("bob", "123"))
	_ = handler.InsertFilter("*", beego.BeforeRouter, NewAuthorizer(casbin.NewEnforcer(authCfg, authCsv)))

	handler.Any("*", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
	})

	const d2r1 = "/dataset2/resource1"
	testRequest(t, handler, "bob", d2r1, "GET", 200)
	testRequest(t, handler, "bob", d2r1, "POST", 200)
	testRequest(t, handler, "bob", d2r1, "DELETE", 200)
	const d2r2 = "/dataset2/resource2"
	testRequest(t, handler, "bob", d2r2, "GET", 200)
	testRequest(t, handler, "bob", d2r2, "POST", 403)
	testRequest(t, handler, "bob", d2r2, "DELETE", 403)

	const item1 = "/dataset2/folder1/item1"
	testRequest(t, handler, "bob", item1, "GET", 403)
	testRequest(t, handler, "bob", item1, "POST", 200)
	testRequest(t, handler, "bob", item1, "DELETE", 403)
	const item2 = "/dataset2/folder1/item2"
	testRequest(t, handler, "bob", item2, "GET", 403)
	testRequest(t, handler, "bob", item2, "POST", 200)
	testRequest(t, handler, "bob", item2, "DELETE", 403)
}

func TestRBAC(t *testing.T) {
	handler := beego.NewControllerRegister()

	_ = handler.InsertFilter("*", beego.BeforeRouter, auth.Basic("cathy", "123"))
	e := casbin.NewEnforcer(authCfg, authCsv)
	_ = handler.InsertFilter("*", beego.BeforeRouter, NewAuthorizer(e))

	handler.Any("*", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
	})

	// cathy can access all /dataset1/* resources via all methods because it has the dataset1_admin role.
	const dataSet1 = "/dataset1/item"
	testRequest(t, handler, "cathy", dataSet1, "GET", 200)
	testRequest(t, handler, "cathy", dataSet1, "POST", 200)
	testRequest(t, handler, "cathy", dataSet1, "DELETE", 200)
	const dataSet2 = "/dataset2/item"
	testRequest(t, handler, "cathy", dataSet2, "GET", 403)
	testRequest(t, handler, "cathy", dataSet2, "POST", 403)
	testRequest(t, handler, "cathy", dataSet2, "DELETE", 403)

	// delete all roles on user cathy, so cathy cannot access any resources now.
	e.DeleteRolesForUser("cathy")

	testRequest(t, handler, "cathy", dataSet1, "GET", 403)
	testRequest(t, handler, "cathy", dataSet1, "POST", 403)
	testRequest(t, handler, "cathy", dataSet1, "DELETE", 403)
	testRequest(t, handler, "cathy", dataSet2, "GET", 403)
	testRequest(t, handler, "cathy", dataSet2, "POST", 403)
	testRequest(t, handler, "cathy", dataSet2, "DELETE", 403)
}
