package session

import (
	"github.com/garyburd/redigo/redis"
)

var redispder = &RedisProvider{}

var MAX_POOL_SIZE = 20

var redisPool chan redis.Conn

type RedisSessionStore struct {
	c   redis.Conn
	sid string
}

func (rs *RedisSessionStore) Set(key, value interface{}) error {
	//_, err := rs.c.Do("HSET", rs.sid, key, value)
	_, err := rs.c.Do("HSET", rs.sid, key, value)
	return err
}

func (rs *RedisSessionStore) Get(key interface{}) interface{} {
	reply, err := rs.c.Do("HGET", rs.sid, key)
	if err != nil {
		return nil
	}
	return reply
}

func (rs *RedisSessionStore) Delete(key interface{}) error {
	//_, err := rs.c.Do("HDEL", rs.sid, key)
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
	/*c, err := redis.Dial("tcp", rp.savePath)
	if err != nil {
		return nil
	}
	return c*/
	//if redisPool == nil {
	redisPool = make(chan redis.Conn, MAX_POOL_SIZE)
	//}
	if len(redisPool) == 0 {
		go func() {
			for i := 0; i < MAX_POOL_SIZE/2; i++ {
				c, err := redis.Dial("tcp", rp.savePath)
				if err != nil {
					panic(err)
				}
				putRedis(c)
			}
		}()
	}
	return <-redisPool
}

func putRedis(conn redis.Conn) {
	if redisPool == nil {
		redisPool = make(chan redis.Conn, MAX_POOL_SIZE)
	}
	if len(redisPool) >= MAX_POOL_SIZE {
		conn.Close()
		return
	}
	redisPool <- conn
}

func (rp *RedisProvider) SessionInit(maxlifetime int64, savePath string) error {
	rp.maxlifetime = maxlifetime
	rp.savePath = savePath
	return nil
}

func (rp *RedisProvider) SessionRead(sid string) (SessionStore, error) {
	c := rp.connectInit()
	//if str, err := redis.String(c.Do("GET", sid)); err != nil || str == "" {
	if str, err := redis.String(c.Do("HGET", sid, sid)); err != nil || str == "" {
		//c.Do("SET", sid, sid, rp.maxlifetime)
		c.Do("HSET", sid, sid, rp.maxlifetime)
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
