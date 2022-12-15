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
	"errors"
	"testing"
	"time"

	"github.com/beego/beego/v2/core/berror"
	"github.com/stretchr/testify/assert"
)

func TestReadThroughCache_Memory_Get(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	assert.Nil(t, err)
	testReadThroughCacheGet(t, bm)
}

func TestReadThroughCache_file_Get(t *testing.T) {
	fc := NewFileCache().(*FileCache)
	fc.CachePath = "////aaa"
	err := fc.Init()
	assert.NotNil(t, err)
	fc.CachePath = getTestCacheFilePath()
	err = fc.Init()
	assert.Nil(t, err)
	testReadThroughCacheGet(t, fc)
}

func testReadThroughCacheGet(t *testing.T, bm Cache) {
	testCases := []struct {
		name    string
		key     string
		value   string
		cache   Cache
		wantErr error
	}{
		{
			name: "Get load err",
			key:  "key0",
			cache: func() Cache {
				kvs := map[string]any{"key0": "value0"}
				db := &MockOrm{kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc)
				assert.Nil(t, err)
				return c
			}(),
			wantErr: func() error {
				err := errors.New("the key not exist")
				return berror.Wrap(
					err, LoadFuncFailed, "cache unable to load data")
			}(),
		},
		{
			name:  "Get cache exist",
			key:   "key1",
			value: "value1",
			cache: func() Cache {
				keysMap := map[string]int{"key1": 1}
				kvs := map[string]any{"key1": "value1"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc)
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
			cache: func() Cache {
				keysMap := map[string]int{"key2": 1}
				kvs := map[string]any{"key2": "value2"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc)
				assert.Nil(t, err)
				return c
			}(),
		},
	}
	_, err := NewReadThroughCache(bm, 3*time.Second, nil)
	assert.Equal(t, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil"), err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.cache
			val, err := c.Get(context.Background(), tc.key)
			if err != nil {
				assert.EqualError(t, tc.wantErr, err.Error())
				return
			}
			assert.Equal(t, tc.value, val)
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
