package cache

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPreloadCache_Put(t *testing.T) {
	mc, err := NewCache("memory", `{"interval":1}`)
	assert.Nil(t, err)

	// mock save to db
	mockDb := make(map[string]any)
	mockDb["hello"] = "world"

	cache, err := NewPreloadCache(mc, func(ctx context.Context, key string) (any, error) {
		val, ok := mockDb[key]
		if !ok {
			return nil, errors.New("not found in db")
		}
		return val, nil
	}, 4*time.Second, 3)
	assert.Nil(t, err)

	err = cache.Put(context.Background(), "hello", "world", 4*time.Second)
	val, err := cache.Get(context.Background(), "hello")
	assert.Nil(t, err)
	assert.Equal(t, "world", val)

	time.Sleep(5 * time.Second)
	val, err = cache.Get(context.Background(), "hello")
	assert.Nil(t, err)
	assert.Equal(t, "world", val)
}
