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

	"github.com/beego/beego/v2/client/orm"

	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/qb/errs"

	"reflect"
	"strings"
)

// The Selector is used to construct a SELECT statement
type Selector[T any] struct {
	sb        strings.Builder
	args      []any
	orderBy   []Column
	where     []Predicate
	offset    int
	limit     int
	model     *models.ModelInfo
	columns   []Selectable
	tableName string

	cache *models.ModelCache
	db    orm.Ormer
}

func NewSelector[T any](cache *models.ModelCache) *Selector[T] {
	return &Selector[T]{
		cache: cache,
	}
}
func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)

	s.model, _ = s.cache.GetByMd(&t)
	if s.model == nil {
		//orm.BootStrap()
		err = s.cache.Register("", true, &t)
		if err != nil {
			return nil, err
		}
		s.model, _ = s.cache.GetByMd(&t)
	}

	s.sb.WriteString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")
	if s.tableName != "" {
		if s.tableName[0] == '`' && s.tableName[len(s.tableName)-1] == '`' {
			s.sb.WriteString(s.tableName)
		} else {
			s.sb.WriteByte('`')
			s.sb.WriteString(s.tableName)
			s.sb.WriteByte('`')
		}
	} else {
		if s.model.Table == "" {
			typ := reflect.TypeOf(t)
			s.sb.WriteByte('`')
			s.sb.WriteString(typ.Name())
			s.sb.WriteByte('`')
		} else {
			s.sb.WriteByte('`')
			s.sb.WriteString(s.model.Table)
			s.sb.WriteByte('`')
		}
	}
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}
	if len(s.orderBy) > 0 {
		s.sb.WriteString(" ORDER BY ")
		for i, c := range s.orderBy {
			if i > 0 {
				s.sb.WriteByte(',')
			}
			if err = s.buildColumn(c, false); err != nil {
				return nil, err
			}
		}
	}
	if s.limit > 0 {
		s.sb.WriteString(" LIMIT ?")
		s.addArgs(s.limit)
	}

	if s.offset > 0 {
		s.sb.WriteString(" OFFSET ?")
		s.addArgs(s.offset)
	}

	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}
func (s *Selector[T]) addArgs(args ...any) {
	if s.args == nil {
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, args...)
}

func (s *Selector[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(exp.name)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.args = append(s.args, exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}
	default:
		return errs.NewErrUnsupportedExpressionType(exp)
	}
	return nil
}
func (s *Selector[T]) buildColumn(c Column, useAlias bool) error {
	s.sb.WriteByte('`')
	fd, ok := s.model.Fields.Fields[c.name]
	if !ok {
		return errs.NewErrUnknownField(c.name)
	}
	s.sb.WriteString(fd.Column)
	s.sb.WriteByte('`')
	if c.order != "" {
		s.sb.WriteString(c.order)
	}
	if useAlias {
		s.buildAs(c.alias)
	}
	return nil
}
func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.sb.WriteByte('*')
		return nil
	}
	for i, c := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
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
			s.sb.WriteString(val.raw)
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
	s.sb.WriteString(a.fn)
	s.sb.WriteString("(`")
	fd, ok := s.model.Fields.Fields[a.arg]
	if !ok {
		return errs.NewErrUnknownField(a.arg)
	}
	s.sb.WriteString(fd.Column)
	s.sb.WriteString("`)")
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
		s.sb.WriteString(" AS ")
		s.sb.WriteByte('`')
		s.sb.WriteString(alias)
		s.sb.WriteByte('`')
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
