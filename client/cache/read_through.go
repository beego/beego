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
	Expiration time.Duration
	LoadFunc   func(ctx context.Context, key string) (any, error)
}

// NewReadThroughCache create readThroughCache
func NewReadThroughCache(cache Cache, expiration time.Duration,
	loadFunc func(ctx context.Context, key string) (any, error)) (Cache, error) {
	if loadFunc == nil {
		return nil, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil")
	}
	return &readThroughCache{
		Cache:      cache,
		Expiration: expiration,
		LoadFunc:   loadFunc,
	}, nil
}

// Get cache from readThroughCache
func (c *readThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := c.Cache.Get(ctx, key)
	if val == nil || err != nil {
		val, err = c.LoadFunc(ctx, key)
		if err != nil {
			return nil, berror.Wrap(
				err, KeyNotExist, "cache unable to load data")
		}
		err = c.Cache.Put(ctx, key, val, c.Expiration)
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// GetMulti cache from readThroughCache
func (c *readThroughCache) GetMulti(ctx context.Context, keys []string) ([]any, error) {
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
