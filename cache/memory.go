package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	DefaultEvery int = 60 // 1 minute
)

type MemoryItem struct {
	val        interface{}
	Lastaccess time.Time
	expired    int64
}

type MemoryCache struct {
	lock  sync.RWMutex
	dur   time.Duration
	items map[string]*MemoryItem
	Every int // Run an expiration check Every seconds
}

// NewDefaultCache returns a new FileCache with sane defaults.
func NewMemoryCache() *MemoryCache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

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

func (bc *MemoryCache) Put(name string, value interface{}, expired int64) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	t := MemoryItem{
		val:        value,
		Lastaccess: time.Now(),
		expired:    expired,
	}
	if _, ok := bc.items[name]; ok {
		return errors.New("the key is exist")
	} else {
		bc.items[name] = &t
	}
	return nil
}

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

func (bc *MemoryCache) IsExist(name string) bool {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	_, ok := bc.items[name]
	return ok
}

func (bc *MemoryCache) ClearAll() error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	bc.items = make(map[string]*MemoryItem)
	return nil
}

// Start activates the file cache; it will
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
