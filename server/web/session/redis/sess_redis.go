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
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/session/redis"
//   "github.com/astaxie/beego/session"
// )
//
// 	func init() {
// 		globalSessions, _ = session.NewManager("redis", ``{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:7070"}``)
// 		go globalSessions.GC()
// 	}
//
// more docs: http://beego.me/docs/module/session.md
package redis

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/astaxie/beego/server/web/session"
)

var redispder = &Provider{}

// MaxPoolSize redis max pool size
var MaxPoolSize = 100

// SessionStore redis session store
type SessionStore struct {
	p           *redis.Client
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Set value in redis session
func (rs *SessionStore) Set(ctx context.Context, key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in redis session
func (rs *SessionStore) Get(ctx context.Context, key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in redis session
func (rs *SessionStore) Delete(ctx context.Context, key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in redis session
func (rs *SessionStore) Flush(context.Context) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get redis session id
func (rs *SessionStore) SessionID(context.Context) string {
	return rs.sid
}

// SessionRelease save session values to redis
func (rs *SessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	b, err := session.EncodeGob(rs.values)
	if err != nil {
		return
	}
	c := rs.p
	c.Set(rs.sid, string(b), time.Duration(rs.maxlifetime)*time.Second)
}

// Provider redis session provider
type Provider struct {
	maxlifetime int64
	SavePath    string `json:"save_path"`
	Poolsize    int    `json:"poolsize"`
	Password    string `json:"password"`
	DbNum       int    `json:"db_num"`

	idleTimeout    time.Duration
	IdleTimeoutStr string `json:"idle_timeout"`

	idleCheckFrequency    time.Duration
	IdleCheckFrequencyStr string `json:"idle_check_frequency"`
	MaxRetries            int    `json:"max_retries"`
	poollist              *redis.Client
}

// SessionInit init redis session
// savepath like redis server addr,pool size,password,dbnum,IdleTimeout second
// v1.x e.g. 127.0.0.1:6379,100,astaxie,0,30
// v2.0 you should pass json string
func (rp *Provider) SessionInit(ctx context.Context, maxlifetime int64, cfgStr string) error {
	rp.maxlifetime = maxlifetime

	cfgStr = strings.TrimSpace(cfgStr)
	// we think cfgStr is v2.0, using json to init the session
	if strings.HasPrefix(cfgStr, "{") {
		err := json.Unmarshal([]byte(cfgStr), rp)
		if err != nil {
			return err
		}
		rp.idleTimeout, err = time.ParseDuration(rp.IdleTimeoutStr)
		if err != nil {
			return err
		}

		rp.idleCheckFrequency, err = time.ParseDuration(rp.IdleCheckFrequencyStr)
		if err != nil {
			return err
		}

	} else {
		rp.initOldStyle(cfgStr)
	}

	rp.poollist = redis.NewClient(&redis.Options{
		Addr:               rp.SavePath,
		Password:           rp.Password,
		PoolSize:           rp.Poolsize,
		DB:                 rp.DbNum,
		IdleTimeout:        rp.idleTimeout,
		IdleCheckFrequency: rp.idleCheckFrequency,
		MaxRetries:         rp.MaxRetries,
	})

	return rp.poollist.Ping().Err()
}

func (rp *Provider) initOldStyle(savePath string) {
	configs := strings.Split(savePath, ",")
	if len(configs) > 0 {
		rp.SavePath = configs[0]
	}
	if len(configs) > 1 {
		poolsize, err := strconv.Atoi(configs[1])
		if err != nil || poolsize < 0 {
			rp.Poolsize = MaxPoolSize
		} else {
			rp.Poolsize = poolsize
		}
	} else {
		rp.Poolsize = MaxPoolSize
	}
	if len(configs) > 2 {
		rp.Password = configs[2]
	}
	if len(configs) > 3 {
		dbnum, err := strconv.Atoi(configs[3])
		if err != nil || dbnum < 0 {
			rp.DbNum = 0
		} else {
			rp.DbNum = dbnum
		}
	} else {
		rp.DbNum = 0
	}
	if len(configs) > 4 {
		timeout, err := strconv.Atoi(configs[4])
		if err == nil && timeout > 0 {
			rp.idleTimeout = time.Duration(timeout) * time.Second
		}
	}
	if len(configs) > 5 {
		checkFrequency, err := strconv.Atoi(configs[5])
		if err == nil && checkFrequency > 0 {
			rp.idleCheckFrequency = time.Duration(checkFrequency) * time.Second
		}
	}
	if len(configs) > 6 {
		retries, err := strconv.Atoi(configs[6])
		if err == nil && retries > 0 {
			rp.MaxRetries = retries
		}
	}
}

// SessionRead read redis session by sid
func (rp *Provider) SessionRead(ctx context.Context, sid string) (session.Store, error) {
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

// SessionExist check redis session exist by sid
func (rp *Provider) SessionExist(ctx context.Context, sid string) (bool, error) {
	c := rp.poollist

	if existed, err := c.Exists(sid).Result(); err != nil || existed == 0 {
		return false, err
	}
	return true, nil
}

// SessionRegenerate generate new sid for redis session
func (rp *Provider) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	c := rp.poollist
	if existed, _ := c.Exists(oldsid).Result(); existed == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
		c.Do(c.Context(), "SET", sid, "", "EX", rp.maxlifetime)
	} else {
		c.Rename(oldsid, sid)
		c.Expire(sid, time.Duration(rp.maxlifetime)*time.Second)
	}
	return rp.SessionRead(context.Background(), sid)
}

// SessionDestroy delete redis session by id
func (rp *Provider) SessionDestroy(ctx context.Context, sid string) error {
	c := rp.poollist

	c.Del(sid)
	return nil
}

// SessionGC Impelment method, no used.
func (rp *Provider) SessionGC(context.Context) {
}

// SessionAll return all activeSession
func (rp *Provider) SessionAll(context.Context) int {
	return 0
}

func init() {
	session.Register("redis", redispder)
}
