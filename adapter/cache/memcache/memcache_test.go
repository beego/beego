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

package memcache

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/adapter/cache"
)

const (
	initError = "init err"
	setError = "set Error"
	checkError = "check err"
	getError = "get err"
	getMultiError = "GetMulti Error"
)

func TestMemcacheCache(t *testing.T) {

	addr := os.Getenv("MEMCACHE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11211"
	}

	bm, err := cache.NewCache("memcache", fmt.Sprintf(`{"conn": "%s"}`, addr))
	assert.Nil(t, err)
	timeoutDuration := 5 * time.Second

	assert.Nil(t, bm.Put("astaxie", "1", timeoutDuration))

	assert.True(t, bm.IsExist("astaxie"))

	time.Sleep(11 * time.Second)

	assert.False(t, bm.IsExist("astaxie"))

	assert.Nil(t, bm.Put("astaxie", "1", timeoutDuration))
	v, err := strconv.Atoi(string(bm.Get("astaxie").([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	assert.Nil(t, bm.Incr("astaxie"))

	v, err = strconv.Atoi(string(bm.Get("astaxie").([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	assert.Nil(t, bm.Decr("astaxie"))

	v, err = strconv.Atoi(string(bm.Get("astaxie").([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	assert.Nil(t, bm.Delete("astaxie"))

	assert.False(t,  bm.IsExist("astaxie"))

	assert.Nil(t, bm.Put("astaxie", "author", timeoutDuration))

	assert.True(t, bm.IsExist("astaxie"))

	assert.Equal(t, []byte("author"), bm.Get("astaxie"))

	assert.Nil(t, bm.Put("astaxie1", "author1", timeoutDuration))

	assert.True(t, bm.IsExist("astaxie1"))

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Equal(t, []byte("author"), vv[0])
	assert.Equal(t, []byte("author1"), vv[1])

	assert.Nil(t, bm.ClearAll())
}
