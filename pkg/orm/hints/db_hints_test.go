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
	"github.com/astaxie/beego/pkg/common"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHint_time(t *testing.T) {
	key1 := "key1"
	value1 := time.Second
	hint1 := NewHint(key1, value1)

	assert.Equal(t, hint1.GetKey(), key1)
	assert.Equal(t, hint1.GetValue(), value1)

	key2 := "key2"
	value2 := 281230
	hint2 := NewHint(key2, value2)

	assert.Equal(t, hint2.GetKey(), key2)
	assert.Equal(t, hint2.GetValue(), value2)

	key3 := "key3"
	value3 := 21.2459753
	hint3 := NewHint(key3, value3)

	assert.Equal(t, hint3.GetKey(), key3)
	assert.Equal(t, hint3.GetValue(), value3)
}

func TestNewHintFunc(t *testing.T) {
	kvs := common.NewKVs()

	key1 := "key1"
	value1 := time.Second
	key2 := "key2"
	value2 := 281230
	key3 := "key3"
	value3 := 21.2459753

	funcList := []HintFunc{
		NewHintFunc(key1, value1),
		NewHintFunc(key2, value2),
		NewHintFunc(key3, value3),
	}

	for _, tFunc := range funcList {
		tFunc(kvs)
	}

	assert.Equal(t, kvs.Contains(key1), true)
	assert.Equal(t, kvs.Contains(key2), true)
	assert.Equal(t, kvs.Contains(key3), true)

}

func TestMaxOpenConnections(t *testing.T) {
	i := 887423
	kvs := common.NewKVs()
	MaxOpenConnections(i)(kvs)

	value := kvs.GetValueOr(KeyMaxOpenConnections, nil)
	assert.Equal(t, value, i)
}

func TestConnMaxLifetime(t *testing.T) {
	i := time.Hour
	kvs := common.NewKVs()
	ConnMaxLifetime(i)(kvs)

	value := kvs.GetValueOr(KeyConnMaxLifetime, nil)
	assert.Equal(t, value, i)
}

func TestMaxIdleConnections(t *testing.T) {
	i := 42316
	kvs := common.NewKVs()
	MaxIdleConnections(i)(kvs)

	value := kvs.GetValueOr(KeyMaxIdleConnections, nil)
	assert.Equal(t, value, i)
}

func TestMaxStmtCacheSize(t *testing.T) {
	i := 94157
	kvs := common.NewKVs()
	MaxStmtCacheSize(i)(kvs)

	value := kvs.GetValueOr(KeyMaxStmtCacheSize, nil)
	assert.Equal(t, value, i)
}

func TestForceIndex(t *testing.T) {
	s := []string{`f_index1`, `f_index2`, `f_index3`}
	kvs := common.NewKVs()
	ForceIndex(s...)(kvs)

	value := kvs.GetValueOr(KeyForceIndex, nil)
	assert.Equal(t, value, s)
}

func TestForceIndex_0(t *testing.T) {
	var s []string
	kvs := common.NewKVs()
	ForceIndex(s...)(kvs)

	value := kvs.GetValueOr(KeyForceIndex, nil)
	assert.Equal(t, value, s)
}

func TestIgnoreIndex(t *testing.T) {
	s := []string{`i_index1`, `i_index2`, `i_index3`}
	kvs := common.NewKVs()
	IgnoreIndex(s...)(kvs)

	value := kvs.GetValueOr(KeyIgnoreIndex, nil)
	assert.Equal(t, value, s)
}

func TestIgnoreIndex_0(t *testing.T) {
	var s []string
	kvs := common.NewKVs()
	IgnoreIndex(s...)(kvs)

	value := kvs.GetValueOr(KeyIgnoreIndex, nil)
	assert.Equal(t, value, s)
}

func TestUseIndex(t *testing.T) {
	s := []string{`u_index1`, `u_index2`, `u_index3`}
	kvs := common.NewKVs()
	UseIndex(s...)(kvs)

	value := kvs.GetValueOr(KeyUseIndex, nil)
	assert.Equal(t, value, s)
}

func TestUseIndex_0(t *testing.T) {
	var s []string
	kvs := common.NewKVs()
	UseIndex(s...)(kvs)

	value := kvs.GetValueOr(KeyUseIndex, nil)
	assert.Equal(t, value, s)
}

func TestForUpdate(t *testing.T) {
	kvs := common.NewKVs()
	ForUpdate()(kvs)

	value := kvs.GetValueOr(KeyForUpdate, nil)
	assert.Equal(t, value, true)
}

func TestDefaultRelDepth(t *testing.T) {
	kvs := common.NewKVs()
	DefaultRelDepth()(kvs)

	value := kvs.GetValueOr(KeyRelDepth, nil)
	assert.Equal(t, value, true)
}

func TestRelDepth(t *testing.T) {
	i := 1579625
	kvs := common.NewKVs()
	RelDepth(i)(kvs)

	value := kvs.GetValueOr(KeyRelDepth, nil)
	assert.Equal(t, value, i)
}

func TestLimit(t *testing.T) {
	i := int64(1579625)
	kvs := common.NewKVs()
	Limit(i)(kvs)

	value := kvs.GetValueOr(KeyLimit, nil)
	assert.Equal(t, value, i)
}

func TestOffset(t *testing.T) {
	i := int64(1572123965)
	kvs := common.NewKVs()
	Offset(i)(kvs)

	value := kvs.GetValueOr(KeyOffset, nil)
	assert.Equal(t, value, i)
}

func TestOrderBy(t *testing.T) {
	kvs := common.NewKVs()
	OrderBy(`-ID`)(kvs)

	value := kvs.GetValueOr(KeyOrderBy, nil)
	assert.Equal(t, value, `-ID`)
}
