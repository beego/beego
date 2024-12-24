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
	"fmt"
)

func ExampleDB_BeginTx() {
	db, _ := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	if err == nil {
		fmt.Println("Begin")
	}
	// 或者 tx.Rollback()
	err = tx.Commit()
	if err == nil {
		fmt.Println("Commit")
	}
	// Output:
	// Begin
	// Commit
}

func ExampleOpen() {
	// case1 without DBOption
	db, _ := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	fmt.Printf("case1 dialect: %s\n", db.dialect.Name())

	// Output:
	// case1 dialect: SQLite
}

func ExampleNewSelector() {
	tm := &TestModel{}
	db := memoryDB()
	query, _ := NewSelector[TestModel](db).From(tm).Build()
	fmt.Printf("SQL: %s", query.SQL)
	// Output:
	// SQL: SELECT * FROM `test_model`;
}

func ExampleNewDeleter() {
	tm := &TestModel{}
	db := memoryDB()
	query, _ := NewDeleter[TestModel](db).From(tm).Build()
	fmt.Printf("SQL: %s", query.SQL)
	// Output:
	// SQL: DELETE FROM `test_model`;
}

// memoryDB 返回一个基于内存的 ORM，它使用的是 sqlite3 内存模式。
func memoryDB() *DB {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		panic(err)
	}
	return db
}
