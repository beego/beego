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

package httplib

import (
	"strings"
	"testing"
)

func TestSimpleGet(t *testing.T) {
	str, err := Get("http://httpbin.org/get").String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
}

func TestSimplePost(t *testing.T) {
	v := "smallfish"
	req := Post("http://httpbin.org/post")
	req.Param("username", v)
	str, err := req.String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in post")
	}
}

func TestPostFile(t *testing.T) {
	v := "smallfish"
	req := Post("http://httpbin.org/post")
	req.Param("username", v)
	req.PostFile("uploadfile", "httplib_test.go")
	str, err := req.String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in post")
	}
}

func TestWithCookie(t *testing.T) {
	v := "smallfish"
	str, err := Get("http://httpbin.org/cookies/set?k1=" + v).SetEnableCookie(true).String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
	str, err = Get("http://httpbin.org/cookies").SetEnableCookie(true).String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in cookie")
	}
}

func TestWithUserAgent(t *testing.T) {
	v := "beego"
	str, err := Get("http://httpbin.org/headers").SetUserAgent(v).String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in user-agent")
	}
}

func TestWithSetting(t *testing.T) {
	v := "beego"
	var setting BeegoHttpSettings
	setting.EnableCookie = true
	setting.UserAgent = v
	setting.Transport = nil
	SetDefaultSetting(setting)

	str, err := Get("http://httpbin.org/get").String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in user-agent")
	}
}
