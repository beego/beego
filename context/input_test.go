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

package context

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=astaxie", nil)
	beegoInput := NewInput()
	beegoInput.Context = NewContext()
	beegoInput.Context.Reset(httptest.NewRecorder(), r)
	beegoInput.ParseFormOrMulitForm(1 << 20)

	var id int
	err := beegoInput.Bind(&id, "id")
	if id != 123 || err != nil {
		t.Fatal("id should has int value")
	}

	var isok bool
	err = beegoInput.Bind(&isok, "isok")
	if !isok || err != nil {
		t.Fatal("isok should be true")
	}

	var float float64
	err = beegoInput.Bind(&float, "ft")
	if float != 1.2 || err != nil {
		t.Fatal("float should be equal to 1.2")
	}

	ol := make([]int, 0, 2)
	err = beegoInput.Bind(&ol, "ol")
	if len(ol) != 2 || err != nil || ol[0] != 1 || ol[1] != 2 {
		t.Fatal("ol should has two elements")
	}

	ul := make([]string, 0, 2)
	err = beegoInput.Bind(&ul, "ul")
	if len(ul) != 2 || err != nil || ul[0] != "str" || ul[1] != "array" {
		t.Fatal("ul should has two elements")
	}

	type User struct {
		Name string
	}
	user := User{}
	err = beegoInput.Bind(&user, "user")
	if err != nil || user.Name != "astaxie" {
		t.Fatal("user should has name")
	}
}

func TestParse2(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?user[0][Username]=Raph&user[1].Username=Leo&user[0].Password=123456&user[1][Password]=654321", nil)
	beegoInput := NewInput()
	beegoInput.Context = NewContext()
	beegoInput.Context.Reset(httptest.NewRecorder(), r)
	beegoInput.ParseFormOrMulitForm(1 << 20)
	type User struct {
		Username string
		Password string
	}
	var users []User
	err := beegoInput.Bind(&users, "user")
	fmt.Println(users)
	if err != nil || users[0].Username != "Raph" || users[0].Password != "123456" || users[1].Username != "Leo" || users[1].Password != "654321" {
		t.Fatal("users info wrong")
	}
}

func TestSubDomain(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://www.example.com/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=astaxie", nil)
	beegoInput := NewInput()
	beegoInput.Context = NewContext()
	beegoInput.Context.Reset(httptest.NewRecorder(), r)

	subdomain := beegoInput.SubDomains()
	if subdomain != "www" {
		t.Fatal("Subdomain parse error, got" + subdomain)
	}

	r, _ = http.NewRequest("GET", "http://localhost/", nil)
	beegoInput.Context.Request = r
	if beegoInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, should be empty, got " + beegoInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.example.com/", nil)
	beegoInput.Context.Request = r
	if beegoInput.SubDomains() != "aa.bb" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}

	/* TODO Fix this
	r, _ = http.NewRequest("GET", "http://127.0.0.1/", nil)
	beegoInput.Context.Request = r
	if beegoInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}
	*/

	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	beegoInput.Context.Request = r
	if beegoInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.cc.dd.example.com/", nil)
	beegoInput.Context.Request = r
	if beegoInput.SubDomains() != "aa.bb.cc.dd" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}
}

func TestParams(t *testing.T) {
	inp := NewInput()

	inp.SetParam("p1", "val1_ver1")
	inp.SetParam("p2", "val2_ver1")
	inp.SetParam("p3", "val3_ver1")
	if l := inp.ParamsLen(); l != 3 {
		t.Fatalf("Input.ParamsLen wrong value: %d, expected %d", l, 3)
	}

	if val := inp.Param("p1"); val != "val1_ver1" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val1_ver1")
	}
	if val := inp.Param("p3"); val != "val3_ver1" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val3_ver1")
	}
	vals := inp.Params()
	expected := map[string]string{
		"p1": "val1_ver1",
		"p2": "val2_ver1",
		"p3": "val3_ver1",
	}
	if !reflect.DeepEqual(vals, expected) {
		t.Fatalf("Input.Params wrong value: %s, expected %s", vals, expected)
	}

	// overwriting existing params
	inp.SetParam("p1", "val1_ver2")
	inp.SetParam("p2", "val2_ver2")
	expected = map[string]string{
		"p1": "val1_ver2",
		"p2": "val2_ver2",
		"p3": "val3_ver1",
	}
	vals = inp.Params()
	if !reflect.DeepEqual(vals, expected) {
		t.Fatalf("Input.Params wrong value: %s, expected %s", vals, expected)
	}

	if l := inp.ParamsLen(); l != 3 {
		t.Fatalf("Input.ParamsLen wrong value: %d, expected %d", l, 3)
	}

	if val := inp.Param("p1"); val != "val1_ver2" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val1_ver2")
	}

	if val := inp.Param("p2"); val != "val2_ver2" {
		t.Fatalf("Input.Param wrong value: %s, expected %s", val, "val1_ver2")
	}

}
