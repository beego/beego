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
	"math/rand"
	"time"
)

// RandomExpireCacheOption implement genreate random time offset expired option
type RandomExpireCacheOption func(*RandomExpireCache)

// RandomExpireCache prevent cache batch invalidation
// Cache random time offset expired
type RandomExpireCache struct {
	Cache
	offset func() time.Duration
}

// Put random time offset expired
func (rec *RandomExpireCache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	timeout += rec.offset()
	return rec.Cache.Put(ctx, key, val, timeout)
}

// NewRandomExpireCache return random expire cache struct
func NewRandomExpireCache(adapter Cache, opts ...RandomExpireCacheOption) Cache {
	var rec RandomExpireCache
	rec.Cache = adapter
	for _, fn := range opts {
		fn(&rec)
	}
	return &rec
}

// defaultExpiredFunc genreate random time offset expired
func defaultExpiredFunc() time.Duration {
	return time.Duration(rand.Intn(5)+3) * time.Second
}
