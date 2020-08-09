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
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterOrmDecorator_Read(t *testing.T) {

	register()

	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "ReadWithCtx", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			next(ctx, inv)
		}
	})

	fte := &FilterTestEntity{}
	err := od.Read(fte)
	assert.NotNil(t, err)
	assert.Equal(t, "read error", err.Error())
}

func TestFilterOrmDecorator_BeginTx(t *testing.T) {
	register()

	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			if inv.Method == "BeginWithCtxAndOpts" {
				assert.Equal(t, 1, len(inv.Args))
				assert.Equal(t, "", inv.GetTableName())
				assert.False(t, inv.InsideTx)
			} else if inv.Method == "Commit" {
				assert.Equal(t, 0, len(inv.Args))
				assert.Equal(t, "Commit_tx", inv.TxName)
				assert.Equal(t, "", inv.GetTableName())
				assert.True(t, inv.InsideTx)
			} else if inv.Method == "Rollback" {
				assert.Equal(t, 0, len(inv.Args))
				assert.Equal(t, "Rollback_tx", inv.TxName)
				assert.Equal(t, "", inv.GetTableName())
				assert.True(t, inv.InsideTx)
			} else {
				t.Fail()
			}

			next(ctx, inv)
		}
	})
	to, err := od.Begin()
	assert.True(t, validateBeginResult(t, to, err))

	to, err = od.BeginWithOpts(nil)
	assert.True(t, validateBeginResult(t, to, err))

	ctx := context.WithValue(context.Background(), TxNameKey, "Commit_tx")
	to, err = od.BeginWithCtx(ctx)
	assert.True(t, validateBeginResult(t, to, err))

	err = to.Commit()
	assert.NotNil(t, err)
	assert.Equal(t, "commit", err.Error())

	ctx = context.WithValue(context.Background(), TxNameKey, "Rollback_tx")
	to, err = od.BeginWithCtxAndOpts(ctx, nil)
	assert.True(t, validateBeginResult(t, to, err))

	err = to.Rollback()
	assert.NotNil(t, err)
	assert.Equal(t, "rollback", err.Error())
}

func TestFilterOrmDecorator_DBStats(t *testing.T) {
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "DBStats", inv.Method)
			assert.Equal(t, 0, len(inv.Args))
			assert.Equal(t, "", inv.GetTableName())
			next(ctx, inv)
		}
	})
	res := od.DBStats()
	assert.NotNil(t, res)
	assert.Equal(t, -1, res.MaxOpenConnections)
}

func TestFilterOrmDecorator_Delete(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "DeleteWithCtx", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			next(ctx, inv)
		}
	})
	res, err := od.Delete(&FilterTestEntity{})
	assert.NotNil(t, err)
	assert.Equal(t, "delete error", err.Error())
	assert.Equal(t, int64(-2), res)
}

func TestFilterOrmDecorator_DoTx(t *testing.T) {
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "DoTxWithCtxAndOpts", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})

	err := od.DoTx(func(txOrm TxOrmer) error {
		return errors.New("tx error")
	})
	assert.NotNil(t, err)
	assert.Equal(t, "tx error", err.Error())

	err = od.DoTxWithCtx(context.Background(), func(txOrm TxOrmer) error {
		return errors.New("tx ctx error")
	})
	assert.NotNil(t, err)
	assert.Equal(t, "tx ctx error", err.Error())

	err = od.DoTxWithOpts(nil, func(txOrm TxOrmer) error {
		return errors.New("tx opts error")
	})
	assert.NotNil(t, err)
	assert.Equal(t, "tx opts error", err.Error())

	od = NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "DoTxWithCtxAndOpts", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "", inv.GetTableName())
			assert.Equal(t, "do tx name", inv.TxName)
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})

	ctx := context.WithValue(context.Background(), TxNameKey, "do tx name")
	err = od.DoTxWithCtxAndOpts(ctx, nil, func(txOrm TxOrmer) error {
		return errors.New("tx ctx opts error")
	})
	assert.NotNil(t, err)
	assert.Equal(t, "tx ctx opts error", err.Error())
}

func TestFilterOrmDecorator_Driver(t *testing.T) {
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "Driver", inv.Method)
			assert.Equal(t, 0, len(inv.Args))
			assert.Equal(t, "", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	res := od.Driver()
	assert.Nil(t, res)
}

func TestFilterOrmDecorator_Insert(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "InsertWithCtx", inv.Method)
			assert.Equal(t, 1, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})

	i, err := od.Insert(&FilterTestEntity{})
	assert.NotNil(t, err)
	assert.Equal(t, "insert error", err.Error())
	assert.Equal(t, int64(100), i)
}

func TestFilterOrmDecorator_InsertMulti(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "InsertMultiWithCtx", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})

	bulk := []*FilterTestEntity{&FilterTestEntity{}, &FilterTestEntity{}}
	i, err := od.InsertMulti(2, bulk)
	assert.NotNil(t, err)
	assert.Equal(t, "insert multi error", err.Error())
	assert.Equal(t, int64(2), i)
}

func TestFilterOrmDecorator_InsertOrUpdate(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "InsertOrUpdateWithCtx", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	i, err := od.InsertOrUpdate(&FilterTestEntity{})
	assert.NotNil(t, err)
	assert.Equal(t, "insert or update error", err.Error())
	assert.Equal(t, int64(1), i)
}

func TestFilterOrmDecorator_LoadRelated(t *testing.T) {
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "LoadRelatedWithCtx", inv.Method)
			assert.Equal(t, 3, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	i, err := od.LoadRelated(&FilterTestEntity{}, "hello")
	assert.NotNil(t, err)
	assert.Equal(t, "load related error", err.Error())
	assert.Equal(t, int64(99), i)
}

func TestFilterOrmDecorator_QueryM2M(t *testing.T) {
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "QueryM2MWithCtx", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	res := od.QueryM2M(&FilterTestEntity{}, "hello")
	assert.Nil(t, res)
}

func TestFilterOrmDecorator_QueryTable(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "QueryTableWithCtx", inv.Method)
			assert.Equal(t, 1, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	res := od.QueryTable(&FilterTestEntity{})
	assert.Nil(t, res)
}

func TestFilterOrmDecorator_Raw(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "RawWithCtx", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	res := od.Raw("hh")
	assert.Nil(t, res)
}

func TestFilterOrmDecorator_ReadForUpdate(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "ReadForUpdateWithCtx", inv.Method)
			assert.Equal(t, 2, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	err := od.ReadForUpdate(&FilterTestEntity{})
	assert.NotNil(t, err)
	assert.Equal(t, "read for update error", err.Error())
}

func TestFilterOrmDecorator_ReadOrCreate(t *testing.T) {
	register()
	o := &filterMockOrm{}
	od := NewFilterOrmDecorator(o, func(next Filter) Filter {
		return func(ctx context.Context, inv *Invocation) {
			assert.Equal(t, "ReadOrCreateWithCtx", inv.Method)
			assert.Equal(t, 3, len(inv.Args))
			assert.Equal(t, "FILTER_TEST", inv.GetTableName())
			assert.False(t, inv.InsideTx)
			next(ctx, inv)
		}
	})
	ok, i, err := od.ReadOrCreate(&FilterTestEntity{}, "name")
	assert.NotNil(t, err)
	assert.Equal(t, "read or create error", err.Error())
	assert.True(t, ok)
	assert.Equal(t, int64(13), i)
}

// filterMockOrm is only used in this test file
type filterMockOrm struct {
	DoNothingOrm
}

func (f *filterMockOrm) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return true, 13, errors.New("read or create error")
}

func (f *filterMockOrm) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	return errors.New("read for update error")
}

func (f *filterMockOrm) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...interface{}) (int64, error) {
	return 99, errors.New("load related error")
}

func (f *filterMockOrm) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	return 1, errors.New("insert or update error")
}

func (f *filterMockOrm) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	return 2, errors.New("insert multi error")
}

func (f *filterMockOrm) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	return 100, errors.New("insert error")
}

func (f *filterMockOrm) DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(txOrm TxOrmer) error) error {
	return task(nil)
}

func (f *filterMockOrm) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	return -2, errors.New("delete error")
}

func (f *filterMockOrm) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error) {
	return &filterMockOrm{}, errors.New("begin tx")
}

func (f *filterMockOrm) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	return errors.New("read error")
}

func (f *filterMockOrm) Commit() error {
	return errors.New("commit")
}

func (f *filterMockOrm) Rollback() error {
	return errors.New("rollback")
}

func (f *filterMockOrm) DBStats() *sql.DBStats {
	return &sql.DBStats{
		MaxOpenConnections: -1,
	}
}

func validateBeginResult(t *testing.T, to TxOrmer, err error) bool {
	assert.NotNil(t, err)
	assert.Equal(t, "begin tx", err.Error())
	_, ok := to.(*filterOrmDecorator).TxCommitter.(*filterMockOrm)
	assert.True(t, ok)
	return true
}

var filterTestEntityRegisterOnce sync.Once

type FilterTestEntity struct {
	ID   int
	Name string
}

func register() {
	filterTestEntityRegisterOnce.Do(func() {
		RegisterModel(&FilterTestEntity{})
	})
}

func (f *FilterTestEntity) TableName() string {
	return "FILTER_TEST"
}
