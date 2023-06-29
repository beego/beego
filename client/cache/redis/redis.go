// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package redis for cache provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//
//	_ "github.com/beego/beego/v2/client/cache/redis"
//	"github.com/beego/beego/v2/client/cache"
//
// )
//
//	bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/berror"
)

const (
	// DefaultKey defines the collection name of redis for the cache adapter.
	DefaultKey = "beecacheRedis"
	// DefaultMaxIdle defines the default max idle connection number.
	DefaultMaxIdle = 3
	// DefaultTimeout defines the default timeout .
	DefaultTimeout = time.Second * 180
)

// Cache is Redis cache adapter.
type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	// key actually is prefix.
	key      string
	password string
	maxIdle  int

	// skipEmptyPrefix for backward compatible,
	// check function associate
	// see https://github.com/beego/beego/issues/5248
	skipEmptyPrefix bool

	// Timeout value (less than the redis server's timeout value)
	timeout time.Duration
}

// NewRedisCache creates a new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return &Cache{key: DefaultKey}
}

// Execute the redis commands. args[0] must be the key name
func (rc *Cache) do(commandName string, args ...interface{}) (interface{}, error) {
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()

	reply, err := c.Do(commandName, args...)
	if err != nil {
		return nil, berror.Wrapf(err, cache.RedisCacheCurdFailed,
			"could not execute this command: %s", commandName)
	}

	return reply, nil
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	if rc.key == "" && rc.skipEmptyPrefix {
		return fmt.Sprintf("%s", originKey)
	}
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	if v, err := rc.do("GET", key); err == nil {
		return v, nil
	} else {
		return nil, err
	}
}

// GetMulti gets cache from redis.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	return redis.Values(c.Do("MGET", args...))
}

// Put puts cache into redis.
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete deletes a key's cache in redis.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist checks cache's existence in redis.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return v, nil
}

// Incr increases a key's counter in redis.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decreases a key's counter in redis.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// ClearAll deletes all cache in the redis collection
// Be careful about this method, because it scans all keys and the delete them one by one
func (rc *Cache) ClearAll(context.Context) error {
	cachedKeys, err := rc.Scan(rc.key + ":*")
	if err != nil {
		return err
	}
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// Scan scans all keys matching a given pattern.
func (rc *Cache) Scan(pattern string) (keys []string, err error) {
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()
	var (
		cursor uint64 = 0 // start
		result []interface{}
		list   []string
	)
	for {
		result, err = redis.Values(c.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 1024))
		if err != nil {
			return
		}
		list, err = redis.Strings(result[1], nil)
		if err != nil {
			return
		}
		keys = append(keys, list...)
		cursor, err = redis.Uint64(result[0], nil)
		if err != nil {
			return
		}
		if cursor == 0 { // over
			return
		}
	}
}

// StartAndGC starts the redis cache adapter.
// config: must be in this format {"key":"collection key","conn":"connection info","dbNum":"0", "skipEmptyPrefix":"true"}
// Cached items in redis are stored forever, no garbage collection happens
func (rc *Cache) StartAndGC(config string) error {
	err := rc.parseConf(config)
	if err != nil {
		return err
	}

	rc.connectInit()

	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()

	// test connection
	if err = c.Err(); err != nil {
		return berror.Wrapf(err, cache.InvalidConnection,
			"can not connect to remote redis server, please check the connection info and network state: %s", config)
	}
	return nil
}

func (rc *Cache) parseConf(config string) error {
	var cf redisConfig
	err := json.Unmarshal([]byte(config), &cf)
	if err != nil {
		return berror.Wrapf(err, cache.InvalidRedisCacheCfg, "could not unmarshal the config: %s", config)
	}

	err = cf.parse()
	if err != nil {
		return err
	}

	rc.dbNum = cf.DbNum
	rc.key = cf.Key
	rc.conninfo = cf.Conn
	rc.password = cf.password
	rc.maxIdle = cf.MaxIdle
	rc.timeout = cf.timeout
	rc.skipEmptyPrefix = cf.SkipEmptyPrefix

	return nil
}

type redisConfig struct {
	DbNum           int    `json:"dbNum"`
	SkipEmptyPrefix bool   `json:"skipEmptyPrefix"`
	Key             string `json:"key"`
	// Format redis://<password>@<host>:<port>
	Conn       string `json:"conn"`
	MaxIdle    int    `json:"maxIdle"`
	TimeoutStr string `json:"timeout"`

	// parse from Conn
	password string
	// parse from TimeoutStr
	timeout time.Duration
}

// parse parses the config.
// If the necessary settings have not been set, it will return an error.
// It will fill the default values if some fields are missing.
func (cf *redisConfig) parse() error {
	if cf.Conn == "" {
		return berror.Error(cache.InvalidRedisCacheCfg, "config missing conn field")
	}

	if cf.Key == "" {
		cf.Key = DefaultKey
	}

	// Format redis://<password>@<host>:<port>
	cf.Conn = strings.Replace(cf.Conn, "redis://", "", 1)
	if i := strings.Index(cf.Conn, "@"); i > -1 {
		cf.password = cf.Conn[0:i]
		cf.Conn = cf.Conn[i+1:]
	}

	if cf.MaxIdle == 0 {
		cf.MaxIdle = DefaultMaxIdle
	}

	if v, err := time.ParseDuration(cf.TimeoutStr); err == nil {
		cf.timeout = v
	} else {
		cf.timeout = DefaultTimeout
	}

	return nil
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, berror.Wrapf(err, cache.DialFailed,
				"could not dial to remote server: %s ", rc.conninfo)
		}

		if rc.password != "" {
			if _, err = c.Do("AUTH", rc.password); err != nil {
				_ = c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			_ = c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: rc.timeout,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("redis", NewRedisCache)
}
