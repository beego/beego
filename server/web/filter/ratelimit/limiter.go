// Copyright 2020 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// Limiter is an interface used to ratelimit
type Limiter interface {
	take(amount uint, r *http.Request) bool
}

// limiterOption is constructor option
type limiterOption func(l *limiter)

type limiter struct {
	sync.RWMutex
	capacity      uint
	rate          time.Duration
	buckets       map[string]bucket
	bucketFactory func(opts ...bucketOption) bucket
	sessionKey    func(r *http.Request) string
	resp          RejectionResponse
}

// RejectionResponse stores response information
// for the request rejected by limiter
type RejectionResponse struct {
	code int
	body string
}

const perRequestConsumedAmount = 1

var defaultRejectionResponse = RejectionResponse{
	code: 429,
	body: "too many requests",
}

// NewLimiter return FilterFunc, the limiter enables rate limit
// according to the configuration.
func NewLimiter(opts ...limiterOption) web.FilterFunc {
	l := &limiter{
		buckets: make(map[string]bucket),
		sessionKey: func(r *http.Request) string {
			return defaultSessionKey(r)
		},
		bucketFactory: NewTokenBucket,
		resp:          defaultRejectionResponse,
	}
	for _, o := range opts {
		o(l)
	}

	return func(ctx *context.Context) {
		if !l.take(perRequestConsumedAmount, ctx.Request) {
			ctx.ResponseWriter.WriteHeader(l.resp.code)
			ctx.WriteString(l.resp.body)
		}
	}
}

// WithSessionKey return limiterOption. WithSessionKey config func
// which defines the request characteristic againstthe limit is applied
func WithSessionKey(f func(r *http.Request) string) limiterOption {
	return func(l *limiter) {
		l.sessionKey = f
	}
}

// WithRate return limiterOption. WithRate config how long it takes to
// generate a token.
func WithRate(r time.Duration) limiterOption {
	return func(l *limiter) {
		l.rate = r
	}
}

// WithCapacity return limiterOption. WithCapacity config the capacity size.
// The bucket with a capacity of n has n tokens after initialization. The capacity
// defines how many requests a client can make in excess of the rate.
func WithCapacity(c uint) limiterOption {
	return func(l *limiter) {
		l.capacity = c
	}
}

// WithBucketFactory return limiterOption. WithBucketFactory customize the
// implementation of Bucket.
func WithBucketFactory(f func(opts ...bucketOption) bucket) limiterOption {
	return func(l *limiter) {
		l.bucketFactory = f
	}
}

// WithRejectionResponse return limiterOption. WithRejectionResponse
// customize the response for the request rejected by the limiter.
func WithRejectionResponse(resp RejectionResponse) limiterOption {
	return func(l *limiter) {
		l.resp = resp
	}
}

func (l *limiter) take(amount uint, r *http.Request) bool {
	bucket := l.getBucket(r)
	if bucket == nil {
		return true
	}
	return bucket.take(amount)
}

func (l *limiter) getBucket(r *http.Request) bucket {
	key := l.sessionKey(r)
	l.RLock()
	b, ok := l.buckets[key]
	l.RUnlock()
	if !ok {
		b = l.createBucket(key)
	}

	return b
}

func (l *limiter) createBucket(key string) bucket {
	l.Lock()
	defer l.Unlock()
	// double check avoid overwriting
	b, ok := l.buckets[key]
	if ok {
		return b
	}
	b = l.bucketFactory(withCapacity(l.capacity), withRate(l.rate))
	l.buckets[key] = b
	return b
}


func defaultSessionKey(r *http.Request) string {
	return ""
}

func RemoteIPSessionKey(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
