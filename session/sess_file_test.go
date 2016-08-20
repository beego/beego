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

package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	conf := &ManagerConfig{
		CookieName:      "gosessionid",
		Gclifetime:      180,
		EnableSetCookie: true,
		ProviderConfig:  "/tmp",
	}

	globalSessions, _ := NewManager("file", conf)
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("set error,", err)
	}
	err = sess.Set("username", "astaxie")
	if err != nil {
		t.Fatal("set error,", err)
	}
	if username := sess.Get("username"); username != "astaxie" {
		t.Fatal("get username error")
	}
	if cookiestr := w.Header().Get("Set-Cookie"); cookiestr == "" {
		t.Fatal("setcookie error")
	} else {
		parts := strings.Split(strings.TrimSpace(cookiestr), ";")
		for k, v := range parts {
			nameval := strings.Split(v, "=")
			if k == 0 && nameval[0] != "gosessionid" {
				t.Fatal("error")
			}
		}
	}
	sess.SessionRelease(w)
	r2, _ := http.NewRequest("GET", "/", nil)
	w2 := httptest.NewRecorder()
	r3, _ := http.NewRequest("GET", "/", nil)
	w3 := httptest.NewRecorder()
	r2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	r3.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	sess2, err2 := globalSessions.SessionStart(w2, r2)
	if err2 != nil {
		t.Fatal("set error,", err2)
	}
	sess3, err3 := globalSessions.SessionStart(w3, r3)
	if err3 != nil {
		t.Fatal("set error,", err3)
	}
	err2 = sess2.Set("sex", "man")
	if err2 != nil {
		t.Fatal("set error,", err2)
	}
	err3 = sess3.Set("coutry", "China")
	if err3 != nil {
		t.Fatal("set error,", err3)
	}

	if username := sess2.Get("username"); username != "astaxie" {
		t.Fatal("get username error")
	}

	if sex := sess2.Get("sex"); sex != "man" {
		t.Fatal("get sex error")
	}

	if coutry := sess2.Get("coutry"); coutry != "China" {
		t.Fatal("get coutry error")
	}

	sess2.SessionRelease(w2)

	if username := sess3.Get("username"); username != "astaxie" {
		t.Fatal("get username error")
	}

	if sex := sess3.Get("sex"); sex != "man" {
		t.Fatal("get sex error")
	}

	if coutry := sess3.Get("coutry"); coutry != "China" {
		t.Fatal("get coutry error")
	}

	sess3.SessionRelease(w3)

	r4, _ := http.NewRequest("GET", "/", nil)
	w4 := httptest.NewRecorder()
	r4.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	sess4, err4 := globalSessions.SessionStart(w4, r4)
	if err4 != nil {
		t.Fatal("set error,", err4)
	}

	if username := sess4.Get("username"); username != "astaxie" {
		t.Fatal("get username error")
	}

	if sex := sess4.Get("sex"); sex != "man" {
		t.Fatal("get sex error")
	}

	if coutry := sess4.Get("coutry"); coutry != "China" {
		t.Fatal("get coutry error")
	}

}
