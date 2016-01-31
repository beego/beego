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

func TestCookie(t *testing.T) {
	config := `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`
	globalSessions, err := NewManager("cookie", config)
	if err != nil {
		t.Fatal("init cookie session err", err)
	}
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
	sess.SessionRelease(w)
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
}

func TestDestorySessionCookie(t *testing.T) {
	config := `{"cookieName":"gosessionid","enableSetCookie":true,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`
	globalSessions, err := NewManager("cookie", config)
	if err != nil {
		t.Fatal("init cookie session err", err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	session, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start err,", err)
	}

	// request again ,will get same sesssion id .
	r1, _ := http.NewRequest("GET", "/", nil)
	r1.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	newSession, err := globalSessions.SessionStart(w, r1)
	if err != nil {
		t.Fatal("session start err,", err)
	}
	if newSession.SessionID() != session.SessionID() {
		t.Fatal("get cookie session id is not the same again.")
	}

	// After destroy session , will get a new session id .
	globalSessions.SessionDestroy(w, r1)
	r2, _ := http.NewRequest("GET", "/", nil)
	r2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	w = httptest.NewRecorder()
	newSession, err = globalSessions.SessionStart(w, r2)
	if err != nil {
		t.Fatal("session start error")
	}
	if newSession.SessionID() == session.SessionID() {
		t.Fatal("after destroy session and reqeust again ,get cookie session id is same.")
	}
}
