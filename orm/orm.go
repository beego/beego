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
	DefaultRelsDepth = 2
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

func (o *orm) getMiInd(md interface{}) (mi *modelInfo, ind reflect.Value) {
	val := reflect.ValueOf(md)
	ind = reflect.Indirect(val)
	typ := ind.Type()
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Ormer> cannot use non-ptr model struct `%s`", getFullName(typ)))
	}
	name := getFullName(typ)
	if mi, ok := modelCache.getByFN(name); ok {
		return mi, ind
	}
	panic(fmt.Errorf("<Ormer> table: `%s` not found, maybe not RegisterModel", name))
}

func (o *orm) getFieldInfo(mi *modelInfo, name string) *fieldInfo {
	fi, ok := mi.fields.GetByAny(name)
	if !ok {
		panic(fmt.Errorf("<Ormer> cannot find field `%s` for model `%s`", name, mi.fullName))
	}
	return fi
}

func (o *orm) Read(md interface{}, cols ...string) error {
	mi, ind := o.getMiInd(md)
	err := o.alias.DbBaser.Read(o.db, mi, ind, o.alias.TZ, cols)
	if err != nil {
		return err
	}
	return nil
}

func (o *orm) Insert(md interface{}) (int64, error) {
	mi, ind := o.getMiInd(md)
	id, err := o.alias.DbBaser.Insert(o.db, mi, ind, o.alias.TZ)
	if err != nil {
		return id, err
	}
	if id > 0 {
		if mi.fields.pk.auto {
			if mi.fields.pk.fieldType&IsPostiveIntegerField > 0 {
				ind.Field(mi.fields.pk.fieldIndex).SetUint(uint64(id))
			} else {
				ind.Field(mi.fields.pk.fieldIndex).SetInt(id)
			}
		}
	}
	return id, nil
}

func (o *orm) Update(md interface{}, cols ...string) (int64, error) {
	mi, ind := o.getMiInd(md)
	num, err := o.alias.DbBaser.Update(o.db, mi, ind, o.alias.TZ, cols)
	if err != nil {
		return num, err
	}
	return num, nil
}

func (o *orm) Delete(md interface{}) (int64, error) {
	mi, ind := o.getMiInd(md)
	num, err := o.alias.DbBaser.Delete(o.db, mi, ind, o.alias.TZ)
	if err != nil {
		return num, err
	}
	if num > 0 {
		if mi.fields.pk.auto {
			if mi.fields.pk.fieldType&IsPostiveIntegerField > 0 {
				ind.Field(mi.fields.pk.fieldIndex).SetUint(0)
			} else {
				ind.Field(mi.fields.pk.fieldIndex).SetInt(0)
			}
		}
	}
	return num, nil
}

func (o *orm) QueryM2M(md interface{}, name string) QueryM2Mer {
	mi, ind := o.getMiInd(md)
	fi := o.getFieldInfo(mi, name)

	switch {
	case fi.fieldType == RelManyToMany:
	case fi.fieldType == RelReverseMany && fi.reverseFieldInfo.mi.isThrough:
	default:
		panic(fmt.Errorf("<Ormer.QueryM2M> model `%s` . name `%s` is not a m2m field", fi.name, mi.fullName))
	}

	return newQueryM2M(md, o, mi, fi, ind)
}

func (o *orm) LoadRelated(md interface{}, name string, args ...interface{}) (int64, error) {
	_, fi, ind, qseter := o.queryRelated(md, name)

	qs := qseter.(*querySet)

	var relDepth int
	var limit, offset int64
	var order string
	for i, arg := range args {
		switch i {
		case 0:
			if v, ok := arg.(bool); ok {
				if v {
					relDepth = DefaultRelsDepth
				}
			} else if v, ok := arg.(int); ok {
				relDepth = v
			}
		case 1:
			limit = ToInt64(arg)
		case 2:
			offset = ToInt64(arg)
		case 3:
			order, _ = arg.(string)
		}
	}

	switch fi.fieldType {
	case RelOneToOne, RelForeignKey, RelReverseOne:
		limit = 1
		offset = 0
	}

	qs.limit = limit
	qs.offset = offset
	qs.relDepth = relDepth

	if len(order) > 0 {
		qs.orders = []string{order}
	}

	find := ind.Field(fi.fieldIndex)

	var nums int64
	var err error
	switch fi.fieldType {
	case RelOneToOne, RelForeignKey, RelReverseOne:
		val := reflect.New(find.Type().Elem())
		container := val.Interface()
		err = qs.One(container)
		if err == nil {
			find.Set(val)
			nums = 1
		}
	default:
		nums, err = qs.All(find.Addr().Interface())
	}

	return nums, err
}

func (o *orm) QueryRelated(md interface{}, name string) QuerySeter {
	// is this api needed ?
	_, _, _, qs := o.queryRelated(md, name)
	return qs
}

func (o *orm) queryRelated(md interface{}, name string) (*modelInfo, *fieldInfo, reflect.Value, QuerySeter) {
	mi, ind := o.getMiInd(md)
	fi := o.getFieldInfo(mi, name)

	_, _, exist := getExistPk(mi, ind)
	if exist == false {
		panic(ErrMissPK)
	}

	var qs *querySet

	switch fi.fieldType {
	case RelOneToOne, RelForeignKey, RelManyToMany:
		if !fi.inModel {
			break
		}
		qs = o.getRelQs(md, mi, fi)
	case RelReverseOne, RelReverseMany:
		if !fi.inModel {
			break
		}
		qs = o.getReverseQs(md, mi, fi)
	}

	if qs == nil {
		panic(fmt.Errorf("<Ormer> name `%s` for model `%s` is not an available rel/reverse field"))
	}

	return mi, fi, ind, qs
}

func (o *orm) getReverseQs(md interface{}, mi *modelInfo, fi *fieldInfo) *querySet {
	switch fi.fieldType {
	case RelReverseOne, RelReverseMany:
	default:
		panic(fmt.Errorf("<Ormer> name `%s` for model `%s` is not an available reverse field", fi.name, mi.fullName))
	}

	var q *querySet

	if fi.fieldType == RelReverseMany && fi.reverseFieldInfo.mi.isThrough {
		q = newQuerySet(o, fi.relModelInfo).(*querySet)
		q.cond = NewCondition().And(fi.reverseFieldInfoM2M.column+ExprSep+fi.reverseFieldInfo.column, md)
	} else {
		q = newQuerySet(o, fi.reverseFieldInfo.mi).(*querySet)
		q.cond = NewCondition().And(fi.reverseFieldInfo.column, md)
	}

	return q
}

func (o *orm) getRelQs(md interface{}, mi *modelInfo, fi *fieldInfo) *querySet {
	switch fi.fieldType {
	case RelOneToOne, RelForeignKey, RelManyToMany:
	default:
		panic(fmt.Errorf("<Ormer> name `%s` for model `%s` is not an available rel field", fi.name, mi.fullName))
	}

	q := newQuerySet(o, fi.relModelInfo).(*querySet)
	q.cond = NewCondition()

	if fi.fieldType == RelManyToMany {
		q.cond = q.cond.And(fi.reverseFieldInfoM2M.column+ExprSep+fi.reverseFieldInfo.column, md)
	} else {
		q.cond = q.cond.And(fi.reverseFieldInfo.column, md)
	}

	return q
}

func (o *orm) QueryTable(ptrStructOrTableName interface{}) (qs QuerySeter) {
	name := ""
	if table, ok := ptrStructOrTableName.(string); ok {
		name = snakeString(table)
		if mi, ok := modelCache.get(name); ok {
			qs = newQuerySet(o, mi)
		}
	} else {
		name = getFullName(indirectType(reflect.TypeOf(ptrStructOrTableName)))
		if mi, ok := modelCache.getByFN(name); ok {
			qs = newQuerySet(o, mi)
		}
	}
	if qs == nil {
		panic(fmt.Errorf("<Ormer.QueryTable> table name: `%s` not exists", name))
	}
	return
}

func (o *orm) Using(name string) error {
	if o.isTx {
		panic(fmt.Errorf("<Ormer.Using> transaction has been start, cannot change db"))
	}
	if al, ok := dataBaseCache.get(name); ok {
		o.alias = al
		if Debug {
			o.db = newDbQueryLog(al, al.DB)
		} else {
			o.db = al.DB
		}
	} else {
		return fmt.Errorf("<Ormer.Using> unknown db alias name `%s`", name)
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
