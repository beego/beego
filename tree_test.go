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
	routers = append(routers, testinfo{"/", "/", nil})
	routers = append(routers, testinfo{"/customer/login", "/customer/login", nil})
	routers = append(routers, testinfo{"/*", "/customer/123", map[string]string{":splat": "customer/123"}})
	routers = append(routers, testinfo{"/*.*", "/nice/api.json", map[string]string{":path": "nice/api", ":ext": "json"}})
	routers = append(routers, testinfo{"/v1/shop/:id:int", "/v1/shop/123", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/v1/shop/:id/:name", "/v1/shop/123/nike", map[string]string{":id": "123", ":name": "nike"}})
	routers = append(routers, testinfo{"/v1/shop/:id/account", "/v1/shop/123/account", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/v1/shop/:name:string", "/v1/shop/nike", map[string]string{":name": "nike"}})
	routers = append(routers, testinfo{"/v1/shop/:id([0-9]+)", "/v1/shop//123", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/v1/shop/:id([0-9]+)_:name", "/v1/shop/123_nike", map[string]string{":id": "123", ":name": "nike"}})
	routers = append(routers, testinfo{"/v1/shop/:id_cms.html", "/v1/shop/123_cms.html", map[string]string{":id": "123"}})
	routers = append(routers, testinfo{"/v1/shop/cms_:id_:page.html", "/v1/shop/cms_123_1.html", map[string]string{":id": "123", ":page": "1"}})
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
					t.Fatal(r.url + r.requesturl + " get param empty:" + k)
				} else if vv != v {
					t.Fatal(r.url + " " + r.requesturl + " should be:" + v + " get param:" + vv)
				}
			}
		}
	}
}

func TestSplitPath(t *testing.T) {
	a := splitPath("/")
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
	b, w, r = splitSegment(":id_cms.html")
	if !b || len(w) != 1 || w[0] != ":id" || r != `(.+)_cms.html` {
		t.Fatal(":id_cms.html should return true, [:id], '(.+)_cms.html'")
	}
	b, w, r = splitSegment("cms_:id_:page.html")
	if !b || len(w) != 2 || w[0] != ":id" || w[1] != ":page" || r != `cms_(.+)_(.+).html` {
		t.Fatal(":id_cms.html should return true, [:id :page], cms_(.+)_(.+).html")
	}
}
