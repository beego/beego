// Copyright 2020
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

package config

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseConfiger_DefaultBool(t *testing.T) {
	bc := newBaseConfier("true")
	assert.True(t, bc.DefaultBool(context.Background(), "key1", false))
	assert.True(t, bc.DefaultBool(context.Background(), "key2", true))
}

func TestBaseConfiger_DefaultFloat(t *testing.T) {
	bc := newBaseConfier("12.3")
	assert.Equal(t, 12.3, bc.DefaultFloat(context.Background(), "key1", 0.1))
	assert.Equal(t, 0.1, bc.DefaultFloat(context.Background(), "key2", 0.1))
}

func TestBaseConfiger_DefaultInt(t *testing.T) {
	bc := newBaseConfier("10")
	assert.Equal(t, 10, bc.DefaultInt(context.Background(), "key1", 8))
	assert.Equal(t, 8, bc.DefaultInt(context.Background(), "key2", 8))
}

func TestBaseConfiger_DefaultInt64(t *testing.T) {
	bc := newBaseConfier("64")
	assert.Equal(t, int64(64), bc.DefaultInt64(context.Background(), "key1", int64(8)))
	assert.Equal(t, int64(8), bc.DefaultInt64(context.Background(), "key2", int64(8)))
}

func TestBaseConfiger_DefaultString(t *testing.T) {
	bc := newBaseConfier("Hello")
	assert.Equal(t, "Hello", bc.DefaultString(context.Background(), "key1", "world"))
	assert.Equal(t, "world", bc.DefaultString(context.Background(), "key2", "world"))
}

func TestBaseConfiger_DefaultStrings(t *testing.T) {
	bc := newBaseConfier("Hello;world")
	assert.Equal(t, []string{"Hello", "world"}, bc.DefaultStrings(context.Background(), "key1", []string{"world"}))
	assert.Equal(t, []string{"world"}, bc.DefaultStrings(context.Background(), "key2", []string{"world"}))
}

func newBaseConfier(str1 string) *BaseConfiger {
	return &BaseConfiger{
		reader: func(ctx context.Context, key string) (string, error) {
			if key == "key1" {
				return str1, nil
			} else {
				return "", errors.New("mock error")
			}

		},
	}
}
