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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomExpireCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	assert.Nil(t, err)

	// cache := NewRandomExpireCache(bm)
	cache := NewRandomExpireCache(bm, func(opt *RandomExpireCache) {
		opt.Offset = defaultExpiredFunc
	})

	timeoutDuration := 3 * time.Second

	if err = cache.Put(context.Background(), "Leon Ding", 22, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if res, _ := cache.IsExist(context.Background(), "Leon Ding"); !res {
		t.Error("check err")
	}

	if v, _ := cache.Get(context.Background(), "Leon Ding"); v.(int) != 22 {
		t.Error("get err")
	}

	cache.Delete(context.Background(), "Leon Ding")
	res, _ := cache.IsExist(context.Background(), "Leon Ding")
	assert.False(t, res)

	assert.Nil(t, cache.Put(context.Background(), "Leon Ding", "author", timeoutDuration))

	cache.Delete(context.Background(), "astaxie")
	res, _ = cache.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	assert.Nil(t, cache.Put(context.Background(), "astaxie", "author", timeoutDuration))

	res, _ = cache.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	v, _ := cache.Get(context.Background(), "astaxie")
	assert.Equal(t, "author", v)

	assert.Nil(t, cache.Put(context.Background(), "astaxie1", "author1", timeoutDuration))

	res, _ = cache.IsExist(context.Background(), "astaxie1")
	assert.True(t, res)

	vv, _ := cache.GetMulti(context.Background(), []string{"astaxie", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Equal(t, "author", vv[0])
	assert.Equal(t, "author1", vv[1])

	vv, err = cache.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Nil(t, vv[0])
	assert.Equal(t, "author1", vv[1])

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "key isn't exist"))

}
