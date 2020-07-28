// Copyright 2020 beego-dev
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

package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewHint_time(t *testing.T) {
	key := "qweqwe"
	value := time.Second
	hint := NewHint(key, value)

	assert.Equal(t, hint.GetKey(), key)
	assert.Equal(t, hint.GetValue(), value)
}

func TestNewHint_int(t *testing.T) {
	key := "qweqwe"
	value := 281230
	hint := NewHint(key, value)

	assert.Equal(t, hint.GetKey(), key)
	assert.Equal(t, hint.GetValue(), value)
}

func TestNewHint_float(t *testing.T) {
	key := "qweqwe"
	value := 21.2459753
	hint := NewHint(key, value)

	assert.Equal(t, hint.GetKey(), key)
	assert.Equal(t, hint.GetValue(), value)
}

func TestMaxOpenConnections(t *testing.T) {
	i := 887423
	hint := MaxOpenConnections(i)
	assert.Equal(t, hint.GetValue(), i)
	assert.Equal(t, hint.GetKey(), maxOpenConnectionsKey)
}

func TestConnMaxLifetime(t *testing.T) {
	i := time.Hour
	hint := ConnMaxLifetime(i)
	assert.Equal(t, hint.GetValue(), i)
	assert.Equal(t, hint.GetKey(), connMaxLifetimeKey)
}

func TestMaxIdleConnections(t *testing.T) {
	i := 42316
	hint := MaxIdleConnections(i)
	assert.Equal(t, hint.GetValue(), i)
	assert.Equal(t, hint.GetKey(), maxIdleConnectionsKey)
}
