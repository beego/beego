// Copyright 2023 beego. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qb

import (
	"context"
	"database/sql"
	"github.com/beego/beego/v2/client/orm"
)

var _ Session = (*Tx)(nil)

type Tx struct {
	tx orm.TxOrmer
	db *DB
}

func (t *Tx) getCore() core {
	return t.db.core
}

func (t *Tx) queryContext(ctx context.Context, md any, sql string, args ...any) error {
	return t.tx.ReadRaw(ctx, md, sql, args)
}

func (t *Tx) execContext(ctx context.Context, md any, sql string, args ...any) (sql.Result, error) {
	return t.tx.ExecRaw(ctx, md, sql, args)
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
