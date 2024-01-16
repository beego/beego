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
	"errors"
	"reflect"

	"github.com/valyala/bytebufferpool"

	"github.com/beego/beego/v2/client/orm"

	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/qb/errs"
)

var _ QueryBuilder = &Selector[any]{}

// The Selector is used to construct a SELECT statement
type Selector[T any] struct {
	builder
	orderBy   []Column
	where     []Predicate
	offset    int
	limit     int
	columns   []Selectable
	tableName string
	db        orm.Ormer
}

func NewSelector[T any](db orm.Ormer) *Selector[T] {
	return &Selector[T]{
		db: db,
		builder: builder{
			buffer: bytebufferpool.Get(),
		},
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	defer bytebufferpool.Put(s.buffer)
	registry := models.DefaultModelCache
	s.model, _ = registry.GetByMd(&t)
	if s.model == nil {
		//orm.BootStrap()
		err = registry.Register("", true, &t)
		if err != nil {
			return nil, err
		}
		s.model, _ = registry.GetByMd(&t)
	}
	s.writeString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.writeString(" FROM ")
	s.buildTable()
	if len(s.where) > 0 {
		s.writeString(" WHERE ")
		err = s.buildPredicates(s.where)
		if err != nil {
			return nil, err
		}
	}
	if len(s.orderBy) > 0 {
		s.writeString(" ORDER BY ")
		for i, c := range s.orderBy {
			if i > 0 {
				s.comma()
			}
			if err = s.buildColumn(c, false); err != nil {
				return nil, err
			}
		}
	}
	if s.limit > 0 {
		s.writeString(" LIMIT ?")
		s.addArgs(s.limit)
	}
	if s.offset > 0 {
		s.writeString(" OFFSET ?")
		s.addArgs(s.offset)
	}

	s.end()
	return &Query{
		SQL:  s.buffer.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildTable() {
	if s.tableName != "" {
		if s.tableName[0] == '`' && s.tableName[len(s.tableName)-1] == '`' {
			s.writeString(s.tableName)
		} else {
			s.writeByte('`')
			s.writeString(s.tableName)
			s.writeByte('`')
		}
	} else {
		if s.model.Table == "" {
			var t T
			typ := reflect.TypeOf(t)
			s.writeByte('`')
			s.writeString(typ.Name())
			s.writeByte('`')
		} else {
			s.writeByte('`')
			s.writeString(s.model.Table)
			s.writeByte('`')
		}
	}
}

func (s *Selector[T]) addArgs(args ...any) {
	if s.args == nil {
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, args...)
}

func (s *Selector[T]) buildColumn(c Column, useAlias bool) error {
	s.writeByte('`')
	fd, ok := s.model.Fields.Fields[c.name]
	if !ok {
		return errs.NewErrUnknownField(c.name)
	}
	s.writeString(fd.Column)
	s.writeByte('`')
	if c.order != "" {
		s.writeString(c.order)
	}
	if useAlias {
		s.buildAs(c.alias)
	}
	return nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.writeByte('*')
		return nil
	}
	for i, c := range s.columns {
		if i > 0 {
			s.comma()
		}
		switch val := c.(type) {
		case Column:
			if err := s.buildColumn(val, true); err != nil {
				return err
			}
		case Aggregate:
			if err := s.buildAggregate(val, true); err != nil {
				return err
			}
		case RawExpr:
			s.writeString(val.raw)
			if len(val.args) != 0 {
				s.addArgs(val.args...)
			}
		default:
			return errors.New("orm: Unsupported target column")
		}
	}
	return nil
}

func (s *Selector[T]) buildAggregate(a Aggregate, useAlias bool) error {
	s.writeString(a.fn)
	s.writeString("(`")
	fd, ok := s.model.Fields.Fields[a.arg]
	if !ok {
		return errs.NewErrUnknownField(a.arg)
	}
	s.writeString(fd.Column)
	s.writeString("`)")
	if useAlias {
		s.buildAs(a.alias)
	}
	return nil
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.tableName = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) OrderBy(cols ...Column) *Selector[T] {
	s.orderBy = cols
	return s
}

func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}

func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

type Selectable interface {
	selectable()
}

func (s *Selector[T]) buildAs(alias string) {
	if alias != "" {
		s.writeString(" AS ")
		s.writeByte('`')
		s.writeString(alias)
		s.writeByte('`')
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	t := new(T)
	err = s.db.ReadRaw(ctx, t, q.SQL, q.Args...)
	return t, nil
}
