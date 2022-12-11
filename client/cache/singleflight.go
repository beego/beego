package cache

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"golang.org/x/sync/singleflight"
	"time"
)

// SingleflightCache
// This is a very simple decorator mode
type SingleflightCache struct {
	Cache
	group      *singleflight.Group
	Expiration time.Duration
	LoadFunc   func(ctx context.Context, key string) (any, error)
}

// NewSingleflightCache create SingleflightCache
func NewSingleflightCache(c Cache, expiration time.Duration,
	loadFunc func(ctx context.Context, key string) (any, error)) (Cache, error) {
	if loadFunc == nil {
		return nil, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil")
	}
	return &SingleflightCache{
		Cache:      c,
		group:      &singleflight.Group{},
		Expiration: expiration,
		LoadFunc:   loadFunc,
	}, nil
}

// Get In the Get method, single flight is used to load data and write back the cache.
func (s *SingleflightCache) Get(ctx context.Context, key string) (any, error) {
	val, err := s.Cache.Get(ctx, key)
	fmt.Println(val)
	if val == nil || err != nil {
		val, err, _ = s.group.Do(key, func() (interface{}, error) {
			v, er := s.LoadFunc(ctx, key)
			fmt.Println(v)
			if er != nil {
				return nil, berror.Wrap(er, KeyNotExist, "cache unable to load data")
			}
			er = s.Cache.Put(ctx, key, v, s.Expiration)
			return v, er
		})
	}
	return val, err
}
