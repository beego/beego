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
	"fmt"
	"net/http"
	"testing"
)

func TestParse(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=astaxie", nil)
	beegoInput := NewInput(r)
	beegoInput.ParseFormOrMulitForm(1 << 20)

	var id int
	err := beegoInput.Bind(&id, "id")
	if id != 123 || err != nil {
		t.Fatal("id should has int value")
	}
	fmt.Println(id)

	var isok bool
	err = beegoInput.Bind(&isok, "isok")
	if !isok || err != nil {
		t.Fatal("isok should be true")
	}
	fmt.Println(isok)

	var float float64
	err = beegoInput.Bind(&float, "ft")
	if float != 1.2 || err != nil {
		t.Fatal("float should be equal to 1.2")
	}
	fmt.Println(float)

	ol := make([]int, 0, 2)
	err = beegoInput.Bind(&ol, "ol")
	if len(ol) != 2 || err != nil || ol[0] != 1 || ol[1] != 2 {
		t.Fatal("ol should has two elements")
	}
	fmt.Println(ol)

	ul := make([]string, 0, 2)
	err = beegoInput.Bind(&ul, "ul")
	if len(ul) != 2 || err != nil || ul[0] != "str" || ul[1] != "array" {
		t.Fatal("ul should has two elements")
	}
	fmt.Println(ul)

	type User struct {
		Name string
	}
	user := User{}
	err = beegoInput.Bind(&user, "user")
	if err != nil || user.Name != "astaxie" {
		t.Fatal("user should has name")
	}
	fmt.Println(user)
}

func TestSubDomain(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://www.example.com/?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=astaxie", nil)
	beegoInput := NewInput(r)

	subdomain := beegoInput.SubDomains()
	if subdomain != "www" {
		t.Fatal("Subdomain parse error, got" + subdomain)
	}

	r, _ = http.NewRequest("GET", "http://localhost/", nil)
	beegoInput.Request = r
	if beegoInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, should be empty, got " + beegoInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.example.com/", nil)
	beegoInput.Request = r
	if beegoInput.SubDomains() != "aa.bb" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}

	/* TODO Fix this
	r, _ = http.NewRequest("GET", "http://127.0.0.1/", nil)
	beegoInput.Request = r
	if beegoInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}
	*/

	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	beegoInput.Request = r
	if beegoInput.SubDomains() != "" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}

	r, _ = http.NewRequest("GET", "http://aa.bb.cc.dd.example.com/", nil)
	beegoInput.Request = r
	if beegoInput.SubDomains() != "aa.bb.cc.dd" {
		t.Fatal("Subdomain parse error, got " + beegoInput.SubDomains())
	}
}
