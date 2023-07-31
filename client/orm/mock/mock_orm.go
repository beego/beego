// Copyright 2020 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mock

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/beego/beego/v2/client/orm"
)

func init() {
	RegisterMockDB("default")
}

// RegisterMockDB create an "virtual DB" by using sqllite
// you should not
func RegisterMockDB(name string) {
	source := filepath.Join(os.TempDir(), name+".db")
	_ = orm.RegisterDataBase(name, "sqlite3", source)
}

// MockTable only check table name
func MockTable(tableName string, resp ...interface{}) *Mock {
	return NewMock(NewSimpleCondition(tableName, ""), resp, nil)
}

// MockMethod only check method name
func MockMethod(method string, resp ...interface{}) *Mock {
	return NewMock(NewSimpleCondition("", method), resp, nil)
}

// MockRead support orm.Read and orm.ReadWithCtx
// cb is used to mock read data from DB
func MockRead(tableName string, cb func(data interface{}), err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "ReadWithCtx"), []interface{}{err}, func(inv *orm.Invocation) {
		if cb != nil {
			cb(inv.Args[0])
		}
	})
}

// MockReadForUpdateWithCtx support ReadForUpdate and ReadForUpdateWithCtx
// cb is used to mock read data from DB
func MockReadForUpdateWithCtx(tableName string, cb func(data interface{}), err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "ReadForUpdateWithCtx"),
		[]interface{}{err},
		func(inv *orm.Invocation) {
			cb(inv.Args[0])
		})
}

// MockReadOrCreateWithCtx support ReadOrCreate and ReadOrCreateWithCtx
// cb is used to mock read data from DB
func MockReadOrCreateWithCtx(tableName string,
	cb func(data interface{}),
	insert bool, id int64, err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "ReadOrCreateWithCtx"),
		[]interface{}{insert, id, err},
		func(inv *orm.Invocation) {
			cb(inv.Args[0])
		})
}

// MockInsertWithCtx support Insert and InsertWithCtx
func MockInsertWithCtx(tableName string, id int64, err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "InsertWithCtx"), []interface{}{id, err}, nil)
}

// MockInsertMultiWithCtx support InsertMulti and InsertMultiWithCtx
func MockInsertMultiWithCtx(tableName string, cnt int64, err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "InsertMultiWithCtx"), []interface{}{cnt, err}, nil)
}

// MockInsertOrUpdateWithCtx support InsertOrUpdate and InsertOrUpdateWithCtx
func MockInsertOrUpdateWithCtx(tableName string, id int64, err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "InsertOrUpdateWithCtx"), []interface{}{id, err}, nil)
}

// MockUpdateWithCtx support UpdateWithCtx and Update
func MockUpdateWithCtx(tableName string, affectedRow int64, err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "UpdateWithCtx"), []interface{}{affectedRow, err}, nil)
}

// MockDeleteWithCtx support Delete and DeleteWithCtx
func MockDeleteWithCtx(tableName string, affectedRow int64, err error) *Mock {
	return NewMock(NewSimpleCondition(tableName, "DeleteWithCtx"), []interface{}{affectedRow, err}, nil)
}

// MockQueryM2MWithCtx support QueryM2MWithCtx and QueryM2M
// Now you may be need to use golang/mock to generate QueryM2M mock instance
// Or use DoNothingQueryM2Mer
// for example:
//
//	post := Post{Id: 4}
//	m2m := Ormer.QueryM2M(&post, "Tags")
//
// when you write test code:
// MockQueryM2MWithCtx("post", "Tags", mockM2Mer)
// "post" is the table name of model Post structure
// TODO provide orm.QueryM2Mer
func MockQueryM2MWithCtx(tableName string, name string, res orm.QueryM2Mer) *Mock {
	return NewMock(NewQueryM2MerCondition(tableName, name), []interface{}{res}, nil)
}

// MockLoadRelatedWithCtx support LoadRelatedWithCtx and LoadRelated
func MockLoadRelatedWithCtx(tableName string, name string, rows int64, err error) *Mock {
	return NewMock(NewQueryM2MerCondition(tableName, name), []interface{}{rows, err}, nil)
}

// MockQueryTableWithCtx support QueryTableWithCtx and QueryTable
func MockQueryTableWithCtx(tableName string, qs orm.QuerySeter) *Mock {
	return NewMock(NewSimpleCondition(tableName, "QueryTable"), []interface{}{qs}, nil)
}

// MockRawWithCtx support RawWithCtx and Raw
func MockRawWithCtx(rs orm.RawSeter) *Mock {
	return NewMock(NewSimpleCondition("", "RawWithCtx"), []interface{}{rs}, nil)
}

// MockDriver support Driver
// func MockDriver(driver orm.Driver) *Mock {
// 	return NewMock(NewSimpleCondition("", "Driver"), []interface{}{driver})
// }

// MockDBStats support DBStats
func MockDBStats(stats *sql.DBStats) *Mock {
	return NewMock(NewSimpleCondition("", "DBStats"), []interface{}{stats}, nil)
}

// MockBeginWithCtxAndOpts support Begin, BeginWithCtx, BeginWithOpts, BeginWithCtxAndOpts
// func MockBeginWithCtxAndOpts(txOrm *orm.TxOrmer, err error) *Mock {
// 	return NewMock(NewSimpleCondition("", "BeginWithCtxAndOpts"), []interface{}{txOrm, err})
// }

// MockDoTxWithCtxAndOpts support DoTx, DoTxWithCtx, DoTxWithOpts, DoTxWithCtxAndOpts
// func MockDoTxWithCtxAndOpts(txOrm *orm.TxOrmer, err error) *Mock {
// 	return MockBeginWithCtxAndOpts(txOrm, err)
// }

// MockCommit support Commit
func MockCommit(err error) *Mock {
	return NewMock(NewSimpleCondition("", "Commit"), []interface{}{err}, nil)
}

// MockRollback support Rollback
func MockRollback(err error) *Mock {
	return NewMock(NewSimpleCondition("", "Rollback"), []interface{}{err}, nil)
}

// MockRollbackUnlessCommit support RollbackUnlessCommit
func MockRollbackUnlessCommit(err error) *Mock {
	return NewMock(NewSimpleCondition("", "RollbackUnlessCommit"), []interface{}{err}, nil)
}
