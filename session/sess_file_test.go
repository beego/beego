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
	"testing"
	"time"
)

func TestFileSession(t *testing.T) {
	conf := new(ManagerConfig)
	conf.CookieName = "gosessionid"
	conf.EnableSetCookie = true
	conf.Gclifetime = 1
	conf.Maxlifetime = 1
	conf.ProviderConfig = "./tmp"
	globalSessions, _ := NewManager("file", conf)
	go globalSessions.GC()
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start error,", err)
	}
	err = sess.Set("username", "astaxie")
	if err != nil {
		t.Fatal("set error,", err)
	}
	if v := sess.Get("username"); v == nil || v.(string) != "astaxie" {
		t.Fatal("get username error")
	}
	err = sess.Delete("username")
	if err != nil {
		t.Fatal("delete error,", err)
	}
	if v := sess.Get("username"); v != nil {
		t.Fatal("get deleted session not return nil")
	}
	err = sess.Set("username", "astaxie")
	if err != nil {
		t.Fatal("set error,", err)
	}
	err = sess.Flush()
	if err != nil {
		t.Fatal("flush error,", err)
	}
	if v := sess.Get("username"); v != nil {
		t.Fatal("not really flush session")
	}
	err = sess.Set("username", "astaxie")
	if err != nil {
		t.Fatal("set error,", err)
	}
	sess.SessionRelease(w)
	t.Log("sleep 3 seconds for session expired")
	time.Sleep(3 * time.Second)
	t.Log("test session expiring")
	sess, err = globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start error,", err)
	}
	if v := sess.Get("username"); v != nil {
		t.Fatal("session hasn't expired")
	}
	sess.SessionRelease(w)
}
