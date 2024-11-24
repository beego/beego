// Copyright 2020 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package context

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	sliceOfInts    = reflect.TypeOf([]int(nil))
	sliceOfStrings = reflect.TypeOf([]string(nil))
)

// ParseForm will parse form values to struct via tag.
// Support for anonymous struct.
func parseFormToStruct(form url.Values, objT reflect.Type, objV reflect.Value) error {
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}

		fieldT := objT.Field(i)
		if fieldT.Anonymous && fieldT.Type.Kind() == reflect.Struct {
			err := parseFormToStruct(form, fieldT.Type, fieldV)
			if err != nil {
				return err
			}
			continue
		}

		tag, ok := formTagName(fieldT)
		if !ok {
			continue
		}

		value, ok := formValue(tag, form, fieldT)
		if !ok {
			continue
		}

		switch fieldT.Type.Kind() {
		case reflect.Bool:
			b, err := parseFormBoolValue(value)
			if err != nil {
				return err
			}
			fieldV.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			fieldV.SetFloat(x)
		case reflect.Interface:
			fieldV.Set(reflect.ValueOf(value))
		case reflect.String:
			fieldV.SetString(value)
		case reflect.Struct:
			if fieldT.Type.String() == "time.Time" {
				t, err := parseFormTime(value)
				if err != nil {
					return err
				}
				fieldV.Set(reflect.ValueOf(t))
			}
		case reflect.Slice:
			if fieldT.Type == sliceOfInts {
				formVals := form[tag]
				fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(int(1))), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					val, err := strconv.Atoi(formVals[i])
					if err != nil {
						return err
					}
					fieldV.Index(i).SetInt(int64(val))
				}
			} else if fieldT.Type == sliceOfStrings {
				formVals := form[tag]
				fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf("")), len(formVals), len(formVals)))
				for i := 0; i < len(formVals); i++ {
					fieldV.Index(i).SetString(formVals[i])
				}
			}
		}
	}
	return nil
}

// nolint
func parseFormTime(value string) (time.Time, error) {
	var pattern string
	if len(value) >= 25 {
		value = value[:25]
		pattern = time.RFC3339
	} else if strings.HasSuffix(strings.ToUpper(value), "Z") {
		pattern = time.RFC3339
	} else if len(value) >= 19 {
		if strings.Contains(value, "T") {
			pattern = formatDateTimeT
		} else {
			pattern = formatDateTime
		}
		value = value[:19]
	} else if len(value) >= 10 {
		if len(value) > 10 {
			value = value[:10]
		}
		pattern = formatDate
	} else if len(value) >= 8 {
		if len(value) > 8 {
			value = value[:8]
		}
		pattern = formatTime
	}
	return time.ParseInLocation(pattern, value, time.Local)
}

func parseFormBoolValue(value string) (bool, error) {
	if strings.ToLower(value) == "on" || strings.ToLower(value) == "1" || strings.ToLower(value) == "yes" {
		return true, nil
	}
	if strings.ToLower(value) == "off" || strings.ToLower(value) == "0" || strings.ToLower(value) == "no" {
		return false, nil
	}
	return strconv.ParseBool(value)
}

// nolint
func formTagName(fieldT reflect.StructField) (string, bool) {
	tags := strings.Split(fieldT.Tag.Get("form"), ",")
	var tag string
	if len(tags) == 0 || tags[0] == "" {
		tag = fieldT.Name
	} else if tags[0] == "-" {
		return "", false
	} else {
		tag = tags[0]
	}
	return tag, true
}

func formValue(tag string, form url.Values, fieldT reflect.StructField) (string, bool) {
	formValues := form[tag]
	if len(formValues) == 0 {
		defaultValue := fieldT.Tag.Get("default")
		return defaultValue, defaultValue != ""
	}
	return formValues[0], true
}
