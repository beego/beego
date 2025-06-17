package valuer

import (
	"database/sql"
	"github.com/beego/beego/v2/client/orm/internal/models"
	"reflect"
)

type Value interface {
	Field(name string) (reflect.Value, error)
	SetColumns(rows *sql.Rows) error
}

type Creator func(val any, meta *models.ModelInfo) Value
