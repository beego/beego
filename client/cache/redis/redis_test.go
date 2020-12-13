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

	"github.com/beego/beego/client/cache"
)

func TestRedisCache(t *testing.T) {

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	bm, err := cache.NewCache("redis", fmt.Sprintf(`{"conn": "%s"}`, redisAddr))
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put(context.Background(), "astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "astaxie"); !res {
		t.Error("check err")
	}

	time.Sleep(11 * time.Second)

	if res, _ := bm.IsExist(context.Background(), "astaxie"); res {
		t.Error("check err")
	}
	if err = bm.Put(context.Background(), "astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	val, _ := bm.Get(context.Background(), "astaxie")
	if v, _ := redis.Int(val, err); v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr(context.Background(), "astaxie"); err != nil {
		t.Error("Incr Error", err)
	}
	val, _ = bm.Get(context.Background(), "astaxie")
	if v, _ := redis.Int(val, err); v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr(context.Background(), "astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	val, _ = bm.Get(context.Background(), "astaxie")
	if v, _ := redis.Int(val, err); v != 1 {
		t.Error("get err")
	}
	bm.Delete(context.Background(), "astaxie")
	if res, _ := bm.IsExist(context.Background(), "astaxie"); res {
		t.Error("delete err")
	}

	// test string
	if err = bm.Put(context.Background(), "astaxie", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "astaxie"); !res {
		t.Error("check err")
	}

	val, _ = bm.Get(context.Background(), "astaxie")
	if v, _ := redis.String(val, err); v != "author" {
		t.Error("get err")
	}

	// test GetMulti
	if err = bm.Put(context.Background(), "astaxie1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "astaxie1"); !res {
		t.Error("check err")
	}

	vv, _ := bm.GetMulti(context.Background(), []string{"astaxie", "astaxie1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if v, _ := redis.String(vv[0], nil); v != "author" {
		t.Error("GetMulti ERROR")
	}
	if v, _ := redis.String(vv[1], nil); v != "author1" {
		t.Error("GetMulti ERROR")
	}

	vv, _ = bm.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	if vv[0] != nil {
		t.Error("GetMulti ERROR")
	}
	if v, _ := redis.String(vv[1], nil); v != "author1" {
		t.Error("GetMulti ERROR")
	}

	// test clear all
	if err = bm.ClearAll(context.Background()); err != nil {
		t.Error("clear all err")
	}
}

func TestCache_Scan(t *testing.T) {
	timeoutDuration := 10 * time.Second

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	// init
	bm, err := cache.NewCache("redis", fmt.Sprintf(`{"conn": "%s"}`, addr))
	if err != nil {
		t.Error("init err")
	}
	// insert all
	for i := 0; i < 100; i++ {
		if err = bm.Put(context.Background(), fmt.Sprintf("astaxie%d", i), fmt.Sprintf("author%d", i), timeoutDuration); err != nil {
			t.Error("set Error", err)
		}
	}
	time.Sleep(time.Second)
	// scan all for the first time
	keys, err := bm.(*Cache).Scan(DefaultKey + ":*")
	if err != nil {
		t.Error("scan Error", err)
	}

	assert.Equal(t, 100, len(keys), "scan all error")

	// clear all
	if err = bm.ClearAll(context.Background()); err != nil {
		t.Error("clear all err")
	}

	// scan all for the second time
	keys, err = bm.(*Cache).Scan(DefaultKey + ":*")
	if err != nil {
		t.Error("scan Error", err)
	}
	if len(keys) != 0 {
		t.Error("scan all err")
	}
}
