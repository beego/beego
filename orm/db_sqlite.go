package orm

import (
	"fmt"
)

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

type dbBaseSqlite struct {
	dbBase
}

var _ dbBaser = new(dbBaseSqlite)

func (d *dbBaseSqlite) OperatorSql(operator string) string {
	return sqliteOperators[operator]
}

func (d *dbBaseSqlite) GenerateOperatorLeftCol(fi *fieldInfo, operator string, leftCol *string) {
	if fi.fieldType == TypeDateField {
		*leftCol = fmt.Sprintf("DATE(%s)", *leftCol)
	}
}

func (d *dbBaseSqlite) SupportUpdateJoin() bool {
	return false
}

func (d *dbBaseSqlite) MaxLimit() uint64 {
	return 9223372036854775807
}

func (d *dbBaseSqlite) DbTypes() map[string]string {
	return sqliteTypes
}

func newdbBaseSqlite() dbBaser {
	b := new(dbBaseSqlite)
	b.ins = b
	return b
}
