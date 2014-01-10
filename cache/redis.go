package cache

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/beego/redigo/redis"
)

var (
	// the collection name of redis for cache adapter.
	DefaultKey string = "beecacheRedis"
)

// Redis cache adapter.
type RedisCache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	key      string
	mu       sync.Mutex
}

// create new redis cache with default collection name.
func NewRedisCache() *RedisCache {
	return &RedisCache{key: DefaultKey}
}

// actually do the redis cmds
func (rc *RedisCache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// Get cache from redis.
func (rc *RedisCache) Get(key string) interface{} {
	v, err := rc.do("HGET", rc.key, key)
	if err != nil {
		return nil
	}

	return v
}

// put cache to redis.
// timeout is ignored.
func (rc *RedisCache) Put(key string, val interface{}, timeout int64) error {
	_, err := rc.do("HSET", rc.key, key, val)
	return err
}

// delete cache in redis.
func (rc *RedisCache) Delete(key string) error {
	_, err := rc.do("HDEL", rc.key, key)
	return err
}

// check cache exist in redis.
func (rc *RedisCache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("HEXISTS", rc.key, key))
	if err != nil {
		return false
	}

	return v
}

// increase counter in redis.
func (rc *RedisCache) Incr(key string) error {
	_, err := redis.Bool(rc.do("HINCRBY", rc.key, key, 1))
	return err
}

// decrease counter in redis.
func (rc *RedisCache) Decr(key string) error {
	_, err := redis.Bool(rc.do("HINCRBY", rc.key, key, -1))
	return err
}

// clean all cache in redis. delete this redis collection.
func (rc *RedisCache) ClearAll() error {
	_, err := rc.do("DEL", rc.key)
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
	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()
	if err := c.Err(); err != nil {
		return err
	}

	return nil
}

// connect to redis.
func (rc *RedisCache) connectInit() {
	rc.mu.Lock()

	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rc.conninfo)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	rc.mu.Unlock()
}

func init() {
	Register("redis", NewRedisCache())
}
