package orm

import (
	"fmt"
	"reflect"
	"time"
)

func getDbAlias(name string) *alias {
	if al, ok := dataBaseCache.get(name); ok {
		return al
	} else {
		panic(fmt.Errorf("unknown DataBase alias name %s", name))
	}
	return nil
}

func getExistPk(mi *modelInfo, ind reflect.Value) (column string, value interface{}, exist bool) {
	fi := mi.fields.pk

	v := ind.Field(fi.fieldIndex)
	if fi.fieldType&IsPostiveIntegerField > 0 {
		vu := v.Uint()
		exist = vu > 0
		value = vu
	} else if fi.fieldType&IsIntegerField > 0 {
		vu := v.Int()
		exist = vu > 0
		value = vu
	} else {
		vu := v.String()
		exist = vu != ""
		value = vu
	}

	column = fi.column
	return
}

func getFlatParams(fi *fieldInfo, args []interface{}, tz *time.Location) (params []interface{}) {

outFor:
	for _, arg := range args {
		val := reflect.ValueOf(arg)

		if arg == nil {
			params = append(params, arg)
			continue
		}

		switch v := arg.(type) {
		case []byte:
		case string:
			if fi != nil {
				if fi.fieldType == TypeDateField || fi.fieldType == TypeDateTimeField {
					var t time.Time
					var err error
					if len(v) >= 19 {
						s := v[:19]
						t, err = time.ParseInLocation(format_DateTime, s, DefaultTimeLoc)
					} else {
						s := v
						if len(v) > 10 {
							s = v[:10]
						}
						t, err = time.ParseInLocation(format_Date, s, tz)
					}
					if err == nil {
						if fi.fieldType == TypeDateField {
							v = t.In(tz).Format(format_Date)
						} else {
							v = t.In(tz).Format(format_DateTime)
						}
					}
				}
			}
			arg = v
		case time.Time:
			if fi != nil && fi.fieldType == TypeDateField {
				arg = v.In(tz).Format(format_Date)
			} else {
				arg = v.In(tz).Format(format_DateTime)
			}
		default:
			kind := val.Kind()
			switch kind {
			case reflect.Slice, reflect.Array:

				var args []interface{}
				for i := 0; i < val.Len(); i++ {
					v := val.Index(i)

					var vu interface{}
					if v.CanInterface() {
						vu = v.Interface()
					}

					if vu == nil {
						continue
					}

					args = append(args, vu)
				}

				if len(args) > 0 {
					p := getFlatParams(fi, args, tz)
					params = append(params, p...)
				}
				continue outFor

			case reflect.Ptr, reflect.Struct:
				ind := reflect.Indirect(val)

				if ind.Kind() == reflect.Struct {
					typ := ind.Type()
					name := getFullName(typ)
					var value interface{}
					if mmi, ok := modelCache.getByFN(name); ok {
						if _, vu, exist := getExistPk(mmi, ind); exist {
							value = vu
						}
					}
					arg = value

					if arg == nil {
						panic(fmt.Errorf("need a valid args value, unknown table or value `%s`", name))
					}
				} else {
					arg = ind.Interface()
				}
			}
		}
		params = append(params, arg)
	}
	return
}
