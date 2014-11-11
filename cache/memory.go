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
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// clock time of recycling the expired cache items in memory.
	DefaultEvery int = 60 // 1 minute
)

// Memory cache item.
type MemoryItem struct {
	val        interface{}
	Lastaccess time.Time
	expired    int64
}

// Memory cache adapter.
// it contains a RW locker for safe map storage.
type MemoryCache struct {
	lock  sync.RWMutex
	dur   time.Duration
	items map[string]*MemoryItem
	Every int // run an expiration check Every clock time
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() *MemoryCache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

// Get cache from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) Get(name string) interface{} {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	if itm, ok := bc.items[name]; ok {
		if (time.Now().Unix() - itm.Lastaccess.Unix()) > itm.expired {
			go bc.Delete(name)
			return nil
		}
		return itm.val
	}
	return nil
}

// Put cache to memory.
// if expired is 0, it will be cleaned by next gc operation ( default gc clock is 1 minute).
func (bc *MemoryCache) Put(name string, value interface{}, expired int64) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	bc.items[name] = &MemoryItem{
		val:        value,
		Lastaccess: time.Now(),
		expired:    expired,
	}
	return nil
}

/// Delete cache in memory.
func (bc *MemoryCache) Delete(name string) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if _, ok := bc.items[name]; !ok {
		return errors.New("key not exist")
	}
	delete(bc.items, name)
	if _, ok := bc.items[name]; ok {
		return errors.New("delete key error")
	}
	return nil
}

func convertToCounter(value interface{}) (converted int64, err error) {
	switch value.(type) {
	case uint64:
		converted = int64(value.(uint64))
	case int:
		converted = int64(value.(int))
	case int32:
		converted = int64(value.(int32))
	case uint:
		converted = int64(value.(uint))
	case uint32:
		converted = int64(value.(uint32))
	case int64:
		converted = value.(int64)
	default:
		err = errors.New("value cannot be converted to int64")
	}
	return
}

// Increase cache counter in memory.
// it supports int,int64,int32,uint,uint64,uint32.
func (bc *MemoryCache) Incr(key string) (int64, error) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return 0, errors.New("key not exist")
	}

	counter, err := convertToCounter(itm.val)
	if err != nil {
		return 0, err
	} else {
		counter += 1
		itm.val = counter
		return counter, nil
	}
}

// Decrease counter in memory.
func (bc *MemoryCache) Decr(key string) (int64, error) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return 0, errors.New("key not exist")
	}
	counter, err := convertToCounter(itm.val)
	if err != nil {
		return 0, err
	} else {
		counter -= 1
		itm.val = counter
		return counter, nil
	}
}

// check cache exist in memory.
func (bc *MemoryCache) IsExist(name string) bool {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	_, ok := bc.items[name]
	return ok
}

// delete all cache in memory.
func (bc *MemoryCache) ClearAll() error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	bc.items = make(map[string]*MemoryItem)
	return nil
}

// start memory cache. it will check expiration in every clock time.
func (bc *MemoryCache) StartAndGC(config string) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultEvery
	}
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", cf["interval"]))
	if err != nil {
		return err
	}
	bc.Every = cf["interval"]
	bc.dur = dur
	go bc.vaccuum()
	return nil
}

// check expiration.
func (bc *MemoryCache) vaccuum() {
	if bc.Every < 1 {
		return
	}
	for {
		<-time.After(bc.dur)
		if bc.items == nil {
			return
		}
		for name, _ := range bc.items {
			bc.item_expired(name)
		}
	}
}

// item_expired returns true if an item is expired.
func (bc *MemoryCache) item_expired(name string) bool {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	itm, ok := bc.items[name]
	if !ok {
		return true
	}
	if time.Now().Unix()-itm.Lastaccess.Unix() >= itm.expired {
		delete(bc.items, name)
		return true
	}
	return false
}

func init() {
	Register("memory", NewMemoryCache())
}
