package orm

import (
	"fmt"
	"strings"
	"time"
)

type dbTable struct {
	id    int
	index string
	name  string
	names []string
	sel   bool
	inner bool
	mi    *modelInfo
	fi    *fieldInfo
	jtl   *dbTable
}

type dbTables struct {
	tablesM map[string]*dbTable
	tables  []*dbTable
	mi      *modelInfo
	base    dbBaser
}

func (t *dbTables) set(names []string, mi *modelInfo, fi *fieldInfo, inner bool) *dbTable {
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

func (t *dbTables) add(names []string, mi *modelInfo, fi *fieldInfo, inner bool) (*dbTable, bool) {
	name := strings.Join(names, ExprSep)
	if _, ok := t.tablesM[name]; ok == false {
		i := len(t.tables) + 1
		jt := &dbTable{i, fmt.Sprintf("T%d", i), name, names, false, inner, mi, fi, nil}
		t.tablesM[name] = jt
		t.tables = append(t.tables, jt)
		return jt, true
	}
	return t.tablesM[name], false
}

func (t *dbTables) get(name string) (*dbTable, bool) {
	j, ok := t.tablesM[name]
	return j, ok
}

func (t *dbTables) loopDepth(depth int, prefix string, fi *fieldInfo, related []string) []string {
	if depth < 0 || fi.fieldType == RelManyToMany {
		return related
	}

	if prefix == "" {
		prefix = fi.name
	} else {
		prefix = prefix + ExprSep + fi.name
	}
	related = append(related, prefix)

	depth--
	for _, fi := range fi.relModelInfo.fields.fieldsRel {
		related = t.loopDepth(depth, prefix, fi, related)
	}

	return related
}

func (t *dbTables) parseRelated(rels []string, depth int) {

	relsNum := len(rels)
	related := make([]string, relsNum)
	copy(related, rels)

	relDepth := depth

	if relsNum != 0 {
		relDepth = 0
	}

	relDepth--
	for _, fi := range t.mi.fields.fieldsRel {
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
			if fi, ok := mmi.fields.GetByAny(ex); ok && fi.rel && fi.fieldType != RelManyToMany {
				names = append(names, fi.name)
				mmi = fi.relModelInfo

				if fi.null {
					inner = false
				}

				jt := t.set(names, mmi, fi, inner)
				jt.jtl = jtl

				if fi.reverse {
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

func (t *dbTables) getJoinSql() (join string) {
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
		table = jt.mi.table

		switch {
		case jt.fi.fieldType == RelManyToMany || jt.fi.reverse && jt.fi.reverseFieldInfo.fieldType == RelManyToMany:
			c1 = jt.fi.mi.fields.pk.column
			for _, ffi := range jt.mi.fields.fieldsRel {
				if jt.fi.mi == ffi.relModelInfo {
					c2 = ffi.column
					break
				}
			}
		default:
			c1 = jt.fi.column
			c2 = jt.fi.relModelInfo.fields.pk.column

			if jt.fi.reverse {
				c1 = jt.mi.fields.pk.column
				c2 = jt.fi.reverseFieldInfo.column
			}
		}

		join += fmt.Sprintf("%s%s%s %s ON %s.%s%s%s = %s.%s%s%s ", Q, table, Q, t2,
			t2, Q, c2, Q, t1, Q, c1, Q)
	}
	return
}

func (t *dbTables) parseExprs(mi *modelInfo, exprs []string) (index, name string, info *fieldInfo, success bool) {
	var (
		jtl *dbTable
		mmi = mi
	)

	num := len(exprs) - 1
	names := make([]string, 0)

	inner := true

	for i, ex := range exprs {

		fi, ok := mmi.fields.GetByAny(ex)

		if ok {

			isRel := fi.rel || fi.reverse

			names = append(names, fi.name)

			switch {
			case fi.rel:
				mmi = fi.relModelInfo
				if fi.fieldType == RelManyToMany {
					mmi = fi.relThroughModelInfo
				}
			case fi.reverse:
				mmi = fi.reverseFieldInfo.mi
			}

			if isRel && (fi.mi.isThrough == false || num != i) {
				if fi.null {
					inner = false
				}

				jt, _ := t.add(names, mmi, fi, inner)
				jt.jtl = jtl
				jtl = jt
			}

			if num == i {
				if i == 0 || jtl == nil {
					index = "T0"
				} else {
					index = jtl.index
				}

				info = fi

				if jtl == nil {
					name = fi.name
				} else {
					name = jtl.name + ExprSep + fi.name
				}

				switch {
				case fi.rel:

				case fi.reverse:
					switch fi.reverseFieldInfo.fieldType {
					case RelOneToOne, RelForeignKey:
						index = jtl.index
						info = fi.reverseFieldInfo.mi.fields.pk
						name = info.name
					}
				}
			}

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

func (t *dbTables) getCondSql(cond *Condition, sub bool, tz *time.Location) (where string, params []interface{}) {
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
			w, ps := t.getCondSql(p.cond, true, tz)
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
			if suc == false {
				panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(p.exprs, ExprSep)))
			}

			if operator == "" {
				operator = "exact"
			}

			operSql, args := t.base.GenerateOperatorSql(mi, fi, operator, p.args, tz)

			leftCol := fmt.Sprintf("%s.%s%s%s", index, Q, fi.column, Q)
			t.base.GenerateOperatorLeftCol(fi, operator, &leftCol)

			where += fmt.Sprintf("%s %s ", leftCol, operSql)
			params = append(params, args...)

		}
	}

	if sub == false && where != "" {
		where = "WHERE " + where
	}

	return
}

func (t *dbTables) getOrderSql(orders []string) (orderSql string) {
	if len(orders) == 0 {
		return
	}

	Q := t.base.TableQuote()

	orderSqls := make([]string, 0, len(orders))
	for _, order := range orders {
		asc := "ASC"
		if order[0] == '-' {
			asc = "DESC"
			order = order[1:]
		}
		exprs := strings.Split(order, ExprSep)

		index, _, fi, suc := t.parseExprs(t.mi, exprs)
		if suc == false {
			panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(exprs, ExprSep)))
		}

		orderSqls = append(orderSqls, fmt.Sprintf("%s.%s%s%s %s", index, Q, fi.column, Q, asc))
	}

	orderSql = fmt.Sprintf("ORDER BY %s ", strings.Join(orderSqls, ", "))
	return
}

func (t *dbTables) getLimitSql(mi *modelInfo, offset int64, limit int64) (limits string) {
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

func newDbTables(mi *modelInfo, base dbBaser) *dbTables {
	tables := &dbTables{}
	tables.tablesM = make(map[string]*dbTable)
	tables.mi = mi
	tables.base = base
	return tables
}
