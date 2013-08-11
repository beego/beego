package orm

import (
	"fmt"
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

func (d *dbBasePostgres) GenerateOperatorLeftCol(operator string, leftCol *string) {
	switch operator {
	case "contains", "startswith", "endswith":
		*leftCol = fmt.Sprintf("%s::text", *leftCol)
	case "iexact", "icontains", "istartswith", "iendswith":
		*leftCol = fmt.Sprintf("UPPER(%s::text)", *leftCol)
	}
}

func (d *dbBasePostgres) SupportUpdateJoin() bool {
	return false
}

func (d *dbBasePostgres) MaxLimit() uint64 {
	return 0
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

func (d *dbBasePostgres) HasReturningID(mi *modelInfo, query *string) (has bool) {
	if mi.fields.pk.auto {
		if query != nil {
			*query = fmt.Sprintf(`%s RETURNING "%s"`, *query, mi.fields.pk.column)
		}
		has = true
	}
	return
}

func newdbBasePostgres() dbBaser {
	b := new(dbBasePostgres)
	b.ins = b
	return b
}
