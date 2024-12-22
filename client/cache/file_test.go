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

package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFileCacheStartAndGC(t *testing.T) {
	fc := NewFileCache().(*FileCache)
	err := fc.StartAndGC(`{`)
	assert.NotNil(t, err)
	err = fc.StartAndGC(`{}`)
	assert.Nil(t, err)
	_, err = fc.getCacheFileName("key1")
	assert.Nil(t, err)

	assert.Equal(t, fc.CachePath, FileCachePath)
	assert.Equal(t, fc.DirectoryLevel, FileCacheDirectoryLevel)
	assert.Equal(t, fc.EmbedExpiry, int(FileCacheEmbedExpiry))
	assert.Equal(t, fc.FileSuffix, FileCacheFileSuffix)

	err = fc.StartAndGC(`{"CachePath":"/cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)
	// could not create dir
	assert.NotNil(t, err)

	str := getTestCacheFilePath()
	err = fc.StartAndGC(fmt.Sprintf(`{"CachePath":"%s","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`, str))
	assert.Nil(t, err)
	assert.Equal(t, fc.CachePath, str)
	assert.Equal(t, fc.DirectoryLevel, 2)
	assert.Equal(t, fc.EmbedExpiry, 0)
	assert.Equal(t, fc.FileSuffix, ".bin")
	_, err = fc.getCacheFileName("key1")
	assert.Nil(t, err)

	err = fc.StartAndGC(fmt.Sprintf(`{"CachePath":"%s","FileSuffix":".bin","DirectoryLevel":"aaa","EmbedExpiry":"0"}`, str))
	assert.NotNil(t, err)

	err = fc.StartAndGC(fmt.Sprintf(`{"CachePath":"%s","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"aaa"}`, str))
	assert.NotNil(t, err)

	_, err = fc.getCacheFileName("key1")
	assert.Nil(t, err)
}

func TestFileCacheInit(t *testing.T) {
	fc := NewFileCache().(*FileCache)
	fc.CachePath = "////aaa"
	err := fc.Init()
	assert.NotNil(t, err)
	fc.CachePath = getTestCacheFilePath()
	err = fc.Init()
	assert.Nil(t, err)
}

func TestFileGetContents(t *testing.T) {
	_, err := FileGetContents("/bin/aaa")
	assert.NotNil(t, err)
	fn := filepath.Join(os.TempDir(), "fileCache.txt")
	f, err := os.Create(fn)
	assert.Nil(t, err)
	_, err = f.WriteString("text")
	assert.Nil(t, err)
	data, err := FileGetContents(fn)
	assert.Nil(t, err)
	assert.Equal(t, "text", string(data))
}

func TestGobEncodeDecode(t *testing.T) {
	_, err := GobEncode(func() {
		fmt.Print("test func")
	})
	assert.NotNil(t, err)
	data, err := GobEncode(&FileCacheItem{
		Data: "hello",
	})
	assert.Nil(t, err)
	err = GobDecode([]byte("wrong data"), &FileCacheItem{})
	assert.NotNil(t, err)
	dci := &FileCacheItem{}
	err = GobDecode(data, dci)
	assert.Nil(t, err)
	assert.Equal(t, "hello", dci.Data)
}

func TestFileCacheDelete(t *testing.T) {
	fc := NewFileCache()
	err := fc.StartAndGC(`{}`)
	assert.Nil(t, err)
	err = fc.Delete(context.Background(), "my-key")
	assert.Nil(t, err)
}

func getTestCacheFilePath() string {
	return filepath.Join(os.TempDir(), "test", "file.txt")
}
