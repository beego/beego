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

type MySQLQueryBuilder struct {
	QueryTokens []string
}

func (qb *MySQLQueryBuilder) Select(fields ...string) QueryBuilder {
	segment := fmt.Sprintf("SELECT %s", strings.Join(fields, ", "))
	qb.QueryTokens = append(qb.QueryTokens, segment)
	return qb
}

func (qb *MySQLQueryBuilder) From(tables ...string) QueryBuilder {
	segment := fmt.Sprintf("FROM %s", strings.Join(tables, ", "))
	qb.QueryTokens = append(qb.QueryTokens, segment)
	return qb
}

func (qb *MySQLQueryBuilder) InnerJoin(table string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "INNER JOIN "+table)
	return qb
}

func (qb *MySQLQueryBuilder) LeftJoin(table string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "LEFT JOIN "+table)
	return qb
}

func (qb *MySQLQueryBuilder) RightJoin(table string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "RIGHT JOIN "+table)
	return qb
}

func (qb *MySQLQueryBuilder) On(cond string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "ON "+cond)
	return qb
}

func (qb *MySQLQueryBuilder) Where(cond string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "WHERE "+cond)
	return qb
}

func (qb *MySQLQueryBuilder) And(cond string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "AND "+cond)
	return qb
}

func (qb *MySQLQueryBuilder) Or(cond string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "OR "+cond)
	return qb
}

func (qb *MySQLQueryBuilder) In(vals ...string) QueryBuilder {
	segment := fmt.Sprintf("IN (%s)", strings.Join(vals, ", "))
	qb.QueryTokens = append(qb.QueryTokens, segment)
	return qb
}

func (qb *MySQLQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "ORDER BY "+strings.Join(fields, ", "))
	return qb
}

func (qb *MySQLQueryBuilder) Asc() QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "ASC")
	return qb
}

func (qb *MySQLQueryBuilder) Desc() QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "DESC")
	return qb
}

func (qb *MySQLQueryBuilder) Limit(limit int) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "LIMIT "+strconv.Itoa(limit))
	return qb
}

func (qb *MySQLQueryBuilder) Offset(offset int) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "OFFSET "+strconv.Itoa(offset))
	return qb
}

func (qb *MySQLQueryBuilder) GroupBy(fields ...string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "GROUP BY "+strings.Join(fields, ", "))
	return qb
}

func (qb *MySQLQueryBuilder) Having(cond string) QueryBuilder {
	qb.QueryTokens = append(qb.QueryTokens, "HAVING "+cond)
	return qb
}

func (qb *MySQLQueryBuilder) Subquery(sub string, alias string) string {
	return fmt.Sprintf("(%s) AS %s", sub, alias)
}

func (qb *MySQLQueryBuilder) String() string {
	return strings.Join(qb.QueryTokens, " ")
}
