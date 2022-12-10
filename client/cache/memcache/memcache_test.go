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
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"os"
	"strconv"
	"strings"
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

	assert.Nil(t, bm.Put(context.Background(), "astaxie", "author", timeoutDuration))
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
	assert.True(t, strings.Contains(err.Error(), "key not exist"))

	assert.Nil(t, bm.ClearAll(context.Background()))
	// test clear all
}

func TestReadThroughCache_Memcache_Get(t *testing.T) {
	bm, err := cache.NewCache("memcache", fmt.Sprintf(`{"conn": "%s"}`, "127.0.0.1:11211"))
	assert.Nil(t, err)

	testReadThroughCacheGet(t, bm)

	testReadThroughCacheGetMulti(t, bm)

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
					err, cache.KeyNotExist, "cache unable to load data")
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

func testReadThroughCacheGetMulti(t *testing.T, bm cache.Cache) {
	testCases := []struct {
		name    string
		keys    []string
		values  []any
		cache   cache.Cache
		wantErr error
	}{
		{
			name: "GetMulti load err",
			keys: []string{"key0", "key01"},
			cache: func() cache.Cache {
				db := &MockOrm{}
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
				keysErr := make([]string, 0)
				err1 := berror.Wrap(
					errors.New("the key not exist"),
					cache.KeyNotExist, "cache unable to load data")
				err2 := berror.Wrap(
					errors.New("the key not exist"),
					cache.KeyNotExist, "cache unable to load data")
				keys := []string{"key0", "key01"}
				keyErrMap := map[string]error{"key0": err1, "key01": err2}
				for _, ki := range keys {
					err := keyErrMap[ki]
					keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, err.Error()))
				}
				return berror.Error(cache.MultiGetFailed, strings.Join(keysErr, "; "))

			}(),
		},
		{
			name:   "GetMulti cache exist",
			keys:   []string{"key1", "key2"},
			values: []any{[]byte("value1"), []byte("value2")},
			cache: func() cache.Cache {
				keysMap := map[string]int{"key1": 1, "key2": 1}
				kvs := map[string]any{"key1": "value1", "key2": "value2"}
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
				for key, value := range kvs {
					err = c.Put(context.Background(), key, value, 3*time.Second)
					assert.Nil(t, err)
				}
				return c
			}(),
		},
		{
			name:   "GetMulti loadFunc exist",
			keys:   []string{"key3", "key4"},
			values: []any{[]byte("value3"), []byte("value4")},
			cache: func() cache.Cache {
				keysMap := map[string]int{"key3": 1, "key4": 1}
				kvs := map[string]any{"key3": "value3", "key4": "value4"}
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
			c := tc.cache
			val, err := c.GetMulti(context.Background(), tc.keys)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.EqualValues(t, tc.values, val)
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
