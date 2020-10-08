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

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/astaxie/beego/client/cache"
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
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, err
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return item.Value, nil
	} else {
		return nil, err
	}
}

// GetMulti gets a value from a key in memcache.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	var rv []interface{}
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return rv, err
		}
	}
	mv, err := rc.conn.GetMulti(keys)
	if err == nil {
		for _, v := range mv {
			rv = append(rv, v.Value)
		}
	}
	return rv, err
}

// Put puts a value into memcache.
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	item := memcache.Item{Key: key, Expiration: int32(timeout / time.Second)}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return errors.New("val only support string and []byte")
	}
	return rc.conn.Set(&item)
}

// Delete deletes a value in memcache.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// Incr increases counter.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// Decr decreases counter.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// IsExist checks if a value exists in memcache.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false, err
		}
	}
	_, err := rc.conn.Get(key)
	return err == nil, err
}

// ClearAll clears all cache in memcache.
func (rc *Cache) ClearAll(context.Context) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
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
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *Cache) connectInit() error {
	rc.conn = memcache.New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache)
}
