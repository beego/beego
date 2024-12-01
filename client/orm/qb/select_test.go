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
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/qb/errs"
)

func TestSelector_RawAndWhereMap(t *testing.T) {
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCase := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "WhereRaw",
			q:    NewSelector[TestModel](db).WhereRaw("`age` = ? and `first_name` = ?", 18, "sep"),
			wantQuery: &Query{
				// There are two spaces at the end because we use predicate but not predicate.op
				SQL:  "SELECT * FROM `test_model` WHERE `age` = ? and `first_name` = ?  ;",
				Args: []any{18, "sep"},
			},
		},
		// The WhereMap test might fail because the traversal of the map is unordered.
		{
			name: "WhereMap",
			q:    NewSelector[TestModel](db).WhereMap(map[string]any{"Age": 18}),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name: "WhereMapAndWhereRaw",
			q:    NewSelector[TestModel](db).WhereMap(map[string]any{"Age": 18}).WhereRaw("`id` = ? and `last_name` = ?", 1, "join"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`id` = ? and `last_name` = ?  );",
				Args: []any{18, 1, "join"},
			},
		},

		{
			name: "Where_WhereMap_WhereRaw",
			q:    NewSelector[TestModel](db).Where(C("LastName").EQ("join")).WhereMap(map[string]any{"Age": 18}).WhereRaw("`id` = ?", 1),
			wantQuery: &Query{
				// There are two spaces before NOT because we did not perform any special processing on NOT
				SQL:  "SELECT * FROM `test_model` WHERE ((`last_name` = ?) AND (`age` = ?)) AND (`id` = ?  );",
				Args: []any{"join", 18, 1},
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

func TestSelector_Build(t *testing.T) {
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCase := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "no from",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		}, {
			name: "from",
			q:    NewSelector[TestModel](db).From("from_test"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `from_test`;",
				Args: nil,
			},
		},
		{
			name: "no from",
			q:    NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		}, {
			name: "test_db",
			q:    NewSelector[TestModel](db).From("`test_db`.`db_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`db_model`;",
				Args: nil,
			},
		}, {
			name: "single and simple predicate",
			q: NewSelector[TestModel](db).From("`test_model_t`").
				Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `Id` = ?;",
				Args: []any{1},
			},
		},
		{
			name: "multiple predicates",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "and",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "or",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`Age` > ?) OR (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "not",
			q:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				// There are two spaces before NOT because we did not perform any special processing on NOT
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`Age` > ?);",
				Args: []any{18},
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

func TestSelector_OffsetLimit(t *testing.T) {
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "offset only",
			q:    NewSelector[TestModel](db).Offset(10),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` OFFSET ?;",
				Args: []any{10},
			},
		},
		{
			name: "limit only",
			q:    NewSelector[TestModel](db).Limit(10),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` LIMIT ?;",
				Args: []any{10},
			},
		},
		{
			name: "limit offset",
			q:    NewSelector[TestModel](db).Limit(20).Offset(10),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` LIMIT ? OFFSET ?;",
				Args: []any{20, 10},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_OrderBy(t *testing.T) {
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "none",
			q:    NewSelector[TestModel](db).OrderBy(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name: "single",
			q:    NewSelector[TestModel](db).OrderBy(C("Age")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age`;",
			},
		},
		{
			name: "single asc",
			q:    NewSelector[TestModel](db).OrderBy(C("Age").Asc()),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age` ASC;",
			},
		},
		{
			name: "single desc",
			q:    NewSelector[TestModel](db).OrderBy(C("Age").Desc()),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age` DESC;",
			},
		},
		{
			name: "multiple",
			q:    NewSelector[TestModel](db).OrderBy(C("Age").Asc(), C("FirstName").Desc()),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age` ASC,`first_name` DESC;",
			},
		},
		{
			name: "multiple asc",
			q:    NewSelector[TestModel](db).OrderBy(C("Age"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` ORDER BY `age`,`first_name`;",
			},
		},
		{
			name:    "invalid column",
			q:       NewSelector[TestModel](db).OrderBy(C("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestSelector_Select(t *testing.T) {
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "all",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name:    "invalid column",
			q:       NewSelector[TestModel](db).Select(Avg("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "partial columns",
			q:    NewSelector[TestModel](db).Select(C("Id"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT `id`,`first_name` FROM `test_model`;",
			},
		},
		{
			name: "avg",
			q:    NewSelector[TestModel](db).Select(Avg("Age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "raw expression",
			q:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},
		{
			name: "alias",
			q: NewSelector[TestModel](db).
				Select(C("Id").As("my_id"),
					Avg("Age").As("avg_age")),
			wantQuery: &Query{
				SQL: "SELECT `id` AS `my_id`,AVG(`age`) AS `avg_age` FROM `test_model`;",
			},
		},
		{
			name: "where ignore alias",
			q: NewSelector[TestModel](db).
				Where(C("Id").As("my_id").LT(100)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` < ?;",
				Args: []any{100},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

type TestModel struct {
	Id int64
	// ""
	FirstName string
	Age       int8
	LastName  sql.NullString
}
