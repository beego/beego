package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/stretchr/testify/assert"
)

type Mock struct {
	Db      map[string]any
	loadCnt int64
}

func TestBloomFilterCache(t *testing.T) {
	var (
		mockDb = &Mock{Db: make(map[string]any)}
		blm    = bloom.NewWithEstimates(1000, 0.9)
		l      = sync.RWMutex{}
	)
	mockDb.Db = map[string]any{
		"key1": "val1",
		"key2": "val2",
	}
	bfc, err := NewBloomFilterCache(
		NewMemoryCache(),
		func(ctx context.Context, key string) (any, error) {
			defer l.Unlock()
			l.Lock()
			mockDb.loadCnt += 1 // flag of number load data from db
			v, ok := mockDb.Db[key]
			if !ok {
				return nil, errors.New("fail")
			}
			return v, nil
		},
		blm,
	)
	assert.Nil(t, err)

	// case: keys not exist in cache, not exist in bloom, but exist in db,
	// want: not load data from db
	_, err = bfc.Get(context.Background(), "key1")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), mockDb.loadCnt)

	// case: set cache
	// want: add key to bloom
	err = bfc.Put(context.Background(), "key1", "val1", time.Minute)
	assert.Nil(t, err)
	assert.True(t, blm.TestString("key1"))

	// case: keys exist in cache
	// want: not load from db
	v1, err := bfc.Get(context.Background(), "key1")
	assert.Equal(t, mockDb.Db["key1"], v1)
	assert.Equal(t, int64(0), mockDb.loadCnt)

	// case: keys not exist in cache, but exist in bloom, exist in db,
	// want: load data from db, and set cache
	bfc.AddString("key2")
	_, err = bfc.Get(context.Background(), "key2")
	assert.Nil(t, err)
	v2, err := bfc.Get(context.Background(), "key2")
	assert.Nil(t, err)
	assert.Equal(t, mockDb.Db["key2"], v2)
	assert.Equal(t, int64(1), mockDb.loadCnt)

	// case: keys not exist in cache, not exist in bloom
	// want: not load data from db
	_, err = bfc.Get(context.Background(), "key3")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mockDb.loadCnt)

	// case: keys not exist in cache, not exist in bloom, execute Get single key 100 times
	// want: not load data from db
	for i := 0; i < 100; i++ {
		_, err = bfc.Get(context.Background(), "not_exist_key")
	}
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mockDb.loadCnt)

	// case: keys not exist in cache, not exist in bloom, execute Get single key 100 times concurrently
	// want: not load data from db
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			_, err = bfc.Get(context.Background(), "not_exist_key")
		}()
	}
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mockDb.loadCnt)

	// case: keys in bloom, execute Get different keys concurrently
	// want: load data from db, but not set cache
	bfc.BloomFilter.ClearAll()
	for i := 1; i <= 100; i++ {
		bfc.AddString(fmt.Sprintf("key%d", i))
	}
	wg1 := sync.WaitGroup{}
	wg1.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg1.Done()
			_, err = bfc.Get(context.Background(), fmt.Sprintf("key%d", i))
		}()
	}
	wg1.Wait()
	assert.NotNil(t, err)
	assert.Equal(t, int64(101), mockDb.loadCnt)

	// case: keys in bloom, execute Get different keys concurrently
	// want: load data from db, and set cache
	for i := 1; i <= 100; i++ {
		mockDb.Db[fmt.Sprintf("key%d", i)] = fmt.Sprintf("val%d", i)
	}

	wg2 := sync.WaitGroup{}
	wg2.Add(100)
	for i := 1; i <= 100; i++ {
		go func(i int, t *testing.T) {
			defer wg2.Done()
			kkk := fmt.Sprintf("key%d", i)
			vvv, er := bfc.Get(context.Background(), kkk)
			assert.Nil(t, er)
			assert.Equal(t, fmt.Sprintf("val%d", i), vvv)
		}(i, t)
	}
	wg2.Wait()
	assert.Equal(t, int64(199), mockDb.loadCnt)
}
