package orm

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

type dbBaseSqlite struct {
	dbBase
}

var _ dbBaser = new(dbBaseSqlite)

func (d *dbBaseSqlite) OperatorSql(operator string) string {
	return sqliteOperators[operator]
}

func (d *dbBaseSqlite) SupportUpdateJoin() bool {
	return false
}

func (d *dbBaseSqlite) MaxLimit() uint64 {
	return 9223372036854775807
}

func newdbBaseSqlite() dbBaser {
	b := new(dbBaseSqlite)
	b.ins = b
	return b
}
