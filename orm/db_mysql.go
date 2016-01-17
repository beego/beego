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
	"fmt"
)

// mysql operators.
var mysqlOperators = map[string]string{
	"exact":     "= ?",
	"iexact":    "LIKE ?",
	"contains":  "LIKE BINARY ?",
	"icontains": "LIKE ?",
	// "regex":       "REGEXP BINARY ?",
	// "iregex":      "REGEXP ?",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"eq":          "= ?",
	"ne":          "!= ?",
	"startswith":  "LIKE BINARY ?",
	"endswith":    "LIKE BINARY ?",
	"istartswith": "LIKE ?",
	"iendswith":   "LIKE ?",
}

// mysql column field types.
var mysqlTypes = map[string]string{
	"auto":            "AUTO_INCREMENT NOT NULL PRIMARY KEY",
	"pk":              "NOT NULL PRIMARY KEY",
	"bool":            "bool",
	"string":          "varchar(%d)",
	"string-text":     "longtext",
	"time.Time-date":  "date",
	"time.Time":       "datetime",
	"int8":            "tinyint",
	"int16":           "smallint",
	"int32":           "integer",
	"int64":           "bigint",
	"uint8":           "tinyint unsigned",
	"uint16":          "smallint unsigned",
	"uint32":          "integer unsigned",
	"uint64":          "bigint unsigned",
	"float64":         "double precision",
	"float64-decimal": "numeric(%d, %d)",
}

// mysql dbBaser implementation.
type dbBaseMysql struct {
	dbBase
}

var _ dbBaser = new(dbBaseMysql)

// get mysql operator.
func (d *dbBaseMysql) OperatorSQL(operator string) string {
	return mysqlOperators[operator]
}

// get mysql table field types.
func (d *dbBaseMysql) DbTypes() map[string]string {
	return mysqlTypes
}

// show table sql for mysql.
func (d *dbBaseMysql) ShowTablesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = DATABASE()"
}

// show columns sql of table for mysql.
func (d *dbBaseMysql) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE FROM information_schema.columns "+
		"WHERE table_schema = DATABASE() AND table_name = '%s'", table)
}

// execute sql to check index exist.
func (d *dbBaseMysql) IndexExists(db dbQuerier, table string, name string) bool {
	row := db.QueryRow("SELECT count(*) FROM information_schema.statistics "+
		"WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", table, name)
	var cnt int
	row.Scan(&cnt)
	return cnt > 0
}

// create new mysql dbBaser.
func newdbBaseMysql() dbBaser {
	b := new(dbBaseMysql)
	b.ins = b
	return b
}
