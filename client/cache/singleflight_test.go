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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleflight_Memory_Get(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	assert.Nil(t, err)

	testSingleflightCacheConcurrencyGet(t, bm)
}

func TestSingleflight_file_Get(t *testing.T) {
	fc := NewFileCache().(*FileCache)
	fc.CachePath = "////aaa"
	err := fc.Init()
	assert.NotNil(t, err)
	fc.CachePath = getTestCacheFilePath()
	err = fc.Init()
	assert.Nil(t, err)

	testSingleflightCacheConcurrencyGet(t, fc)
}

func testSingleflightCacheConcurrencyGet(t *testing.T, bm Cache) {
	key, value := "key3", "value3"
	db := &MockOrm{keysMap: map[string]int{key: 1}, kvs: map[string]any{key: value}}
	c, err := NewSingleflightCache(bm, 10*time.Second,
		func(ctx context.Context, key string) (any, error) {
			val, er := db.Load(key)
			if er != nil {
				return nil, er
			}
			return val, nil
		})
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			val, err := c.Get(context.Background(), key)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, value, val)
		}()
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()
}
