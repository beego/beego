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
	if err != nil {
		t.Error(initError)
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie") {
		t.Error(checkError)
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist("astaxie") {
		t.Error(checkError)
	}
	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error(setError, err)
	}

	if v, _ := redis.Int(bm.Get("astaxie"), err); v != 1 {
		t.Error(getError)
	}

	if err = bm.Incr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v, _ := redis.Int(bm.Get("astaxie"), err); v != 2 {
		t.Error(getError)
	}

	if err = bm.Decr("astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	if v, _ := redis.Int(bm.Get("astaxie"), err); v != 1 {
		t.Error(getError)
	}
	bm.Delete("astaxie")
	if bm.IsExist("astaxie") {
		t.Error("delete err")
	}

	// test string
	if err = bm.Put("astaxie", "author", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie") {
		t.Error(checkError)
	}

	if v, _ := redis.String(bm.Get("astaxie"), err); v != "author" {
		t.Error(getError)
	}

	// test GetMulti
	if err = bm.Put("astaxie1", "author1", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie1") {
		t.Error(checkError)
	}

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})
	if len(vv) != 2 {
		t.Error(getMultiError)
	}
	if v, _ := redis.String(vv[0], nil); v != "author" {
		t.Error(getMultiError)
	}
	if v, _ := redis.String(vv[1], nil); v != "author1" {
		t.Error(getMultiError)
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}

func TestCache_Scan(t *testing.T) {
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
