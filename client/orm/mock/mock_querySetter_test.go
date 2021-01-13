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

func TestDoNothingQuerySetter(t *testing.T) {
	setter := &DoNothingQuerySetter{}
	setter.GroupBy().Filter("").Limit(10).
		Distinct().Exclude("a").FilterRaw("", "").
		ForceIndex().ForUpdate().IgnoreIndex().
		Offset(11).OrderBy().RelatedSel().SetCond(nil).UseIndex()

	assert.True(t, setter.Exist())
	err := setter.One(nil)
	assert.Nil(t, err)
	i, err := setter.Count()
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.Delete()
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.All(nil)
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.Update(nil)
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.RowsToMap(nil, "", "")
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.RowsToStruct(nil, "", "")
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.Values(nil)
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.ValuesFlat(nil, "")
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	i, err = setter.ValuesList(nil)
	assert.Equal(t, int64(0), i)
	assert.Nil(t, err)

	ins, err := setter.PrepareInsert()
	assert.Nil(t, err)
	assert.Nil(t, ins)

	assert.NotNil(t, setter.GetCond())
}
