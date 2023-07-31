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
	"time"

	"github.com/beego/beego/v2/client/orm/internal/models"

	"github.com/beego/beego/v2/client/orm/clauses"
	"github.com/beego/beego/v2/client/orm/clauses/order_clause"
)

// table info struct.
type dbTable struct {
	id    int
	index string
	name  string
	names []string
	sel   bool
	inner bool
	mi    *models.ModelInfo
	fi    *models.FieldInfo
	jtl   *dbTable
}

// tables collection struct, contains some tables.
type dbTables struct {
	tablesM map[string]*dbTable
	tables  []*dbTable
	mi      *models.ModelInfo
	base    dbBaser
	skipEnd bool
}

// set table info to collection.
// if not exist, create new.
func (t *dbTables) set(names []string, mi *models.ModelInfo, fi *models.FieldInfo, inner bool) *dbTable {
	name := strings.Join(names, ExprSep)
	if j, ok := t.tablesM[name]; ok {
		j.name = name
		j.mi = mi
		j.fi = fi
		j.inner = inner
	} else {
		i := len(t.tables) + 1
		jt := &dbTable{i, fmt.Sprintf("T%d", i), name, names, false, inner, mi, fi, nil}
		t.tablesM[name] = jt
		t.tables = append(t.tables, jt)
	}
	return t.tablesM[name]
}

// add table info to collection.
func (t *dbTables) add(names []string, mi *models.ModelInfo, fi *models.FieldInfo, inner bool) (*dbTable, bool) {
	name := strings.Join(names, ExprSep)
	if _, ok := t.tablesM[name]; !ok {
		i := len(t.tables) + 1
		jt := &dbTable{i, fmt.Sprintf("T%d", i), name, names, false, inner, mi, fi, nil}
		t.tablesM[name] = jt
		t.tables = append(t.tables, jt)
		return jt, true
	}
	return t.tablesM[name], false
}

// get table info in collection.
func (t *dbTables) get(name string) (*dbTable, bool) {
	j, ok := t.tablesM[name]
	return j, ok
}

// get related Fields info in recursive depth loop.
// loop once, depth decreases one.
func (t *dbTables) loopDepth(depth int, prefix string, fi *models.FieldInfo, related []string) []string {
	if depth < 0 || fi.FieldType == RelManyToMany {
		return related
	}

	if prefix == "" {
		prefix = fi.Name
	} else {
		prefix = prefix + ExprSep + fi.Name
	}
	related = append(related, prefix)

	depth--
	for _, fi := range fi.RelModelInfo.Fields.FieldsRel {
		related = t.loopDepth(depth, prefix, fi, related)
	}

	return related
}

// parse related Fields.
func (t *dbTables) parseRelated(rels []string, depth int) {
	relsNum := len(rels)
	related := make([]string, relsNum)
	copy(related, rels)

	relDepth := depth

	if relsNum != 0 {
		relDepth = 0
	}

	relDepth--
	for _, fi := range t.mi.Fields.FieldsRel {
		related = t.loopDepth(relDepth, "", fi, related)
	}

	for i, s := range related {
		var (
			exs    = strings.Split(s, ExprSep)
			names  = make([]string, 0, len(exs))
			mmi    = t.mi
			cancel = true
			jtl    *dbTable
		)

		inner := true

		for _, ex := range exs {
			if fi, ok := mmi.Fields.GetByAny(ex); ok && fi.Rel && fi.FieldType != RelManyToMany {
				names = append(names, fi.Name)
				mmi = fi.RelModelInfo

				if fi.Null || t.skipEnd {
					inner = false
				}

				jt := t.set(names, mmi, fi, inner)
				jt.jtl = jtl

				if fi.Reverse {
					cancel = false
				}

				if cancel {
					jt.sel = depth > 0

					if i < relsNum {
						jt.sel = true
					}
				}

				jtl = jt

			} else {
				panic(fmt.Errorf("unknown model/table name `%s`", ex))
			}
		}
	}
}

// generate join string.
func (t *dbTables) getJoinSQL() (join string) {
	Q := t.base.TableQuote()

	for _, jt := range t.tables {
		if jt.inner {
			join += "INNER JOIN "
		} else {
			join += "LEFT OUTER JOIN "
		}
		var (
			table  string
			t1, t2 string
			c1, c2 string
		)
		t1 = "T0"
		if jt.jtl != nil {
			t1 = jt.jtl.index
		}
		t2 = jt.index
		table = jt.mi.Table

		switch {
		case jt.fi.FieldType == RelManyToMany || jt.fi.FieldType == RelReverseMany || jt.fi.Reverse && jt.fi.ReverseFieldInfo.FieldType == RelManyToMany:
			c1 = jt.fi.Mi.Fields.Pk.Column
			for _, ffi := range jt.mi.Fields.FieldsRel {
				if jt.fi.Mi == ffi.RelModelInfo {
					c2 = ffi.Column
					break
				}
			}
		default:
			c1 = jt.fi.Column
			c2 = jt.fi.RelModelInfo.Fields.Pk.Column

			if jt.fi.Reverse {
				c1 = jt.mi.Fields.Pk.Column
				c2 = jt.fi.ReverseFieldInfo.Column
			}
		}

		join += fmt.Sprintf("%s%s%s %s ON %s.%s%s%s = %s.%s%s%s ", Q, table, Q, t2,
			t2, Q, c2, Q, t1, Q, c1, Q)
	}
	return
}

// parse orm model struct field tag expression.
func (t *dbTables) parseExprs(mi *models.ModelInfo, exprs []string) (index, name string, info *models.FieldInfo, success bool) {
	var (
		jtl *dbTable
		fi  *models.FieldInfo
		fiN *models.FieldInfo
		mmi = mi
	)

	num := len(exprs) - 1
	var names []string

	inner := true

loopFor:
	for i, ex := range exprs {

		var ok, okN bool

		if fiN != nil {
			fi = fiN
			ok = true
			fiN = nil
		}

		if i == 0 {
			fi, ok = mmi.Fields.GetByAny(ex)
		}

		_ = okN

		if ok {

			isRel := fi.Rel || fi.Reverse

			names = append(names, fi.Name)

			switch {
			case fi.Rel:
				mmi = fi.RelModelInfo
				if fi.FieldType == RelManyToMany {
					mmi = fi.RelThroughModelInfo
				}
			case fi.Reverse:
				mmi = fi.ReverseFieldInfo.Mi
			}

			if i < num {
				fiN, okN = mmi.Fields.GetByAny(exprs[i+1])
			}

			if isRel && (!fi.Mi.IsThrough || num != i) {
				if fi.Null || t.skipEnd {
					inner = false
				}

				if t.skipEnd && okN || !t.skipEnd {
					if t.skipEnd && okN && fiN.Pk {
						goto loopEnd
					}

					jt, _ := t.add(names, mmi, fi, inner)
					jt.jtl = jtl
					jtl = jt
				}

			}

			if num != i {
				continue
			}

		loopEnd:

			if i == 0 || jtl == nil {
				index = "T0"
			} else {
				index = jtl.index
			}

			info = fi

			if jtl == nil {
				name = fi.Name
			} else {
				name = jtl.name + ExprSep + fi.Name
			}

			switch {
			case fi.Rel:

			case fi.Reverse:
				switch fi.ReverseFieldInfo.FieldType {
				case RelOneToOne, RelForeignKey:
					index = jtl.index
					info = fi.ReverseFieldInfo.Mi.Fields.Pk
					name = info.Name
				}
			}

			break loopFor

		} else {
			index = ""
			name = ""
			info = nil
			success = false
			return
		}
	}

	success = index != "" && info != nil
	return
}

// generate condition sql.
func (t *dbTables) getCondSQL(cond *Condition, sub bool, tz *time.Location) (where string, params []interface{}) {
	if cond == nil || cond.IsEmpty() {
		return
	}

	Q := t.base.TableQuote()

	mi := t.mi

	for i, p := range cond.params {
		if i > 0 {
			if p.isOr {
				where += "OR "
			} else {
				where += "AND "
			}
		}
		if p.isNot {
			where += "NOT "
		}
		if p.isCond {
			w, ps := t.getCondSQL(p.cond, true, tz)
			if w != "" {
				w = fmt.Sprintf("( %s) ", w)
			}
			where += w
			params = append(params, ps...)
		} else {
			exprs := p.exprs

			num := len(exprs) - 1
			operator := ""
			if operators[exprs[num]] {
				operator = exprs[num]
				exprs = exprs[:num]
			}

			index, _, fi, suc := t.parseExprs(mi, exprs)
			if !suc {
				panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(p.exprs, ExprSep)))
			}

			if operator == "" {
				operator = "exact"
			}

			var operSQL string
			var args []interface{}
			if p.isRaw {
				operSQL = p.sql
			} else {
				operSQL, args = t.base.GenerateOperatorSQL(mi, fi, operator, p.args, tz)
			}

			leftCol := fmt.Sprintf("%s.%s%s%s", index, Q, fi.Column, Q)
			t.base.GenerateOperatorLeftCol(fi, operator, &leftCol)

			where += fmt.Sprintf("%s %s ", leftCol, operSQL)
			params = append(params, args...)

		}
	}

	if !sub && where != "" {
		where = "WHERE " + where
	}

	return
}

// generate group sql.
func (t *dbTables) getGroupSQL(groups []string) (groupSQL string) {
	if len(groups) == 0 {
		return
	}

	Q := t.base.TableQuote()

	groupSqls := make([]string, 0, len(groups))
	for _, group := range groups {
		exprs := strings.Split(group, ExprSep)

		index, _, fi, suc := t.parseExprs(t.mi, exprs)
		if !suc {
			panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(exprs, ExprSep)))
		}

		groupSqls = append(groupSqls, fmt.Sprintf("%s.%s%s%s", index, Q, fi.Column, Q))
	}

	groupSQL = fmt.Sprintf("GROUP BY %s ", strings.Join(groupSqls, ", "))
	return
}

// generate order sql.
func (t *dbTables) getOrderSQL(orders []*order_clause.Order) (orderSQL string) {
	if len(orders) == 0 {
		return
	}

	Q := t.base.TableQuote()

	orderSqls := make([]string, 0, len(orders))
	for _, order := range orders {
		column := order.GetColumn()
		clause := strings.Split(column, clauses.ExprDot)

		if order.IsRaw() {
			if len(clause) == 2 {
				orderSqls = append(orderSqls, fmt.Sprintf("%s.%s%s%s %s", clause[0], Q, clause[1], Q, order.SortString()))
			} else if len(clause) == 1 {
				orderSqls = append(orderSqls, fmt.Sprintf("%s%s%s %s", Q, clause[0], Q, order.SortString()))
			} else {
				panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(clause, ExprSep)))
			}
		} else {
			index, _, fi, suc := t.parseExprs(t.mi, clause)
			if !suc {
				panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(clause, ExprSep)))
			}

			orderSqls = append(orderSqls, fmt.Sprintf("%s.%s%s%s %s", index, Q, fi.Column, Q, order.SortString()))
		}
	}

	orderSQL = fmt.Sprintf("ORDER BY %s ", strings.Join(orderSqls, ", "))
	return
}

// generate limit sql.
func (t *dbTables) getLimitSQL(mi *models.ModelInfo, offset int64, limit int64) (limits string) {
	if limit == 0 {
		limit = int64(DefaultRowsLimit)
	}
	if limit < 0 {
		// no limit
		if offset > 0 {
			maxLimit := t.base.MaxLimit()
			if maxLimit == 0 {
				limits = fmt.Sprintf("OFFSET %d", offset)
			} else {
				limits = fmt.Sprintf("LIMIT %d OFFSET %d", maxLimit, offset)
			}
		}
	} else if offset <= 0 {
		limits = fmt.Sprintf("LIMIT %d", limit)
	} else {
		limits = fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	return
}

// getIndexSql generate index sql.
func (t *dbTables) getIndexSql(tableName string, useIndex int, indexes []string) (clause string) {
	if len(indexes) == 0 {
		return
	}

	return t.base.GenerateSpecifyIndex(tableName, useIndex, indexes)
}

// crete new tables collection.
func newDbTables(mi *models.ModelInfo, base dbBaser) *dbTables {
	tables := &dbTables{}
	tables.tablesM = make(map[string]*dbTable)
	tables.mi = mi
	tables.base = base
	return tables
}
