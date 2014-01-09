package cache

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/beego/redigo/redis"
)

var (
	// the collection name of redis for cache adapter.
	DefaultKey string = "beecacheRedis"
)

// Redis cache adapter.
type RedisCache struct {
	c        redis.Conn
	conninfo string
	key      string
}

// create new redis cache with default collection name.
func NewRedisCache() *RedisCache {
	return &RedisCache{key: DefaultKey}
}

// Get cache from redis.
func (rc *RedisCache) Get(key string) interface{} {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return nil
		}
	}

	v, err := rc.c.Do("HGET", rc.key, key)
	// write to closed socket, reset rc.c to nil
	if err == io.EOF {
		rc.c = nil
		return nil
	}

	if err != nil {
		return nil
	}

	return v
}

// put cache to redis.
// timeout is ignored.
func (rc *RedisCache) Put(key string, val interface{}, timeout int64) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}

	_, err := rc.c.Do("HSET", rc.key, key, val)
	// write to closed socket, reset rc.c to nil
	if err == io.EOF {
		rc.c = nil
		return err
	}

	return err
}

// delete cache in redis.
func (rc *RedisCache) Delete(key string) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}

	_, err := rc.c.Do("HDEL", rc.key, key)
	// write to closed socket, reset rc.c to nil
	if err == io.EOF {
		rc.c = nil
		return err
	}

	return err
}

// check cache exist in redis.
func (rc *RedisCache) IsExist(key string) bool {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return false
		}
	}

	v, err := redis.Bool(rc.c.Do("HEXISTS", rc.key, key))
	// write to closed socket, reset rc.c to nil
	if err == io.EOF {
		rc.c = nil
		return false
	}

	if err != nil {
		return false
	}

	return v
}

// increase counter in redis.
func (rc *RedisCache) Incr(key string) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}

	_, err := redis.Bool(rc.c.Do("HINCRBY", rc.key, key, 1))
	// write to closed socket
	if err == io.EOF {
		rc.c = nil
	}

	return err
}

// decrease counter in redis.
func (rc *RedisCache) Decr(key string) error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}

	_, err := redis.Bool(rc.c.Do("HINCRBY", rc.key, key, -1))

	// write to closed socket
	if err == io.EOF {
		rc.c = nil
	}

	return err
}

// clean all cache in redis. delete this redis collection.
func (rc *RedisCache) ClearAll() error {
	if rc.c == nil {
		var err error
		rc.c, err = rc.connectInit()
		if err != nil {
			return err
		}
	}

	_, err := rc.c.Do("DEL", rc.key)
	// write to closed socket
	if err == io.EOF {
		rc.c = nil
	}

	return err
}

// start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info"}
// the cache item in redis are stored forever,
// so no gc operation.
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

// connect to redis.
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
