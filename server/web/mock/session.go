// Copyright 2021 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mock

import (
	"context"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/google/uuid"
	"net/http"
)

// NewSessionProvider create new SessionProvider
// and you could use it to mock data
// Parameter "name" is the real SessionProvider you used
func NewSessionProvider(name string) *SessionProvider {
	sp := newSessionProvider()
	session.Register(name, sp)
	web.GlobalSessions, _ = session.NewManager(name, session.NewManagerConfig())
	return sp
}

// SessionProvider will replace session provider with "mock" provider
type SessionProvider struct {
	Store *SessionStore
}

func newSessionProvider() *SessionProvider {
	return &SessionProvider{
		Store: newSessionStore(),
	}
}

// SessionInit do nothing
func (s *SessionProvider) SessionInit(ctx context.Context, gclifetime int64, config string) error {
	return nil
}

// SessionRead return Store
func (s *SessionProvider) SessionRead(ctx context.Context, sid string) (session.Store, error) {
	return s.Store, nil
}

// SessionExist always return true
func (s *SessionProvider) SessionExist(ctx context.Context, sid string) (bool, error) {
	return true, nil
}

// SessionRegenerate create new Store
func (s *SessionProvider) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	s.Store = newSessionStore()
	return s.Store, nil
}

// SessionDestroy reset Store to nil
func (s *SessionProvider) SessionDestroy(ctx context.Context, sid string) error {
	s.Store = nil;
	return nil
}

// SessionAll return 0
func (s *SessionProvider) SessionAll(ctx context.Context) int {
	return 0
}

// SessionGC do nothing
func (s *SessionProvider) SessionGC(ctx context.Context) {
	// we do anything since we don't need to mock GC
}


type SessionStore struct {
	sid string
	values map[interface{}]interface{}
}

func (s *SessionStore) Set(ctx context.Context, key, value interface{}) error {
	s.values[key]=value
	return nil
}

func (s *SessionStore) Get(ctx context.Context, key interface{}) interface{} {
	return s.values[key]
}

func (s *SessionStore) Delete(ctx context.Context, key interface{}) error {
	delete(s.values, key)
	return nil
}

func (s *SessionStore) SessionID(ctx context.Context) string {
	return s.sid
}

// SessionRelease do nothing
func (s *SessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	// Support in the future if necessary, now I think we don't need to implement this
}

func (s *SessionStore) Flush(ctx context.Context) error {
	s.values = make(map[interface{}]interface{}, 4)
	return nil
}

func newSessionStore() *SessionStore {
	return &SessionStore{
		sid: uuid.New().String(),
		values: make(map[interface{}]interface{}, 4),
	}
}



