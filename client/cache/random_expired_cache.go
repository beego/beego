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
	"sync/atomic"
	"time"
)

// RandomExpireCacheOption implement genreate random time offset expired option
type RandomExpireCacheOption func(*RandomExpireCache)

// RandomExpireCache prevent cache batch invalidation
// Cache random time offset expired
type RandomExpireCache struct {
	Cache
	Offset func() time.Duration
}

// Put random time offset expired
func (rec *RandomExpireCache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	timeout += rec.Offset()
	return rec.Cache.Put(ctx, key, val, timeout)
}

// NewRandomExpireCache return random expire cache struct
func NewRandomExpireCache(adapter Cache, opts ...RandomExpireCacheOption) Cache {
	var rec RandomExpireCache
	rec.Cache = adapter
	for _, fn := range opts {
		fn(&rec)
	}
	if rec.Offset == nil {
		rec.Offset = defaultExpiredFunc()
	}
	return &rec
}

// defaultExpiredFunc return a func that used to generate random time offset (range: [3s,8s)) expired
func defaultExpiredFunc() func() time.Duration {
	const size = 5
	var randTimes [size]time.Duration
	for i := range randTimes {
		randTimes[i] = time.Duration(i + 3)
	}
	// shuffle values
	for i := range randTimes {
		n := rand.Intn(size)
		randTimes[i], randTimes[n] = randTimes[n], randTimes[i]
	}
	var i uint64
	return func() time.Duration {
		return randTimes[atomic.AddUint64(&i, 1)%size]
	}
}
