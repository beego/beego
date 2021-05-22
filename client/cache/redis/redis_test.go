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
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/cache"
)

func TestRedisCache(t *testing.T) {

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	bm, err := cache.NewCache("redis", fmt.Sprintf(`{"conn": "%s"}`, redisAddr))
	assert.Nil(t, err)
	timeoutDuration := 3 * time.Second

	assert.Nil(t, bm.Put(context.Background(), "astaxie", 1, timeoutDuration))

	res, _ := bm.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	time.Sleep(5 * time.Second)

	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	assert.Nil(t, bm.Put(context.Background(), "astaxie", 1, timeoutDuration))

	val, _ := bm.Get(context.Background(), "astaxie")
	v, _ := redis.Int(val, err)
	assert.Equal(t, 1, v)

	assert.Nil(t, bm.Incr(context.Background(), "astaxie"))
	val, _ = bm.Get(context.Background(), "astaxie")
	v, _ = redis.Int(val, err)
	assert.Equal(t, 2, v)

	assert.Nil(t, bm.Decr(context.Background(), "astaxie"))

	val, _ = bm.Get(context.Background(), "astaxie")
	v, _ = redis.Int(val, err)
	assert.Equal(t, 1, v)
	bm.Delete(context.Background(), "astaxie")

	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	assert.Nil(t, bm.Put(context.Background(), "astaxie", "author", timeoutDuration))
	// test string

	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	val, _ = bm.Get(context.Background(), "astaxie")
	vs, _ := redis.String(val, err)
	assert.Equal(t, "author", vs)

	// test GetMulti
	assert.Nil(t, bm.Put(context.Background(), "astaxie1", "author1", timeoutDuration))

	res, _ = bm.IsExist(context.Background(), "astaxie1")
	assert.True(t, res)

	vv, _ := bm.GetMulti(context.Background(), []string{"astaxie", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	vs, _ = redis.String(vv[0], nil)
	assert.Equal(t, "author", vs)

	vs, _ = redis.String(vv[1], nil)
	assert.Equal(t, "author1", vs)

	vv, _ = bm.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	assert.Nil(t, vv[0])

	vs, _ = redis.String(vv[1], nil)
	assert.Equal(t, "author1", vs)

	// test clear all
	assert.Nil(t, bm.ClearAll(context.Background()))
}

func TestCache_Scan(t *testing.T) {
	timeoutDuration := 10 * time.Second

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	// init
	bm, err := cache.NewCache("redis", fmt.Sprintf(`{"conn": "%s"}`, addr))

	assert.Nil(t, err)
	// insert all
	for i := 0; i < 100; i++ {
		assert.Nil(t, bm.Put(context.Background(), fmt.Sprintf("astaxie%d", i), fmt.Sprintf("author%d", i), timeoutDuration))
	}
	time.Sleep(time.Second)
	// scan all for the first time
	keys, err := bm.(*Cache).Scan(DefaultKey + ":*")
	assert.Nil(t, err)

	assert.Equal(t, 100, len(keys), "scan all error")

	// clear all
	assert.Nil(t, bm.ClearAll(context.Background()))

	// scan all for the second time
	keys, err = bm.(*Cache).Scan(DefaultKey + ":*")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(keys))
}
