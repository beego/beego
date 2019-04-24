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

// Package redis for session provider
//
// depend on github.com/go-redis/redis
//
// go install github.com/go-redis/redis
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/session/redis_sentinel"
//   "github.com/astaxie/beego/session"
// )
//
//	func init() {
//		globalSessions, _ = session.NewManager("redis_sentinel", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:26379;127.0.0.2:26379"}``)
//		go globalSessions.GC()
//	}
//
// more detail about params: please check the notes on the function SessionInit in this package
package redis_sentinel

import (
	"github.com/astaxie/beego/session"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var redispder = &Provider{}

// DefaultPoolSize redis_sentinel default pool size
var DefaultPoolSize = 100

// SessionStore redis_sentinel session store
type SessionStore struct {
	p           *redis.Client
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Set value in redis_sentinel session
func (rs *SessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in redis_sentinel session
func (rs *SessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in redis_sentinel session
func (rs *SessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in redis_sentinel session
func (rs *SessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get redis_sentinel session id
func (rs *SessionStore) SessionID() string {
	return rs.sid
}

// SessionRelease save session values to redis_sentinel
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	b, err := session.EncodeGob(rs.values)
	if err != nil {
		return
	}
	c := rs.p
	c.Set(rs.sid, string(b), time.Duration(rs.maxlifetime)*time.Second)
}

// Provider redis_sentinel session provider
type Provider struct {
	maxlifetime int64
	savePath    string
	poolsize    int
	password    string
	dbNum       int
	poollist    *redis.Client
	masterName  string
}

// SessionInit init redis_sentinel session
// savepath like redis sentinel addr,pool size,password,dbnum,masterName
// e.g. 127.0.0.1:26379;127.0.0.2:26379,100,1qaz2wsx,0,mymaster
func (rp *Provider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	configs := strings.Split(savePath, ",")
	if len(configs) > 0 {
		rp.savePath = configs[0]
	}
	if len(configs) > 1 {
		poolsize, err := strconv.Atoi(configs[1])
		if err != nil || poolsize < 0 {
			rp.poolsize = DefaultPoolSize
		} else {
			rp.poolsize = poolsize
		}
	} else {
		rp.poolsize = DefaultPoolSize
	}
	if len(configs) > 2 {
		rp.password = configs[2]
	}
	if len(configs) > 3 {
		dbnum, err := strconv.Atoi(configs[3])
		if err != nil || dbnum < 0 {
			rp.dbNum = 0
		} else {
			rp.dbNum = dbnum
		}
	} else {
		rp.dbNum = 0
	}
	if len(configs) > 4 {
		if configs[4] != "" {
			rp.masterName = configs[4]
		} else {
			rp.masterName = "mymaster"
		}
	} else {
		rp.masterName = "mymaster"
	}

	rp.poollist = redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: strings.Split(rp.savePath, ";"),
		Password:      rp.password,
		PoolSize:      rp.poolsize,
		DB:            rp.dbNum,
		MasterName:    rp.masterName,
	})

	return rp.poollist.Ping().Err()
}

// SessionRead read redis_sentinel session by sid
func (rp *Provider) SessionRead(sid string) (session.Store, error) {
	var kv map[interface{}]interface{}
	kvs, err := rp.poollist.Get(sid).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		if kv, err = session.DecodeGob([]byte(kvs)); err != nil {
			return nil, err
		}
	}

	rs := &SessionStore{p: rp.poollist, sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionExist check redis_sentinel session exist by sid
func (rp *Provider) SessionExist(sid string) bool {
	c := rp.poollist
	if existed, err := c.Exists(sid).Result(); err != nil || existed == 0 {
		return false
	}
	return true
}

// SessionRegenerate generate new sid for redis_sentinel session
func (rp *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	c := rp.poollist

	if existed, err := c.Exists(oldsid).Result(); err != nil || existed == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
		c.Set(sid, "", time.Duration(rp.maxlifetime)*time.Second)
	} else {
		c.Rename(oldsid, sid)
		c.Expire(sid, time.Duration(rp.maxlifetime)*time.Second)
	}
	return rp.SessionRead(sid)
}

// SessionDestroy delete redis session by id
func (rp *Provider) SessionDestroy(sid string) error {
	c := rp.poollist
	c.Del(sid)
	return nil
}

// SessionGC Impelment method, no used.
func (rp *Provider) SessionGC() {
}

// SessionAll return all activeSession
func (rp *Provider) SessionAll() int {
	return 0
}

func init() {
	session.Register("redis_sentinel", redispder)
}
