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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheIncr(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	assert.Nil(t, err)
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
	assert.Nil(t, err)
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
	res, _ := bm.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	assert.Nil(t, bm.Put(context.Background(), "astaxie", "author", timeoutDuration))

	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	v, _ := bm.Get(context.Background(), "astaxie")
	assert.Equal(t, "author", v)

	assert.Nil(t, bm.Put(context.Background(), "astaxie1", "author1", timeoutDuration))

	res, _ = bm.IsExist(context.Background(), "astaxie1")
	assert.True(t, res)

	vv, _ := bm.GetMulti(context.Background(), []string{"astaxie", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Equal(t, "author", vv[0])
	assert.Equal(t, "author1", vv[1])

	vv, err = bm.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Nil(t, vv[0])
	assert.Equal(t, "author1", vv[1])

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "key isn't exist"))
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)
	assert.Nil(t, err)
	timeoutDuration := 10 * time.Second
	assert.Nil(t, bm.Put(context.Background(), "astaxie", 1, timeoutDuration))

	res, _ := bm.IsExist(context.Background(), "astaxie")
	assert.True(t, res)
	v, _ := bm.Get(context.Background(), "astaxie")
	assert.Equal(t, 1, v)

	// test different integer type for incr & decr
	testMultiTypeIncrDecr(t, bm, timeoutDuration)

	// test overflow of incr&decr
	testIncrOverFlow(t, bm, timeoutDuration)
	testDecrOverFlow(t, bm, timeoutDuration)

	bm.Delete(context.Background(), "astaxie")
	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	// test string
	assert.Nil(t, bm.Put(context.Background(), "astaxie", "author", timeoutDuration))
	res, _ = bm.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	v, _ = bm.Get(context.Background(), "astaxie")
	assert.Equal(t, "author", v)

	// test GetMulti
	assert.Nil(t, bm.Put(context.Background(), "astaxie1", "author1", timeoutDuration))

	res, _ = bm.IsExist(context.Background(), "astaxie1")
	assert.True(t, res)

	vv, _ := bm.GetMulti(context.Background(), []string{"astaxie", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Equal(t, "author", vv[0])
	assert.Equal(t, "author1", vv[1])

	vv, err = bm.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	assert.Equal(t, 2, len(vv))

	assert.Nil(t, vv[0])

	assert.Equal(t, "author1", vv[1])
	assert.NotNil(t, err)
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
	ctx := context.Background()
	key := "incDecKey"

	assert.Nil(t, c.Put(ctx, key, beforeIncr, timeout))
	assert.Nil(t, c.Incr(ctx, key))

	v, _ := c.Get(ctx, key)
	assert.Equal(t, afterIncr, v)

	assert.Nil(t, c.Decr(ctx, key))

	v, _ = c.Get(ctx, key)
	assert.Equal(t, v, beforeIncr)
	assert.Nil(t, c.Delete(ctx, key))
}

func testIncrOverFlow(t *testing.T, c Cache, timeout time.Duration) {
	ctx := context.Background()
	key := "incKey"

	assert.Nil(t, c.Put(ctx, key, int64(math.MaxInt64), timeout))
	// int64
	defer func() {
		assert.Nil(t, c.Delete(ctx, key))
	}()
	assert.NotNil(t, c.Incr(ctx, key))
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

func TestRandomExpireCache(t *testing.T) {

	bm, err := NewCache("memory", `{"interval":20}`)
	assert.Nil(t, err)

	cache := NewRandomExpireCache(bm)

	timeoutDuration := 5 * time.Second

	if err = cache.Put(context.Background(), "Leon Ding", 22, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if res, _ := bm.IsExist(context.Background(), "Leon Ding"); !res {
		t.Error("check err")
	}

	if v, _ := bm.Get(context.Background(), "Leon Ding"); v.(int) != 22 {
		t.Error("get err")
	}
}
