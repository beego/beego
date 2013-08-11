package orm

import (
	"database/sql"
	"fmt"
	"reflect"
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

	args := getFlatParams(nil, o.args)
	return o.orm.db.Exec(query, args...)
}

func (o *rawSet) QueryRow(...interface{}) error {
	//TODO
	return nil
}

func (o *rawSet) QueryRows(...interface{}) (int64, error) {
	//TODO
	return 0, nil
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
		panic(fmt.Sprintf("unsupport read values type `%T`", container))
	}

	query := o.query
	o.orm.alias.DbBaser.ReplaceMarks(&query)

	args := getFlatParams(nil, o.args)

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
				params[cols[i]] = value.String
			}
			maps = append(maps, params)
		case 2:
			params := make(ParamsList, 0, len(cols))
			for _, ref := range refs {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				params = append(params, value.String)
			}
			lists = append(lists, params)
		case 3:
			for _, ref := range refs {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				list = append(list, value.String)
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
