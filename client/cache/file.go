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
	"context"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/berror"
)

// FileCacheItem is basic unit of file cache adapter which
// contains data and expire time.
type FileCacheItem struct {
	Data       interface{}
	Lastaccess time.Time
	Expired    time.Time
}

// FileCache Config
var (
	FileCachePath           = "cache"     // cache directory
	FileCacheFileSuffix     = ".bin"      // cache file suffix
	FileCacheDirectoryLevel = 2           // cache file deep level if auto generated cache files.
	FileCacheEmbedExpiry    time.Duration // cache expire time, default is no expire forever.
)

// FileCache is cache adapter for file storage.
type FileCache struct {
	CachePath      string
	FileSuffix     string
	DirectoryLevel int
	EmbedExpiry    int
}

// NewFileCache creates a new file cache with no config.
// The level and expiry need to be set in the method StartAndGC as config string.
func NewFileCache() Cache {
	//    return &FileCache{CachePath:FileCachePath, FileSuffix:FileCacheFileSuffix}
	return &FileCache{}
}

// StartAndGC starts gc for file cache.
// config must be in the format {CachePath:"/cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}
func (fc *FileCache) StartAndGC(config string) error {

	cfg := make(map[string]string)
	err := json.Unmarshal([]byte(config), &cfg)
	if err != nil {
		return err
	}

	const cpKey = "CachePath"
	const fsKey = "FileSuffix"
	const dlKey = "DirectoryLevel"
	const eeKey = "EmbedExpiry"

	if _, ok := cfg[cpKey]; !ok {
		cfg[cpKey] = FileCachePath
	}

	if _, ok := cfg[fsKey]; !ok {
		cfg[fsKey] = FileCacheFileSuffix
	}

	if _, ok := cfg[dlKey]; !ok {
		cfg[dlKey] = strconv.Itoa(FileCacheDirectoryLevel)
	}

	if _, ok := cfg[eeKey]; !ok {
		cfg[eeKey] = strconv.FormatInt(int64(FileCacheEmbedExpiry.Seconds()), 10)
	}
	fc.CachePath = cfg[cpKey]
	fc.FileSuffix = cfg[fsKey]
	fc.DirectoryLevel, err = strconv.Atoi(cfg[dlKey])
	if err != nil {
		return berror.Wrapf(err, InvalidFileCacheDirectoryLevelCfg,
			"invalid directory level config, please check your input, it must be integer: %s", cfg[dlKey])
	}
	fc.EmbedExpiry, err = strconv.Atoi(cfg[eeKey])
	if err != nil {
		return berror.Wrapf(err, InvalidFileCacheEmbedExpiryCfg,
			"invalid embed expiry config, please check your input, it must be integer: %s", cfg[eeKey])
	}
	return fc.Init()
}

// Init makes new a dir for file cache if it does not already exist
func (fc *FileCache) Init() error {
	ok, err := exists(fc.CachePath)
	if err != nil || ok {
		return err
	}
	err = os.MkdirAll(fc.CachePath, os.ModePerm)
	if err != nil {
		return berror.Wrapf(err, CreateFileCacheDirFailed,
			"could not create directory, please check the config [%s] and file mode.", fc.CachePath)
	}
	return nil
}

// getCachedFilename returns an md5 encoded file name.
func (fc *FileCache) getCacheFileName(key string) (string, error) {
	m := md5.New()
	_, _ = io.WriteString(m, key)
	keyMd5 := hex.EncodeToString(m.Sum(nil))
	cachePath := fc.CachePath
	switch fc.DirectoryLevel {
	case 2:
		cachePath = filepath.Join(cachePath, keyMd5[0:2], keyMd5[2:4])
	case 1:
		cachePath = filepath.Join(cachePath, keyMd5[0:2])
	}
	ok, err := exists(cachePath)
	if err != nil {
		return "", err
	}
	if !ok {
		err  = os.MkdirAll(cachePath, os.ModePerm)
		if err != nil {
			return "", berror.Wrapf(err, CreateFileCacheDirFailed,
				"could not create the directory: %s", cachePath)
		}
	}

	return filepath.Join(cachePath, fmt.Sprintf("%s%s", keyMd5, fc.FileSuffix)), nil
}

// Get value from file cache.
// if nonexistent or expired return an empty string.
func (fc *FileCache) Get(ctx context.Context, key string) (interface{}, error) {
	fn, err := fc.getCacheFileName(key)
	if err != nil {
		return nil, err
	}
	fileData, err := FileGetContents(fn)
	if err != nil {
		return nil, err
	}

	var to FileCacheItem
	err = GobDecode(fileData, &to)
	if err != nil {
		return nil, err
	}

	if to.Expired.Before(time.Now()) {
		return nil, ErrKeyExpired
	}
	return to.Data, nil
}

// GetMulti gets values from file cache.
// if nonexistent or expired return an empty string.
func (fc *FileCache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	rc := make([]interface{}, len(keys))
	keysErr := make([]string, 0)

	for i, ki := range keys {
		val, err := fc.Get(context.Background(), ki)
		if err != nil {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, err.Error()))
			continue
		}
		rc[i] = val
	}

	if len(keysErr) == 0 {
		return rc, nil
	}
	return rc, berror.Error(MultiGetFailed, strings.Join(keysErr, "; "))
}

// Put value into file cache.
// timeout: how long this file should be kept in ms
// if timeout equals fc.EmbedExpiry(default is 0), cache this item forever.
func (fc *FileCache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	gob.Register(val)

	item := FileCacheItem{Data: val}
	if timeout == time.Duration(fc.EmbedExpiry) {
		item.Expired = time.Now().Add((86400 * 365 * 10) * time.Second) // ten years
	} else {
		item.Expired = time.Now().Add(timeout)
	}
	item.Lastaccess = time.Now()
	data, err := GobEncode(item)
	if err != nil {
		return err
	}

	fn, err := fc.getCacheFileName(key)
	if err != nil {
		return err
	}
	return FilePutContents(fn, data)
}

// Delete file cache value.
func (fc *FileCache) Delete(ctx context.Context, key string) error {
	filename, err := fc.getCacheFileName(key)
	if err != nil {
		return err
	}
	if ok, _ := exists(filename); ok {
		err = os.Remove(filename)
		if err != nil {
			return berror.Wrapf(err, DeleteFileCacheItemFailed,
				"can not delete this file cache key-value, key is %s and file name is %s", key, filename)
		}
	}
	return nil
}

// Incr increases cached int value.
// fc value is saved forever unless deleted.
func (fc *FileCache) Incr(ctx context.Context, key string) error {
	data, err := fc.Get(context.Background(), key)
	if err != nil {
		return err
	}

	val, err := incr(data)
	if err != nil {
		return err
	}

	return fc.Put(context.Background(), key, val, time.Duration(fc.EmbedExpiry))
}

// Decr decreases cached int value.
func (fc *FileCache) Decr(ctx context.Context, key string) error {
	data, err := fc.Get(context.Background(), key)
	if err != nil {
		return err
	}

	val, err := decr(data)
	if err != nil {
		return err
	}

	return fc.Put(context.Background(), key, val, time.Duration(fc.EmbedExpiry))
}

// IsExist checks if value exists.
func (fc *FileCache) IsExist(ctx context.Context, key string) (bool, error) {
	fn, err := fc.getCacheFileName(key)
	if err != nil {
		return false, err
	}
	return exists(fn)
}

// ClearAll cleans cached files (not implemented)
func (fc *FileCache) ClearAll(context.Context) error {
	return nil
}

// Check if a file exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, berror.Wrapf(err, InvalidFileCachePath, "file cache path is invalid: %s", path)
}

// FileGetContents Reads bytes from a file.
// if non-existent, create this file.
func FileGetContents(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, berror.Wrapf(err, ReadFileCacheContentFailed,
			"could not read the data from the file: %s, " +
			"please confirm that file exist and Beego has the permission to read the content.", filename)
	}
	return data, nil
}

// FilePutContents puts bytes into a file.
// if non-existent, create this file.
func FilePutContents(filename string, content []byte) error {
	return ioutil.WriteFile(filename, content, os.ModePerm)
}

// GobEncode Gob encodes a file cache item.
func GobEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, berror.Wrap(err, GobEncodeDataFailed, "could not encode this data")
	}
	return buf.Bytes(), nil
}

// GobDecode Gob decodes a file cache item.
func GobDecode(data []byte, to *FileCacheItem) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&to)
	if err != nil {
		return berror.Wrap(err, InvalidGobEncodedData,
			"could not decode this data to FileCacheItem. Make sure that the data is encoded by GOB.")
	}
	return nil
}

func init() {
	Register("file", NewFileCache)
}
