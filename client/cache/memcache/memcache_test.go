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

	"github.com/astaxie/beego/client/cache"
)

func TestMemcacheCache(t *testing.T) {

	addr := os.Getenv("MEMCACHE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11211"
	}

	bm, err := cache.NewCache("memcache", fmt.Sprintf(`{"conn": "%s"}`, addr))
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put(context.Background(), "astaxie", "1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "astaxie"); !res {
		t.Error("check err")
	}

	time.Sleep(11 * time.Second)

	if res, _ := bm.IsExist(context.Background(), "astaxie"); res {
		t.Error("check err")
	}
	if err = bm.Put(context.Background(), "astaxie", "1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	val, _ := bm.Get(context.Background(), "astaxie")
	if v, err := strconv.Atoi(string(val.([]byte))); err != nil || v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr(context.Background(), "astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	val, _ = bm.Get(context.Background(), "astaxie")
	if v, err := strconv.Atoi(string(val.([]byte))); err != nil || v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr(context.Background(), "astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	val, _ = bm.Get(context.Background(), "astaxie")
	if v, err := strconv.Atoi(string(val.([]byte))); err != nil || v != 1 {
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
	if v := val.([]byte); string(v) != "author" {
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
	if string(vv[0].([]byte)) != "author" && string(vv[0].([]byte)) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if string(vv[1].([]byte)) != "author1" && string(vv[1].([]byte)) != "author" {
		t.Error("GetMulti ERROR")
	}

	// test clear all
	if err = bm.ClearAll(context.Background()); err != nil {
		t.Error("clear all err")
	}
}
