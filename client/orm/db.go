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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm/internal/buffers"

	"github.com/beego/beego/v2/client/orm/internal/logs"

	"github.com/beego/beego/v2/client/orm/internal/utils"

	"github.com/beego/beego/v2/client/orm/internal/models"

	"github.com/beego/beego/v2/client/orm/hints"
)

// ErrMissPK missing pk error
var ErrMissPK = errors.New("missed pk value")

var operators = map[string]bool{
	"exact":       true,
	"iexact":      true,
	"strictexact": true,
	"contains":    true,
	"icontains":   true,
	// "regex":       true,
	// "iregex":      true,
	"gt":          true,
	"gte":         true,
	"lt":          true,
	"lte":         true,
	"eq":          true,
	"nq":          true,
	"ne":          true,
	"startswith":  true,
	"endswith":    true,
	"istartswith": true,
	"iendswith":   true,
	"in":          true,
	"between":     true,
	// "year":        true,
	// "month":       true,
	// "day":         true,
	// "week_day":    true,
	"isnull": true,
	// "search":      true,
}

// an instance of dbBaser interface/
type dbBase struct {
	ins dbBaser
}

// check dbBase implements dbBaser interface.
var _ dbBaser = new(dbBase)

// get struct Columns values as interface slice.
func (d *dbBase) collectValues(mi *models.ModelInfo, ind reflect.Value, cols []string, skipAuto bool, insert bool, names *[]string, tz *time.Location) (values []interface{}, autoFields []string, err error) {
	if names == nil {
		ns := make([]string, 0, len(cols))
		names = &ns
	}
	values = make([]interface{}, 0, len(cols))

	for _, column := range cols {
		var fi *models.FieldInfo
		if fi, _ = mi.Fields.GetByAny(column); fi != nil {
			column = fi.Column
		} else {
			panic(fmt.Errorf("wrong db field/column name `%s` for model `%s`", column, mi.FullName))
		}
		if !fi.DBcol || fi.Auto && skipAuto {
			continue
		}
		value, err := d.collectFieldValue(mi, fi, ind, insert, tz)
		if err != nil {
			return nil, nil, err
		}

		// ignore empty value auto field
		if insert && fi.Auto {
			if fi.FieldType&IsPositiveIntegerField > 0 {
				if vu, ok := value.(uint64); !ok || vu == 0 {
					continue
				}
			} else {
				if vu, ok := value.(int64); !ok || vu == 0 {
					continue
				}
			}
			autoFields = append(autoFields, fi.Column)
		}

		*names, values = append(*names, column), append(values, value)
	}

	return
}

// get one field value in struct column as interface.
func (d *dbBase) collectFieldValue(mi *models.ModelInfo, fi *models.FieldInfo, ind reflect.Value, insert bool, tz *time.Location) (interface{}, error) {
	var value interface{}
	if fi.Pk {
		_, value, _ = getExistPk(mi, ind)
	} else {
		field := ind.FieldByIndex(fi.FieldIndex)
		if fi.IsFielder {
			f := field.Addr().Interface().(models.Fielder)
			value = f.RawValue()
		} else {
			switch fi.FieldType {
			case TypeBooleanField:
				if nb, ok := field.Interface().(sql.NullBool); ok {
					value = nil
					if nb.Valid {
						value = nb.Bool
					}
				} else if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						value = nil
					} else {
						value = field.Elem().Bool()
					}
				} else {
					value = field.Bool()
				}
			case TypeVarCharField, TypeCharField, TypeTextField, TypeJSONField, TypeJsonbField:
				if ns, ok := field.Interface().(sql.NullString); ok {
					value = nil
					if ns.Valid {
						value = ns.String
					}
				} else if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						value = nil
					} else {
						value = field.Elem().String()
					}
				} else {
					value = field.String()
				}
			case TypeFloatField, TypeDecimalField:
				if nf, ok := field.Interface().(sql.NullFloat64); ok {
					value = nil
					if nf.Valid {
						value = nf.Float64
					}
				} else if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						value = nil
					} else {
						value = field.Elem().Float()
					}
				} else {
					vu := field.Interface()
					if _, ok := vu.(float32); ok {
						value, _ = utils.StrTo(utils.ToStr(vu)).Float64()
					} else {
						value = field.Float()
					}
				}
			case TypeTimeField, TypeDateField, TypeDateTimeField:
				value = field.Interface()
				if t, ok := value.(time.Time); ok {
					d.ins.TimeToDB(&t, tz)
					if t.IsZero() {
						value = nil
					} else {
						value = t
					}
				}
			default:
				switch {
				case fi.FieldType&IsPositiveIntegerField > 0:
					if field.Kind() == reflect.Ptr {
						if field.IsNil() {
							value = nil
						} else {
							value = field.Elem().Uint()
						}
					} else {
						value = field.Uint()
					}
				case fi.FieldType&IsIntegerField > 0:
					if ni, ok := field.Interface().(sql.NullInt64); ok {
						value = nil
						if ni.Valid {
							value = ni.Int64
						}
					} else if field.Kind() == reflect.Ptr {
						if field.IsNil() {
							value = nil
						} else {
							value = field.Elem().Int()
						}
					} else {
						value = field.Int()
					}
				case fi.FieldType&IsRelField > 0:
					if field.IsNil() {
						value = nil
					} else {
						if _, vu, ok := getExistPk(fi.RelModelInfo, reflect.Indirect(field)); ok {
							value = vu
						} else {
							value = nil
						}
					}
					if !fi.Null && value == nil {
						return nil, fmt.Errorf("field `%s` cannot be NULL", fi.FullName)
					}
				}
			}
		}
		switch fi.FieldType {
		case TypeTimeField, TypeDateField, TypeDateTimeField:
			if fi.AutoNow || fi.AutoNowAdd && insert {
				if insert {
					if t, ok := value.(time.Time); ok && !t.IsZero() {
						break
					}
				}
				tnow := time.Now()
				d.ins.TimeToDB(&tnow, tz)
				value = tnow
				if fi.IsFielder {
					f := field.Addr().Interface().(models.Fielder)
					f.SetRaw(tnow.In(DefaultTimeLoc))
				} else if field.Kind() == reflect.Ptr {
					v := tnow.In(DefaultTimeLoc)
					field.Set(reflect.ValueOf(&v))
				} else {
					field.Set(reflect.ValueOf(tnow.In(DefaultTimeLoc)))
				}
			}
		case TypeJSONField, TypeJsonbField:
			if s, ok := value.(string); (ok && len(s) == 0) || value == nil {
				if fi.ColDefault && fi.Initial.Exist() {
					value = fi.Initial.String()
				} else {
					value = nil
				}
			}
		}
	}
	return value, nil
}

// PrepareInsert create insert sql preparation statement object.
func (d *dbBase) PrepareInsert(ctx context.Context, q dbQuerier, mi *models.ModelInfo) (stmtQuerier, string, error) {
	Q := d.ins.TableQuote()

	dbcols := make([]string, 0, len(mi.Fields.DBcols))
	marks := make([]string, 0, len(mi.Fields.DBcols))
	for _, fi := range mi.Fields.FieldsDB {
		if !fi.Auto {
			dbcols = append(dbcols, fi.Column)
			marks = append(marks, "?")
		}
	}
	qmarks := strings.Join(marks, ", ")
	sep := fmt.Sprintf("%s, %s", Q, Q)
	columns := strings.Join(dbcols, sep)

	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, mi.Table, Q, Q, columns, Q, qmarks)

	d.ins.ReplaceMarks(&query)

	d.ins.HasReturningID(mi, &query)

	stmt, err := q.PrepareContext(ctx, query)
	return stmt, query, err
}

// InsertStmt insert struct with prepared statement and given struct reflect value.
func (d *dbBase) InsertStmt(ctx context.Context, stmt stmtQuerier, mi *models.ModelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	values, _, err := d.collectValues(mi, ind, mi.Fields.DBcols, true, true, nil, tz)
	if err != nil {
		return 0, err
	}

	if d.ins.HasReturningID(mi, nil) {
		row := stmt.QueryRow(values...)
		var id int64
		err := row.Scan(&id)
		return id, err
	}
	res, err := stmt.ExecContext(ctx, values...)
	if err == nil {
		return res.LastInsertId()
	}
	return 0, err
}

// query sql ,read records and persist in dbBaser.
func (d *dbBase) Read(ctx context.Context, q dbQuerier, mi *models.ModelInfo, ind reflect.Value, tz *time.Location, cols []string, isForUpdate bool) error {
	var whereCols []string
	var args []interface{}

	// if specify cols length > 0, then use it for where condition.
	if len(cols) > 0 {
		var err error
		whereCols = make([]string, 0, len(cols))
		args, _, err = d.collectValues(mi, ind, cols, false, false, &whereCols, tz)
		if err != nil {
			return err
		}
	} else {
		// default use pk value as where condtion.
		pkColumn, pkValue, ok := getExistPk(mi, ind)
		if !ok {
			return ErrMissPK
		}
		whereCols = []string{pkColumn}
		args = append(args, pkValue)
	}

	Q := d.ins.TableQuote()

	sep := fmt.Sprintf("%s, %s", Q, Q)
	sels := strings.Join(mi.Fields.DBcols, sep)
	colsNum := len(mi.Fields.DBcols)

	sep = fmt.Sprintf("%s = ? AND %s", Q, Q)
	wheres := strings.Join(whereCols, sep)

	forUpdate := ""
	if isForUpdate {
		forUpdate = "FOR UPDATE"
	}

	query := fmt.Sprintf("SELECT %s%s%s FROM %s%s%s WHERE %s%s%s = ? %s", Q, sels, Q, Q, mi.Table, Q, Q, wheres, Q, forUpdate)

	refs := make([]interface{}, colsNum)
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	d.ins.ReplaceMarks(&query)

	row := q.QueryRowContext(ctx, query, args...)
	if err := row.Scan(refs...); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRows
		}
		return err
	}
	elm := reflect.New(mi.AddrField.Elem().Type())
	mind := reflect.Indirect(elm)
	d.setColsValues(mi, &mind, mi.Fields.DBcols, refs, tz)
	ind.Set(mind)
	return nil
}

// Insert execute insert sql dbQuerier with given struct reflect.Value.
func (d *dbBase) Insert(ctx context.Context, q dbQuerier, mi *models.ModelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	names := make([]string, 0, len(mi.Fields.DBcols))
	values, autoFields, err := d.collectValues(mi, ind, mi.Fields.DBcols, false, true, &names, tz)
	if err != nil {
		return 0, err
	}

	id, err := d.InsertValue(ctx, q, mi, false, names, values)
	if err != nil {
		return 0, err
	}

	if len(autoFields) > 0 {
		err = d.ins.setval(ctx, q, mi, autoFields)
	}
	return id, err
}

// InsertMulti multi-insert sql with given slice struct reflect.Value.
func (d *dbBase) InsertMulti(ctx context.Context, q dbQuerier, mi *models.ModelInfo, sind reflect.Value, bulk int, tz *time.Location) (int64, error) {
	var (
		cnt    int64
		nums   int
		values []interface{}
		names  []string
	)

	length, autoFields := sind.Len(), make([]string, 0, 1)

	for i := 1; i <= length; i++ {

		ind := reflect.Indirect(sind.Index(i - 1))

		if i == 1 {
			var (
				vus []interface{}
				err error
			)
			vus, autoFields, err = d.collectValues(mi, ind, mi.Fields.DBcols, false, true, &names, tz)
			if err != nil {
				return cnt, err
			}
			values = make([]interface{}, bulk*len(vus))
			nums += copy(values, vus)
		} else {
			vus, _, err := d.collectValues(mi, ind, mi.Fields.DBcols, false, true, nil, tz)
			if err != nil {
				return cnt, err
			}

			if len(vus) != len(names) {
				return cnt, ErrArgs
			}

			nums += copy(values[nums:], vus)
		}

		if i > 1 && i%bulk == 0 || length == i {
			num, err := d.InsertValue(ctx, q, mi, true, names, values[:nums])
			if err != nil {
				return cnt, err
			}
			cnt += num
			nums = 0
		}
	}

	var err error
	if len(autoFields) > 0 {
		err = d.ins.setval(ctx, q, mi, autoFields)
	}

	return cnt, err
}

// InsertValue execute insert sql with given struct and given values.
// insert the given values, not the field values in struct.
func (d *dbBase) InsertValue(ctx context.Context, q dbQuerier, mi *models.ModelInfo, isMulti bool, names []string, values []interface{}) (int64, error) {
	query := d.InsertValueSQL(names, values, isMulti, mi)

	if isMulti || !d.ins.HasReturningID(mi, &query) {
		res, err := q.ExecContext(ctx, query, values...)
		if err == nil {
			if isMulti {
				return res.RowsAffected()
			}

			lastInsertId, err := res.LastInsertId()
			if err != nil {
				logs.DebugLog.Println(ErrLastInsertIdUnavailable, ':', err)
				return lastInsertId, ErrLastInsertIdUnavailable
			} else {
				return lastInsertId, nil
			}
		}
		return 0, err
	}
	row := q.QueryRowContext(ctx, query, values...)
	var id int64
	err := row.Scan(&id)
	return id, err
}

func (d *dbBase) InsertValueSQL(names []string, values []interface{}, isMulti bool, mi *models.ModelInfo) string {
	buf := buffers.Get()
	defer buffers.Put(buf)

	Q := d.ins.TableQuote()

	_, _ = buf.WriteString("INSERT INTO ")
	_, _ = buf.WriteString(Q)
	_, _ = buf.WriteString(mi.Table)
	_, _ = buf.WriteString(Q)

	_, _ = buf.WriteString(" (")
	for i, name := range names {
		if i > 0 {
			_, _ = buf.WriteString(", ")
		}
		_, _ = buf.WriteString(Q)
		_, _ = buf.WriteString(name)
		_, _ = buf.WriteString(Q)
	}
	_, _ = buf.WriteString(") VALUES (")

	marks := make([]string, len(names))
	for i := range marks {
		marks[i] = "?"
	}
	qmarks := strings.Join(marks, ", ")

	_, _ = buf.WriteString(qmarks)

	multi := len(values) / len(names)

	if isMulti && multi > 1 {
		for i := 0; i < multi-1; i++ {
			_, _ = buf.WriteString("), (")
			_, _ = buf.WriteString(qmarks)
		}
	}

	_ = buf.WriteByte(')')

	query := buf.String()
	d.ins.ReplaceMarks(&query)

	return query
}

// InsertOrUpdate a row
// If your primary key or unique column conflict will update
// If no will insert
func (d *dbBase) InsertOrUpdate(ctx context.Context, q dbQuerier, mi *models.ModelInfo, ind reflect.Value, a *alias, args ...string) (int64, error) {
	args0 := ""
	iouStr := ""
	argsMap := map[string]string{}
	switch a.Driver {
	case DRMySQL:
		iouStr = "ON DUPLICATE KEY UPDATE"
	case DRPostgres:
		if len(args) == 0 {
			return 0, fmt.Errorf("`%s` use InsertOrUpdate must have a conflict column", a.DriverName)
		}
		args0 = strings.ToLower(args[0])
		iouStr = fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET", args0)
	default:
		return 0, fmt.Errorf("`%s` nonsupport InsertOrUpdate in beego", a.DriverName)
	}

	// Get on the key-value pairs
	for _, v := range args {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			argsMap[strings.ToLower(kv[0])] = kv[1]
		}
	}

	names := make([]string, 0, len(mi.Fields.DBcols)-1)
	Q := d.ins.TableQuote()
	values, _, err := d.collectValues(mi, ind, mi.Fields.DBcols, true, true, &names, a.TZ)
	if err != nil {
		return 0, err
	}

	marks := make([]string, len(names))
	updateValues := make([]interface{}, 0)
	updates := make([]string, len(names))
	var conflitValue interface{}
	for i, v := range names {
		// identifier in database may not be case-sensitive, so quote it
		v = fmt.Sprintf("%s%s%s", Q, v, Q)
		marks[i] = "?"
		valueStr := argsMap[strings.ToLower(v)]
		if v == args0 {
			conflitValue = values[i]
		}
		if valueStr != "" {
			switch a.Driver {
			case DRMySQL:
				updates[i] = v + "=" + valueStr
			case DRPostgres:
				if conflitValue != nil {
					// postgres ON CONFLICT DO UPDATE SET can`t use colu=colu+values
					updates[i] = fmt.Sprintf("%s=(select %s from %s where %s = ? )", v, valueStr, mi.Table, args0)
					updateValues = append(updateValues, conflitValue)
				} else {
					return 0, fmt.Errorf("`%s` must be in front of `%s` in your struct", args0, v)
				}
			}
		} else {
			updates[i] = v + "=?"
			updateValues = append(updateValues, values[i])
		}
	}

	values = append(values, updateValues...)

	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	qupdates := strings.Join(updates, ", ")
	columns := strings.Join(names, sep)

	// conflitValue maybe is a int,can`t use fmt.Sprintf
	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s) %s "+qupdates, Q, mi.Table, Q, Q, columns, Q, qmarks, iouStr)

	d.ins.ReplaceMarks(&query)

	if !d.ins.HasReturningID(mi, &query) {
		res, err := q.ExecContext(ctx, query, values...)
		if err == nil {
			lastInsertId, err := res.LastInsertId()
			if err != nil {
				logs.DebugLog.Println(ErrLastInsertIdUnavailable, ':', err)
				return lastInsertId, ErrLastInsertIdUnavailable
			} else {
				return lastInsertId, nil
			}
		}
		return 0, err
	}

	row := q.QueryRowContext(ctx, query, values...)
	var id int64
	err = row.Scan(&id)
	if err != nil && err.Error() == `pq: syntax error at or near "ON"` {
		err = fmt.Errorf("postgres version must 9.5 or higher")
	}
	return id, err
}

// Update execute update sql dbQuerier with given struct reflect.Value.
func (d *dbBase) Update(ctx context.Context, q dbQuerier, mi *models.ModelInfo, ind reflect.Value, tz *time.Location, cols []string) (int64, error) {
	pkName, pkValue, ok := getExistPk(mi, ind)
	if !ok {
		return 0, ErrMissPK
	}

	var setNames []string

	// if specify cols length is zero, then commit all Columns.
	if len(cols) == 0 {
		cols = mi.Fields.DBcols
		setNames = make([]string, 0, len(mi.Fields.DBcols)-1)
	} else {
		setNames = make([]string, 0, len(cols))
	}

	setValues, _, err := d.collectValues(mi, ind, cols, true, false, &setNames, tz)
	if err != nil {
		return 0, err
	}

	var findAutoNowAdd, findAutoNow bool
	var index int
	for i, col := range setNames {
		if mi.Fields.GetByColumn(col).AutoNowAdd {
			index = i
			findAutoNowAdd = true
		}
		if mi.Fields.GetByColumn(col).AutoNow {
			findAutoNow = true
		}
	}
	if findAutoNowAdd {
		setNames = append(setNames[0:index], setNames[index+1:]...)
		setValues = append(setValues[0:index], setValues[index+1:]...)
	}

	if !findAutoNow {
		for col, info := range mi.Fields.Columns {
			if info.AutoNow {
				setNames = append(setNames, col)
				setValues = append(setValues, time.Now())
			}
		}
	}

	setValues = append(setValues, pkValue)

	query := d.UpdateSQL(setNames, pkName, mi)

	res, err := q.ExecContext(ctx, query, setValues...)
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

func (d *dbBase) UpdateSQL(setNames []string, pkName string, mi *models.ModelInfo) string {
	buf := buffers.Get()
	defer buffers.Put(buf)

	Q := d.ins.TableQuote()

	_, _ = buf.WriteString("UPDATE ")
	_, _ = buf.WriteString(Q)
	_, _ = buf.WriteString(mi.Table)
	_, _ = buf.WriteString(Q)
	_, _ = buf.WriteString(" SET ")

	for i, name := range setNames {
		if i > 0 {
			_, _ = buf.WriteString(", ")
		}
		_, _ = buf.WriteString(Q)
		_, _ = buf.WriteString(name)
		_, _ = buf.WriteString(Q)
		_, _ = buf.WriteString(" = ?")
	}

	_, _ = buf.WriteString(" WHERE ")
	_, _ = buf.WriteString(Q)
	_, _ = buf.WriteString(pkName)
	_, _ = buf.WriteString(Q)
	_, _ = buf.WriteString(" = ?")

	query := buf.String()
	d.ins.ReplaceMarks(&query)

	return query
}

// Delete execute delete sql dbQuerier with given struct reflect.Value.
// delete index is pk.
func (d *dbBase) Delete(ctx context.Context, q dbQuerier, mi *models.ModelInfo, ind reflect.Value, tz *time.Location, cols []string) (int64, error) {
	var whereCols []string
	var args []interface{}
	// if specify cols length > 0, then use it for where condition.
	if len(cols) > 0 {
		var err error
		whereCols = make([]string, 0, len(cols))
		args, _, err = d.collectValues(mi, ind, cols, false, false, &whereCols, tz)
		if err != nil {
			return 0, err
		}
	} else {
		// default use pk value as where condtion.
		pkColumn, pkValue, ok := getExistPk(mi, ind)
		if !ok {
			return 0, ErrMissPK
		}
		whereCols = []string{pkColumn}
		args = append(args, pkValue)
	}

	query := d.DeleteSQL(whereCols, mi)

	res, err := q.ExecContext(ctx, query, args...)
	if err == nil {
		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		if num > 0 {
			err := d.deleteRels(ctx, q, mi, args, tz)
			if err != nil {
				return num, err
			}
		}
		return num, err
	}
	return 0, err
}

func (d *dbBase) DeleteSQL(whereCols []string, mi *models.ModelInfo) string {
	buf := buffers.Get()
	defer buffers.Put(buf)

	Q := d.ins.TableQuote()

	_, _ = buf.WriteString("DELETE FROM ")
	_, _ = buf.WriteString(Q)
	_, _ = buf.WriteString(mi.Table)
	_, _ = buf.WriteString(Q)
	_, _ = buf.WriteString(" WHERE ")

	for i, col := range whereCols {
		if i > 0 {
			_, _ = buf.WriteString(" AND ")
		}
		_, _ = buf.WriteString(Q)
		_, _ = buf.WriteString(col)
		_, _ = buf.WriteString(Q)
		_, _ = buf.WriteString(" = ?")
	}

	query := buf.String()

	d.ins.ReplaceMarks(&query)

	return query
}

// UpdateBatch update table-related record by querySet.
// need querySet not struct reflect.Value to update related records.
func (d *dbBase) UpdateBatch(ctx context.Context, q dbQuerier, qs *querySet, mi *models.ModelInfo, cond *Condition, params Params, tz *time.Location) (int64, error) {
	columns := make([]string, 0, len(params))
	values := make([]interface{}, 0, len(params))
	for col, val := range params {
		if fi, ok := mi.Fields.GetByAny(col); !ok || !fi.DBcol {
			panic(fmt.Errorf("wrong field/column name `%s`", col))
		} else {
			columns = append(columns, fi.Column)
			values = append(values, val)
		}
	}

	if len(columns) == 0 {
		panic(fmt.Errorf("update params cannot empty"))
	}

	tables := newDbTables(mi, d.ins)
	var specifyIndexes string
	if qs != nil {
		tables.parseRelated(qs.related, qs.relDepth)
		specifyIndexes = tables.getIndexSql(mi.Table, qs.useIndex, qs.indexes)
	}

	where, args := tables.getCondSQL(cond, false, tz)

	values = append(values, args...)

	join := tables.getJoinSQL()

	var query, T string

	Q := d.ins.TableQuote()

	if d.ins.SupportUpdateJoin() {
		T = "T0."
	}

	cols := make([]string, 0, len(columns))

	for i, v := range columns {
		col := fmt.Sprintf("%s%s%s%s", T, Q, v, Q)
		if c, ok := values[i].(colValue); ok {
			switch c.opt {
			case ColAdd:
				cols = append(cols, col+" = "+col+" + ?")
			case ColMinus:
				cols = append(cols, col+" = "+col+" - ?")
			case ColMultiply:
				cols = append(cols, col+" = "+col+" * ?")
			case ColExcept:
				cols = append(cols, col+" = "+col+" / ?")
			case ColBitAnd:
				cols = append(cols, col+" = "+col+" & ?")
			case ColBitRShift:
				cols = append(cols, col+" = "+col+" >> ?")
			case ColBitLShift:
				cols = append(cols, col+" = "+col+" << ?")
			case ColBitXOR:
				cols = append(cols, col+" = "+col+" ^ ?")
			case ColBitOr:
				cols = append(cols, col+" = "+col+" | ?")
			}
			values[i] = c.value
		} else {
			cols = append(cols, col+" = ?")
		}
	}

	sets := strings.Join(cols, ", ") + " "

	if d.ins.SupportUpdateJoin() {
		query = fmt.Sprintf("UPDATE %s%s%s T0 %s%sSET %s%s", Q, mi.Table, Q, specifyIndexes, join, sets, where)
	} else {
		supQuery := fmt.Sprintf("SELECT T0.%s%s%s FROM %s%s%s T0 %s%s%s",
			Q, mi.Fields.Pk.Column, Q,
			Q, mi.Table, Q,
			specifyIndexes, join, where)
		query = fmt.Sprintf("UPDATE %s%s%s SET %sWHERE %s%s%s IN ( %s )", Q, mi.Table, Q, sets, Q, mi.Fields.Pk.Column, Q, supQuery)
	}

	d.ins.ReplaceMarks(&query)
	res, err := q.ExecContext(ctx, query, values...)
	if err == nil {
		return res.RowsAffected()
	}
	return 0, err
}

// delete related records.
// do UpdateBanch or DeleteBanch by condition of tables' relationship.
func (d *dbBase) deleteRels(ctx context.Context, q dbQuerier, mi *models.ModelInfo, args []interface{}, tz *time.Location) error {
	for _, fi := range mi.Fields.FieldsReverse {
		fi = fi.ReverseFieldInfo
		switch fi.OnDelete {
		case models.OdCascade:
			cond := NewCondition().And(fmt.Sprintf("%s__in", fi.Name), args...)
			_, err := d.DeleteBatch(ctx, q, nil, fi.Mi, cond, tz)
			if err != nil {
				return err
			}
		case models.OdSetDefault, models.OdSetNULL:
			cond := NewCondition().And(fmt.Sprintf("%s__in", fi.Name), args...)
			params := Params{fi.Column: nil}
			if fi.OnDelete == models.OdSetDefault {
				params[fi.Column] = fi.Initial.String()
			}
			_, err := d.UpdateBatch(ctx, q, nil, fi.Mi, cond, params, tz)
			if err != nil {
				return err
			}
		case models.OdDoNothing:
		}
	}
	return nil
}

// DeleteBatch delete table-related records.
func (d *dbBase) DeleteBatch(ctx context.Context, q dbQuerier, qs *querySet, mi *models.ModelInfo, cond *Condition, tz *time.Location) (int64, error) {
	tables := newDbTables(mi, d.ins)
	tables.skipEnd = true

	var specifyIndexes string
	if qs != nil {
		tables.parseRelated(qs.related, qs.relDepth)
		specifyIndexes = tables.getIndexSql(mi.Table, qs.useIndex, qs.indexes)
	}

	if cond == nil || cond.IsEmpty() {
		panic(fmt.Errorf("delete operation cannot execute without condition"))
	}

	Q := d.ins.TableQuote()

	where, args := tables.getCondSQL(cond, false, tz)
	join := tables.getJoinSQL()

	cols := fmt.Sprintf("T0.%s%s%s", Q, mi.Fields.Pk.Column, Q)
	query := fmt.Sprintf("SELECT %s FROM %s%s%s T0 %s%s%s", cols, Q, mi.Table, Q, specifyIndexes, join, where)

	d.ins.ReplaceMarks(&query)

	var rs *sql.Rows
	r, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	rs = r
	defer rs.Close()

	var ref interface{}
	args = make([]interface{}, 0)
	cnt := 0
	for rs.Next() {
		if err := rs.Scan(&ref); err != nil {
			return 0, err
		}
		pkValue, err := d.convertValueFromDB(mi.Fields.Pk, reflect.ValueOf(ref).Interface(), tz)
		if err != nil {
			return 0, err
		}
		args = append(args, pkValue)
		cnt++
	}

	if err = rs.Err(); err != nil {
		return 0, err
	}

	if cnt == 0 {
		return 0, nil
	}

	marks := make([]string, len(args))
	for i := range marks {
		marks[i] = "?"
	}
	sqlIn := fmt.Sprintf("IN (%s)", strings.Join(marks, ", "))
	query = fmt.Sprintf("DELETE FROM %s%s%s WHERE %s%s%s %s", Q, mi.Table, Q, Q, mi.Fields.Pk.Column, Q, sqlIn)

	d.ins.ReplaceMarks(&query)
	res, err := q.ExecContext(ctx, query, args...)
	if err == nil {
		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		if num > 0 {
			err := d.deleteRels(ctx, q, mi, args, tz)
			if err != nil {
				return num, err
			}
		}
		return num, nil
	}
	return 0, err
}

// ReadBatch read related records.
func (d *dbBase) ReadBatch(ctx context.Context, q dbQuerier, qs *querySet, mi *models.ModelInfo, cond *Condition, container interface{}, tz *time.Location, cols []string) (int64, error) {
	val := reflect.ValueOf(container)
	ind := reflect.Indirect(val)

	unregister := true
	one := true
	isPtr := true
	name := ""

	if val.Kind() == reflect.Ptr {
		fn := ""
		if ind.Kind() == reflect.Slice {
			one = false
			typ := ind.Type().Elem()
			switch typ.Kind() {
			case reflect.Ptr:
				fn = models.GetFullName(typ.Elem())
			case reflect.Struct:
				isPtr = false
				fn = models.GetFullName(typ)
				name = models.GetTableName(reflect.New(typ))
			}
		} else {
			fn = models.GetFullName(ind.Type())
			name = models.GetTableName(ind)
		}
		unregister = fn != mi.FullName
	}

	if unregister {
		RegisterModel(container)
	}

	rlimit := qs.limit
	offset := qs.offset

	Q := d.ins.TableQuote()

	var tCols []string
	if len(cols) > 0 {
		hasRel := len(qs.related) > 0 || qs.relDepth > 0
		tCols = make([]string, 0, len(cols))
		var maps map[string]bool
		if hasRel {
			maps = make(map[string]bool)
		}
		for _, col := range cols {
			if fi, ok := mi.Fields.GetByAny(col); ok {
				tCols = append(tCols, fi.Column)
				if hasRel {
					maps[fi.Column] = true
				}
			} else {
				return 0, fmt.Errorf("wrong field/column name `%s`", col)
			}
		}
		if hasRel {
			for _, fi := range mi.Fields.FieldsDB {
				if fi.FieldType&IsRelField > 0 {
					if !maps[fi.Column] {
						tCols = append(tCols, fi.Column)
					}
				}
			}
		}
	} else {
		tCols = mi.Fields.DBcols
	}

	colsNum := len(tCols)
	sep := fmt.Sprintf("%s, T0.%s", Q, Q)
	sels := fmt.Sprintf("T0.%s%s%s", Q, strings.Join(tCols, sep), Q)

	tables := newDbTables(mi, d.ins)
	tables.parseRelated(qs.related, qs.relDepth)

	where, args := tables.getCondSQL(cond, false, tz)
	groupBy := tables.getGroupSQL(qs.groups)
	orderBy := tables.getOrderSQL(qs.orders)
	limit := tables.getLimitSQL(mi, offset, rlimit)
	join := tables.getJoinSQL()
	specifyIndexes := tables.getIndexSql(mi.Table, qs.useIndex, qs.indexes)

	for _, tbl := range tables.tables {
		if tbl.sel {
			colsNum += len(tbl.mi.Fields.DBcols)
			sep := fmt.Sprintf("%s, %s.%s", Q, tbl.index, Q)
			sels += fmt.Sprintf(", %s.%s%s%s", tbl.index, Q, strings.Join(tbl.mi.Fields.DBcols, sep), Q)
		}
	}

	sqlSelect := "SELECT"
	if qs.distinct {
		sqlSelect += " DISTINCT"
	}
	if qs.aggregate != "" {
		sels = qs.aggregate
	}
	query := fmt.Sprintf("%s %s FROM %s%s%s T0 %s%s%s%s%s%s",
		sqlSelect, sels, Q, mi.Table, Q,
		specifyIndexes, join, where, groupBy, orderBy, limit)

	if qs.forUpdate {
		query += " FOR UPDATE"
	}

	d.ins.ReplaceMarks(&query)

	rs, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	defer rs.Close()

	slice := ind
	if unregister {
		mi, _ = defaultModelCache.get(name)
		tCols = mi.Fields.DBcols
		colsNum = len(tCols)
	}

	refs := make([]interface{}, colsNum)
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}
	var cnt int64
	for rs.Next() {
		if one && cnt == 0 || !one {
			if err := rs.Scan(refs...); err != nil {
				return 0, err
			}

			elm := reflect.New(mi.AddrField.Elem().Type())
			mind := reflect.Indirect(elm)

			cacheV := make(map[string]*reflect.Value)
			cacheM := make(map[string]*models.ModelInfo)
			trefs := refs

			d.setColsValues(mi, &mind, tCols, refs[:len(tCols)], tz)
			trefs = refs[len(tCols):]

			for _, tbl := range tables.tables {
				// loop selected tables
				if tbl.sel {
					last := mind
					names := ""
					mmi := mi
					// loop cascade models
					for _, name := range tbl.names {
						names += name
						if val, ok := cacheV[names]; ok {
							last = *val
							mmi = cacheM[names]
						} else {
							fi := mmi.Fields.GetByName(name)
							lastm := mmi
							mmi = fi.RelModelInfo
							field := last
							if last.Kind() != reflect.Invalid {
								field = reflect.Indirect(last.FieldByIndex(fi.FieldIndex))
								if field.IsValid() {
									d.setColsValues(mmi, &field, mmi.Fields.DBcols, trefs[:len(mmi.Fields.DBcols)], tz)
									for _, fi := range mmi.Fields.FieldsReverse {
										if fi.InModel && fi.ReverseFieldInfo.Mi == lastm {
											if fi.ReverseFieldInfo != nil {
												f := field.FieldByIndex(fi.FieldIndex)
												if f.Kind() == reflect.Ptr {
													f.Set(last.Addr())
												}
											}
										}
									}
									last = field
								}
							}
							cacheV[names] = &field
							cacheM[names] = mmi
						}
					}
					trefs = trefs[len(mmi.Fields.DBcols):]
				}
			}

			if one {
				ind.Set(mind)
			} else {
				if cnt == 0 {
					// you can use an empty & caped container list
					// orm will not replace it
					if ind.Len() != 0 {
						// if container is not empty
						// create a new one
						slice = reflect.New(ind.Type()).Elem()
					}
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

	if err = rs.Err(); err != nil {
		return 0, err
	}

	if !one {
		if cnt > 0 {
			ind.Set(slice)
		} else {
			// when a result is empty and container is nil
			// to set an empty container
			if ind.IsNil() {
				ind.Set(reflect.MakeSlice(ind.Type(), 0, 0))
			}
		}
	}

	return cnt, nil
}

// Count excute count sql and return count result int64.
func (d *dbBase) Count(ctx context.Context, q dbQuerier, qs *querySet, mi *models.ModelInfo, cond *Condition, tz *time.Location) (cnt int64, err error) {
	tables := newDbTables(mi, d.ins)
	tables.parseRelated(qs.related, qs.relDepth)

	where, args := tables.getCondSQL(cond, false, tz)
	groupBy := tables.getGroupSQL(qs.groups)
	tables.getOrderSQL(qs.orders)
	join := tables.getJoinSQL()
	specifyIndexes := tables.getIndexSql(mi.Table, qs.useIndex, qs.indexes)

	Q := d.ins.TableQuote()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s%s%s T0 %s%s%s%s",
		Q, mi.Table, Q,
		specifyIndexes, join, where, groupBy)

	if groupBy != "" {
		query = fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS T", query)
	}

	d.ins.ReplaceMarks(&query)

	row := q.QueryRowContext(ctx, query, args...)
	err = row.Scan(&cnt)
	return
}

// GenerateOperatorSQL generate sql with replacing operator string placeholders and replaced values.
func (d *dbBase) GenerateOperatorSQL(mi *models.ModelInfo, fi *models.FieldInfo, operator string, args []interface{}, tz *time.Location) (string, []interface{}) {
	var sql string
	params := getFlatParams(fi, args, tz)

	if len(params) == 0 {
		panic(fmt.Errorf("operator `%s` need at least one args", operator))
	}
	arg := params[0]

	switch operator {
	case "in":
		marks := make([]string, len(params))
		for i := range marks {
			marks[i] = "?"
		}
		sql = fmt.Sprintf("IN (%s)", strings.Join(marks, ", "))
	case "between":
		if len(params) != 2 {
			panic(fmt.Errorf("operator `%s` need 2 args not %d", operator, len(params)))
		}
		sql = "BETWEEN ? AND ?"
	default:
		if len(params) > 1 {
			panic(fmt.Errorf("operator `%s` need 1 args not %d", operator, len(params)))
		}
		sql = d.ins.OperatorSQL(operator)
		switch operator {
		case "exact":
			if arg == nil {
				params[0] = "IS NULL"
			}
		case "iexact", "contains", "icontains", "startswith", "endswith", "istartswith", "iendswith":
			param := strings.Replace(utils.ToStr(arg), `%`, `\%`, -1)
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
				panic(fmt.Errorf("operator `%s` need a bool value not `%T`", operator, arg))
			}
		}
	}
	return sql, params
}

// GenerateOperatorLeftCol gernerate sql string with inner function, such as UPPER(text).
func (d *dbBase) GenerateOperatorLeftCol(*models.FieldInfo, string, *string) {
	// default not use
}

// set values to struct column.
func (d *dbBase) setColsValues(mi *models.ModelInfo, ind *reflect.Value, cols []string, values []interface{}, tz *time.Location) {
	for i, column := range cols {
		val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()

		fi := mi.Fields.GetByColumn(column)

		field := ind.FieldByIndex(fi.FieldIndex)

		value, err := d.convertValueFromDB(fi, val, tz)
		if err != nil {
			panic(fmt.Errorf("Raw value: `%v` %s", val, err.Error()))
		}

		_, err = d.setFieldValue(fi, value, field)

		if err != nil {
			panic(fmt.Errorf("Raw value: `%v` %s", val, err.Error()))
		}
	}
}

// convert value from database result to value following in field type.
func (d *dbBase) convertValueFromDB(fi *models.FieldInfo, val interface{}, tz *time.Location) (interface{}, error) {
	if val == nil {
		return nil, nil
	}

	var value interface{}
	var tErr error

	var str *utils.StrTo
	switch v := val.(type) {
	case []byte:
		s := utils.StrTo(string(v))
		str = &s
	case string:
		s := utils.StrTo(v)
		str = &s
	}

	fieldType := fi.FieldType

setValue:
	switch {
	case fieldType == TypeBooleanField:
		if str == nil {
			switch v := val.(type) {
			case int64:
				b := v == 1
				value = b
			default:
				s := utils.StrTo(utils.ToStr(v))
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
	case fieldType == TypeVarCharField || fieldType == TypeCharField || fieldType == TypeTextField || fieldType == TypeJSONField || fieldType == TypeJsonbField:
		if str == nil {
			value = utils.ToStr(val)
		} else {
			value = str.String()
		}
	case fieldType == TypeTimeField || fieldType == TypeDateField || fieldType == TypeDateTimeField:
		if str == nil {
			switch t := val.(type) {
			case time.Time:
				d.ins.TimeFromDB(&t, tz)
				value = t
			default:
				s := utils.StrTo(utils.ToStr(t))
				str = &s
			}
		}
		if str != nil {
			s := str.String()
			var (
				t   time.Time
				err error
			)

			if fi.TimePrecision != nil && len(s) >= (20+*fi.TimePrecision) {
				layout := utils.FormatDateTime + "."
				for i := 0; i < *fi.TimePrecision; i++ {
					layout += "0"
				}
				t, err = time.ParseInLocation(layout, s[:20+*fi.TimePrecision], tz)
			} else if len(s) >= 19 {
				s = s[:19]
				t, err = time.ParseInLocation(utils.FormatDateTime, s, tz)
			} else if len(s) >= 10 {
				if len(s) > 10 {
					s = s[:10]
				}
				t, err = time.ParseInLocation(utils.FormatDate, s, tz)
			} else if len(s) >= 8 {
				if len(s) > 8 {
					s = s[:8]
				}
				t, err = time.ParseInLocation(utils.FormatTime, s, tz)
			}
			t = t.In(DefaultTimeLoc)

			if err != nil && s != "00:00:00" && s != "0000-00-00" && s != "0000-00-00 00:00:00" {
				tErr = err
				goto end
			}
			value = t
		}
	case fieldType&IsIntegerField > 0:
		if str == nil {
			s := utils.StrTo(utils.ToStr(val))
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
			if fieldType&IsPositiveIntegerField > 0 {
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
				s := utils.StrTo(utils.ToStr(v))
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
		fi = fi.RelModelInfo.Fields.Pk
		fieldType = fi.FieldType
		goto setValue
	}

end:
	if tErr != nil {
		err := fmt.Errorf("convert to `%s` failed, field: %s err: %s", fi.AddrValue.Type(), fi.FullName, tErr)
		return nil, err
	}

	return value, nil
}

// set one value to struct column field.
func (d *dbBase) setFieldValue(fi *models.FieldInfo, value interface{}, field reflect.Value) (interface{}, error) {
	fieldType := fi.FieldType
	isNative := !fi.IsFielder

setValue:
	switch {
	case fieldType == TypeBooleanField:
		if isNative {
			if nb, ok := field.Interface().(sql.NullBool); ok {
				if value == nil {
					nb.Valid = false
				} else {
					nb.Bool = value.(bool)
					nb.Valid = true
				}
				field.Set(reflect.ValueOf(nb))
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					v := value.(bool)
					field.Set(reflect.ValueOf(&v))
				}
			} else {
				if value == nil {
					value = false
				}
				field.SetBool(value.(bool))
			}
		}
	case fieldType == TypeVarCharField || fieldType == TypeCharField || fieldType == TypeTextField || fieldType == TypeJSONField || fieldType == TypeJsonbField:
		if isNative {
			if ns, ok := field.Interface().(sql.NullString); ok {
				if value == nil {
					ns.Valid = false
				} else {
					ns.String = value.(string)
					ns.Valid = true
				}
				field.Set(reflect.ValueOf(ns))
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					v := value.(string)
					field.Set(reflect.ValueOf(&v))
				}
			} else {
				if value == nil {
					value = ""
				}
				field.SetString(value.(string))
			}
		}
	case fieldType == TypeTimeField || fieldType == TypeDateField || fieldType == TypeDateTimeField:
		if isNative {
			if value == nil {
				value = time.Time{}
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					v := value.(time.Time)
					field.Set(reflect.ValueOf(&v))
				}
			} else {
				field.Set(reflect.ValueOf(value))
			}
		}
	case fieldType == TypePositiveBitField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := uint8(value.(uint64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypePositiveSmallIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := uint16(value.(uint64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypePositiveIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			if field.Type() == reflect.TypeOf(new(uint)) {
				v := uint(value.(uint64))
				field.Set(reflect.ValueOf(&v))
			} else {
				v := uint32(value.(uint64))
				field.Set(reflect.ValueOf(&v))
			}
		}
	case fieldType == TypePositiveBigIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := value.(uint64)
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypeBitField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := int8(value.(int64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypeSmallIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := int16(value.(int64))
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType == TypeIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			if field.Type() == reflect.TypeOf(new(int)) {
				v := int(value.(int64))
				field.Set(reflect.ValueOf(&v))
			} else {
				v := int32(value.(int64))
				field.Set(reflect.ValueOf(&v))
			}
		}
	case fieldType == TypeBigIntegerField && field.Kind() == reflect.Ptr:
		if value != nil {
			v := value.(int64)
			field.Set(reflect.ValueOf(&v))
		}
	case fieldType&IsIntegerField > 0:
		if fieldType&IsPositiveIntegerField > 0 {
			if isNative {
				if value == nil {
					value = uint64(0)
				}
				field.SetUint(value.(uint64))
			}
		} else {
			if isNative {
				if ni, ok := field.Interface().(sql.NullInt64); ok {
					if value == nil {
						ni.Valid = false
					} else {
						ni.Int64 = value.(int64)
						ni.Valid = true
					}
					field.Set(reflect.ValueOf(ni))
				} else {
					if value == nil {
						value = int64(0)
					}
					field.SetInt(value.(int64))
				}
			}
		}
	case fieldType == TypeFloatField || fieldType == TypeDecimalField:
		if isNative {
			if nf, ok := field.Interface().(sql.NullFloat64); ok {
				if value == nil {
					nf.Valid = false
				} else {
					nf.Float64 = value.(float64)
					nf.Valid = true
				}
				field.Set(reflect.ValueOf(nf))
			} else if field.Kind() == reflect.Ptr {
				if value != nil {
					if field.Type() == reflect.TypeOf(new(float32)) {
						v := float32(value.(float64))
						field.Set(reflect.ValueOf(&v))
					} else {
						v := value.(float64)
						field.Set(reflect.ValueOf(&v))
					}
				}
			} else {

				if value == nil {
					value = float64(0)
				}
				field.SetFloat(value.(float64))
			}
		}
	case fieldType&IsRelField > 0:
		if value != nil {
			fieldType = fi.RelModelInfo.Fields.Pk.FieldType
			mf := reflect.New(fi.RelModelInfo.AddrField.Elem().Type())
			field.Set(mf)
			f := mf.Elem().FieldByIndex(fi.RelModelInfo.Fields.Pk.FieldIndex)
			field = f
			goto setValue
		}
	}

	if !isNative {
		fd := field.Addr().Interface().(models.Fielder)
		err := fd.SetRaw(value)
		if err != nil {
			err = fmt.Errorf("converted value `%v` set to Fielder `%s` failed, err: %s", value, fi.FullName, err)
			return nil, err
		}
	}

	return value, nil
}

// ReadValues query sql, read values , save to *[]ParamList.
func (d *dbBase) ReadValues(ctx context.Context, q dbQuerier, qs *querySet, mi *models.ModelInfo, cond *Condition, exprs []string, container interface{}, tz *time.Location) (int64, error) {
	var (
		maps  []Params
		lists []ParamsList
		list  ParamsList
	)

	typ := 0
	switch v := container.(type) {
	case *[]Params:
		d := *v
		if len(d) == 0 {
			maps = d
		}
		typ = 1
	case *[]ParamsList:
		d := *v
		if len(d) == 0 {
			lists = d
		}
		typ = 2
	case *ParamsList:
		d := *v
		if len(d) == 0 {
			list = d
		}
		typ = 3
	default:
		panic(fmt.Errorf("unsupport read values type `%T`", container))
	}

	tables := newDbTables(mi, d.ins)

	var (
		cols  []string
		infos []*models.FieldInfo
	)

	hasExprs := len(exprs) > 0

	Q := d.ins.TableQuote()

	if hasExprs {
		cols = make([]string, 0, len(exprs))
		infos = make([]*models.FieldInfo, 0, len(exprs))
		for _, ex := range exprs {
			index, name, fi, suc := tables.parseExprs(mi, strings.Split(ex, ExprSep))
			if !suc {
				panic(fmt.Errorf("unknown field/column name `%s`", ex))
			}
			cols = append(cols, fmt.Sprintf("%s.%s%s%s %s%s%s", index, Q, fi.Column, Q, Q, name, Q))
			infos = append(infos, fi)
		}
	} else {
		cols = make([]string, 0, len(mi.Fields.DBcols))
		infos = make([]*models.FieldInfo, 0, len(exprs))
		for _, fi := range mi.Fields.FieldsDB {
			cols = append(cols, fmt.Sprintf("T0.%s%s%s %s%s%s", Q, fi.Column, Q, Q, fi.Name, Q))
			infos = append(infos, fi)
		}
	}

	where, args := tables.getCondSQL(cond, false, tz)
	groupBy := tables.getGroupSQL(qs.groups)
	orderBy := tables.getOrderSQL(qs.orders)
	limit := tables.getLimitSQL(mi, qs.offset, qs.limit)
	join := tables.getJoinSQL()
	specifyIndexes := tables.getIndexSql(mi.Table, qs.useIndex, qs.indexes)

	sels := strings.Join(cols, ", ")

	sqlSelect := "SELECT"
	if qs.distinct {
		sqlSelect += " DISTINCT"
	}
	query := fmt.Sprintf("%s %s FROM %s%s%s T0 %s%s%s%s%s%s",
		sqlSelect, sels,
		Q, mi.Table, Q,
		specifyIndexes, join, where, groupBy, orderBy, limit)

	d.ins.ReplaceMarks(&query)

	rs, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	refs := make([]interface{}, len(cols))
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}

	defer rs.Close()

	var (
		cnt     int64
		columns []string
	)
	for rs.Next() {
		if cnt == 0 {
			cols, err := rs.Columns()
			if err != nil {
				return 0, err
			}
			columns = cols
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
					panic(fmt.Errorf("db value convert failed `%v` %s", val, err.Error()))
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
					panic(fmt.Errorf("db value convert failed `%v` %s", val, err.Error()))
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
					panic(fmt.Errorf("db value convert failed `%v` %s", val, err.Error()))
				}

				list = append(list, value)
			}
		}

		cnt++
	}

	if err = rs.Err(); err != nil {
		return 0, err
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

// SupportUpdateJoin flag of update joined record.
func (d *dbBase) SupportUpdateJoin() bool {
	return true
}

func (d *dbBase) MaxLimit() uint64 {
	return 18446744073709551615
}

// TableQuote return quote.
func (d *dbBase) TableQuote() string {
	return "`"
}

// ReplaceMarks replace value placeholder in parametered sql string.
func (d *dbBase) ReplaceMarks(query *string) {
	// default use `?` as mark, do nothing
}

// flag of RETURNING sql.
func (d *dbBase) HasReturningID(*models.ModelInfo, *string) bool {
	return false
}

// sync auto key
func (d *dbBase) setval(ctx context.Context, db dbQuerier, mi *models.ModelInfo, autoFields []string) error {
	return nil
}

// TimeFromDB convert time from db.
func (d *dbBase) TimeFromDB(t *time.Time, tz *time.Location) {
	*t = t.In(tz)
}

// TimeToDB convert time to db.
func (d *dbBase) TimeToDB(t *time.Time, tz *time.Location) {
	*t = t.In(tz)
}

// DbTypes get database types.
func (d *dbBase) DbTypes() map[string]string {
	return nil
}

// GetTables gt all tables.
func (d *dbBase) GetTables(db dbQuerier) (map[string]bool, error) {
	tables := make(map[string]bool)
	query := d.ins.ShowTablesQuery()
	rows, err := db.Query(query)
	if err != nil {
		return tables, err
	}

	defer rows.Close()

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

	return tables, rows.Err()
}

// GetColumns get all cloumns in table.
func (d *dbBase) GetColumns(ctx context.Context, db dbQuerier, table string) (map[string][3]string, error) {
	columns := make(map[string][3]string)
	query := d.ins.ShowColumnsQuery(table)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return columns, err
	}

	defer rows.Close()

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

	return columns, rows.Err()
}

// not implement.
func (d *dbBase) OperatorSQL(operator string) string {
	panic(ErrNotImplement)
}

// not implement.
func (d *dbBase) ShowTablesQuery() string {
	panic(ErrNotImplement)
}

// not implement.
func (d *dbBase) ShowColumnsQuery(table string) string {
	panic(ErrNotImplement)
}

// not implement.
func (d *dbBase) IndexExists(context.Context, dbQuerier, string, string) bool {
	panic(ErrNotImplement)
}

// GenerateSpecifyIndex return a specifying index clause
func (d *dbBase) GenerateSpecifyIndex(tableName string, useIndex int, indexes []string) string {
	var s []string
	Q := d.TableQuote()
	for _, index := range indexes {
		tmp := fmt.Sprintf(`%s%s%s`, Q, index, Q)
		s = append(s, tmp)
	}

	var useWay string

	switch useIndex {
	case hints.KeyUseIndex:
		useWay = `USE`
	case hints.KeyForceIndex:
		useWay = `FORCE`
	case hints.KeyIgnoreIndex:
		useWay = `IGNORE`
	default:
		logs.DebugLog.Println("[WARN] Not a valid specifying action, so that action is ignored")
		return ``
	}

	return fmt.Sprintf(` %s INDEX(%s) `, useWay, strings.Join(s, `,`))
}
