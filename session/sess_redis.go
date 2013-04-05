package session

import (
	"github.com/garyburd/redigo/redis"
)

var redispder = &RedisProvider{}

type RedisSessionStore struct {
	c   redis.Conn
	sid string
}

func (rs *RedisSessionStore) Set(key, value interface{}) error {
	_, err := rs.c.Do("HSET", rs.sid, key, value)
	return err
}

func (rs *RedisSessionStore) Get(key interface{}) interface{} {
	v, err := rs.c.Do("GET", rs.sid, key)
	if err != nil {
		return nil
	}
	return v
}

func (rs *RedisSessionStore) Delete(key interface{}) error {
	_, err := rs.c.Do("HDEL", rs.sid, key)
	return err
}

func (rs *RedisSessionStore) SessionID() string {
	return rs.sid
}

func (rs *RedisSessionStore) SessionRelease() {
	rs.c.Close()
}

type RedisProvider struct {
	maxlifetime int64
	savePath    string
}

func (rp *RedisProvider) connectInit() redis.Conn {
	c, err := redis.Dial("tcp", rp.savePath)
	if err != nil {
		return nil
	}
	return c
}

func (rp *RedisProvider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	rp.savePath = savePath
	return nil
}

func (rp *RedisProvider) SessionRead(sid string) (SessionStore, error) {
	c := rp.connectInit()
	if str, err := redis.String(c.Do("GET", sid)); err != nil || str == "" {
		c.Do("SET", sid, sid, rp.maxlifetime)
	}
	rs := &RedisSessionStore{c: c, sid: sid}
	return rs, nil
}

func (rp *RedisProvider) SessionDestroy(sid string) error {
	c := rp.connectInit()
	c.Do("DEL", sid)
	return nil
}

func (rp *RedisProvider) SessionGC() {
	return
}

func init() {
	Register("redis", redispder)
}
