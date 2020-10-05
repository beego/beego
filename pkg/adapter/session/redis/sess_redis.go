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
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/session/redis"
//   "github.com/astaxie/beego/session"
// )
//
// 	func init() {
// 		globalSessions, _ = session.NewManager("redis", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:7070"}``)
// 		go globalSessions.GC()
// 	}
//
// more docs: http://beego.me/docs/module/session.md
package redis

import (
	"context"
	"net/http"

	"github.com/astaxie/beego/pkg/adapter/session"

	beeRedis "github.com/astaxie/beego/pkg/core/session/redis"
)

// MaxPoolSize redis max pool size
var MaxPoolSize = beeRedis.MaxPoolSize

// SessionStore redis session store
type SessionStore beeRedis.SessionStore

// Set value in redis session
func (rs *SessionStore) Set(key, value interface{}) error {
	return (*beeRedis.SessionStore)(rs).Set(context.Background(), key, value)
}

// Get value in redis session
func (rs *SessionStore) Get(key interface{}) interface{} {
	return (*beeRedis.SessionStore)(rs).Get(context.Background(), key)
}

// Delete value in redis session
func (rs *SessionStore) Delete(key interface{}) error {
	return (*beeRedis.SessionStore)(rs).Delete(context.Background(), key)
}

// Flush clear all values in redis session
func (rs *SessionStore) Flush() error {
	return (*beeRedis.SessionStore)(rs).Flush(context.Background())
}

// SessionID get redis session id
func (rs *SessionStore) SessionID() string {
	return (*beeRedis.SessionStore)(rs).SessionID(context.Background())
}

// SessionRelease save session values to redis
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*beeRedis.SessionStore)(rs).SessionRelease(context.Background(), w)
}

// Provider redis session provider
type Provider beeRedis.Provider

// SessionInit init redis session
// savepath like redis server addr,pool size,password,dbnum,IdleTimeout second
// e.g. 127.0.0.1:6379,100,astaxie,0,30
func (rp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*beeRedis.Provider)(rp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead read redis session by sid
func (rp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*beeRedis.Provider)(rp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check redis session exist by sid
func (rp *Provider) SessionExist(sid string) bool {
	res, _ := (*beeRedis.Provider)(rp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for redis session
func (rp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*beeRedis.Provider)(rp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete redis session by id
func (rp *Provider) SessionDestroy(sid string) error {
	return (*beeRedis.Provider)(rp).SessionDestroy(context.Background(), sid)
}

// SessionGC Impelment method, no used.
func (rp *Provider) SessionGC() {
	(*beeRedis.Provider)(rp).SessionGC(context.Background())
}

// SessionAll return all activeSession
func (rp *Provider) SessionAll() int {
	return (*beeRedis.Provider)(rp).SessionAll(context.Background())
}
