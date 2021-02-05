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
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/cache"
)

func TestMemcacheCache(t *testing.T) {
	addr := os.Getenv("MEMCACHE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11211"
	}

	bm, err := cache.NewCache("memcache", fmt.Sprintf(`{"conn": "%s"}`, addr))
	assert.Nil(t, err)

	timeoutDuration := 10 * time.Second

	assert.Nil(t, bm.Put(context.Background(), "astaxie", "1", timeoutDuration))
	res, _ := bm.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	time.Sleep(11 * time.Second)

	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	assert.Nil(t, bm.Put(context.Background(), "astaxie", "1", timeoutDuration))

	val, _ := bm.Get(context.Background(), "astaxie")
	v, err := strconv.Atoi(string(val.([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	assert.Nil(t, bm.Incr(context.Background(), "astaxie"))

	val, _ = bm.Get(context.Background(), "astaxie")
	v, err = strconv.Atoi(string(val.([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	assert.Nil(t, bm.Decr(context.Background(), "astaxie"))

	val, _ = bm.Get(context.Background(), "astaxie")
	v, err = strconv.Atoi(string(val.([]byte)))
	assert.Nil(t, err)
	assert.Equal(t, 1, v)
	bm.Delete(context.Background(), "astaxie")

	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	assert.Nil(t,bm.Put(context.Background(), "astaxie", "author", timeoutDuration) )
	// test string
	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	val, _ = bm.Get(context.Background(), "astaxie")
	vs := val.([]byte)
	assert.Equal(t, "author", string(vs))

	// test GetMulti
	assert.Nil(t, bm.Put(context.Background(), "astaxie1", "author1", timeoutDuration))


	res, _ = bm.IsExist(context.Background(), "astaxie1")
	assert.True(t, res)

	vv, _ := bm.GetMulti(context.Background(), []string{"astaxie", "astaxie1"})
	assert.Equal(t, 2, len(vv))

	if string(vv[0].([]byte)) != "author" && string(vv[0].([]byte)) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if string(vv[1].([]byte)) != "author1" && string(vv[1].([]byte)) != "author" {
		t.Error("GetMulti ERROR")
	}

	vv, err = bm.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Nil(t, vv[0])

	assert.Equal(t, "author1", string(vv[1].([]byte)))

	assert.NotNil(t, err)
	assert.Equal(t, "key [astaxie0] error: key not exist", err.Error())

	assert.Nil(t, bm.ClearAll(context.Background()))
	// test clear all
}
