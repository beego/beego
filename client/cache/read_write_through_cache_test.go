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

package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadThroughCache_New(t *testing.T) {
	bm, err := cache.NewCache("memory", `{"interval":20}`)
	assert.NoError(t, err)

	c := cache.NewReadThroughCache(bm, func(ctx context.Context, key string) (any, error) {
		return nil, nil
	}, time.Millisecond)
	assert.NotNil(t, c)
}

func TestReadThroughCache_Get(t *testing.T) {

	fakeLoadFuncError := errors.New("loadFunc: fake error")

	testCases := map[string]struct {
		newUnderlyingCacheFunc func() cache.Cache
		loadFunc               func(ctx context.Context, key string) (any, error)
		keyExpiration          time.Duration

		key string

		wantVal any
		wantErr error
	}{
		`key exists, get key from underlying cache`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				require.NoError(t, c.Put(context.Background(), "Key1", "Val1", time.Millisecond))
				return c
			},
			keyExpiration: time.Millisecond,
			key:           "Key1",
			wantVal:       "Val1",
		},

		`key exists, get key from loadFunc`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				db := map[string]string{"Key2": "Val2"}
				return db[key], nil
			},
			keyExpiration: time.Millisecond,
			key:           "Key2",
			wantVal:       "Val2",
		},

		`key doesn't exist`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				return nil, fakeLoadFuncError
			},
			keyExpiration: time.Millisecond,
			key:           "Key3",
			wantErr:       fakeLoadFuncError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			underlyingCache := tc.newUnderlyingCacheFunc()

			readThroughCache := cache.NewReadThroughCache(underlyingCache, tc.loadFunc, tc.keyExpiration)
			assert.NotNil(t, readThroughCache)

			val, err := readThroughCache.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)

			val, err = underlyingCache.Get(context.Background(), tc.key)
			require.NoError(t, err)
			assert.Equal(t, tc.wantVal, val)

			<-time.After(tc.keyExpiration)
			_, err = underlyingCache.Get(context.Background(), tc.key)
			assert.Error(t, err)
		})
	}
}

func TestReadThroughCache_GetMulti(t *testing.T) {

	fakeLoadFuncError := errors.New("loadFunc: fake error")

	testCases := map[string]struct {
		newUnderlyingCacheFunc func() cache.Cache
		loadFunc               func(ctx context.Context, key string) (any, error)
		keyExpiration          time.Duration

		ctx  context.Context
		keys []string

		wantValues []any
		wantErr    error
	}{
		`all keys are exist, all keys are from underlying cache`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				require.NoError(t, c.Put(context.Background(), "Key100", "Val100", time.Millisecond))
				return c
			},
			keyExpiration: time.Millisecond,
			keys:          []string{"Key100"},
			wantValues:    []any{"Val100"},
		},

		`all keys are exist, all keys are from loadFunc`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				db := map[string]string{
					"Key103": "Val103",
					"Key104": "Val104",
					"Key105": "Val105",
				}
				return db[key], nil
			},
			keyExpiration: time.Millisecond,
			keys:          []string{"Key103", "Key104", "Key105"},
			wantValues:    []any{"Val103", "Val104", "Val105"},
		},

		`all keys are exist, all keys are from underlying cache and loadFunc`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				require.NoError(t, c.Put(context.Background(), "Key111", "Val111", time.Millisecond))
				require.NoError(t, c.Put(context.Background(), "Key112", "Val112", time.Millisecond))
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				db := map[string]string{
					"Key113": "Val113",
					"Key114": "Val114",
					"Key115": "Val115",
				}
				return db[key], nil
			},
			keyExpiration: time.Millisecond,
			keys:          []string{"Key111", "Key112", "Key113", "Key114", "Key115"},
			wantValues:    []any{"Val111", "Val112", "Val113", "Val114", "Val115"},
		},

		`all keys aren't exist`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				return nil, fakeLoadFuncError
			},
			keyExpiration: time.Millisecond,
			keys:          []string{"Key126"},
			wantErr:       fakeLoadFuncError,
		},

		`some keys aren't exist, some keys are from underlying cache`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				require.NoError(t, c.Put(context.Background(), "Key221", "Val221", time.Millisecond))
				require.NoError(t, c.Put(context.Background(), "Key222", "Val222", time.Millisecond))
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				return nil, fakeLoadFuncError
			},
			keyExpiration: time.Millisecond,
			keys:          []string{"Key221", "Key222", "Key223", "Key224"},
			wantErr:       fakeLoadFuncError,
		},

		`some keys aren't exist, some keys are from loadFunc`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				db := map[string]string{
					"Key334": "Val334",
					"Key335": "Val335",
				}
				val, ok := db[key]
				if !ok {
					return nil, fakeLoadFuncError
				}
				return val, nil
			},
			keyExpiration: time.Millisecond,
			keys:          []string{"Key331", "Key332", "Key333", "Key334", "Key335"},
			wantErr:       fakeLoadFuncError,
		},
		`some keys aren't exist, some keys are from underlying cache and loadFunc`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				require.NoError(t, c.Put(context.Background(), "Key444", "Val444", time.Millisecond))
				return c
			},
			loadFunc: func(ctx context.Context, key string) (any, error) {
				db := map[string]string{
					"Key448": "Val448",
				}
				val, ok := db[key]
				if !ok {
					return nil, fakeLoadFuncError
				}
				return val, nil
			},
			keyExpiration: time.Millisecond,
			keys:          []string{"Key444", "Key445", "Key446", "Key447", "Key448"},
			wantErr:       fakeLoadFuncError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			underlyingCache := tc.newUnderlyingCacheFunc()

			readThroughCache := cache.NewReadThroughCache(underlyingCache, tc.loadFunc, tc.keyExpiration)
			assert.NotNil(t, readThroughCache)

			values, err := readThroughCache.GetMulti(tc.ctx, tc.keys)

			if err != nil {
				assert.Contains(t, err.Error(), tc.wantErr.Error())
				return
			}
			assert.Equal(t, tc.wantValues, values)

			values, err = underlyingCache.GetMulti(tc.ctx, tc.keys)
			require.NoError(t, err)
			assert.Equal(t, tc.wantValues, values)

			<-time.After(tc.keyExpiration)
			_, err = underlyingCache.GetMulti(tc.ctx, tc.keys)
			assert.Error(t, err)
		})
	}
}

func TestWriteThroughCache_New(t *testing.T) {
	bm, err := cache.NewCache("memory", `{"interval":20}`)
	assert.NoError(t, err)

	c := cache.NewWriteThroughCache(bm, func(ctx context.Context, key string, val any) error {
		return nil
	})

	assert.NotNil(t, c)
}

func TestWriteThroughCache_Put(t *testing.T) {

	fakeRemoteDB := cache.NewMemoryCache()
	fakeErrUnderlyingCacheFailedToPutKey := errors.New("cache: failed to put key")
	fakeErrStoreFuncFailedToPutKey := errors.New("remoteDB: failed to put key")

	testCases := map[string]struct {
		newUnderlyingCacheFunc func() cache.Cache
		storeFunc              func(ctx context.Context, key string, val any) error

		key        string
		val        any
		expiration time.Duration

		wantErr error
	}{
		`put key succeed`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				return c
			},
			storeFunc: func(ctx context.Context, key string, val any) error {
				return fakeRemoteDB.Put(ctx, key, val, time.Second)
			},
			key:        "Key-1",
			val:        "Val-1",
			expiration: time.Millisecond,
		},

		`put key failed, error from underlying cache`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				return &PutFailedCacheStub{
					err: fakeErrUnderlyingCacheFailedToPutKey,
				}
			},
			storeFunc: func(ctx context.Context, key string, val any) error {
				return nil
			},
			key:        "Key-2",
			val:        "Val-2",
			expiration: time.Millisecond,
			wantErr:    fakeErrUnderlyingCacheFailedToPutKey,
		},

		`put key failed, error from storeFunc`: {
			newUnderlyingCacheFunc: func() cache.Cache {
				c, err := cache.NewCache("memory", `{"interval":20}`)
				require.NoError(t, err)
				return c
			},
			storeFunc: func(ctx context.Context, key string, val any) error {
				return fakeErrStoreFuncFailedToPutKey
			},
			key:        "Key-3",
			val:        "Val-3",
			expiration: time.Millisecond,
			wantErr:    fakeErrStoreFuncFailedToPutKey,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			underlyingCache := tc.newUnderlyingCacheFunc()
			writeThroughCache := cache.NewWriteThroughCache(underlyingCache, tc.storeFunc)

			err := writeThroughCache.Put(context.Background(), tc.key, tc.val, tc.expiration)
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}

			val, err := writeThroughCache.Get(context.Background(), tc.key)
			require.NoError(t, err)
			require.Equal(t, tc.val, val)

			val, err = fakeRemoteDB.Get(context.Background(), tc.key)
			require.NoError(t, err)
			require.Equal(t, tc.val, val)

			<-time.After(tc.expiration)
			_, err = writeThroughCache.Get(context.Background(), tc.key)
			require.Error(t, err)
		})
	}
}

type PutFailedCacheStub struct {
	cache.Cache
	err error
}

func (p *PutFailedCacheStub) Put(ctx context.Context, key string, val any, expiration time.Duration) error {
	return p.err
}
