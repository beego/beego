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
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
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

func TestCacheScan(t *testing.T) {
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

func TestRedisCache_WriteThough_Set(t *testing.T) {
	bm, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379"}`)
	assert.Nil(t, err)

	var mockDbStore = make(map[string]any)
	testCases := []struct {
		name      string
		storeFunc func(ctx context.Context, key string, val any) error
		key       string
		value     any
		wantErr   error
	}{
		{
			name:    "storeFunc nil",
			wantErr: berror.Error(cache.InvalidStoreFunc, "storeFunc can not be nil"),
		},
		{
			name: "set error",
			storeFunc: func(ctx context.Context, key string, val any) error {
				return errors.New("failed")
			},
			wantErr: berror.Wrap(errors.New("failed"), cache.PersistCacheFailed,
				fmt.Sprintf("key: %s, val: %v", "", nil)),
		},
		{
			name: "memory set success",
			storeFunc: func(ctx context.Context, key string, val any) error {
				mockDbStore[key] = val
				return nil
			},
			key:   "hello",
			value: []byte("world"),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := &cache.WriteThoughCache{
				Cache:     bm,
				StoreFunc: tt.storeFunc,
			}
			err := w.Set(context.Background(), tt.key, tt.value, 60*time.Second)
			if err != nil {
				assert.EqualError(t, tt.wantErr, err.Error())
				return
			}

			val, err := w.Get(context.Background(), tt.key)
			assert.Nil(t, err)
			assert.Equal(t, tt.value, val)

			vv, ok := mockDbStore[tt.key]
			assert.True(t, ok)
			assert.Equal(t, tt.value, vv)
		})
	}
}
