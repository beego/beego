// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package session

import (
	"net/http"
	"sync"

	"github.com/astaxie/beego/session"

	"github.com/beego/memcache"
)

var mempder = &MemProvider{}

// memcache session store
type MemcacheSessionStore struct {
	c           *memcache.Connection
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// set value in memcache session
func (rs *MemcacheSessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// get value in memcache session
func (rs *MemcacheSessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	} else {
		return nil
	}
}

// delete value in memcache session
func (rs *MemcacheSessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// clear all values in memcache session
func (rs *MemcacheSessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// get redis session id
func (rs *MemcacheSessionStore) SessionID() string {
	return rs.sid
}

// save session values to redis
func (rs *MemcacheSessionStore) SessionRelease(w http.ResponseWriter) {
	defer rs.c.Close()
	// if rs.values is empty, return directly
	if len(rs.values) < 1 {
		rs.c.Delete(rs.sid)
		return
	}

	b, err := session.EncodeGob(rs.values)
	if err != nil {
		return
	}
	rs.c.Set(rs.sid, 0, uint64(rs.maxlifetime), b)
}

// redis session provider
type MemProvider struct {
	maxlifetime int64
	savePath    string
	poolsize    int
	password    string
}

// init redis session
// savepath like
// e.g. 127.0.0.1:9090
func (rp *MemProvider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	rp.savePath = savePath
	return nil
}

// read redis session by sid
func (rp *MemProvider) SessionRead(sid string) (session.SessionStore, error) {
	conn, err := rp.connectInit()
	if err != nil {
		return nil, err
	}
	kvs, err := conn.Get(sid)
	if err != nil {
		return nil, err
	}
	var contain []byte
	if len(kvs) > 0 {
		contain = kvs[0].Value
	}
	var kv map[interface{}]interface{}
	if len(contain) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(contain)
		if err != nil {
			return nil, err
		}
	}

	rs := &MemcacheSessionStore{c: conn, sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// check redis session exist by sid
func (rp *MemProvider) SessionExist(sid string) bool {
	conn, err := rp.connectInit()
	if err != nil {
		return false
	}
	defer conn.Close()
	if kvs, err := conn.Get(sid); err != nil || len(kvs) == 0 {
		return false
	} else {
		return true
	}
}

// generate new sid for redis session
func (rp *MemProvider) SessionRegenerate(oldsid, sid string) (session.SessionStore, error) {
	conn, err := rp.connectInit()
	if err != nil {
		return nil, err
	}
	var contain []byte
	if kvs, err := conn.Get(sid); err != nil || len(kvs) == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
		conn.Set(sid, 0, uint64(rp.maxlifetime), []byte(""))
	} else {
		conn.Delete(oldsid)
		conn.Set(sid, 0, uint64(rp.maxlifetime), kvs[0].Value)
		contain = kvs[0].Value
	}

	var kv map[interface{}]interface{}
	if len(contain) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(contain)
		if err != nil {
			return nil, err
		}
	}

	rs := &MemcacheSessionStore{c: conn, sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// delete redis session by id
func (rp *MemProvider) SessionDestroy(sid string) error {
	conn, err := rp.connectInit()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Delete(sid)
	if err != nil {
		return err
	}
	return nil
}

// Impelment method, no used.
func (rp *MemProvider) SessionGC() {
	return
}

// @todo
func (rp *MemProvider) SessionAll() int {
	return 0
}

// connect to memcache and keep the connection.
func (rp *MemProvider) connectInit() (*memcache.Connection, error) {
	c, err := memcache.Connect(rp.savePath)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func init() {
	session.Register("memcache", mempder)
}
