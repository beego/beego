// Copyright 2016 beego Author. All Rights Reserved.
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

package redis

import (
	"testing"
)

func TestSessionRelease(t *testing.T) {

	provider := Provider{}
	if err := provider.SessionInit(3, "127.0.0.1:6379"); err != nil {
		t.Fatal("init session err,", err)
	}

	sessionID := "beegosessionid_00001"

	session, err := provider.SessionRegenerate("", sessionID)
	if err != nil {
		t.Fatal("new session error,", err)
	}

	// set item.
	session.Set("k1", "v1")
	// update.
	session.SessionRelease(nil)

	session, err = provider.SessionRead(sessionID)
	if err != nil {
		t.Fatal("read session error,", err)
	}
	if v1 := session.Get("k1"); v1 == nil {
		t.Fatal("want v1 got nil")
	} else if v, _ := v1.(string); v != "v1" {
		t.Fatalf("want v1 got %s", v)
	}

	// delete
	provider.SessionDestroy(sessionID)
	session.Set("k2", "v2")

	session.SessionRelease(nil)

	session, err = provider.SessionRead(sessionID)
	if err != nil {
		t.Fatal("read session error,", err)
	}
	if session.Get("k1") != nil || session.Get("k2") != nil {
		t.Fatalf("want emtpy session value,got %s,%s", session.Get("k1"), session.Get("k2"))
	}

}
