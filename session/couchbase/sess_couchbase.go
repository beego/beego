// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package session

import (
	"net/http"
	"strings"
	"sync"

	"github.com/couchbaselabs/go-couchbase"

	"github.com/astaxie/beego/session"
)

var couchbpder = &CouchbaseProvider{}

type CouchbaseSessionStore struct {
	b           *couchbase.Bucket
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

type CouchbaseProvider struct {
	maxlifetime int64
	savePath    string
	pool        string
	bucket      string
	b           *couchbase.Bucket
}

func (cs *CouchbaseSessionStore) Set(key, value interface{}) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.values[key] = value
	return nil
}

func (cs *CouchbaseSessionStore) Get(key interface{}) interface{} {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	if v, ok := cs.values[key]; ok {
		return v
	} else {
		return nil
	}
}

func (cs *CouchbaseSessionStore) Delete(key interface{}) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	delete(cs.values, key)
	return nil
}

func (cs *CouchbaseSessionStore) Flush() error {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.values = make(map[interface{}]interface{})
	return nil
}

func (cs *CouchbaseSessionStore) SessionID() string {
	return cs.sid
}

func (cs *CouchbaseSessionStore) SessionRelease(w http.ResponseWriter) {
	defer cs.b.Close()

	// if rs.values is empty, return directly
	if len(cs.values) < 1 {
		cs.b.Delete(cs.sid)
		return
	}

	bo, err := session.EncodeGob(cs.values)
	if err != nil {
		return
	}

	cs.b.Set(cs.sid, int(cs.maxlifetime), bo)
}

func (cp *CouchbaseProvider) getBucket() *couchbase.Bucket {
	c, err := couchbase.Connect(cp.savePath)
	if err != nil {
		return nil
	}

	pool, err := c.GetPool(cp.pool)
	if err != nil {
		return nil
	}

	bucket, err := pool.GetBucket(cp.bucket)
	if err != nil {
		return nil
	}

	return bucket
}

// init couchbase session
// savepath like couchbase server REST/JSON URL
// e.g. http://host:port/, Pool, Bucket
func (cp *CouchbaseProvider) SessionInit(maxlifetime int64, savePath string) error {
	cp.maxlifetime = maxlifetime
	configs := strings.Split(savePath, ",")
	if len(configs) > 0 {
		cp.savePath = configs[0]
	}
	if len(configs) > 1 {
		cp.pool = configs[1]
	}
	if len(configs) > 2 {
		cp.bucket = configs[2]
	}

	return nil
}

// read couchbase session by sid
func (cp *CouchbaseProvider) SessionRead(sid string) (session.SessionStore, error) {
	cp.b = cp.getBucket()

	var doc []byte

	err := cp.b.Get(sid, &doc)
	var kv map[interface{}]interface{}
	if doc == nil {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(doc)
		if err != nil {
			return nil, err
		}
	}

	cs := &CouchbaseSessionStore{b: cp.b, sid: sid, values: kv, maxlifetime: cp.maxlifetime}
	return cs, nil
}

func (cp *CouchbaseProvider) SessionExist(sid string) bool {
	cp.b = cp.getBucket()
	defer cp.b.Close()

	var doc []byte

	if err := cp.b.Get(sid, &doc); err != nil || doc == nil {
		return false
	} else {
		return true
	}
}

func (cp *CouchbaseProvider) SessionRegenerate(oldsid, sid string) (session.SessionStore, error) {
	cp.b = cp.getBucket()

	var doc []byte
	if err := cp.b.Get(oldsid, &doc); err != nil || doc == nil {
		cp.b.Set(sid, int(cp.maxlifetime), "")
	} else {
		err := cp.b.Delete(oldsid)
		if err != nil {
			return nil, err
		}
		_, _ = cp.b.Add(sid, int(cp.maxlifetime), doc)
	}

	err := cp.b.Get(sid, &doc)
	if err != nil {
		return nil, err
	}
	var kv map[interface{}]interface{}
	if doc == nil {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob(doc)
		if err != nil {
			return nil, err
		}
	}

	cs := &CouchbaseSessionStore{b: cp.b, sid: sid, values: kv, maxlifetime: cp.maxlifetime}
	return cs, nil
}

func (cp *CouchbaseProvider) SessionDestroy(sid string) error {
	cp.b = cp.getBucket()
	defer cp.b.Close()

	cp.b.Delete(sid)
	return nil
}

func (cp *CouchbaseProvider) SessionGC() {
	return
}

func (cp *CouchbaseProvider) SessionAll() int {
	return 0
}

func init() {
	session.Register("couchbase", couchbpder)
}
