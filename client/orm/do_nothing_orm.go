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

	"github.com/beego/beego/v2/core/utils"
)

// DoNothingOrm won't do anything, usually you use this to custom your mock Ormer implementation
// I think golang mocking interface is hard to use
// this may help you to integrate with Ormer

var _ Ormer = new(DoNothingOrm)

type DoNothingOrm struct{}

func (d *DoNothingOrm) Read(md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadForUpdate(md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingOrm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return false, 0, nil
}

func (d *DoNothingOrm) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return false, 0, nil
}

func (d *DoNothingOrm) LoadRelated(md interface{}, name string, args ...utils.KV) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...utils.KV) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) QueryM2M(md interface{}, name string) QueryM2Mer {
	return nil
}

// NOTE: this method is deprecated, context parameter will not take effect.
func (d *DoNothingOrm) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) QueryM2Mer {
	return nil
}

func (d *DoNothingOrm) QueryTable(ptrStructOrTableName interface{}) QuerySeter {
	return nil
}

// NOTE: this method is deprecated, context parameter will not take effect.
func (d *DoNothingOrm) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) QuerySeter {
	return nil
}

func (d *DoNothingOrm) DBStats() *sql.DBStats {
	return nil
}

func (d *DoNothingOrm) Insert(md interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) Update(md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) Delete(md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingOrm) Raw(query string, args ...interface{}) RawSeter {
	return nil
}

func (d *DoNothingOrm) RawWithCtx(ctx context.Context, query string, args ...interface{}) RawSeter {
	return nil
}

func (d *DoNothingOrm) Driver() Driver {
	return nil
}

func (d *DoNothingOrm) Begin() (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) BeginWithCtx(ctx context.Context) (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) BeginWithOpts(opts *sql.TxOptions) (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error) {
	return nil, nil
}

func (d *DoNothingOrm) DoTx(task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

func (d *DoNothingOrm) DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

func (d *DoNothingOrm) DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

func (d *DoNothingOrm) DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return nil
}

// DoNothingTxOrm is similar with DoNothingOrm, usually you use it to test
type DoNothingTxOrm struct {
	DoNothingOrm
}

func (d *DoNothingTxOrm) Commit() error {
	return nil
}

func (d *DoNothingTxOrm) Rollback() error {
	return nil
}
