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

//import (
//	"context"
//	"database/sql"
//	"github.com/beego/beego/v2/client/orm/internal/condition"
//	"github.com/beego/beego/v2/client/orm/internal/queryset"
//	"github.com/beego/beego/v2/client/orm/internal/session"
//	"reflect"
//	"time"
//
//	"github.com/beego/beego/v2/client/orm/internal/models"
//)
//
//// TableNameI is usually used by model
//// when you custom your table name, please implement this interfaces
//// for example:
////
////	type User struct {
////	  ...
////	}
////
////	func (u *User) TableName() string {
////	   return "USER_TABLE"
////	}
//type TableNameI interface {
//	TableName() string
//}
//
//// TableEngineI is usually used by model
//// when you want to use specific engine, like myisam, you can implement this interface
//// for example:
////
////	type User struct {
////	  ...
////	}
////
////	func (u *User) TableEngine() string {
////	   return "myisam"
////	}
//type TableEngineI interface {
//	TableEngine() string
//}
//
//// TableIndexI is usually used by model
//// when you want to create indexes, you can implement this interface
//// for example:
////
////	type User struct {
////	  ...
////	}
////
////	func (u *User) TableIndex() [][]string {
////	   return [][]string{{"Name"}}
////	}
//type TableIndexI interface {
//	TableIndex() [][]string
//}
//
//// TableUniqueI is usually used by model
//// when you want to create unique indexes, you can implement this interface
//// for example:
////
////	type User struct {
////	  ...
////	}
////
////	func (u *User) TableUnique() [][]string {
////	   return [][]string{{"Email"}}
////	}
//type TableUniqueI interface {
//	TableUnique() [][]string
//}
//
//// IsApplicableTableForDB if return false, we won't create table to this db
//type IsApplicableTableForDB interface {
//	IsApplicableTableForDB(db string) bool
//}
//
//// Driver define database driver
//type Driver interface {
//	Name() string
//	Type() session.DriverType
//}
//
//type Fielder = models.Fielder
//
//type TxBeginner interface {
//	// Begin self control transaction
//	Begin() (TxOrmer, error)
//	BeginWithCtx(ctx context.Context) (TxOrmer, error)
//	BeginWithOpts(opts *sql.TxOptions) (TxOrmer, error)
//	BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (TxOrmer, error)
//
//	// DoTx closure control transaction
//	DoTx(task func(ctx context.Context, txOrm TxOrmer) error) error
//	DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm TxOrmer) error) error
//	DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error
//	DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm TxOrmer) error) error
//}
//
//type TxCommitter interface {
//	txEnder
//}
//
//// transaction beginner
//type txer interface {
//	Begin() (*sql.Tx, error)
//	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
//}
//
//// transaction ending
//type txEnder interface {
//	Commit() error
//	Rollback() error
//
//	// RollbackUnlessCommit if the transaction has been committed, do nothing, or transaction will be rollback
//	// For example:
//	// ```go
//	//    txOrm := orm.Begin()
//	//    defer txOrm.RollbackUnlessCommit()
//	//    err := txOrm.Insert() // do something
//	//    if err != nil {
//	//       return err
//	//    }
//	//    txOrm.Commit()
//	// ```
//	RollbackUnlessCommit() error
//}
//
//// DML Data Manipulation Language
//type DML interface {
//	// Insert insert model data to database
//	// for example:
//	//  user := new(User)
//	//  id, err = Ormer.Insert(user)
//	//  user must be a pointer and Insert will Set user's pk field
//	Insert(md interface{}) (int64, error)
//	InsertWithCtx(ctx context.Context, md interface{}) (int64, error)
//	// InsertOrUpdate mysql:InsertOrUpdate(model) or InsertOrUpdate(model,"colu=colu+value")
//	// if colu type is integer : can use(+-*/), string : convert(colu,"value")
//	// postgres: InsertOrUpdate(model,"conflictColumnName") or InsertOrUpdate(model,"conflictColumnName","colu=colu+value")
//	// if colu type is integer : can use(+-*/), string : colu || "value"
//	InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error)
//	InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error)
//	// InsertMulti inserts some models to database
//	InsertMulti(bulk int, mds interface{}) (int64, error)
//	InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error)
//	// Update updates model to database.
//	// cols Set the Columns those want to update.
//	// find model by Id(pk) field and update Columns specified by Fields, if cols is null then update All Columns
//	// for example:
//	// user := User{Id: 2}
//	//	user.Langs = append(user.Langs, "zh-CN", "en-US")
//	//	user.Extra.Name = "beego"
//	//	user.Extra.Data = "orm"
//	//	num, err = Ormer.Update(&user, "Langs", "Extra")
//	Update(md interface{}, cols ...string) (int64, error)
//	UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error)
//	// Delete deletes model in database
//	Delete(md interface{}, cols ...string) (int64, error)
//	DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error)
//
//	// Raw return a raw query seter for raw sql string.
//	// for example:
//	//	 ormer.Raw("UPDATE `user` SET `user_name` = ? WHERE `user_name` = ?", "slene", "testing").Exec()
//	//	// update user testing's name to slene
//	Raw(query string, args ...interface{}) RawSeter
//	RawWithCtx(ctx context.Context, query string, args ...interface{}) RawSeter
//	ExecRaw(ctx context.Context, md interface{}, query string, args ...any) (sql.Result, error)
//}
//
//type DriverGetter interface {
//	Driver() Driver
//}
//
//type ormer interface {
//	DQL
//	DML
//	DriverGetter
//}
//
//// QueryExecutor wrapping for ormer
//type QueryExecutor interface {
//	ormer
//}
//
//type Ormer interface {
//	QueryExecutor
//	TxBeginner
//}
//
//type TxOrmer interface {
//	QueryExecutor
//	TxCommitter
//}
//
//// Inserter insert prepared statement
//type Inserter interface {
//	Insert(interface{}) (int64, error)
//	InsertWithCtx(context.Context, interface{}) (int64, error)
//	Close() error
//}
//
//// QueryM2Mer model to model query struct
//// All operations are on the m2m table only, will not affect the origin model table
//type QueryM2Mer interface {
//	// Add adds models to origin models when creating queryM2M.
//	// example:
//	// 	m2m := orm.QueryM2M(post,"Tag")
//	// 	m2m.Add(&Tag1{},&Tag2{})
//	//  	for _,tag := range post.Tags{}{ ... }
//	// param could also be any of the follow
//	// 	[]*Tag{{Id:3,Name: "TestTag1"}, {Id:4,Name: "TestTag2"}}
//	//	&Tag{Id:5,Name: "TestTag3"}
//	//	[]interface{}{&Tag{Id:6,Name: "TestTag4"}}
//	// insert one or more rows to m2m table
//	// make sure the relation is defined in post model struct tag.
//	Add(...interface{}) (int64, error)
//	AddWithCtx(context.Context, ...interface{}) (int64, error)
//	// Remove removes models following the origin model relationship
//	// only delete rows from m2m table
//	// for example:
//	// tag3 := &Tag{Id:5,Name: "TestTag3"}
//	// num, err = m2m.Remove(tag3)
//	Remove(...interface{}) (int64, error)
//	RemoveWithCtx(context.Context, ...interface{}) (int64, error)
//	// Exist checks model is existed in relationship of origin model
//	Exist(interface{}) bool
//	ExistWithCtx(context.Context, interface{}) bool
//	// Clear cleans All models in related of origin model
//	Clear() (int64, error)
//	ClearWithCtx(context.Context) (int64, error)
//	// Count counts All related models of origin model
//	Count() (int64, error)
//	CountWithCtx(context.Context) (int64, error)
//}
//
//// RawPreparer raw query statement
//type RawPreparer interface {
//	Exec(...interface{}) (sql.Result, error)
//	Close() error
//}
//
//// RawSeter raw query seter
//// create From Ormer.Raw
//// for example:
////
////	sql := fmt.Sprintf("SELECT %sid%s,%sname%s FROM %suser%s WHERE id = ?",Q,Q,Q,Q,Q,Q)
////	rs := Ormer.Raw(sql, 1)
//type RawSeter interface {
//	// Exec execute sql and Get result
//	Exec() (sql.Result, error)
//	// QueryRow query data and map to container
//	// for example:
//	//	var name string
//	//	var id int
//	//	rs.QueryRow(&id,&name) // id==2 name=="slene"
//	QueryRow(containers ...interface{}) error
//
//	// QueryRows query data rows and map to container
//	//	var ids []int
//	//	var names []int
//	//	query = fmt.Sprintf("SELECT 'id','name' FROM %suser%s", Q, Q)
//	//	num, err = dORM.Raw(query).QueryRows(&ids,&names) // ids=>{1,2},names=>{"nobody","slene"}
//	QueryRows(containers ...interface{}) (int64, error)
//	SetArgs(...interface{}) RawSeter
//	// Values query data to []map[string]interface
//	// see QuerySeter's Values
//	Values(container *[]session.Params, cols ...string) (int64, error)
//	// ValuesList query data to [][]interface
//	// see QuerySeter's ValuesList
//	ValuesList(container *[]session.ParamsList, cols ...string) (int64, error)
//	// ValuesFlat query data to []interface
//	// see QuerySeter's ValuesFlat
//	ValuesFlat(container *session.ParamsList, cols ...string) (int64, error)
//	// RowsToMap query All rows into map[string]interface with specify key and value column name.
//	// keyCol = "name", valueCol = "value"
//	// table data
//	// name  | value
//	// total | 100
//	// found | 200
//	// to map[string]interface{}{
//	// 	"total": 100,
//	// 	"found": 200,
//	// }
//	RowsToMap(result *session.Params, keyCol, valueCol string) (int64, error)
//	// RowsToStruct query All rows into struct with specify key and value column name.
//	// keyCol = "name", valueCol = "value"
//	// table data
//	// name  | value
//	// total | 100
//	// found | 200
//	// to struct {
//	// 	Total int
//	// 	Found int
//	// }
//	RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error)
//
//	// Prepare return prepared raw statement for used in times.
//	// for example:
//	// 	pre, err := dORM.Raw("INSERT INTO tag (name) VALUES (?)").Prepare()
//	// 	r, err := pre.Exec("name1") // INSERT INTO tag (name) VALUES (`name1`)
//	Prepare() (RawPreparer, error)
//}
//
//// stmtQuerier statement querier
//type stmtQuerier interface {
//	Close() error
//	Exec(args ...interface{}) (sql.Result, error)
//	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
//	Query(args ...interface{}) (*sql.Rows, error)
//	QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
//	QueryRow(args ...interface{}) *sql.Row
//	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
//}
//
//// db querier
//type dbQuerier interface {
//	Prepare(query string) (*sql.Stmt, error)
//	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
//	Exec(query string, args ...interface{}) (sql.Result, error)
//	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
//	Query(query string, args ...interface{}) (*sql.Rows, error)
//	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
//	QueryRow(query string, args ...interface{}) *sql.Row
//	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
//}
//
//// type DB interface {
//// 	Begin() (*sql.Tx, error)
//// 	Prepare(query string) (stmtQuerier, error)
//// 	Exec(query string, args ...interface{}) (sql.Result, error)
//// 	Query(query string, args ...interface{}) (*sql.Rows, error)
//// 	QueryRow(query string, args ...interface{}) *sql.Row
//// }
//
//// base database struct
//type dbBaser interface {
//	Read(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location, []string, bool) error
//	ReadRaw(ctx context.Context, q dbQuerier, mi *models.ModelInfo, ind reflect.Value, tz *time.Location, query string, args ...any) error
//	ReadBatch(context.Context, dbQuerier, queryset.querySet, *models.ModelInfo, *condition.Condition, interface{}, *time.Location, []string) (int64, error)
//	Count(context.Context, dbQuerier, queryset.querySet, *models.ModelInfo, *condition.Condition, *time.Location) (int64, error)
//	ReadValues(context.Context, dbQuerier, queryset.querySet, *models.ModelInfo, *condition.Condition, []string, interface{}, *time.Location) (int64, error)
//
//	ExecRaw(ctx context.Context, q dbQuerier, query string, args ...any) (sql.Result, error)
//	Insert(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location) (int64, error)
//	InsertOrUpdate(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *session.DB, ...string) (int64, error)
//	InsertMulti(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, int, *time.Location) (int64, error)
//	InsertValue(context.Context, dbQuerier, *models.ModelInfo, bool, []string, []interface{}) (int64, error)
//	InsertStmt(context.Context, stmtQuerier, *models.ModelInfo, reflect.Value, *time.Location) (int64, error)
//
//	Update(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location, []string) (int64, error)
//	UpdateBatch(context.Context, dbQuerier, *queryset.querySet, *models.ModelInfo, *condition.Condition, session.Params, *time.Location) (int64, error)
//
//	Delete(context.Context, dbQuerier, *models.ModelInfo, reflect.Value, *time.Location, []string) (int64, error)
//	DeleteBatch(context.Context, dbQuerier, *queryset.querySet, *models.ModelInfo, *condition.Condition, *time.Location) (int64, error)
//
//	SupportUpdateJoin() bool
//	OperatorSQL(string) string
//	GenerateOperatorSQL(*models.ModelInfo, *models.FieldInfo, string, []interface{}, *time.Location) (string, []interface{})
//	GenerateOperatorLeftCol(*models.FieldInfo, string, *string)
//	PrepareInsert(context.Context, dbQuerier, *models.ModelInfo) (stmtQuerier, string, error)
//	MaxLimit() uint64
//	TableQuote() string
//	ReplaceMarks(*string)
//	HasReturningID(*models.ModelInfo, *string) bool
//	TimeFromDB(*time.Time, *time.Location)
//	TimeToDB(*time.Time, *time.Location)
//	DbTypes() map[string]string
//	GetTables(dbQuerier) (map[string]bool, error)
//	GetColumns(context.Context, dbQuerier, string) (map[string][3]string, error)
//	ShowTablesQuery() string
//	ShowColumnsQuery(string) string
//	IndexExists(context.Context, dbQuerier, string, string) bool
//	collectFieldValue(*models.ModelInfo, *models.FieldInfo, reflect.Value, bool, *time.Location) (interface{}, error)
//	setval(context.Context, dbQuerier, *models.ModelInfo, []string) error
//
//	GenerateSpecifyIndex(tableName string, useIndex int, indexes []string) string
//}
