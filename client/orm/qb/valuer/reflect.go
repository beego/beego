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
	return nil
}

// Field 返回字段值
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
