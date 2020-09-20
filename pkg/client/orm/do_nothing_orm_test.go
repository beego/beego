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

package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoNothingOrm(t *testing.T) {
	o := &DoNothingOrm{}
	err := o.DoTxWithCtxAndOpts(nil, nil, nil)
	assert.Nil(t, err)

	err = o.DoTxWithCtx(nil, nil)
	assert.Nil(t, err)

	err = o.DoTx(nil)
	assert.Nil(t, err)

	err = o.DoTxWithOpts(nil, nil)
	assert.Nil(t, err)

	assert.Nil(t, o.Driver())

	assert.Nil(t, o.QueryM2MWithCtx(nil, nil, ""))
	assert.Nil(t, o.QueryM2M(nil, ""))
	assert.Nil(t, o.ReadWithCtx(nil, nil))
	assert.Nil(t, o.Read(nil))

	txOrm, err := o.BeginWithCtxAndOpts(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	txOrm, err = o.BeginWithCtx(nil)
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	txOrm, err = o.BeginWithOpts(nil)
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	txOrm, err = o.Begin()
	assert.Nil(t, err)
	assert.Nil(t, txOrm)

	assert.Nil(t, o.RawWithCtx(nil, ""))
	assert.Nil(t, o.Raw(""))

	i, err := o.InsertMulti(0, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.Insert(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertOrUpdateWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertOrUpdate(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.InsertMultiWithCtx(nil, 0, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.LoadRelatedWithCtx(nil, nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.LoadRelated(nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	assert.Nil(t, o.QueryTableWithCtx(nil, nil))
	assert.Nil(t, o.QueryTable(nil))

	assert.Nil(t, o.Read(nil))
	assert.Nil(t, o.ReadWithCtx(nil, nil))
	assert.Nil(t, o.ReadForUpdateWithCtx(nil, nil))
	assert.Nil(t, o.ReadForUpdate(nil))

	ok, i, err := o.ReadOrCreate(nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)
	assert.False(t, ok)

	ok, i, err = o.ReadOrCreateWithCtx(nil, nil, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)
	assert.False(t, ok)

	i, err = o.Delete(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.DeleteWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.Update(nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	i, err = o.UpdateWithCtx(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), i)

	assert.Nil(t, o.DBStats())

	to := &DoNothingTxOrm{}
	assert.Nil(t, to.Commit())
	assert.Nil(t, to.Rollback())
}
