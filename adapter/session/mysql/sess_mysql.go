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

// Package mysql for session provider
//
// depends on github.com/go-sql-driver/mysql:
//
// go install github.com/go-sql-driver/mysql
//
// mysql session support need create table as sql:
//	CREATE TABLE `session` (
//	`session_key` char(64) NOT NULL,
//	`session_data` blob,
//	`session_expiry` int(11) unsigned NOT NULL,
//	PRIMARY KEY (`session_key`)
//	) ENGINE=MyISAM DEFAULT CHARSET=utf8;
//
// Usage:
// import(
//   _ "github.com/beego/beego/session/mysql"
//   "github.com/beego/beego/session"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("mysql", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]"}``)
//		go globalSessions.GC()
//	}
//
// more docs: http://beego.me/docs/module/session.md
package mysql

import (
	"context"
	"net/http"

	"github.com/beego/beego/adapter/session"
	"github.com/beego/beego/server/web/session/mysql"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	// TableName store the session in MySQL
	TableName = mysql.TableName
	mysqlpder = &Provider{}
)

// SessionStore mysql session store
type SessionStore mysql.SessionStore

// Set value in mysql session.
// it is temp value in map.
func (st *SessionStore) Set(key, value interface{}) error {
	return (*mysql.SessionStore)(st).Set(context.Background(), key, value)
}

// Get value from mysql session
func (st *SessionStore) Get(key interface{}) interface{} {
	return (*mysql.SessionStore)(st).Get(context.Background(), key)
}

// Delete value in mysql session
func (st *SessionStore) Delete(key interface{}) error {
	return (*mysql.SessionStore)(st).Delete(context.Background(), key)
}

// Flush clear all values in mysql session
func (st *SessionStore) Flush() error {
	return (*mysql.SessionStore)(st).Flush(context.Background())
}

// SessionID get session id of this mysql session store
func (st *SessionStore) SessionID() string {
	return (*mysql.SessionStore)(st).SessionID(context.Background())
}

// SessionRelease save mysql session values to database.
// must call this method to save values to database.
func (st *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*mysql.SessionStore)(st).SessionRelease(context.Background(), w)
}

// Provider mysql session provider
type Provider mysql.Provider

// SessionInit init mysql session.
// savepath is the connection string of mysql.
func (mp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	return (*mysql.Provider)(mp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead get mysql session by sid
func (mp *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*mysql.Provider)(mp).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist check mysql session exist
func (mp *Provider) SessionExist(sid string) bool {
	res, _ := (*mysql.Provider)(mp).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for mysql session
func (mp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*mysql.Provider)(mp).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy delete mysql session by sid
func (mp *Provider) SessionDestroy(sid string) error {
	return (*mysql.Provider)(mp).SessionDestroy(context.Background(), sid)
}

// SessionGC delete expired values in mysql session
func (mp *Provider) SessionGC() {
	(*mysql.Provider)(mp).SessionGC(context.Background())
}

// SessionAll count values in mysql session
func (mp *Provider) SessionAll() int {
	return (*mysql.Provider)(mp).SessionAll(context.Background())
}
