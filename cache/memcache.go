package cache

import (
	"encoding/json"
	"errors"

	"github.com/beego/memcache"
)

// Memcache adapter.
type MemcacheCache struct {
	c        *memcache.Connection
	conninfo string
}

// create new memcache adapter.
func NewMemCache() *MemcacheCache {
	return &MemcacheCache{}
}

// get value from memcache.
func (rc *MemcacheCache) Get(key string) interface{} {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := rc.c.Get(key)
	if err != nil {
		return nil
	}
	var contain interface{}
	if len(v) > 0 {
		contain = string(v[0].Value)
	} else {
		contain = nil
	}
	return contain
}

// put value to memcache. only support string.
func (rc *MemcacheCache) Put(key string, val interface{}, timeout int64) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, ok := val.(string)
	if !ok {
		return errors.New("val must string")
	}
	stored, err := rc.c.Set(key, 0, uint64(timeout), []byte(v))
	if err == nil && stored == false {
		return errors.New("stored fail")
	}
	return err
}

// delete value in memcache.
func (rc *MemcacheCache) Delete(key string) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Delete(key)
	return err
}

// [Not Support]
// increase counter.
func (rc *MemcacheCache) Incr(key string) error {
	return errors.New("not support in memcache")
}

// [Not Support]
// decrease counter.
func (rc *MemcacheCache) Decr(key string) error {
	return errors.New("not support in memcache")
}

// check value exists in memcache.
func (rc *MemcacheCache) IsExist(key string) bool {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := rc.c.Get(key)
	if err != nil {
		return false
	}
	if len(v) == 0 {
		return false
	} else {
		return true
	}
	return true
}

// clear all cached in memcache.
func (rc *MemcacheCache) ClearAll() error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	err := rc.c.FlushAll()
	return err
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
	rc.c = rc.connectInit()
	if rc.c == nil {
		return errors.New("dial tcp conn error")
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *MemcacheCache) connectInit() *memcache.Connection {
	c, err := memcache.Connect(rc.conninfo)
	if err != nil {
		return nil
	}
	return c
}

func init() {
	Register("memcache", NewMemCache())
}
