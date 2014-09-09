/**
implements the HashCache interface
*/

package redis

import (
	"github.com/garyburd/redigo/redis"
)

// get the field value from a hash indexed by the key
func (rc *RedisCache) HGet(key, field string) interface{} {
	if v, err := rc.do("HGET", key, field); err == nil {
		return v
	}
	return nil
}

// save value in to a hash record indexed by the key
func (rc *RedisCache) HPut(key, field string, val interface{}) error {
	if _, err := rc.do("HSET", key, field, val); err != nil {
		return err
	}

	// save a record of the key so we can delete it by calling ClearAll
	if _, err := rc.do("HSET", rc.key, key, true); err != nil {
		return err
	}

	return nil
}

// get the field value from a hash indexed by the key
func (rc *RedisCache) HGetAll(key string) ([]interface{}, error) {
	ret, err := rc.do("HGetAll", key)
	if err == nil {
		return ret.([]interface{}), nil
	} else {
		return nil, err
	}
}

// delete value from a hash record
func (rc *RedisCache) HDelete(key, field string) error {
	_, err := rc.do("HDEL", key, field)
	return err
}

// increment the value under a particular field from a hash recored by the specified amount
func (rc *RedisCache) HIncrBy(key, field string, amount uint64) (uint64, error) {
	return redis.Uint64(rc.do("HIncrBy", key, field, amount))
}

// decrement the value under a particular field from a hash recored by the specified amount
func (rc *RedisCache) HDecrBy(key, field string, amount uint64) (uint64, error) {
	return redis.Uint64(rc.do("HIncrBy", key, field, -int64(amount)))
}
