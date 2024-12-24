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
	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/qb/errs"
)

var _ Session = (*DB)(nil)

type DBOption func(db *DB)

func DBWithDialect(d Dialect) DBOption {
	return func(db *DB) {
		db.dialect = d
	}
}

type DB struct {
	core
	db orm.Ormer
}

func (db *DB) getCore() core {
	return db.core
}

func (db *DB) queryContext(ctx context.Context, md any, sql string, args ...any) error {
	return db.db.ReadRaw(ctx, md, sql, args)
}

func (db *DB) execContext(ctx context.Context, md any, sql string, args ...any) (sql.Result, error) {
	return db.db.ExecRaw(ctx, md, sql, args)
}

//func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
//	db := orm.NewOrmUsingDB(driver)
//	return OpenDB(driver, db, opts...)
//}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	err := orm.RegisterDB("default", driver, dsn, orm.MaxIdleConnections(20))
	if err != nil {
		return nil, err
	}
	return OpenDB(driver, orm.NewOrm(), opts...)
}

func OpenDB(driver string, db orm.Ormer, opts ...DBOption) (*DB, error) {
	dl, err := Of(driver)
	if err != nil {
		return nil, err
	}

	res := &DB{
		core: core{
			dialect:  dl,
			registry: models.NewModelCacheHandler(),
		},
		db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginWithCtxAndOpts(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

func (db *DB) Transaction(ctx context.Context,
	fn func(ctx context.Context, tx *Tx) error,
	opts *sql.TxOptions) (err error) {
	var tx *Tx
	tx, err = db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		if panicked || err != nil {
			exc := tx.Rollback()
			if exc != nil {
				err = errs.NewErrFailToRollbackTx(err, exc, panicked)
			}
		} else {
			err = tx.Commit()
		}
	}()
	err = fn(ctx, tx)
	panicked = false
	return err
}
