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

// nolint
package cache

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/beego/beego/v2/core/berror"
	"github.com/stretchr/testify/assert"
)

func TestWriteThoughCache_Set(t *testing.T) {
	var mockDbStore = make(map[string]any)

	testCases := []struct {
		name      string
		cache     Cache
		storeFunc func(ctx context.Context, key string, val any) error
		key       string
		value     any
		wantErr   error
	}{
		{
			name:  "store key/value in db fail",
			cache: NewMemoryCache(),
			storeFunc: func(ctx context.Context, key string, val any) error {
				return errors.New("failed")
			},
			wantErr: berror.Wrap(errors.New("failed"), PersistCacheFailed,
				fmt.Sprintf("key: %s, val: %v", "", nil)),
		},
		{
			name:  "store key/value success",
			cache: NewMemoryCache(),
			storeFunc: func(ctx context.Context, key string, val any) error {
				mockDbStore[key] = val
				return nil
			},
			key:   "hello",
			value: "world",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWriteThoughCache(tt.cache, tt.storeFunc)
			if err != nil {
				assert.EqualError(t, tt.wantErr, err.Error())
				return
			}

			err = w.Set(context.Background(), tt.key, tt.value, 60*time.Second)
			if err != nil {
				assert.EqualError(t, tt.wantErr, err.Error())
				return
			}

			val, err := w.Get(context.Background(), tt.key)
			assert.Nil(t, err)
			assert.Equal(t, tt.value, val)

			vv, ok := mockDbStore[tt.key]
			assert.True(t, ok)
			assert.Equal(t, tt.value, vv)
		})
	}
}

func TestNewWriteThoughCache(t *testing.T) {
	underlyingCache := NewMemoryCache()
	storeFunc := func(ctx context.Context, key string, val any) error { return nil }

	type args struct {
		cache Cache
		fn    func(ctx context.Context, key string, val any) error
	}
	tests := []struct {
		name    string
		args    args
		wantRes *WriteThoughCache
		wantErr error
	}{
		{
			name: "nil cache parameters",
			args: args{
				cache: nil,
				fn:    storeFunc,
			},
			wantErr: berror.Error(InvalidInitParameters, "cache or storeFunc can not be nil"),
		},
		{
			name: "nil storeFunc parameters",
			args: args{
				cache: underlyingCache,
				fn:    nil,
			},
			wantErr: berror.Error(InvalidInitParameters, "cache or storeFunc can not be nil"),
		},
		{
			name: "init write-though cache success",
			args: args{
				cache: underlyingCache,
				fn:    storeFunc,
			},
			wantRes: &WriteThoughCache{
				Cache:     underlyingCache,
				storeFunc: storeFunc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWriteThoughCache(tt.args.cache, tt.args.fn)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
		})
	}
}
