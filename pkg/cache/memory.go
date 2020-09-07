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

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"sync"
	"time"
)

var (
	//DefaultEvery Timer for how often to recycle the expired cache items in memory (in seconds)
	DefaultEvery = 60 // 1 minute

	//ErrKeyExpire error
	ErrKeyExpire = errors.New("memory: with an expire key")
	//ErrKeyNotFound error
	ErrKeyNotFound = errors.New("memory: key not exist")
	//ErrValNotLogicType (u)int (u)int32 (u)int64
	ErrValNotLogicType = errors.New("memory: val not Logic operations")
)

// MemoryItem stores memory cache item.
type MemoryItem struct {
	val         interface{}
	createdTime time.Time
	lifespan    time.Duration
}

func (mi *MemoryItem) isExpire() bool {
	// 0 means forever
	if mi.lifespan == 0 {
		return false
	}
	return time.Now().Sub(mi.createdTime) > mi.lifespan
}

// MemoryCache is a memory cache adapter.
// Contains a RW locker for safe map storage.
type MemoryCache struct {
	sync.RWMutex
	dur   time.Duration
	items map[string]*MemoryItem
	Every int // run an expiration check Every clock time
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() Cache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

// Get If non-existent or expired, return nil.
func (bc *MemoryCache) Get(key string) (interface{}, error) {
	return bc.GetWithCtx(context.Background(), key)
}

//GetWithCtx If non-existent or expired, return nil.
func (bc *MemoryCache) GetWithCtx(ctx context.Context, key string) (interface{}, error) {
	bc.RLock()
	defer bc.RUnlock()
	if itm, ok := bc.items[key]; ok {
		if itm.isExpire() {
			return itm.val, ErrKeyExpire
		}
		return itm.val, nil
	}
	return nil, ErrKeyNotFound
}

// GetMulti is a batch version of Get.
// add error as return value
func (bc *MemoryCache) GetMulti(keys []string) ([]interface{}, error) {
	return bc.GetMultiWithCtx(context.Background(), keys)
}

//GetMultiWithCtx gets caches from memory.
// If non-existent or expired, return nil.
func (bc *MemoryCache) GetMultiWithCtx(ctx context.Context, keys []string) ([]interface{}, error) {
	var rc []interface{}
	for _, k := range keys {
		v, err := bc.GetWithCtx(ctx, k)
		if err != nil {
			return rc, err
		}
		rc = append(rc, v)
	}
	return rc, nil
}

// Put puts cache into memory.
// If lifespan is 0, it will never overwrite this value unless restarted
func (bc *MemoryCache) Put(key string, val interface{}, timeout time.Duration) error {
	return bc.PutWithCtx(context.Background(), key, val, timeout)
}

// PutWithCtx puts cache into memory.
// If lifespan is 0, it will never overwrite this value unless restarted
func (bc *MemoryCache) PutWithCtx(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items[key] = &MemoryItem{
		val:         val,
		createdTime: time.Now(),
		lifespan:    timeout,
	}
	return nil
}

// Delete cached value by key.
func (bc *MemoryCache) Delete(key string) error {
	return bc.DeleteWithCtx(context.Background(), key)
}

//DeleteWithCtx ..
func (bc *MemoryCache) DeleteWithCtx(ctx context.Context, key string) error {
	bc.Lock()
	defer bc.Unlock()
	delete(bc.items, key)
	return nil
}

// IncrBy a cached int value by key, as a counter.
// int indicates current value after increasing
func (bc *MemoryCache) IncrBy(key string, n int) (int, error) {
	return bc.IncrByWithCtx(context.Background(), key, n)
}

//IncrByWithCtx ..
func (bc *MemoryCache) IncrByWithCtx(ctx context.Context, key string, n int) (int, error) {
	bc.Lock()
	defer bc.Unlock()
	itm, ok := bc.items[key]
	if !ok {
		return n, bc.Put(key, n, 0)
	}
	switch val := itm.val.(type) {
	case int:
		itm.val = val + int(n)
		return itm.val.(int), nil
	case int32:
		itm.val = val + int32(n)
		return int(itm.val.(int32)), nil
	case int64:
		itm.val = val + int64(n)
		return int(itm.val.(int64)), nil
	case uint:
		if n > 0 {
			itm.val = val + uint(n)
		} else {
			itm.val = val - uint(n)
		}
		return int(itm.val.(uint)), nil
	case uint32:
		if n > 0 {
			itm.val = val + uint32(n)
		} else {
			itm.val = val - uint32(n)
		}
		return int(itm.val.(uint32)), nil
	case uint64:
		if n > 0 {
			itm.val = val + uint64(n)
		} else {
			itm.val = val - uint64(n)
		}
		return int(itm.val.(uint64)), nil
	default:
		return 0, ErrValNotLogicType
	}
}

// Incr a cached int value by key, as a counter.
// int indicates current value after increasing
// Supports int,int32,int64,uint,uint32,uint64.
func (bc *MemoryCache) Incr(key string) (int, error) {
	return bc.IncrBy(key, 1)
}

//IncrWithCtx ..
// Supports int,int32,int64,uint,uint32,uint64.
func (bc *MemoryCache) IncrWithCtx(ctx context.Context, key string) (int, error) {
	return bc.IncrByWithCtx(context.Background(), key, 1)
}

// Decr a cached int value by key, as a counter.
// int indicates current value after decreasing
func (bc *MemoryCache) Decr(key string) (int, error) {
	return bc.IncrBy(key, -1)
}

// DecrWithCtx ..
// Supports int,int32,int64,uint,uint32,uint64.
func (bc *MemoryCache) DecrWithCtx(ctx context.Context, key string) (int, error) {
	return bc.IncrByWithCtx(context.Background(), key, -1)
}

// IsExist Check if a cached value exists or not.
// add error as return value
func (bc *MemoryCache) IsExist(key string) (bool, error) {
	return bc.IsExistWithCtx(context.Background(), key)
}

//IsExistWithCtx ..
func (bc *MemoryCache) IsExistWithCtx(ctx context.Context, key string) (bool, error) {
	bc.RLock()
	defer bc.RUnlock()
	if v, ok := bc.items[key]; ok {
		return !v.isExpire(), ErrKeyExpire
	}
	return false, ErrKeyNotFound
}

// ClearAll Clear all cache.
func (bc *MemoryCache) ClearAll() error {
	return bc.ClearAllWithCtx(context.Background())
}

//ClearAllWithCtx deletes all cache in memory.
func (bc *MemoryCache) ClearAllWithCtx(ctx context.Context) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items = make(map[string]*MemoryItem)
	runtime.GC() //todo ? need gc?
	return nil
}

// StartAndGC starts memory cache. Checks expiration in every clock time.
func (bc *MemoryCache) StartAndGC(config string) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultEvery
	}
	dur := time.Duration(cf["interval"]) * time.Second
	bc.Every = cf["interval"]
	bc.dur = dur
	go bc.vacuum()
	return nil
}

// check expiration.
func (bc *MemoryCache) vacuum() {
	bc.RLock()
	every := bc.Every
	bc.RUnlock()

	if every < 1 {
		return
	}
	for {
		<-time.After(bc.dur)
		bc.RLock()
		if bc.items == nil {
			bc.RUnlock()
			return
		}
		bc.RUnlock()
		if keys := bc.expiredKeys(); len(keys) != 0 {
			bc.clearItems(keys)
		}
	}
}

// expiredKeys returns keys list which are expired.
func (bc *MemoryCache) expiredKeys() (keys []string) {
	bc.RLock()
	defer bc.RUnlock()
	for key, itm := range bc.items {
		if itm.isExpire() {
			keys = append(keys, key)
		}
	}
	return
}

// ClearItems removes all items who's key is in keys
func (bc *MemoryCache) clearItems(keys []string) {
	bc.Lock()
	defer bc.Unlock()
	for _, key := range keys {
		delete(bc.items, key)
	}
}

func init() {
	Register("memory", NewMemoryCache)
}
