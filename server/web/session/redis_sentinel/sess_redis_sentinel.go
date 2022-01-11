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
//   _ "github.com/beego/beego/v2/server/web/session/redis_sentinel"
//   "github.com/beego/beego/v2/server/web/session"
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
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/beego/beego/v2/server/web/session"
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
func (rs *SessionStore) Set(ctx context.Context, key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in redis_sentinel session
func (rs *SessionStore) Get(ctx context.Context, key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in redis_sentinel session
func (rs *SessionStore) Delete(ctx context.Context, key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in redis_sentinel session
func (rs *SessionStore) Flush(context.Context) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get redis_sentinel session id
func (rs *SessionStore) SessionID(context.Context) string {
	return rs.sid
}

// SessionRelease save session values to redis_sentinel
func (rs *SessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
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
	MasterName            string `json:"master_name"`
}

// SessionInit init redis_sentinel session
// cfgStr like redis sentinel addr,pool size,password,dbnum,masterName
// e.g. 127.0.0.1:26379;127.0.0.2:26379,100,1qaz2wsx,0,mymaster
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

	rp.poollist = redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs:      strings.Split(rp.SavePath, ";"),
		Password:           rp.Password,
		PoolSize:           rp.Poolsize,
		DB:                 rp.DbNum,
		MasterName:         rp.MasterName,
		IdleTimeout:        rp.idleTimeout,
		IdleCheckFrequency: rp.idleCheckFrequency,
		MaxRetries:         rp.MaxRetries,
	})

	return rp.poollist.Ping().Err()
}

// for v1.x
func (rp *Provider) initOldStyle(savePath string) {
	configs := strings.Split(savePath, ",")
	if len(configs) > 0 {
		rp.SavePath = configs[0]
	}
	if len(configs) > 1 {
		poolsize, err := strconv.Atoi(configs[1])
		if err != nil || poolsize < 0 {
			rp.Poolsize = DefaultPoolSize
		} else {
			rp.Poolsize = poolsize
		}
	} else {
		rp.Poolsize = DefaultPoolSize
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
		if configs[4] != "" {
			rp.MasterName = configs[4]
		} else {
			rp.MasterName = "mymaster"
		}
	} else {
		rp.MasterName = "mymaster"
	}
	if len(configs) > 5 {
		timeout, err := strconv.Atoi(configs[4])
		if err == nil && timeout > 0 {
			rp.idleTimeout = time.Duration(timeout) * time.Second
		}
	}
	if len(configs) > 6 {
		checkFrequency, err := strconv.Atoi(configs[5])
		if err == nil && checkFrequency > 0 {
			rp.idleCheckFrequency = time.Duration(checkFrequency) * time.Second
		}
	}
	if len(configs) > 7 {
		retries, err := strconv.Atoi(configs[6])
		if err == nil && retries > 0 {
			rp.MaxRetries = retries
		}
	}
}

// SessionRead read redis_sentinel session by sid
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

// SessionExist check redis_sentinel session exist by sid
func (rp *Provider) SessionExist(ctx context.Context, sid string) (bool, error) {
	c := rp.poollist
	if existed, err := c.Exists(sid).Result(); err != nil || existed == 0 {
		return false, err
	}
	return true, nil
}

// SessionRegenerate generate new sid for redis_sentinel session
func (rp *Provider) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
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
	session.Register("redis_sentinel", redispder)
}
