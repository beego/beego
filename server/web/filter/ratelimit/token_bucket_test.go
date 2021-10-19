package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	b := newTokenBucket(withRate(1 * time.Second)).(*tokenBucket)
	assert.Equal(t, b.getRate(), 1*time.Second)
}

func TestGetRemainingAndCapacity(t *testing.T) {
	b := newTokenBucket(withCapacity(10))
	assert.Equal(t, b.getRemaining(), uint(10))
	assert.Equal(t, b.getCapacity(), uint(10))
}

func TestTake(t *testing.T) {
	b := newTokenBucket(withCapacity(10), withRate(10*time.Millisecond)).(*tokenBucket)
	for i := 0; i < 10; i++ {
		assert.True(t, b.take(1))
	}
	assert.False(t, b.take(1))
	assert.Equal(t, b.getRemaining(), uint(0))
	b = newTokenBucket(withCapacity(1), withRate(1*time.Millisecond)).(*tokenBucket)
	assert.True(t, b.take(1))
	time.Sleep(2 * time.Millisecond)
	assert.True(t, b.take(1))
}
