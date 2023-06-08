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

package orm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/beego/beego/v2/client/orm/internal/utils"

	"github.com/beego/beego/v2/client/orm/internal/models"
)

// get table alias.
func getDbAlias(name string) *alias {
	if al, ok := dataBaseCache.get(name); ok {
		return al
	}
	panic(fmt.Errorf("unknown DataBase alias name %s", name))
}

// get pk column info.
func getExistPk(mi *models.ModelInfo, ind reflect.Value) (column string, value interface{}, exist bool) {
	fi := mi.Fields.Pk

	v := ind.FieldByIndex(fi.FieldIndex)
	if fi.FieldType&IsPositiveIntegerField > 0 {
		vu := v.Uint()
		exist = vu > 0
		value = vu
	} else if fi.FieldType&IsIntegerField > 0 {
		vu := v.Int()
		exist = true
		value = vu
	} else if fi.FieldType&IsRelField > 0 {
		_, value, exist = getExistPk(fi.RelModelInfo, reflect.Indirect(v))
	} else {
		vu := v.String()
		exist = vu != ""
		value = vu
	}

	column = fi.Column
	return
}

// get Fields description as flatted string.
func getFlatParams(fi *models.FieldInfo, args []interface{}, tz *time.Location) (params []interface{}) {
outFor:
	for _, arg := range args {
		if arg == nil {
			params = append(params, arg)
			continue
		}

		val := reflect.ValueOf(arg)
		kind := val.Kind()
		if kind == reflect.Ptr {
			val = val.Elem()
			kind = val.Kind()
			arg = val.Interface()
		}

		switch kind {
		case reflect.String:
			v := val.String()
			if fi != nil {
				if fi.FieldType == TypeTimeField || fi.FieldType == TypeDateField || fi.FieldType == TypeDateTimeField {
					var t time.Time
					var err error
					if len(v) >= 19 {
						s := v[:19]
						t, err = time.ParseInLocation(utils.FormatDateTime, s, DefaultTimeLoc)
					} else if len(v) >= 10 {
						s := v
						if len(v) > 10 {
							s = v[:10]
						}
						t, err = time.ParseInLocation(utils.FormatDate, s, tz)
					} else {
						s := v
						if len(s) > 8 {
							s = v[:8]
						}
						t, err = time.ParseInLocation(utils.FormatTime, s, tz)
					}
					if err == nil {
						if fi.FieldType == TypeDateField {
							v = t.In(tz).Format(utils.FormatDate)
						} else if fi.FieldType == TypeDateTimeField {
							v = t.In(tz).Format(utils.FormatDateTime)
						} else {
							v = t.In(tz).Format(utils.FormatTime)
						}
					}
				}
			}
			arg = v
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			arg = val.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			arg = val.Uint()
		case reflect.Float32:
			arg, _ = utils.StrTo(utils.ToStr(arg)).Float64()
		case reflect.Float64:
			arg = val.Float()
		case reflect.Bool:
			arg = val.Bool()
		case reflect.Slice, reflect.Array:
			if _, ok := arg.([]byte); ok {
				continue outFor
			}

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
		case reflect.Struct:
			if v, ok := arg.(time.Time); ok {
				if fi != nil && fi.FieldType == TypeDateField {
					arg = v.In(tz).Format(utils.FormatDate)
				} else if fi != nil && fi.FieldType == TypeDateTimeField {
					arg = v.In(tz).Format(utils.FormatDateTime)
				} else if fi != nil && fi.FieldType == TypeTimeField {
					arg = v.In(tz).Format(utils.FormatTime)
				} else {
					arg = v.In(tz).Format(utils.FormatDateTime)
				}
			} else {
				typ := val.Type()
				name := models.GetFullName(typ)
				var value interface{}
				if mmi, ok := defaultModelCache.getByFullName(name); ok {
					if _, vu, exist := getExistPk(mmi, val); exist {
						value = vu
					}
				}
				arg = value

				if arg == nil {
					panic(fmt.Errorf("need a valid args value, unknown table or value `%s`", name))
				}
			}
		}

		params = append(params, arg)
	}
	return
}
