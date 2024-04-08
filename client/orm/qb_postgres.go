package orm

import (
	"fmt"
	"strconv"
	"strings"
)

var quote string = `"`

// PostgresQueryBuilder is the SQL build
type PostgresQueryBuilder struct {
	tokens []string
}

func processingStr(str []string) string {
	s := strings.Join(str, `","`)
	s = fmt.Sprintf("%s%s%s", quote, s, quote)
	return s
}

// Select will join the Fields
func (qb *PostgresQueryBuilder) Select(fields ...string) QueryBuilder {
	var str string
	n := len(fields)

	if fields[0] == "*" {
		str = "*"
	} else {
		for i := 0; i < n; i++ {
			sli := strings.Split(fields[i], ".")
			s := strings.Join(sli, `"."`)
			s = fmt.Sprintf("%s%s%s", quote, s, quote)
			if n == 1 || i == n-1 {
				str += s
			} else {
				str += s + ","
			}
		}
	}

	qb.tokens = append(qb.tokens, "SELECT", str)
	return qb
}

// ForUpdate add the FOR UPDATE clause
func (qb *PostgresQueryBuilder) ForUpdate() QueryBuilder {
	qb.tokens = append(qb.tokens, "FOR UPDATE")
	return qb
}

// From join the tables
func (qb *PostgresQueryBuilder) From(tables ...string) QueryBuilder {
	str := processingStr(tables)
	qb.tokens = append(qb.tokens, "FROM", str)
	return qb
}

// InnerJoin INNER JOIN the table
func (qb *PostgresQueryBuilder) InnerJoin(table string) QueryBuilder {
	str := fmt.Sprintf("%s%s%s", quote, table, quote)
	qb.tokens = append(qb.tokens, "INNER JOIN", str)
	return qb
}

// LeftJoin LEFT JOIN the table
func (qb *PostgresQueryBuilder) LeftJoin(table string) QueryBuilder {
	str := fmt.Sprintf("%s%s%s", quote, table, quote)
	qb.tokens = append(qb.tokens, "LEFT JOIN", str)
	return qb
}

// RightJoin RIGHT JOIN the table
func (qb *PostgresQueryBuilder) RightJoin(table string) QueryBuilder {
	str := fmt.Sprintf("%s%s%s", quote, table, quote)
	qb.tokens = append(qb.tokens, "RIGHT JOIN", str)
	return qb
}

// On join with on cond
func (qb *PostgresQueryBuilder) On(cond string) QueryBuilder {
	var str string
	cond = strings.Replace(cond, " ", "", -1)
	slice := strings.Split(cond, "=")
	for i := 0; i < len(slice); i++ {
		sli := strings.Split(slice[i], ".")
		s := strings.Join(sli, `"."`)
		s = fmt.Sprintf("%s%s%s", quote, s, quote)
		if i == 0 {
			str = s + " =" + " "
		} else {
			str += s
		}
	}

	qb.tokens = append(qb.tokens, "ON", str)
	return qb
}

// Where join the Where cond
func (qb *PostgresQueryBuilder) Where(cond string) QueryBuilder {
	qb.tokens = append(qb.tokens, "WHERE", cond)
	return qb
}

// And join the and cond
func (qb *PostgresQueryBuilder) And(cond string) QueryBuilder {
	qb.tokens = append(qb.tokens, "AND", cond)
	return qb
}

// Or join the or cond
func (qb *PostgresQueryBuilder) Or(cond string) QueryBuilder {
	qb.tokens = append(qb.tokens, "OR", cond)
	return qb
}

// In join the IN (vals)
func (qb *PostgresQueryBuilder) In(vals ...string) QueryBuilder {
	qb.tokens = append(qb.tokens, "IN", "(", strings.Join(vals, CommaSpace), ")")
	return qb
}

// OrderBy join the Order by Fields
func (qb *PostgresQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	str := processingStr(fields)
	qb.tokens = append(qb.tokens, "ORDER BY", str)
	return qb
}

// Asc join the asc
func (qb *PostgresQueryBuilder) Asc() QueryBuilder {
	qb.tokens = append(qb.tokens, "ASC")
	return qb
}

// Desc join the desc
func (qb *PostgresQueryBuilder) Desc() QueryBuilder {
	qb.tokens = append(qb.tokens, "DESC")
	return qb
}

// Limit join the limit num
func (qb *PostgresQueryBuilder) Limit(limit int) QueryBuilder {
	qb.tokens = append(qb.tokens, "LIMIT", strconv.Itoa(limit))
	return qb
}

// Offset join the offset num
func (qb *PostgresQueryBuilder) Offset(offset int) QueryBuilder {
	qb.tokens = append(qb.tokens, "OFFSET", strconv.Itoa(offset))
	return qb
}

// GroupBy join the Group by Fields
func (qb *PostgresQueryBuilder) GroupBy(fields ...string) QueryBuilder {
	str := processingStr(fields)
	qb.tokens = append(qb.tokens, "GROUP BY", str)
	return qb
}

// Having join the Having cond
func (qb *PostgresQueryBuilder) Having(cond string) QueryBuilder {
	qb.tokens = append(qb.tokens, "HAVING", cond)
	return qb
}

// Update join the update table
func (qb *PostgresQueryBuilder) Update(tables ...string) QueryBuilder {
	str := processingStr(tables)
	qb.tokens = append(qb.tokens, "UPDATE", str)
	return qb
}

// Set join the Set kv
func (qb *PostgresQueryBuilder) Set(kv ...string) QueryBuilder {
	qb.tokens = append(qb.tokens, "SET", strings.Join(kv, CommaSpace))
	return qb
}

// Delete join the Delete tables
func (qb *PostgresQueryBuilder) Delete(tables ...string) QueryBuilder {
	qb.tokens = append(qb.tokens, "DELETE")
	if len(tables) != 0 {
		str := processingStr(tables)
		qb.tokens = append(qb.tokens, str)
	}
	return qb
}

// InsertInto join the insert SQL
func (qb *PostgresQueryBuilder) InsertInto(table string, fields ...string) QueryBuilder {
	str := fmt.Sprintf("%s%s%s", quote, table, quote)
	qb.tokens = append(qb.tokens, "INSERT INTO", str)
	if len(fields) != 0 {
		fieldsStr := strings.Join(fields, CommaSpace)
		qb.tokens = append(qb.tokens, "(", fieldsStr, ")")
	}
	return qb
}

// Values join the Values(vals)
func (qb *PostgresQueryBuilder) Values(vals ...string) QueryBuilder {
	valsStr := strings.Join(vals, CommaSpace)
	qb.tokens = append(qb.tokens, "VALUES", "(", valsStr, ")")
	return qb
}

// Subquery join the sub as alias
func (qb *PostgresQueryBuilder) Subquery(sub string, alias string) string {
	return fmt.Sprintf("(%s) AS %s", sub, alias)
}

// String join All tokens
func (qb *PostgresQueryBuilder) String() string {
	s := strings.Join(qb.tokens, " ")
	qb.tokens = qb.tokens[:0]
	return s
}
