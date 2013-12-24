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
	itm, ok := bc.items[name]
	if !ok {
		return nil
	}
	if (time.Now().Unix() - itm.Lastaccess.Unix()) > itm.expired {
		go bc.Delete(name)
		return nil
	}
	return itm.val
}

// Put cache to memory.
// if expired is 0, it will be cleaned by next gc operation ( default gc clock is 1 minute).
func (bc *MemoryCache) Put(name string, value interface{}, expired int64) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	t := MemoryItem{
		val:        value,
		Lastaccess: time.Now(),
		expired:    expired,
	}
	bc.items[name] = &t
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
	_, valid := bc.items[name]
	if valid {
		return errors.New("delete key error")
	}
	return nil
}

// Increase cache counter in memory.
// it supports int,int64,int32,uint,uint64,uint32.
func (bc *MemoryCache) Incr(key string) error {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) + 1
	case int64:
		itm.val = itm.val.(int64) + 1
	case int32:
		itm.val = itm.val.(int32) + 1
	case uint:
		itm.val = itm.val.(uint) + 1
	case uint32:
		itm.val = itm.val.(uint32) + 1
	case uint64:
		itm.val = itm.val.(uint64) + 1
	default:
		return errors.New("item val is not int int64 int32")
	}
	return nil
}

// Decrease counter in memory.
func (bc *MemoryCache) Decr(key string) error {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) - 1
	case int64:
		itm.val = itm.val.(int64) - 1
	case int32:
		itm.val = itm.val.(int32) - 1
	case uint:
		if itm.val.(uint) > 0 {
			itm.val = itm.val.(uint) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint32:
		if itm.val.(uint32) > 0 {
			itm.val = itm.val.(uint32) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint64:
		if itm.val.(uint64) > 0 {
			itm.val = itm.val.(uint64) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	default:
		return errors.New("item val is not int int64 int32")
	}
	return nil
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
	sec := time.Now().Unix() - itm.Lastaccess.Unix()
	if sec >= itm.expired {
		delete(bc.items, name)
		return true
	}
	return false
}

func init() {
	Register("memory", NewMemoryCache())
}
