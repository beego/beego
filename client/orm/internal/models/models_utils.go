// Copyright 2014 beego Author. All Rights Reserved.
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

package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// 1 is attr
// 2 is tag
var supportTag = map[string]int{
	"-":            1,
	"null":         1,
	"index":        1,
	"unique":       1,
	"pk":           1,
	"auto":         1,
	"auto_now":     1,
	"auto_now_add": 1,
	"size":         2,
	"column":       2,
	"default":      2,
	"rel":          2,
	"reverse":      2,
	"rel_table":    2,
	"rel_through":  2,
	"digits":       2,
	"decimals":     2,
	"on_delete":    2,
	"type":         2,
	"description":  2,
	"precision":    2,
	"db_type":      2,
}

type fn func(string) string

var (
	NameStrategyMap = map[string]fn{
		DefaultNameStrategy:      SnakeString,
		SnakeAcronymNameStrategy: SnakeStringWithAcronym,
	}
	DefaultNameStrategy      = "snakeString"
	SnakeAcronymNameStrategy = "snakeStringWithAcronym"
	NameStrategy             = DefaultNameStrategy
	defaultStructTagDelim    = ";"
	DefaultStructTagName     = "orm"
)

// GetFullName get reflect.Type name with package path.
func GetFullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}

// GetTableName get struct table name.
// If the struct implement the TableName, then get the result as tablename
// else use the struct name which will apply snakeString.
func GetTableName(val reflect.Value) string {
	if fun := val.MethodByName("TableName"); fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		// has return and the first val is string
		if len(vals) > 0 && vals[0].Kind() == reflect.String {
			return vals[0].String()
		}
	}
	return SnakeString(reflect.Indirect(val).Type().Name())
}

// GetTableEngine get table engine, myisam or innodb.
func GetTableEngine(val reflect.Value) string {
	fun := val.MethodByName("TableEngine")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 && vals[0].Kind() == reflect.String {
			return vals[0].String()
		}
	}
	return ""
}

// GetTableIndex get table index from method.
func GetTableIndex(val reflect.Value) [][]string {
	fun := val.MethodByName("TableIndex")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 && vals[0].CanInterface() {
			if d, ok := vals[0].Interface().([][]string); ok {
				return d
			}
		}
	}
	return nil
}

// GetTableUnique get table unique from method
func GetTableUnique(val reflect.Value) [][]string {
	fun := val.MethodByName("TableUnique")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 && vals[0].CanInterface() {
			if d, ok := vals[0].Interface().([][]string); ok {
				return d
			}
		}
	}
	return nil
}

// IsApplicableTableForDB get whether the table needs to be created for the database alias
func IsApplicableTableForDB(val reflect.Value, db string) bool {
	if !val.IsValid() {
		return true
	}
	fun := val.MethodByName("IsApplicableTableForDB")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{reflect.ValueOf(db)})
		if len(vals) > 0 && vals[0].Kind() == reflect.Bool {
			return vals[0].Bool()
		}
	}
	return true
}

// get snaked column name
func getColumnName(ft int, addrField reflect.Value, sf reflect.StructField, col string) string {
	column := col
	if col == "" {
		column = NameStrategyMap[NameStrategy](sf.Name)
	}
	switch ft {
	case RelForeignKey, RelOneToOne:
		if len(col) == 0 {
			column = column + "_id"
		}
	case RelManyToMany, RelReverseMany, RelReverseOne:
		column = sf.Name
	}
	return column
}

// return field type as type constant from reflect.Value
func getFieldType(val reflect.Value) (ft int, err error) {
	switch val.Type() {
	case reflect.TypeOf(new(int8)):
		ft = TypeBitField
	case reflect.TypeOf(new(int16)):
		ft = TypeSmallIntegerField
	case reflect.TypeOf(new(int32)),
		reflect.TypeOf(new(int)):
		ft = TypeIntegerField
	case reflect.TypeOf(new(int64)):
		ft = TypeBigIntegerField
	case reflect.TypeOf(new(uint8)):
		ft = TypePositiveBitField
	case reflect.TypeOf(new(uint16)):
		ft = TypePositiveSmallIntegerField
	case reflect.TypeOf(new(uint32)),
		reflect.TypeOf(new(uint)):
		ft = TypePositiveIntegerField
	case reflect.TypeOf(new(uint64)):
		ft = TypePositiveBigIntegerField
	case reflect.TypeOf(new(float32)),
		reflect.TypeOf(new(float64)):
		ft = TypeFloatField
	case reflect.TypeOf(new(bool)):
		ft = TypeBooleanField
	case reflect.TypeOf(new(string)):
		ft = TypeVarCharField
	case reflect.TypeOf(new(time.Time)):
		ft = TypeDateTimeField
	default:
		elm := reflect.Indirect(val)
		switch elm.Kind() {
		case reflect.Int8:
			ft = TypeBitField
		case reflect.Int16:
			ft = TypeSmallIntegerField
		case reflect.Int32, reflect.Int:
			ft = TypeIntegerField
		case reflect.Int64:
			ft = TypeBigIntegerField
		case reflect.Uint8:
			ft = TypePositiveBitField
		case reflect.Uint16:
			ft = TypePositiveSmallIntegerField
		case reflect.Uint32, reflect.Uint:
			ft = TypePositiveIntegerField
		case reflect.Uint64:
			ft = TypePositiveBigIntegerField
		case reflect.Float32, reflect.Float64:
			ft = TypeFloatField
		case reflect.Bool:
			ft = TypeBooleanField
		case reflect.String:
			ft = TypeVarCharField
		default:
			if elm.Interface() == nil {
				panic(fmt.Errorf("%s is nil pointer, may be miss setting tag", val))
			}
			switch elm.Interface().(type) {
			case sql.NullInt64:
				ft = TypeBigIntegerField
			case sql.NullFloat64:
				ft = TypeFloatField
			case sql.NullBool:
				ft = TypeBooleanField
			case sql.NullString:
				ft = TypeVarCharField
			case time.Time:
				ft = TypeDateTimeField
			}
		}
	}
	if ft&IsFieldType == 0 {
		err = fmt.Errorf("unsupport field type %s, may be miss setting tag", val)
	}
	return
}

// ParseStructTag parse struct tag string
func ParseStructTag(data string) (attrs map[string]bool, tags map[string]string) {
	attrs = make(map[string]bool)
	tags = make(map[string]string)
	for _, v := range strings.Split(data, defaultStructTagDelim) {
		if v == "" {
			continue
		}
		v = strings.TrimSpace(v)
		if t := strings.ToLower(v); supportTag[t] == 1 {
			attrs[t] = true
		} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
			name := t[:i]
			if supportTag[name] == 2 {
				v = v[i+1 : len(v)-1]
				tags[name] = v
			}
		}
	}
	return
}

func SnakeStringWithAcronym(s string) string {
	data := make([]byte, 0, len(s)*2)
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		before := false
		after := false
		if i > 0 {
			before = s[i-1] >= 'a' && s[i-1] <= 'z'
		}
		if i+1 < num {
			after = s[i+1] >= 'a' && s[i+1] <= 'z'
		}
		if i > 0 && d >= 'A' && d <= 'Z' && (before || after) {
			data = append(data, '_')
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data))
}

// SnakeString snake string, XxYy to xx_yy , XxYY to xx_y_y
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data))
}

// CamelString camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data)
}

const (
	OdCascade    = "cascade"
	OdSetNULL    = "set_null"
	OdSetDefault = "set_default"
	OdDoNothing  = "do_nothing"
)
