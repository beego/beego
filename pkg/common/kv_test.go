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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVs(t *testing.T) {
	key := "my-key"
	kvs := NewKVs(KV{
		Key:   key,
		Value: 12,
	})

	assert.True(t, kvs.Contains(key))

	kvs.IfContains(key, func(value interface{}) {
		kvs.Put("my-key1", "")
	})

	assert.True(t, kvs.Contains("my-key1"))

	v := kvs.GetValueOr(key, 13)
	assert.Equal(t, 12, v)
}
