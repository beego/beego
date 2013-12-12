package cache

import (
	"encoding/json"
	"errors"

	"github.com/beego/redigo/redis"
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
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return nil
		}
	}
	v, err := rc.c.Do("HGET", rc.key, key)
	if err != nil {
		return nil
	}
	return v
}

func (rc *RedisCache) GetString(key string) (string, bool) {
	var contain string

	data := rc.Get(key)
	if data == nil {
		return contain, false
	}

	if d, ok := data.([]byte); ok {
		contain = string(d)
	}
	return contain, true
}

func (rc *RedisCache) Put(key string, val interface{}, timeout int64) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := rc.c.Do("HSET", rc.key, key, val)
	return err
}

func (rc *RedisCache) Delete(key string) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := rc.c.Do("HDEL", rc.key, key)
	return err
}

func (rc *RedisCache) IsExist(key string) bool {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return false
		}
	}
	v, err := redis.Bool(rc.c.Do("HEXISTS", rc.key, key))
	if err != nil {
		return false
	}
	return v
}

func (rc *RedisCache) Incr(key string) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := redis.Bool(rc.c.Do("HINCRBY", rc.key, key, 1))
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisCache) Decr(key string) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}
	_, err := redis.Bool(rc.c.Do("HINCRBY", rc.key, key, -1))
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisCache) ClearAll() error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
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
	var err error
	rc.c, err = rc.connectInit()
	if err != nil {
		return err
	}
	if rc.c == nil {
		return errors.New("dial tcp conn error")
	}
	return nil
}

func (rc *RedisCache) connectInit() (redis.Conn, error) {
	c, err := redis.Dial("tcp", rc.conninfo)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func init() {
	Register("redis", NewRedisCache())
}
