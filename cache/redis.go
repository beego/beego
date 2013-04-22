package cache

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
)

var (
	DefaultKey string = "beecacheRedis"
)

type RedisCache struct {
	c        redis.Conn
	conninfo string
	key      string
}

func NewRedisCache() *RedisCache {
	return &RedisCache{key: DefaultKey}
}

func (rc *RedisCache) Get(key string) interface{} {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := rc.c.Do("HGET", rc.key, key)
	if err != nil {
		return nil
	}
	return v
}

func (rc *RedisCache) Put(key string, val interface{}, timeout int) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Do("HSET", rc.key, key, val)
	return err
}

func (rc *RedisCache) Delete(key string) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Do("HDEL", rc.key, key)
	return err
}

func (rc *RedisCache) IsExist(key string) bool {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := redis.Bool(rc.c.Do("HEXISTS", rc.key, key))
	if err != nil {
		return false
	}
	return v
}

func (rc *RedisCache) ClearAll() error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Do("DEL", rc.key)
	return err
}

func (rc *RedisCache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.c = rc.connectInit()
	if rc.c == nil {
		return errors.New("dial tcp conn error")
	}
	return nil
}

func (rc *RedisCache) connectInit() redis.Conn {
	c, err := redis.Dial("tcp", rc.conninfo)
	if err != nil {
		return nil
	}
	return c
}

func init() {
	Register("redis", NewRedisCache())
}
