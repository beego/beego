// Copyright 2023 beego. All Rights Reserved.
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

package qb

import (
	"context"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/valyala/bytebufferpool"
	"reflect"
)

var _ QueryBuilder = &Deleter[any]{}

// Deleter builds DELETE query
type Deleter[T any] struct {
	builder
	table interface{}
	db    orm.Ormer
	where []Predicate
}

// NewDeleter 开始构建一个 DELETE 查询
func NewDeleter[T any](db orm.Ormer) *Deleter[T] {
	return &Deleter[T]{
		db: db,
		builder: builder{
			buffer: bytebufferpool.Get(),
		},
	}
}

func (d *Deleter[T]) Build() (*Query, error) {
	defer bytebufferpool.Put(d.buffer)
	d.writeString("DELETE FROM ")
	var err error
	if d.table == nil {
		d.table = new(T)
	}
	registry := models.DefaultModelCache
	d.model, _ = registry.GetByMd(d.table)
	if d.model == nil {
		err = registry.Register("", true, d.table)
		if err != nil {
			return nil, err
		}
		d.model, _ = registry.GetByMd(d.table)
	}
	if d.model.Table == "" {
		typ := reflect.TypeOf(d.table)
		d.writeByte('`')
		d.writeString(typ.Name())
		d.writeByte('`')
	} else {
		d.writeByte('`')
		d.writeString(d.model.Table)
		d.writeByte('`')
	}
	if len(d.where) > 0 {
		d.writeString(" WHERE ")
		err = d.buildPredicates(d.where)
		if err != nil {
			return nil, err
		}
	}
	d.end()
	return &Query{SQL: d.buffer.String(), Args: d.args}, nil
}

// From accepts model definition
func (d *Deleter[T]) From(table interface{}) *Deleter[T] {
	d.table = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.where = predicates
	return d
}

// Exec sql
func (d *Deleter[T]) Exec(ctx context.Context) Result {
	q, err := d.Build()
	if err != nil {
		return Result{err: err}
	}
	t := new(T)
	res, err := d.db.ExecRaw(ctx, t, q.SQL, q.Args...)
	return Result{res: res, err: err}
}
