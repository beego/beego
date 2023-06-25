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
	"errors"
	"fmt"

	"github.com/beego/beego/v2/core/berror"
)

type WriteDeleteCache struct {
	Cache
	storeFunc func(ctx context.Context, key string, val any) error
}

// NewWriteDeleteCache creates a write delete cache pattern decorator.
// The fn is the function that persistent the key and val.
func NewWriteDeleteCache(cache Cache, fn func(ctx context.Context, key string, val any) error) (*WriteDeleteCache, error) {
	if fn == nil || cache == nil {
		return nil, berror.Error(InvalidInitParameters, "cache or storeFunc can not be nil")
	}

	w := &WriteDeleteCache{
		Cache:     cache,
		storeFunc: fn,
	}
	return w, nil
}

func (w *WriteDeleteCache) Set(ctx context.Context, key string, val any) error {
	err := w.storeFunc(ctx, key, val)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return berror.Wrap(err, PersistCacheFailed, fmt.Sprintf("key: %s, val: %v", key, val))
	}
	return w.Cache.Delete(ctx, key)
}
