package orm

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

type dbBaseMysql struct {
	dbBase
}

var _ dbBaser = new(dbBaseMysql)

func (d *dbBaseMysql) OperatorSql(operator string) string {
	return mysqlOperators[operator]
}

func newdbBaseMysql() dbBaser {
	b := new(dbBaseMysql)
	b.ins = b
	return b
}
