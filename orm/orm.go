package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTXHasBegin    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTXNotBegin    = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter.One> return multi rows")
	ErrStmtClosed    = errors.New("<QuerySeter.Insert> stmt already closed")
	DefaultRowsLimit = 1000
	DefaultRelsDepth = 5
	DefaultTimeLoc   = time.Local
)

type Params map[string]interface{}
type ParamsList []interface{}

type orm struct {
	alias *alias
	db    dbQuerier
	isTx  bool
}

func (o *orm) Object(md Modeler) ObjectSeter {
	md.Init(md, true)
	name := md.GetTableName()
	if mi, ok := modelCache.get(name); ok {
		return newObject(o, mi, md)
	}
	panic(fmt.Sprintf("<orm.Object> table name: `%s` not exists", name))
}

func (o *orm) QueryTable(ptrStructOrTableName interface{}) QuerySeter {
	name := ""
	if table, ok := ptrStructOrTableName.(string); ok {
		name = snakeString(table)
	} else if md, ok := ptrStructOrTableName.(Modeler); ok {
		md.Init(md, true)
		name = md.GetTableName()
	}
	if mi, ok := modelCache.get(name); ok {
		return newQuerySet(o, mi)
	}
	panic(fmt.Sprintf("<orm.SetTable> table name: `%s` not exists", name))
}

func (o *orm) Using(name string) error {
	if o.isTx {
		panic("<orm.Using> transaction has been start, cannot change db")
	}
	if al, ok := dataBaseCache.get(name); ok {
		o.alias = al
		o.db = al.DB
	} else {
		return errors.New(fmt.Sprintf("<orm.Using> unknown db alias name `%s`", name))
	}
	return nil
}

func (o *orm) Begin() error {
	if o.isTx {
		return ErrTXHasBegin
	}
	tx, err := o.alias.DB.Begin()
	if err != nil {
		return err
	}
	o.isTx = true
	o.db = tx
	return nil
}

func (o *orm) Commit() error {
	if o.isTx == false {
		return ErrTXNotBegin
	}
	err := o.db.(*sql.Tx).Commit()
	if err == nil {
		o.isTx = false
		o.db = o.alias.DB
	}
	return err
}

func (o *orm) Rollback() error {
	if o.isTx == false {
		return ErrTXNotBegin
	}
	err := o.db.(*sql.Tx).Rollback()
	if err == nil {
		o.isTx = false
		o.db = o.alias.DB
	}
	return err
}

func (o *orm) Raw(query string, args ...interface{}) RawSeter {
	return newRawSet(o, query, args)
}

func NewOrm() Ormer {
	o := new(orm)
	err := o.Using("default")
	if err != nil {
		panic(err)
	}
	return o
}
