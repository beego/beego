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

type rawSQlExecutorSelector interface {
	Select(ctx context.Context, sql string, args ...interface{}) orm.Ormer
}

func (s *SimpleReadWriteOrm) Read(ctx context.Context, md interface{}, cols ...string) error {
	hint := ctx.Value("db")
	if hint == "master"{
		return s.master.Read(ctx, md, cols...)
	}
	return s.follower.Read(ctx, md, cols...)
}

func (s *SimpleReadWriteOrm) ReadForUpdate(ctx context.Context, md interface{}, cols ...string) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ReadOrCreate(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Insert(ctx context.Context, md interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) InsertOrUpdate(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) InsertMulti(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Update(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Delete(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) LoadRelated(ctx context.Context, md interface{}, name string, args ...interface{}) (int64, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) QueryM2M(ctx context.Context, md interface{}, name string) orm.QueryM2Mer {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) QueryTable(ctx context.Context, ptrStructOrTableName interface{}) orm.QuerySeter {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) BeginTx(ctx context.Context) (*sql.Tx, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) BeginTxWithOpts(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ExecuteTx(ctx context.Context, task func() error) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) ExecuteTxWithOpts(ctx context.Context, opts *sql.TxOptions, task func() error) error {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Raw(ctx context.Context, query string, args ...interface{}) orm.RawSeter {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) Driver() orm.Driver {
	panic("implement me")
}

func (s *SimpleReadWriteOrm) DBStats() *sql.DBStats {
	panic("implement me")
}




