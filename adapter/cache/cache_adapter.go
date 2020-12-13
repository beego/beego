// Copyright 2020
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/beego/beego/client/cache"
)

type newToOldCacheAdapter struct {
	delegate cache.Cache
}

func (c *newToOldCacheAdapter) Get(key string) interface{} {
	res, _ := c.delegate.Get(context.Background(), key)
	return res
}

func (c *newToOldCacheAdapter) GetMulti(keys []string) []interface{} {
	res, _ := c.delegate.GetMulti(context.Background(), keys)
	return res
}

func (c *newToOldCacheAdapter) Put(key string, val interface{}, timeout time.Duration) error {
	return c.delegate.Put(context.Background(), key, val, timeout)
}

func (c *newToOldCacheAdapter) Delete(key string) error {
	return c.delegate.Delete(context.Background(), key)
}

func (c *newToOldCacheAdapter) Incr(key string) error {
	return c.delegate.Incr(context.Background(), key)
}

func (c *newToOldCacheAdapter) Decr(key string) error {
	return c.delegate.Decr(context.Background(), key)
}

func (c *newToOldCacheAdapter) IsExist(key string) bool {
	res, err := c.delegate.IsExist(context.Background(), key)
	return res && err == nil
}

func (c *newToOldCacheAdapter) ClearAll() error {
	return c.delegate.ClearAll(context.Background())
}

func (c *newToOldCacheAdapter) StartAndGC(config string) error {
	return c.delegate.StartAndGC(config)
}

func CreateNewToOldCacheAdapter(delegate cache.Cache) Cache {
	return &newToOldCacheAdapter{
		delegate: delegate,
	}
}

type oldToNewCacheAdapter struct {
	old Cache
}

func (o *oldToNewCacheAdapter) Get(ctx context.Context, key string) (interface{}, error) {
	return o.old.Get(key), nil
}

func (o *oldToNewCacheAdapter) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	return o.old.GetMulti(keys), nil
}

func (o *oldToNewCacheAdapter) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return o.old.Put(key, val, timeout)
}

func (o *oldToNewCacheAdapter) Delete(ctx context.Context, key string) error {
	return o.old.Delete(key)
}

func (o *oldToNewCacheAdapter) Incr(ctx context.Context, key string) error {
	return o.old.Incr(key)
}

func (o *oldToNewCacheAdapter) Decr(ctx context.Context, key string) error {
	return o.old.Decr(key)
}

func (o *oldToNewCacheAdapter) IsExist(ctx context.Context, key string) (bool, error) {
	return o.old.IsExist(key), nil
}

func (o *oldToNewCacheAdapter) ClearAll(ctx context.Context) error {
	return o.old.ClearAll()
}

func (o *oldToNewCacheAdapter) StartAndGC(config string) error {
	return o.old.StartAndGC(config)
}

func CreateOldToNewAdapter(old Cache) cache.Cache {
	return &oldToNewCacheAdapter{
		old: old,
	}
}
