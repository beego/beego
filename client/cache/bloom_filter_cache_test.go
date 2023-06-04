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

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/core/berror"
)

type MockDB struct {
	Db      Cache
	loadCnt int64
}

type BloomFilterMock struct {
	*bloom.BloomFilter
	lock       *sync.RWMutex
	concurrent bool
}

func (b *BloomFilterMock) Add(data string) {
	if b.concurrent {
		b.lock.Lock()
		defer b.lock.Unlock()
	}
	b.BloomFilter.AddString(data)
}

func (b *BloomFilterMock) Test(data string) bool {
	if b.concurrent {
		b.lock.Lock()
		defer b.lock.Unlock()
	}
	return b.BloomFilter.TestString(data)
}

var (
	mockDB    = MockDB{Db: NewMemoryCache(), loadCnt: 0}
	mockBloom = &BloomFilterMock{
		BloomFilter: bloom.NewWithEstimates(20000, 0.01),
		lock:        &sync.RWMutex{},
		concurrent:  false,
	}
	loadFunc = func(ctx context.Context, key string) (any, error) {
		mockDB.loadCnt += 1 // flag of number load data from db
		v, err := mockDB.Db.Get(context.Background(), key)
		if err != nil {
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
				mockDB.loadCnt = 0
				_ = mockDB.Db.ClearAll(context.Background())
			},
		},
		// case: keys not exist in cache, not exist in bloom
		// want: not load data from db
		{
			name: "not_load_db",
			before: func() {
				_ = mockDB.Db.ClearAll(context.Background())
				_ = mockDB.Db.Put(context.Background(), "exist_in_DB", "exist_in_DB", 0)
				mockBloom.AddString("other")
			},
			key: "exist_in_DB",
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(0))
				mockBloom.ClearAll()
				mockDB.loadCnt = 0
				_ = mockDB.Db.ClearAll(context.Background())
			},
		},
		// case: keys not exist in cache, exist in bloom, exist in db,
		// want: load data from db, and set cache
		{
			name: "load_db",
			before: func() {
				_ = mockDB.Db.ClearAll(context.Background())
				_ = mockDB.Db.Put(context.Background(), "exist_in_DB", "exist_in_DB", 0)
				mockBloom.Add("exist_in_DB")
			},
			key:     "exist_in_DB",
			wantVal: "exist_in_DB",
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(1))
				_ = cacheUnderlying.Delete(context.Background(), "exist_in_DB")
				mockBloom.ClearAll()
				mockDB.loadCnt = 0
				_ = mockDB.Db.ClearAll(context.Background())
			},
		},
		// case: keys not exist in cache, exist in bloom, not exist in db,
		// want: load func error
		{
			name: "load db fail",
			before: func() {
				mockBloom.Add("not_exist_in_DB")
			},
			after: func() {
				assert.Equal(t, mockDB.loadCnt, int64(1))
				mockBloom.ClearAll()
				mockDB.loadCnt = 0
				_ = mockDB.Db.ClearAll(context.Background())
			},
			key:         "not_exist_in_DB",
			wantErrCode: LoadFuncFailed.Code(),
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

func ExampleNewBloomFilterCache() {
	c := NewMemoryCache()
	c, err := NewBloomFilterCache(c, func(ctx context.Context, key string) (any, error) {
		return fmt.Sprintf("hello, %s", key), nil
	}, &AlwaysExist{}, time.Minute)
	if err != nil {
		panic(err)
	}

	val, err := c.Get(context.Background(), "Beego")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output:
	// hello, Beego
}

type AlwaysExist struct {
}

func (*AlwaysExist) Test(string) bool {
	return true
}

func (*AlwaysExist) Add(string) {

}
