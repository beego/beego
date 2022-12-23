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

// Package postgres for session provider
//
// depends on github.com/lib/pq:
//
// go install github.com/lib/pq
//
//
// needs this table in your database:
//
// CREATE TABLE session (
// session_key	char(64) NOT NULL,
// session_data	bytea,
// session_expiry	timestamp NOT NULL,
// CONSTRAINT session_key PRIMARY KEY(session_key)
// );
//
// will be activated with these settings in app.conf:
//
// SessionOn = true
// SessionProvider = postgresql
// SessionSavePath = "user=a password=b dbname=c sslmode=disable"
// SessionName = session
//
//
// Usage:
// import(
//   _ "github.com/beego/beego/v2/server/web/session/postgresql"
//   "github.com/beego/beego/v2/server/web/session"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("postgresql", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"user=pqgotest dbname=pqgotest sslmode=verify-full"}``)
//		go globalSessions.GC()
//	}
//
package postgres

import (
	"context"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/beego/beego/v2/adapter/session"
	"github.com/beego/beego/v2/server/web/session/postgres"
)

// SessionStore postgresql session store
type SessionStore postgres.SessionStore

// Set value in postgresql session.
// it is temp value in map.
func (st *SessionStore) Set(key, value interface{}) error {
	return (*postgres.SessionStore)(st).Set(context.Background(), key, value)
}

// Get value from postgresql session
func (st *SessionStore) Get(key interface{}) interface{} {
	return (*postgres.SessionStore)(st).Get(context.Background(), key)
}

// Delete value in postgresql session
func (st *SessionStore) Delete(key interface{}) error {
	return (*postgres.SessionStore)(st).Delete(context.Background(), key)
}

// Flush clear all values in postgresql session
func (st *SessionStore) Flush() error {
	return (*postgres.SessionStore)(st).Flush(context.Background())
}

// SessionID get session id of this postgresql session store
func (st *SessionStore) SessionID() string {
	return (*postgres.SessionStore)(st).SessionID(context.Background())
}

// SessionRelease save postgresql session values to database.
// must call this method to save values to database.
func (st *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*postgres.SessionStore)(st).SessionRelease(context.Background(), w)
}

// Provider postgresql session provider
type Provider postgres.Provider

// SessionInit init postgresql session.
// savepath is the connection string of postgresql.
func (mp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*postgres.Provider)(mp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead get postgresql session by sid
func (mp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*postgres.Provider)(mp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check postgresql session exist
func (mp *Provider) SessionExist(sid string) bool {
	res, _ := (*postgres.Provider)(mp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for postgresql session
func (mp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*postgres.Provider)(mp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete postgresql session by sid
func (mp *Provider) SessionDestroy(sid string) error {
	return (*postgres.Provider)(mp).SessionDestroy(context.Background(), sid)
}

// SessionGC delete expired values in postgresql session
func (mp *Provider) SessionGC() {
	(*postgres.Provider)(mp).SessionGC(context.Background())
}

// SessionAll count values in postgresql session
func (mp *Provider) SessionAll() int {
	return (*postgres.Provider)(mp).SessionAll(context.Background())
}
