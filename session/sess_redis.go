package session

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"sync"
)

var redispder = &RedisProvider{}

var MAX_POOL_SIZE = 100

var redisPool chan redis.Conn

type RedisSessionStore struct {
	c      redis.Conn
	sid    string
	lock   sync.RWMutex
	values map[interface{}]interface{}
}

func (rs *RedisSessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

func (rs *RedisSessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (rs *RedisSessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

func (rs *RedisSessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

func (rs *RedisSessionStore) SessionID() string {
	return rs.sid
}

func (rs *RedisSessionStore) SessionRelease() {
	defer rs.c.Close()
	keys, err := redis.Values(rs.c.Do("HKEYS", rs.sid))
	if err == nil {
		for _, key := range keys {
			if val, ok := rs.values[key]; ok {
				rs.c.Do("HSET", rs.sid, key, val)
				rs.Delete(key)
			} else {
				rs.c.Do("HDEL", rs.sid, key)
			}
		}
	}
	if len(rs.values) > 0 {
		for k, v := range rs.values {
			rs.c.Do("HSET", rs.sid, k, v)
		}
	}
}

type RedisProvider struct {
	maxlifetime int64
	savePath    string
	poolsize    int
	password    string
	poollist    *redis.Pool
}

//savepath like redisserveraddr,poolsize,password
//127.0.0.1:6379,100,astaxie
func (rp *RedisProvider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	configs := strings.Split(savePath, ",")
	if len(configs) > 0 {
		rp.savePath = configs[0]
	}
	if len(configs) > 1 {
		poolsize, err := strconv.Atoi(configs[1])
		if err != nil || poolsize <= 0 {
			rp.poolsize = MAX_POOL_SIZE
		} else {
			rp.poolsize = poolsize
		}
	} else {
		rp.poolsize = MAX_POOL_SIZE
	}
	if len(configs) > 2 {
		rp.password = configs[2]
	}
	rp.poollist = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", rp.savePath)
		if err != nil {
			return nil, err
		}
		if rp.password != "" {
			if _, err := c.Do("AUTH", rp.password); err != nil {
				c.Close()
				return nil, err
			}
		}
		return c, err
	}, rp.poolsize)
	return nil
}

func (rp *RedisProvider) SessionRead(sid string) (SessionStore, error) {
	c := rp.poollist.Get()
	//if str, err := redis.String(c.Do("GET", sid)); err != nil || str == "" {
	if str, err := redis.String(c.Do("HGET", sid, sid)); err != nil || str == "" {
		//c.Do("SET", sid, sid, rp.maxlifetime)
		c.Do("HSET", sid, sid, rp.maxlifetime)
	}
	c.Do("EXPIRE", sid, rp.maxlifetime)
	kvs, err := redis.Values(c.Do("HGETALL", sid))
	vals := make(map[interface{}]interface{})
	var key interface{}
	if err == nil {
		for k, v := range kvs {
			if k%2 == 0 {
				key = v
			} else {
				vals[key] = v
			}
		}
	}
	rs := &RedisSessionStore{c: c, sid: sid, values: vals}
	return rs, nil
}

func (rp *RedisProvider) SessionRegenerate(oldsid, sid string) (SessionStore, error) {
	c := rp.poollist.Get()
	if str, err := redis.String(c.Do("HGET", oldsid, oldsid)); err != nil || str == "" {
		c.Do("HSET", oldsid, oldsid, rp.maxlifetime)
	}
	c.Do("RENAME", oldsid, sid)
	c.Do("EXPIRE", sid, rp.maxlifetime)
	kvs, err := redis.Values(c.Do("HGETALL", sid))
	vals := make(map[interface{}]interface{})
	var key interface{}
	if err == nil {
		for k, v := range kvs {
			if k%2 == 0 {
				key = v
			} else {
				vals[key] = v
			}
		}
	}
	rs := &RedisSessionStore{c: c, sid: sid, values: vals}
	return rs, nil
}

func (rp *RedisProvider) SessionDestroy(sid string) error {
	c := rp.poollist.Get()
	c.Do("DEL", sid)
	return nil
}

func (rp *RedisProvider) SessionGC() {
	return
}

//@todo
func (rp *RedisProvider) SessionAll() int {

	return 0
}

func init() {
	Register("redis", redispder)
}
