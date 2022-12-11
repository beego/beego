package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemoryCache_WriteBack(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":1}`)
	assert.Nil(t, err)
	timeoutDuration := 2 * time.Second

	// mock save to db
	mockDb := make(map[string]any)

	wc := NewWriteBackCache(bm, func(key string, val any) {
		mockDb[key] = val
	})
	assert.NotNil(t, wc)

	if err = wc.Put(context.Background(), "astaxie", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := wc.IsExist(context.Background(), "astaxie"); !res {
		t.Error("check err")
	}

	if v, _ := bm.Get(context.Background(), "astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(3 * time.Second)

	if res, _ := bm.IsExist(context.Background(), "astaxie"); res {
		t.Error("check err")
	}

	assert.Equal(t, mockDb["astaxie"], 1)

	if err = wc.Put(context.Background(), "hello1", "world1", 10*time.Second); err != nil {
		t.Error("set Error", err)
	}

	if err = wc.Put(context.Background(), "hello2", "world2", 10*time.Second); err != nil {
		t.Error("set Error", err)
	}

	err = wc.Close()
	assert.Nil(t, err)

	assert.Equal(t, mockDb, map[string]any{"astaxie": 1, "hello1": "world1", "hello2": "world2"})
}
