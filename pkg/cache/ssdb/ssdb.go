package ssdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ssdb/gossdb/ssdb"

	"github.com/astaxie/beego/pkg/cache"
)

// Cache SSDB adapter
type Cache struct {
	conn     *ssdb.Client
	conninfo []string
}

//NewSsdbCache creates new ssdb adapter.
func NewSsdbCache() cache.Cache {
	return &Cache{}
}

// Get gets a key's value from memcache.
func (rc *Cache) Get(key string) (interface{}, error) {
	return rc.GetWithCtx(context.Background(), key)
}

// GetWithCtx gets a key's value from memcache.
func (rc *Cache) GetWithCtx(ctx context.Context, key string) (interface{}, error) {
	rc.conn.Do()
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, nil
		}
	}
	return rc.conn.Get(key)
}

// GetMulti gets one or keys values from memcache.
func (rc *Cache) GetMulti(keys []string) ([]interface{}, error) {
	return rc.GetMultiWithCtx(context.Background(), keys)
}

// GetMultiWithCtx gets one or keys values from memcache.
func (rc *Cache) GetMultiWithCtx(ctx context.Context, keys []string) ([]interface{}, error) {
	size := len(keys)
	var values []interface{}
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			for i := 0; i < size; i++ {
				values = append(values, nil)
			}
			return values, err
		}
	}
	res, err := rc.conn.Do("multi_get", keys)
	resSize := len(res)
	if err == nil {
		for i := 0; i < size; i++ {
			values = append(values, nil)
		}
		return values, err
	}
	for i := 1; i < resSize; i += 2 {
		values = append(values, res[i+1])
	}
	return values, nil
}

// Put puts value into memcache.
// value:  must be of type string
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	return rc.PutWithCtx(context.Background(), key, val, timeout)
}

//PutWithCtx ..
func (rc *Cache) PutWithCtx(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	v, ok := val.(string)
	if !ok {
		return errors.New("value must string")
	}
	var resp []string
	var err error
	ttl := int(timeout / time.Second)
	if ttl < 0 {
		resp, err = rc.conn.Do("set", key, v)
	} else {
		resp, err = rc.conn.Do("setx", key, v, ttl)
	}
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return errors.New("bad response")
}

// Delete cached value by key.
func (rc *Cache) Delete(key string) error {
	return rc.DeleteWithCtx(context.Background(), key)
}

// DeleteWithCtx deletes a value in memcache.
func (rc *Cache) DeleteWithCtx(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Del(key)
	return err
}

// IncrBy increases a key's counter.
func (rc *Cache) IncrBy(key string, n int) (int, error) {
	return rc.IncrByWithCtx(context.Background(), key, n)
}

//IncrByWithCtx increases a key's counter.
func (rc *Cache) IncrByWithCtx(ctx context.Context, key string, n int) (int, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return 0, err
		}
	}
	resp, err := rc.conn.Do("incr", key, n)
	if err != nil {
		return 0, err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return strconv.Atoi(resp[1])
	}
	if resp[0] == "not_found" {
		return 0, fmt.Errorf("not_found")
	}
	return 0, fmt.Errorf(resp[0])
}

// Incr a cached int value by key, as a counter.
// int indicates current value after increasing
func (rc *Cache) Incr(key string) (int, error) {
	return rc.IncrWithCtx(context.Background(), key)
}

// IncrWithCtx ..
func (rc *Cache) IncrWithCtx(ctx context.Context, key string) (int, error) {
	return rc.IncrByWithCtx(context.Background(), key, 1)
}

// Decr a cached int value by key, as a counter.
// int indicates current value after decreasing
func (rc *Cache) Decr(key string) (int, error) {
	return rc.DecrWithCtx(context.Background(), key)
}

// DecrWithCtx ..
func (rc *Cache) DecrWithCtx(ctx context.Context, key string) (int, error) {
	return rc.IncrByWithCtx(context.Background(), key, -1)
}

// IsExist checks if a key exists
func (rc *Cache) IsExist(key string) (bool, error) {
	return rc.IsExistWithCtx(context.Background(), key)
}

// IsExistWithCtx checks if a key exists
func (rc *Cache) IsExistWithCtx(ctx context.Context, key string) (bool, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false, err
		}
	}
	resp, err := rc.conn.Do("exists", key)
	if err != nil {
		return false, err
	}
	if len(resp) == 2 && resp[1] == "1" {
		return true, nil
	}
	return false, nil
}

// ClearAll clears all cached items in memcache.
func (rc *Cache) ClearAll() error {
	return rc.ClearAllWithCtx(context.Background())
}

//ClearAllWithCtx ..
func (rc *Cache) ClearAllWithCtx(ctx context.Context) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	keyStart, keyEnd, limit := "", "", 50
	resp, err := rc.Scan(keyStart, keyEnd, limit)
	for err == nil {
		size := len(resp)
		if size == 1 {
			return nil
		}
		keys := []string{}
		for i := 1; i < size; i += 2 {
			keys = append(keys, resp[i])
		}
		_, e := rc.conn.Do("multi_del", keys)
		if e != nil {
			return e
		}
		keyStart = resp[size-2]
		resp, err = rc.Scan(keyStart, keyEnd, limit)
	}
	return err
}

// Scan key all cached in ssdb.
func (rc *Cache) Scan(keyStart string, keyEnd string, limit int) ([]string, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, err
		}
	}
	resp, err := rc.conn.Do("scan", keyStart, keyEnd, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// StartAndGC starts the memcache adapter.
// config: must be in the format {"conn":"connection info"}.
// If an error occurs during connection, an error is returned
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *Cache) connectInit() error {
	conninfoArray := strings.Split(rc.conninfo[0], ":")
	host := conninfoArray[0]
	port, e := strconv.Atoi(conninfoArray[1])
	if e != nil {
		return e
	}
	var err error
	rc.conn, err = ssdb.Connect(host, port)
	return err
}

func init() {
	cache.Register("ssdb", NewSsdbCache)
}
