package session

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"fmt"
	"time"
	"sync"
	"net/http"
	"strings"
)

type redisProvider struct {
	maxLifeTime time.Duration
	redisClient *redis.Client
	generator   IdGeneration
	op          RedisOptions
}

var provider = redisProvider{}

const (
	STORE_PREFIX = "GO_SESSION"
)

type RedisOptions struct {
	Addr               string `json:"addr",omitempty`
	Password           string `json:"password",omitempty`
	DB                 int    `json:"db",omitempty`
	MaxRetries         int    `json:"max_retries",omitempty`
	MinRetryBackoff    int64  `json:"min_retry_backoff",omitempty`
	MaxRetryBackoff    int64  `json:"max_retry_backoff",omitempty`
	DialTimeout        int64  `json:"dial_timeout",omitempty`
	ReadTimeout        int64  `json:"read_timeout",omitempty`
	WriteTimeout       int64  `json:"write_timeout",omitempty`
	PoolSize           int    `json:"pool_size",omitempty`
	PoolTimeout        int64  `json:"pool_timeout",omitempty`
	IdleTimeout        int64  `json:"idle_timeout",omitempty`
	IdleCheckFrequency int64  `json:"idle_check_frequency",omitempty`
}

func (r *redisProvider) SessionInit(gclifetime int64, config string, generator IdGeneration) error {
	op := RedisOptions{}
	if err := json.Unmarshal([]byte(config), op); nil != err {
		fmt.Printf("json decode error %s , start config %s", err, config)
		return err
	}
	provider.generator = generator
	provider.op = op
	provider.maxLifeTime = time.Duration(gclifetime) * time.Second
	provider.redisClient = redis.NewClient(&redis.Options{
		Addr:               op.Addr,
		Password:           op.Password,
		DB:                 op.DB,
		MaxRetries:         op.MaxRetries,
		MinRetryBackoff:    time.Duration(op.MinRetryBackoff),
		MaxRetryBackoff:    time.Duration(op.MaxRetryBackoff),
		DialTimeout:        time.Duration(op.DialTimeout),
		ReadTimeout:        time.Duration(op.ReadTimeout),
		WriteTimeout:       time.Duration(op.WriteTimeout),
		PoolSize:           op.PoolSize,
		PoolTimeout:        time.Duration(op.PoolTimeout),
		IdleTimeout:        time.Duration(op.IdleTimeout),
		IdleCheckFrequency: time.Duration(op.IdleCheckFrequency),
	})
	_, err := provider.redisClient.Ping().Result()
	return err
}

func (r *redisProvider) SessionRead(rawSid string) (Store, error) {
	tsid, e := r.generator.GetSessionID(rawSid)
	if nil != e {
		return nil, e
	}
	var kvs map[interface{}]interface{}
	b, e := r.redisClient.Get(sessionStoreionKey(tsid)).Bytes()
	if kvs = decodeRedisSave(b); nil == kvs {
		kvs = make(map[interface{}]interface{})
	}
	rs := &SessionStore{redisClient: r.redisClient, sid: tsid, values: kvs, maxlifetime: r.maxLifeTime}
	return rs, nil
}

func (r *redisProvider) SessionExist(sid string) bool {
	tsid, e := r.generator.GetSessionID(sid)
	if nil != e {
		return false
	}
	if exists, err := r.redisClient.Exists(sessionStoreionKey(tsid)).Result(); nil != err || exists == 0 {
		return false
	}
	return true
}

func (r *redisProvider) SessionRegenerate(oldrawSid, rawSid string) (Store, error) {
	oldSid, e := r.generator.GetSessionID(oldrawSid)
	if nil != e {
		return nil, e
	}
	newSid, e := r.generator.GetSessionID(rawSid)
	if nil != e {
		return nil, e
	}
	oldKey := sessionStoreionKey(oldSid)
	newKey := sessionStoreionKey(newSid)
	if existed, _ := r.redisClient.Exists(oldKey).Result(); existed == 0 {
		r.redisClient.SetXX(newKey, "", r.maxLifeTime)
	} else {
		p := r.redisClient.Pipeline()
		defer func() {
			if nil != p {
				p.Close()
			}
		}()
		p.Rename(oldKey, newKey)
		p.Expire(newKey, r.maxLifeTime)
		p.Exec()
	}
	return r.SessionRead(newKey)
}

func (r *redisProvider) SessionDestroy(rawSid string) error {
	sid, e := r.generator.GetSessionID(rawSid)
	if nil != e {
		return e
	}
	_, e = r.redisClient.Del(sessionStoreionKey(sid)).Result()
	return e
}
func (r *redisProvider) SessionAll() int {
	return 0
}
func (r *redisProvider) SessionGC() {

}

// SessionStore redis session store
type SessionStore struct {
	redisClient *redis.Client
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime time.Duration
}

// Set value in redis session
func (rs *SessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in redis session
func (rs *SessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in redis session
func (rs *SessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in redis session
func (rs *SessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get redis session id
func (rs *SessionStore) SessionID() string {
	return rs.sid
}

// SessionRelease save session values to redis
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	b, err := EncodeGob(rs.values)
	if err != nil {
		return
	}
	rs.redisClient.SetXX(sessionStoreionKey(rs.sid), b, rs.maxlifetime)
}

func sessionStoreionKey(trueSid string) string {
	return strings.Join([]string{STORE_PREFIX, trueSid}, ":")
}

//if nil return nil
func decodeRedisSave(b []byte) (map[interface{}]interface{}) {
	if len(b) == 0 {
		return nil
	} else {
		kvs, e := DecodeGob(b)
		if nil != e {
			return nil
		}
		return kvs
	}
}

func init() {
	Register("redis", &provider)
}
