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
	"fmt"
	"reflect"

	"github.com/beego/beego/v2/client/orm/internal/models"
)

// an insert queryer struct
type insertSet struct {
	mi     *models.ModelInfo
	orm    *ormBase
	stmt   stmtQuerier
	closed bool
}

var _ Inserter = new(insertSet)

// insert model ignore it's registered or not.
func (o *insertSet) Insert(md interface{}) (int64, error) {
	return o.InsertWithCtx(context.Background(), md)
}

func (o *insertSet) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	if o.closed {
		return 0, ErrStmtClosed
	}
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	name := models.GetFullName(typ)
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Inserter.Insert> cannot use non-ptr model struct `%s`", name))
	}
	if name != o.mi.FullName {
		panic(fmt.Errorf("<Inserter.Insert> need model `%s` but found `%s`", o.mi.FullName, name))
	}
	id, err := o.orm.alias.DbBaser.InsertStmt(ctx, o.stmt, o.mi, ind, o.orm.alias.TZ)
	if err != nil {
		return id, err
	}
	if id > 0 {
		if o.mi.Fields.Pk.Auto {
			if o.mi.Fields.Pk.FieldType&IsPositiveIntegerField > 0 {
				ind.FieldByIndex(o.mi.Fields.Pk.FieldIndex).SetUint(uint64(id))
			} else {
				ind.FieldByIndex(o.mi.Fields.Pk.FieldIndex).SetInt(id)
			}
		}
	}
	return id, nil
}

// close insert queryer statement
func (o *insertSet) Close() error {
	if o.closed {
		return ErrStmtClosed
	}
	o.closed = true
	return o.stmt.Close()
}

// create new insert queryer.
func newInsertSet(ctx context.Context, orm *ormBase, mi *models.ModelInfo) (Inserter, error) {
	bi := new(insertSet)
	bi.orm = orm
	bi.mi = mi
	st, query, err := orm.alias.DbBaser.PrepareInsert(ctx, orm.db, mi)
	if err != nil {
		return nil, err
	}
	if Debug {
		bi.stmt = newStmtQueryLog(orm.alias, st, query)
	} else {
		bi.stmt = st
	}
	return bi, nil
}
