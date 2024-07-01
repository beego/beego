// Copyright 2014 beego Author. All Rights Reserved.
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

package orm

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm/internal/logs"
)

type Log = logs.Log

// NewLog Set io.Writer to create a Logger.
func NewLog(out io.Writer) *logs.Log {
	return logs.NewLog(out)
}

// LogFunc costomer log func
var LogFunc func(query map[string]interface{})

func debugLogQueies(alias *alias, operation, query string, t time.Time, err error, args ...interface{}) {
	logMap := make(map[string]interface{})
	sub := time.Since(t) / 1e5
	elsp := float64(int(sub)) / 10.0
	logMap["cost_time"] = elsp
	flag := "OK"
	if err != nil {
		flag = "FAIL"
	}

	logMap["flag"] = flag
	con := fmt.Sprintf(" -[Queries/%s] - [  %s / %11s / %7.1fms] - [%s]", alias.Name, flag, operation, elsp, query)
	logMap["alias_name"] = alias.Name
	logMap["operation"] = operation
	logMap["query"] = query
	cons := make([]string, 0, len(args))
	for _, arg := range args {
		cons = append(cons, fmt.Sprintf("%v", arg))
	}
	if len(cons) > 0 {
		con += fmt.Sprintf(" - `%s`", strings.Join(cons, "`, `"))
		logMap["cons"] = cons
	}
	if err != nil {
		con += " - " + err.Error()
		logMap["err"] = err
	}
	logMap["sql"] = fmt.Sprintf("%s-`%s`", query, strings.Join(cons, "`, `"))
	if LogFunc != nil {
		LogFunc(logMap)
	}
	DebugLog.Println(con)
}

// statement query logger struct.
// if dev mode, use stmtQueryLog, or use stmtQuerier.
type stmtQueryLog struct {
	alias *alias
	query string
	stmt  stmtQuerier
}

var _ stmtQuerier = new(stmtQueryLog)

func (d *stmtQueryLog) Close() error {
	a := time.Now()
	err := d.stmt.Close()
	debugLogQueies(d.alias, "st.Close", d.query, a, err)
	return err
}

func (d *stmtQueryLog) Exec(args ...interface{}) (sql.Result, error) {
	return d.ExecContext(context.Background(), args...)
}

func (d *stmtQueryLog) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.stmt.ExecContext(ctx, args...)
	debugLogQueies(d.alias, "st.Exec", d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) Query(args ...interface{}) (*sql.Rows, error) {
	return d.QueryContext(context.Background(), args...)
}

func (d *stmtQueryLog) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.stmt.QueryContext(ctx, args...)
	debugLogQueies(d.alias, "st.Query", d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) QueryRow(args ...interface{}) *sql.Row {
	return d.QueryRowContext(context.Background(), args...)
}

func (d *stmtQueryLog) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.stmt.QueryRow(args...)
	debugLogQueies(d.alias, "st.QueryRow", d.query, a, nil, args...)
	return res
}

func newStmtQueryLog(alias *alias, stmt stmtQuerier, query string) stmtQuerier {
	d := new(stmtQueryLog)
	d.stmt = stmt
	d.alias = alias
	d.query = query
	return d
}

// database query logger struct.
// if dev mode, use dbQueryLog, or use dbQuerier.
type dbQueryLog struct {
	alias *alias
	db    dbQuerier
	tx    txer
	txe   txEnder
}

var (
	_ dbQuerier = new(dbQueryLog)
	_ txer      = new(dbQueryLog)
	_ txEnder   = new(dbQueryLog)
)

func (d *dbQueryLog) Prepare(query string) (*sql.Stmt, error) {
	return d.PrepareContext(context.Background(), query)
}

func (d *dbQueryLog) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	a := time.Now()
	stmt, err := d.db.PrepareContext(ctx, query)
	debugLogQueies(d.alias, "db.Prepare", query, a, err)
	return stmt, err
}

func (d *dbQueryLog) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.ExecContext(context.Background(), query, args...)
}

func (d *dbQueryLog) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.db.ExecContext(ctx, query, args...)
	debugLogQueies(d.alias, "db.Exec", query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.QueryContext(context.Background(), query, args...)
}

func (d *dbQueryLog) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.db.QueryContext(ctx, query, args...)
	debugLogQueies(d.alias, "db.Query", query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.QueryRowContext(context.Background(), query, args...)
}

func (d *dbQueryLog) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.db.QueryRowContext(ctx, query, args...)
	debugLogQueies(d.alias, "db.QueryRow", query, a, nil, args...)
	return res
}

func (d *dbQueryLog) Begin() (*sql.Tx, error) {
	return d.BeginTx(context.Background(), nil)
}

func (d *dbQueryLog) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	a := time.Now()
	tx, err := d.db.(txer).BeginTx(ctx, opts)
	debugLogQueies(d.alias, "db.BeginTx", "START TRANSACTION", a, err)
	return tx, err
}

func (d *dbQueryLog) Commit() error {
	a := time.Now()
	err := d.db.(txEnder).Commit()
	debugLogQueies(d.alias, "tx.Commit", "COMMIT", a, err)
	return err
}

func (d *dbQueryLog) Rollback() error {
	a := time.Now()
	err := d.db.(txEnder).Rollback()
	debugLogQueies(d.alias, "tx.Rollback", "ROLLBACK", a, err)
	return err
}

func (d *dbQueryLog) RollbackUnlessCommit() error {
	a := time.Now()
	err := d.db.(txEnder).RollbackUnlessCommit()
	debugLogQueies(d.alias, "tx.RollbackUnlessCommit", "ROLLBACK UNLESS COMMIT", a, err)
	return err
}

func (d *dbQueryLog) SetDB(db dbQuerier) {
	d.db = db
}

func newDbQueryLog(alias *alias, db dbQuerier) dbQuerier {
	d := new(dbQueryLog)
	d.alias = alias
	d.db = db
	return d
}
