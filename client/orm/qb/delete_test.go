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
	"github.com/beego/beego/v2/client/orm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	err := orm.RegisterDataBase("default", "sqlite3", "")
	if err != nil {
		return
	}
	db := orm.NewOrm()
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no where",
			builder: NewDeleter[TestModel](db).From(&TestModel{}),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: NewDeleter[TestModel](db).Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id` = ?;",
				Args: []interface{}{16},
			},
		},
		{
			name:    "no where combination",
			builder: NewDeleter[TestCombinedModel](db).From(&TestCombinedModel{}),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_combined_model`;",
			},
		},
		{
			name:    "where combination",
			builder: NewDeleter[TestCombinedModel](db).Where(C("CreateTime").EQ(uint64(1000))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_combined_model` WHERE `create_time` = ?;",
				Args: []interface{}{uint64(1000)},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

type BaseEntity struct {
	CreateTime uint64
	UpdateTime uint64
}

type TestCombinedModel struct {
	BaseEntity
	Id        int64 `eorm:"primary_key"`
	FirstName string
	Age       int8
	LastName  *string
}
