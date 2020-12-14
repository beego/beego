// Copyright 2015 TiDB Author. All Rights Reserved.
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
	"github.com/beego/beego/v2/client/orm"
)

// TiDBQueryBuilder is the SQL build
type TiDBQueryBuilder orm.TiDBQueryBuilder

// Select will join the fields
func (qb *TiDBQueryBuilder) Select(fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Select(fields...)
}

// ForUpdate add the FOR UPDATE clause
func (qb *TiDBQueryBuilder) ForUpdate() QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).ForUpdate()
}

// From join the tables
func (qb *TiDBQueryBuilder) From(tables ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).From(tables...)
}

// InnerJoin INNER JOIN the table
func (qb *TiDBQueryBuilder) InnerJoin(table string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).InnerJoin(table)
}

// LeftJoin LEFT JOIN the table
func (qb *TiDBQueryBuilder) LeftJoin(table string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).LeftJoin(table)
}

// RightJoin RIGHT JOIN the table
func (qb *TiDBQueryBuilder) RightJoin(table string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).RightJoin(table)
}

// On join with on cond
func (qb *TiDBQueryBuilder) On(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).On(cond)
}

// Where join the Where cond
func (qb *TiDBQueryBuilder) Where(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Where(cond)
}

// And join the and cond
func (qb *TiDBQueryBuilder) And(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).And(cond)
}

// Or join the or cond
func (qb *TiDBQueryBuilder) Or(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Or(cond)
}

// In join the IN (vals)
func (qb *TiDBQueryBuilder) In(vals ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).In(vals...)
}

// OrderBy join the Order by fields
func (qb *TiDBQueryBuilder) OrderBy(fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).OrderBy(fields...)
}

// Asc join the asc
func (qb *TiDBQueryBuilder) Asc() QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Asc()
}

// Desc join the desc
func (qb *TiDBQueryBuilder) Desc() QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Desc()
}

// Limit join the limit num
func (qb *TiDBQueryBuilder) Limit(limit int) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Limit(limit)
}

// Offset join the offset num
func (qb *TiDBQueryBuilder) Offset(offset int) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Offset(offset)
}

// GroupBy join the Group by fields
func (qb *TiDBQueryBuilder) GroupBy(fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).GroupBy(fields...)
}

// Having join the Having cond
func (qb *TiDBQueryBuilder) Having(cond string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Having(cond)
}

// Update join the update table
func (qb *TiDBQueryBuilder) Update(tables ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Update(tables...)
}

// Set join the set kv
func (qb *TiDBQueryBuilder) Set(kv ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Set(kv...)
}

// Delete join the Delete tables
func (qb *TiDBQueryBuilder) Delete(tables ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Delete(tables...)
}

// InsertInto join the insert SQL
func (qb *TiDBQueryBuilder) InsertInto(table string, fields ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).InsertInto(table, fields...)
}

// Values join the Values(vals)
func (qb *TiDBQueryBuilder) Values(vals ...string) QueryBuilder {
	return (*orm.TiDBQueryBuilder)(qb).Values(vals...)
}

// Subquery join the sub as alias
func (qb *TiDBQueryBuilder) Subquery(sub string, alias string) string {
	return (*orm.TiDBQueryBuilder)(qb).Subquery(sub, alias)
}

// String join all Tokens
func (qb *TiDBQueryBuilder) String() string {
	return (*orm.TiDBQueryBuilder)(qb).String()
}
