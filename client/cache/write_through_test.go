package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
			name:    "nil init parameters",
			wantErr: berror.Error(InvalidInitParameters, "cache or storeFunc can not be nil"),
		},
		{
			name:  "set error",
			cache: NewMemoryCache(),
			storeFunc: func(ctx context.Context, key string, val any) error {
				return errors.New("failed")
			},
			wantErr: berror.Wrap(errors.New("failed"), PersistCacheFailed,
				fmt.Sprintf("key: %s, val: %v", "", nil)),
		},
		{
			name:  "memory set success",
			cache: NewMemoryCache(),
			storeFunc: func(ctx context.Context, key string, val any) error {
				mockDbStore[key] = val
				return nil
			},
			key:   "hello",
			value: "world",
		},
		{
			name: "file set success",
			cache: func() Cache {
				fc := NewFileCache().(*FileCache)
				fc.CachePath = getTestCacheFilePath()
				return fc
			}(),
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
