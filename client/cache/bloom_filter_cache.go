package cache

import (
	"context"
	"errors"
	"time"

	"github.com/beego/beego/v2/core/berror"
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilterCache struct {
	Cache
	*bloom.BloomFilter
	loadFunc       func(ctx context.Context, key string) (any, error)
	expiration     time.Duration // set cache expiration, default never expire
	notUpdateBloom bool          // update bloom key after put cache, default update
}

type BloomFilterCacheOption func(bfc *BloomFilterCache)

func WithExpirationOpt(t time.Duration) BloomFilterCacheOption {
	return func(bfc *BloomFilterCache) {
		bfc.expiration = t
	}
}

func WithUpdateBloomOpt(update bool) BloomFilterCacheOption {
	return func(bfc *BloomFilterCache) {
		bfc.notUpdateBloom = update
	}
}

func NewBloomFilterCache(cache Cache, ln func(context.Context, string) (any, error), blm *bloom.BloomFilter,
	opts ...BloomFilterCacheOption) (*BloomFilterCache, error) {

	if cache == nil || ln == nil || blm == nil {
		return nil, berror.Error(InvalidInitParameters, "missing required parameters")
	}

	bfc := &BloomFilterCache{
		Cache:       cache,
		BloomFilter: blm,
		loadFunc:    ln,
	}

	for _, opt := range opts {
		opt(bfc)
	}

	return bfc, nil
}

func (bfc *BloomFilterCache) Get(ctx context.Context, key string) (any, error) {
	val, err := bfc.Cache.Get(ctx, key)
	if err != nil && !errors.Is(err, ErrKeyNotExist) {
		return nil, err
	}
	if errors.Is(err, ErrKeyNotExist) {
		exist := bfc.BloomFilter.TestString(key)
		if exist {
			val, err = bfc.loadFunc(ctx, key)
			if err != nil {
				return nil, berror.Wrap(err, LoadFuncFailed, "cache unable to load data")
			}
			err = bfc.Put(ctx, key, val, bfc.expiration)
			if err != nil {
				return val, err
			}
		}
	}
	return val, nil
}

func (bfc *BloomFilterCache) Put(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := bfc.Cache.Put(ctx, key, val, expiration)
	if err != nil {
		return err
	}
	if !bfc.notUpdateBloom {
		bfc.AddString(key)
	}
	return nil
}
