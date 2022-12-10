package ssdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/berror"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/cache"
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

func TestSsdbCache_WriteThough_Set(t *testing.T) {
	bm, err := cache.NewCache("ssdb", `{"conn": "127.0.0.1:8888"}`)
	assert.Nil(t, err)

	var mockDbStore = make(map[string]any)
	testCases := []struct {
		name      string
		storeFunc func(ctx context.Context, key string, val any) error
		key       string
		value     any
		wantErr   error
	}{
		{
			name:    "storeFunc nil",
			wantErr: berror.Error(cache.InvalidStoreFunc, "storeFunc can not be nil"),
		},
		{
			name: "set error",
			storeFunc: func(ctx context.Context, key string, val any) error {
				return errors.New("failed")
			},
			wantErr: berror.Wrap(errors.New("failed"), cache.PersistCacheFailed,
				fmt.Sprintf("key: %s, val: %v", "", nil)),
		},
		{
			name: "memory set success",
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
			w := &cache.WriteThoughCache{
				Cache:     bm,
				StoreFunc: tt.storeFunc,
			}
			err := w.Set(context.Background(), tt.key, tt.value, 60*time.Second)
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
