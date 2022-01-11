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
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	sid         = "Session_id"
	sidNew      = "Session_id_new"
	sessionPath = "./_session_runtime"
)

var mutex sync.Mutex

func TestFileProviderSessionExist(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	if fp.SessionExist(sid) {
		t.Error()
	}

	_, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	if !fp.SessionExist(sid) {
		t.Error()
	}
}

func TestFileProviderSessionExist2(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	if fp.SessionExist(sid) {
		t.Error()
	}

	if fp.SessionExist("") {
		t.Error()
	}

	if fp.SessionExist("1") {
		t.Error()
	}
}

func TestFileProviderSessionRead(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	s, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	_ = s.Set("sessionValue", 18975)
	v := s.Get("sessionValue")

	if v.(int) != 18975 {
		t.Error()
	}
}

func TestFileProviderSessionRead1(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	_, err := fp.SessionRead("")
	if err == nil {
		t.Error(err)
	}

	_, err = fp.SessionRead("1")
	if err == nil {
		t.Error(err)
	}
}

func TestFileProviderSessionAll(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 546

	for i := 1; i <= sessionCount; i++ {
		_, err := fp.SessionRead(fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
	}

	if fp.SessionAll() != sessionCount {
		t.Error()
	}
}

func TestFileProviderSessionRegenerate(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	_, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	if !fp.SessionExist(sid) {
		t.Error()
	}

	_, err = fp.SessionRegenerate(sid, sidNew)
	if err != nil {
		t.Error(err)
	}

	if fp.SessionExist(sid) {
		t.Error()
	}

	if !fp.SessionExist(sidNew) {
		t.Error()
	}
}

func TestFileProviderSessionDestroy(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	_, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	if !fp.SessionExist(sid) {
		t.Error()
	}

	err = fp.SessionDestroy(sid)
	if err != nil {
		t.Error(err)
	}

	if fp.SessionExist(sid) {
		t.Error()
	}
}

func TestFileProviderSessionGC(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(1, sessionPath)

	sessionCount := 412

	for i := 1; i <= sessionCount; i++ {
		_, err := fp.SessionRead(fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
	}

	time.Sleep(2 * time.Second)

	fp.SessionGC()
	if fp.SessionAll() != 0 {
		t.Error()
	}
}

func TestFileSessionStoreSet(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(sid)
	for i := 1; i <= sessionCount; i++ {
		err := s.Set(i, i)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFileSessionStoreGet(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(sid)
	for i := 1; i <= sessionCount; i++ {
		_ = s.Set(i, i)

		v := s.Get(i)
		if v.(int) != i {
			t.Error()
		}
	}
}

func TestFileSessionStoreDelete(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	s, _ := fp.SessionRead(sid)
	s.Set("1", 1)

	if s.Get("1") == nil {
		t.Error()
	}

	s.Delete("1")

	if s.Get("1") != nil {
		t.Error()
	}
}

func TestFileSessionStoreFlush(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(sid)
	for i := 1; i <= sessionCount; i++ {
		_ = s.Set(i, i)
	}

	_ = s.Flush()

	for i := 1; i <= sessionCount; i++ {
		if s.Get(i) != nil {
			t.Error()
		}
	}
}

func TestFileSessionStoreSessionID(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 85

	for i := 1; i <= sessionCount; i++ {
		s, err := fp.SessionRead(fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
		if s.SessionID() != fmt.Sprintf("%s_%d", sid, i) {
			t.Error(err)
		}
	}
}
