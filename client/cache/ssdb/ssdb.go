package ssdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ssdb/gossdb/ssdb"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/berror"
)

// Cache SSDB adapter
type Cache struct {
	conn     *ssdb.Client
	conninfo []string
}

// NewSsdbCache creates new ssdb adapter.
func NewSsdbCache() cache.Cache {
	return &Cache{}
}

// Get gets a key's value from memcache.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	value, err := rc.conn.Get(key)
	if err == nil {
		return value, nil
	}
	return nil, berror.Wrapf(err, cache.SsdbCacheCurdFailed, "could not get value, key: %s", key)
}

// GetMulti gets one or keys values from ssdb.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	size := len(keys)
	values := make([]interface{}, size)

	res, err := rc.conn.Do("multi_get", keys)
	if err != nil {
		return values, berror.Wrapf(err, cache.SsdbCacheCurdFailed, "multi_get failed, key: %v", keys)
	}

	resSize := len(res)
	keyIdx := make(map[string]int)
	for i := 1; i < resSize; i += 2 {
		keyIdx[res[i]] = i
	}

	keysErr := make([]string, 0)
	for i, ki := range keys {
		if _, ok := keyIdx[ki]; !ok {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, "key not exist"))
			continue
		}
		values[i] = res[keyIdx[ki]+1]
	}

	if len(keysErr) != 0 {
		return values, berror.Error(cache.MultiGetFailed, strings.Join(keysErr, "; "))
	}

	return values, nil
}

// DelMulti deletes one or more keys from memcache
func (rc *Cache) DelMulti(keys []string) error {
	_, err := rc.conn.Do("multi_del", keys)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "multi_del failed: %v", keys)
}

// Put puts value into memcache.
// value:  must be of type string
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	v, ok := val.(string)
	if !ok {
		return berror.Errorf(cache.InvalidSsdbCacheValue, "value must be string: %v", val)
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
		return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "set or setx failed, key: %s", key)
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return berror.Errorf(cache.SsdbBadResponse, "the response from SSDB server is invalid: %v", resp)
}

// Delete deletes a value in memcache.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	_, err := rc.conn.Del(key)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "del failed: %s", key)
}

// Incr increases a key's counter.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	_, err := rc.conn.Do("incr", key, 1)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "increase failed: %s", key)
}

// Decr decrements a key's counter.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	_, err := rc.conn.Do("incr", key, -1)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "decrease failed: %s", key)
}

// IsExist checks if a key exists in memcache.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	resp, err := rc.conn.Do("exists", key)
	if err != nil {
		return false, berror.Wrapf(err, cache.SsdbCacheCurdFailed, "exists failed: %s", key)
	}
	if len(resp) == 2 && resp[1] == "1" {
		return true, nil
	}
	return false, nil

}

// ClearAll clears all cached items in ssdb.
// If there are many keys, this method may spent much time.
func (rc *Cache) ClearAll(context.Context) error {
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
			return berror.Wrapf(e, cache.SsdbCacheCurdFailed, "multi_del failed: %v", keys)
		}
		keyStart = resp[size-2]
		resp, err = rc.Scan(keyStart, keyEnd, limit)
	}
	return berror.Wrap(err, cache.SsdbCacheCurdFailed, "scan failed")
}

// Scan key all cached in ssdb.
func (rc *Cache) Scan(keyStart string, keyEnd string, limit int) ([]string, error) {
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
	err := json.Unmarshal([]byte(config), &cf)
	if err != nil {
		return berror.Wrapf(err, cache.InvalidSsdbCacheCfg,
			"unmarshal this config failed, it must be a valid json string: %s", config)
	}
	if _, ok := cf["conn"]; !ok {
		return berror.Wrapf(err, cache.InvalidSsdbCacheCfg,
			"Missing conn field: %s", config)
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	return rc.connectInit()
}

// connect to memcache and keep the connection.
func (rc *Cache) connectInit() error {
	conninfoArray := strings.Split(rc.conninfo[0], ":")
	if len(conninfoArray) < 2 {
		return berror.Errorf(cache.InvalidSsdbCacheCfg, "The value of conn should be host:port: %s", rc.conninfo[0])
	}
	host := conninfoArray[0]
	port, e := strconv.Atoi(conninfoArray[1])
	if e != nil {
		return berror.Errorf(cache.InvalidSsdbCacheCfg, "Port is invalid. It must be integer, %s", rc.conninfo[0])
	}
	var err error
	if rc.conn, err = ssdb.Connect(host, port); err != nil {
		return berror.Wrapf(err, cache.InvalidConnection,
			"could not connect to SSDB, please check your connection info, network and firewall: %s", rc.conninfo[0])
	}
	return nil
}

func init() {
	cache.Register("ssdb", NewSsdbCache)
}
