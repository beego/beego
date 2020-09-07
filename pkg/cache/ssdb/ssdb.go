package ssdb

import (
	"encoding/json"
	"errors"
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
func (rc *Cache) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	value, err := rc.conn.Get(key)
	if err == nil {
		return value
	}
	return nil
}

// GetMulti gets one or keys values from memcache.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var values []interface{}
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			for i := 0; i < size; i++ {
				values = append(values, err)
			}
			return values
		}
	}
	res, err := rc.conn.Do("multi_get", keys)
	resSize := len(res)
	if err == nil {
		for i := 1; i < resSize; i += 2 {
			values = append(values, res[i+1])
		}
		return values
	}
	for i := 0; i < size; i++ {
		values = append(values, err)
	}
	return values
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
func (rc *Cache) Put(key string, value interface{}, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	v, ok := value.(string)
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
func (rc *Cache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Del(key)
	return err
}

// Incr increases a key's counter.
func (rc *Cache) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("incr", key, 1)
	return err
}

// Decr decrements a key's counter.
func (rc *Cache) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Do("incr", key, -1)
	return err
}

// IsExist checks if a key exists in memcache.
func (rc *Cache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	resp, err := rc.conn.Do("exists", key)
	if err != nil {
		return false
	}
	if len(resp) == 2 && resp[1] == "1" {
		return true
	}
	return false

}

// ClearAll clears all cached items in memcache.
func (rc *Cache) ClearAll() error {
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
	// cache.Register("ssdb", NewSsdbCache)
}
