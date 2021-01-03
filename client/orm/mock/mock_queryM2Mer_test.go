// Copyright 2020 beego
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

package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/orm"
)

func TestDoNothingQueryM2Mer(t *testing.T) {
	m2m := &DoNothingQueryM2Mer{}

	i, err := m2m.Clear()
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = m2m.Count()
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = m2m.Add()
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = m2m.Remove()
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	assert.True(t, m2m.Exist(nil))
}

func TestNewQueryM2MerCondition(t *testing.T) {
	cond := NewQueryM2MerCondition("", "")
	res := cond.Match(context.Background(), &orm.Invocation{})
	assert.True(t, res)
	cond = NewQueryM2MerCondition("hello", "")
	assert.False(t, cond.Match(context.Background(), &orm.Invocation{}))

	cond = NewQueryM2MerCondition("", "A")
	assert.False(t, cond.Match(context.Background(), &orm.Invocation{
		Args: []interface{}{0, "B"},
	}))

	assert.True(t, cond.Match(context.Background(), &orm.Invocation{
		Args: []interface{}{0, "A"},
	}))
}