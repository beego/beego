package cache

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ValueReader reads data from data source, such as database
type ValueReader interface {
	Query(key string) (interface{}, error)
}

// ValueWriter writes data to data source, such as database
type ValueWriter interface {
	Update(key string, val interface{}) error
}

// CacheAside is a decorator of Cache, it implements "cache aside" pattern
type CacheAside struct {
	Cache
	reader   ValueReader
	writer   ValueWriter
	lifeSpan time.Duration
}

func NewCacheAside(cache Cache, reader ValueReader, writer ValueWriter, lifeSpan time.Duration) *CacheAside {
	return &CacheAside{Cache: cache, reader: reader, writer: writer, lifeSpan: lifeSpan}
}

func (c *CacheAside) Get(ctx context.Context, key string) (interface{}, error) {
	if val, _ := c.Cache.Get(ctx, key); val != nil {
		return val, nil
	}

	val, err := c.reader.Query(key)
	if err != nil {
		return nil, err
	}

	return val, c.Cache.Put(ctx, key, val, c.lifeSpan)
}

func (c *CacheAside) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	keysErr := make([]string, 0)
	vals := make([]interface{}, len(keys))

	for i := range keys {
		val, err := c.Get(ctx, keys[i])
		vals[i] = val
		if err != nil {
			keysErr = append(keysErr, err.Error())
		}
	}

	if len(keysErr) == 0 {
		return vals, nil
	}

	return vals, fmt.Errorf(strings.Join(keysErr, "; "))
}

func (c *CacheAside) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	if err := c.writer.Update(key, val); err != nil {
		return err
	}

	if ok, _ := c.IsExist(ctx, key); !ok {
		return nil
	}

	return c.Cache.Delete(ctx, key)
}

func (c *CacheAside) Incr(ctx context.Context, key string) error {
	originVal, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	updateVal, err := incr(originVal)
	if err != nil {
		return err
	}

	if err := c.writer.Update(key, updateVal); err != nil {
		return err
	}

	return c.Cache.Delete(ctx, key)
}

func (c *CacheAside) Decr(ctx context.Context, key string) error {
	originVal, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	updateVal, err := decr(originVal)
	if err != nil {
		return err
	}

	if err := c.writer.Update(key, updateVal); err != nil {
		return err
	}

	return c.Cache.Delete(ctx, key)
}