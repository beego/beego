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
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	initError = "init err"
	setError = "set Error"
	checkError = "check err"
	getError = "get err"
	getMultiError = "GetMulti Error"
)

func TestCacheIncr(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error(initError)
	}
	// timeoutDuration := 10 * time.Second

	bm.Put("edwardhey", 0, time.Second*20)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			bm.Incr("edwardhey")
		}()
	}
	wg.Wait()
	if bm.Get("edwardhey").(int) != 10 {
		t.Error("Incr err")
	}
}

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)

	assert.Nil(t, err)

	timeoutDuration := 5 * time.Second
	err = bm.Put("astaxie", 1, timeoutDuration)
	assert.Nil(t, err)

	assert.True(t, bm.IsExist("astaxie"))

	assert.Equal(t, 1, bm.Get("astaxie"))

	time.Sleep(10 * time.Second)

	assert.False(t, bm.IsExist("astaxie"))

	err = bm.Put("astaxie", 1, timeoutDuration)
	assert.Nil(t, err)

	err = bm.Incr("astaxie")
	assert.Nil(t, err)

	assert.Equal(t, 2, bm.Get("astaxie"))

	assert.Nil(t, bm.Decr("astaxie"))

	assert.Equal(t, 1, bm.Get("astaxie"))

	assert.Nil(t, bm.Delete("astaxie"))

	assert.False(t, bm.IsExist("astaxie"))

	assert.Nil(t, bm.Put("astaxie", "author", timeoutDuration))

	assert.True(t, bm.IsExist("astaxie"))

	assert.Equal(t, "author", bm.Get("astaxie"))

	assert.Nil(t, bm.Put("astaxie1", "author1", timeoutDuration))

	assert.True(t, bm.IsExist("astaxie1"))

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})

	assert.Equal(t, 2, len(vv))


	assert.Equal(t, "author", vv[0])

	assert.Equal(t, "author1", vv[1])
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)

	assert.Nil(t, err)
	timeoutDuration := 5 * time.Second

	assert.Nil(t, bm.Put("astaxie", 1, timeoutDuration))

	assert.True(t, bm.IsExist("astaxie"))

	assert.Equal(t, 1, bm.Get("astaxie"))

	assert.Nil(t, bm.Incr("astaxie"))

	assert.Equal(t, 2, bm.Get("astaxie"))

	assert.Nil(t, bm.Decr("astaxie"))

	assert.Equal(t, 1, bm.Get("astaxie"))
	assert.Nil(t, bm.Delete("astaxie"))

	assert.False(t, bm.IsExist("astaxie"))

	assert.Nil(t, bm.Put("astaxie", "author", timeoutDuration))

	assert.True(t, bm.IsExist("astaxie"))

	assert.Equal(t, "author", bm.Get("astaxie"))

	assert.Nil(t, bm.Put("astaxie1", "author1", timeoutDuration))

	assert.True(t, bm.IsExist("astaxie1"))

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})

	assert.Equal(t, 2, len(vv))

	assert.Equal(t, "author", vv[0])
	assert.Equal(t, "author1", vv[1])
	assert.Nil(t, os.RemoveAll("cache"))
}
