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

// nolint
package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/beego/beego/v2/core/berror"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	Db      map[string]any
	loadCnt int64
}

var (
	mockDB    = MockDB{Db: make(map[string]any), loadCnt: 0}
	mockBloom = bloom.NewWithEstimates(20000, 0.99)
	loadFunc  = func(ctx context.Context, key string) (any, error) {
		mockDB.loadCnt += 1 // flag of number load data from db
		v, ok := mockDB.Db[key]
		if !ok {
			return nil, errors.New("fail")
		}
		return v, nil
	}
	cacheUnderlying = NewMemoryCache()
)

func TestBloomFilterCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		wantVal any

		before func()
		after  func()

		wantErrCode uint32
	}{
		// case: keys exist in cache
		// want: not load data from db
		{
			name: "not_load_db",
			before: func() {
				_ = cacheUnderlying.Put(context.Background(), "exist_in_cache", "123", time.Minute)
			},
			key: "exist_in_DB",
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(0))
				_ = cacheUnderlying.Delete(context.Background(), "exist_in_cache")
				mockDB = MockDB{
					Db:      make(map[string]any),
					loadCnt: 0,
				}
			},
		},
		// case: keys not exist in cache, not exist in bloom
		// want: not load data from db
		{
			name: "not_load_db",
			before: func() {
				mockDB.Db = map[string]any{
					"exist_in_DB": "exist_in_DB",
				}
				mockBloom.AddString("other")
			},
			key: "exist_in_DB",
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(0))
				mockBloom.ClearAll()
				mockDB = MockDB{
					Db:      make(map[string]any),
					loadCnt: 0,
				}
			},
		},
		// case: keys not exist in cache, exist in bloom, exist in db,
		// want: load data from db, and set cache
		{
			name: "load_db",
			before: func() {
				mockDB.Db = map[string]any{
					"exist_in_DB": "exist_in_DB",
				}
				mockBloom.AddString("exist_in_DB")
			},
			key:     "exist_in_DB",
			wantVal: "exist_in_DB",
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(1))
				_ = cacheUnderlying.Delete(context.Background(), "exist_in_DB")
				mockBloom.ClearAll()
				mockDB = MockDB{
					Db:      make(map[string]any),
					loadCnt: 0,
				}
			},
		},
		// case: keys not exist in cache, exist in bloom, not exist in db,
		// want: load func error
		{
			name: "load db fail",
			before: func() {
				mockBloom.AddString("not_exist_in_DB")
			},
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(1))
				mockBloom.ClearAll()
				mockDB = MockDB{
					Db:      make(map[string]any),
					loadCnt: 0,
				}
			},
			key:         "not_exist_in_DB",
			wantErrCode: LoadFuncFailed.Code(),
		},
		// case: keys not exist in cache, not exist in bloom, execute Get single key 100 times concurrently
		// want: not load data from db
		{
			name: "Concurrency_Get",
			before: func() {
				mockBloom.AddString("exist_key")
			},
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(1))
				mockBloom.ClearAll()
				mockDB = MockDB{
					Db:      make(map[string]any),
					loadCnt: 0,
				}
			},
			key: "not_exist_in_DB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			bfc, err := NewBloomFilterCache(cacheUnderlying, loadFunc, mockBloom, time.Minute)
			assert.Nil(t, err)

			got, err := bfc.Get(context.Background(), tc.key)
			if tc.wantErrCode != 0 {
				errCode, _ := berror.FromError(err)
				assert.Equal(t, tc.wantErrCode, errCode.Code())
				return
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.wantVal, got)

			cacheVal, _ := bfc.Cache.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantVal, cacheVal)
			tc.after()
		})
	}
}

func TestBloomFilterCache_Get_Concurrency(t *testing.T) {
	bfc, err := NewBloomFilterCache(cacheUnderlying, loadFunc, mockBloom, time.Minute)
	assert.Nil(t, err)

	mockDB.Db = map[string]any{
		"key_11": "value_11",
	}
	mockBloom.AddString("key_11")

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%d", i)
		go func(key string) {
			defer wg.Done()
			val, _ := bfc.Get(context.Background(), key)

			if val != nil {
				assert.Equal(t, "value_11", val)
			}
		}(key)
	}
	wg.Wait()
	assert.Equal(t, int64(1), mockDB.loadCnt)
}
