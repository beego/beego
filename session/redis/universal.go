// Arthur xie <arthur.xie@unosys.io>
package redis

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/astaxie/beego/session"
	"github.com/go-redis/redis"
)

const (
	versionKey = "$version"
)

// SessionStore redis session store
type RedisUniversalSessionStore struct {
	client      redis.UniversalClient
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime time.Duration
	changed     bool
}

// Set value in redis session
func (rs *RedisUniversalSessionStore) Set(key, value interface{}) error {
	if key == versionKey {
		return nil
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	rs.changed = true
	return nil
}

// Get value in redis session
func (rs *RedisUniversalSessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in redis session
func (rs *RedisUniversalSessionStore) Delete(key interface{}) error {
	if key == versionKey {
		return nil
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	rs.changed = true
	return nil
}

// Flush clear all values in redis session
func (rs *RedisUniversalSessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	version := rs.values[versionKey]
	rs.values = make(map[interface{}]interface{})
	rs.values[versionKey] = version
	rs.changed = true
	return nil
}

// SessionID get redis session id
func (rs *RedisUniversalSessionStore) SessionID() string {
	return rs.sid
}

func getVersion(values map[interface{}]interface{}) int64 {
	v, _ := values[versionKey].(int64)
	return v
}

// SessionRelease save session values to redis
func (rs *RedisUniversalSessionStore) SessionRelease(w http.ResponseWriter) {
	if rs.changed {
		version := getVersion(rs.values)
		rs.values[versionKey] = version + 1

		kvs, err := rs.client.Get(rs.sid).Result()
		if len(kvs) != 0 {
			kv, err := session.DecodeGob([]byte(kvs))
			if err == nil && getVersion(kv) > version {
				for k, v := range kv {
					rs.values[k] = v
				}
				rs.values[versionKey] = getVersion(kv) + 1
			}
		}

		b, err := session.EncodeGob(rs.values)
		if err != nil {
			return
		}

		rs.client.Set(rs.sid, string(b), rs.maxlifetime)
		rs.changed = false
	} else {
		rs.client.Expire(rs.sid, rs.maxlifetime)
	}
}

// Provider redis session provider
type RedisUniversalProvider struct {
	maxlifetime time.Duration
	options     redis.UniversalOptions
	client      redis.UniversalClient
}

// SessionInit init redis session
// savepath MUST be an UniversalOptions(http://godoc.org/github.com/go-redis/redis#UniversalOptions)
// json. e.g. {"Addrs": ["localhost:6379"], "DB": 0}
func (rp *RedisUniversalProvider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = time.Duration(maxlifetime) * time.Second
	err := json.Unmarshal([]byte(savePath), &rp.options)
	if err != nil {
		return err
	}
	rp.client = redis.NewUniversalClient(&rp.options)
	return nil
}

// SessionRead read redis session by sid
func (rp *RedisUniversalProvider) SessionRead(sid string) (session.Store, error) {
	changed := false
	kvs, err := rp.client.Get(sid).Result()
	var kv map[interface{}]interface{}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
		if err == redis.Nil {
			//Treat new one as changed
			changed = true
		}
	} else {
		kv, err = session.DecodeGob([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}

	rs := &RedisUniversalSessionStore{
		client:      rp.client,
		sid:         sid,
		values:      kv,
		changed:     changed,
		maxlifetime: rp.maxlifetime,
	}
	return rs, nil
}

// SessionExist check redis session exist by sid
func (rp *RedisUniversalProvider) SessionExist(sid string) bool {
	val, err := rp.client.Exists(sid).Result()
	if err != nil {
		return false
	} else {
		return val != 0
	}
}

// SessionRegenerate generate new sid for redis session
func (rp *RedisUniversalProvider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	if existed, _ := rp.client.Exists(oldsid).Result(); existed == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
		rp.client.Set(sid, "", rp.maxlifetime)
	} else {
		rp.client.Rename(oldsid, sid)
		rp.client.Expire(sid, rp.maxlifetime)
	}

	changed := false
	kvs, err := rp.client.Get(sid).Result()
	var kv map[interface{}]interface{}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
		if err == redis.Nil {
			//Treat new one as changed
			changed = true
		}
	} else {
		kv, err = session.DecodeGob([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}

	rs := &RedisUniversalSessionStore{
		client:      rp.client,
		sid:         sid,
		values:      kv,
		changed:     changed,
		maxlifetime: rp.maxlifetime,
	}
	return rs, nil
}

// SessionDestroy delete redis session by id
func (rp *RedisUniversalProvider) SessionDestroy(sid string) error {
	return rp.client.Del(sid).Err()
}

// SessionGC Impelment method, no used.
func (rp *RedisUniversalProvider) SessionGC() {
	return
}

// SessionAll return all activeSession
func (rp *RedisUniversalProvider) SessionAll() int {
	return 0
}

func init() {
	session.Register("redis_universal", &RedisUniversalProvider{})
}
