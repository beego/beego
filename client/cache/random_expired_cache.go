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

// ExpiredFunc implement genreate random time offset expired
type RandomExpireCacheOptions func(*RandomExpireCache)

// RandomExpireCache prevent cache batch invalidation
// Cache random time offset expired
type RandomExpireCache struct {
	cache  Cache
	Offset func() time.Duration
}

// Put random time offset expired
func (rec *RandomExpireCache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	timeout += rec.Offset()
	return rec.cache.Put(ctx, key, val, timeout)
}

// NewRandomExpireCache return random expire cache struct
func NewRandomExpireCache(adapter Cache, opts ...RandomExpireCacheOptions) Cache {
	var cache RandomExpireCache
	if len(opts) > 0 {
		for _, fn := range opts {
			fn(&cache)
		}
	}
	if cache.Offset == nil {
		cache.Offset = defaultExpiredFunc
	}
	cache.cache = adapter
	return &cache
}

// defaultExpiredFunc genreate random time offset expired
func defaultExpiredFunc() time.Duration {
	offs := (time.Duration(rand.Intn(5)) * time.Second)

	for (offs < offs+(2*time.Second)) && (offs > offs+(8*time.Second)) {
		offs = (time.Duration(rand.Intn(5)) * time.Second)
	}

	return offs
}

// Get get value from memcache.
func (rec *RandomExpireCache) Get(ctx context.Context, key string) (interface{}, error) {
	return rec.cache.Get(ctx, key)
}

// GetMulti gets a value from a key in memcache.
func (rec *RandomExpireCache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	return rec.cache.GetMulti(ctx, keys)
}

// Delete deletes a value in memcache.
func (rec *RandomExpireCache) Delete(ctx context.Context, key string) error {
	return rec.cache.Delete(ctx, key)
}

// Incr increases counter.
func (rec *RandomExpireCache) Incr(ctx context.Context, key string) error {
	return rec.cache.Incr(ctx, key)
}

// Decr decreases counter.
func (rec *RandomExpireCache) Decr(ctx context.Context, key string) error {
	return rec.cache.Decr(ctx, key)
}

// IsExist checks if a value exists in memcache.
func (rec *RandomExpireCache) IsExist(ctx context.Context, key string) (bool, error) {
	return rec.cache.IsExist(ctx, key)
}

// ClearAll clears all cache in memcache.
func (rec *RandomExpireCache) ClearAll(ctx context.Context) error {
	return rec.cache.ClearAll(ctx)
}

// StartAndGC starts the memcache adapter.
// config: must be in the format {"conn":"connection info"}.
// If an error occurs during connecting, an error is returned
func (rec *RandomExpireCache) StartAndGC(config string) error {
	return rec.cache.StartAndGC(config)
}
