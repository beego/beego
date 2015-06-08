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

package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

func init() {
	Register("file", NewFileCache())
}

// FileCacheItem is basic unit of file cache adapter.
// it contains data and expire time.
type FileCacheItem struct {
	Data       interface{}
	Lastaccess int64
	Expired    int64
}

var (
	FileCachePath           string = "cache" // cache directory
	FileCacheFileSuffix     string = ".bin"  // cache file suffix
	FileCacheDirectoryLevel int    = 2       // cache file deep level if auto generated cache files.
	FileCacheEmbedExpiry    int64  = 0       // cache expire time, default is no expire forever.
)

// FileCache is cache adapter for file storage.
type FileCache struct {
	CachePath      string
	FileSuffix     string
	DirectoryLevel int
	EmbedExpiry    int
}

// Create new file cache with no config.
// the level and expiry need set in method StartAndGC as config string.
func NewFileCache() *FileCache {
	//    return &FileCache{CachePath:FileCachePath, FileSuffix:FileCacheFileSuffix}
	return &FileCache{}
}

// Start and begin gc for file cache.
// the config need to be like {CachePath:"/cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}
func (fc *FileCache) StartAndGC(config string) error {

	var cfg map[string]string
	json.Unmarshal([]byte(config), &cfg)
	if _, ok := cfg["CachePath"]; !ok {
		cfg["CachePath"] = FileCachePath
	}
	if _, ok := cfg["FileSuffix"]; !ok {
		cfg["FileSuffix"] = FileCacheFileSuffix
	}
	if _, ok := cfg["DirectoryLevel"]; !ok {
		cfg["DirectoryLevel"] = strconv.Itoa(FileCacheDirectoryLevel)
	}
	if _, ok := cfg["EmbedExpiry"]; !ok {
		cfg["EmbedExpiry"] = strconv.FormatInt(FileCacheEmbedExpiry, 10)
	}
	fc.CachePath = cfg["CachePath"]
	fc.FileSuffix = cfg["FileSuffix"]
	fc.DirectoryLevel, _ = strconv.Atoi(cfg["DirectoryLevel"])
	fc.EmbedExpiry, _ = strconv.Atoi(cfg["EmbedExpiry"])

	fc.Init()
	return nil
}

// Init will make new dir for file cache if not exist.
func (fc *FileCache) Init() {
	if ok, _ := exists(fc.CachePath); !ok { // todo : error handle
		_ = os.MkdirAll(fc.CachePath, os.ModePerm) // todo : error handle
	}
}

// get cached file name. it's md5 encoded.
func (fc *FileCache) getCacheFileName(key string) string {
	m := md5.New()
	io.WriteString(m, key)
	keyMd5 := hex.EncodeToString(m.Sum(nil))
	cachePath := fc.CachePath
	switch fc.DirectoryLevel {
	case 2:
		cachePath = filepath.Join(cachePath, keyMd5[0:2], keyMd5[2:4])
	case 1:
		cachePath = filepath.Join(cachePath, keyMd5[0:2])
	}

	if ok, _ := exists(cachePath); !ok { // todo : error handle
		_ = os.MkdirAll(cachePath, os.ModePerm) // todo : error handle
	}

	return filepath.Join(cachePath, fmt.Sprintf("%s%s", keyMd5, fc.FileSuffix))
}

// Get value from file cache.
// if non-exist or expired, return empty string.
func (fc *FileCache) Get(key string) interface{} {
	fileData, err := File_get_contents(fc.getCacheFileName(key))
	if err != nil {
		return ""
	}
	var to FileCacheItem
	Gob_decode(fileData, &to)
	if to.Expired < time.Now().Unix() {
		return ""
	}
	return to.Data
}

// GetMulti gets values from file cache.
// if non-exist or expired, return empty string.
func (fc *FileCache) GetMulti(keys []string) []interface{} {
	var rc []interface{}
	for _, key := range keys {
		rc = append(rc, fc.Get(key))
	}
	return rc
}

// Put value into file cache.
// timeout means how long to keep this file, unit of ms.
// if timeout equals FileCacheEmbedExpiry(default is 0), cache this item forever.
func (fc *FileCache) Put(key string, val interface{}, timeout int64) error {
	gob.Register(val)

	item := FileCacheItem{Data: val}
	if timeout == FileCacheEmbedExpiry {
		item.Expired = time.Now().Unix() + (86400 * 365 * 10) // ten years
	} else {
		item.Expired = time.Now().Unix() + timeout
	}
	item.Lastaccess = time.Now().Unix()
	data, err := Gob_encode(item)
	if err != nil {
		return err
	}
	return File_put_contents(fc.getCacheFileName(key), data)
}

// Delete file cache value.
func (fc *FileCache) Delete(key string) error {
	filename := fc.getCacheFileName(key)
	if ok, _ := exists(filename); ok {
		return os.Remove(filename)
	}
	return nil
}

// Increase cached int value.
// fc value is saving forever unless Delete.
func (fc *FileCache) Incr(key string) error {
	data := fc.Get(key)
	var incr int
	if reflect.TypeOf(data).Name() != "int" {
		incr = 0
	} else {
		incr = data.(int) + 1
	}
	fc.Put(key, incr, FileCacheEmbedExpiry)
	return nil
}

// Decrease cached int value.
func (fc *FileCache) Decr(key string) error {
	data := fc.Get(key)
	var decr int
	if reflect.TypeOf(data).Name() != "int" || data.(int)-1 <= 0 {
		decr = 0
	} else {
		decr = data.(int) - 1
	}
	fc.Put(key, decr, FileCacheEmbedExpiry)
	return nil
}

// Check value is exist.
func (fc *FileCache) IsExist(key string) bool {
	ret, _ := exists(fc.getCacheFileName(key))
	return ret
}

// Clean cached files.
// not implemented.
func (fc *FileCache) ClearAll() error {
	return nil
}

// check file exist.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Get bytes to file.
// if non-exist, create this file.
func File_get_contents(filename string) (data []byte, e error) {
	f, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if e != nil {
		return
	}
	defer f.Close()
	stat, e := f.Stat()
	if e != nil {
		return
	}
	data = make([]byte, stat.Size())
	result, e := f.Read(data)
	if e != nil || int64(result) != stat.Size() {
		return nil, e
	}
	return
}

// Put bytes to file.
// if non-exist, create this file.
func File_put_contents(filename string, content []byte) error {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(content)
	return err
}

// Gob encodes file cache item.
func Gob_encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Gob decodes file cache item.
func Gob_decode(data []byte, to *FileCacheItem) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&to)
}
