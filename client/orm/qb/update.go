package qb

import (
	"context"
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
	var err error
	t := new(T)
	if u.table == nil {
		u.table = t
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
	for _, c := range u.model.Fields.Columns {
		refVal, _ := u.val.Field(c.Name)
		if u.ignoreZeroVal && isZeroValue(refVal) {
			continue
		}
		if u.ignoreNilVal && isNilValue(refVal) {
			continue
		}
		if has {
			_ = u.buffer.WriteByte(',')
		}
		u.writeByte('`')
		u.writeString(c.Column)
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

func (u *Updater[T]) SkipNilValue() *Updater[T] {
	u.ignoreNilVal = true
	return u
}

func isNilValue(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return val.IsNil()
	}
	return false
}

func (u *Updater[T]) SkipZeroValue() *Updater[T] {
	u.ignoreZeroVal = true
	return u
}

func isZeroValue(val reflect.Value) bool {
	return val.IsZero()
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
