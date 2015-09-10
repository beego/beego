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
	"strconv"
	"strings"
)

type TiDBQueryBuilder struct {
	Tokens []string
}

func (qb *TiDBQueryBuilder) Select(fields ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "SELECT", strings.Join(fields, COMMA_SPACE))
	return qb
}

func (qb *TiDBQueryBuilder) From(tables ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "FROM", strings.Join(tables, COMMA_SPACE))
	return qb
}

func (qb *TiDBQueryBuilder) InnerJoin(table string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "INNER JOIN", table)
	return qb
}

func (qb *TiDBQueryBuilder) LeftJoin(table string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "LEFT JOIN", table)
	return qb
}

func (qb *TiDBQueryBuilder) RightJoin(table string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "RIGHT JOIN", table)
	return qb
}

func (qb *TiDBQueryBuilder) On(cond string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "ON", cond)
	return qb
}

func (qb *TiDBQueryBuilder) Where(cond string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "WHERE", cond)
	return qb
}

func (qb *TiDBQueryBuilder) And(cond string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "AND", cond)
	return qb
}

func (qb *TiDBQueryBuilder) Or(cond string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "OR", cond)
	return qb
}

func (qb *TiDBQueryBuilder) In(vals ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "IN", "(", strings.Join(vals, COMMA_SPACE), ")")
	return qb
}

func (qb *TiDBQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "ORDER BY", strings.Join(fields, COMMA_SPACE))
	return qb
}

func (qb *TiDBQueryBuilder) Asc() QueryBuilder {
	qb.Tokens = append(qb.Tokens, "ASC")
	return qb
}

func (qb *TiDBQueryBuilder) Desc() QueryBuilder {
	qb.Tokens = append(qb.Tokens, "DESC")
	return qb
}

func (qb *TiDBQueryBuilder) Limit(limit int) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "LIMIT", strconv.Itoa(limit))
	return qb
}

func (qb *TiDBQueryBuilder) Offset(offset int) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "OFFSET", strconv.Itoa(offset))
	return qb
}

func (qb *TiDBQueryBuilder) GroupBy(fields ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "GROUP BY", strings.Join(fields, COMMA_SPACE))
	return qb
}

func (qb *TiDBQueryBuilder) Having(cond string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "HAVING", cond)
	return qb
}

func (qb *TiDBQueryBuilder) Update(tables ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "UPDATE", strings.Join(tables, COMMA_SPACE))
	return qb
}

func (qb *TiDBQueryBuilder) Set(kv ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "SET", strings.Join(kv, COMMA_SPACE))
	return qb
}

func (qb *TiDBQueryBuilder) Delete(tables ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "DELETE")
	if len(tables) != 0 {
		qb.Tokens = append(qb.Tokens, strings.Join(tables, COMMA_SPACE))
	}
	return qb
}

func (qb *TiDBQueryBuilder) InsertInto(table string, fields ...string) QueryBuilder {
	qb.Tokens = append(qb.Tokens, "INSERT INTO", table)
	if len(fields) != 0 {
		fieldsStr := strings.Join(fields, COMMA_SPACE)
		qb.Tokens = append(qb.Tokens, "(", fieldsStr, ")")
	}
	return qb
}

func (qb *TiDBQueryBuilder) Values(vals ...string) QueryBuilder {
	valsStr := strings.Join(vals, COMMA_SPACE)
	qb.Tokens = append(qb.Tokens, "VALUES", "(", valsStr, ")")
	return qb
}

func (qb *TiDBQueryBuilder) Subquery(sub string, alias string) string {
	return fmt.Sprintf("(%s) AS %s", sub, alias)
}

func (qb *TiDBQueryBuilder) String() string {
	return strings.Join(qb.Tokens, " ")
}
