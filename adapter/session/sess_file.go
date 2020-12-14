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

	"github.com/beego/beego/v2/server/web/session"
)

// FileSessionStore File session store
type FileSessionStore session.FileSessionStore

// Set value to file session
func (fs *FileSessionStore) Set(key, value interface{}) error {
	return (*session.FileSessionStore)(fs).Set(context.Background(), key, value)
}

// Get value from file session
func (fs *FileSessionStore) Get(key interface{}) interface{} {
	return (*session.FileSessionStore)(fs).Get(context.Background(), key)
}

// Delete value in file session by given key
func (fs *FileSessionStore) Delete(key interface{}) error {
	return (*session.FileSessionStore)(fs).Delete(context.Background(), key)
}

// Flush Clean all values in file session
func (fs *FileSessionStore) Flush() error {
	return (*session.FileSessionStore)(fs).Flush(context.Background())
}

// SessionID Get file session store id
func (fs *FileSessionStore) SessionID() string {
	return (*session.FileSessionStore)(fs).SessionID(context.Background())
}

// SessionRelease Write file session to local file with Gob string
func (fs *FileSessionStore) SessionRelease(w http.ResponseWriter) {
	(*session.FileSessionStore)(fs).SessionRelease(context.Background(), w)
}

// FileProvider File session provider
type FileProvider session.FileProvider

// SessionInit Init file session provider.
// savePath sets the session files path.
func (fp *FileProvider) SessionInit(maxlifetime int64, savePath string) error {
	return (*session.FileProvider)(fp).SessionInit(context.Background(), maxlifetime, savePath)
}

// SessionRead Read file session by sid.
// if file is not exist, create it.
// the file path is generated from sid string.
func (fp *FileProvider) SessionRead(sid string) (Store, error) {
	s, err := (*session.FileProvider)(fp).SessionRead(context.Background(), sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

// SessionExist Check file session exist.
// it checks the file named from sid exist or not.
func (fp *FileProvider) SessionExist(sid string) bool {
	res, _ := (*session.FileProvider)(fp).SessionExist(context.Background(), sid)
	return res
}

// SessionDestroy Remove all files in this save path
func (fp *FileProvider) SessionDestroy(sid string) error {
	return (*session.FileProvider)(fp).SessionDestroy(context.Background(), sid)
}

// SessionGC Recycle files in save path
func (fp *FileProvider) SessionGC() {
	(*session.FileProvider)(fp).SessionGC(context.Background())
}

// SessionAll Get active file session number.
// it walks save path to count files.
func (fp *FileProvider) SessionAll() int {
	return (*session.FileProvider)(fp).SessionAll(context.Background())
}

// SessionRegenerate Generate new sid for file session.
// it delete old file and create new file named from new sid.
func (fp *FileProvider) SessionRegenerate(oldsid, sid string) (Store, error) {
	s, err := (*session.FileProvider)(fp).SessionRegenerate(context.Background(), oldsid, sid)
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}
