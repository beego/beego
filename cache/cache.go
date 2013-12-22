package cache

import (
	"fmt"
)

// Cache interface contains all behaviors for cache adapter.
// usage:
//	cache.Register("file",cache.NewFileCache()) // this operation is run in init method of file.go.
//	c := cache.NewCache("file","{....}")
//	c.Put("key",value,3600)
//	v := c.Get("key")
//
//	c.Incr("counter")  // now is 1
//	c.Incr("counter")  // now is 2
//	count := c.Get("counter").(int)
type Cache interface {
	// get cached value by key.
	Get(key string) interface{}
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout int64) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check cached value is existed or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
	// start gc routine via config string setting.
	StartAndGC(config string) error
}

var adapters = make(map[string]Cache)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Cache) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// Create a new cache driver by adapter and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(adapterName, config string) (Cache, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adaptername %q (forgotten import?)", adapterName)
	}
	err := adapter.StartAndGC(config)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}
