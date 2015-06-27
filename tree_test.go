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

package beego

import "testing"

type testinfo struct {
	url        string
	requesturl string
	params     map[string]string
}

var routers []testinfo

func init() {
	routers = make([]testinfo, 0)
	routers = append(routers, testinfo{"/:id", "/123", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/hello/?:id", "/hello", map[string]string{":id": ""}})
	routers = append(routers, testinfo{"/", "/", nil})
	routers = append(routers, testinfo{"/customer/login", "/customer/login", nil})
	routers = append(routers, testinfo{"/customer/login", "/customer/login.json", map[string]string{":ext": "json"}})
	routers = append(routers, testinfo{"/*", "/customer/123", map[string]string{":splat": "customer/123"}})
	routers = append(routers, testinfo{"/*", "/customer/2009/12/11", map[string]string{":splat": "customer/2009/12/11"}})
	routers = append(routers, testinfo{"/aa/*/bb", "/aa/2009/bb", map[string]string{":splat": "2009"}})
	routers = append(routers, testinfo{"/cc/*/dd", "/cc/2009/11/dd", map[string]string{":splat": "2009/11"}})
	routers = append(routers, testinfo{"/ee/:year/*/ff", "/ee/2009/11/ff", map[string]string{":year": "2009", ":splat": "11"}})
	routers = append(routers, testinfo{"/thumbnail/:size/uploads/*",
		"/thumbnail/100x100/uploads/items/2014/04/20/dPRCdChkUd651t1Hvs18.jpg",
		map[string]string{":size": "100x100", ":splat": "items/2014/04/20/dPRCdChkUd651t1Hvs18.jpg"}})
	routers = append(routers, testinfo{"/*.*", "/nice/api.json", map[string]string{":path": "nice/api", ":ext": "json"}})
	routers = append(routers, testinfo{"/:name/*.*", "/nice/api.json", map[string]string{":name": "nice", ":path": "api", ":ext": "json"}})
	routers = append(routers, testinfo{"/:name/test/*.*", "/nice/test/api.json", map[string]string{":name": "nice", ":path": "api", ":ext": "json"}})
	routers = append(routers, testinfo{"/dl/:width:int/:height:int/*.*",
		"/dl/48/48/05ac66d9bda00a3acf948c43e306fc9a.jpg",
		map[string]string{":width": "48", ":height": "48", ":ext": "jpg", ":path": "05ac66d9bda00a3acf948c43e306fc9a"}})
	routers = append(routers, testinfo{"/v1/shop/:id:int", "/v1/shop/123", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/:year:int/:month:int/:id/:endid", "/1111/111/aaa/aaa", map[string]string{":year": "1111", ":month": "111", ":id": "aaa", ":endid": "aaa"}})
	routers = append(routers, testinfo{"/v1/shop/:id/:name", "/v1/shop/123/nike", map[string]string{":id": "123", ":name": "nike"}})
	routers = append(routers, testinfo{"/v1/shop/:id/account", "/v1/shop/123/account", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/v1/shop/:name:string", "/v1/shop/nike", map[string]string{":name": "nike"}})
	routers = append(routers, testinfo{"/v1/shop/:id([0-9]+)", "/v1/shop//123", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/v1/shop/:id([0-9]+)_:name", "/v1/shop/123_nike", map[string]string{":id": "123", ":name": "nike"}})
	routers = append(routers, testinfo{"/v1/shop/:id(.+)_cms.html", "/v1/shop/123_cms.html", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/v1/shop/cms_:id(.+)_:page(.+).html", "/v1/shop/cms_123_1.html", map[string]string{":id": "123", ":page": "1"}})
	routers = append(routers, testinfo{"/v1/:v/cms/aaa_:id(.+)_:page(.+).html", "/v1/2/cms/aaa_123_1.html", map[string]string{":v": "2", ":id": "123", ":page": "1"}})
	routers = append(routers, testinfo{"/v1/:v/cms_:id(.+)_:page(.+).html", "/v1/2/cms_123_1.html", map[string]string{":v": "2", ":id": "123", ":page": "1"}})
	routers = append(routers, testinfo{"/v1/:v(.+)_cms/ttt_:id(.+)_:page(.+).html", "/v1/2_cms/ttt_123_1.html", map[string]string{":v": "2", ":id": "123", ":page": "1"}})
}

func TestTreeRouters(t *testing.T) {
	for _, r := range routers {
		tr := NewTree()
		tr.AddRouter(r.url, "astaxie")
		obj, param := tr.Match(r.requesturl)
		if obj == nil || obj.(string) != "astaxie" {
			t.Fatal(r.url + " can't get obj ")
		}
		if r.params != nil {
			for k, v := range r.params {
				if vv, ok := param[k]; !ok {
					t.Fatal(r.url + "    " + r.requesturl + " get param empty:" + k)
				} else if vv != v {
					t.Fatal(r.url + "     " + r.requesturl + " should be:" + v + " get param:" + vv)
				}
			}
		}
	}
}

func TestAddTree(t *testing.T) {
	tr := NewTree()
	tr.AddRouter("/shop/:id/account", "astaxie")
	tr.AddRouter("/shop/:sd/ttt_:id(.+)_:page(.+).html", "astaxie")
	t1 := NewTree()
	t1.AddTree("/v1/zl", tr)
	obj, param := t1.Match("/v1/zl/shop/123/account")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/v1/zl/shop/:id/account can't get obj ")
	}
	if param == nil {
		t.Fatal("get param error")
	}
	if param[":id"] != "123" {
		t.Fatal("get :id param error")
	}
	obj, param = t1.Match("/v1/zl/shop/123/ttt_1_12.html")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/v1/zl//shop/:sd/ttt_:id(.+)_:page(.+).html can't get obj ")
	}
	if param == nil {
		t.Fatal("get param error")
	}
	if param[":sd"] != "123" || param[":id"] != "1" || param[":page"] != "12" {
		t.Fatal("get :sd :id :page param error")
	}

	t2 := NewTree()
	t2.AddTree("/v1/:shopid", tr)
	obj, param = t2.Match("/v1/zl/shop/123/account")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/v1/:shopid/shop/:id/account can't get obj ")
	}
	if param == nil {
		t.Fatal("get param error")
	}
	if param[":id"] != "123" || param[":shopid"] != "zl" {
		t.Fatal("get :id :shopid param error")
	}
	obj, param = t2.Match("/v1/zl/shop/123/ttt_1_12.html")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/v1/:shopid/shop/:sd/ttt_:id(.+)_:page(.+).html can't get obj ")
	}
	if param == nil {
		t.Fatal("get :shopid param error")
	}
	if param[":sd"] != "123" || param[":id"] != "1" || param[":page"] != "12" || param[":shopid"] != "zl" {
		t.Fatal("get :sd :id :page :shopid param error")
	}
}

func TestAddTree2(t *testing.T) {
	tr := NewTree()
	tr.AddRouter("/shop/:id/account", "astaxie")
	tr.AddRouter("/shop/:sd/ttt_:id(.+)_:page(.+).html", "astaxie")
	t3 := NewTree()
	t3.AddTree("/:version(v1|v2)/:prefix", tr)
	obj, param := t3.Match("/v1/zl/shop/123/account")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/:version(v1|v2)/:prefix/shop/:id/account can't get obj ")
	}
	if param == nil {
		t.Fatal("get param error")
	}
	if param[":id"] != "123" || param[":prefix"] != "zl" || param[":version"] != "v1" {
		t.Fatal("get :id :prefix :version param error")
	}
}

func TestAddTree3(t *testing.T) {
	tr := NewTree()
	tr.AddRouter("/create", "astaxie")
	tr.AddRouter("/shop/:sd/account", "astaxie")
	t3 := NewTree()
	t3.AddTree("/table/:num", tr)
	obj, param := t3.Match("/table/123/shop/123/account")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/table/:num/shop/:sd/account can't get obj ")
	}
	if param == nil {
		t.Fatal("get param error")
	}
	if param[":num"] != "123" || param[":sd"] != "123" {
		t.Fatal("get :num :sd param error")
	}
	obj, param = t3.Match("/table/123/create")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/table/:num/create can't get obj ")
	}
}

func TestAddTree4(t *testing.T) {
	tr := NewTree()
	tr.AddRouter("/create", "astaxie")
	tr.AddRouter("/shop/:sd/:account", "astaxie")
	t4 := NewTree()
	t4.AddTree("/:info:int/:num/:id", tr)
	obj, param := t4.Match("/12/123/456/shop/123/account")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/:info:int/:num/:id/shop/:sd/:account can't get obj ")
	}
	if param == nil {
		t.Fatal("get param error")
	}
	if param[":info"] != "12" || param[":num"] != "123" || param[":id"] != "456" || param[":sd"] != "123" || param[":account"] != "account" {
		t.Fatal("get :info :num :id :sd :account param error")
	}
	obj, param = t4.Match("/12/123/456/create")
	if obj == nil || obj.(string) != "astaxie" {
		t.Fatal("/:info:int/:num/:id/create can't get obj ")
	}
}

func TestSplitPath(t *testing.T) {
	a := splitPath("")
	if len(a) != 0 {
		t.Fatal("/ should retrun []")
	}
	a = splitPath("/")
	if len(a) != 0 {
		t.Fatal("/ should retrun []")
	}
	a = splitPath("/admin")
	if len(a) != 1 || a[0] != "admin" {
		t.Fatal("/admin should retrun [admin]")
	}
	a = splitPath("/admin/")
	if len(a) != 1 || a[0] != "admin" {
		t.Fatal("/admin/ should retrun [admin]")
	}
	a = splitPath("/admin/users")
	if len(a) != 2 || a[0] != "admin" || a[1] != "users" {
		t.Fatal("/admin should retrun [admin users]")
	}
	a = splitPath("/admin/:id:int")
	if len(a) != 2 || a[0] != "admin" || a[1] != ":id:int" {
		t.Fatal("/admin should retrun [admin :id:int]")
	}
}

func TestSplitSegment(t *testing.T) {
	b, w, r := splitSegment("admin")
	if b || len(w) != 0 || r != "" {
		t.Fatal("admin should return false, nil, ''")
	}
	b, w, r = splitSegment("*")
	if !b || len(w) != 1 || w[0] != ":splat" || r != "" {
		t.Fatal("* should return true, [:splat], ''")
	}
	b, w, r = splitSegment("*.*")
	if !b || len(w) != 3 || w[1] != ":path" || w[2] != ":ext" || w[0] != "." || r != "" {
		t.Fatal("admin should return true,[. :path :ext], ''")
	}
	b, w, r = splitSegment(":id")
	if !b || len(w) != 1 || w[0] != ":id" || r != "" {
		t.Fatal(":id should return true, [:id], ''")
	}
	b, w, r = splitSegment("?:id")
	if !b || len(w) != 2 || w[0] != ":" || w[1] != ":id" || r != "" {
		t.Fatal("?:id should return true, [: :id], ''")
	}
	b, w, r = splitSegment(":id:int")
	if !b || len(w) != 1 || w[0] != ":id" || r != "([0-9]+)" {
		t.Fatal(":id:int should return true, [:id], '([0-9]+)'")
	}
	b, w, r = splitSegment(":name:string")
	if !b || len(w) != 1 || w[0] != ":name" || r != `([\w]+)` {
		t.Fatal(`:name:string should return true, [:name], '([\w]+)'`)
	}
	b, w, r = splitSegment(":id([0-9]+)")
	if !b || len(w) != 1 || w[0] != ":id" || r != `([0-9]+)` {
		t.Fatal(`:id([0-9]+) should return true, [:id], '([0-9]+)'`)
	}
	b, w, r = splitSegment(":id([0-9]+)_:name")
	if !b || len(w) != 2 || w[0] != ":id" || w[1] != ":name" || r != `([0-9]+)_(.+)` {
		t.Fatal(`:id([0-9]+)_:name should return true, [:id :name], '([0-9]+)_(.+)'`)
	}
	b, w, r = splitSegment(":id(.+)_cms.html")
	if !b || len(w) != 1 || w[0] != ":id" || r != `(.+)_cms.html` {
		t.Fatal(":id_cms.html should return true, [:id], '(.+)_cms.html'")
	}
	b, w, r = splitSegment("cms_:id(.+)_:page(.+).html")
	if !b || len(w) != 2 || w[0] != ":id" || w[1] != ":page" || r != `cms_(.+)_(.+).html` {
		t.Fatal(":id_cms.html should return true, [:id :page], cms_(.+)_(.+).html")
	}
}
