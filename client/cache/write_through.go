package cache

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"time"
)

type WriteThoughCache struct {
	Cache
	StoreFunc func(ctx context.Context, key string, val any) error
}

func (w *WriteThoughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	if w.StoreFunc == nil {
		return berror.Error(InvalidStoreFunc, "storeFunc can not be nil")
	}
	err := w.StoreFunc(ctx, key, val)
	if err != nil {
		return berror.Wrap(err, PersistCacheFailed, fmt.Sprintf("key: %s, val: %v", key, val))
	}
	return w.Cache.Put(ctx, key, val, expiration)
}
