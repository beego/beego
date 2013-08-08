package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	format_Date     = "2006-01-02"
	format_DateTime = "2006-01-02 15:04:05"
)

var (
	ErrMissPK = errors.New("missed pk value")
)

var (
	operators = map[string]bool{
		"exact":     true,
		"iexact":    true,
		"contains":  true,
		"icontains": true,
		// "regex":       true,
		// "iregex":      true,
		"gt":          true,
		"gte":         true,
		"lt":          true,
		"lte":         true,
		"startswith":  true,
		"endswith":    true,
		"istartswith": true,
		"iendswith":   true,
		"in":          true,
		// "range":       true,
		// "year":        true,
		// "month":       true,
		// "day":         true,
		// "week_day":    true,
		"isnull": true,
		// "search":      true,
	}
	operatorsSQL = map[string]string{
		"exact":     "= ?",
		"iexact":    "LIKE ?",
		"contains":  "LIKE BINARY ?",
		"icontains": "LIKE ?",
		// "regex":       "REGEXP BINARY ?",
		// "iregex":      "REGEXP ?",
		"gt":          "> ?",
		"gte":         ">= ?",
		"lt":          "< ?",
		"lte":         "<= ?",
		"startswith":  "LIKE BINARY ?",
		"endswith":    "LIKE BINARY ?",
		"istartswith": "LIKE ?",
		"iendswith":   "LIKE ?",
	}
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
			cansel = true
			jtl    *dbTable
		)
		for _, ex := range exs {
			if fi, ok := mmi.fields.GetByAny(ex); ok && fi.rel && fi.fieldType != RelManyToMany {
				names = append(names, fi.name)
				mmi = fi.relModelInfo

				jt := t.set(names, mmi, fi, fi.null == false)
				jt.jtl = jtl

				if fi.reverse {
					cansel = false
				}

				if cansel {
					jt.sel = depth > 0

					if i < relsNum {
						jt.sel = true
					}
				}

				jtl = jt

			} else {
				panic(fmt.Sprintf("unknown model/table name `%s`", ex))
			}
		}
	}
}

func (t *dbTables) getJoinSql() (join string) {
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

		join += fmt.Sprintf("`%s` %s ON %s.`%s` = %s.`%s` ", table, t2,
			t2, c2, t1, c1)
	}
	return
}

func (d *dbTables) parseExprs(mi *modelInfo, exprs []string) (index, column, name string, info *fieldInfo, success bool) {
	var (
		ffi *fieldInfo
		jtl *dbTable
		mmi = mi
	)

	num := len(exprs) - 1
	names := make([]string, 0)

	for i, ex := range exprs {
		exist := false

	check:
		fi, ok := mmi.fields.GetByAny(ex)

		if ok {

			if num != i {
				names = append(names, fi.name)

				switch {
				case fi.rel:
					mmi = fi.relModelInfo
					if fi.fieldType == RelManyToMany {
						mmi = fi.relThroughModelInfo
					}
				case fi.reverse:
					mmi = fi.reverseFieldInfo.mi
					if fi.reverseFieldInfo.fieldType == RelManyToMany {
						mmi = fi.reverseFieldInfo.relThroughModelInfo
					}
				default:
					return
				}

				jt, _ := d.add(names, mmi, fi, fi.null == false)
				jt.jtl = jtl
				jtl = jt

				if fi.rel && fi.fieldType == RelManyToMany {
					ex = fi.relModelInfo.name
					goto check
				}

				if fi.reverse && fi.reverseFieldInfo.fieldType == RelManyToMany {
					ex = fi.reverseFieldInfo.mi.name
					goto check
				}

				exist = true

			} else {

				if ffi == nil {
					index = "T0"
				} else {
					index = jtl.index
				}
				column = fi.column
				info = fi
				if jtl != nil {
					name = jtl.name + ExprSep + fi.name
				} else {
					name = fi.name
				}

				switch fi.fieldType {
				case RelManyToMany, RelReverseMany:
				default:
					exist = true
				}
			}

			ffi = fi
		}

		if exist == false {
			index = ""
			column = ""
			name = ""
			success = false
			return
		}
	}

	success = index != "" && column != ""
	return
}

func (d *dbTables) getCondSql(cond *Condition, sub bool) (where string, params []interface{}) {
	if cond == nil || cond.IsEmpty() {
		return
	}

	mi := d.mi

	// outFor:
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
			w, ps := d.getCondSql(p.cond, true)
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

			index, column, _, _, suc := d.parseExprs(mi, exprs)
			if suc == false {
				panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(p.exprs, ExprSep)))
			}

			if operator == "" {
				operator = "exact"
			}

			operSql, args := d.base.GetOperatorSql(mi, operator, p.args)

			where += fmt.Sprintf("%s.`%s` %s ", index, column, operSql)
			params = append(params, args...)

		}
	}

	if sub == false && where != "" {
		where = "WHERE " + where
	}

	return
}

func (d *dbTables) getOrderSql(orders []string) (orderSql string) {
	if len(orders) == 0 {
		return
	}

	orderSqls := make([]string, 0, len(orders))
	for _, order := range orders {
		asc := "ASC"
		if order[0] == '-' {
			asc = "DESC"
			order = order[1:]
		}
		exprs := strings.Split(order, ExprSep)

		index, column, _, _, suc := d.parseExprs(d.mi, exprs)
		if suc == false {
			panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(exprs, ExprSep)))
		}

		orderSqls = append(orderSqls, fmt.Sprintf("%s.`%s` %s", index, column, asc))
	}

	orderSql = fmt.Sprintf("ORDER BY %s ", strings.Join(orderSqls, ", "))
	return
}

func (d *dbTables) getLimitSql(offset int64, limit int) (limits string) {
	if limit == 0 {
		limit = DefaultRowsLimit
	}
	if limit < 0 {
		// no limit
		if offset > 0 {
			limits = fmt.Sprintf("LIMIT 18446744073709551615 OFFSET %d", offset)
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

type dbBase struct {
	ins dbBaser
}

func (d *dbBase) existPk(mi *modelInfo, ind reflect.Value) (column string, value interface{}, exist bool) {

	fi := mi.fields.pk

	v := ind.Field(fi.fieldIndex)
	if fi.fieldType&IsIntegerField > 0 {
		vu := v.Int()
		exist = vu > 0
		value = vu
	} else {
		vu := v.String()
		exist = vu != ""
		value = vu
	}

	column = fi.column

	return
}

func (d *dbBase) collectValues(mi *modelInfo, ind reflect.Value, skipAuto bool, insert bool) (columns []string, values []interface{}, err error) {
	_, pkValue, _ := d.existPk(mi, ind)
	for _, column := range mi.fields.orders {
		fi := mi.fields.columns[column]
		if fi.dbcol == false || fi.auto && skipAuto {
			continue
		}
		var value interface{}
		if fi.pk {
			value = pkValue
		} else {
			field := ind.Field(fi.fieldIndex)
			if fi.isFielder {
				f := field.Addr().Interface().(Fielder)
				value = f.RawValue()
			} else {
				switch fi.fieldType {
				case TypeBooleanField:
					value = field.Bool()
				case TypeCharField, TypeTextField:
					value = field.String()
				case TypeFloatField, TypeDecimalField:
					value = field.Float()
				case TypeDateField, TypeDateTimeField:
					value = field.Interface()
				default:
					switch {
					case fi.fieldType&IsPostiveIntegerField > 0:
						value = field.Uint()
					case fi.fieldType&IsIntegerField > 0:
						value = field.Int()
					case fi.fieldType&IsRelField > 0:
						if field.IsNil() {
							value = nil
						} else {
							if _, vu, ok := d.existPk(fi.relModelInfo, reflect.Indirect(field)); ok {
								value = vu
							} else {
								value = nil
							}
						}
						if fi.null == false && value == nil {
							return nil, nil, errors.New(fmt.Sprintf("field `%s` cannot be NULL", fi.fullName))
						}
					}
				}
			}
			switch fi.fieldType {
			case TypeDateField, TypeDateTimeField:
				if fi.auto_now || fi.auto_now_add && insert {
					tnow := time.Now()
					if fi.fieldType == TypeDateField {
						value = timeFormat(tnow, format_Date)
					} else {
						value = timeFormat(tnow, format_DateTime)
					}
					if fi.isFielder {
						f := field.Addr().Interface().(Fielder)
						f.SetRaw(tnow)
					} else {
						field.Set(reflect.ValueOf(tnow))
					}
				}
			}
		}
		columns = append(columns, column)
		values = append(values, value)
	}
	return
}

func (d *dbBase) PrepareInsert(q dbQuerier, mi *modelInfo) (*sql.Stmt, error) {
	dbcols := make([]string, 0, len(mi.fields.dbcols))
	marks := make([]string, 0, len(mi.fields.dbcols))
	for _, fi := range mi.fields.fieldsDB {
		if fi.auto == false {
			dbcols = append(dbcols, fi.column)
			marks = append(marks, "?")
		}
	}
	qmarks := strings.Join(marks, ", ")
	columns := strings.Join(dbcols, "`,`")

	query := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES (%s)", mi.table, columns, qmarks)
	return q.Prepare(query)
}

func (d *dbBase) InsertStmt(stmt *sql.Stmt, mi *modelInfo, ind reflect.Value) (int64, error) {
	_, values, err := d.collectValues(mi, ind, true, true)
	if err != nil {
		return 0, err
	}

	if res, err := stmt.Exec(values...); err == nil {
		return res.LastInsertId()
	} else {
		return 0, err
	}
}

func (d *dbBase) Read(q dbQuerier, mi *modelInfo, ind reflect.Value) error {
	pkColumn, pkValue, ok := d.existPk(mi, ind)
	if ok == false {
		return ErrMissPK
	}

	sels := strings.Join(mi.fields.dbcols, "`, `")
	colsNum := len(mi.fields.dbcols)

	query := fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s` = ?", sels, mi.table, pkColumn)

	refs := make([]interface{}, colsNum)
	for i, _ := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	row := q.QueryRow(query, pkValue)
	if err := row.Scan(refs...); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRows
		}
		return err
	} else {
		elm := reflect.New(mi.addrField.Elem().Type())
		md := elm.Interface().(Modeler)
		md.Init(md)
		mind := reflect.Indirect(elm)

		d.setColsValues(mi, &mind, mi.fields.dbcols, refs)

		ind.Set(mind)
	}

	return nil
}

func (d *dbBase) Insert(q dbQuerier, mi *modelInfo, ind reflect.Value) (int64, error) {
	names, values, err := d.collectValues(mi, ind, true, true)
	if err != nil {
		return 0, err
	}

	marks := make([]string, len(names))
	for i, _ := range marks {
		marks[i] = "?"
	}
	qmarks := strings.Join(marks, ", ")
	columns := strings.Join(names, "`,`")

	query := fmt.Sprintf("INSERT INTO `%s` (`%s`) VALUES (%s)", mi.table, columns, qmarks)

	if res, err := q.Exec(query, values...); err == nil {
		return res.LastInsertId()
	} else {
		return 0, err
	}
}

func (d *dbBase) Update(q dbQuerier, mi *modelInfo, ind reflect.Value) (int64, error) {
	pkName, pkValue, ok := d.existPk(mi, ind)
	if ok == false {
		return 0, ErrMissPK
	}
	setNames, setValues, err := d.collectValues(mi, ind, true, false)
	if err != nil {
		return 0, err
	}

	setColumns := strings.Join(setNames, "` = ?, `")

	query := fmt.Sprintf("UPDATE `%s` SET `%s` = ? WHERE `%s` = ?", mi.table, setColumns, pkName)

	setValues = append(setValues, pkValue)

	if res, err := q.Exec(query, setValues...); err == nil {
		return res.RowsAffected()
	} else {
		return 0, err
	}
	return 0, nil
}

func (d *dbBase) Delete(q dbQuerier, mi *modelInfo, ind reflect.Value) (int64, error) {
	pkName, pkValue, ok := d.existPk(mi, ind)
	if ok == false {
		return 0, ErrMissPK
	}

	query := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` = ?", mi.table, pkName)

	if res, err := q.Exec(query, pkValue); err == nil {

		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}

		if num > 0 {
			if mi.fields.pk.auto {
				ind.Field(mi.fields.pk.fieldIndex).SetInt(0)
			}

			err := d.deleteRels(q, mi, []interface{}{pkValue})
			if err != nil {
				return num, err
			}
		}

		return num, err
	} else {
		return 0, err
	}
	return 0, nil
}

func (d *dbBase) UpdateBatch(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, params Params) (int64, error) {
	columns := make([]string, 0, len(params))
	values := make([]interface{}, 0, len(params))
	for col, val := range params {
		if fi, ok := mi.fields.GetByAny(col); ok == false || fi.dbcol == false {
			panic(fmt.Sprintf("wrong field/column name `%s`", col))
		} else {
			columns = append(columns, fi.column)
			values = append(values, val)
		}
	}

	if len(columns) == 0 {
		panic("update params cannot empty")
	}

	tables := newDbTables(mi, d.ins)
	if qs != nil {
		tables.parseRelated(qs.related, qs.relDepth)
	}

	where, args := tables.getCondSql(cond, false)

	join := tables.getJoinSql()

	query := fmt.Sprintf("UPDATE `%s` T0 %sSET T0.`%s` = ? %s", mi.table, join, strings.Join(columns, "` = ?, T0.`"), where)

	values = append(values, args...)

	if res, err := q.Exec(query, values...); err == nil {
		return res.RowsAffected()
	} else {
		return 0, err
	}
	return 0, nil
}

func (d *dbBase) deleteRels(q dbQuerier, mi *modelInfo, args []interface{}) error {
	for _, fi := range mi.fields.fieldsReverse {
		fi = fi.reverseFieldInfo
		switch fi.onDelete {
		case od_CASCADE:
			cond := NewCondition().And(fmt.Sprintf("%s__in", fi.name), args...)
			_, err := d.DeleteBatch(q, nil, fi.mi, cond)
			if err != nil {
				return err
			}
		case od_SET_DEFAULT, od_SET_NULL:
			cond := NewCondition().And(fmt.Sprintf("%s__in", fi.name), args...)
			params := Params{fi.column: nil}
			if fi.onDelete == od_SET_DEFAULT {
				params[fi.column] = fi.initial.String()
			}
			_, err := d.UpdateBatch(q, nil, fi.mi, cond, params)
			if err != nil {
				return err
			}
		case od_DO_NOTHING:
		}
	}
	return nil
}

func (d *dbBase) DeleteBatch(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition) (int64, error) {
	tables := newDbTables(mi, d.ins)
	if qs != nil {
		tables.parseRelated(qs.related, qs.relDepth)
	}

	if cond == nil || cond.IsEmpty() {
		panic("delete operation cannot execute without condition")
	}

	where, args := tables.getCondSql(cond, false)
	join := tables.getJoinSql()

	cols := fmt.Sprintf("T0.`%s`", mi.fields.pk.column)
	query := fmt.Sprintf("SELECT %s FROM `%s` T0 %s%s", cols, mi.table, join, where)

	var rs *sql.Rows
	if r, err := q.Query(query, args...); err != nil {
		return 0, err
	} else {
		rs = r
	}

	var ref interface{}

	args = make([]interface{}, 0)
	cnt := 0
	for rs.Next() {
		if err := rs.Scan(&ref); err != nil {
			return 0, err
		}
		args = append(args, reflect.ValueOf(ref).Interface())
		cnt++
	}

	if cnt == 0 {
		return 0, nil
	}

	sql, args := d.ins.GetOperatorSql(mi, "in", args)
	query = fmt.Sprintf("DELETE FROM `%s` WHERE `%s` %s", mi.table, mi.fields.pk.column, sql)

	if res, err := q.Exec(query, args...); err == nil {
		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}

		if num > 0 {
			err := d.deleteRels(q, mi, args)
			if err != nil {
				return num, err
			}
		}

		return num, nil
	} else {
		return 0, err
	}

	return 0, nil
}

func (d *dbBase) ReadBatch(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, container interface{}) (int64, error) {

	val := reflect.ValueOf(container)
	ind := reflect.Indirect(val)
	typ := ind.Type()

	errTyp := true

	one := true

	if val.Kind() == reflect.Ptr {
		tp := typ
		if ind.Kind() == reflect.Slice {
			one = false
			if ind.Type().Elem().Kind() == reflect.Ptr {
				tp = ind.Type().Elem().Elem()
			}
		}
		errTyp = tp.PkgPath()+"."+tp.Name() != mi.fullName
	}

	if errTyp {
		panic(fmt.Sprintf("wrong object type `%s` for rows scan, need *[]*%s or *%s", val.Type(), mi.fullName, mi.fullName))
	}

	rlimit := qs.limit
	offset := qs.offset
	if one {
		rlimit = 0
		offset = 0
	}

	tables := newDbTables(mi, d.ins)
	tables.parseRelated(qs.related, qs.relDepth)

	where, args := tables.getCondSql(cond, false)
	orderBy := tables.getOrderSql(qs.orders)
	limit := tables.getLimitSql(offset, rlimit)
	join := tables.getJoinSql()

	colsNum := len(mi.fields.dbcols)
	cols := fmt.Sprintf("T0.`%s`", strings.Join(mi.fields.dbcols, "`, T0.`"))
	for _, tbl := range tables.tables {
		if tbl.sel {
			colsNum += len(tbl.mi.fields.dbcols)
			cols += fmt.Sprintf(", %s.`%s`", tbl.index, strings.Join(tbl.mi.fields.dbcols, "`, "+tbl.index+".`"))
		}
	}

	query := fmt.Sprintf("SELECT %s FROM `%s` T0 %s%s%s%s", cols, mi.table, join, where, orderBy, limit)

	var rs *sql.Rows
	if r, err := q.Query(query, args...); err != nil {
		return 0, err
	} else {
		rs = r
	}

	refs := make([]interface{}, colsNum)
	for i, _ := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	slice := ind

	var cnt int64
	for rs.Next() {
		if one && cnt == 0 || one == false {
			if err := rs.Scan(refs...); err != nil {
				return 0, err
			}

			elm := reflect.New(mi.addrField.Elem().Type())
			md := elm.Interface().(Modeler)
			md.Init(md)
			mind := reflect.Indirect(elm)

			cacheV := make(map[string]*reflect.Value)
			cacheM := make(map[string]*modelInfo)
			trefs := refs

			d.setColsValues(mi, &mind, mi.fields.dbcols, refs[:len(mi.fields.dbcols)])
			trefs = refs[len(mi.fields.dbcols):]

			for _, tbl := range tables.tables {
				if tbl.sel {
					last := mind
					names := ""
					mmi := mi
					for _, name := range tbl.names {
						names += name
						if val, ok := cacheV[names]; ok {
							last = *val
							mmi = cacheM[names]
						} else {
							fi := mmi.fields.GetByName(name)
							lastm := mmi
							mmi := fi.relModelInfo
							field := reflect.Indirect(last.Field(fi.fieldIndex))
							if field.IsValid() {
								d.setColsValues(mmi, &field, mmi.fields.dbcols, trefs[:len(mmi.fields.dbcols)])
								for _, fi := range mmi.fields.fieldsReverse {
									if fi.reverseFieldInfo.mi == lastm {
										if fi.reverseFieldInfo != nil {
											field.Field(fi.fieldIndex).Set(last.Addr())
										}
									}
								}
								cacheV[names] = &field
								cacheM[names] = mmi
								last = field
							}
							trefs = trefs[len(mmi.fields.dbcols):]
						}
					}
				}
			}

			if one {
				ind.Set(mind)
			} else {
				slice = reflect.Append(slice, mind.Addr())
			}
		}
		cnt++
	}

	if one == false {
		ind.Set(slice)
	}

	return cnt, nil
}

func (d *dbBase) Count(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition) (cnt int64, err error) {
	tables := newDbTables(mi, d.ins)
	tables.parseRelated(qs.related, qs.relDepth)

	where, args := tables.getCondSql(cond, false)
	tables.getOrderSql(qs.orders)
	join := tables.getJoinSql()

	query := fmt.Sprintf("SELECT COUNT(*) FROM `%s` T0 %s%s", mi.table, join, where)

	row := q.QueryRow(query, args...)

	err = row.Scan(&cnt)
	return
}

func (d *dbBase) GetOperatorSql(mi *modelInfo, operator string, args []interface{}) (string, []interface{}) {
	params := make([]interface{}, len(args))
	copy(params, args)
	sql := ""
	for i, arg := range args {
		if md, ok := arg.(Modeler); ok {
			ind := reflect.Indirect(reflect.ValueOf(md))
			if _, vu, exist := d.existPk(mi, ind); exist {
				arg = vu
			} else {
				panic(fmt.Sprintf("`%s` need a valid args value", operator))
			}
		}
		params[i] = arg
	}
	if operator == "in" {
		marks := make([]string, len(params))
		for i, _ := range marks {
			marks[i] = "?"
		}
		sql = fmt.Sprintf("IN (%s)", strings.Join(marks, ", "))
	} else {
		if len(params) > 1 {
			panic(fmt.Sprintf("operator `%s` need 1 args not %d", operator, len(params)))
		}
		sql = operatorsSQL[operator]
		arg := params[0]
		switch operator {
		case "exact":
			if arg == nil {
				params[0] = "IS NULL"
			}
		case "iexact", "contains", "icontains", "startswith", "endswith", "istartswith", "iendswith":
			param := strings.Replace(ToStr(arg), `%`, `\%`, -1)
			switch operator {
			case "iexact":
			case "contains", "icontains":
				param = fmt.Sprintf("%%%s%%", param)
			case "startswith", "istartswith":
				param = fmt.Sprintf("%s%%", param)
			case "endswith", "iendswith":
				param = fmt.Sprintf("%%%s", param)
			}
			params[0] = param
		case "isnull":
			if b, ok := arg.(bool); ok {
				if b {
					sql = "IS NULL"
				} else {
					sql = "IS NOT NULL"
				}
				params = nil
			} else {
				panic(fmt.Sprintf("operator `%s` need a bool value not `%T`", operator, arg))
			}
		}
	}
	return sql, params
}

func (d *dbBase) setColsValues(mi *modelInfo, ind *reflect.Value, cols []string, values []interface{}) {
	for i, column := range cols {
		val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()

		fi := mi.fields.GetByColumn(column)

		field := ind.Field(fi.fieldIndex)

		value, err := d.getValue(fi, val)
		if err != nil {
			panic(fmt.Sprintf("db value convert failed `%v` %s", val, err.Error()))
		}

		_, err = d.setValue(fi, value, &field)

		if err != nil {
			panic(fmt.Sprintf("db value convert failed `%v` %s", val, err.Error()))
		}
	}
}

func (d *dbBase) getValue(fi *fieldInfo, val interface{}) (interface{}, error) {
	if val == nil {
		return nil, nil
	}

	var value interface{}

	var str *StrTo
	switch v := val.(type) {
	case []byte:
		s := StrTo(string(v))
		str = &s
	case string:
		s := StrTo(v)
		str = &s
	}

	fieldType := fi.fieldType

setValue:
	switch {
	case fieldType == TypeBooleanField:
		if str == nil {
			switch v := val.(type) {
			case int64:
				b := v == 1
				value = b
			default:
				s := StrTo(ToStr(v))
				str = &s
			}
		}
		if str != nil {
			b, err := str.Bool()
			if err != nil {
				return nil, err
			}
			value = b
		}
	case fieldType == TypeCharField || fieldType == TypeTextField:
		if str == nil {
			value = ToStr(val)
		} else {
			value = str.String()
		}
	case fieldType == TypeDateField || fieldType == TypeDateTimeField:
		if str == nil {
			switch v := val.(type) {
			case time.Time:
				value = v
			default:
				s := StrTo(ToStr(v))
				str = &s
			}
		}
		if str != nil {
			format := format_DateTime
			if fi.fieldType == TypeDateField {
				format = format_Date
			}
			s := str.String()
			t, err := timeParse(s, format)
			if err != nil && s != "0000-00-00" && s != "0000-00-00 00:00:00" {
				return nil, err
			}
			value = t
		}
	case fieldType&IsIntegerField > 0:
		if str == nil {
			s := StrTo(ToStr(val))
			str = &s
		}
		if str != nil {
			var err error
			switch fieldType {
			case TypeSmallIntegerField:
				_, err = str.Int16()
			case TypeIntegerField:
				_, err = str.Int32()
			case TypeBigIntegerField:
				_, err = str.Int64()
			case TypePositiveSmallIntegerField:
				_, err = str.Uint16()
			case TypePositiveIntegerField:
				_, err = str.Uint32()
			case TypePositiveBigIntegerField:
				_, err = str.Uint64()
			}
			if err != nil {
				return nil, err
			}
			if fieldType&IsPostiveIntegerField > 0 {
				v, _ := str.Uint64()
				value = v
			} else {
				v, _ := str.Int64()
				value = v
			}
		}
	case fieldType == TypeFloatField || fieldType == TypeDecimalField:
		if str == nil {
			switch v := val.(type) {
			case float64:
				value = v
			default:
				s := StrTo(ToStr(v))
				str = &s
			}
		}
		if str != nil {
			v, err := str.Float64()
			if err != nil {
				return nil, err
			}
			value = v
		}
	case fieldType&IsRelField > 0:
		fieldType = fi.relModelInfo.fields.pk.fieldType
		goto setValue
	}

	return value, nil

}

func (d *dbBase) setValue(fi *fieldInfo, value interface{}, field *reflect.Value) (interface{}, error) {

	fieldType := fi.fieldType
	isNative := fi.isFielder == false

setValue:
	switch {
	case fieldType == TypeBooleanField:
		if isNative {
			if value == nil {
				value = false
			}
			field.SetBool(value.(bool))
		}
	case fieldType == TypeCharField || fieldType == TypeTextField:
		if isNative {
			if value == nil {
				value = ""
			}
			field.SetString(value.(string))
		}
	case fieldType == TypeDateField || fieldType == TypeDateTimeField:
		if isNative {
			if value == nil {
				value = time.Time{}
			}
			field.Set(reflect.ValueOf(value))
		}
	case fieldType&IsIntegerField > 0:
		if fieldType&IsPostiveIntegerField > 0 {
			if isNative {
				if value == nil {
					value = uint64(0)
				}
				field.SetUint(value.(uint64))
			}
		} else {
			if isNative {
				if value == nil {
					value = int64(0)
				}
				field.SetInt(value.(int64))
			}
		}
	case fieldType == TypeFloatField || fieldType == TypeDecimalField:
		if isNative {
			if value == nil {
				value = float64(0)
			}
			field.SetFloat(value.(float64))
		}
	case fieldType&IsRelField > 0:
		if value != nil {
			fieldType = fi.relModelInfo.fields.pk.fieldType
			mf := reflect.New(fi.relModelInfo.addrField.Elem().Type())
			md := mf.Interface().(Modeler)
			md.Init(md)
			field.Set(mf)
			f := mf.Elem().Field(fi.relModelInfo.fields.pk.fieldIndex)
			field = &f
			goto setValue
		}
	}

	if isNative == false {
		fd := field.Addr().Interface().(Fielder)
		err := fd.SetRaw(value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (d *dbBase) ReadValues(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, exprs []string, container interface{}) (int64, error) {

	var (
		maps  []Params
		lists []ParamsList
		list  ParamsList
	)

	typ := 0
	switch container.(type) {
	case *[]Params:
		typ = 1
	case *[]ParamsList:
		typ = 2
	case *ParamsList:
		typ = 3
	default:
		panic(fmt.Sprintf("unsupport read values type `%T`", container))
	}

	tables := newDbTables(mi, d.ins)

	var (
		cols  []string
		infos []*fieldInfo
	)

	hasExprs := len(exprs) > 0

	if hasExprs {
		cols = make([]string, 0, len(exprs))
		infos = make([]*fieldInfo, 0, len(exprs))
		for _, ex := range exprs {
			index, col, name, fi, suc := tables.parseExprs(mi, strings.Split(ex, ExprSep))
			if suc == false {
				panic(fmt.Errorf("unknown field/column name `%s`", ex))
			}
			cols = append(cols, fmt.Sprintf("%s.`%s` `%s`", index, col, name))
			infos = append(infos, fi)
		}
	} else {
		cols = make([]string, 0, len(mi.fields.dbcols))
		infos = make([]*fieldInfo, 0, len(exprs))
		for _, fi := range mi.fields.fieldsDB {
			cols = append(cols, fmt.Sprintf("T0.`%s` `%s`", fi.column, fi.name))
			infos = append(infos, fi)
		}
	}

	where, args := tables.getCondSql(cond, false)
	orderBy := tables.getOrderSql(qs.orders)
	limit := tables.getLimitSql(qs.offset, qs.limit)
	join := tables.getJoinSql()

	sels := strings.Join(cols, ", ")

	query := fmt.Sprintf("SELECT %s FROM `%s` T0 %s%s%s%s", sels, mi.table, join, where, orderBy, limit)

	var rs *sql.Rows
	if r, err := q.Query(query, args...); err != nil {
		return 0, err
	} else {
		rs = r
	}

	refs := make([]interface{}, len(cols))
	for i, _ := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	var (
		cnt     int64
		columns []string
	)
	for rs.Next() {
		if cnt == 0 {
			if cols, err := rs.Columns(); err != nil {
				return 0, err
			} else {
				columns = cols
			}
		}

		if err := rs.Scan(refs...); err != nil {
			return 0, err
		}

		switch typ {
		case 1:
			params := make(Params, len(cols))
			for i, ref := range refs {
				fi := infos[i]

				val := reflect.Indirect(reflect.ValueOf(ref)).Interface()

				value, err := d.getValue(fi, val)
				if err != nil {
					panic(fmt.Sprintf("db value convert failed `%v` %s", val, err.Error()))
				}

				params[columns[i]] = value
			}
			maps = append(maps, params)
		case 2:
			params := make(ParamsList, 0, len(cols))
			for i, ref := range refs {
				fi := infos[i]

				val := reflect.Indirect(reflect.ValueOf(ref)).Interface()

				value, err := d.getValue(fi, val)
				if err != nil {
					panic(fmt.Sprintf("db value convert failed `%v` %s", val, err.Error()))
				}

				params = append(params, value)
			}
			lists = append(lists, params)
		case 3:
			for i, ref := range refs {
				fi := infos[i]

				val := reflect.Indirect(reflect.ValueOf(ref)).Interface()

				value, err := d.getValue(fi, val)
				if err != nil {
					panic(fmt.Sprintf("db value convert failed `%v` %s", val, err.Error()))
				}

				list = append(list, value)
			}
		}

		cnt++
	}

	switch v := container.(type) {
	case *[]Params:
		*v = maps
	case *[]ParamsList:
		*v = lists
	case *ParamsList:
		*v = list
	}

	return cnt, nil
}
