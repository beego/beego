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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	filepder      = &FileProvider{}
	gcmaxlifetime int64
)

// FileSessionStore File session store
type FileSessionStore struct {
	sid    string
	lock   sync.RWMutex
	values map[interface{}]interface{}
}

// Set value to file session
func (fs *FileSessionStore) Set(ctx context.Context, key, value interface{}) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.values[key] = value
	return nil
}

// Get value from file session
func (fs *FileSessionStore) Get(ctx context.Context, key interface{}) interface{} {
	fs.lock.RLock()
	defer fs.lock.RUnlock()
	if v, ok := fs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in file session by given key
func (fs *FileSessionStore) Delete(ctx context.Context, key interface{}) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	delete(fs.values, key)
	return nil
}

// Flush Clean all values in file session
func (fs *FileSessionStore) Flush(context.Context) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID Get file session store id
func (fs *FileSessionStore) SessionID(context.Context) string {
	return fs.sid
}

// SessionRelease Write file session to local file with Gob string
func (fs *FileSessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	filepder.lock.Lock()
	defer filepder.lock.Unlock()
	b, err := EncodeGob(fs.values)
	if err != nil {
		SLogger.Println(err)
		return
	}
	_, err = os.Stat(path.Join(filepder.savePath, string(fs.sid[0]), string(fs.sid[1]), fs.sid))
	var f *os.File
	if err == nil {
		f, err = os.OpenFile(path.Join(filepder.savePath, string(fs.sid[0]), string(fs.sid[1]), fs.sid), os.O_RDWR, 0o777)
		if err != nil {
			SLogger.Println(err)
			return
		}
	} else if os.IsNotExist(err) {
		f, err = os.Create(path.Join(filepder.savePath, string(fs.sid[0]), string(fs.sid[1]), fs.sid))
		if err != nil {
			SLogger.Println(err)
			return
		}
	} else {
		return
	}
	f.Truncate(0)
	f.Seek(0, 0)
	f.Write(b)
	f.Close()
}

// FileProvider File session provider
type FileProvider struct {
	lock        sync.RWMutex
	maxlifetime int64
	savePath    string
}

// SessionInit Init file session provider.
// savePath sets the session files path.
func (fp *FileProvider) SessionInit(ctx context.Context, maxlifetime int64, savePath string) error {
	fp.maxlifetime = maxlifetime
	fp.savePath = savePath
	return nil
}

// SessionRead Read file session by sid.
// if file is not exist, create it.
// the file path is generated from sid string.
func (fp *FileProvider) SessionRead(ctx context.Context, sid string) (Store, error) {
	invalidChars := "./"
	if strings.ContainsAny(sid, invalidChars) {
		return nil, errors.New("the sid shouldn't have following characters: " + invalidChars)
	}
	if len(sid) < 2 {
		return nil, errors.New("length of the sid is less than 2")
	}
	filepder.lock.Lock()
	defer filepder.lock.Unlock()

	sessionPath := filepath.Join(fp.savePath, string(sid[0]), string(sid[1]))
	sidPath := filepath.Join(sessionPath, sid)
	err := os.MkdirAll(sessionPath, 0o755)
	if err != nil {
		SLogger.Println(err.Error())
	}
	var f *os.File
	_, err = os.Stat(sidPath)
	switch {
	case err == nil:
		f, err = os.OpenFile(sidPath, os.O_RDWR, 0o777)
		if err != nil {
			return nil, err
		}
	case os.IsNotExist(err):
		f, err = os.Create(sidPath)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	defer f.Close()

	os.Chtimes(sidPath, time.Now(), time.Now())
	var kv map[interface{}]interface{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = DecodeGob(b)
		if err != nil {
			return nil, err
		}
	}

	ss := &FileSessionStore{sid: sid, values: kv}
	return ss, nil
}

// SessionExist Check file session exist.
// it checks the file named from sid exist or not.
func (fp *FileProvider) SessionExist(ctx context.Context, sid string) (bool, error) {
	filepder.lock.Lock()
	defer filepder.lock.Unlock()

	if len(sid) < 2 {
		SLogger.Println("min length of session id is 2 but got length: ", sid)
		return false, errors.New("min length of session id is 2")
	}

	_, err := os.Stat(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	return err == nil, nil
}

// SessionDestroy Remove all files in this save path
func (fp *FileProvider) SessionDestroy(ctx context.Context, sid string) error {
	filepder.lock.Lock()
	defer filepder.lock.Unlock()
	os.Remove(path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid))
	return nil
}

// SessionGC Recycle files in save path
func (fp *FileProvider) SessionGC(context.Context) {
	filepder.lock.Lock()
	defer filepder.lock.Unlock()

	gcmaxlifetime = fp.maxlifetime
	filepath.Walk(fp.savePath, gcpath)
}

// SessionAll Get active file session number.
// it walks save path to count files.
func (fp *FileProvider) SessionAll(context.Context) int {
	a := &activeSession{}
	err := filepath.Walk(fp.savePath, a.visit)
	if err != nil {
		SLogger.Printf("filepath.Walk() returned %v\n", err)
		return 0
	}
	return a.total
}

// SessionRegenerate Generate new sid for file session.
// it delete old file and create new file named from new sid.
func (fp *FileProvider) SessionRegenerate(ctx context.Context, oldsid, sid string) (Store, error) {
	filepder.lock.Lock()
	defer filepder.lock.Unlock()

	oldPath := path.Join(fp.savePath, string(oldsid[0]), string(oldsid[1]))
	oldSidFile := path.Join(oldPath, oldsid)
	newPath := path.Join(fp.savePath, string(sid[0]), string(sid[1]))
	newSidFile := path.Join(newPath, sid)

	// new sid file is exist
	_, err := os.Stat(newSidFile)
	if err == nil {
		return nil, fmt.Errorf("newsid %s exist", newSidFile)
	}

	err = os.MkdirAll(newPath, 0o755)
	if err != nil {
		SLogger.Println(err.Error())
	}

	// if old sid file exist
	// 1.read and parse file content
	// 2.write content to new sid file
	// 3.remove old sid file, change new sid file atime and ctime
	// 4.return FileSessionStore
	_, err = os.Stat(oldSidFile)
	if err == nil {
		b, err := ioutil.ReadFile(oldSidFile)
		if err != nil {
			return nil, err
		}

		var kv map[interface{}]interface{}
		if len(b) == 0 {
			kv = make(map[interface{}]interface{})
		} else {
			kv, err = DecodeGob(b)
			if err != nil {
				return nil, err
			}
		}

		ioutil.WriteFile(newSidFile, b, 0o777)
		os.Remove(oldSidFile)
		os.Chtimes(newSidFile, time.Now(), time.Now())
		ss := &FileSessionStore{sid: sid, values: kv}
		return ss, nil
	}

	// if old sid file not exist, just create new sid file and return
	newf, err := os.Create(newSidFile)
	if err != nil {
		return nil, err
	}
	newf.Close()
	ss := &FileSessionStore{sid: sid, values: make(map[interface{}]interface{})}
	return ss, nil
}

// remove file in save path if expired
func gcpath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if (info.ModTime().Unix() + gcmaxlifetime) < time.Now().Unix() {
		os.Remove(path)
	}
	return nil
}

type activeSession struct {
	total int
}

func (as *activeSession) visit(paths string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	as.total = as.total + 1
	return nil
}

func init() {
	Register("file", filepder)
}
