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
//   _ "github.com/beego/beego/v2/client/cache/memcache"
//   "github.com/beego/beego/v2/client/cache"
// )
//
//  bm, err := cache.NewCache("memcache", `{"conn":"127.0.0.1:11211"}`)
//
package memcache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/berror"
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

// Get get value from memcache.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	if item, err := rc.conn.Get(key); err == nil {
		return item.Value, nil
	} else {
		return nil, berror.Wrapf(err, cache.MemCacheCurdFailed,
			"could not read data from memcache, please check your key, network and connection. Root cause: %s",
			err.Error())
	}
}

// GetMulti gets a value from a key in memcache.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	rv := make([]interface{}, len(keys))

	mv, err := rc.conn.GetMulti(keys)
	if err != nil {
		return rv, berror.Wrapf(err, cache.MemCacheCurdFailed,
			"could not read multiple key-values from memcache, "+
				"please check your keys, network and connection. Root cause: %s",
			err.Error())
	}

	keysErr := make([]string, 0)
	for i, ki := range keys {
		if _, ok := mv[ki]; !ok {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, "key not exist"))
			continue
		}
		rv[i] = mv[ki].Value
	}

	if len(keysErr) == 0 {
		return rv, nil
	}
	return rv, berror.Error(cache.MultiGetFailed, strings.Join(keysErr, "; "))
}

// Put puts a value into memcache.
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	item := memcache.Item{Key: key, Expiration: int32(timeout / time.Second)}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return berror.Errorf(cache.InvalidMemCacheValue,
			"the value must be string or byte[]. key: %s, value:%v", key, val)
	}
	return berror.Wrapf(rc.conn.Set(&item), cache.MemCacheCurdFailed,
		"could not put key-value to memcache, key: %s", key)
}

// Delete deletes a value in memcache.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	return berror.Wrapf(rc.conn.Delete(key), cache.MemCacheCurdFailed,
		"could not delete key-value from memcache, key: %s", key)
}

// Incr increases counter.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	_, err := rc.conn.Increment(key, 1)
	return berror.Wrapf(err, cache.MemCacheCurdFailed,
		"could not increase value for key: %s", key)
}

// Decr decreases counter.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	_, err := rc.conn.Decrement(key, 1)
	return berror.Wrapf(err, cache.MemCacheCurdFailed,
		"could not decrease value for key: %s", key)
}

// IsExist checks if a value exists in memcache.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	_, err := rc.Get(ctx, key)
	return err == nil, err
}

// ClearAll clears all cache in memcache.
func (rc *Cache) ClearAll(context.Context) error {
	return berror.Wrap(rc.conn.FlushAll(), cache.MemCacheCurdFailed,
		"try to clear all key-value pairs failed")
}

// StartAndGC starts the memcache adapter.
// config: must be in the format {"conn":"connection info"}.
// If an error occurs during connecting, an error is returned
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	if err := json.Unmarshal([]byte(config), &cf); err != nil {
		return berror.Wrapf(err, cache.InvalidMemCacheCfg,
			"could not unmarshal this config, it must be valid json stringP: %s", config)
	}

	if _, ok := cf["conn"]; !ok {
		return berror.Errorf(cache.InvalidMemCacheCfg, `config must contains "conn" field: %s`, config)
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	rc.conn = memcache.New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache)
}
