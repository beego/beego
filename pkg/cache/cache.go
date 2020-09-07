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

// Package cache provide a Cache interface and some implement engine
// Usage:
//
// import(
//   "github.com/astaxie/beego/cache"
// )
//
// bm, err := cache.NewCache("memory", `{"interval":60}`)
//
// Use it like this:
//
//	bm.Put("astaxie", 1, 10 * time.Second)
//	bm.Get("astaxie")
//	bm.IsExist("astaxie")
//	bm.Delete("astaxie")
//
//  more docs http://beego.me/docs/module/cache.md
package cache

import (
	"context"
	"fmt"
	"time"
)

// Cache interface contains all behaviors for cache adapter.
// usage:
//	cache.Register("file",cache.NewFileCache) // this operation is run in init method of file.go.
//	c,err := cache.NewCache("file","{....}")
//	c.Put("key",value, 3600 * time.Second)
//	v := c.Get("key")
//
//	c.Incr("counter")  // now is 1
//	c.Incr("counter")  // now is 2
//	count := c.Get("counter").(int)
type Cache interface {
	// Get a cached value by key.
	// add error as return value
	Get(key string) (interface{}, error)
	GetWithCtx(ctx context.Context, key string) (interface{}, error)

	// those two methods won't return error.
	// default value will be return if error or not found
	// if key not found, the defValue won't be set into cache
	// GetOr(key string, defValue interface{}) interface{}
	// GetOrWithCtx(ctx context.Context, key string, defValue interface{}) interface{}

	// those two methods' behavior look like:
	// if IsExit(key)
	//    return false, Get(key)
	// else
	//    Put(key, candidate)
	//    return ...
	// bool -> whether use candidate
	// interface{} -> final result related to the key. If first return value is true, this is candidate
	// error -> any error
	// all implementations must be thread-safe
	// GetOrSet(key string, candidate interface{}) (bool, interface{}, error)
	// GetOrSetWithCtx(ctx context.Context, key string, candidate interface{}) (bool, interface{}, error)

	// GetMulti is a batch version of Get.
	// add error as return value
	GetMulti(keys []string) ([]interface{}, error)
	GetMultiWithCtx(ctx context.Context, keys []string) ([]interface{}, error)
	// Set a cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error
	PutWithCtx(ctx context.Context, key string, val interface{}, timeout time.Duration) error

	// Delete cached value by key.
	Delete(key string) error
	DeleteWithCtx(ctx context.Context, key string) error

	// DeleteAndRet try to delete this key, and return old value
	// if interface{} is nil(or zero value), means key is absent
	// DeleteAndRet(key string) (interface{}, error)
	// DeleteAndRetWithCtx(key string) (interface{}, error)

	// // int indicates that how many keys are deleted
	// // error == nil && int == 0 means that all keys are absent
	// DeleteMulti(keys []string) (int, error)
	// DeleteMultiWithCtx(keys []string) (int, error)

	// Increment a cached int value by key, as a counter.
	// int indicates current value after increasing
	IncrBy(key string, n int) (int, error)
	IncrByWithCtx(ctx context.Context, key string, n int) (int, error)

	// Increment a cached int value by key, as a counter.
	// int indicates current value after increasing
	Incr(key string) (int, error)
	IncrWithCtx(ctx context.Context, key string) (int, error)

	// Decrement a cached int value by key, as a counter.
	// int indicates current value after decreasing
	Decr(key string) (int, error)
	DecrWithCtx(ctx context.Context, key string) (int, error)

	// Check if a cached value exists or not.
	// add error as return value
	IsExist(key string) (bool, error)
	IsExistWithCtx(ctx context.Context, key string) (bool, error)

	// Clear all cache.
	ClearAll() error
	ClearAllWithCtx(ctx context.Context) error
	// Start gc routine based on config string settings.
	StartAndGC(config string) error
	// StartAndGCWithCtx(ctx context.Context, config string) error
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache creates a new cache driver by adapter name and config string.
// config: must be in JSON format such as {"interval":360}.
// Starts gc automatically.
func NewCache(adapterName, config string) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
