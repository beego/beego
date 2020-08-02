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
	"time"
)

const (
	//db level
	KeyMaxIdleConnections = iota
	KeyMaxOpenConnections
	KeyConnMaxLifetime
	KeyMaxStmtCacheSize

	//query level
	KeyForceIndex
	KeyUseIndex
	KeyIgnoreIndex
	KeyForUpdate
)

type Hint struct {
	key   interface{}
	value interface{}
}

var _ common.KV = new(Hint)

// GetKey return key
func (s *Hint) GetKey() interface{} {
	return s.key
}

// GetValue return value
func (s *Hint) GetValue() interface{} {
	return s.value
}

var _ common.KV = new(Hint)

// MaxIdleConnections return a hint about MaxIdleConnections
func MaxIdleConnections(v int) *Hint {
	return NewHint(KeyMaxIdleConnections, v)
}

// MaxOpenConnections return a hint about MaxOpenConnections
func MaxOpenConnections(v int) *Hint {
	return NewHint(KeyMaxOpenConnections, v)
}

// ConnMaxLifetime return a hint about ConnMaxLifetime
func ConnMaxLifetime(v time.Duration) *Hint {
	return NewHint(KeyConnMaxLifetime, v)
}

// MaxStmtCacheSize return a hint about MaxStmtCacheSize
func MaxStmtCacheSize(v int) *Hint {
	return NewHint(KeyMaxStmtCacheSize, v)
}

// ForceIndex return a hint about ForceIndex
func ForceIndex(index ...string) *Hint {
	return NewHint(KeyForceIndex, index)
}

// UseIndex return a hint about UseIndex
func UseIndex(index ...string) *Hint {
	return NewHint(KeyUseIndex, index)
}

// IgnoreIndex return a hint about IgnoreIndex
func IgnoreIndex(index ...string) *Hint {
	return NewHint(KeyIgnoreIndex, index)
}

// ForUpdate return a hint about ForUpdate
func ForUpdate() *Hint {
	return NewHint(KeyForUpdate, true)
}

// NewHint return a hint
func NewHint(key interface{}, value interface{}) *Hint {
	return &Hint{
		key:   key,
		value: value,
	}
}
