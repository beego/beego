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

// Package orm provide ORM for MySQL/PostgreSQL/sqlite
// Simple Usage
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/beego/beego/v2/client/orm"
//		_ "github.com/go-sql-driver/mysql" // import your used driver
//	)
//
//	// Model Struct
//	type User struct {
//		Id   int    `orm:"auto"`
//		Name string `orm:"size(100)"`
//	}
//
//	func init() {
//		orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
//	}
//
//	func main() {
//		o := orm.NewOrm()
//		user := User{Name: "slene"}
//		// insert
//		id, err := o.Insert(&user)
//		// update
//		user.Name = "astaxie"
//		num, err := o.Update(&user)
//		// read one
//		u := User{Id: user.Id}
//		err = o.Read(&u)
//		// delete
//		num, err = o.Delete(&u)
//	}
package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	ilogs "github.com/beego/beego/v2/client/orm/internal/logs"
	iutils "github.com/beego/beego/v2/client/orm/internal/utils"

	"github.com/beego/beego/v2/client/orm/internal/models"

	"github.com/beego/beego/v2/client/orm/clauses/order_clause"
	"github.com/beego/beego/v2/client/orm/hints"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
)

// DebugQueries define the debug
const (
	DebugQueries = iota
)

// Define common vars
var (
	Debug            = false
	DebugLog         = ilogs.DebugLog
	DefaultRowsLimit = -1
	DefaultRelsDepth = 2
	DefaultTimeLoc   = iutils.DefaultTimeLoc
	ErrTxDone        = errors.New("<TxOrmer.Commit/Rollback> transaction already done")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> no row found")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrArgs          = errors.New("<Ormer> args error may be empty")
	ErrNotImplement  = errors.New("have not implement")

	ErrLastInsertIdUnavailable = errors.New("<Ormer> last insert id is unavailable")
)

// Params stores the Params
type Params map[string]interface{}

// ParamsList stores paramslist
type ParamsList []interface{}

type ormBase struct {
	alias *alias
	db    dbQuerier
}

var (
	_ DQL          = new(ormBase)
	_ DML          = new(ormBase)
	_ DriverGetter = new(ormBase)
)

// Get model info and model reflect value
func (*ormBase) getMi(md interface{}) (mi *models.ModelInfo) {
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	mi = getTypeMi(typ)
	return
}

// Get need ptr model info and model reflect value
func (*ormBase) getPtrMiInd(md interface{}) (mi *models.ModelInfo, ind reflect.Value) {
	val := reflect.ValueOf(md)
	ind = reflect.Indirect(val)
	typ := ind.Type()
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Ormer> cannot use non-ptr model struct `%s`", models.GetFullName(typ)))
	}
	mi = getTypeMi(typ)
	return
}

func getTypeMi(mdTyp reflect.Type) *models.ModelInfo {
	name := models.GetFullName(mdTyp)
	if mi, ok := defaultModelCache.GetByFullName(name); ok {
		return mi
	}
	panic(fmt.Errorf("<Ormer> table: `%s` not found, make sure it was registered with `RegisterModel()`", name))
}

// Get field info from model info by given field name
func (*ormBase) getFieldInfo(mi *models.ModelInfo, name string) *models.FieldInfo {
	fi, ok := mi.Fields.GetByAny(name)
	if !ok {
		panic(fmt.Errorf("<Ormer> cannot find field `%s` for model `%s`", name, mi.FullName))
	}
	return fi
}

// read data to model
func (o *ormBase) Read(md interface{}, cols ...string) error {
	return o.ReadWithCtx(context.Background(), md, cols...)
}

func (o *ormBase) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	mi, ind := o.getPtrMiInd(md)
	return o.alias.DbBaser.Read(ctx, o.db, mi, ind, o.alias.TZ, cols, false)
}

// read data to model, like Read(), but use "SELECT FOR UPDATE" form
func (o *ormBase) ReadForUpdate(md interface{}, cols ...string) error {
	return o.ReadForUpdateWithCtx(context.Background(), md, cols...)
}

func (o *ormBase) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	mi, ind := o.getPtrMiInd(md)
	return o.alias.DbBaser.Read(ctx, o.db, mi, ind, o.alias.TZ, cols, true)
}

// Try to read a row from the database, or insert one if it doesn't exist
func (o *ormBase) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	return o.ReadOrCreateWithCtx(context.Background(), md, col1, cols...)
}

func (o *ormBase) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	cols = append([]string{col1}, cols...)
	mi, ind := o.getPtrMiInd(md)
	err := o.alias.DbBaser.Read(ctx, o.db, mi, ind, o.alias.TZ, cols, false)
	if err == ErrNoRows {
		// Create
		id, err := o.InsertWithCtx(ctx, md)
		return err == nil, id, err
	}

	id, vid := int64(0), ind.FieldByIndex(mi.Fields.Pk.FieldIndex)
	if mi.Fields.Pk.FieldType&IsPositiveIntegerField > 0 {
		id = int64(vid.Uint())
	} else if mi.Fields.Pk.Rel {
		return o.ReadOrCreateWithCtx(ctx, vid.Interface(), mi.Fields.Pk.RelModelInfo.Fields.Pk.Name)
	} else {
		id = vid.Int()
	}

	return false, id, err
}

// insert model data to database
func (o *ormBase) Insert(md interface{}) (int64, error) {
	return o.InsertWithCtx(context.Background(), md)
}

func (o *ormBase) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	mi, ind := o.getPtrMiInd(md)
	id, err := o.alias.DbBaser.Insert(ctx, o.db, mi, ind, o.alias.TZ)
	if err != nil {
		return id, err
	}

	o.setPk(mi, ind, id)

	return id, nil
}

// Set auto pk field
func (*ormBase) setPk(mi *models.ModelInfo, ind reflect.Value, id int64) {
	if mi.Fields.Pk != nil && mi.Fields.Pk.Auto {
		if mi.Fields.Pk.FieldType&IsPositiveIntegerField > 0 {
			ind.FieldByIndex(mi.Fields.Pk.FieldIndex).SetUint(uint64(id))
		} else {
			ind.FieldByIndex(mi.Fields.Pk.FieldIndex).SetInt(id)
		}
	}
}

// insert some models to database
func (o *ormBase) InsertMulti(bulk int, mds interface{}) (int64, error) {
	return o.InsertMultiWithCtx(context.Background(), bulk, mds)
}

func (o *ormBase) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	var cnt int64

	sind := reflect.Indirect(reflect.ValueOf(mds))

	switch sind.Kind() {
	case reflect.Array, reflect.Slice:
		if sind.Len() == 0 {
			return cnt, ErrArgs
		}
	default:
		return cnt, ErrArgs
	}

	if bulk <= 1 {
		for i := 0; i < sind.Len(); i++ {
			ind := reflect.Indirect(sind.Index(i))
			mi := o.getMi(ind.Interface())
			id, err := o.alias.DbBaser.Insert(ctx, o.db, mi, ind, o.alias.TZ)
			if err != nil {
				return cnt, err
			}

			o.setPk(mi, ind, id)

			cnt++
		}
	} else {
		mi := o.getMi(sind.Index(0).Interface())
		return o.alias.DbBaser.InsertMulti(ctx, o.db, mi, sind, bulk, o.alias.TZ)
	}
	return cnt, nil
}

// InsertOrUpdate data to database
func (o *ormBase) InsertOrUpdate(md interface{}, colConflictAndArgs ...string) (int64, error) {
	return o.InsertOrUpdateWithCtx(context.Background(), md, colConflictAndArgs...)
}

func (o *ormBase) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	mi, ind := o.getPtrMiInd(md)
	id, err := o.alias.DbBaser.InsertOrUpdate(ctx, o.db, mi, ind, o.alias, colConflitAndArgs...)
	if err != nil {
		return id, err
	}

	o.setPk(mi, ind, id)

	return id, nil
}

// update model to database.
// cols Set the Columns those want to update.
func (o *ormBase) Update(md interface{}, cols ...string) (int64, error) {
	return o.UpdateWithCtx(context.Background(), md, cols...)
}

func (o *ormBase) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	mi, ind := o.getPtrMiInd(md)
	return o.alias.DbBaser.Update(ctx, o.db, mi, ind, o.alias.TZ, cols)
}

// delete model in database
// cols shows the delete conditions values read from. default is pk
func (o *ormBase) Delete(md interface{}, cols ...string) (int64, error) {
	return o.DeleteWithCtx(context.Background(), md, cols...)
}

func (o *ormBase) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	mi, ind := o.getPtrMiInd(md)
	num, err := o.alias.DbBaser.Delete(ctx, o.db, mi, ind, o.alias.TZ, cols)
	return num, err
}

// create a models to models queryer
func (o *ormBase) QueryM2M(md interface{}, name string) QueryM2Mer {
	mi, ind := o.getPtrMiInd(md)
	fi := o.getFieldInfo(mi, name)

	switch {
	case fi.FieldType == RelManyToMany:
	case fi.FieldType == RelReverseMany && fi.ReverseFieldInfo.Mi.IsThrough:
	default:
		panic(fmt.Errorf("<Ormer.QueryM2M> model `%s` . name `%s` is not a m2m field", fi.Name, mi.FullName))
	}

	return newQueryM2M(md, o, mi, fi, ind)
}

// NOTE: this method is deprecated, context parameter will not take effect.
func (o *ormBase) QueryM2MWithCtx(_ context.Context, md interface{}, name string) QueryM2Mer {
	logs.Warn("QueryM2MWithCtx is DEPRECATED. Use methods with `WithCtx` suffix on QueryM2M as replacement please.")
	return o.QueryM2M(md, name)
}

// load related models to md model.
// args are limit, offset int and order string.
//
// example:
//
//	orm.LoadRelated(post,"Tags")
//	for _,tag := range post.Tags{...}
//
// make sure the relation is defined in model struct tags.
func (o *ormBase) LoadRelated(md interface{}, name string, args ...utils.KV) (int64, error) {
	return o.LoadRelatedWithCtx(context.Background(), md, name, args...)
}

func (o *ormBase) LoadRelatedWithCtx(_ context.Context, md interface{}, name string, args ...utils.KV) (int64, error) {
	_, fi, ind, qs := o.queryRelated(md, name)

	var relDepth int
	var limit, offset int64
	var order string

	kvs := utils.NewKVs(args...)
	kvs.IfContains(hints.KeyRelDepth, func(value interface{}) {
		if v, ok := value.(bool); ok {
			if v {
				relDepth = DefaultRelsDepth
			}
		} else if v, ok := value.(int); ok {
			relDepth = v
		}
	}).IfContains(hints.KeyLimit, func(value interface{}) {
		if v, ok := value.(int64); ok {
			limit = v
		}
	}).IfContains(hints.KeyOffset, func(value interface{}) {
		if v, ok := value.(int64); ok {
			offset = v
		}
	}).IfContains(hints.KeyOrderBy, func(value interface{}) {
		if v, ok := value.(string); ok {
			order = v
		}
	})

	switch fi.FieldType {
	case RelOneToOne, RelForeignKey, RelReverseOne:
		limit = 1
		offset = 0
	}

	qs.limit = limit
	qs.offset = offset
	qs.relDepth = relDepth

	if len(order) > 0 {
		qs.orders = order_clause.ParseOrder(order)
	}

	find := ind.FieldByIndex(fi.FieldIndex)

	var nums int64
	var err error
	switch fi.FieldType {
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

// Get QuerySeter for related models to md model
func (o *ormBase) queryRelated(md interface{}, name string) (*models.ModelInfo, *models.FieldInfo, reflect.Value, *querySet) {
	mi, ind := o.getPtrMiInd(md)
	fi := o.getFieldInfo(mi, name)

	_, _, exist := getExistPk(mi, ind)
	if !exist {
		panic(ErrMissPK)
	}

	var qs *querySet

	switch fi.FieldType {
	case RelOneToOne, RelForeignKey, RelManyToMany:
		if !fi.InModel {
			break
		}
		qs = o.getRelQs(md, mi, fi)
	case RelReverseOne, RelReverseMany:
		if !fi.InModel {
			break
		}
		qs = o.getReverseQs(md, mi, fi)
	}

	if qs == nil {
		panic(fmt.Errorf("<Ormer> name `%s` for model `%s` is not an available rel/reverse field", md, name))
	}

	return mi, fi, ind, qs
}

// Get reverse relation QuerySeter
func (o *ormBase) getReverseQs(md interface{}, mi *models.ModelInfo, fi *models.FieldInfo) *querySet {
	switch fi.FieldType {
	case RelReverseOne, RelReverseMany:
	default:
		panic(fmt.Errorf("<Ormer> name `%s` for model `%s` is not an available reverse field", fi.Name, mi.FullName))
	}

	var q *querySet

	if fi.FieldType == RelReverseMany && fi.ReverseFieldInfo.Mi.IsThrough {
		q = newQuerySet(o, fi.RelModelInfo).(*querySet)
		q.cond = NewCondition().And(fi.ReverseFieldInfoM2M.Column+ExprSep+fi.ReverseFieldInfo.Column, md)
	} else {
		q = newQuerySet(o, fi.ReverseFieldInfo.Mi).(*querySet)
		q.cond = NewCondition().And(fi.ReverseFieldInfo.Column, md)
	}

	return q
}

// Get relation QuerySeter
func (o *ormBase) getRelQs(md interface{}, mi *models.ModelInfo, fi *models.FieldInfo) *querySet {
	switch fi.FieldType {
	case RelOneToOne, RelForeignKey, RelManyToMany:
	default:
		panic(fmt.Errorf("<Ormer> name `%s` for model `%s` is not an available rel field", fi.Name, mi.FullName))
	}

	q := newQuerySet(o, fi.RelModelInfo).(*querySet)
	q.cond = NewCondition()

	if fi.FieldType == RelManyToMany {
		q.cond = q.cond.And(fi.ReverseFieldInfoM2M.Column+ExprSep+fi.ReverseFieldInfo.Column, md)
	} else {
		q.cond = q.cond.And(fi.ReverseFieldInfo.Column, md)
	}

	return q
}

// return a QuerySeter for table operations.
// table name can be string or struct.
// e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil)),
func (o *ormBase) QueryTable(ptrStructOrTableName interface{}) (qs QuerySeter) {
	var name string
	if table, ok := ptrStructOrTableName.(string); ok {
		name = models.NameStrategyMap[models.DefaultNameStrategy](table)
		if mi, ok := defaultModelCache.Get(name); ok {
			qs = newQuerySet(o, mi)
		}
	} else {
		name = models.GetFullName(iutils.IndirectType(reflect.TypeOf(ptrStructOrTableName)))
		if mi, ok := defaultModelCache.GetByFullName(name); ok {
			qs = newQuerySet(o, mi)
		}
	}
	if qs == nil {
		panic(fmt.Errorf("<Ormer.QueryTable> table name: `%s` not exists", name))
	}
	return qs
}

// Deprecated: QueryTableWithCtx is deprecated, context parameter will not take effect.
func (o *ormBase) QueryTableWithCtx(_ context.Context, ptrStructOrTableName interface{}) (qs QuerySeter) {
	logs.Warn("QueryTableWithCtx is DEPRECATED. Use methods with `WithCtx` suffix on QuerySeter as replacement please.")
	return o.QueryTable(ptrStructOrTableName)
}

// Raw return a raw query seter for raw sql string.
func (o *ormBase) Raw(query string, args ...interface{}) RawSeter {
	return o.RawWithCtx(context.Background(), query, args...)
}

func (o *ormBase) RawWithCtx(_ context.Context, query string, args ...interface{}) RawSeter {
	return newRawSet(o, query, args)
}

// Driver return current using database Driver
func (o *ormBase) Driver() Driver {
	return driver(o.alias.Name)
}

// DBStats return sql.DBStats for current database
func (o *ormBase) DBStats() *sql.DBStats {
	if o.alias != nil && o.alias.DB != nil {
		stats := o.alias.DB.DB.Stats()
		return &stats
	}
	return nil
}

type orm struct {
	ormBase
}

var _ Ormer = new(orm)

func (o *orm) Begin() (TxOrmer, error) {
	return o.BeginWithCtx(context.Background())
}

func (o *orm) BeginWithCtx(ctx context.Context) (TxOrmer, error) {
	return o.BeginWithCtxAndOpts(ctx, nil)
}

func (o *orm) BeginWithOpts(opts *sql.TxOptions) (TxOrmer, error) {
	return o.BeginWithCtxAndOpts(context.Background(), opts)
}

func (o *orm) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error) {
	tx, err := o.db.(txer).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	_txOrm := &txOrm{
		ormBase: ormBase{
			alias: o.alias,
			db:    &TxDB{tx: tx},
		},
	}

	if Debug {
		_txOrm.db = newDbQueryLog(o.alias, _txOrm.db)
	}

	var taskTxOrm TxOrmer = _txOrm
	return taskTxOrm, nil
}

func (o *orm) DoTx(task func(ctx context.Context, txOrm TxOrmer) error) error {
	return o.DoTxWithCtx(context.Background(), task)
}

func (o *orm) DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return o.DoTxWithCtxAndOpts(ctx, nil, task)
}

func (o *orm) DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return o.DoTxWithCtxAndOpts(context.Background(), opts, task)
}

func (o *orm) DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error {
	return doTxTemplate(ctx, o, opts, task)
}

func doTxTemplate(ctx context.Context, o TxBeginner, opts *sql.TxOptions,
	task func(ctx context.Context, txOrm TxOrmer) error) error {
	_txOrm, err := o.BeginWithCtxAndOpts(ctx, opts)
	if err != nil {
		return err
	}
	panicked := true
	defer func() {
		if panicked || err != nil {
			e := _txOrm.Rollback()
			if e != nil {
				logs.Error("rollback transaction failed: %v,%v", e, panicked)
			}
		} else {
			e := _txOrm.Commit()
			if e != nil {
				logs.Error("commit transaction failed: %v,%v", e, panicked)
			}
		}
	}()
	taskTxOrm := _txOrm
	err = task(ctx, taskTxOrm)
	panicked = false
	return err
}

type txOrm struct {
	ormBase
}

var _ TxOrmer = new(txOrm)

func (t *txOrm) Commit() error {
	return t.db.(txEnder).Commit()
}

func (t *txOrm) Rollback() error {
	return t.db.(txEnder).Rollback()
}

func (t *txOrm) RollbackUnlessCommit() error {
	return t.db.(txEnder).RollbackUnlessCommit()
}

// NewOrm create new orm
func NewOrm() Ormer {
	BootStrap() // execute only once
	return NewOrmUsingDB(`default`)
}

// NewOrmUsingDB create new orm with the name
func NewOrmUsingDB(aliasName string) Ormer {
	if al, ok := dataBaseCache.get(aliasName); ok {
		return newDBWithAlias(al)
	}
	panic(fmt.Errorf("<Ormer.Using> unknown db alias name `%s`", aliasName))
}

// NewOrmWithDB create a new ormer object with specify *sql.DB for query
func NewOrmWithDB(driverName, aliasName string, db *sql.DB, params ...DBOption) (Ormer, error) {
	al, err := newAliasWithDb(aliasName, driverName, db, params...)
	if err != nil {
		return nil, err
	}

	return newDBWithAlias(al), nil
}

func newDBWithAlias(al *alias) Ormer {
	o := new(orm)
	o.alias = al

	if Debug {
		o.db = newDbQueryLog(al, al.DB)
	} else {
		o.db = al.DB
	}

	if len(globalFilterChains) > 0 {
		return NewFilterOrmDecorator(o, globalFilterChains...)
	}
	return o
}
