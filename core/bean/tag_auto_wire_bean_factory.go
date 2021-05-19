// Copyright 2020
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

package bean

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"

	"github.com/beego/beego/v2/core/logs"
)

const DefaultValueTagKey = "default"

// TagAutoWireBeanFactory wire the bean based on Fields' tag
// if field's value is "zero value", we will execute injection
// see reflect.Value.IsZero()
// If field's kind is one of(reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice
// reflect.UnsafePointer, reflect.Array, reflect.Uintptr, reflect.Complex64, reflect.Complex128
// reflect.Ptr, reflect.Struct),
// it will be ignored
type TagAutoWireBeanFactory struct {
	// we allow user register their TypeAdapter
	Adapters map[string]TypeAdapter

	// FieldTagParser is an extension point which means that you can custom how to read field's metadata from tag
	FieldTagParser func(field reflect.StructField) *FieldMetadata
}

// NewTagAutoWireBeanFactory create an instance of TagAutoWireBeanFactory
// by default, we register Time adapter, the time will be parse by using layout "2006-01-02 15:04:05"
// If you need more adapter, you can implement interface TypeAdapter
func NewTagAutoWireBeanFactory() *TagAutoWireBeanFactory {
	return &TagAutoWireBeanFactory{
		Adapters: map[string]TypeAdapter{
			"Time": &TimeTypeAdapter{Layout: "2006-01-02 15:04:05"},
		},

		FieldTagParser: func(field reflect.StructField) *FieldMetadata {
			return &FieldMetadata{
				DftValue: field.Tag.Get(DefaultValueTagKey),
			}
		},
	}
}

// AutoWire use value from appCtx to wire the bean, or use default value, or do nothing
func (t *TagAutoWireBeanFactory) AutoWire(ctx context.Context, appCtx ApplicationContext, bean interface{}) error {
	if bean == nil {
		return nil
	}

	v := reflect.Indirect(reflect.ValueOf(bean))

	bm := t.getConfig(v)

	// field name, field metadata
	for fn, fm := range bm.Fields {
		fValue := v.FieldByName(fn)
		if len(fm.DftValue) == 0 || !t.needInject(fValue) || !fValue.CanSet() {
			continue
		}

		// handle type adapter
		typeName := fValue.Type().Name()
		if adapter, ok := t.Adapters[typeName]; ok {
			dftValue, err := adapter.DefaultValue(ctx, fm.DftValue)
			if err == nil {
				fValue.Set(reflect.ValueOf(dftValue))
				continue
			} else {
				return err
			}
		}

		switch fValue.Kind() {
		case reflect.Bool:
			if v, err := strconv.ParseBool(fm.DftValue); err != nil {
				return errors.WithMessage(err,
					fmt.Sprintf("can not convert the field[%s]'s default value[%s] to bool value",
						fn, fm.DftValue))
			} else {
				fValue.SetBool(v)
				continue
			}
		case reflect.Int:
			if err := t.setIntXValue(fm.DftValue, 0, fn, fValue); err != nil {
				return err
			}
			continue
		case reflect.Int8:
			if err := t.setIntXValue(fm.DftValue, 8, fn, fValue); err != nil {
				return err
			}
			continue
		case reflect.Int16:
			if err := t.setIntXValue(fm.DftValue, 16, fn, fValue); err != nil {
				return err
			}
			continue

		case reflect.Int32:
			if err := t.setIntXValue(fm.DftValue, 32, fn, fValue); err != nil {
				return err
			}
			continue

		case reflect.Int64:
			if err := t.setIntXValue(fm.DftValue, 64, fn, fValue); err != nil {
				return err
			}
			continue

		case reflect.Uint:
			if err := t.setUIntXValue(fm.DftValue, 0, fn, fValue); err != nil {
				return err
			}

		case reflect.Uint8:
			if err := t.setUIntXValue(fm.DftValue, 8, fn, fValue); err != nil {
				return err
			}
			continue

		case reflect.Uint16:
			if err := t.setUIntXValue(fm.DftValue, 16, fn, fValue); err != nil {
				return err
			}
			continue
		case reflect.Uint32:
			if err := t.setUIntXValue(fm.DftValue, 32, fn, fValue); err != nil {
				return err
			}
			continue

		case reflect.Uint64:
			if err := t.setUIntXValue(fm.DftValue, 64, fn, fValue); err != nil {
				return err
			}
			continue

		case reflect.Float32:
			if err := t.setFloatXValue(fm.DftValue, 32, fn, fValue); err != nil {
				return err
			}
			continue
		case reflect.Float64:
			if err := t.setFloatXValue(fm.DftValue, 64, fn, fValue); err != nil {
				return err
			}
			continue

		case reflect.String:
			fValue.SetString(fm.DftValue)
			continue

		// case reflect.Ptr:
		// case reflect.Struct:
		default:
			logs.Warn("this field[%s] has default setting, but we don't support this type: %s",
				fn, fValue.Kind().String())
		}
	}
	return nil
}

func (t *TagAutoWireBeanFactory) setFloatXValue(dftValue string, bitSize int, fn string, fv reflect.Value) error {
	if v, err := strconv.ParseFloat(dftValue, bitSize); err != nil {
		return errors.WithMessage(err,
			fmt.Sprintf("can not convert the field[%s]'s default value[%s] to float%d value",
				fn, dftValue, bitSize))
	} else {
		fv.SetFloat(v)
		return nil
	}
}

func (t *TagAutoWireBeanFactory) setUIntXValue(dftValue string, bitSize int, fn string, fv reflect.Value) error {
	if v, err := strconv.ParseUint(dftValue, 10, bitSize); err != nil {
		return errors.WithMessage(err,
			fmt.Sprintf("can not convert the field[%s]'s default value[%s] to uint%d value",
				fn, dftValue, bitSize))
	} else {
		fv.SetUint(v)
		return nil
	}
}

func (t *TagAutoWireBeanFactory) setIntXValue(dftValue string, bitSize int, fn string, fv reflect.Value) error {
	if v, err := strconv.ParseInt(dftValue, 10, bitSize); err != nil {
		return errors.WithMessage(err,
			fmt.Sprintf("can not convert the field[%s]'s default value[%s] to int%d value",
				fn, dftValue, bitSize))
	} else {
		fv.SetInt(v)
		return nil
	}
}

func (t *TagAutoWireBeanFactory) needInject(fValue reflect.Value) bool {
	return fValue.IsZero()
}

// getConfig never return nil
func (t *TagAutoWireBeanFactory) getConfig(beanValue reflect.Value) *BeanMetadata {
	fms := make(map[string]*FieldMetadata, beanValue.NumField())
	for i := 0; i < beanValue.NumField(); i++ {
		// f => StructField
		f := beanValue.Type().Field(i)
		fms[f.Name] = t.FieldTagParser(f)
	}
	return &BeanMetadata{
		Fields: fms,
	}
}
