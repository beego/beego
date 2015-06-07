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

// package memcahe for cache provider
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
	"encoding/json"
	"errors"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/astaxie/beego/cache"
)

// Memcache adapter.
type MemcacheCache struct {
	conn     *memcache.Client
	conninfo []string
}

// create new memcache adapter.
func NewMemCache() *MemcacheCache {
	return &MemcacheCache{}
}

// get value from memcache.
func (rc *MemcacheCache) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return string(item.Value)
	}
	return nil
}

// get value from memcache.
func (rc *MemcacheCache) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var rv []interface{}
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			for i := 0; i < size; i++ {
				rv = append(rv, err)
			}
			return rv
		}
	}
	mv, err := rc.conn.GetMulti(keys)
	if err == nil {
		for _, v := range mv {
			rv = append(rv, string(v.Value))
		}
		return rv
	} else {
		for i := 0; i < size; i++ {
			rv = append(rv, err)
		}
		return rv
	}
}

// put value to memcache. only support string.
func (rc *MemcacheCache) Put(key string, val interface{}, timeout int64) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	v, ok := val.(string)
	if !ok {
		return errors.New("val must string")
	}
	item := memcache.Item{Key: key, Value: []byte(v), Expiration: int32(timeout)}
	return rc.conn.Set(&item)
}

// delete value in memcache.
func (rc *MemcacheCache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// increase counter.
func (rc *MemcacheCache) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// decrease counter.
func (rc *MemcacheCache) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// check value exists in memcache.
func (rc *MemcacheCache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	_, err := rc.conn.Get(key)
	if err != nil {
		return false
	}
	return true
}

// clear all cached in memcache.
func (rc *MemcacheCache) ClearAll() error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.FlushAll()
}

// start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *MemcacheCache) StartAndGC(config string) error {
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
func (rc *MemcacheCache) connectInit() error {
	rc.conn = memcache.New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache())
}
