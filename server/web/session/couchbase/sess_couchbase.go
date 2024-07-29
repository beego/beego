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

// Package couchbase for session provider
//
// depend on github.com/couchbaselabs/go-couchbasee
//
// go install github.com/couchbaselabs/go-couchbase
//
// Usage:
// import(
//
//	_ "github.com/beego/beego/v2/server/web/session/couchbase"
//	"github.com/beego/beego/v2/server/web/session"
//
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("couchbase", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"http://host:port/, Pool, Bucket"}``)
//		go globalSessions.GC()
//	}
package couchbase

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	couchbase "github.com/couchbase/go-couchbase"

	"github.com/beego/beego/v2/server/web/session"
)

var couchbpder = &Provider{}

// SessionStore store each session
type SessionStore struct {
	b           *couchbase.Bucket
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Provider couchabse provided
type Provider struct {
	maxlifetime int64
	SavePath    string `json:"save_path"`
	Pool        string `json:"pool"`
	Bucket      string `json:"bucket"`
	b           *couchbase.Bucket
}

// Set value to couchabse session
func (cs *SessionStore) Set(ctx context.Context, key, value interface{}) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.values[key] = value
	return nil
}

// Get value from couchabse session
func (cs *SessionStore) Get(ctx context.Context, key interface{}) interface{} {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	if v, ok := cs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in couchbase session by given key
func (cs *SessionStore) Delete(ctx context.Context, key interface{}) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	delete(cs.values, key)
	return nil
}

// Flush Clean all values in couchbase session
func (cs *SessionStore) Flush(context.Context) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID Get couchbase session store id
func (cs *SessionStore) SessionID(context.Context) string {
	return cs.sid
}

// SessionRelease Write couchbase session with Gob string
func (cs *SessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	defer cs.b.Close()
	cs.lock.RLock()
	values := cs.values
	cs.lock.RUnlock()
	bo, err := session.EncodeGob(values)
	if err != nil {
		return
	}

	cs.b.Set(cs.sid, int(cs.maxlifetime), bo)
}

// SessionReleaseIfPresent Write couchbase session with Gob string
func (cs *SessionStore) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {
}

func (cp *Provider) getBucket() *couchbase.Bucket {
	c, err := couchbase.Connect(cp.SavePath)
	if err != nil {
		return nil
	}

	pool, err := c.GetPool(cp.Pool)
	if err != nil {
		return nil
	}

	bucket, err := pool.GetBucket(cp.Bucket)
	if err != nil {
		return nil
	}

	return bucket
}

// SessionInit init couchbase session
// savepath like couchbase server REST/JSON URL
// For v1.x e.g. http://host:port/, Pool, Bucket
// For v2.x, you should pass json string.
// e.g. { "save_path": "http://host:port/", "pool": "mypool", "bucket": "mybucket"}
func (cp *Provider) SessionInit(ctx context.Context, maxlifetime int64, cfg string) error {
	cp.maxlifetime = maxlifetime
	cfg = strings.TrimSpace(cfg)
	// we think this is v2.0, using json to init the session
	if strings.HasPrefix(cfg, "{") {
		return json.Unmarshal([]byte(cfg), cp)
	} else {
		return cp.initOldStyle(cfg)
	}
}

// initOldStyle keep compatible with v1.x
func (cp *Provider) initOldStyle(savePath string) error {
	configs := strings.Split(savePath, ",")
	if len(configs) > 0 {
		cp.SavePath = configs[0]
	}
	if len(configs) > 1 {
		cp.Pool = configs[1]
	}
	if len(configs) > 2 {
		cp.Bucket = configs[2]
	}

	return nil
}

// SessionRead read couchbase session by sid
func (cp *Provider) SessionRead(ctx context.Context, sid string) (session.Store, error) {
	cp.b = cp.getBucket()

	var (
		kv  map[interface{}]interface{}
		err error
		doc []byte
	)

	err = cp.b.Get(sid, &doc)
	if err != nil {
		return nil, err
	} else if doc == nil {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(doc)
		if err != nil {
			return nil, err
		}
	}

	cs := &SessionStore{b: cp.b, sid: sid, values: kv, maxlifetime: cp.maxlifetime}
	return cs, nil
}

// SessionExist Check couchbase session exist.
// it checkes sid exist or not.
func (cp *Provider) SessionExist(ctx context.Context, sid string) (bool, error) {
	cp.b = cp.getBucket()
	defer cp.b.Close()

	var doc []byte

	if err := cp.b.Get(sid, &doc); err != nil || doc == nil {
		return false, err
	}
	return true, nil
}

// SessionRegenerate remove oldsid and use sid to generate new session
func (cp *Provider) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	cp.b = cp.getBucket()

	var doc []byte
	if err := cp.b.Get(oldsid, &doc); err != nil || doc == nil {
		cp.b.Set(sid, int(cp.maxlifetime), "")
	} else {
		err := cp.b.Delete(oldsid)
		if err != nil {
			return nil, err
		}
		_, _ = cp.b.Add(sid, int(cp.maxlifetime), doc)
	}

	err := cp.b.Get(sid, &doc)
	if err != nil {
		return nil, err
	}
	var kv map[interface{}]interface{}
	if doc == nil {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(doc)
		if err != nil {
			return nil, err
		}
	}

	cs := &SessionStore{b: cp.b, sid: sid, values: kv, maxlifetime: cp.maxlifetime}
	return cs, nil
}

// SessionDestroy Remove Bucket in this couchbase
func (cp *Provider) SessionDestroy(ctx context.Context, sid string) error {
	cp.b = cp.getBucket()
	defer cp.b.Close()

	cp.b.Delete(sid)
	return nil
}

// SessionGC Recycle
func (cp *Provider) SessionGC(context.Context) {
}

// SessionAll return all active session
func (cp *Provider) SessionAll(context.Context) int {
	return 0
}

func init() {
	session.Register("couchbase", couchbpder)
}
