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
	"reflect"
	"time"
)

const TxNameKey = "TxName"

type filterOrmDecorator struct {
	ormer
	TxBeginner
	TxCommitter

	root Filter

	insideTx    bool
	txStartTime time.Time
	txName      string
}

func NewFilterOrmDecorator(delegate Ormer, filterChains ...FilterChain) Ormer {
	res := &filterOrmDecorator{
		ormer:      delegate,
		TxBeginner: delegate,
		root: func(ctx context.Context, inv *Invocation) {
			inv.execute()
		},
	}

	for i := len(filterChains) - 1; i >= 0; i-- {
		node := filterChains[i]
		res.root = node(res.root)
	}
	return res
}

func NewFilterTxOrmDecorator(delegate TxOrmer, root Filter, txName string) TxOrmer {
	res := &filterOrmDecorator{
		ormer:       delegate,
		TxCommitter: delegate,
		root:        root,
		insideTx:    true,
		txStartTime: time.Now(),
		txName: txName,
	}
	return res
}

func (f *filterOrmDecorator) Read(md interface{}, cols ...string) error {
	return f.ReadWithCtx(context.Background(), md, cols...)
}

func (f *filterOrmDecorator) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) (err error) {
	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "ReadWithCtx",
		Args:        []interface{}{md, cols},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			err = f.ormer.ReadWithCtx(ctx, md, cols...)
		},
	}
	f.root(ctx, inv)
	return err
}

func (f *filterOrmDecorator) ReadForUpdate(md interface{}, cols ...string) error {
	return f.ReadForUpdateWithCtx(context.Background(), md, cols...)
}

func (f *filterOrmDecorator) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	var err error
	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "ReadForUpdateWithCtx",
		Args:        []interface{}{md, cols},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			err = f.ormer.ReadForUpdateWithCtx(ctx, md, cols...)
		},
	}
	f.root(ctx, inv)
	return err
}

func (f *filterOrmDecorator) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return f.ReadOrCreateWithCtx(context.Background(), md, col1, cols...)
}

func (f *filterOrmDecorator) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	var (
		ok  bool
		res int64
		err error
	)

	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "ReadOrCreateWithCtx",
		Args:        []interface{}{md, col1, cols},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			ok, res, err = f.ormer.ReadOrCreateWithCtx(ctx, md, col1, cols...)
		},
	}
	f.root(ctx, inv)
	return ok, res, err
}

func (f *filterOrmDecorator) LoadRelated(md interface{}, name string, args ...interface{}) (int64, error) {
	return f.LoadRelatedWithCtx(context.Background(), md, name, args...)
}

func (f *filterOrmDecorator) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...interface{}) (int64, error) {
	var (
		res int64
		err error
	)

	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "LoadRelatedWithCtx",
		Args:        []interface{}{md, name, args},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res, err = f.ormer.LoadRelatedWithCtx(ctx, md, name, args...)
		},
	}
	f.root(ctx, inv)
	return res, err
}

func (f *filterOrmDecorator) QueryM2M(md interface{}, name string) QueryM2Mer {
	return f.QueryM2MWithCtx(context.Background(), md, name)
}

func (f *filterOrmDecorator) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) QueryM2Mer {
	var (
		res QueryM2Mer
	)

	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "QueryM2MWithCtx",
		Args:        []interface{}{md, name},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res = f.ormer.QueryM2MWithCtx(ctx, md, name)
		},
	}
	f.root(ctx, inv)
	return res
}

func (f *filterOrmDecorator) QueryTable(ptrStructOrTableName interface{}) QuerySeter {
	return f.QueryTableWithCtx(context.Background(), ptrStructOrTableName)
}

func (f *filterOrmDecorator) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) QuerySeter {
	var (
		res QuerySeter
		name string
		md interface{}
		mi *modelInfo
	)

	if table, ok := ptrStructOrTableName.(string); ok {
		name = table
	} else {
		name = getFullName(indirectType(reflect.TypeOf(ptrStructOrTableName)))
		md = ptrStructOrTableName
	}

	if m, ok := modelCache.getByFullName(name); ok {
		mi = m
	}

	inv := &Invocation{
		Method:      "QueryTableWithCtx",
		Args:        []interface{}{ptrStructOrTableName},
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		Md: md,
		mi: mi,
		f: func() {
			res = f.ormer.QueryTableWithCtx(ctx, ptrStructOrTableName)
		},
	}
	f.root(ctx, inv)
	return res
}

func (f *filterOrmDecorator) DBStats() *sql.DBStats {
	var (
		res *sql.DBStats
	)
	inv := &Invocation{
		Method:      "DBStats",
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res = f.ormer.DBStats()
		},
	}
	f.root(context.Background(), inv)
	return res
}

func (f *filterOrmDecorator) Insert(md interface{}) (int64, error) {
	return f.InsertWithCtx(context.Background(), md)
}

func (f *filterOrmDecorator) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	var (
		res int64
		err error
	)
	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "InsertWithCtx",
		Args:        []interface{}{md},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res, err = f.ormer.InsertWithCtx(ctx, md)
		},
	}
	f.root(ctx, inv)
	return res, err
}

func (f *filterOrmDecorator) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	return f.InsertOrUpdateWithCtx(context.Background(), md, colConflitAndArgs...)
}

func (f *filterOrmDecorator) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	var (
		res int64
		err error
	)
	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "InsertOrUpdateWithCtx",
		Args:        []interface{}{md, colConflitAndArgs},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res, err = f.ormer.InsertOrUpdateWithCtx(ctx, md, colConflitAndArgs...)
		},
	}
	f.root(ctx, inv)
	return res, err
}

func (f *filterOrmDecorator) InsertMulti(bulk int, mds interface{}) (int64, error) {
	return f.InsertMultiWithCtx(context.Background(), bulk, mds)
}

// InsertMultiWithCtx uses the first element's model info
func (f *filterOrmDecorator) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	var (
		res int64
		err error
		md  interface{}
		mi  *modelInfo
	)

	sind := reflect.Indirect(reflect.ValueOf(mds))

	if (sind.Kind() == reflect.Array || sind.Kind() == reflect.Slice) && sind.Len() > 0 {
		ind := reflect.Indirect(sind.Index(0))
		md = ind.Interface()
		mi, _ = modelCache.getByMd(md)
	}

	inv := &Invocation{
		Method:      "InsertMultiWithCtx",
		Args:        []interface{}{bulk, mds},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res, err = f.ormer.InsertMultiWithCtx(ctx, bulk, mds)
		},
	}
	f.root(ctx, inv)
	return res, err
}

func (f *filterOrmDecorator) Update(md interface{}, cols ...string) (int64, error) {
	return f.UpdateWithCtx(context.Background(), md, cols...)
}

func (f *filterOrmDecorator) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	var (
		res int64
		err error
	)
	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "UpdateWithCtx",
		Args:        []interface{}{md, cols},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res, err = f.ormer.UpdateWithCtx(ctx, md, cols...)
		},
	}
	f.root(ctx, inv)
	return res, err
}

func (f *filterOrmDecorator) Delete(md interface{}, cols ...string) (int64, error) {
	return f.DeleteWithCtx(context.Background(), md, cols...)
}

func (f *filterOrmDecorator) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	var (
		res int64
		err error
	)
	mi, _ := modelCache.getByMd(md)
	inv := &Invocation{
		Method:      "DeleteWithCtx",
		Args:        []interface{}{md, cols},
		Md:          md,
		mi:          mi,
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res, err = f.ormer.DeleteWithCtx(ctx, md, cols...)
		},
	}
	f.root(ctx, inv)
	return res, err
}

func (f *filterOrmDecorator) Raw(query string, args ...interface{}) RawSeter {
	return f.RawWithCtx(context.Background(), query, args...)
}

func (f *filterOrmDecorator) RawWithCtx(ctx context.Context, query string, args ...interface{}) RawSeter {
	var (
		res RawSeter
	)
	inv := &Invocation{
		Method:      "RawWithCtx",
		Args:        []interface{}{query, args},
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res = f.ormer.RawWithCtx(ctx, query, args...)
		},
	}
	f.root(ctx, inv)
	return res
}

func (f *filterOrmDecorator) Driver() Driver {
	var (
		res Driver
	)
	inv := &Invocation{
		Method:      "Driver",
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res = f.ormer.Driver()
		},
	}
	f.root(context.Background(), inv)
	return res
}

func (f *filterOrmDecorator) Begin() (TxOrmer, error) {
	return f.BeginWithCtxAndOpts(context.Background(), nil)
}

func (f *filterOrmDecorator) BeginWithCtx(ctx context.Context) (TxOrmer, error) {
	return f.BeginWithCtxAndOpts(ctx, nil)
}

func (f *filterOrmDecorator) BeginWithOpts(opts *sql.TxOptions) (TxOrmer, error) {
	return f.BeginWithCtxAndOpts(context.Background(), opts)
}

func (f *filterOrmDecorator) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error) {
	var (
		res TxOrmer
		err error
	)
	inv := &Invocation{
		Method:      "BeginWithCtxAndOpts",
		Args:        []interface{}{opts},
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		f: func() {
			res, err = f.TxBeginner.BeginWithCtxAndOpts(ctx, opts)
			res = NewFilterTxOrmDecorator(res, f.root, getTxNameFromCtx(ctx))
		},
	}
	f.root(ctx, inv)
	return res, err
}

func (f *filterOrmDecorator) DoTx(task func(txOrm TxOrmer) error) error {
	return f.DoTxWithCtxAndOpts(context.Background(), nil, task)
}

func (f *filterOrmDecorator) DoTxWithCtx(ctx context.Context, task func(txOrm TxOrmer) error) error {
	return f.DoTxWithCtxAndOpts(ctx, nil, task)
}

func (f *filterOrmDecorator) DoTxWithOpts(opts *sql.TxOptions, task func(txOrm TxOrmer) error) error {
	return f.DoTxWithCtxAndOpts(context.Background(), opts, task)
}

func (f *filterOrmDecorator) DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(txOrm TxOrmer) error) error {
	var (
		err error
	)

	inv := &Invocation{
		Method:      "DoTxWithCtxAndOpts",
		Args:        []interface{}{opts, task},
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		TxName:      getTxNameFromCtx(ctx),
		f: func() {
			err = f.TxBeginner.DoTxWithCtxAndOpts(ctx, opts, task)
		},
	}
	f.root(ctx, inv)
	return err
}

func (f *filterOrmDecorator) Commit() error {
	var (
		err error
	)
	inv := &Invocation{
		Method:      "Commit",
		Args:        []interface{}{},
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		TxName:      f.txName,
		f: func() {
			err = f.TxCommitter.Commit()
		},
	}
	f.root(context.Background(), inv)
	return err
}

func (f *filterOrmDecorator) Rollback() error {
	var (
		err error
	)
	inv := &Invocation{
		Method:      "Rollback",
		Args:        []interface{}{},
		InsideTx:    f.insideTx,
		TxStartTime: f.txStartTime,
		TxName:      f.txName,
		f: func() {
			err = f.TxCommitter.Rollback()
		},
	}
	f.root(context.Background(), inv)
	return err
}

func getTxNameFromCtx(ctx context.Context) string {
	txName := ""
	if n, ok := ctx.Value(TxNameKey).(string); ok {
		txName = n
	}
	return txName
}

