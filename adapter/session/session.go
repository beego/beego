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

// Package session provider
//
// Usage:
// import(
//   "github.com/beego/beego/v2/server/web/session"
// )
//
//	func init() {
//      globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "cookieLifeTime": 3600, "providerConfig": ""}`)
//		go globalSessions.GC()
//	}
//
package session

import (
	"io"
	"net/http"
	"os"

	"github.com/beego/beego/v2/server/web/session"
)

// Store contains all data for one session process with specific id.
type Store interface {
	Set(key, value interface{}) error     // set session value
	Get(key interface{}) interface{}      // get session value
	Delete(key interface{}) error         // delete session value
	SessionID() string                    // back current sessionID
	SessionRelease(w http.ResponseWriter) // release the resource & save data to provider & return the data
	Flush() error                         // delete all data
}

// Provider contains global session methods and saved SessionStores.
// it can operate a SessionStore by its id.
type Provider interface {
	SessionInit(gclifetime int64, config string) error
	SessionRead(sid string) (Store, error)
	SessionExist(sid string) bool
	SessionRegenerate(oldsid, sid string) (Store, error)
	SessionDestroy(sid string) error
	SessionAll() int // get all active session
	SessionGC()
}

// SLogger a helpful variable to log information about session
var SLogger = NewSessionLog(os.Stderr)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provide Provider) {
	session.Register(name, &oldToNewProviderAdapter{
		delegate: provide,
	})
}

// GetProvider
func GetProvider(name string) (Provider, error) {
	res, err := session.GetProvider(name)
	if adt, ok := res.(*oldToNewProviderAdapter); err == nil && ok {
		return adt.delegate, err
	}

	return &newToOldProviderAdapter{
		delegate: res,
	}, err
}

// ManagerConfig define the session config
type ManagerConfig session.ManagerConfig

// Manager contains Provider and its configuration.
type Manager session.Manager

// NewManager Create new Manager with provider name and json config string.
// provider name:
// 1. cookie
// 2. file
// 3. memory
// 4. redis
// 5. mysql
// json config:
// 1. is https  default false
// 2. hashfunc  default sha1
// 3. hashkey default beegosessionkey
// 4. maxage default is none
func NewManager(provideName string, cf *ManagerConfig) (*Manager, error) {
	m, err := session.NewManager(provideName, (*session.ManagerConfig)(cf))
	return (*Manager)(m), err
}

// GetProvider return current manager's provider
func (manager *Manager) GetProvider() Provider {
	return &newToOldProviderAdapter{
		delegate: (*session.Manager)(manager).GetProvider(),
	}
}

// SessionStart generate or read the session id from http request.
// if session id exists, return SessionStore with this id.
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (Store, error) {
	s, err := (*session.Manager)(manager).SessionStart(w, r)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionDestroy Destroy session by its id in http request cookie.
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	(*session.Manager)(manager).SessionDestroy(w, r)
}

// GetSessionStore Get SessionStore by its id.
func (manager *Manager) GetSessionStore(sid string) (Store, error) {
	s, err := (*session.Manager)(manager).GetSessionStore(sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// GC Start session gc process.
// it can do gc in times after gc lifetime.
func (manager *Manager) GC() {
	(*session.Manager)(manager).GC()
}

// SessionRegenerateID Regenerate a session id for this SessionStore who's id is saving in http request.
func (manager *Manager) SessionRegenerateID(w http.ResponseWriter, r *http.Request) Store {
	s, _ := (*session.Manager)(manager).SessionRegenerateID(w, r)
	return &NewToOldStoreAdapter{
		delegate: s,
	}
}

// GetActiveSession Get all active sessions count number.
func (manager *Manager) GetActiveSession() int {
	return (*session.Manager)(manager).GetActiveSession()
}

// SetSecure Set cookie with https.
func (manager *Manager) SetSecure(secure bool) {
	(*session.Manager)(manager).SetSecure(secure)
}

// Log implement the log.Logger
type Log session.Log

// NewSessionLog set io.Writer to create a Logger for session.
func NewSessionLog(out io.Writer) *Log {
	return (*Log)(session.NewSessionLog(out))
}
