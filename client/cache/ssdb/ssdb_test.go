package ssdb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/berror"
)

func TestSsdbcacheCache(t *testing.T) {
	ssdbAddr := os.Getenv("SSDB_ADDR")
	if ssdbAddr == "" {
		ssdbAddr = "127.0.0.1:8888"
	}

	ssdb, err := cache.NewCache("ssdb", fmt.Sprintf(`{"conn": "%s"}`, ssdbAddr))
	assert.Nil(t, err)

	// test put and exist
	res, _ := ssdb.IsExist(context.Background(), "ssdb")
	assert.False(t, res)
	timeoutDuration := 3 * time.Second
	// timeoutDuration := -10*time.Second   if timeoutDuration is negtive,it means permanent

	assert.Nil(t, ssdb.Put(context.Background(), "ssdb", "ssdb", timeoutDuration))

	res, _ = ssdb.IsExist(context.Background(), "ssdb")
	assert.True(t, res)

	// Get test done
	assert.Nil(t, ssdb.Put(context.Background(), "ssdb", "ssdb", timeoutDuration))

	v, _ := ssdb.Get(context.Background(), "ssdb")
	assert.Equal(t, "ssdb", v)

	// inc/dec test done
	assert.Nil(t, ssdb.Put(context.Background(), "ssdb", "2", timeoutDuration))

	assert.Nil(t, ssdb.Incr(context.Background(), "ssdb"))

	val, _ := ssdb.Get(context.Background(), "ssdb")
	v, err = strconv.Atoi(val.(string))
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	assert.Nil(t, ssdb.Decr(context.Background(), "ssdb"))

	// test del
	assert.Nil(t, ssdb.Put(context.Background(), "ssdb", "3", timeoutDuration))

	val, _ = ssdb.Get(context.Background(), "ssdb")
	v, err = strconv.Atoi(val.(string))
	assert.Equal(t, 3, v)
	assert.Nil(t, err)

	assert.Nil(t, ssdb.Delete(context.Background(), "ssdb"))
	assert.Nil(t, ssdb.Put(context.Background(), "ssdb", "ssdb", -10*time.Second))
	// test string

	res, _ = ssdb.IsExist(context.Background(), "ssdb")
	assert.True(t, res)

	v, _ = ssdb.Get(context.Background(), "ssdb")
	assert.Equal(t, "ssdb", v.(string))

	// test GetMulti done
	assert.Nil(t, ssdb.Put(context.Background(), "ssdb1", "ssdb1", -10*time.Second))

	res, _ = ssdb.IsExist(context.Background(), "ssdb1")
	assert.True(t, res)
	vv, _ := ssdb.GetMulti(context.Background(), []string{"ssdb", "ssdb1"})
	assert.Equal(t, 2, len(vv))

	assert.Equal(t, "ssdb", vv[0])
	assert.Equal(t, "ssdb1", vv[1])

	vv, err = ssdb.GetMulti(context.Background(), []string{"ssdb", "ssdb11"})

	assert.Equal(t, 2, len(vv))

	assert.Equal(t, "ssdb", vv[0])
	assert.Nil(t, vv[1])

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "key not exist"))

	// test clear all done
	assert.Nil(t, ssdb.ClearAll(context.Background()))
	e1, _ := ssdb.IsExist(context.Background(), "ssdb")
	e2, _ := ssdb.IsExist(context.Background(), "ssdb1")
	assert.False(t, e1)
	assert.False(t, e2)
}

func TestReadThroughCache_ssdb_Get(t *testing.T) {
	bm, err := cache.NewCache("ssdb", fmt.Sprintf(`{"conn": "%s"}`, "127.0.0.1:8888"))
	assert.Nil(t, err)

	testReadThroughCacheGet(t, bm)
}

func testReadThroughCacheGet(t *testing.T, bm cache.Cache) {
	testCases := []struct {
		name    string
		key     string
		value   string
		cache   cache.Cache
		wantErr error
	}{
		{
			name: "Get load err",
			key:  "key0",
			cache: func() cache.Cache {
				kvs := map[string]any{"key0": "value0"}
				db := &MockOrm{kvs: kvs}
				loadfunc := func(ctx context.Context, key string) (any, error) {
					val, er := db.Load(key)
					if er != nil {
						return nil, er
					}
					return val, nil
				}
				c, err := cache.NewReadThroughCache(bm, 3*time.Second, loadfunc)
				assert.Nil(t, err)
				return c
			}(),
			wantErr: func() error {
				err := errors.New("the key not exist")
				return berror.Wrap(
					err, cache.LoadFuncFailed, "cache unable to load data")
			}(),
		},
		{
			name:  "Get cache exist",
			key:   "key1",
			value: "value1",
			cache: func() cache.Cache {
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
				c, err := cache.NewReadThroughCache(bm, 3*time.Second, loadfunc)
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
			cache: func() cache.Cache {
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
				c, err := cache.NewReadThroughCache(bm, 3*time.Second, loadfunc)
				assert.Nil(t, err)
				return c
			}(),
		},
	}
	_, err := cache.NewReadThroughCache(bm, 3*time.Second, nil)
	assert.Equal(t, berror.Error(cache.InvalidLoadFunc, "loadFunc cannot be nil"), err)
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
