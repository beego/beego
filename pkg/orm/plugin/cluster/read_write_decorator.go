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
	master   orm.Ormer
	follower orm.Ormer
}

var _ orm.Ormer = new(SimpleReadWriteOrm)

func (s *SimpleReadWriteOrm) Read(md interface{}, cols ...string) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ReadForUpdate(md interface{}, cols ...string) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) LoadRelated(md interface{}, name string, args ...interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) QueryM2M(md interface{}, name string) orm.QueryM2Mer {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) orm.QueryM2Mer {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) QueryTable(ptrStructOrTableName interface{}) orm.QuerySeter {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) orm.QuerySeter {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) DBStats() *sql.DBStats {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Insert(md interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Update(md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Delete(md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Raw(query string, args ...interface{}) orm.RawSeter {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) RawWithCtx(ctx context.Context, query string, args ...interface{}) orm.RawSeter {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Driver() orm.Driver {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Begin() (orm.TxOrmer, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) BeginWithCtx(ctx context.Context) (orm.TxOrmer, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) BeginWithOpts(opts *sql.TxOptions) (orm.TxOrmer, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (orm.TxOrmer, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) DoTransaction(task func(txOrm orm.TxOrmer) error) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) DoTransactionWithCtx(ctx context.Context, task func(txOrm orm.TxOrmer) error) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) DoTransactionWithOpts(opts *sql.TxOptions, task func(txOrm orm.TxOrmer) error) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) DoTransactionWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(txOrm orm.TxOrmer) error) error {
	panic("implement me")
}

type MultipleReadWriteOrm struct {
	master    orm.Ormer
	followers map[string]orm.Ormer
}

type rawSQlExecutorSelector interface {
	Select(ctx context.Context, sql string, args ...interface{}) orm.Ormer
}


