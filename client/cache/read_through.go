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
)

// readThroughCache is a decorator
// add the read through function to the original Cache function
type readThroughCache struct {
	Cache
	expiration       time.Duration
	loadFunc         func(ctx context.Context, key string) (any, error)
	applyForGetMulti bool
}

// NewReadThroughCache create readThroughCache
func NewReadThroughCache(cache Cache, expiration time.Duration,
	loadFunc func(ctx context.Context, key string) (any, error), applyForGetMulti bool) (Cache, error) {
	if loadFunc == nil {
		return nil, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil")
	}
	return &readThroughCache{
		Cache:            cache,
		expiration:       expiration,
		loadFunc:         loadFunc,
		applyForGetMulti: applyForGetMulti,
	}, nil
}

// Get will try to call the LoadFunc to load data if the Cache returns value nil or non-nil error.
func (c *readThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := c.Cache.Get(ctx, key)
	if val == nil || err != nil {
		val, err = c.loadFunc(ctx, key)
		if err != nil {
			return nil, berror.Wrap(
				err, LoadFuncFailed, "cache unable to load data")
		}
		err = c.Cache.Put(ctx, key, val, c.expiration)
		if err != nil {
			return val, err
		}
	}
	return val, nil
}
