package cache

import (
	"fmt"
)

type Cache interface {
	Get(key string) interface{}
	Put(key string, val interface{}, timeout int64) error
	Delete(key string) error
	Incr(key string) error
	Decr(key string) error
	IsExist(key string) bool
	ClearAll() error
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

// config need to be correct JSON as string: {"interval":360}
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
