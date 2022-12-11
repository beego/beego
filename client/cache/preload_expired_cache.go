package cache

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"log"
	"time"
)

// PreloadCache can use in always-in-cache scene. Via sentinelCache,
// when the main cache key is going to expire, execute the function to flash the key expired time
type PreloadCache struct {
	Cache
	sentinelCache *MemoryCache
	// sentinelCache ahead of main cache expired time, default 3s
	expiredAhead time.Duration
	// main cache key expired time
	expired time.Duration
}

type PreloadExpireCacheOption func(*PreloadCache)

// WithSentinelAheadExpiredTime returns a PreloadExpireCacheOption that
// configures the expired time before main cache key expired
func WithSentinelAheadExpiredTime(before time.Duration) PreloadExpireCacheOption {
	return func(cache *PreloadCache) {
		cache.expiredAhead = before
	}
}

func NewPreloadCache(c Cache, load func(ctx context.Context, key string) (any, error),
	expired time.Duration, loadFailRetry int, opts ...PreloadExpireCacheOption) (*PreloadCache, error) {
	memCache, err := NewCache("memory", `{"interval":1}`)
	if err != nil {
		return nil, err
	}
	sentinel, ok := memCache.(*MemoryCache)
	if !ok {
		return nil, err
	}

	// register load function that will load val from main cache
	sentinel.RegisterEvictedFunc(func(key string, val any) {
		ctx := context.Background()
		var (
			mainVal any
			err     error
		)
		mainVal, err = load(ctx, key)
		if err != nil {
			for i := 0; i < loadFailRetry; i++ {
				if mainVal, err = load(ctx, key); err == nil {
					break
				}
			}
		}
		if err != nil {
			log.Fatalln(err)
		}

		sentinel.Unlock()
		// set sentinel cache
		err = sentinel.Put(ctx, key, "", expired-time.Second*3)
		sentinel.Lock()
		if err != nil {
			log.Fatalln(err)
		}
		// set main cache
		err = c.Put(ctx, key, mainVal, expired)
		if err != nil {
			log.Fatalln(err)
		}
	})

	pc := &PreloadCache{
		Cache:         c,
		sentinelCache: sentinel,
		expiredAhead:  3 * time.Second,
		expired:       expired,
	}

	for _, opt := range opts {
		opt(pc)
	}

	if expired-pc.expiredAhead <= 0 {
		return nil, berror.Error(InvalidPreloadCacheCfg, "main cache expired time too short")
	}

	return pc, nil
}

func (p *PreloadCache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	// if expired time of key < 3s, no need set sentinel cache
	err := p.sentinelCache.Put(ctx, key, "", p.expired-time.Second*3)
	if err != nil {
		return berror.Error(MemoryCacheCurdFailed, fmt.Sprintf("sentinel put key fail, key: %s", key))
	}
	return p.Cache.Put(ctx, key, val, timeout)
}
