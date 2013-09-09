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
)

type dbBase struct {
	ins dbBaser
}

func (d *dbBase) collectValues(mi *modelInfo, ind reflect.Value, skipAuto bool, insert bool, tz *time.Location) (columns []string, values []interface{}, err error) {
	_, pkValue, _ := getExistPk(mi, ind)
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
					vu := field.Interface()
					if _, ok := vu.(float32); ok {
						value, _ = StrTo(ToStr(vu)).Float64()
					} else {
						value = field.Float()
					}
				case TypeDateField, TypeDateTimeField:
					value = field.Interface()
					if t, ok := value.(time.Time); ok {
						if fi.fieldType == TypeDateField {
							d.ins.TimeToDB(&t, DefaultTimeLoc)
						} else {
							d.ins.TimeToDB(&t, tz)
						}
						value = t
					}
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
							if _, vu, ok := getExistPk(fi.relModelInfo, reflect.Indirect(field)); ok {
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
						d.ins.TimeToDB(&tnow, DefaultTimeLoc)
					} else {
						d.ins.TimeToDB(&tnow, tz)
					}
					value = tnow
					if fi.isFielder {
						f := field.Addr().Interface().(Fielder)
						f.SetRaw(tnow.In(DefaultTimeLoc))
					} else {
						field.Set(reflect.ValueOf(tnow.In(DefaultTimeLoc)))
					}
				}
			}
		}
		columns = append(columns, column)
		values = append(values, value)
	}
	return
}

func (d *dbBase) PrepareInsert(q dbQuerier, mi *modelInfo) (stmtQuerier, string, error) {
	Q := d.ins.TableQuote()

	dbcols := make([]string, 0, len(mi.fields.dbcols))
	marks := make([]string, 0, len(mi.fields.dbcols))
	for _, fi := range mi.fields.fieldsDB {
		if fi.auto == false {
			dbcols = append(dbcols, fi.column)
			marks = append(marks, "?")
		}
	}
	qmarks := strings.Join(marks, ", ")
	sep := fmt.Sprintf("%s, %s", Q, Q)
	columns := strings.Join(dbcols, sep)

	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, mi.table, Q, Q, columns, Q, qmarks)

	d.ins.ReplaceMarks(&query)

	d.ins.HasReturningID(mi, &query)

	stmt, err := q.Prepare(query)
	return stmt, query, err
}

func (d *dbBase) InsertStmt(stmt stmtQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	_, values, err := d.collectValues(mi, ind, true, true, tz)
	if err != nil {
		return 0, err
	}

	if d.ins.HasReturningID(mi, nil) {
		row := stmt.QueryRow(values...)
		var id int64
		err := row.Scan(&id)
		return id, err
	} else {
		if res, err := stmt.Exec(values...); err == nil {
			return res.LastInsertId()
		} else {
			return 0, err
		}
	}
}

func (d *dbBase) Read(q dbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) error {
	pkColumn, pkValue, ok := getExistPk(mi, ind)
	if ok == false {
		return ErrMissPK
	}

	Q := d.ins.TableQuote()

	sep := fmt.Sprintf("%s, %s", Q, Q)
	sels := strings.Join(mi.fields.dbcols, sep)
	colsNum := len(mi.fields.dbcols)

	query := fmt.Sprintf("SELECT %s%s%s FROM %s%s%s WHERE %s%s%s = ?", Q, sels, Q, Q, mi.table, Q, Q, pkColumn, Q)

	refs := make([]interface{}, colsNum)
	for i, _ := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	d.ins.ReplaceMarks(&query)

	row := q.QueryRow(query, pkValue)
	if err := row.Scan(refs...); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRows
		}
		return err
	} else {
		elm := reflect.New(mi.addrField.Elem().Type())
		mind := reflect.Indirect(elm)

		d.setColsValues(mi, &mind, mi.fields.dbcols, refs, tz)

		ind.Set(mind)
	}

	return nil
}

func (d *dbBase) Insert(q dbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	names, values, err := d.collectValues(mi, ind, true, true, tz)
	if err != nil {
		return 0, err
	}

	Q := d.ins.TableQuote()

	marks := make([]string, len(names))
	for i, _ := range marks {
		marks[i] = "?"
	}

	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	columns := strings.Join(names, sep)

	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, mi.table, Q, Q, columns, Q, qmarks)

	d.ins.ReplaceMarks(&query)

	if d.ins.HasReturningID(mi, &query) {
		row := q.QueryRow(query, values...)
		var id int64
		err := row.Scan(&id)
		return id, err
	} else {
		if res, err := q.Exec(query, values...); err == nil {
			return res.LastInsertId()
		} else {
			return 0, err
		}
	}
}

func (d *dbBase) Update(q dbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	pkName, pkValue, ok := getExistPk(mi, ind)
	if ok == false {
		return 0, ErrMissPK
	}
	setNames, setValues, err := d.collectValues(mi, ind, true, false, tz)
	if err != nil {
		return 0, err
	}

	setValues = append(setValues, pkValue)

	Q := d.ins.TableQuote()

	sep := fmt.Sprintf("%s = ?, %s", Q, Q)
	setColumns := strings.Join(setNames, sep)

	query := fmt.Sprintf("UPDATE %s%s%s SET %s%s%s = ? WHERE %s%s%s = ?", Q, mi.table, Q, Q, setColumns, Q, Q, pkName, Q)

	d.ins.ReplaceMarks(&query)

	if res, err := q.Exec(query, setValues...); err == nil {
		return res.RowsAffected()
	} else {
		return 0, err
	}
	return 0, nil
}

func (d *dbBase) Delete(q dbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	pkName, pkValue, ok := getExistPk(mi, ind)
	if ok == false {
		return 0, ErrMissPK
	}

	Q := d.ins.TableQuote()

	query := fmt.Sprintf("DELETE FROM %s%s%s WHERE %s%s%s = ?", Q, mi.table, Q, Q, pkName, Q)

	d.ins.ReplaceMarks(&query)

	if res, err := q.Exec(query, pkValue); err == nil {

		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}

		if num > 0 {
			if mi.fields.pk.auto {
				if mi.fields.pk.fieldType&IsPostiveIntegerField > 0 {
					ind.Field(mi.fields.pk.fieldIndex).SetUint(0)
				} else {
					ind.Field(mi.fields.pk.fieldIndex).SetInt(0)
				}
			}

			err := d.deleteRels(q, mi, []interface{}{pkValue}, tz)
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

func (d *dbBase) UpdateBatch(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, params Params, tz *time.Location) (int64, error) {
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

	where, args := tables.getCondSql(cond, false, tz)

	values = append(values, args...)

	join := tables.getJoinSql()

	var query string

	Q := d.ins.TableQuote()

	if d.ins.SupportUpdateJoin() {
		cols := strings.Join(columns, fmt.Sprintf("%s = ?, T0.%s", Q, Q))
		query = fmt.Sprintf("UPDATE %s%s%s T0 %sSET T0.%s%s%s = ? %s", Q, mi.table, Q, join, Q, cols, Q, where)
	} else {
		cols := strings.Join(columns, fmt.Sprintf("%s = ?, %s", Q, Q))
		supQuery := fmt.Sprintf("SELECT T0.%s%s%s FROM %s%s%s T0 %s%s", Q, mi.fields.pk.column, Q, Q, mi.table, Q, join, where)
		query = fmt.Sprintf("UPDATE %s%s%s SET %s%s%s = ? WHERE %s%s%s IN ( %s )", Q, mi.table, Q, Q, cols, Q, Q, mi.fields.pk.column, Q, supQuery)
	}

	d.ins.ReplaceMarks(&query)

	if res, err := q.Exec(query, values...); err == nil {
		return res.RowsAffected()
	} else {
		return 0, err
	}
	return 0, nil
}

func (d *dbBase) deleteRels(q dbQuerier, mi *modelInfo, args []interface{}, tz *time.Location) error {
	for _, fi := range mi.fields.fieldsReverse {
		fi = fi.reverseFieldInfo
		switch fi.onDelete {
		case od_CASCADE:
			cond := NewCondition().And(fmt.Sprintf("%s__in", fi.name), args...)
			_, err := d.DeleteBatch(q, nil, fi.mi, cond, tz)
			if err != nil {
				return err
			}
		case od_SET_DEFAULT, od_SET_NULL:
			cond := NewCondition().And(fmt.Sprintf("%s__in", fi.name), args...)
			params := Params{fi.column: nil}
			if fi.onDelete == od_SET_DEFAULT {
				params[fi.column] = fi.initial.String()
			}
			_, err := d.UpdateBatch(q, nil, fi.mi, cond, params, tz)
			if err != nil {
				return err
			}
		case od_DO_NOTHING:
		}
	}
	return nil
}

func (d *dbBase) DeleteBatch(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, tz *time.Location) (int64, error) {
	tables := newDbTables(mi, d.ins)
	if qs != nil {
		tables.parseRelated(qs.related, qs.relDepth)
	}

	if cond == nil || cond.IsEmpty() {
		panic("delete operation cannot execute without condition")
	}

	Q := d.ins.TableQuote()

	where, args := tables.getCondSql(cond, false, tz)
	join := tables.getJoinSql()

	cols := fmt.Sprintf("T0.%s%s%s", Q, mi.fields.pk.column, Q)
	query := fmt.Sprintf("SELECT %s FROM %s%s%s T0 %s%s", cols, Q, mi.table, Q, join, where)

	d.ins.ReplaceMarks(&query)

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

	marks := make([]string, len(args))
	for i, _ := range marks {
		marks[i] = "?"
	}
	sql := fmt.Sprintf("IN (%s)", strings.Join(marks, ", "))
	query = fmt.Sprintf("DELETE FROM %s%s%s WHERE %s%s%s %s", Q, mi.table, Q, Q, mi.fields.pk.column, Q, sql)

	d.ins.ReplaceMarks(&query)

	if res, err := q.Exec(query, args...); err == nil {
		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}

		if num > 0 {
			err := d.deleteRels(q, mi, args, tz)
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

func (d *dbBase) ReadBatch(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, container interface{}, tz *time.Location) (int64, error) {

	val := reflect.ValueOf(container)
	ind := reflect.Indirect(val)

	errTyp := true
	one := true
	isPtr := true

	if val.Kind() == reflect.Ptr {
		fn := ""
		if ind.Kind() == reflect.Slice {
			one = false
			typ := ind.Type().Elem()
			switch typ.Kind() {
			case reflect.Ptr:
				fn = getFullName(typ.Elem())
			case reflect.Struct:
				isPtr = false
				fn = getFullName(typ)
			}
		} else {
			fn = getFullName(ind.Type())
		}
		errTyp = fn != mi.fullName
	}

	if errTyp {
		if one {
			panic(fmt.Sprintf("wrong object type `%s` for rows scan, need *%s", val.Type(), mi.fullName))
		} else {
			panic(fmt.Sprintf("wrong object type `%s` for rows scan, need *[]*%s or *[]%s", val.Type(), mi.fullName, mi.fullName))
		}
	}

	rlimit := qs.limit
	offset := qs.offset

	Q := d.ins.TableQuote()

	tables := newDbTables(mi, d.ins)
	tables.parseRelated(qs.related, qs.relDepth)

	where, args := tables.getCondSql(cond, false, tz)
	orderBy := tables.getOrderSql(qs.orders)
	limit := tables.getLimitSql(mi, offset, rlimit)
	join := tables.getJoinSql()

	colsNum := len(mi.fields.dbcols)
	sep := fmt.Sprintf("%s, T0.%s", Q, Q)
	cols := fmt.Sprintf("T0.%s%s%s", Q, strings.Join(mi.fields.dbcols, sep), Q)
	for _, tbl := range tables.tables {
		if tbl.sel {
			colsNum += len(tbl.mi.fields.dbcols)
			sep := fmt.Sprintf("%s, %s.%s", Q, tbl.index, Q)
			cols += fmt.Sprintf(", %s.%s%s%s", tbl.index, Q, strings.Join(tbl.mi.fields.dbcols, sep), Q)
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s%s%s T0 %s%s%s%s", cols, Q, mi.table, Q, join, where, orderBy, limit)

	d.ins.ReplaceMarks(&query)

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
			mind := reflect.Indirect(elm)

			cacheV := make(map[string]*reflect.Value)
			cacheM := make(map[string]*modelInfo)
			trefs := refs

			d.setColsValues(mi, &mind, mi.fields.dbcols, refs[:len(mi.fields.dbcols)], tz)
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
								d.setColsValues(mmi, &field, mmi.fields.dbcols, trefs[:len(mmi.fields.dbcols)], tz)
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
				if cnt == 0 {
					slice = reflect.New(ind.Type()).Elem()
				}

				if isPtr {
					slice = reflect.Append(slice, mind.Addr())
				} else {
					slice = reflect.Append(slice, mind)
				}
			}
		}
		cnt++
	}

	if one == false && cnt > 0 {
		ind.Set(slice)
	}

	return cnt, nil
}

func (d *dbBase) Count(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, tz *time.Location) (cnt int64, err error) {
	tables := newDbTables(mi, d.ins)
	tables.parseRelated(qs.related, qs.relDepth)

	where, args := tables.getCondSql(cond, false, tz)
	tables.getOrderSql(qs.orders)
	join := tables.getJoinSql()

	Q := d.ins.TableQuote()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s%s%s T0 %s%s", Q, mi.table, Q, join, where)

	d.ins.ReplaceMarks(&query)

	row := q.QueryRow(query, args...)

	err = row.Scan(&cnt)
	return
}

func (d *dbBase) GenerateOperatorSql(mi *modelInfo, fi *fieldInfo, operator string, args []interface{}, tz *time.Location) (string, []interface{}) {
	sql := ""
	params := getFlatParams(fi, args, tz)

	if len(params) == 0 {
		panic(fmt.Sprintf("operator `%s` need at least one args", operator))
	}
	arg := params[0]

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
		sql = d.ins.OperatorSql(operator)
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

func (d *dbBase) GenerateOperatorLeftCol(*fieldInfo, string, *string) {
	// default not use
}

func (d *dbBase) setColsValues(mi *modelInfo, ind *reflect.Value, cols []string, values []interface{}, tz *time.Location) {
	for i, column := range cols {
		val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()

		fi := mi.fields.GetByColumn(column)

		field := ind.Field(fi.fieldIndex)

		value, err := d.convertValueFromDB(fi, val, tz)
		if err != nil {
			panic(fmt.Sprintf("Raw value: `%v` %s", val, err.Error()))
		}

		_, err = d.setFieldValue(fi, value, field)

		if err != nil {
			panic(fmt.Sprintf("Raw value: `%v` %s", val, err.Error()))
		}
	}
}

func (d *dbBase) convertValueFromDB(fi *fieldInfo, val interface{}, tz *time.Location) (interface{}, error) {
	if val == nil {
		return nil, nil
	}

	var value interface{}
	var tErr error

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
				tErr = err
				goto end
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
			switch t := val.(type) {
			case time.Time:
				d.ins.TimeFromDB(&t, tz)
				value = t
			default:
				s := StrTo(ToStr(t))
				str = &s
			}
		}
		if str != nil {
			s := str.String()
			var (
				t   time.Time
				err error
			)
			if fi.fieldType == TypeDateField {
				if len(s) > 10 {
					s = s[:10]
				}
				t, err = time.ParseInLocation(format_Date, s, DefaultTimeLoc)
			} else {
				if len(s) > 19 {
					s = s[:19]
				}
				t, err = time.ParseInLocation(format_DateTime, s, tz)
				t = t.In(DefaultTimeLoc)
			}
			if err != nil && s != "0000-00-00" && s != "0000-00-00 00:00:00" {
				tErr = err
				goto end
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
			case TypeBitField:
				_, err = str.Int8()
			case TypeSmallIntegerField:
				_, err = str.Int16()
			case TypeIntegerField:
				_, err = str.Int32()
			case TypeBigIntegerField:
				_, err = str.Int64()
			case TypePositiveBitField:
				_, err = str.Uint8()
			case TypePositiveSmallIntegerField:
				_, err = str.Uint16()
			case TypePositiveIntegerField:
				_, err = str.Uint32()
			case TypePositiveBigIntegerField:
				_, err = str.Uint64()
			}
			if err != nil {
				tErr = err
				goto end
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
				tErr = err
				goto end
			}
			value = v
		}
	case fieldType&IsRelField > 0:
		fi = fi.relModelInfo.fields.pk
		fieldType = fi.fieldType
		goto setValue
	}

end:
	if tErr != nil {
		err := fmt.Errorf("convert to `%s` failed, field: %s err: %s", fi.addrValue.Type(), fi.fullName, tErr)
		return nil, err
	}

	return value, nil

}

func (d *dbBase) setFieldValue(fi *fieldInfo, value interface{}, field reflect.Value) (interface{}, error) {

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
			field.Set(mf)
			f := mf.Elem().Field(fi.relModelInfo.fields.pk.fieldIndex)
			field = f
			goto setValue
		}
	}

	if isNative == false {
		fd := field.Addr().Interface().(Fielder)
		err := fd.SetRaw(value)
		if err != nil {
			err = fmt.Errorf("converted value `%v` set to Fielder `%s` failed, err: %s", value, fi.fullName, err)
			return nil, err
		}
	}

	return value, nil
}

func (d *dbBase) ReadValues(q dbQuerier, qs *querySet, mi *modelInfo, cond *Condition, exprs []string, container interface{}, tz *time.Location) (int64, error) {

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

	Q := d.ins.TableQuote()

	if hasExprs {
		cols = make([]string, 0, len(exprs))
		infos = make([]*fieldInfo, 0, len(exprs))
		for _, ex := range exprs {
			index, name, fi, suc := tables.parseExprs(mi, strings.Split(ex, ExprSep))
			if suc == false {
				panic(fmt.Errorf("unknown field/column name `%s`", ex))
			}
			cols = append(cols, fmt.Sprintf("%s.%s%s%s %s%s%s", index, Q, fi.column, Q, Q, name, Q))
			infos = append(infos, fi)
		}
	} else {
		cols = make([]string, 0, len(mi.fields.dbcols))
		infos = make([]*fieldInfo, 0, len(exprs))
		for _, fi := range mi.fields.fieldsDB {
			cols = append(cols, fmt.Sprintf("T0.%s%s%s %s%s%s", Q, fi.column, Q, Q, fi.name, Q))
			infos = append(infos, fi)
		}
	}

	where, args := tables.getCondSql(cond, false, tz)
	orderBy := tables.getOrderSql(qs.orders)
	limit := tables.getLimitSql(mi, qs.offset, qs.limit)
	join := tables.getJoinSql()

	sels := strings.Join(cols, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s%s%s T0 %s%s%s%s", sels, Q, mi.table, Q, join, where, orderBy, limit)

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

				value, err := d.convertValueFromDB(fi, val, tz)
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

				value, err := d.convertValueFromDB(fi, val, tz)
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

				value, err := d.convertValueFromDB(fi, val, tz)
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

func (d *dbBase) SupportUpdateJoin() bool {
	return true
}

func (d *dbBase) MaxLimit() uint64 {
	return 18446744073709551615
}

func (d *dbBase) TableQuote() string {
	return "`"
}

func (d *dbBase) ReplaceMarks(query *string) {
	// default use `?` as mark, do nothing
}

func (d *dbBase) HasReturningID(*modelInfo, *string) bool {
	return false
}

func (d *dbBase) TimeFromDB(t *time.Time, tz *time.Location) {
	*t = t.In(tz)
}

func (d *dbBase) TimeToDB(t *time.Time, tz *time.Location) {
	*t = t.In(tz)
}

func (d *dbBase) DbTypes() map[string]string {
	return nil
}

func (d *dbBase) GetTables(db dbQuerier) (map[string]bool, error) {
	tables := make(map[string]bool)
	query := d.ins.ShowTablesQuery()
	rows, err := db.Query(query)
	if err != nil {
		return tables, err
	}

	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return tables, err
		}
		if table != "" {
			tables[table] = true
		}
	}

	return tables, nil
}

func (d *dbBase) GetColumns(db dbQuerier, table string) (map[string][3]string, error) {
	columns := make(map[string][3]string)
	query := d.ins.ShowColumnsQuery(table)
	rows, err := db.Query(query)
	if err != nil {
		return columns, err
	}

	for rows.Next() {
		var (
			name string
			typ  string
			null string
		)
		err := rows.Scan(&name, &typ, &null)
		if err != nil {
			return columns, err
		}
		columns[name] = [3]string{name, typ, null}
	}

	return columns, nil
}

func (d *dbBase) ShowTablesQuery() string {
	panic(ErrNotImplement)
}

func (d *dbBase) ShowColumnsQuery(table string) string {
	panic(ErrNotImplement)
}

func (d *dbBase) IndexExists(dbQuerier, string, string) bool {
	panic(ErrNotImplement)
}
