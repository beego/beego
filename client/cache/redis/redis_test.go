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
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/berror"
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
	assert.Nil(t, bm.Delete(context.Background(), "astaxie"))

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

func TestReadThroughCache_redis_Get(t *testing.T) {
	bm, err := cache.NewCache("redis", fmt.Sprintf(`{"conn": "%s"}`, "127.0.0.1:6379"))
	assert.Nil(t, err)

	testReadThroughCacheGet(t, bm)
}

func testReadThroughCacheGet(t *testing.T, bm cache.Cache) {
	testCases := []struct {
		name    string
		key     string
		value   string
		cache   cache.Cache
		wantErr error
	}{
		{
			name: "Get load err",
			key:  "key0",
			cache: func() cache.Cache {
				kvs := map[string]any{"key0": "value0"}
				db := &MockOrm{kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					v, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					val := []byte(v.(string))
					return val, nil
				}
				c, err := cache.NewReadThroughCache(bm, 3*time.Second, loadfunc)
				assert.Nil(t, err)
				return c
			}(),
			wantErr: func() error {
				err := errors.New("the key not exist")
				return berror.Wrap(
					err, cache.LoadFuncFailed, "cache unable to load data")
			}(),
		},
		{
			name:  "Get cache exist",
			key:   "key1",
			value: "value1",
			cache: func() cache.Cache {
				keysMap := map[string]int{"key1": 1}
				kvs := map[string]any{"key1": "value1"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					v, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					val := []byte(v.(string))
					return val, nil
				}
				c, err := cache.NewReadThroughCache(bm, 3*time.Second, loadfunc)
				assert.Nil(t, err)
				err = c.Put(context.Background(), "key1", "value1", 3*time.Second)
				assert.Nil(t, err)
				return c
			}(),
		},
		{
			name:  "Get loadFunc exist",
			key:   "key2",
			value: "value2",
			cache: func() cache.Cache {
				keysMap := map[string]int{"key2": 1}
				kvs := map[string]any{"key2": "value2"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					v, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					val := []byte(v.(string))
					return val, nil
				}
				c, err := cache.NewReadThroughCache(bm, 3*time.Second, loadfunc)
				assert.Nil(t, err)
				return c
			}(),
		},
	}
	_, err := cache.NewReadThroughCache(bm, 3*time.Second, nil)
	assert.Equal(t, berror.Error(cache.InvalidLoadFunc, "loadFunc cannot be nil"), err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs := []byte(tc.value)
			c := tc.cache
			val, err := c.Get(context.Background(), tc.key)
			if err != nil {
				assert.EqualError(t, tc.wantErr, err.Error())
				return
			}
			assert.Equal(t, bs, val)
		})
	}
}

type MockOrm struct {
	keysMap map[string]int
	kvs     map[string]any
}

func (m *MockOrm) Load(key string) (any, error) {
	_, ok := m.keysMap[key]
	if !ok {
		return nil, errors.New("the key not exist")
	}
	return m.kvs[key], nil
}

func TestCache_associate(t *testing.T) {
	testCases := []struct {
		name            string
		skipEmptyPrefix bool
		prefix          string
		input           string
		wantRes         string
	}{
		{
			name:            "skip prefix",
			skipEmptyPrefix: true,
			prefix:          "",
			input:           "my-key",
			wantRes:         "my-key",
		},
		{
			name:            "skip prefix but prefix not empty",
			skipEmptyPrefix: true,
			prefix:          "abc",
			input:           "my-key",
			wantRes:         "abc:my-key",
		},
		{
			name:            "using empty prefix",
			skipEmptyPrefix: false,
			prefix:          "",
			input:           "my-key",
			wantRes:         ":my-key",
		},
		{
			name:    "using prefix",
			prefix:  "abc",
			input:   "my-key",
			wantRes: "abc:my-key",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewRedisCache().(*Cache)
			c.skipEmptyPrefix = tc.skipEmptyPrefix
			c.key = tc.prefix
			res := c.associate(tc.input)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestCache_parseConf(t *testing.T) {
	tests := []struct {
		name string

		configStr string

		wantCache Cache
		wantErr   error
	}{
		{
			name: "just conn",
			configStr: `{
  "conn": "127.0.0.1:6379"
}`,

			wantCache: Cache{
				conninfo:        "127.0.0.1:6379",
				dbNum:           0,
				key:             DefaultKey,
				password:        "",
				maxIdle:         defaultMaxIdle,
				skipEmptyPrefix: false,
				timeout:         defaultTimeout,
			},
			wantErr: nil,
		},

		{
			name: "all",
			configStr: `{
  "dbNum": "2",
  "skipEmptyPrefix": "true",
  "key": "mykey",
  "conn": "redis://mypwd@127.0.0.1:6379",
  "maxIdle": "10",
  "timeout": "30s"
}`,

			wantCache: Cache{
				conninfo:        "127.0.0.1:6379",
				dbNum:           2,
				key:             "mykey",
				password:        "mypwd",
				maxIdle:         10,
				skipEmptyPrefix: true,
				timeout:         time.Second * 30,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{}
			err := c.parseConf(tt.configStr)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantCache, c)
		})
	}
}
