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
	"context"
	"math"
	"os"
	"sync"
	"testing"
	"time"
)

func TestCacheIncr(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	// timeoutDuration := 10 * time.Second

	bm.Put(context.Background(), "edwardhey", 0, time.Second*20)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			bm.Incr(context.Background(), "edwardhey")
		}()
	}
	wg.Wait()
	val, _ := bm.Get(context.Background(), "edwardhey")
	if val.(int) != 10 {
		t.Error("Incr err")
	}
}

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":1}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 5 * time.Second
	if err = bm.Put(context.Background(), "astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "astaxie"); !res {
		t.Error("check err")
	}

	if v, _ := bm.Get(context.Background(), "astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(7 * time.Second)

	if res, _ := bm.IsExist(context.Background(), "astaxie"); res {
		t.Error("check err")
	}

	if err = bm.Put(context.Background(), "astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	// test different integer type for incr & decr
	testMultiTypeIncrDecr(t, bm, timeoutDuration)

	// test overflow of incr&decr
	testIncrOverFlow(t, bm, timeoutDuration)
	testDecrOverFlow(t, bm, timeoutDuration)

	bm.Delete(context.Background(), "astaxie")
	if res, _ := bm.IsExist(context.Background(), "astaxie"); res {
		t.Error("delete err")
	}

	// test GetMulti
	if err = bm.Put(context.Background(), "astaxie", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "astaxie"); !res {
		t.Error("check err")
	}
	if v, _ := bm.Get(context.Background(), "astaxie"); v.(string) != "author" {
		t.Error("get err")
	}

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
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}

	vv, err = bm.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0] != nil {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if err != nil && err.Error() != "key [astaxie0] error: the key isn't exist" {
		t.Error("GetMulti ERROR")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)
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

	if v, _ := bm.Get(context.Background(), "astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	// test different integer type for incr & decr
	testMultiTypeIncrDecr(t, bm, timeoutDuration)

	// test overflow of incr&decr
	testIncrOverFlow(t, bm, timeoutDuration)
	testDecrOverFlow(t, bm, timeoutDuration)

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
	if v, _ := bm.Get(context.Background(), "astaxie"); v.(string) != "author" {
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
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}

	vv, err = bm.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0] != nil {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if err == nil {
		t.Error("GetMulti ERROR")
	}

	os.RemoveAll("cache")
}

func testMultiTypeIncrDecr(t *testing.T, c Cache, timeout time.Duration) {
	testIncrDecr(t, c, 1, 2, timeout)
	testIncrDecr(t, c, int32(1), int32(2), timeout)
	testIncrDecr(t, c, int64(1), int64(2), timeout)
	testIncrDecr(t, c, uint(1), uint(2), timeout)
	testIncrDecr(t, c, uint32(1), uint32(2), timeout)
	testIncrDecr(t, c, uint64(1), uint64(2), timeout)
}

func testIncrDecr(t *testing.T, c Cache, beforeIncr interface{}, afterIncr interface{}, timeout time.Duration) {
	var err error
	ctx := context.Background()
	key := "incDecKey"
	if err = c.Put(ctx, key, beforeIncr, timeout); err != nil {
		t.Error("Get Error", err)
	}

	if err = c.Incr(ctx, key); err != nil {
		t.Error("Incr Error", err)
	}

	if v, _ := c.Get(ctx, key); v != afterIncr {
		t.Error("Get Error")
	}

	if err = c.Decr(ctx, key); err != nil {
		t.Error("Decr Error", err)
	}

	if v, _ := c.Get(ctx, key); v != beforeIncr {
		t.Error("Get Error")
	}

	if err := c.Delete(ctx, key); err != nil {
		t.Error("Delete Error")
	}
}

func testIncrOverFlow(t *testing.T, c Cache, timeout time.Duration) {
	var err error
	ctx := context.Background()
	key := "incKey"

	// int64
	if err = c.Put(ctx, key, int64(math.MaxInt64), timeout); err != nil {
		t.Error("Put Error: ", err.Error())
		return
	}
	defer func() {
		if err = c.Delete(ctx, key); err != nil {
			t.Errorf("Delete error: %s", err.Error())
		}
	}()
	if err = c.Incr(ctx, key); err == nil {
		t.Error("Incr error")
		return
	}
}

func testDecrOverFlow(t *testing.T, c Cache, timeout time.Duration) {
	var err error
	ctx := context.Background()
	key := "decKey"

	// int64
	if err = c.Put(ctx, key, int64(math.MinInt64), timeout); err != nil {
		t.Error("Put Error: ", err.Error())
		return
	}
	defer func() {
		if err = c.Delete(ctx, key); err != nil {
			t.Errorf("Delete error: %s", err.Error())
		}
	}()
	if err = c.Decr(ctx, key); err == nil {
		t.Error("Decr error")
		return
	}
}
