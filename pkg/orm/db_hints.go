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
	"github.com/astaxie/beego/pkg/common"
	"time"
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

const (
	maxIdleConnectionsKey = "MaxIdleConnections"
	maxOpenConnectionsKey = "MaxOpenConnections"
	connMaxLifetimeKey    = "ConnMaxLifetime"
)

var _ common.KV = new(Hint)

// MaxIdleConnections return a hint about MaxIdleConnections
func MaxIdleConnections(v int) *Hint {
	return NewHint(maxIdleConnectionsKey, v)
}

// MaxOpenConnections return a hint about MaxOpenConnections
func MaxOpenConnections(v int) *Hint {
	return NewHint(maxOpenConnectionsKey, v)
}

// ConnMaxLifetime return a hint about ConnMaxLifetime
func ConnMaxLifetime(v time.Duration) *Hint {
	return NewHint(connMaxLifetimeKey, v)
}

// NewHint return a hint
func NewHint(key interface{}, value interface{}) *Hint {
	return &Hint{
		key:   key,
		value: value,
	}
}
