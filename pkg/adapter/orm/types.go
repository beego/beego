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

	"github.com/astaxie/beego/pkg/client/orm"
)

// Params stores the Params
type Params orm.Params

// ParamsList stores paramslist
type ParamsList orm.ParamsList

// Driver define database driver
type Driver orm.Driver

// Fielder define field info
type Fielder orm.Fielder

// Ormer define the orm interface
type Ormer interface {
	// read data to model
	// for example:
	//	this will find User by Id field
	// 	u = &User{Id: user.Id}
	// 	err = Ormer.Read(u)
	//	this will find User by UserName field
	// 	u = &User{UserName: "astaxie", Password: "pass"}
	//	err = Ormer.Read(u, "UserName")
	Read(md interface{}, cols ...string) error
	// Like Read(), but with "FOR UPDATE" clause, useful in transaction.
	// Some databases are not support this feature.
	ReadForUpdate(md interface{}, cols ...string) error
	// Try to read a row from the database, or insert one if it doesn't exist
	ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error)
	// insert model data to database
	// for example:
	//  user := new(User)
	//  id, err = Ormer.Insert(user)
	//  user must be a pointer and Insert will set user's pk field
	Insert(interface{}) (int64, error)
	// mysql:InsertOrUpdate(model) or InsertOrUpdate(model,"colu=colu+value")
	// if colu type is integer : can use(+-*/), string : convert(colu,"value")
	// postgres: InsertOrUpdate(model,"conflictColumnName") or InsertOrUpdate(model,"conflictColumnName","colu=colu+value")
	// if colu type is integer : can use(+-*/), string : colu || "value"
	InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error)
	// insert some models to database
	InsertMulti(bulk int, mds interface{}) (int64, error)
	// update model to database.
	// cols set the columns those want to update.
	// find model by Id(pk) field and update columns specified by fields, if cols is null then update all columns
	// for example:
	// user := User{Id: 2}
	//	user.Langs = append(user.Langs, "zh-CN", "en-US")
	//	user.Extra.Name = "beego"
	//	user.Extra.Data = "orm"
	//	num, err = Ormer.Update(&user, "Langs", "Extra")
	Update(md interface{}, cols ...string) (int64, error)
	// delete model in database
	Delete(md interface{}, cols ...string) (int64, error)
	// load related models to md model.
	// args are limit, offset int and order string.
	//
	// example:
	// 	Ormer.LoadRelated(post,"Tags")
	// 	for _,tag := range post.Tags{...}
	// args[0] bool true useDefaultRelsDepth ; false  depth 0
	// args[0] int  loadRelationDepth
	// args[1] int limit default limit 1000
	// args[2] int offset default offset 0
	// args[3] string order  for example : "-Id"
	// make sure the relation is defined in model struct tags.
	LoadRelated(md interface{}, name string, args ...interface{}) (int64, error)
	// create a models to models queryer
	// for example:
	// 	post := Post{Id: 4}
	// 	m2m := Ormer.QueryM2M(&post, "Tags")
	QueryM2M(md interface{}, name string) QueryM2Mer
	// return a QuerySeter for table operations.
	// table name can be string or struct.
	// e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil)),
	QueryTable(ptrStructOrTableName interface{}) QuerySeter
	// switch to another registered database driver by given name.
	Using(name string) error
	// begin transaction
	// for example:
	// 	o := NewOrm()
	// 	err := o.Begin()
	// 	...
	// 	err = o.Rollback()
	Begin() error
	// begin transaction with provided context and option
	// the provided context is used until the transaction is committed or rolled back.
	// if the context is canceled, the transaction will be rolled back.
	// the provided TxOptions is optional and may be nil if defaults should be used.
	// if a non-default isolation level is used that the driver doesn't support, an error will be returned.
	// for example:
	//  o := NewOrm()
	// 	err := o.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	//  ...
	//  err = o.Rollback()
	BeginTx(ctx context.Context, opts *sql.TxOptions) error
	// commit transaction
	Commit() error
	// rollback transaction
	Rollback() error
	// return a raw query seter for raw sql string.
	// for example:
	//	 ormer.Raw("UPDATE `user` SET `user_name` = ? WHERE `user_name` = ?", "slene", "testing").Exec()
	//	// update user testing's name to slene
	Raw(query string, args ...interface{}) RawSeter
	Driver() Driver
	DBStats() *sql.DBStats
}

// Inserter insert prepared statement
type Inserter orm.Inserter

// QuerySeter query seter
type QuerySeter orm.QuerySeter

// QueryM2Mer model to model query struct
// all operations are on the m2m table only, will not affect the origin model table
type QueryM2Mer orm.QueryM2Mer

// RawPreparer raw query statement
type RawPreparer orm.RawPreparer

// RawSeter raw query seter
// create From Ormer.Raw
// for example:
//  sql := fmt.Sprintf("SELECT %sid%s,%sname%s FROM %suser%s WHERE id = ?",Q,Q,Q,Q,Q,Q)
//  rs := Ormer.Raw(sql, 1)
type RawSeter orm.RawSeter
