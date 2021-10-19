package ratelimit

import (
	"sync"
	"time"
)

type tokenBucket struct {
	sync.RWMutex
	remaining   uint
	capacity    uint
	lastCheckAt time.Time
	rate        time.Duration
}

// newTokenBucket return an bucket that implements token bucket
func newTokenBucket(opts ...bucketOption) bucket {
	b := &tokenBucket{lastCheckAt: time.Now()}
	for _, o := range opts {
		o(b)
	}
	return b
}

func withCapacity(capacity uint) bucketOption {
	return func(b bucket) {
		bucket := b.(*tokenBucket)
		bucket.capacity = capacity
		bucket.remaining = capacity
	}
}

func withRate(rate time.Duration) bucketOption {
	return func(b bucket) {
		bucket := b.(*tokenBucket)
		bucket.rate = rate
	}
}

func (b *tokenBucket) getRemaining() uint {
	b.RLock()
	defer b.RUnlock()
	return b.remaining
}

func (b *tokenBucket) getRate() time.Duration {
	b.RLock()
	defer b.RUnlock()
	return b.rate
}

func (b *tokenBucket) getCapacity() uint {
	b.RLock()
	defer b.RUnlock()
	return b.capacity
}

func (b *tokenBucket) take(amount uint) bool {
	if b.rate <= 0 {
		return true
	}
	b.Lock()
	defer b.Unlock()
	now := time.Now()
	times := uint(now.Sub(b.lastCheckAt) / b.rate)
	b.lastCheckAt = b.lastCheckAt.Add(time.Duration(times) * b.rate)
	b.remaining += times
	if b.remaining < amount {
		return false
	}
	b.remaining -= amount
	if b.remaining > b.capacity {
		b.remaining = b.capacity
	}
	return true
}
