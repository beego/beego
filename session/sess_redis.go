package session

import (
	"strconv"
	"strings"
	"sync"

	"github.com/beego/redigo/redis"
)

var redispder = &RedisProvider{}

var MAX_POOL_SIZE = 100

var redisPool chan redis.Conn

type RedisSessionStore struct {
	c           redis.Conn
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
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
	if len(rs.values) > 0 {
		b, err := encodeGob(rs.values)
		if err != nil {
			return
		}
		rs.c.Do("SET", rs.sid, string(b))
		rs.c.Do("EXPIRE", rs.sid, rs.maxlifetime)
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
	if existed, err := redis.Int(c.Do("EXISTS", sid)); err != nil || existed == 0 {
		c.Do("SET", sid)
	}
	c.Do("EXPIRE", sid, rp.maxlifetime)
	kvs, err := redis.String(c.Do("GET", sid))
	var kv map[interface{}]interface{}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}
	rs := &RedisSessionStore{c: c, sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

func (rp *RedisProvider) SessionExist(sid string) bool {
	c := rp.poollist.Get()
	defer c.Close()
	if existed, err := redis.Int(c.Do("EXISTS", sid)); err != nil || existed == 0 {
		return false
	} else {
		return true
	}
}

func (rp *RedisProvider) SessionRegenerate(oldsid, sid string) (SessionStore, error) {
	c := rp.poollist.Get()
	if existed, err := redis.Int(c.Do("EXISTS", oldsid)); err != nil || existed == 0 {
		c.Do("SET", oldsid)
	}
	c.Do("RENAME", oldsid, sid)
	c.Do("EXPIRE", sid, rp.maxlifetime)
	kvs, err := redis.String(c.Do("GET", sid))
	var kv map[interface{}]interface{}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}
	rs := &RedisSessionStore{c: c, sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

func (rp *RedisProvider) SessionDestroy(sid string) error {
	c := rp.poollist.Get()
	defer c.Close()
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
