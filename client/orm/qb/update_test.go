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

package qb

import (
	"database/sql"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/qb/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdater_SetForCombination(t *testing.T) {
	u := &User{
		Id: 12,
		Person: Person{
			FirstName: "Tom",
			Age:       int8(18),
			LastName:  sql.NullString{String: "Jerry", Valid: true},
		},
	}
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCases := []CommonTestCase{
		{
			name:     "no set",
			builder:  NewUpdater[User](db).Update(u),
			wantSql:  "UPDATE `user` SET `id`=?,`first_name`=?,`age`=?,`last_name`=?;",
			wantArgs: []interface{}{int64(12), "Tom", int8(18), sql.NullString{String: "Jerry", Valid: true}},
		},
		{
			name:     "set columns",
			builder:  NewUpdater[User](db).Update(u).Set(Columns("FirstName", "Age")),
			wantSql:  "UPDATE `user` SET `first_name`=?,`age`=?;",
			wantArgs: []interface{}{"Tom", int8(18)},
		},
		{
			name:    "set invalid columns",
			builder: NewUpdater[User](db).Update(u).Set(Columns("FirstNameInvalid", "Age")),
			wantErr: errs.NewErrUnknownField("FirstNameInvalid"),
		},
		{
			name:     "set c2",
			builder:  NewUpdater[User](db).Update(u).Set(C("FirstName"), C("Age")),
			wantSql:  "UPDATE `user` SET `first_name`=?,`age`=?;",
			wantArgs: []interface{}{"Tom", int8(18)},
		},

		{
			name:    "set invalid c2",
			builder: NewUpdater[User](db).Update(u).Set(C("FirstNameInvalid"), C("Age")),
			wantErr: errs.NewErrUnknownField("FirstNameInvalid"),
		},

		{
			name:     "set assignment",
			builder:  NewUpdater[User](db).Update(u).Set(C("FirstName"), Assign("Age", 30)),
			wantSql:  "UPDATE `user` SET `first_name`=?,`age`=?;",
			wantArgs: []interface{}{"Tom", 30},
		},
		{
			name:    "set invalid assignment",
			builder: NewUpdater[User](db).Update(u).Set(C("FirstName"), Assign("InvalidAge", 30)),
			wantErr: errs.NewErrUnknownField("InvalidAge"),
		},
		{
			name:     "set age+1",
			builder:  NewUpdater[User](db).Update(u).Set(C("FirstName"), Assign("Age", C("Age").Add(1))),
			wantSql:  "UPDATE `user` SET `first_name`=?,`age`=(`age`+?);",
			wantArgs: []interface{}{"Tom", 1},
		},
		{
			name:     "set age=id+1",
			builder:  NewUpdater[User](db).Update(u).Set(C("FirstName"), Assign("Age", C("Id").Add(10))),
			wantSql:  "UPDATE `user` SET `first_name`=?,`age`=(`id`+?);",
			wantArgs: []interface{}{"Tom", 10},
		},
		{
			name:     "set age=id+(age*100)",
			builder:  NewUpdater[User](db).Update(u).Set(C("FirstName"), Assign("Age", C("Id").Add(C("Age").Multi(100)))),
			wantSql:  "UPDATE `user` SET `first_name`=?,`age`=(`id`+(`age`*?));",
			wantArgs: []interface{}{"Tom", 100},
		},
		{
			name:     "set age=(id+(age*100))*110",
			builder:  NewUpdater[User](db).Update(u).Set(C("FirstName"), Assign("Age", C("Id").Add(C("Age").Multi(100)).Multi(110))),
			wantSql:  "UPDATE `user` SET `first_name`=?,`age`=((`id`+(`age`*?))*?);",
			wantArgs: []interface{}{"Tom", 100, 110},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := tc.builder.Build()
			assert.Equal(t, err, c.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, c.wantSql, query.SQL)
			assert.Equal(t, c.wantArgs, query.Args)
		})
	}
}

func TestUpdater_Set(t *testing.T) {
	type UserPerson struct {
		FirstName string
		Age       *int8
		LastName  sql.NullString
	}

	tm := &TestModel{
		Id:        12,
		FirstName: "Tom",
		Age:       18,
		LastName:  sql.NullString{String: "Jerry", Valid: true},
	}
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCases := []CommonTestCase{
		{
			name:     "no set and update",
			builder:  NewUpdater[TestModel](db),
			wantSql:  "UPDATE `test_model` SET `id`=?,`first_name`=?,`age`=?,`last_name`=?;",
			wantArgs: []interface{}{int64(0), "", int8(0), sql.NullString{}},
		},
		{
			name: "no set",
			builder: NewUpdater[TestModel](db).Update(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       18,
			}),
			wantSql:  "UPDATE `test_model` SET `id`=?,`first_name`=?,`age`=?,`last_name`=?;",
			wantArgs: []interface{}{int64(12), "Tom", int8(18), sql.NullString{}},
		},
		{
			name:     "set columns",
			builder:  NewUpdater[TestModel](db).Update(tm).Set(Columns("FirstName", "Age")),
			wantSql:  "UPDATE `test_model` SET `first_name`=?,`age`=?;",
			wantArgs: []interface{}{"Tom", int8(18)},
		},
		{
			name:    "set invalid columns",
			builder: NewUpdater[TestModel](db).Update(tm).Set(Columns("FirstNameInvalid", "Age")),
			wantErr: errs.NewErrUnknownField("FirstNameInvalid"),
		},
		{
			name:     "set c2",
			builder:  NewUpdater[TestModel](db).Update(tm).Set(C("FirstName"), C("Age")),
			wantSql:  "UPDATE `test_model` SET `first_name`=?,`age`=?;",
			wantArgs: []interface{}{"Tom", int8(18)},
		},
		{
			name:    "set invalid c2",
			builder: NewUpdater[TestModel](db).Update(tm).Set(C("FirstNameInvalid"), C("Age")),
			wantErr: errs.NewErrUnknownField("FirstNameInvalid"),
		},
		{
			name:     "set assignment",
			builder:  NewUpdater[TestModel](db).Update(tm).Set(C("FirstName"), Assign("Age", 30)),
			wantSql:  "UPDATE `test_model` SET `first_name`=?,`age`=?;",
			wantArgs: []interface{}{"Tom", 30},
		},
		{
			name:    "set invalid assignment",
			builder: NewUpdater[TestModel](db).Update(tm).Set(C("FirstName"), Assign("InvalidAge", 30)),
			wantErr: errs.NewErrUnknownField("InvalidAge"),
		},
		{
			name:     "set age+1",
			builder:  NewUpdater[TestModel](db).Update(tm).Set(C("FirstName"), Assign("Age", C("Age").Add(1))),
			wantSql:  "UPDATE `test_model` SET `first_name`=?,`age`=(`age`+?);",
			wantArgs: []interface{}{"Tom", 1},
		},
		{
			name:     "set age=id+1",
			builder:  NewUpdater[TestModel](db).Update(tm).Set(C("FirstName"), Assign("Age", C("Id").Add(10))),
			wantSql:  "UPDATE `test_model` SET `first_name`=?,`age`=(`id`+?);",
			wantArgs: []interface{}{"Tom", 10},
		},
		{
			name:     "set age=id+(age*100)+10",
			builder:  NewUpdater[TestModel](db).Update(tm).Set(C("FirstName"), Assign("Age", C("Id").Add(C("Age").Multi(100)).Add(10))),
			wantSql:  "UPDATE `test_model` SET `first_name`=?,`age`=((`id`+(`age`*?))+?);",
			wantArgs: []interface{}{"Tom", 100, 10},
		},
		{
			name:     "set age=(id+(age*100))*110",
			builder:  NewUpdater[TestModel](db).Update(tm).Set(C("FirstName"), Assign("Age", C("Id").Add(C("Age").Multi(100)).Multi(110))),
			wantSql:  "UPDATE `test_model` SET `first_name`=?,`age`=((`id`+(`age`*?))*?);",
			wantArgs: []interface{}{"Tom", 100, 110},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := tc.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, c.wantSql, query.SQL)
			assert.Equal(t, c.wantArgs, query.Args)
		})
	}
}

type Person struct {
	FirstName string
	Age       int8
	LastName  sql.NullString
}

type User struct {
	Id int64 `eorm:"auto_increment,primary_key"`
	Person
}

type CommonTestCase struct {
	name     string
	builder  QueryBuilder
	wantArgs []interface{}
	wantSql  string
	wantErr  error
}
