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

	"github.com/beego/beego/client/cache"
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
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, err
		}
	}
	value, err := rc.conn.Get(key)
	if err == nil {
		return value, nil
	}
	return nil, err
}

// GetMulti gets one or keys values from ssdb.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	size := len(keys)
	values := make([]interface{}, size)
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return values, err
		}
	}

	res, err := rc.conn.Do("multi_get", keys)
	if err != nil {
		return values, err
	}

	resSize := len(res)
	keyIdx := make(map[string]int)
	for i := 1; i < resSize; i += 2 {
		keyIdx[res[i]] = i
	}

	keysErr := make([]string, 0)
	for i, ki := range keys {
		if _, ok := keyIdx[ki]; !ok {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, "the key isn't exist"))
			continue
		}
		values[i] = res[keyIdx[ki]+1]
	}

	if len(keysErr) != 0 {
		return values, fmt.Errorf(strings.Join(keysErr, "; "))
	}

	return values, nil
}

// DelMulti deletes one or more keys from memcache
func (rc *Cache) DelMulti(keys []string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("multi_del", keys)
	return err
}

// Put puts value into memcache.
// value:  must be of type string
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
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

// Delete deletes a value in memcache.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Del(key)
	return err
}

// Incr increases a key's counter.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("incr", key, 1)
	return err
}

// Decr decrements a key's counter.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("incr", key, -1)
	return err
}

// IsExist checks if a key exists in memcache.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
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
func (rc *Cache) ClearAll(context.Context) error {
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
