// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	var t1 = "test1"

	assert.Equal(t, "test1", GetString(t1))
	var t2 = []byte("test2")
	assert.Equal(t, "test2", GetString(t2))
	var t3 = 1
	assert.Equal(t, "1", GetString(t3))
	var t4 int64 = 1
	assert.Equal(t, "1", GetString(t4))
	var t5 = 1.1
	assert.Equal(t, "1.1", GetString(t5))
	assert.Equal(t, "", GetString(nil))
}

func TestGetInt(t *testing.T) {
	var t1 = 1
	assert.Equal(t, 1, GetInt(t1))
	var t2 int32 = 32
	assert.Equal(t, 32, GetInt(t2))

	var t3 int64 = 64
	assert.Equal(t, 64, GetInt(t3))
	var t4 = "128"

	assert.Equal(t, 128, GetInt(t4))
	assert.Equal(t, 0, GetInt(nil))
}

func TestGetInt64(t *testing.T) {
	var i int64 = 1
	var t1 = 1
	assert.Equal(t, i, GetInt64(t1))
	var t2 int32 = 1

	assert.Equal(t, i, GetInt64(t2))
	var t3 int64 = 1
	assert.Equal(t, i, GetInt64(t3))
	var t4 = "1"
	assert.Equal(t, i, GetInt64(t4))
	assert.Equal(t, int64(0), GetInt64(nil))
}

func TestGetFloat64(t *testing.T) {
	var f = 1.11
	var t1 float32 = 1.11
	assert.Equal(t, f, GetFloat64(t1))
	var t2 = 1.11
	assert.Equal(t, f, GetFloat64(t2))
	var t3 = "1.11"
	assert.Equal(t, f, GetFloat64(t3))

	var f2 float64 = 1
	var t4 = 1
	assert.Equal(t, f2, GetFloat64(t4))

	assert.Equal(t, float64(0), GetFloat64(nil))
}

func TestGetBool(t *testing.T) {
	var t1 = true
	assert.True(t, GetBool(t1))
	var t2 = "true"
	assert.True(t, GetBool(t2))

	assert.False(t, GetBool(nil))
}
