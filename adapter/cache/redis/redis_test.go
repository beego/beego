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

package redis

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
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

func TestRedisCache(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	bm, err := cache.NewCache("redis", fmt.Sprintf(`{"conn": "%s"}`, redisAddr))
	assert.Nil(t, err)
	timeoutDuration := 5 * time.Second

	assert.Nil(t, bm.Put("astaxie", 1, timeoutDuration))

	assert.True(t, bm.IsExist("astaxie"))

	time.Sleep(7 * time.Second)

	assert.False(t, bm.IsExist("astaxie"))

	assert.Nil(t,  bm.Put("astaxie", 1, timeoutDuration))

	v, err := redis.Int(bm.Get("astaxie"), err)
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	assert.Nil(t, bm.Incr("astaxie"))

	v, err = redis.Int(bm.Get("astaxie"), err)
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	assert.Nil(t, bm.Decr("astaxie"))

	v, err = redis.Int(bm.Get("astaxie"), err)
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	assert.Nil(t, bm.Delete("astaxie"))

	assert.False(t, bm.IsExist("astaxie"))

	assert.Nil(t, bm.Put("astaxie", "author", timeoutDuration))
	assert.True(t, bm.IsExist("astaxie"))

	vs, err := redis.String(bm.Get("astaxie"), err)
	assert.Nil(t, err)
	assert.Equal(t, "author", vs)

	assert.Nil(t, bm.Put("astaxie1", "author1", timeoutDuration))

	assert.False(t, bm.IsExist("astaxie1"))

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})

	assert.Equal(t, 2, len(vv))

	vs, err = redis.String(vv[0], nil)

	assert.Nil(t, err)
	assert.Equal(t, "author", vs)

	vs, err = redis.String(vv[1], nil)

	assert.Nil(t, err)
	assert.Equal(t, "author1", vs)

	assert.Nil(t, bm.ClearAll())
	// test clear all
}

func TestCacheScan(t *testing.T) {
	timeoutDuration := 10 * time.Second
	// init
	bm, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379"}`)
	if err != nil {
		t.Error(initError)
	}
	// insert all
	for i := 0; i < 10000; i++ {
		if err = bm.Put(fmt.Sprintf("astaxie%d", i), fmt.Sprintf("author%d", i), timeoutDuration); err != nil {
			t.Error(setError, err)
		}
	}

	// clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}

}
