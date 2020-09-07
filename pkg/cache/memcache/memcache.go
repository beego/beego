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

// Package memcache for cache provider
//
// depend on github.com/bradfitz/gomemcache/memcache
//
// go install github.com/bradfitz/gomemcache/memcache
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/cache/memcache"
//   "github.com/astaxie/beego/cache"
// )
//
//  bm, err := cache.NewCache("memcache", `{"conn":"127.0.0.1:11211"}`)
//
//  more docs http://beego.me/docs/module/cache.md
package memcache

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/astaxie/beego/pkg/cache"
	"github.com/bradfitz/gomemcache/memcache"
)

// Cache Memcache adapter.
type Cache struct {
	conn     *memcache.Client
	conninfo []string
}

// NewMemCache creates a new memcache adapter.
func NewMemCache() cache.Cache {
	return &Cache{}
}

func (rc *Cache) client() *memcache.Client {
	if rc.conn == nil {
		rc.conn = memcache.New(rc.conninfo...)
	}
	return rc.conn
}

//Get get value from memcache
func (rc *Cache) Get(key string) (interface{}, error) {
	return rc.GetWithCtx(context.Background(), key)
}

//GetWithCtx get value with context from memcache
func (rc *Cache) GetWithCtx(ctx context.Context, key string) (interface{}, error) {
	item, err := rc.client().Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, err
}

// GetMulti gets a value from a key in memcache.
func (rc *Cache) GetMulti(keys []string) ([]interface{}, error) {
	return rc.GetMultiWithCtx(context.Background(), keys)
}

// GetMultiWithCtx gets a value from a key in memcache.
func (rc *Cache) GetMultiWithCtx(ctx context.Context, keys []string) ([]interface{}, error) {
	var rv []interface{}
	mv, err := rc.client().GetMulti(keys)
	if err != nil {
		for i := 0; i < len(keys); i++ {
			rv = append(rv, nil)
		}
		return rv, err
	}
	for _, v := range mv {
		rv = append(rv, v.Value)
	}
	return rv, nil
}

// Put puts a value into memcache.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	return rc.PutWithCtx(context.Background(), key, val, timeout)
}

// PutWithCtx puts a value into memcache.
func (rc *Cache) PutWithCtx(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Expiration: int32(timeout / time.Second),
	}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return errors.New("val only support string and []byte")
	}
	return rc.client().Set(item)
}

// Delete cached value by key.
func (rc *Cache) Delete(key string) error {
	return rc.DeleteWithCtx(context.Background(), key)
}

// DeleteWithCtx  cached value by key.
func (rc *Cache) DeleteWithCtx(ctx context.Context, key string) error {
	return rc.client().Delete(key)
}

//IncrBy ..
func (rc *Cache) IncrBy(key string, n int) (int, error) {
	return rc.IncrByWithCtx(context.Background(), key, n)
}

//IncrByWithCtx ..
func (rc *Cache) IncrByWithCtx(ctx context.Context, key string, n int) (int, error) {
	if n > 0 {
		val, err := rc.client().Increment(key, uint64(n))
		return int(val), err
	}
	val, err := rc.client().Decrement(key, uint64(-n))
	return int(val), err
}

// Incr ..
func (rc *Cache) Incr(key string) (int, error) {
	return rc.IncrBy(key, 1)
}

//IncrWithCtx ..
func (rc *Cache) IncrWithCtx(ctx context.Context, key string) (int, error) {
	return rc.IncrByWithCtx(ctx, key, 1)
}

// Decr ..
func (rc *Cache) Decr(key string) (int, error) {
	return rc.IncrBy(key, -1)
}

//DecrWithCtx ..
func (rc *Cache) DecrWithCtx(ctx context.Context, key string) (int, error) {
	return rc.IncrByWithCtx(ctx, key, -1)
}

// IsExist Check if a cached value exists or not.
// add error as return value
func (rc *Cache) IsExist(key string) (bool, error) {
	return rc.IsExistWithCtx(context.Background(), key)
}

//IsExistWithCtx ..
func (rc *Cache) IsExistWithCtx(ctx context.Context, key string) (bool, error) {
	_, err := rc.client().Get(key)
	return err == nil, err
}

// ClearAll Clear all cache.
func (rc *Cache) ClearAll() error {
	return rc.ClearAllWithCtx(context.Background())
}

// ClearAllWithCtx  clears all cache in memcache.
func (rc *Cache) ClearAllWithCtx(ctx context.Context) error {
	return rc.conn.FlushAll()
}

// StartAndGC starts the memcache adapter.
// config: must be in the format {"conn":"connection info"}.
// If an error occurs during connecting, an error is returned
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	rc.client()
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache)
}
