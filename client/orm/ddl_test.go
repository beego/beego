// Copyright 2022
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

	"github.com/beego/beego/v2/client/orm/internal/models"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type ModelWithComments struct {
	ID       int    `orm:"column(id);description(user id)"`
	UserName string `orm:"size(30);unique;description(user name)"`
	Email    string `orm:"size(100);description(email)"`
	Password string `orm:"size(100);description(password)"`
}

type ModelWithoutComments struct {
	ID       int    `orm:"column(id)"`
	UserName string `orm:"size(30);unique"`
	Email    string `orm:"size(100)"`
	Password string `orm:"size(100)"`
}

type ModelWithEmptyComments struct {
	ID       int    `orm:"column(id);description()"`
	UserName string `orm:"size(30);unique;description()"`
	Email    string `orm:"size(100);description()"`
	Password string `orm:"size(100);description()"`
}
type ModelWithDBTypes struct {
	ID       int    `orm:"column(id);description();db_type(bigserial NOT NULL PRIMARY KEY)"`
	UserName string `orm:"size(30);unique;description()"`
	Email    string `orm:"size(100);description()"`
	Password string `orm:"size(100);description()"`
}

func TestGetDbCreateSQLWithComment(t *testing.T) {
	type TestCase struct {
		name    string
		model   interface{}
		wantSQL string
		wantErr error
	}
	al := getDbAlias("default")
	testModelCache := models.NewModelCacheHandler()
	var testCases []TestCase
	switch al.Driver {
	case DRMySQL:
		testCases = append(testCases, TestCase{name: "model with comments for MySQL", model: &ModelWithComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_with_comments` (\n    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY COMMENT 'user id',\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE COMMENT 'user name',\n    `email` varchar(100) NOT NULL DEFAULT ''  COMMENT 'email',\n    `password` varchar(100) NOT NULL DEFAULT ''  COMMENT 'password'\n) ENGINE=INNODB;", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model without comments for MySQL", model: &ModelWithoutComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithoutComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_without_comments` (\n    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n) ENGINE=INNODB;", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model with empty comments for MySQL", model: &ModelWithEmptyComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithEmptyComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_with_empty_comments` (\n    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n) ENGINE=INNODB;", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model with dpType for MySQL", model: &ModelWithDBTypes{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithDBTypes`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_with_d_b_types` (\n    `id` bigserial NOT NULL PRIMARY KEY,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
	case DRPostgres:
		testCases = append(testCases, TestCase{name: "model with comments for Postgres", model: &ModelWithComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS \"model_with_comments\" (\n    \"id\" serial NOT NULL PRIMARY KEY,\n    \"user_name\" varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    \"email\" varchar(100) NOT NULL DEFAULT '' ,\n    \"password\" varchar(100) NOT NULL DEFAULT '' \n);\nCOMMENT ON COLUMN \"model_with_comments\".\"id\" is 'user id';\nCOMMENT ON COLUMN \"model_with_comments\".\"user_name\" is 'user name';\nCOMMENT ON COLUMN \"model_with_comments\".\"email\" is 'email';\nCOMMENT ON COLUMN \"model_with_comments\".\"password\" is 'password';", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model without comments for Postgres", model: &ModelWithoutComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithoutComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS \"model_without_comments\" (\n    \"id\" serial NOT NULL PRIMARY KEY,\n    \"user_name\" varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    \"email\" varchar(100) NOT NULL DEFAULT '' ,\n    \"password\" varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model with empty comments for Postgres", model: &ModelWithEmptyComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithEmptyComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS \"model_with_empty_comments\" (\n    \"id\" serial NOT NULL PRIMARY KEY,\n    \"user_name\" varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    \"email\" varchar(100) NOT NULL DEFAULT '' ,\n    \"password\" varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model with dpType for Postgres", model: &ModelWithDBTypes{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithDBTypes`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_with_d_b_types` (\n    `id` bigserial NOT NULL PRIMARY KEY,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
	case DRSqlite:
		testCases = append(testCases, TestCase{name: "model with comments for Sqlite", model: &ModelWithComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_with_comments` (\n    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model without comments for Sqlite", model: &ModelWithoutComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithoutComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_without_comments` (\n    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model with empty comments for Sqlite", model: &ModelWithEmptyComments{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithEmptyComments`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_with_empty_comments` (\n    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
		testCases = append(testCases, TestCase{name: "model with dpType for Sqlite", model: &ModelWithDBTypes{}, wantSQL: "-- --------------------------------------------------\n--  Table Structure for `github.com/beego/beego/v2/client/orm.ModelWithDBTypes`\n-- --------------------------------------------------\nCREATE TABLE IF NOT EXISTS `model_with_d_b_types` (\n    `id` bigserial NOT NULL PRIMARY KEY,\n    `user_name` varchar(30) NOT NULL DEFAULT ''  UNIQUE,\n    `email` varchar(100) NOT NULL DEFAULT '' ,\n    `password` varchar(100) NOT NULL DEFAULT '' \n);", wantErr: nil})
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testModelCache.Clean()
			err := testModelCache.Register("", true, tc.model)
			assert.NoError(t, err)
			queries, _, err := getDbCreateSQL(testModelCache, al)
			assert.Equal(t, tc.wantSQL, queries[0])
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
