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
	"time"

	"github.com/beego/beego/v2/core/berror"
	"golang.org/x/sync/singleflight"
)

// SingleflightCache
// This is a very simple decorator mode
type SingleflightCache struct {
	Cache
	group      *singleflight.Group
	Expiration time.Duration
	LoadFunc   func(ctx context.Context, key string) (any, error)
}

// NewSingleflightCache create SingleflightCache
func NewSingleflightCache(c Cache, expiration time.Duration,
	loadFunc func(ctx context.Context, key string) (any, error),
) (Cache, error) {
	if loadFunc == nil {
		return nil, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil")
	}
	return &SingleflightCache{
		Cache:      c,
		group:      &singleflight.Group{},
		Expiration: expiration,
		LoadFunc:   loadFunc,
	}, nil
}

// Get In the Get method, single flight is used to load data and write back the cache.
func (s *SingleflightCache) Get(ctx context.Context, key string) (any, error) {
	val, err := s.Cache.Get(ctx, key)
	if val == nil || err != nil {
		val, err, _ = s.group.Do(key, func() (interface{}, error) {
			v, er := s.LoadFunc(ctx, key)
			if er != nil {
				return nil, berror.Wrap(er, LoadFuncFailed, "cache unable to load data")
			}
			er = s.Cache.Put(ctx, key, v, s.Expiration)
			return v, er
		})
	}
	return val, err
}
