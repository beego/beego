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
	"reflect"
	"time"

	"github.com/beego/beego/v2/client/orm/internal/models"

	"github.com/beego/beego/v2/client/orm/clauses/order_clause"
	"github.com/beego/beego/v2/core/utils"
)

// TableNameI is usually used by model
// when you custom your table name, please implement this interfaces
// for example:
//
//	type User struct {
//	  ...
//	}
//
//	func (u *User) TableName() string {
//	   return "USER_TABLE"
//	}
type TableNameI interface {
	TableName() string
}

// TableEngineI is usually used by model
// when you want to use specific engine, like myisam, you can implement this interface
// for example:
//
//	type User struct {
//	  ...
//	}
//
//	func (u *User) TableEngine() string {
//	   return "myisam"
//	}
type TableEngineI interface {
	TableEngine() string
}

// TableIndexI is usually used by model
// when you want to create indexes, you can implement this interface
// for example:
//
//	type User struct {
//	  ...
//	}
//
//	func (u *User) TableIndex() [][]string {
//	   return [][]string{{"Name"}}
//	}
type TableIndexI interface {
	TableIndex() [][]string
}

// TableUniqueI is usually used by model
// when you want to create unique indexes, you can implement this interface
// for example:
//
//	type User struct {
//	  ...
//	}
//
//	func (u *User) TableUnique() [][]string {
//	   return [][]string{{"Email"}}
//	}
type TableUniqueI interface {
	TableUnique() [][]string
}

// IsApplicableTableForDB if return false, we won't create table to this db
type IsApplicableTableForDB interface {
	IsApplicableTableForDB(db string) bool
}

// Driver define database driver
type Driver interface {
	Name() string
	Type() DriverType
}

type Fielder = models.Fielder

type TxBeginner interface {
	// Begin self control transaction
	Begin() (TxOrmer, error)
	BeginWithCtx(ctx context.Context) (TxOrmer, error)
	BeginWithOpts(opts *sql.TxOptions) (TxOrmer, error)
	BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error)

	// DoTx closure control transaction
	DoTx(task func(ctx context.Context, txOrm TxOrmer) error) error
	DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm TxOrmer) error) error
	DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error
	DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error
}

type TxCommitter interface {
	txEnder
}

// transaction beginner
type txer interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// transaction ending
type txEnder interface {
	Commit() error
	Rollback() error

	// RollbackUnlessCommit if the transaction has been committed, do nothing, or transaction will be rollback
	// For example:
	// ```go
	//    txOrm := orm.Begin()
	//    defer txOrm.RollbackUnlessCommit()
	//    err := txOrm.Insert() // do something
	//    if err != nil {
	//       return err
	//    }
	//    txOrm.Commit()
	// ```
	RollbackUnlessCommit() error
}

// DML Data Manipulation Language
type DML interface {
	// Insert insert model data to database
	// for example:
	//  user := new(User)
	//  id, err = Ormer.Insert(user)
	//  user must be a pointer and Insert will Set user's pk field
	Insert(md interface{}) (int64, error)
	InsertWithCtx(ctx context.Context, md interface{}) (int64, error)
	// InsertOrUpdate mysql:InsertOrUpdate(model) or InsertOrUpdate(model,"colu=colu+value")
	// if colu type is integer : can use(+-*/), string : convert(colu,"value")
	// postgres: InsertOrUpdate(model,"conflictColumnName") or InsertOrUpdate(model,"conflictColumnName","colu=colu+value")
	// if colu type is integer : can use(+-*/), string : colu || "value"
	InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error)
	InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error)
	// InsertMulti inserts some models to database
	InsertMulti(bulk int, mds interface{}) (int64, error)
	InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error)
	// Update updates model to database.
	// cols Set the Columns those want to update.
	// find model by Id(pk) field and update Columns specified by Fields, if cols is null then update All Columns
	// for example:
	// user := User{Id: 2}
	//	user.Langs = append(user.Langs, "zh-CN", "en-US")
	//	user.Extra.Name = "beego"
	//	user.Extra.Data = "orm"
	//	num, err = Ormer.Update(&user, "Langs", "Extra")
	Update(md interface{}, cols ...string) (int64, error)
	UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error)
	// Delete deletes model in database
	Delete(md interface{}, cols ...string) (int64, error)
	DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error)

	// Raw return a raw query seter for raw sql string.
	// for example:
	//	 ormer.Raw("UPDATE `user` SET `user_name` = ? WHERE `user_name` = ?", "slene", "testing").Exec()
	//	// update user testing's name to slene
	Raw(query string, args ...interface{}) RawSeter
	RawWithCtx(ctx context.Context, query string, args ...interface{}) RawSeter
}

// DQL Data Query Language
type DQL interface {
	// Read reads data to model
	// for example:
	//	this will find User by Id field
	// 	u = &User{Id: user.Id}
	// 	err = Ormer.Read(u)
	//	this will find User by UserName field
	// 	u = &User{UserName: "astaxie", Password: "pass"}
	//	err = Ormer.Read(u, "UserName")
	Read(md interface{}, cols ...string) error
	ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error

	// ReadForUpdate Like Read(), but with "FOR UPDATE" clause, useful in transaction.
	// Some databases are not support this feature.
	ReadForUpdate(md interface{}, cols ...string) error
	ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error

	// ReadOrCreate Try to read a row from the database, or insert one if it doesn't exist
	ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error)
	ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error)

	// LoadRelated load related models to md model.
	// args are limit, offset int and order string.
	//
	// example:
	// 	Ormer.LoadRelated(post,"Tags")
	// 	for _,tag := range post.Tags{...}
	// hints.DefaultRelDepth useDefaultRelsDepth ; or depth 0
	// hints.RelDepth loadRelationDepth
	// hints.Limit limit default limit 1000
	// hints.Offset int offset default offset 0
	// hints.OrderBy string order  for example : "-Id"
	// make sure the relation is defined in model struct tags.
	LoadRelated(md interface{}, name string, args ...utils.KV) (int64, error)
	LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...utils.KV) (int64, error)

	// QueryM2M create a models to models queryer
	// for example:
	// 	post := Post{Id: 4}
	// 	m2m := Ormer.QueryM2M(&post, "Tags")
	QueryM2M(md interface{}, name string) QueryM2Mer
	// QueryM2MWithCtx NOTE: this method is deprecated, context parameter will not take effect.
	// Use context.Context directly on methods with `WithCtx` suffix such as InsertWithCtx/UpdateWithCtx
	QueryM2MWithCtx(ctx context.Context, md interface{}, name string) QueryM2Mer

	// QueryTable return a QuerySeter for table operations.
	// table name can be string or struct.
	// e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil)),
	QueryTable(ptrStructOrTableName interface{}) QuerySeter
	// QueryTableWithCtx NOTE: this method is deprecated, context parameter will not take effect.
	// Use context.Context directly on methods with `WithCtx` suffix such as InsertWithCtx/UpdateWithCtx
	QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) QuerySeter

	DBStats() *sql.DBStats
}

type DriverGetter interface {
	Driver() Driver
}

type ormer interface {
	DQL
	DML
	DriverGetter
}

// QueryExecutor wrapping for ormer
type QueryExecutor interface {
	ormer
}

type Ormer interface {
	QueryExecutor
	TxBeginner
}

type TxOrmer interface {
	QueryExecutor
	TxCommitter
}

// Inserter insert prepared statement
type Inserter interface {
	Insert(interface{}) (int64, error)
	InsertWithCtx(context.Context, interface{}) (int64, error)
	Close() error
}

// QuerySeter query seter
type QuerySeter interface {
	// Filter add condition expression to QuerySeter.
	// for example:
	//	filter by UserName == 'slene'
	//	qs.Filter("UserName", "slene")
	//	sql : left outer join profile on t0.id1==t1.id2 where t1.age == 28
	//	Filter("profile__Age", 28)
	// 	 // time compare
	//	qs.Filter("created", time.Now())
	Filter(string, ...interface{}) QuerySeter
	// FilterRaw add raw sql to querySeter.
	// for example:
	// qs.FilterRaw("user_id IN (SELECT id FROM profile WHERE age>=18)")
	// //sql-> WHERE user_id IN (SELECT id FROM profile WHERE age>=18)
	FilterRaw(string, string) QuerySeter
	// Exclude add NOT condition to querySeter.
	// have the same usage as Filter
	Exclude(string, ...interface{}) QuerySeter
	// SetCond Set condition to QuerySeter.
	// sql's where condition
	//	cond := orm.NewCondition()
	//	cond1 := cond.And("profile__isnull", false).AndNot("status__in", 1).Or("profile__age__gt", 2000)
	//	//sql-> WHERE T0.`profile_id` IS NOT NULL AND NOT T0.`Status` IN (?) OR T1.`age` >  2000
	//	num, err := qs.SetCond(cond1).Count()
	SetCond(*Condition) QuerySeter
	// GetCond Get condition from QuerySeter.
	// sql's where condition
	//  cond := orm.NewCondition()
	//  cond = cond.And("profile__isnull", false).AndNot("status__in", 1)
	//  qs = qs.SetCond(cond)
	//  cond = qs.GetCond()
	//  cond := cond.Or("profile__age__gt", 2000)
	//  //sql-> WHERE T0.`profile_id` IS NOT NULL AND NOT T0.`Status` IN (?) OR T1.`age` >  2000
	//  num, err := qs.SetCond(cond).Count()
	GetCond() *Condition
	// Limit add LIMIT value.
	// args[0] means offset, e.g. LIMIT num,offset.
	// if Limit <= 0 then Limit will be Set to default limit ,eg 1000
	// if QuerySeter doesn't call Limit, the sql's Limit will be Set to default limit, eg 1000
	//  for example:
	//	qs.Limit(10, 2)
	//	// sql-> limit 10 offset 2
	Limit(limit interface{}, args ...interface{}) QuerySeter
	// Offset add OFFSET value
	// same as Limit function's args[0]
	Offset(offset interface{}) QuerySeter
	// GroupBy add GROUP BY expression
	// for example:
	//	qs.GroupBy("id")
	GroupBy(exprs ...string) QuerySeter
	// OrderBy add ORDER expression.
	// "column" means ASC, "-column" means DESC.
	// for example:
	//	qs.OrderBy("-status")
	OrderBy(exprs ...string) QuerySeter
	// OrderClauses add ORDER expression by order clauses
	// for example:
	//	OrderClauses(
	//		order_clause.Clause(
	//			order.Column("Id"),
	//			order.SortAscending(),
	//		),
	//		order_clause.Clause(
	//			order.Column("status"),
	//			order.SortDescending(),
	//		),
	//	)
	//	OrderClauses(order_clause.Clause(
	//		order_clause.Column(`user__status`),
	//		order_clause.SortDescending(),//default None
	//	))
	//	OrderClauses(order_clause.Clause(
	//		order_clause.Column(`random()`),
	//		order_clause.SortNone(),//default None
	//		order_clause.Raw(),//default false.if true, do not check field is valid or not
	//	))
	OrderClauses(orders ...*order_clause.Order) QuerySeter
	// ForceIndex add FORCE INDEX expression.
	// for example:
	//	qs.ForceIndex(`idx_name1`,`idx_name2`)
	// ForceIndex, UseIndex , IgnoreIndex are mutually exclusive
	ForceIndex(indexes ...string) QuerySeter
	// UseIndex add USE INDEX expression.
	// for example:
	//	qs.UseIndex(`idx_name1`,`idx_name2`)
	// ForceIndex, UseIndex , IgnoreIndex are mutually exclusive
	UseIndex(indexes ...string) QuerySeter
	// IgnoreIndex add IGNORE INDEX expression.
	// for example:
	//	qs.IgnoreIndex(`idx_name1`,`idx_name2`)
	// ForceIndex, UseIndex , IgnoreIndex are mutually exclusive
	IgnoreIndex(indexes ...string) QuerySeter
	// RelatedSel Set relation model to query together.
	// it will query relation models and assign to parent model.
	// for example:
	//	// will load All related Fields use left join .
	// 	qs.RelatedSel().One(&user)
	//	// will  load related field only profile
	//	qs.RelatedSel("profile").One(&user)
	//	user.Profile.Age = 32
	RelatedSel(params ...interface{}) QuerySeter
	// Distinct Set Distinct
	// for example:
	//  o.QueryTable("policy").Filter("Groups__Group__Users__User", user).
	//    Distinct().
	//    All(&permissions)
	Distinct() QuerySeter
	// ForUpdate Set FOR UPDATE to query.
	// for example:
	//  o.QueryTable("user").Filter("uid", uid).ForUpdate().All(&users)
	ForUpdate() QuerySeter
	// Count returns QuerySeter execution result number
	// for example:
	//	num, err = qs.Filter("profile__age__gt", 28).Count()
	Count() (int64, error)
	CountWithCtx(context.Context) (int64, error)
	// Exist check result empty or not after QuerySeter executed
	// the same as QuerySeter.Count > 0
	Exist() bool
	ExistWithCtx(context.Context) bool
	// Update execute update with parameters
	// for example:
	//	num, err = qs.Filter("user_name", "slene").Update(Params{
	//		"Nums": ColValue(Col_Minus, 50),
	//	}) // user slene's Nums will minus 50
	//	num, err = qs.Filter("UserName", "slene").Update(Params{
	//		"user_name": "slene2"
	//	}) // user slene's  name will change to slene2
	Update(values Params) (int64, error)
	UpdateWithCtx(ctx context.Context, values Params) (int64, error)
	// Delete delete from table
	// for example:
	//	num ,err = qs.Filter("user_name__in", "testing1", "testing2").Delete()
	// 	//delete two user  who's name is testing1 or testing2
	Delete() (int64, error)
	DeleteWithCtx(context.Context) (int64, error)
	// PrepareInsert return an insert queryer.
	// it can be used in times.
	// example:
	// 	i,err := sq.PrepareInsert()
	// 	num, err = i.Insert(&user1) // user table will add one record user1 at once
	//	num, err = i.Insert(&user2) // user table will add one record user2 at once
	//	err = i.Close() //don't forget call Close
	PrepareInsert() (Inserter, error)
	PrepareInsertWithCtx(context.Context) (Inserter, error)
	// All query All data and map to containers.
	// cols means the Columns when querying.
	// for example:
	//	var users []*User
	//	qs.All(&users) // users[0],users[1],users[2] ...
	All(container interface{}, cols ...string) (int64, error)
	AllWithCtx(ctx context.Context, container interface{}, cols ...string) (int64, error)
	// One query one row data and map to containers.
	// cols means the Columns when querying.
	// for example:
	//	var user User
	//	qs.One(&user) //user.UserName == "slene"
	One(container interface{}, cols ...string) error
	OneWithCtx(ctx context.Context, container interface{}, cols ...string) error
	// Values query All data and map to []map[string]interface.
	// expres means condition expression.
	// it converts data to []map[column]value.
	// for example:
	//	var maps []Params
	//	qs.Values(&maps) //maps[0]["UserName"]=="slene"
	Values(results *[]Params, exprs ...string) (int64, error)
	ValuesWithCtx(ctx context.Context, results *[]Params, exprs ...string) (int64, error)
	// ValuesList query All data and map to [][]interface
	// it converts data to [][column_index]value
	// for example:
	//	var list []ParamsList
	//	qs.ValuesList(&list) // list[0][1] == "slene"
	ValuesList(results *[]ParamsList, exprs ...string) (int64, error)
	ValuesListWithCtx(ctx context.Context, results *[]ParamsList, exprs ...string) (int64, error)
	// ValuesFlat query All data and map to []interface.
	// it's designed for one column record Set, auto change to []value, not [][column]value.
	// for example:
	//	var list ParamsList
	//	qs.ValuesFlat(&list, "UserName") // list[0] == "slene"
	ValuesFlat(result *ParamsList, expr string) (int64, error)
	ValuesFlatWithCtx(ctx context.Context, result *ParamsList, expr string) (int64, error)
	// RowsToMap query All rows into map[string]interface with specify key and value column name.
	// keyCol = "name", valueCol = "value"
	// table data
	// name  | value
	// total | 100
	// found | 200
	// to map[string]interface{}{
	// 	"total": 100,
	// 	"found": 200,
	// }
	RowsToMap(result *Params, keyCol, valueCol string) (int64, error)
	// RowsToStruct query All rows into struct with specify key and value column name.
	// keyCol = "name", valueCol = "value"
	// table data
	// name  | value
	// total | 100
	// found | 200
	// to struct {
	// 	Total int
	// 	Found int
	// }
	RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error)
	// Aggregate aggregate func.
	// for example:
	// type result struct {
	//  DeptName string
	//	Total    int
	// }
	// var res []result
	//  o.QueryTable("dept_info").Aggregate("dept_name,sum(salary) as total").GroupBy("dept_name").All(&res)
	Aggregate(s string) QuerySeter
}

// QueryM2Mer model to model query struct
// All operations are on the m2m table only, will not affect the origin model table
type QueryM2Mer interface {
	// Add adds models to origin models when creating queryM2M.
	// example:
	// 	m2m := orm.QueryM2M(post,"Tag")
	// 	m2m.Add(&Tag1{},&Tag2{})
	//  	for _,tag := range post.Tags{}{ ... }
	// param could also be any of the follow
	// 	[]*Tag{{Id:3,Name: "TestTag1"}, {Id:4,Name: "TestTag2"}}
	//	&Tag{Id:5,Name: "TestTag3"}
	//	[]interface{}{&Tag{Id:6,Name: "TestTag4"}}
	// insert one or more rows to m2m table
	// make sure the relation is defined in post model struct tag.
	Add(...interface{}) (int64, error)
	AddWithCtx(context.Context, ...interface{}) (int64, error)
	// Remove removes models following the origin model relationship
	// only delete rows from m2m table
	// for example:
	// tag3 := &Tag{Id:5,Name: "TestTag3"}
	// num, err = m2m.Remove(tag3)
	Remove(...interface{}) (int64, error)
	RemoveWithCtx(context.Context, ...interface{}) (int64, error)
	// Exist checks model is existed in relationship of origin model
	Exist(interface{}) bool
	ExistWithCtx(context.Context, interface{}) bool
	// Clear cleans All models in related of origin model
	Clear() (int64, error)
	ClearWithCtx(context.Context) (int64, error)
	// Count counts All related models of origin model
	Count() (int64, error)
	CountWithCtx(context.Context) (int64, error)
}

// RawPreparer raw query statement
type RawPreparer interface {
	Exec(...interface{}) (sql.Result, error)
	Close() error
}

// RawSeter raw query seter
// create From Ormer.Raw
// for example:
//
//	sql := fmt.Sprintf("SELECT %sid%s,%sname%s FROM %suser%s WHERE id = ?",Q,Q,Q,Q,Q,Q)
//	rs := Ormer.Raw(sql, 1)
type RawSeter interface {
	// Exec execute sql and Get result
	Exec() (sql.Result, error)
	// QueryRow query data and map to container
	// for example:
	//	var name string
	//	var id int
	//	rs.QueryRow(&id,&name) // id==2 name=="slene"
	QueryRow(containers ...interface{}) error

	// QueryRows query data rows and map to container
	//	var ids []int
	//	var names []int
	//	query = fmt.Sprintf("SELECT 'id','name' FROM %suser%s", Q, Q)
	//	num, err = dORM.Raw(query).QueryRows(&ids,&names) // ids=>{1,2},names=>{"nobody","slene"}
	QueryRows(containers ...interface{}) (int64, error)
	SetArgs(...interface{}) RawSeter
	// Values query data to []map[string]interface
	// see QuerySeter's Values
	Values(container *[]Params, cols ...string) (int64, error)
	// ValuesList query data to [][]interface
	// see QuerySeter's ValuesList
	ValuesList(container *[]ParamsList, cols ...string) (int64, error)
	// ValuesFlat query data to []interface
	// see QuerySeter's ValuesFlat
	ValuesFlat(container *ParamsList, cols ...string) (int64, error)
	// RowsToMap query All rows into map[string]interface with specify key and value column name.
	// keyCol = "name", valueCol = "value"
	// table data
	// name  | value
	// total | 100
	// found | 200
	// to map[string]interface{}{
	// 	"total": 100,
	// 	"found": 200,
	// }
	RowsToMap(result *Params, keyCol, valueCol string) (int64, error)
	// RowsToStruct query All rows into struct with specify key and value column name.
	// keyCol = "name", valueCol = "value"
	// table data
	// name  | value
	// total | 100
	// found | 200
	// to struct {
	// 	Total int
	// 	Found int
	// }
	RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error)

	// Prepare return prepared raw statement for used in times.
	// for example:
	// 	pre, err := dORM.Raw("INSERT INTO tag (name) VALUES (?)").Prepare()
	// 	r, err := pre.Exec("name1") // INSERT INTO tag (name) VALUES (`name1`)
	Prepare() (RawPreparer, error)
}

// stmtQuerier statement querier
type stmtQuerier interface {
	Close() error
	Exec(args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
	Query(args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
	QueryRow(args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
}

// db querier
type dbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// type DB interface {
// 	Begin() (*sql.Tx, error)
// 	Prepare(query string) (stmtQuerier, error)
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	Query(query string, args ...interface{}) (*sql.Rows, error)
// 	QueryRow(query string, args ...interface{}) *sql.Row
// }

// base database struct
type dbBaser interface {
	Read(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location, []string, bool) error
	ReadBatch(context.Context, dbQuerier, querySet, *models.ModelInfo, *Condition, interface{}, *time.Location, []string) (int64, error)
	Count(context.Context, dbQuerier, querySet, *models.ModelInfo, *Condition, *time.Location) (int64, error)
	ReadValues(context.Context, dbQuerier, querySet, *models.ModelInfo, *Condition, []string, interface{}, *time.Location) (int64, error)

	Insert(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location) (int64, error)
	InsertOrUpdate(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *alias, ...string) (int64, error)
	InsertMulti(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, int, *time.Location) (int64, error)
	InsertValue(context.Context, dbQuerier, *models.ModelInfo, bool, []string, []interface{}) (int64, error)
	InsertStmt(context.Context, stmtQuerier, *models.ModelInfo, reflect.Value, *time.Location) (int64, error)

	Update(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location, []string) (int64, error)
	UpdateBatch(context.Context, dbQuerier, *querySet, *models.ModelInfo, *Condition, Params, *time.Location) (int64, error)

	Delete(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location, []string) (int64, error)
	DeleteBatch(context.Context, dbQuerier, *querySet, *models.ModelInfo, *Condition, *time.Location) (int64, error)

	SupportUpdateJoin() bool
	OperatorSQL(string) string
	GenerateOperatorSQL(*models.ModelInfo, *models.FieldInfo, string, []interface{}, *time.Location) (string, []interface{})
	GenerateOperatorLeftCol(*models.FieldInfo, string, *string)
	PrepareInsert(context.Context, dbQuerier, *models.ModelInfo) (stmtQuerier, string, error)
	MaxLimit() uint64
	TableQuote() string
	ReplaceMarks(*string)
	HasReturningID(*models.ModelInfo, *string) bool
	TimeFromDB(*time.Time, *time.Location)
	TimeToDB(*time.Time, *time.Location)
	DbTypes() map[string]string
	GetTables(dbQuerier) (map[string]bool, error)
	GetColumns(context.Context, dbQuerier, string) (map[string][3]string, error)
	ShowTablesQuery() string
	ShowColumnsQuery(string) string
	IndexExists(context.Context, dbQuerier, string, string) bool
	collectFieldValue(*models.ModelInfo, *models.FieldInfo, reflect.Value, bool, *time.Location) (interface{}, error)
	setval(context.Context, dbQuerier, *models.ModelInfo, []string) error

	GenerateSpecifyIndex(tableName string, useIndex int, indexes []string) string
}
