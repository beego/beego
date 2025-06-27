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

package qb

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/internal/buffers"
	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/qb/errs"
	"github.com/beego/beego/v2/client/orm/qb/valuer"
	"reflect"
)

// Updater is the builder responsible for building UPDATE query
type Updater[T any] struct {
	builder
	val           valuer.Value
	where         []Predicate
	assigns       []Assignable
	table         interface{}
	sess          orm.QueryExecutor
	registry      *models.ModelCache
	valCreator    valuer.Creator
	ignoreNilVal  bool
	ignoreZeroVal bool
}

func (u *Updater[T]) Build() (*Query, error) {
	defer buffers.Put(u.buffer)
	var (
		t   T
		err error
	)
	if u.table == nil {
		u.table = &t
	}
	u.model, err = u.registry.GetOrRegisterByMd(&t)
	if err != nil {
		return nil, err
	}
	u.val = u.valCreator(u.table, u.model)
	u.args = make([]interface{}, 0, len(u.model.Fields.Columns))

	u.writeString("UPDATE ")
	u.buildTable()
	u.writeString(" SET ")
	if len(u.assigns) == 0 {
		err = u.buildDefaultColumns()
	} else {
		err = u.buildAssigns()
	}
	if err != nil {
		return nil, err
	}

	if len(u.where) > 0 {
		u.writeString(" WHERE ")
		err = u.buildPredicates(u.where)
		if err != nil {
			return nil, err
		}
	}

	u.end()
	return &Query{
		SQL:  u.buffer.String(),
		Args: u.args,
	}, nil
}

func (u *Updater[T]) buildDefaultColumns() error {
	has := false
	for _, fi := range u.model.Fields.FieldsDB {
		refVal, _ := u.val.Field(fi.Name)
		if has {
			_ = u.buffer.WriteByte(',')
		}
		u.writeByte('`')
		u.writeString(fi.Column)
		u.writeByte('`')
		_ = u.buffer.WriteByte('=')
		u.parameter(refVal.Interface())
		has = true
	}
	if !has {
		return errs.NewValueNotSetError()
	}
	return nil
}

func (u *Updater[T]) buildTable() {
	if u.model.Table == "" {
		var t T
		typ := reflect.TypeOf(t)
		u.writeByte('`')
		u.writeString(typ.Name())
		u.writeByte('`')
	} else {
		u.writeByte('`')
		u.writeString(u.model.Table)
		u.writeByte('`')
	}
}

func (u *Updater[T]) buildAssigns() error {
	has := false
	for _, assign := range u.assigns {
		if has {
			u.comma()
		}
		switch a := assign.(type) {
		case Column:
			fmt.Print(a.name)
			c, ok := u.model.Fields.Fields[a.name]
			if !ok {
				return errs.NewErrUnknownField(a.name)
			}
			refVal, _ := u.val.Field(a.name)
			u.writeByte('`')
			u.writeString(c.Column)
			u.writeByte('`')
			_ = u.buffer.WriteByte('=')
			u.parameter(refVal.Interface())
			has = true
		case columns:
			for _, name := range a.cs {
				c, ok := u.model.Fields.Fields[name]
				if !ok {
					return errs.NewErrUnknownField(name)
				}
				refVal, _ := u.val.Field(name)
				if has {
					u.comma()
				}
				u.writeByte('`')
				u.writeString(c.Column)
				u.writeByte('`')
				_ = u.buffer.WriteByte('=')
				u.parameter(refVal.Interface())
				has = true
			}
		case Assignment:
			if err := u.buildExpression(binaryExpr(a)); err != nil {
				return err
			}
			has = true
		default:
			return errs.ErrUnsupportedAssignment
		}
	}
	if !has {
		return errs.NewValueNotSetError()
	}
	return nil
}

func (u *Updater[T]) Update(val *T) *Updater[T] {
	u.table = val
	return u
}

// Set represents SET clause
func (u *Updater[T]) Set(assigns ...Assignable) *Updater[T] {
	u.assigns = assigns
	return u
}

// Where represents WHERE clause
func (u *Updater[T]) Where(predicates ...Predicate) *Updater[T] {
	u.where = predicates
	return u
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	q, err := u.Build()
	if err != nil {
		return Result{err: err}
	}
	t := new(T)
	res, err := u.sess.ExecRaw(ctx, t, q.SQL, q.Args...)
	return Result{res: res, err: err}
}

func NewUpdater[T any](sess orm.QueryExecutor) *Updater[T] {
	return &Updater[T]{
		sess: sess,
		builder: builder{
			buffer: buffers.Get(),
		},
		registry:   models.DefaultModelCache,
		valCreator: valuer.NewReflectValue,
	}
}
