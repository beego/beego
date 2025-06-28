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

package valuer

import (
	"database/sql"
	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/qb/errs"
	"reflect"
)

type reflectValue struct {
	val  reflect.Value
	meta *models.ModelInfo
}

func NewReflectValue(t any, model *models.ModelInfo) Value {
	return &reflectValue{
		val:  reflect.ValueOf(t).Elem(),
		meta: model,
	}
}

func (r *reflectValue) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cs) > len(r.meta.Fields.Columns) {
		return errs.ErrTooManyColumns
	}

	colValues := make([]interface{}, len(cs))
	colEleValues := make([]reflect.Value, len(cs))

	for i, c := range cs {
		cm, ok := r.meta.Fields.Columns[c]
		if !ok {
			return errs.NewErrUnknownField(c)
		}
		val := reflect.New(cm.AddrType)
		colValues[i] = val.Interface()
		colEleValues[i] = val.Elem()
	}

	if err = rows.Scan(colValues...); err != nil {
		return err
	}

	for i, c := range cs {
		cm := r.meta.Fields.Columns[c]
		fd, _ := r.fieldByIndex(cm.Name)
		fd.Set(colEleValues[i])
	}
	return nil
}

// Field returns the field value by name.
func (r *reflectValue) Field(name string) (reflect.Value, error) {
	res, ok := r.fieldByIndex(name)
	if !ok {
		return reflect.Value{}, errs.NewErrUnknownField(name)
	}
	return res, nil
}

func (r *reflectValue) fieldByIndex(name string) (reflect.Value, bool) {
	fd, ok := r.meta.Fields.Fields[name]
	if !ok {
		return reflect.Value{}, false
	}
	value := r.val
	for _, i := range fd.FieldIndex {
		value = value.Field(i)
	}
	return value, true
}
