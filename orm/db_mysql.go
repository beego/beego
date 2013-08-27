package orm

import (
	"fmt"
)

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
	"startswith":  "LIKE BINARY ?",
	"endswith":    "LIKE BINARY ?",
	"istartswith": "LIKE ?",
	"iendswith":   "LIKE ?",
}

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

type dbBaseMysql struct {
	dbBase
}

var _ dbBaser = new(dbBaseMysql)

func (d *dbBaseMysql) OperatorSql(operator string) string {
	return mysqlOperators[operator]
}

func (d *dbBaseMysql) DbTypes() map[string]string {
	return mysqlTypes
}

func (d *dbBaseMysql) ShowTablesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = DATABASE()"
}

func (d *dbBaseMysql) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE FROM information_schema.columns "+
		"WHERE table_schema = DATABASE() AND table_name = '%s'", table)
}

func (d *dbBaseMysql) IndexExists(db dbQuerier, table string, name string) bool {
	row := db.QueryRow("SELECT count(*) FROM information_schema.statistics "+
		"WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", table, name)
	var cnt int
	row.Scan(&cnt)
	return cnt > 0
}

func newdbBaseMysql() dbBaser {
	b := new(dbBaseMysql)
	b.ins = b
	return b
}
