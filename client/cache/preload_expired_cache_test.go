package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPreloadCache_Put(t *testing.T) {
	type fields struct {
		Cache         Cache
		sentinelCache *MemoryCache
		expiredAhead  time.Duration
		expired       time.Duration
	}
	type args struct {
		ctx     context.Context
		key     string
		val     interface{}
		timeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PreloadCache{
				Cache:         tt.fields.Cache,
				sentinelCache: tt.fields.sentinelCache,
				expiredAhead:  tt.fields.expiredAhead,
				expired:       tt.fields.expired,
			}
			tt.wantErr(t, p.Put(tt.args.ctx, tt.args.key, tt.args.val, tt.args.timeout), fmt.Sprintf("Put(%v, %v, %v, %v)", tt.args.ctx, tt.args.key, tt.args.val, tt.args.timeout))
		})
	}
}
