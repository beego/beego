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
	"github.com/beego/beego/v2/core/utils"
)

const (
	// query level
	KeyForceIndex = iota
	KeyUseIndex
	KeyIgnoreIndex
	KeyForUpdate
	KeyLimit
	KeyOffset
	KeyOrderBy
	KeyRelDepth
)

type Hint struct {
	key   interface{}
	value interface{}
}

var _ utils.KV = new(Hint)

// GetKey return key
func (s *Hint) GetKey() interface{} {
	return s.key
}

// GetValue return value
func (s *Hint) GetValue() interface{} {
	return s.value
}

var _ utils.KV = new(Hint)

// ForceIndex return a hint about ForceIndex
func ForceIndex(indexes ...string) *Hint {
	return NewHint(KeyForceIndex, indexes)
}

// UseIndex return a hint about UseIndex
func UseIndex(indexes ...string) *Hint {
	return NewHint(KeyUseIndex, indexes)
}

// IgnoreIndex return a hint about IgnoreIndex
func IgnoreIndex(indexes ...string) *Hint {
	return NewHint(KeyIgnoreIndex, indexes)
}

// ForUpdate return a hint about ForUpdate
func ForUpdate() *Hint {
	return NewHint(KeyForUpdate, true)
}

// DefaultRelDepth return a hint about DefaultRelDepth
func DefaultRelDepth() *Hint {
	return NewHint(KeyRelDepth, true)
}

// RelDepth return a hint about RelDepth
func RelDepth(d int) *Hint {
	return NewHint(KeyRelDepth, d)
}

// Limit return a hint about Limit
func Limit(d int64) *Hint {
	return NewHint(KeyLimit, d)
}

// Offset return a hint about Offset
func Offset(d int64) *Hint {
	return NewHint(KeyOffset, d)
}

// OrderBy return a hint about OrderBy
func OrderBy(s string) *Hint {
	return NewHint(KeyOrderBy, s)
}

// NewHint return a hint
func NewHint(key interface{}, value interface{}) *Hint {
	return &Hint{
		key:   key,
		value: value,
	}
}
