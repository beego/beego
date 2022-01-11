package ratelimit

import "time"

// bucket is an interface store ratelimit info
type bucket interface {
	take(amount uint) bool
	getCapacity() uint
	getRemaining() uint
	getRate() time.Duration
}

// bucketOption is constructor option
type bucketOption func(bucket)
