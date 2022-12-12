package cache

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"time"
)

type WriteThoughCache struct {
	Cache
	storeFunc func(ctx context.Context, key string, val any) error
}

func NewWriteThoughCache(cache Cache, fn func(ctx context.Context, key string, val any) error) (*WriteThoughCache, error) {
	if fn == nil || cache == nil {
		return nil, berror.Error(InvalidInitParameters, "cache or storeFunc can not be nil")
	}

	w := &WriteThoughCache{
		Cache:     cache,
		storeFunc: fn,
	}
	return w, nil
}

func (w *WriteThoughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.storeFunc(ctx, key, val)
	if err != nil {
		return berror.Wrap(err, PersistCacheFailed, fmt.Sprintf("key: %s, val: %v", key, val))
	}
	return w.Cache.Put(ctx, key, val, expiration)
}
