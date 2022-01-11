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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoNothingRawSetter(t *testing.T) {
	rs := &DoNothingRawSetter{}
	i, err := rs.ValuesList(nil)
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = rs.Values(nil)
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = rs.ValuesFlat(nil)
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = rs.RowsToStruct(nil, "", "")
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = rs.RowsToMap(nil, "", "")
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = rs.QueryRows()
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	err = rs.QueryRow()
	// assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	s, err := rs.Exec()
	assert.Nil(t, err)
	assert.Nil(t, s)

	p, err := rs.Prepare()
	assert.Nil(t, err)
	assert.Nil(t, p)

	rrs := rs.SetArgs()
	assert.Equal(t, rrs, rs)
}
