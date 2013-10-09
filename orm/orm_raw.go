package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type rawPrepare struct {
	rs     *rawSet
	stmt   stmtQuerier
	closed bool
}

func (o *rawPrepare) Exec(args ...interface{}) (sql.Result, error) {
	if o.closed {
		return nil, ErrStmtClosed
	}
	return o.stmt.Exec(args...)
}

func (o *rawPrepare) Close() error {
	o.closed = true
	return o.stmt.Close()
}

func newRawPreparer(rs *rawSet) (RawPreparer, error) {
	o := new(rawPrepare)
	o.rs = rs

	query := rs.query
	rs.orm.alias.DbBaser.ReplaceMarks(&query)

	st, err := rs.orm.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	if Debug {
		o.stmt = newStmtQueryLog(rs.orm.alias, st, query)
	} else {
		o.stmt = st
	}
	return o, nil
}

type rawSet struct {
	query string
	args  []interface{}
	orm   *orm
}

var _ RawSeter = new(rawSet)

func (o rawSet) SetArgs(args ...interface{}) RawSeter {
	o.args = args
	return &o
}

func (o *rawSet) Exec() (sql.Result, error) {
	query := o.query
	o.orm.alias.DbBaser.ReplaceMarks(&query)

	args := getFlatParams(nil, o.args, o.orm.alias.TZ)
	return o.orm.db.Exec(query, args...)
}

func (o *rawSet) setFieldValue(ind reflect.Value, value interface{}) {
	switch ind.Kind() {
	case reflect.Bool:
		if value == nil {
			ind.SetBool(false)
		} else if v, ok := value.(bool); ok {
			ind.SetBool(v)
		} else {
			v, _ := StrTo(ToStr(value)).Bool()
			ind.SetBool(v)
		}

	case reflect.String:
		if value == nil {
			ind.SetString("")
		} else {
			ind.SetString(ToStr(value))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			ind.SetInt(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ind.SetInt(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ind.SetInt(int64(val.Uint()))
			default:
				v, _ := StrTo(ToStr(value)).Int64()
				ind.SetInt(v)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			ind.SetUint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ind.SetUint(uint64(val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ind.SetUint(val.Uint())
			default:
				v, _ := StrTo(ToStr(value)).Uint64()
				ind.SetUint(v)
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			ind.SetFloat(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				ind.SetFloat(val.Float())
			default:
				v, _ := StrTo(ToStr(value)).Float64()
				ind.SetFloat(v)
			}
		}

	case reflect.Struct:
		if value == nil {
			ind.Set(reflect.Zero(ind.Type()))

		} else if _, ok := ind.Interface().(time.Time); ok {
			var str string
			switch d := value.(type) {
			case time.Time:
				o.orm.alias.DbBaser.TimeFromDB(&d, o.orm.alias.TZ)
				ind.Set(reflect.ValueOf(d))
			case []byte:
				str = string(d)
			case string:
				str = d
			}
			if str != "" {
				if len(str) >= 19 {
					str = str[:19]
					t, err := time.ParseInLocation(format_DateTime, str, o.orm.alias.TZ)
					if err == nil {
						t = t.In(DefaultTimeLoc)
						ind.Set(reflect.ValueOf(t))
					}
				} else if len(str) >= 10 {
					str = str[:10]
					t, err := time.ParseInLocation(format_Date, str, DefaultTimeLoc)
					if err == nil {
						ind.Set(reflect.ValueOf(t))
					}
				}
			}
		}
	}
}

func (o *rawSet) loopInitRefs(typ reflect.Type, refsPtr *[]interface{}, sIdxesPtr *[][]int) {
	sIdxes := *sIdxesPtr
	refs := *refsPtr

	if typ.Kind() == reflect.Struct {
		if typ.String() == "time.Time" {
			var ref interface{}
			refs = append(refs, &ref)
			sIdxes = append(sIdxes, []int{0})
		} else {
			idxs := []int{}
		outFor:
			for idx := 0; idx < typ.NumField(); idx++ {
				ctyp := typ.Field(idx)

				tag := ctyp.Tag.Get(defaultStructTagName)
				for _, v := range strings.Split(tag, defaultStructTagDelim) {
					if v == "-" {
						continue outFor
					}
				}

				tp := ctyp.Type
				if tp.Kind() == reflect.Ptr {
					tp = tp.Elem()
				}

				if tp.String() == "time.Time" {
					var ref interface{}
					refs = append(refs, &ref)

				} else if tp.Kind() != reflect.Struct {
					var ref interface{}
					refs = append(refs, &ref)

				} else {
					// skip other type
					continue
				}

				idxs = append(idxs, idx)
			}
			sIdxes = append(sIdxes, idxs)
		}
	} else {
		var ref interface{}
		refs = append(refs, &ref)
		sIdxes = append(sIdxes, []int{0})
	}

	*sIdxesPtr = sIdxes
	*refsPtr = refs
}

func (o *rawSet) loopSetRefs(refs []interface{}, sIdxes [][]int, sInds []reflect.Value, nIndsPtr *[]reflect.Value, eTyps []reflect.Type, init bool) {
	nInds := *nIndsPtr

	cur := 0
	for i, idxs := range sIdxes {
		sInd := sInds[i]
		eTyp := eTyps[i]

		typ := eTyp
		isPtr := false
		if typ.Kind() == reflect.Ptr {
			isPtr = true
			typ = typ.Elem()
		}
		if typ.Kind() == reflect.Ptr {
			isPtr = true
			typ = typ.Elem()
		}

		var nInd reflect.Value
		if init {
			nInd = reflect.New(sInd.Type()).Elem()
		} else {
			nInd = nInds[i]
		}

		val := reflect.New(typ)
		ind := val.Elem()

		tpName := ind.Type().String()

		if ind.Kind() == reflect.Struct {
			if tpName == "time.Time" {
				value := reflect.ValueOf(refs[cur]).Elem().Interface()
				if isPtr && value == nil {
					val = reflect.New(val.Type()).Elem()
				} else {
					o.setFieldValue(ind, value)
				}
				cur++
			} else {
				hasValue := false
				for _, idx := range idxs {
					tind := ind.Field(idx)
					value := reflect.ValueOf(refs[cur]).Elem().Interface()
					if value != nil {
						hasValue = true
					}
					if tind.Kind() == reflect.Ptr {
						if value == nil {
							tindV := reflect.New(tind.Type()).Elem()
							tind.Set(tindV)
						} else {
							tindV := reflect.New(tind.Type().Elem())
							o.setFieldValue(tindV.Elem(), value)
							tind.Set(tindV)
						}
					} else {
						o.setFieldValue(tind, value)
					}
					cur++
				}
				if hasValue == false && isPtr {
					val = reflect.New(val.Type()).Elem()
				}
			}
		} else {
			value := reflect.ValueOf(refs[cur]).Elem().Interface()
			if isPtr && value == nil {
				val = reflect.New(val.Type()).Elem()
			} else {
				o.setFieldValue(ind, value)
			}
			cur++
		}

		if nInd.Kind() == reflect.Slice {
			if isPtr {
				nInd = reflect.Append(nInd, val)
			} else {
				nInd = reflect.Append(nInd, ind)
			}
		} else {
			if isPtr {
				nInd.Set(val)
			} else {
				nInd.Set(ind)
			}
		}

		nInds[i] = nInd
	}
}

func (o *rawSet) QueryRow(containers ...interface{}) error {
	if len(containers) == 0 {
		panic(fmt.Errorf("<RawSeter.QueryRow> need at least one arg"))
	}

	refs := make([]interface{}, 0, len(containers))
	sIdxes := make([][]int, 0)
	sInds := make([]reflect.Value, 0)
	eTyps := make([]reflect.Type, 0)

	for _, container := range containers {
		val := reflect.ValueOf(container)
		ind := reflect.Indirect(val)

		if val.Kind() != reflect.Ptr {
			panic(fmt.Errorf("<RawSeter.QueryRow> all args must be use ptr"))
		}

		etyp := ind.Type()
		typ := etyp
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		sInds = append(sInds, ind)
		eTyps = append(eTyps, etyp)

		o.loopInitRefs(typ, &refs, &sIdxes)
	}

	query := o.query
	o.orm.alias.DbBaser.ReplaceMarks(&query)

	args := getFlatParams(nil, o.args, o.orm.alias.TZ)
	row := o.orm.db.QueryRow(query, args...)

	if err := row.Scan(refs...); err == sql.ErrNoRows {
		return ErrNoRows
	} else if err != nil {
		return err
	}

	nInds := make([]reflect.Value, len(sInds))
	o.loopSetRefs(refs, sIdxes, sInds, &nInds, eTyps, true)
	for i, sInd := range sInds {
		nInd := nInds[i]
		sInd.Set(nInd)
	}

	return nil
}

func (o *rawSet) QueryRows(containers ...interface{}) (int64, error) {
	refs := make([]interface{}, 0)
	sIdxes := make([][]int, 0)
	sInds := make([]reflect.Value, 0)
	eTyps := make([]reflect.Type, 0)

	for _, container := range containers {
		val := reflect.ValueOf(container)
		sInd := reflect.Indirect(val)
		if val.Kind() != reflect.Ptr || sInd.Kind() != reflect.Slice {
			panic(fmt.Errorf("<RawSeter.QueryRows> all args must be use ptr slice"))
		}

		etyp := sInd.Type().Elem()
		typ := etyp
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		sInds = append(sInds, sInd)
		eTyps = append(eTyps, etyp)

		o.loopInitRefs(typ, &refs, &sIdxes)
	}

	query := o.query
	o.orm.alias.DbBaser.ReplaceMarks(&query)

	args := getFlatParams(nil, o.args, o.orm.alias.TZ)
	rows, err := o.orm.db.Query(query, args...)
	if err != nil {
		return 0, err
	}

	nInds := make([]reflect.Value, len(sInds))

	var cnt int64
	for rows.Next() {
		if err := rows.Scan(refs...); err != nil {
			return 0, err
		}

		o.loopSetRefs(refs, sIdxes, sInds, &nInds, eTyps, cnt == 0)

		cnt++
	}

	if cnt > 0 {
		for i, sInd := range sInds {
			nInd := nInds[i]
			sInd.Set(nInd)
		}
	}

	return cnt, nil
}

func (o *rawSet) readValues(container interface{}) (int64, error) {
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
		panic(fmt.Errorf("<RawSeter> unsupport read values type `%T`", container))
	}

	query := o.query
	o.orm.alias.DbBaser.ReplaceMarks(&query)

	args := getFlatParams(nil, o.args, o.orm.alias.TZ)

	var rs *sql.Rows
	if r, err := o.orm.db.Query(query, args...); err != nil {
		return 0, err
	} else {
		rs = r
	}

	var (
		refs []interface{}
		cnt  int64
		cols []string
	)
	for rs.Next() {
		if cnt == 0 {
			if columns, err := rs.Columns(); err != nil {
				return 0, err
			} else {
				cols = columns
				refs = make([]interface{}, len(cols))
				for i, _ := range refs {
					var ref sql.NullString
					refs[i] = &ref
				}
			}
		}

		if err := rs.Scan(refs...); err != nil {
			return 0, err
		}

		switch typ {
		case 1:
			params := make(Params, len(cols))
			for i, ref := range refs {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params[cols[i]] = value.String
				} else {
					params[cols[i]] = nil
				}
			}
			maps = append(maps, params)
		case 2:
			params := make(ParamsList, 0, len(cols))
			for _, ref := range refs {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params = append(params, value.String)
				} else {
					params = append(params, nil)
				}
			}
			lists = append(lists, params)
		case 3:
			for _, ref := range refs {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					list = append(list, value.String)
				} else {
					list = append(list, nil)
				}
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

func (o *rawSet) Values(container *[]Params) (int64, error) {
	return o.readValues(container)
}

func (o *rawSet) ValuesList(container *[]ParamsList) (int64, error) {
	return o.readValues(container)
}

func (o *rawSet) ValuesFlat(container *ParamsList) (int64, error) {
	return o.readValues(container)
}

func (o *rawSet) Prepare() (RawPreparer, error) {
	return newRawPreparer(o)
}

func newRawSet(orm *orm, query string, args []interface{}) RawSeter {
	o := new(rawSet)
	o.query = query
	o.args = args
	o.orm = orm
	return o
}
