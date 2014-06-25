// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie, slene

package orm

import (
	"database/sql"
	"fmt"
)

// sqlite operators.
var sqliteOperators = map[string]string{
	"exact":       "= ?",
	"iexact":      "LIKE ? ESCAPE '\\'",
	"contains":    "LIKE ? ESCAPE '\\'",
	"icontains":   "LIKE ? ESCAPE '\\'",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"startswith":  "LIKE ? ESCAPE '\\'",
	"endswith":    "LIKE ? ESCAPE '\\'",
	"istartswith": "LIKE ? ESCAPE '\\'",
	"iendswith":   "LIKE ? ESCAPE '\\'",
}

// sqlite column types.
var sqliteTypes = map[string]string{
	"auto":            "integer NOT NULL PRIMARY KEY AUTOINCREMENT",
	"pk":              "NOT NULL PRIMARY KEY",
	"bool":            "bool",
	"string":          "varchar(%d)",
	"string-text":     "text",
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
	"float64":         "real",
	"float64-decimal": "decimal",
}

// sqlite dbBaser.
type dbBaseSqlite struct {
	dbBase
}

var _ dbBaser = new(dbBaseSqlite)

// get sqlite operator.
func (d *dbBaseSqlite) OperatorSql(operator string) string {
	return sqliteOperators[operator]
}

// generate functioned sql for sqlite.
// only support DATE(text).
func (d *dbBaseSqlite) GenerateOperatorLeftCol(fi *fieldInfo, operator string, leftCol *string) {
	if fi.fieldType == TypeDateField {
		*leftCol = fmt.Sprintf("DATE(%s)", *leftCol)
	}
}

// unable updating joined record in sqlite.
func (d *dbBaseSqlite) SupportUpdateJoin() bool {
	return false
}

// max int in sqlite.
func (d *dbBaseSqlite) MaxLimit() uint64 {
	return 9223372036854775807
}

// get column types in sqlite.
func (d *dbBaseSqlite) DbTypes() map[string]string {
	return sqliteTypes
}

// get show tables sql in sqlite.
func (d *dbBaseSqlite) ShowTablesQuery() string {
	return "SELECT name FROM sqlite_master WHERE type = 'table'"
}

// get columns in sqlite.
func (d *dbBaseSqlite) GetColumns(db dbQuerier, table string) (map[string][3]string, error) {
	query := d.ins.ShowColumnsQuery(table)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	columns := make(map[string][3]string)
	for rows.Next() {
		var tmp, name, typ, null sql.NullString
		err := rows.Scan(&tmp, &name, &typ, &null, &tmp, &tmp)
		if err != nil {
			return nil, err
		}
		columns[name.String] = [3]string{name.String, typ.String, null.String}
	}

	return columns, nil
}

// get show columns sql in sqlite.
func (d *dbBaseSqlite) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("pragma table_info('%s')", table)
}

// check index exist in sqlite.
func (d *dbBaseSqlite) IndexExists(db dbQuerier, table string, name string) bool {
	query := fmt.Sprintf("PRAGMA index_list('%s')", table)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tmp, index sql.NullString
		rows.Scan(&tmp, &index, &tmp)
		if name == index.String {
			return true
		}
	}
	return false
}

// create new sqlite dbBaser.
func newdbBaseSqlite() dbBaser {
	b := new(dbBaseSqlite)
	b.ins = b
	return b
}
