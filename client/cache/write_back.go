package cache

// WriteBackCache support that when key is expired,
// it can execute function of saving DB which defined by yourself automatically.
// Because only local memory cache can control the action of the expired keys,
// this mode only support MemoryCache Now.
type WriteBackCache struct {
	*MemoryCache
}

func NewWriteBackCache(cache Cache, onEvict func(key string, val any)) *WriteBackCache {
	memoryCache, ok := cache.(*MemoryCache)
	if !ok {
		return nil
	}

	w := &WriteBackCache{
		MemoryCache: memoryCache,
	}

	w.RegisterEvictedFunc(onEvict)
	return w
}

func (w *WriteBackCache) Close() error {
	w.RLock()
	defer w.RUnlock()
	for key, itm := range w.items {
		if itm.isExpire() {
			continue
		}
		w.onEvicted(key, itm.val)
	}
	return nil
}
