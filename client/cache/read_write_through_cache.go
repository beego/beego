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
	"strings"
	"time"

	"github.com/beego/beego/v2/core/berror"
)

// readThroughCache is a decorator for the underlying cache,
// which enhances the Get and GetMulti methods of the underlying cache by invoking LoadFunc
type readThroughCache struct {
	Cache
	loadFunc      LoadFunc
	keyExpiration time.Duration
}

type LoadFunc func(ctx context.Context, key string) (any, error)

// NewReadThroughCache returns readThroughCache,
// which enhances the Get and GetMulti methods of the underlying cache by invoking LoadFunc.
func NewReadThroughCache(adapter Cache, loadFunc LoadFunc, keyExpiration time.Duration) Cache {
	return &readThroughCache{
		Cache:         adapter,
		loadFunc:      loadFunc,
		keyExpiration: keyExpiration,
	}
}

// Get will try to call the LoadFunc to load data if the Get method of underlying cache returns non-nil error.
func (r *readThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == nil {
		return val, nil
	}
	val, err = r.loadFunc(ctx, key)
	if err != nil {
		return nil, err
	}
	_ = r.Cache.Put(ctx, key, val, r.keyExpiration)
	return val, nil
}

// GetMulti will try to call the LoadFunc to load data if the GetMulti method of underlying cache returns non-nil error.
// You should check the concrete type of underlying cache to learn the cases that the GetMulti function will return non-nil error.
func (r *readThroughCache) GetMulti(ctx context.Context, keys []string) ([]any, error) {
	values, err := r.Cache.GetMulti(ctx, keys)
	if err == nil {
		return values, nil
	}

	values = make([]any, len(keys))
	keysErrs := make([]string, 0)

	for i, key := range keys {
		value, err := r.Get(ctx, key)
		if err != nil {
			keysErrs = append(keysErrs, fmt.Sprintf("key [%d] error: %s", i, err.Error()))
			continue
		}
		values[i] = value
	}

	if len(keysErrs) != 0 {
		return values, berror.Error(MultiGetFailed, strings.Join(keysErrs, "; "))
	}

	return values, nil
}

// writeThroughCache is a decorator for the underlying cache,
// which enhances the Put method of the underlying cache by invoking StoreFunc
type writeThroughCache struct {
	Cache
	storeFunc     StoreFunc
	keyExpiration time.Duration
}

type StoreFunc func(ctx context.Context, key string, val any) error

// NewWriteThroughCache returns writeThroughCache,
// which enhances the Put method of the underlying cache by invoking StoreFunc.
func NewWriteThroughCache(adapter Cache, storeFunc StoreFunc) Cache {
	return &writeThroughCache{
		Cache:     adapter,
		storeFunc: storeFunc,
	}
}

// Put will try to call the StoreFunc to store data before calling the Put method of underlying cache.
func (w *writeThroughCache) Put(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.storeFunc(ctx, key, val)
	if err != nil {
		return err
	}
	return w.Cache.Put(ctx, key, val, expiration)
}
