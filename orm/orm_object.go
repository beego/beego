package orm

import (
	"database/sql"
	"fmt"
	"reflect"
)

type insertSet struct {
	mi     *modelInfo
	orm    *orm
	stmt   *sql.Stmt
	closed bool
}

func (o *insertSet) Insert(md Modeler) (int64, error) {
	if o.closed {
		return 0, ErrStmtClosed
	}
	md.Init(md, true)
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	if val.Type() != o.mi.addrField.Type() {
		panic(fmt.Sprintf("<Inserter.Insert> need type `%s` but found `%s`", o.mi.addrField.Type(), val.Type()))
	}
	id, err := o.orm.alias.DbBaser.InsertStmt(o.stmt, o.mi, ind)
	if err != nil {
		return id, err
	}
	if id > 0 {
		if o.mi.fields.auto != nil {
			ind.Field(o.mi.fields.auto.fieldIndex).SetInt(id)
		}
	}
	return id, nil
}

func (o *insertSet) Close() error {
	o.closed = true
	return o.stmt.Close()
}

func newInsertSet(orm *orm, mi *modelInfo) (Inserter, error) {
	bi := new(insertSet)
	bi.orm = orm
	bi.mi = mi
	st, err := orm.alias.DbBaser.PrepareInsert(orm.db, mi)
	if err != nil {
		return nil, err
	}
	bi.stmt = st
	return bi, nil
}
