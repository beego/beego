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
	"fmt"
	"github.com/beego/beego/v2/client/orm/qb/errs"
	"reflect"
	"time"
)

var (
	_ Dialect = (*mysqlDialect)(nil)
	_ Dialect = (*sqlite3Dialect)(nil)
)

type Dialect interface {
	Name() string
	quoter() byte
	ColTypeOf(typ reflect.Value) string
}

type mysqlDialect struct{}

func (d *mysqlDialect) Name() string {
	return "MySQL"
}

func NewMySQLDialect() Dialect {
	return &mysqlDialect{}
}

func (d *mysqlDialect) quoter() byte {
	return '`'
}

func (d *mysqlDialect) ColTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "int(11)"
	case reflect.Int64, reflect.Uint64:
		return "bigint(11)"
	case reflect.Float32, reflect.Float64:
		return "float(11)"
	case reflect.String:
		return "longtext"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

type sqlite3Dialect struct{}

func (d *sqlite3Dialect) Name() string {
	return "SQLite"
}

func NewSqlite3Dialect() Dialect {
	return &sqlite3Dialect{}
}

func (d *sqlite3Dialect) quoter() byte {
	return '`'
}

func (d *sqlite3Dialect) ColTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func Of(driver string) (Dialect, error) {
	switch driver {
	case "sqlite3":
		return NewSqlite3Dialect(), nil
	case "mysql":
		return NewMySQLDialect(), nil
	default:
		return nil, errs.NewUnsupportedDriverError(driver)
	}
}
