package qb

import (
	"context"
	"github.com/beego/beego/v2/client/orm"
)

type Selector[T any] struct {
	db orm.Ormer
}

func NewSelector[T any](db orm.Ormer) *Selector[T] {
	panic("implement me")
}

func (s *Selector[T]) Build() (Query, error) {
	panic("implement me")
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
