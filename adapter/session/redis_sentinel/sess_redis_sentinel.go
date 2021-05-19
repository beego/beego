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

// Package redis for session provider
//
// depend on github.com/go-redis/redis
//
// go install github.com/go-redis/redis
//
// Usage:
// import(
//   _ "github.com/beego/beego/v2/server/web/session/redis_sentinel"
//   "github.com/beego/beego/v2/server/web/session"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("redis_sentinel", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:26379;127.0.0.2:26379"}``)
//		go globalSessions.GC()
//	}
//
// more detail about params: please check the notes on the function SessionInit in this package
package redis_sentinel

import (
	"context"
	"net/http"

	"github.com/beego/beego/v2/adapter/session"
	sentinel "github.com/beego/beego/v2/server/web/session/redis_sentinel"
)

// DefaultPoolSize redis_sentinel default pool size
var DefaultPoolSize = sentinel.DefaultPoolSize

// SessionStore redis_sentinel session store
type SessionStore sentinel.SessionStore

// Set value in redis_sentinel session
func (rs *SessionStore) Set(key, value interface{}) error {
	return (*sentinel.SessionStore)(rs).Set(context.Background(), key, value)
}

// Get value in redis_sentinel session
func (rs *SessionStore) Get(key interface{}) interface{} {
	return (*sentinel.SessionStore)(rs).Get(context.Background(), key)
}

// Delete value in redis_sentinel session
func (rs *SessionStore) Delete(key interface{}) error {
	return (*sentinel.SessionStore)(rs).Delete(context.Background(), key)
}

// Flush clear all values in redis_sentinel session
func (rs *SessionStore) Flush() error {
	return (*sentinel.SessionStore)(rs).Flush(context.Background())
}

// SessionID get redis_sentinel session id
func (rs *SessionStore) SessionID() string {
	return (*sentinel.SessionStore)(rs).SessionID(context.Background())
}

// SessionRelease save session values to redis_sentinel
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*sentinel.SessionStore)(rs).SessionRelease(context.Background(), w)
}

// Provider redis_sentinel session provider
type Provider sentinel.Provider

// SessionInit init redis_sentinel session
// savepath like redis sentinel addr,pool size,password,dbnum,masterName
// e.g. 127.0.0.1:26379;127.0.0.2:26379,100,1qaz2wsx,0,mymaster
func (rp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*sentinel.Provider)(rp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read redis_sentinel session by sid
func (rp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*sentinel.Provider)(rp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check redis_sentinel session exist by sid
func (rp *Provider) SessionExist(sid string) bool {
	res, _ := (*sentinel.Provider)(rp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for redis_sentinel session
func (rp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*sentinel.Provider)(rp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete redis session by id
func (rp *Provider) SessionDestroy(sid string) error {
	return (*sentinel.Provider)(rp).SessionDestroy(context.Background(), sid)
}

// SessionGC Impelment method, no used.
func (rp *Provider) SessionGC() {
	(*sentinel.Provider)(rp).SessionGC(context.Background())
}

// SessionAll return all activeSession
func (rp *Provider) SessionAll() int {
	return (*sentinel.Provider)(rp).SessionAll(context.Background())
}
