// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package utils

import (
	"sync"
)

type BeeMap struct {
	lock *sync.RWMutex
	bm   map[interface{}]interface{}
}

// NewBeeMap return new safemap
func NewBeeMap() *BeeMap {
	return &BeeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

// Get from maps return the k's value
func (m *BeeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BeeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Returns true if k is exist in the map.
func (m *BeeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

// Delete the given key and value.
func (m *BeeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

// Items returns all items in safemap.
func (m *BeeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.bm
}
