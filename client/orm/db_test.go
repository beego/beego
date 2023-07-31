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
