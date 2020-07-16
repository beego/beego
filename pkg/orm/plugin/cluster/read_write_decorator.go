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

package cluster

import (
	"context"
	"database/sql"
	"github.com/astaxie/beego/pkg/orm"
)

type SimpleReadWriteOrm struct {
	master orm.Ormer
	follower orm.Ormer
}

type MultipleReadWriteOrm struct {
	master orm.Ormer
	followers map[string]orm.Ormer
}

func (m *MultipleReadWriteOrm) Read(md interface{}, cols ...string) error {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) ReadForUpdate(md interface{}, cols ...string) error {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) LoadRelated(md interface{}, name string, args ...interface{}) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...interface{}) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) QueryM2M(md interface{}, name string) orm.QueryM2Mer {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) orm.QueryM2Mer {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) QueryTable(ptrStructOrTableName interface{}) orm.QuerySeter {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) orm.QuerySeter {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) DBStats() *sql.DBStats {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) Insert(md interface{}) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) Update(md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) Delete(md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) Raw(query string, args ...interface{}) orm.RawSeter {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) RawWithCtx(ctx context.Context, query string, args ...interface{}) orm.RawSeter {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) Driver() orm.Driver {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) BeginTx(ctx context.Context) (*orm.TxOrmer, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) BeginTxWithOpts(ctx context.Context, opts *sql.TxOptions) (orm.TxOrmer, error) {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) ExecuteTx(ctx context.Context, task func(txOrm orm.TxOrmer) error) error {
	panic("implement me")
}

func (m *MultipleReadWriteOrm) ExecuteTxWithOpts(ctx context.Context, opts *sql.TxOptions, task func(txOrm orm.TxOrmer) error) error {
	panic("implement me")
}

type rawSQlExecutorSelector interface {
	Select(ctx context.Context, sql string, args ...interface{}) orm.Ormer
}
