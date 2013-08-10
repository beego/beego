package orm

import (
	"strconv"
)

var postgresOperators = map[string]string{
	"exact":       "= ?",
	"iexact":      "= UPPER(?)",
	"contains":    "LIKE ?",
	"icontains":   "LIKE UPPER(?)",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"startswith":  "LIKE ?",
	"endswith":    "LIKE ?",
	"istartswith": "LIKE UPPER(?)",
	"iendswith":   "LIKE UPPER(?)",
}

type dbBasePostgres struct {
	dbBase
}

var _ dbBaser = new(dbBasePostgres)

func (d *dbBasePostgres) OperatorSql(operator string) string {
	return postgresOperators[operator]
}

func (d *dbBasePostgres) TableQuote() string {
	return `"`
}

func (d *dbBasePostgres) ReplaceMarks(query *string) {
	q := *query
	num := 0
	for _, c := range q {
		if c == '?' {
			num += 1
		}
	}
	if num == 0 {
		return
	}
	data := make([]byte, 0, len(q)+num)
	num = 1
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == '?' {
			data = append(data, '$')
			data = append(data, []byte(strconv.Itoa(num))...)
			num += 1
		} else {
			data = append(data, c)
		}
	}
	*query = string(data)
}

// func (d *dbBasePostgres)

func newdbBasePostgres() dbBaser {
	b := new(dbBasePostgres)
	b.ins = b
	return b
}
