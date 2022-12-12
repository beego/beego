package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestReadThroughCache_Memory_Get(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	assert.Nil(t, err)
	testReadThroughCacheGet(t, bm)

	testReadThroughCacheGetMulti(t, bm)
}

func TestReadThroughCache_file_Get(t *testing.T) {
	fc := NewFileCache().(*FileCache)
	fc.CachePath = "////aaa"
	err := fc.Init()
	assert.NotNil(t, err)
	fc.CachePath = getTestCacheFilePath()
	err = fc.Init()
	assert.Nil(t, err)
	testReadThroughCacheGet(t, fc)

	testReadThroughCacheGetMulti(t, fc)
}

func testReadThroughCacheGet(t *testing.T, bm Cache) {
	testCases := []struct {
		name    string
		key     string
		value   string
		cache   Cache
		wantErr error
	}{
		{
			name: "Get load err",
			key:  "key0",
			cache: func() Cache {
				kvs := map[string]any{"key0": "value0"}
				db := &MockOrm{kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc, false)
				assert.Nil(t, err)
				return c
			}(),
			wantErr: func() error {
				err := errors.New("the key not exist")
				return berror.Wrap(
					err, LoadFuncFailed, "cache unable to load data")
			}(),
		},
		{
			name:  "Get cache exist",
			key:   "key1",
			value: "value1",
			cache: func() Cache {
				keysMap := map[string]int{"key1": 1}
				kvs := map[string]any{"key1": "value1"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc, false)
				assert.Nil(t, err)
				err = c.Put(context.Background(), "key1", "value1", 3*time.Second)
				assert.Nil(t, err)
				return c
			}(),
		},
		{
			name:  "Get loadFunc exist",
			key:   "key2",
			value: "value2",
			cache: func() Cache {
				keysMap := map[string]int{"key2": 1}
				kvs := map[string]any{"key2": "value2"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc, false)
				assert.Nil(t, err)
				return c
			}(),
		},
	}
	_, err := NewReadThroughCache(bm, 3*time.Second, nil, false)
	assert.Equal(t, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil"), err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.cache
			val, err := c.Get(context.Background(), tc.key)
			if err != nil {
				assert.EqualError(t, tc.wantErr, err.Error())
				return
			}
			assert.Equal(t, tc.value, val)
		})

	}
}

func testReadThroughCacheGetMulti(t *testing.T, bm Cache) {
	testCases := []struct {
		name    string
		keys    []string
		values  []any
		cache   Cache
		wantErr error
	}{
		{
			name: "GetMulti load err",
			keys: []string{"key0", "key01"},
			cache: func() Cache {
				db := &MockOrm{}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc, true)
				assert.Nil(t, err)
				return c
			}(),
			wantErr: func() error {
				keysErr := make([]string, 0)
				err1 := berror.Wrap(
					errors.New("the key not exist"),
					LoadFuncFailed, "cache unable to load data")
				err2 := berror.Wrap(
					errors.New("the key not exist"),
					LoadFuncFailed, "cache unable to load data")
				keys := []string{"key0", "key01"}
				keyErrMap := map[string]error{"key0": err1, "key01": err2}
				for _, ki := range keys {
					err := keyErrMap[ki]
					keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, err.Error()))
				}
				return berror.Error(MultiGetFailed, strings.Join(keysErr, "; "))

			}(),
		},
		{
			name:   "GetMulti cache exist",
			keys:   []string{"key1", "key2"},
			values: []any{"value1", "value2"},
			cache: func() Cache {
				keysMap := map[string]int{"key1": 1, "key2": 1}
				kvs := map[string]any{"key1": "value1", "key2": "value2"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc, true)
				assert.Nil(t, err)
				for key, value := range kvs {
					err = c.Put(context.Background(), key, value, 3*time.Second)
					assert.Nil(t, err)
				}
				return c
			}(),
		},
		{
			name:   "GetMulti loadFunc exist",
			keys:   []string{"key3", "key4"},
			values: []any{"value3", "value4"},
			cache: func() Cache {
				keysMap := map[string]int{"key3": 1, "key4": 1}
				kvs := map[string]any{"key3": "value3", "key4": "value4"}
				db := &MockOrm{keysMap: keysMap, kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := NewReadThroughCache(bm, 3*time.Second, loadfunc, true)
				assert.Nil(t, err)
				return c
			}(),
		},
	}
	_, err := NewReadThroughCache(bm, 3*time.Second, nil, true)
	assert.Equal(t, berror.Error(InvalidLoadFunc, "loadFunc cannot be nil"), err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.cache
			val, err := c.GetMulti(context.Background(), tc.keys)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.EqualValues(t, tc.values, val)
		})

	}
}

type MockOrm struct {
	keysMap map[string]int
	kvs     map[string]any
}

func (m *MockOrm) Load(key string) (any, error) {
	_, ok := m.keysMap[key]
	if !ok {
		return nil, errors.New("the key not exist")
	}
	return m.kvs[key], nil
}
