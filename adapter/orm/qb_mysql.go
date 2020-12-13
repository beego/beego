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
	"github.com/beego/beego/client/orm"
)

// CommaSpace is the separation
const CommaSpace = orm.CommaSpace

// MySQLQueryBuilder is the SQL build
type MySQLQueryBuilder orm.MySQLQueryBuilder

// Select will join the fields
func (qb *MySQLQueryBuilder) Select(fields ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Select(fields...)
}

// ForUpdate add the FOR UPDATE clause
func (qb *MySQLQueryBuilder) ForUpdate() QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).ForUpdate()
}

// From join the tables
func (qb *MySQLQueryBuilder) From(tables ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).From(tables...)
}

// InnerJoin INNER JOIN the table
func (qb *MySQLQueryBuilder) InnerJoin(table string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).InnerJoin(table)
}

// LeftJoin LEFT JOIN the table
func (qb *MySQLQueryBuilder) LeftJoin(table string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).LeftJoin(table)
}

// RightJoin RIGHT JOIN the table
func (qb *MySQLQueryBuilder) RightJoin(table string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).RightJoin(table)
}

// On join with on cond
func (qb *MySQLQueryBuilder) On(cond string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).On(cond)
}

// Where join the Where cond
func (qb *MySQLQueryBuilder) Where(cond string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Where(cond)
}

// And join the and cond
func (qb *MySQLQueryBuilder) And(cond string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).And(cond)
}

// Or join the or cond
func (qb *MySQLQueryBuilder) Or(cond string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Or(cond)
}

// In join the IN (vals)
func (qb *MySQLQueryBuilder) In(vals ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).In(vals...)
}

// OrderBy join the Order by fields
func (qb *MySQLQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).OrderBy(fields...)
}

// Asc join the asc
func (qb *MySQLQueryBuilder) Asc() QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Asc()
}

// Desc join the desc
func (qb *MySQLQueryBuilder) Desc() QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Desc()
}

// Limit join the limit num
func (qb *MySQLQueryBuilder) Limit(limit int) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Limit(limit)
}

// Offset join the offset num
func (qb *MySQLQueryBuilder) Offset(offset int) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Offset(offset)
}

// GroupBy join the Group by fields
func (qb *MySQLQueryBuilder) GroupBy(fields ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).GroupBy(fields...)
}

// Having join the Having cond
func (qb *MySQLQueryBuilder) Having(cond string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Having(cond)
}

// Update join the update table
func (qb *MySQLQueryBuilder) Update(tables ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Update(tables...)
}

// Set join the set kv
func (qb *MySQLQueryBuilder) Set(kv ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Set(kv...)
}

// Delete join the Delete tables
func (qb *MySQLQueryBuilder) Delete(tables ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Delete(tables...)
}

// InsertInto join the insert SQL
func (qb *MySQLQueryBuilder) InsertInto(table string, fields ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).InsertInto(table, fields...)
}

// Values join the Values(vals)
func (qb *MySQLQueryBuilder) Values(vals ...string) QueryBuilder {
	return (*orm.MySQLQueryBuilder)(qb).Values(vals...)
}

// Subquery join the sub as alias
func (qb *MySQLQueryBuilder) Subquery(sub string, alias string) string {
	return (*orm.MySQLQueryBuilder)(qb).Subquery(sub, alias)
}

// String join all Tokens
func (qb *MySQLQueryBuilder) String() string {
	return (*orm.MySQLQueryBuilder)(qb).String()
}
