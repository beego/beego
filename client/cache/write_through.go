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
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"time"
)

type WriteThoughCache struct {
	Cache
	storeFunc func(ctx context.Context, key string, val any) error
}

func NewWriteThoughCache(cache Cache, fn func(ctx context.Context, key string, val any) error) (*WriteThoughCache, error) {
	if fn == nil || cache == nil {
		return nil, berror.Error(InvalidInitParameters, "cache or storeFunc can not be nil")
	}

	w := &WriteThoughCache{
		Cache:     cache,
		storeFunc: fn,
	}
	return w, nil
}

func (w *WriteThoughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.storeFunc(ctx, key, val)
	if err != nil {
		return berror.Wrap(err, PersistCacheFailed, fmt.Sprintf("key: %s, val: %v", key, val))
	}
	return w.Cache.Put(ctx, key, val, expiration)
}
