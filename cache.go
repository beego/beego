package beego

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	Kb = 1024
	Mb = 1024 * 1024
	Gb = 1024 * 1024 * 1024
)

var (
	DefaultEvery int = 60 // 1 minute
)

var (
	InvalidCacheItem = errors.New("invalid cache item")
	ItemIsDirectory  = errors.New("can't cache a directory")
	ItemNotInCache   = errors.New("item not in cache")
	ItemTooLarge     = errors.New("item too large for cache")
	WriteIncomplete  = errors.New("incomplete write of cache item")
)

type BeeItem struct {
	val        interface{}
	Lastaccess time.Time
	expired    int
}

func (itm *BeeItem) Access() interface{} {
	itm.Lastaccess = time.Now()
	return itm.val
}

type BeeCache struct {
	dur   time.Duration
	items map[string]*BeeItem
	Every int // Run an expiration check Every seconds
}

// NewDefaultCache returns a new FileCache with sane defaults.
func NewBeeCache() *BeeCache {
	cache := BeeCache{time.Since(time.Now()),
		nil,
		DefaultEvery}
	return &cache
}

func (bc *BeeCache) Get(name string) interface{} {
	itm, ok := bc.items[name]
	if !ok {
		return nil
	}
	return itm.Access()
}

func (bc *BeeCache) Put(name string, value interface{}, expired int) error {
	t := BeeItem{val: value, Lastaccess: time.Now(), expired: expired}
	if bc.IsExist(name) {
		return errors.New("the key is exist")
	} else {
		bc.items[name] = &t
	}
	return nil
}

func (bc *BeeCache) Delete(name string) (ok bool, err error) {
	_, ok = bc.items[name]
	if !ok {
		return
	}
	delete(bc.items, name)
	_, valid := bc.items[name]
	if valid {
		ok = false
	}
	return
}

func (bc *BeeCache) IsExist(name string) bool {
	_, ok := bc.items[name]
	return ok
}

// Start activates the file cache; it will 
func (bc *BeeCache) Start() error {
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", bc.Every))
	if err != nil {
		return err
	}
	bc.dur = dur
	bc.items = make(map[string]*BeeItem, 0)
	go bc.vaccuum()
	return nil
}

func (bc *BeeCache) vaccuum() {
	if bc.Every < 1 {
		return
	}
	for {
		<-time.After(time.Duration(bc.dur))
		if bc.items == nil {
			return
		}
		for name, _ := range bc.items {
			if bc.item_expired(name) {
				delete(bc.items, name)
			}
		}
	}
}

// item_expired returns true if an item is expired.
func (bc *BeeCache) item_expired(name string) bool {
	itm, ok := bc.items[name]
	if !ok {
		return true
	}
	dur := time.Now().Sub(itm.Lastaccess)
	sec, err := strconv.Atoi(fmt.Sprintf("%0.0f", dur.Seconds()))
	if err != nil {
		return true
	} else if sec >= itm.expired {
		return true
	}
	return false
}
