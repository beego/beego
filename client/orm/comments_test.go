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
	"testing"

	"github.com/beego/beego/v2/client/orm/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestQueryComments(t *testing.T) {
	qc := NewQueryComments()

	// Test empty comments
	assert.Equal(t, "", qc.String())

	// Test single comment
	qc.AddComment("test comment") // Renamed from Add
	assert.Equal(t, "/* test comment */ ", qc.String())

	// Test multiple comments
	qc.AddComment("another comment") // Renamed from Add
	assert.Equal(t, "/* test comment; another comment */ ", qc.String())

	// Test clear
	qc.ClearComments() // Renamed from Clear
	assert.Equal(t, "", qc.String())
}

type CommentUser struct {
	Id   int    `orm:"auto"`
	Name string `orm:"size(100)"`
}

func TestQueryCommentsWithOrm(t *testing.T) {
	RegisterModel(new(CommentUser))
	err := RegisterDriver("sqlite3", DRSqlite)
	if err != nil {
		t.Fatal(err)
	}

	// Use a unique database name for this test
	err = RegisterDataBase("comments_test_db", "sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create tables
	err = RunSyncdb("comments_test_db", false, false)
	if err != nil {
		t.Fatal(err)
	}

	// Test comments in queries
	DefaultQueryComments.AddComment("Test comment") // Renamed from Add
	defer DefaultQueryComments.ClearComments()      // Renamed from Clear

	// Create test components
	mi := newModelInfo()
	dBase := newDbBase()
	dBase.ins = dBase

	// Test Insert
	sql := dBase.InsertValueSQL([]string{"name"}, []interface{}{"test"}, false, mi)
	if !assert.Contains(t, sql, "/* Test comment */") {
		t.Logf("Insert SQL: %s", sql)
	}

	// Test Update
	sql = dBase.UpdateSQL([]string{"name"}, "id", mi)
	if !assert.Contains(t, sql, "/* Test comment */") {
		t.Logf("Update SQL: %s", sql)
	}

	// Test Delete
	sql = dBase.DeleteSQL([]string{"id"}, mi)
	if !assert.Contains(t, sql, "/* Test comment */") {
		t.Logf("Delete SQL: %s", sql)
	}

	// Test Select
	tables := newDbTables(mi, dBase)
	var tCols []string
	cond := NewCondition()
	qs := newQuerySet(nil, mi).(*querySet)
	sql, _ = dBase.readBatchSQL(tables, tCols, cond, *qs, mi, DefaultTimeLoc)
	if !assert.Contains(t, sql, "/* Test comment */") {
		t.Logf("Select SQL: %s", sql)
	}
}

func newModelInfo() *models.ModelInfo {
	fields := models.NewFields()

	pkField := &models.FieldInfo{
		Name:   "Id",
		Column: "id",
		Auto:   true,
		Pk:     true,
		DBcol:  true,
	}
	fields.Add(pkField)

	nameField := &models.FieldInfo{
		Name:   "Name",
		Column: "name",
		DBcol:  true,
	}
	fields.Add(nameField)

	info := &models.ModelInfo{
		Table:    "comment_user",
		FullName: "orm.CommentUser",
		Fields:   fields,
	}
	return info
}

// newDbBase creates a new dbBase instance for testing
func newDbBase() *dbBase {
	return &dbBase{}
}
