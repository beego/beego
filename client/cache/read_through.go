package cache

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"strings"
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

// GetMulti will try to call the LoadFunc to load data if the GetMulti method of underlying cache returns non-nil error.
// You should check the concrete type of underlying cache to learn the cases that the GetMulti function will return non-nil error.
func (c *readThroughCache) GetMulti(ctx context.Context, keys []string) ([]any, error) {
	if !c.applyForGetMulti {
		return c.Cache.GetMulti(ctx, keys)
	}
	rc := make([]interface{}, len(keys))
	keysErr := make([]string, 0)

	for i, ki := range keys {
		val, err := c.Get(context.Background(), ki)
		if err != nil {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, err.Error()))
			continue
		}
		rc[i] = val
	}
	if len(keysErr) == 0 {
		return rc, nil
	}
	return rc, berror.Error(MultiGetFailed, strings.Join(keysErr, "; "))
}
