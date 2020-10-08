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
	"context"
	"net/http"

	"github.com/astaxie/beego/server/web/session"
)

// MemSessionStore memory session store.
// it saved sessions in a map in memory.
type MemSessionStore session.MemSessionStore

// Set value to memory session
func (st *MemSessionStore) Set(key, value interface{}) error {
	return (*session.MemSessionStore)(st).Set(context.Background(), key, value)
}

// Get value from memory session by key
func (st *MemSessionStore) Get(key interface{}) interface{} {
	return (*session.MemSessionStore)(st).Get(context.Background(), key)
}

// Delete in memory session by key
func (st *MemSessionStore) Delete(key interface{}) error {
	return (*session.MemSessionStore)(st).Delete(context.Background(), key)
}

// Flush clear all values in memory session
func (st *MemSessionStore) Flush() error {
	return (*session.MemSessionStore)(st).Flush(context.Background())
}

// SessionID get this id of memory session store
func (st *MemSessionStore) SessionID() string {
	return (*session.MemSessionStore)(st).SessionID(context.Background())
}

// SessionRelease Implement method, no used.
func (st *MemSessionStore) SessionRelease(w http.ResponseWriter) {
	(*session.MemSessionStore)(st).SessionRelease(context.Background(), w)
}

// MemProvider Implement the provider interface
type MemProvider session.MemProvider

// SessionInit init memory session
func (pder *MemProvider) SessionInit(maxlifetime int64, savePath string) error {
	return (*session.MemProvider)(pder).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead get memory session store by sid
func (pder *MemProvider) SessionRead(sid string) (Store, error) {
	s, err := (*session.MemProvider)(pder).SessionRead(context.Background(), sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionExist check session store exist in memory session by sid
func (pder *MemProvider) SessionExist(sid string) bool {
	res, _ := (*session.MemProvider)(pder).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate generate new sid for session store in memory session
func (pder *MemProvider) SessionRegenerate(oldsid, sid string) (Store, error) {
	s, err := (*session.MemProvider)(pder).SessionRegenerate(context.Background(), oldsid, sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionDestroy delete session store in memory session by id
func (pder *MemProvider) SessionDestroy(sid string) error {
	return (*session.MemProvider)(pder).SessionDestroy(context.Background(), sid)
}

// SessionGC clean expired session stores in memory session
func (pder *MemProvider) SessionGC() {
	(*session.MemProvider)(pder).SessionGC(context.Background())
}

// SessionAll get count number of memory session
func (pder *MemProvider) SessionAll() int {
	return (*session.MemProvider)(pder).SessionAll(context.Background())
}

// SessionUpdate expand time of session store by id in memory session
func (pder *MemProvider) SessionUpdate(sid string) error {
	return (*session.MemProvider)(pder).SessionUpdate(context.Background(), sid)
}
