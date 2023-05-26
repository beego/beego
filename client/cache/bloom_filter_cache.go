// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"errors"
	"time"

	"github.com/beego/beego/v2/core/berror"
)

type BloomFilterCache struct {
	Cache
	blm        BloomFilter
	loadFunc   func(ctx context.Context, key string) (any, error)
	expiration time.Duration // set cache expiration, default never expire
}

type BloomFilter interface {
	Test(data string) bool
	Add(data string)
}

func NewBloomFilterCache(cache Cache, ln func(ctx context.Context, key string) (any, error), blm BloomFilter,
	expiration time.Duration,
) (*BloomFilterCache, error) {
	if cache == nil || ln == nil || blm == nil {
		return nil, berror.Error(InvalidInitParameters, "missing required parameters")
	}

	return &BloomFilterCache{
		Cache:      cache,
		blm:        blm,
		loadFunc:   ln,
		expiration: expiration,
	}, nil
}

func (bfc *BloomFilterCache) Get(ctx context.Context, key string) (any, error) {
	val, err := bfc.Cache.Get(ctx, key)
	if err != nil && !errors.Is(err, ErrKeyNotExist) {
		return nil, err
	}
	if errors.Is(err, ErrKeyNotExist) {
		exist := bfc.blm.Test(key)
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
