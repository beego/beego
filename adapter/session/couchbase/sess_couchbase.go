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
//   _ "github.com/beego/beego/v2/session/couchbase"
//   "github.com/beego/beego/v2/session"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("couchbase", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"http://host:port/, Pool, Bucket"}``)
//		go globalSessions.GC()
//	}
//
// more docs: http://beego.me/docs/module/session.md
package couchbase

import (
	"context"
	"net/http"

	"github.com/beego/beego/v2/adapter/session"
	beecb "github.com/beego/beego/v2/server/web/session/couchbase"
)

// SessionStore store each session
type SessionStore beecb.SessionStore

// Provider couchabse provided
type Provider beecb.Provider

// Set value to couchabse session
func (cs *SessionStore) Set(key, value interface{}) error {
	return (*beecb.SessionStore)(cs).Set(context.Background(), key, value)
}

// Get value from couchabse session
func (cs *SessionStore) Get(key interface{}) interface{} {
	return (*beecb.SessionStore)(cs).Get(context.Background(), key)
}

// Delete value in couchbase session by given key
func (cs *SessionStore) Delete(key interface{}) error {
	return (*beecb.SessionStore)(cs).Delete(context.Background(), key)
}

// Flush Clean all values in couchbase session
func (cs *SessionStore) Flush() error {
	return (*beecb.SessionStore)(cs).Flush(context.Background())
}

// SessionID Get couchbase session store id
func (cs *SessionStore) SessionID() string {
	return (*beecb.SessionStore)(cs).SessionID(context.Background())
}

// SessionRelease Write couchbase session with Gob string
func (cs *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*beecb.SessionStore)(cs).SessionRelease(context.Background(), w)
}

// SessionInit init couchbase session
// savepath like couchbase server REST/JSON URL
// e.g. http://host:port/, Pool, Bucket
func (cp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*beecb.Provider)(cp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read couchbase session by sid
func (cp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*beecb.Provider)(cp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist Check couchbase session exist.
// it checkes sid exist or not.
func (cp *Provider) SessionExist(sid string) bool {
	res, _ := (*beecb.Provider)(cp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate remove oldsid and use sid to generate new session
func (cp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*beecb.Provider)(cp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy Remove bucket in this couchbase
func (cp *Provider) SessionDestroy(sid string) error {
	return (*beecb.Provider)(cp).SessionDestroy(context.Background(), sid)
}

// SessionGC Recycle
func (cp *Provider) SessionGC() {
	(*beecb.Provider)(cp).SessionGC(context.Background())
}

// SessionAll return all active session
func (cp *Provider) SessionAll() int {
	return (*beecb.Provider)(cp).SessionAll(context.Background())
}
