// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie, slene

package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// get reflect.Type name with package path.
func getFullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}

// get table name. method, or field name. auto snaked.
func getTableName(val reflect.Value) string {
	ind := reflect.Indirect(val)
	fun := val.MethodByName("TableName")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 {
			val := vals[0]
			if val.Kind() == reflect.String {
				return val.String()
			}
		}
	}
	return snakeString(ind.Type().Name())
}

// get table engine, mysiam or innodb.
func getTableEngine(val reflect.Value) string {
	fun := val.MethodByName("TableEngine")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 {
			val := vals[0]
			if val.Kind() == reflect.String {
				return val.String()
			}
		}
	}
	return ""
}

// get table index from method.
func getTableIndex(val reflect.Value) [][]string {
	fun := val.MethodByName("TableIndex")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 {
			val := vals[0]
			if val.CanInterface() {
				if d, ok := val.Interface().([][]string); ok {
					return d
				}
			}
		}
	}
	return nil
}

// get table unique from method
func getTableUnique(val reflect.Value) [][]string {
	fun := val.MethodByName("TableUnique")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 {
			val := vals[0]
			if val.CanInterface() {
				if d, ok := val.Interface().([][]string); ok {
					return d
				}
			}
		}
	}
	return nil
}

// get snaked column name
func getColumnName(ft int, addrField reflect.Value, sf reflect.StructField, col string) string {
	column := col
	if col == "" {
		column = snakeString(sf.Name)
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
		ft = TypeCharField
	default:
		switch elm.Interface().(type) {
		case sql.NullInt64:
			ft = TypeBigIntegerField
		case sql.NullFloat64:
			ft = TypeFloatField
		case sql.NullBool:
			ft = TypeBooleanField
		case sql.NullString:
			ft = TypeCharField
		case time.Time:
			ft = TypeDateTimeField
		}
	}
	if ft&IsFieldType == 0 {
		err = fmt.Errorf("unsupport field type %s, may be miss setting tag", val)
	}
	return
}

// parse struct tag string
func parseStructTag(data string, attrs *map[string]bool, tags *map[string]string) {
	attr := make(map[string]bool)
	tag := make(map[string]string)
	for _, v := range strings.Split(data, defaultStructTagDelim) {
		v = strings.TrimSpace(v)
		if supportTag[v] == 1 {
			attr[v] = true
		} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
			name := v[:i]
			if supportTag[name] == 2 {
				v = v[i+1 : len(v)-1]
				tag[name] = v
			}
		}
	}
	*attrs = attr
	*tags = tag
}
