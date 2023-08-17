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
	"errors"
	"github.com/beego/beego/v2/client/orm/internal/buffers"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/orm/internal/models"
)

func TestDbBase_InsertValueSQL(t *testing.T) {

	mi := &models.ModelInfo{
		Table: "test_table",
	}

	testCases := []struct {
		name    string
		db      *dbBase
		isMulti bool
		names   []string
		values  []interface{}

		wantRes string
	}{
		{
			name: "single insert by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},
			isMulti: false,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18},
			wantRes: "INSERT INTO `test_table` (`name`, `age`) VALUES (?, ?)",
		},
		{
			name: "single insert by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			isMulti: false,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18},
			wantRes: "INSERT INTO \"test_table\" (\"name\", \"age\") VALUES ($1, $2)",
		},
		{
			name: "multi insert by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},
			isMulti: true,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18, "test2", 19},
			wantRes: "INSERT INTO `test_table` (`name`, `age`) VALUES (?, ?), (?, ?)",
		},
		{
			name: "multi insert by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			isMulti: true,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18, "test2", 19},
			wantRes: "INSERT INTO \"test_table\" (\"name\", \"age\") VALUES ($1, $2), ($3, $4)",
		},
		{
			name: "multi insert by dbBase but values is not enough",
			db: &dbBase{
				ins: &dbBase{},
			},
			isMulti: true,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18, "test2"},
			wantRes: "INSERT INTO `test_table` (`name`, `age`) VALUES (?, ?)",
		},
		{
			name: "multi insert by dbBasePostgres but values is not enough",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			isMulti: true,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18, "test2"},
			wantRes: "INSERT INTO \"test_table\" (\"name\", \"age\") VALUES ($1, $2)",
		},
		{
			name: "single insert by dbBase but values is double to names",
			db: &dbBase{
				ins: &dbBase{},
			},
			isMulti: false,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18, "test2", 19},
			wantRes: "INSERT INTO `test_table` (`name`, `age`) VALUES (?, ?)",
		},
		{
			name: "single insert by dbBasePostgres but values is double to names",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			isMulti: false,
			names:   []string{"name", "age"},
			values:  []interface{}{"test", 18, "test2", 19},
			wantRes: "INSERT INTO \"test_table\" (\"name\", \"age\") VALUES ($1, $2)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.db.InsertValueSQL(tc.names, tc.values, tc.isMulti, mi)

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestDbBase_UpdateSQL(t *testing.T) {
	mi := &models.ModelInfo{
		Table: "test_table",
	}

	testCases := []struct {
		name string
		db   *dbBase

		setNames []string
		pkName   string

		wantRes string
	}{
		{
			name: "update by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},
			setNames: []string{"name", "age", "sender"},
			pkName:   "id",
			wantRes:  "UPDATE `test_table` SET `name` = ?, `age` = ?, `sender` = ? WHERE `id` = ?",
		},
		{
			name: "update by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			setNames: []string{"name", "age", "sender"},
			pkName:   "id",
			wantRes:  "UPDATE \"test_table\" SET \"name\" = $1, \"age\" = $2, \"sender\" = $3 WHERE \"id\" = $4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.db.UpdateSQL(tc.setNames, tc.pkName, mi)

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestDbBase_DeleteSQL(t *testing.T) {
	mi := &models.ModelInfo{
		Table: "test_table",
	}

	testCases := []struct {
		name string
		db   *dbBase

		whereCols []string

		wantRes string
	}{
		{
			name: "delete by dbBase with id",
			db: &dbBase{
				ins: &dbBase{},
			},
			whereCols: []string{"id"},
			wantRes:   "DELETE FROM `test_table` WHERE `id` = ?",
		},
		{
			name: "delete by dbBase not id",
			db: &dbBase{
				ins: &dbBase{},
			},
			whereCols: []string{"name", "age"},
			wantRes:   "DELETE FROM `test_table` WHERE `name` = ? AND `age` = ?",
		},
		{
			name: "delete by dbBasePostgres with id",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			whereCols: []string{"id"},
			wantRes:   "DELETE FROM \"test_table\" WHERE \"id\" = $1",
		},
		{
			name: "delete by dbBasePostgres not id",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			whereCols: []string{"name", "age"},
			wantRes:   "DELETE FROM \"test_table\" WHERE \"name\" = $1 AND \"age\" = $2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.db.DeleteSQL(tc.whereCols, mi)

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestDbBase_buildSetSQL(t *testing.T) {

	testCases := []struct {
		name string

		db *dbBase

		columns []string
		values  []interface{}

		wantRes    string
		wantValues []interface{}
	}{
		{
			name: "set add/mul operator by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColAdd,
					value: 12,
				},
				colValue{
					opt:   ColMultiply,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET T0.`name` = ?, T0.`age` = T0.`age` + ?, T0.`score` = T0.`score` * ?",
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set min/except operator by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColMinus,
					value: 12,
				},
				colValue{
					opt:   ColExcept,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET T0.`name` = ?, T0.`age` = T0.`age` - ?, T0.`score` = T0.`score` / ?",
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set bitRShift/bitLShift operator by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColBitRShift,
					value: 12,
				},
				colValue{
					opt:   ColBitLShift,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET T0.`name` = ?, T0.`age` = T0.`age` >> ?, T0.`score` = T0.`score` << ?",
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set bitAnd/bitOr/bitXOR operator by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},
			columns: []string{"count", "age", "score"},
			values: []interface{}{
				colValue{
					opt:   ColBitAnd,
					value: 28,
				},
				colValue{
					opt:   ColBitOr,
					value: 12,
				},
				colValue{
					opt:   ColBitXOR,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET T0.`count` = T0.`count` & ?, T0.`age` = T0.`age` | ?, T0.`score` = T0.`score` ^ ?",
			wantValues: []interface{}{int64(28), int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set add/mul operator by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColAdd,
					value: 12,
				},
				colValue{
					opt:   ColMultiply,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    `SET "name" = ?, "age" = "age" + ?, "score" = "score" * ?`,
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set min/except operator by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColMinus,
					value: 12,
				},
				colValue{
					opt:   ColExcept,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    `SET "name" = ?, "age" = "age" - ?, "score" = "score" / ?`,
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set bitRShift/bitLShift operator by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColBitRShift,
					value: 12,
				},
				colValue{
					opt:   ColBitLShift,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    `SET "name" = ?, "age" = "age" >> ?, "score" = "score" << ?`,
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set bitAnd/bitOr/bitXOR operator by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},
			columns: []string{"count", "age", "score"},
			values: []interface{}{
				colValue{
					opt:   ColBitAnd,
					value: 28,
				},
				colValue{
					opt:   ColBitOr,
					value: 12,
				},
				colValue{
					opt:   ColBitXOR,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    `SET "count" = "count" & ?, "age" = "age" | ?, "score" = "score" ^ ?`,
			wantValues: []interface{}{int64(28), int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set add/mul operator by dbBaseSqlite",
			db: &dbBase{
				ins: newdbBaseSqlite(),
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColAdd,
					value: 12,
				},
				colValue{
					opt:   ColMultiply,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET `name` = ?, `age` = `age` + ?, `score` = `score` * ?",
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set min/except operator by dbBaseSqlite",
			db: &dbBase{
				ins: newdbBaseSqlite(),
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColMinus,
					value: 12,
				},
				colValue{
					opt:   ColExcept,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET `name` = ?, `age` = `age` - ?, `score` = `score` / ?",
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set bitRShift/bitLShift operator by dbBaseSqlite",
			db: &dbBase{
				ins: newdbBaseSqlite(),
			},
			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColBitRShift,
					value: 12,
				},
				colValue{
					opt:   ColBitLShift,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET `name` = ?, `age` = `age` >> ?, `score` = `score` << ?",
			wantValues: []interface{}{"test_name", int64(12), int64(2), "test_origin_name", 18},
		},
		{
			name: "set bitAnd/bitOr/bitXOR operator by dbBaseSqlite",
			db: &dbBase{
				ins: newdbBaseSqlite(),
			},
			columns: []string{"count", "age", "score"},
			values: []interface{}{
				colValue{
					opt:   ColBitAnd,
					value: 28,
				},
				colValue{
					opt:   ColBitOr,
					value: 12,
				},
				colValue{
					opt:   ColBitXOR,
					value: 2,
				},
				"test_origin_name",
				18,
			},
			wantRes:    "SET `count` = `count` & ?, `age` = `age` | ?, `score` = `score` ^ ?",
			wantValues: []interface{}{int64(28), int64(12), int64(2), "test_origin_name", 18},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			buf := buffers.Get()
			defer buffers.Put(buf)

			tc.db.buildSetSQL(buf, tc.columns, tc.values)

			assert.Equal(t, tc.wantRes, buf.String())
			assert.Equal(t, tc.wantValues, tc.values)
		})
	}
}

func TestDbBase_UpdateBatchSQL(t *testing.T) {
	mi := &models.ModelInfo{
		Table: "test_tab",
		Fields: &models.Fields{
			Pk: &models.FieldInfo{
				Column: "test_id",
			},
		},
	}

	testCases := []struct {
		name string
		db   *dbBase

		columns []string
		values  []interface{}

		specifyIndexes string
		join           string
		where          string

		wantRes string
	}{
		{
			name: "update batch by dbBase",
			db: &dbBase{
				ins: &dbBase{},
			},

			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColAdd,
					value: 12,
				},
				colValue{
					opt:   ColMultiply,
					value: 2,
				},
				"test_origin_name",
				18,
			},

			specifyIndexes: " USE INDEX(`name`) ",
			join:           "LEFT OUTER JOIN `test_tab_2` T1 ON T1.`id` = T0.`test_id` ",
			where:          "WHERE T0.`name` = ? AND T1.`age` = ?",

			wantRes: "UPDATE `test_tab` T0  USE INDEX(`name`) LEFT OUTER JOIN `test_tab_2` T1 ON T1.`id` = T0.`test_id` SET T0.`name` = ?, T0.`age` = T0.`age` + ?, T0.`score` = T0.`score` * ? WHERE T0.`name` = ? AND T1.`age` = ?",
		},
		{
			name: "update batch by dbBasePostgres",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},

			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColAdd,
					value: 12,
				},
				colValue{
					opt:   ColMultiply,
					value: 2,
				},
				"test_origin_name",
				18,
			},

			specifyIndexes: ` USE INDEX("name") `,
			join:           `LEFT OUTER JOIN "test_tab_2" T1 ON T1."id" = T0."test_id" `,
			where:          `WHERE T0."name" = ? AND T1."age" = ?`,

			wantRes: `UPDATE "test_tab" SET "name" = $1, "age" = "age" + $2, "score" = "score" * $3 WHERE "test_id" IN ( SELECT T0."test_id" FROM "test_tab" T0  USE INDEX("name") LEFT OUTER JOIN "test_tab_2" T1 ON T1."id" = T0."test_id" WHERE T0."name" = $4 AND T1."age" = $5 )`,
		},
		{
			name: "update batch by dbBaseSqlite",
			db: &dbBase{
				ins: newdbBaseSqlite(),
			},

			columns: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				colValue{
					opt:   ColAdd,
					value: 12,
				},
				colValue{
					opt:   ColMultiply,
					value: 2,
				},
				"test_origin_name",
				18,
			},

			specifyIndexes: " USE INDEX(`name`) ",
			join:           "LEFT OUTER JOIN `test_tab_2` T1 ON T1.`id` = T0.`test_id` ",
			where:          "WHERE T0.`name` = ? AND T1.`age` = ?",

			wantRes: "UPDATE `test_tab` SET `name` = ?, `age` = `age` + ?, `score` = `score` * ? WHERE `test_id` IN ( SELECT T0.`test_id` FROM `test_tab` T0  USE INDEX(`name`) LEFT OUTER JOIN `test_tab_2` T1 ON T1.`id` = T0.`test_id` WHERE T0.`name` = ? AND T1.`age` = ? )",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.db.UpdateBatchSQL(mi, tc.columns, tc.values, tc.specifyIndexes, tc.join, tc.where)

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestDbBase_InsertOrUpdateSQL(t *testing.T) {

	mi := &models.ModelInfo{
		Table: "test_tab",
	}

	testCases := []struct {
		name string
		db   *dbBase

		names  []string
		values []interface{}
		a      *alias
		args   []string

		wantRes    string
		wantErr    error
		wantValues []interface{}
	}{
		{
			name: "test nonsupport driver",
			db: &dbBase{
				ins: newdbBaseSqlite(),
			},

			names: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				18,
				12,
			},
			a: &alias{
				Driver:     DRSqlite,
				DriverName: "sqlite3",
			},
			args: []string{
				"`age`=20",
				"`score`=`score`+1",
			},

			wantErr: errors.New("`sqlite3` nonsupport InsertOrUpdate in beego"),
			wantValues: []interface{}{
				"test_name",
				18,
				12,
			},
		},
		{
			name: "insert or update with MySQL",
			db: &dbBase{
				ins: newdbBaseMysql(),
			},

			names: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				18,
				12,
			},
			a: &alias{
				Driver:     DRMySQL,
				DriverName: "mysql",
			},
			args: []string{
				"`age`=20",
				"`score`=`score`+1",
			},

			wantRes: "INSERT INTO `test_tab` (`name`, `age`, `score`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `name`=?, `age`=20, `score`=`score`+1",
			wantValues: []interface{}{
				"test_name",
				18,
				12,
				"test_name",
			},
		},
		{
			name: "insert or update with MySQL with no args",
			db: &dbBase{
				ins: newdbBaseMysql(),
			},

			names: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				18,
				12,
			},
			a: &alias{
				Driver:     DRMySQL,
				DriverName: "mysql",
			},

			wantRes: "INSERT INTO `test_tab` (`name`, `age`, `score`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `name`=?, `age`=?, `score`=?",
			wantValues: []interface{}{
				"test_name",
				18,
				12,
				"test_name",
				18,
				12,
			},
		},
		{
			name: "insert or update with PostgreSQL normal",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},

			names: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				18,
				12,
			},
			a: &alias{
				Driver:     DRPostgres,
				DriverName: "postgres",
			},
			args: []string{
				`"name"`,
				`"score"="score_1"`,
			},

			wantRes: `INSERT INTO "test_tab" ("name", "age", "score") VALUES ($1, $2, $3) ON CONFLICT ("name") DO UPDATE SET "name"=$4, "age"=$5, "score"=(select "score_1" from test_tab where "name" = $6 )`,
			wantValues: []interface{}{
				"test_name",
				18,
				12,
				"test_name",
				18,
				"test_name",
			},
		},
		{
			name: "insert or update with PostgreSQL without conflict column",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},

			names: []string{"name", "age", "score"},
			values: []interface{}{
				"test_name",
				18,
				12,
			},
			a: &alias{
				Driver:     DRPostgres,
				DriverName: "postgres",
			},

			wantErr: errors.New("`postgres` use InsertOrUpdate must have a conflict column"),
			wantValues: []interface{}{
				"test_name",
				18,
				12,
			},
		},
		{
			name: "insert or update with PostgreSQL the conflict column is not in front of the specified column",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},

			names: []string{"score", "name", "age"},
			values: []interface{}{
				12,
				"test_name",
				18,
			},
			a: &alias{
				Driver:     DRPostgres,
				DriverName: "postgres",
			},
			args: []string{
				`"name"`,
				`"score"="score_1"`,
			},

			wantErr: errors.New("`\"name\"` must be in front of `\"score\"` in your struct"),
			wantValues: []interface{}{
				12,
				"test_name",
				18,
			},
		},
		{
			name: "insert or update with PostgreSQL the conflict column is in front of the specified column",
			db: &dbBase{
				ins: newdbBasePostgres(),
			},

			names: []string{"age", "name", "score"},
			values: []interface{}{
				18,
				"test_name",
				12,
			},
			a: &alias{
				Driver:     DRPostgres,
				DriverName: "postgres",
			},
			args: []string{
				`"name"`,
				`"score"="score_1"`,
			},

			wantRes: `INSERT INTO "test_tab" ("age", "name", "score") VALUES ($1, $2, $3) ON CONFLICT ("name") DO UPDATE SET "age"=$4, "name"=$5, "score"=(select "score_1" from test_tab where "name" = $6 )`,
			wantValues: []interface{}{
				18,
				"test_name",
				12,
				18,
				"test_name",
				"test_name",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			res, err := tc.db.InsertOrUpdateSQL(tc.names, &tc.values, mi, tc.a, tc.args...)

			assert.Equal(t, tc.wantValues, tc.values)

			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}

}
