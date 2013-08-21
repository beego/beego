package config

import (
	"fmt"
)

type ConfigContainer interface {
	Set(key, val string) error
	String(key string) string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DIY(key string) (interface{}, error)
}

type Config interface {
	Parse(key string) (ConfigContainer, error)
}

var adapters = make(map[string]Config)

// Register makes a config adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("config: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// config need to be correct JSON as string: {"interval":360}
func NewConfig(adapterName, fileaname string) (ConfigContainer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	return adapter.Parse(fileaname)
}
