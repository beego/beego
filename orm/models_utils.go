package orm

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func getFullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}

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

func getColumnName(ft int, addrField reflect.Value, sf reflect.StructField, col string) string {
	col = strings.ToLower(col)
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
	case reflect.Invalid:
	default:
		if elm.CanInterface() {
			if _, ok := elm.Interface().(time.Time); ok {
				ft = TypeDateTimeField
			}
		}
	}
	if ft&IsFieldType == 0 {
		err = fmt.Errorf("unsupport field type %s, may be miss setting tag", val)
	}
	return
}

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
