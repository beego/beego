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

package hints

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
	assert.Equal(t, hint.GetKey(), KeyMaxOpenConnections)
}

func TestConnMaxLifetime(t *testing.T) {
	i := time.Hour
	hint := ConnMaxLifetime(i)
	assert.Equal(t, hint.GetValue(), i)
	assert.Equal(t, hint.GetKey(), KeyConnMaxLifetime)
}

func TestMaxIdleConnections(t *testing.T) {
	i := 42316
	hint := MaxIdleConnections(i)
	assert.Equal(t, hint.GetValue(), i)
	assert.Equal(t, hint.GetKey(), KeyMaxIdleConnections)
}

func TestMaxStmtCacheSize(t *testing.T) {
	i := 94157
	hint := MaxStmtCacheSize(i)
	assert.Equal(t, hint.GetValue(), i)
	assert.Equal(t, hint.GetKey(), KeyMaxStmtCacheSize)
}

func TestForceIndex(t *testing.T) {
	s := []string{`f_index1`, `f_index2`, `f_index3`}
	hint := ForceIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyForceIndex)
}

func TestForceIndex_0(t *testing.T) {
	var s []string
	hint := ForceIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyForceIndex)
}

func TestIgnoreIndex(t *testing.T) {
	s := []string{`i_index1`, `i_index2`, `i_index3`}
	hint := IgnoreIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyIgnoreIndex)
}

func TestIgnoreIndex_0(t *testing.T) {
	var s []string
	hint := IgnoreIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyIgnoreIndex)
}

func TestUseIndex(t *testing.T) {
	s := []string{`u_index1`, `u_index2`, `u_index3`}
	hint := UseIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyUseIndex)
}

func TestUseIndex_0(t *testing.T) {
	var s []string
	hint := UseIndex(s...)
	assert.Equal(t, hint.GetValue(), s)
	assert.Equal(t, hint.GetKey(), KeyUseIndex)
}

func TestForUpdate(t *testing.T) {
	hint := ForUpdate()
	assert.Equal(t, hint.GetValue(), true)
	assert.Equal(t, hint.GetKey(), KeyForUpdate)
}

func TestDefaultRelDepth(t *testing.T) {
	hint := DefaultRelDepth()
	assert.Equal(t, hint.GetValue(), true)
	assert.Equal(t, hint.GetKey(), KeyRelDepth)
}

func TestRelDepth(t *testing.T) {
	hint := RelDepth(157965)
	assert.Equal(t, hint.GetValue(), 157965)
	assert.Equal(t, hint.GetKey(), KeyRelDepth)
}

func TestLimit(t *testing.T) {
	hint := Limit(1579625)
	assert.Equal(t, hint.GetValue(), int64(1579625))
	assert.Equal(t, hint.GetKey(), KeyLimit)
}

func TestOffset(t *testing.T) {
	hint := Offset(int64(1572123965))
	assert.Equal(t, hint.GetValue(), int64(1572123965))
	assert.Equal(t, hint.GetKey(), KeyOffset)
}

func TestOrderBy(t *testing.T) {
	hint := OrderBy(`-ID`)
	assert.Equal(t, hint.GetValue(), `-ID`)
	assert.Equal(t, hint.GetKey(), KeyOrderBy)
}