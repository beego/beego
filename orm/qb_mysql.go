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
	"strings"
)

type MySQLQueryBuilder struct {
	QueryString []string
}

func (qw *MySQLQueryBuilder) Select(fields ...string) QueryWriter {
	segment := fmt.Sprintf("SELECT %s", strings.Join(fields, ", "))
	qw.QueryString = append(qw.QueryString, segment)
	return qw
}

func (qw *MySQLQueryBuilder) From(tables ...string) QueryWriter {
	segment := fmt.Sprintf("FROM %s", strings.Join(tables, ", "))
	qw.QueryString = append(qw.QueryString, segment)
	return qw
}

func (qw *MySQLQueryBuilder) Where(cond string) QueryWriter {
	qw.QueryString = append(qw.QueryString, "WHERE "+cond)
	return qw
}

func (qw *MySQLQueryBuilder) LimitOffset(limit int, offset int) QueryWriter {
	qw.QueryString = append(qw.QueryString, fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset))
	return qw
}

func (qw *MySQLQueryBuilder) InnerJoin(table string) QueryWriter {
	qw.QueryString = append(qw.QueryString, "INNER JOIN "+table)
	return qw
}

func (qw *MySQLQueryBuilder) LeftJoin(table string) QueryWriter {
	qw.QueryString = append(qw.QueryString, "LEFT JOIN "+table)
	return qw
}

func (qw *MySQLQueryBuilder) On(cond string) QueryWriter {
	qw.QueryString = append(qw.QueryString, "ON "+cond)
	return qw
}

func (qw *MySQLQueryBuilder) And(cond string) QueryWriter {
	qw.QueryString = append(qw.QueryString, "AND "+cond)
	return qw
}

func (qw *MySQLQueryBuilder) Or(cond string) QueryWriter {
	qw.QueryString = append(qw.QueryString, "OR "+cond)
	return qw
}

func (qw *MySQLQueryBuilder) In(vals ...string) QueryWriter {
	segment := fmt.Sprintf("IN (%s)", strings.Join(vals, ", "))
	qw.QueryString = append(qw.QueryString, segment)
	return qw
}

func (qw *MySQLQueryBuilder) Subquery(sub string, alias string) string {
	return fmt.Sprintf("(%s) AS %s", sub, alias)
}

func (qw *MySQLQueryBuilder) String() string {
	return strings.Join(qw.QueryString, " ")
}
