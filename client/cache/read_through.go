package cache

import (
	"context"
	"github.com/beego/beego/v2/core/berror"
	"time"
)

// readThroughCache is a decorator
// add the read through function to the original Cache function
type readThroughCache struct {
	Cache
	expiration       time.Duration
	loadFunc         func(ctx context.Context, key string) (any, error)
	applyForGetMulti bool
}

// NewReadThroughCache create readThroughCache
func NewReadThroughCache(cache Cache, expiration time.Duration,
	loadFunc func(ctx context.Context, key string) (any, error), applyForGetMulti bool) (Cache, error) {
	if loadFunc == nil {
		return nil, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil")
	}
	return &readThroughCache{
		Cache:            cache,
		expiration:       expiration,
		loadFunc:         loadFunc,
		applyForGetMulti: applyForGetMulti,
	}, nil
}

// Get will try to call the LoadFunc to load data if the Cache returns value nil or non-nil error.
func (c *readThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := c.Cache.Get(ctx, key)
	if val == nil || err != nil {
		val, err = c.loadFunc(ctx, key)
		if err != nil {
			return nil, berror.Wrap(
				err, LoadFuncFailed, "cache unable to load data")
		}
		err = c.Cache.Put(ctx, key, val, c.expiration)
		if err != nil {
			return val, err
		}
	}
	return val, nil
}
