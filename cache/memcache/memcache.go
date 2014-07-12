// Beego (http://beego.me/)
//
// @description beego is an open-source, high-performance web framework for the Go programming language.
//
// @link        http://github.com/astaxie/beego for the canonical source repository
//
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
//
// @authors     astaxie
package memcache

import (
	"encoding/json"
	"errors"

	"github.com/beego/memcache"

	"github.com/astaxie/beego/cache"
)

// Memcache adapter.
type MemcacheCache struct {
	conn     *memcache.Connection
	conninfo string
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
	if v, err := rc.conn.Get(key); err == nil {
		if len(v) > 0 {
			return string(v[0].Value)
		}
	}
	return nil
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
	stored, err := rc.conn.Set(key, 0, uint64(timeout), []byte(v))
	if err == nil && stored == false {
		return errors.New("stored fail")
	}
	return err
}

// delete value in memcache.
func (rc *MemcacheCache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Delete(key)
	return err
}

// [Not Support]
// increase counter.
func (rc *MemcacheCache) Incr(_ string) error {
	return errors.New("not support in memcache")
}

// [Not Support]
// decrease counter.
func (rc *MemcacheCache) Decr(_ string) error {
	return errors.New("not support in memcache")
}

// check value exists in memcache.
func (rc *MemcacheCache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	v, err := rc.conn.Get(key)
	if err != nil {
		return false
	}
	if len(v) == 0 {
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
	rc.conninfo = cf["conn"]
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *MemcacheCache) connectInit() error {
	c, err := memcache.Connect(rc.conninfo)
	if err != nil {
		return err
	}
	rc.conn = c
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache())
}
