package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"
)

const (
	Debug_Queries = iota
)

var (
	// DebugLevel       = Debug_Queries
	Debug            = false
	DebugLog         = NewLog(os.Stderr)
	DefaultRowsLimit = 1000
	DefaultRelsDepth = 5
	DefaultTimeLoc   = time.Local
	ErrTxHasBegan    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTxDone        = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> no row found")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrNotImplement  = errors.New("have not implement")
)

type Params map[string]interface{}
type ParamsList []interface{}

type orm struct {
	alias *alias
	db    dbQuerier
	isTx  bool
}

var _ Ormer = new(orm)

func (o *orm) getMiInd(md Modeler) (mi *modelInfo, ind reflect.Value) {
	md.Init(md, true)
	name := md.GetTableName()
	if mi, ok := modelCache.get(name); ok {
		return mi, reflect.Indirect(reflect.ValueOf(md))
	}
	panic(fmt.Sprintf("<orm> table name: `%s` not exists", name))
}

func (o *orm) Read(md Modeler) error {
	mi, ind := o.getMiInd(md)
	err := o.alias.DbBaser.Read(o.db, mi, ind)
	if err != nil {
		return err
	}
	return nil
}

func (o *orm) Insert(md Modeler) (int64, error) {
	mi, ind := o.getMiInd(md)
	id, err := o.alias.DbBaser.Insert(o.db, mi, ind)
	if err != nil {
		return id, err
	}
	if id > 0 {
		if mi.fields.pk.auto {
			ind.Field(mi.fields.pk.fieldIndex).SetInt(id)
		}
	}
	return id, nil
}

func (o *orm) Update(md Modeler) (int64, error) {
	mi, ind := o.getMiInd(md)
	num, err := o.alias.DbBaser.Update(o.db, mi, ind)
	if err != nil {
		return num, err
	}
	return num, nil
}

func (o *orm) Delete(md Modeler) (int64, error) {
	mi, ind := o.getMiInd(md)
	num, err := o.alias.DbBaser.Delete(o.db, mi, ind)
	if err != nil {
		return num, err
	}
	if num > 0 {
		if mi.fields.pk.auto {
			ind.Field(mi.fields.pk.fieldIndex).SetInt(0)
		}
	}
	return num, nil
}

func (o *orm) M2mAdd(md Modeler, name string, mds ...interface{}) (int64, error) {
	// TODO
	panic(ErrNotImplement)
	return 0, nil
}

func (o *orm) M2mDel(md Modeler, name string, mds ...interface{}) (int64, error) {
	// TODO
	panic(ErrNotImplement)
	return 0, nil
}

func (o *orm) LoadRel(md Modeler, name string) (int64, error) {
	// TODO
	panic(ErrNotImplement)
	return 0, nil
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
		if Debug {
			o.db = newDbQueryLog(al, al.DB)
		} else {
			o.db = al.DB
		}
	} else {
		return errors.New(fmt.Sprintf("<orm.Using> unknown db alias name `%s`", name))
	}
	return nil
}

func (o *orm) Begin() error {
	if o.isTx {
		return ErrTxHasBegan
	}
	var tx *sql.Tx
	tx, err := o.db.(txer).Begin()
	if err != nil {
		return err
	}
	o.isTx = true
	if Debug {
		o.db.(*dbQueryLog).SetDB(tx)
	} else {
		o.db = tx
	}
	return nil
}

func (o *orm) Commit() error {
	if o.isTx == false {
		return ErrTxDone
	}
	err := o.db.(txEnder).Commit()
	if err == nil {
		o.isTx = false
		o.Using(o.alias.Name)
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}

func (o *orm) Rollback() error {
	if o.isTx == false {
		return ErrTxDone
	}
	err := o.db.(txEnder).Rollback()
	if err == nil {
		o.isTx = false
		o.Using(o.alias.Name)
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return err
}

func (o *orm) Raw(query string, args ...interface{}) RawSeter {
	return newRawSet(o, query, args)
}

func (o *orm) Driver() Driver {
	return driver(o.alias.Name)
}

func NewOrm() Ormer {
	BootStrap() // execute only once

	o := new(orm)
	err := o.Using("default")
	if err != nil {
		panic(err)
	}
	return o
}
