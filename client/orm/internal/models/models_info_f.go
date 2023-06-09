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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm/internal/utils"
)

var errSkipField = errors.New("skip field")

// Fields field info collection
type Fields struct {
	Pk            *FieldInfo
	Columns       map[string]*FieldInfo
	Fields        map[string]*FieldInfo
	FieldsLow     map[string]*FieldInfo
	FieldsByType  map[int][]*FieldInfo
	FieldsRel     []*FieldInfo
	FieldsReverse []*FieldInfo
	FieldsDB      []*FieldInfo
	Rels          []*FieldInfo
	Orders        []string
	DBcols        []string
}

// Add adds field info
func (f *Fields) Add(fi *FieldInfo) (added bool) {
	if f.Fields[fi.Name] == nil && f.Columns[fi.Column] == nil {
		f.Columns[fi.Column] = fi
		f.Fields[fi.Name] = fi
		f.FieldsLow[strings.ToLower(fi.Name)] = fi
	} else {
		return
	}
	if _, ok := f.FieldsByType[fi.FieldType]; !ok {
		f.FieldsByType[fi.FieldType] = make([]*FieldInfo, 0)
	}
	f.FieldsByType[fi.FieldType] = append(f.FieldsByType[fi.FieldType], fi)
	f.Orders = append(f.Orders, fi.Column)
	if fi.DBcol {
		f.DBcols = append(f.DBcols, fi.Column)
		f.FieldsDB = append(f.FieldsDB, fi)
	}
	if fi.Rel {
		f.FieldsRel = append(f.FieldsRel, fi)
	}
	if fi.Reverse {
		f.FieldsReverse = append(f.FieldsReverse, fi)
	}
	return true
}

// GetByName get field info by name
func (f *Fields) GetByName(name string) *FieldInfo {
	return f.Fields[name]
}

// GetByColumn get field info by column name
func (f *Fields) GetByColumn(column string) *FieldInfo {
	return f.Columns[column]
}

// GetByAny get field info by string, name is prior
func (f *Fields) GetByAny(name string) (*FieldInfo, bool) {
	if fi, ok := f.Fields[name]; ok {
		return fi, ok
	}
	if fi, ok := f.FieldsLow[strings.ToLower(name)]; ok {
		return fi, ok
	}
	if fi, ok := f.Columns[name]; ok {
		return fi, ok
	}
	return nil, false
}

// NewFields create new field info collection
func NewFields() *Fields {
	f := new(Fields)
	f.Fields = make(map[string]*FieldInfo)
	f.FieldsLow = make(map[string]*FieldInfo)
	f.Columns = make(map[string]*FieldInfo)
	f.FieldsByType = make(map[int][]*FieldInfo)
	return f
}

// FieldInfo single field info
type FieldInfo struct {
	DBcol               bool // table column fk and onetoone
	InModel             bool
	Auto                bool
	Pk                  bool
	Null                bool
	Index               bool
	Unique              bool
	ColDefault          bool // whether has default tag
	ToText              bool
	AutoNow             bool
	AutoNowAdd          bool
	Rel                 bool // if type equal to RelForeignKey, RelOneToOne, RelManyToMany then true
	Reverse             bool
	IsFielder           bool // implement Fielder interface
	Mi                  *ModelInfo
	FieldIndex          []int
	FieldType           int
	Name                string
	FullName            string
	Column              string
	AddrValue           reflect.Value
	Sf                  reflect.StructField
	Initial             utils.StrTo // store the default value
	Size                int
	ReverseField        string
	ReverseFieldInfo    *FieldInfo
	ReverseFieldInfoTwo *FieldInfo
	ReverseFieldInfoM2M *FieldInfo
	RelTable            string
	RelThrough          string
	RelThroughModelInfo *ModelInfo
	RelModelInfo        *ModelInfo
	Digits              int
	Decimals            int
	OnDelete            string
	Description         string
	TimePrecision       *int
}

// NewFieldInfo new field info
func NewFieldInfo(mi *ModelInfo, field reflect.Value, sf reflect.StructField, mName string) (fi *FieldInfo, err error) {
	var (
		tag       string
		tagValue  string
		initial   utils.StrTo // store the default value
		fieldType int
		attrs     map[string]bool
		tags      map[string]string
		addrField reflect.Value
	)

	fi = new(FieldInfo)

	// if field which CanAddr is the follow type
	//  A value is addressable if it is an element of a slice,
	//  an element of an addressable array, a field of an
	//  addressable struct, or the result of dereferencing a pointer.
	addrField = field
	if field.CanAddr() && field.Kind() != reflect.Ptr {
		addrField = field.Addr()
		if _, ok := addrField.Interface().(Fielder); !ok {
			if field.Kind() == reflect.Slice {
				addrField = field
			}
		}
	}

	attrs, tags = ParseStructTag(sf.Tag.Get(DefaultStructTagName))

	if _, ok := attrs["-"]; ok {
		return nil, errSkipField
	}

	digits := tags["digits"]
	decimals := tags["decimals"]
	size := tags["size"]
	onDelete := tags["on_delete"]
	precision := tags["precision"]
	initial.Clear()
	if v, ok := tags["default"]; ok {
		initial.Set(v)
	}

checkType:
	switch f := addrField.Interface().(type) {
	case Fielder:
		fi.IsFielder = true
		if field.Kind() == reflect.Ptr {
			err = fmt.Errorf("the model Fielder can not be use ptr")
			goto end
		}
		fieldType = f.FieldType()
		if fieldType&IsRelField > 0 {
			err = fmt.Errorf("unsupport type custom field, please refer to https://github.com/beego/beego/v2/blob/master/orm/models_fields.go#L24-L42")
			goto end
		}
	default:
		tag = "rel"
		tagValue = tags[tag]
		if tagValue != "" {
			switch tagValue {
			case "fk":
				fieldType = RelForeignKey
				break checkType
			case "one":
				fieldType = RelOneToOne
				break checkType
			case "m2m":
				fieldType = RelManyToMany
				if tv := tags["rel_table"]; tv != "" {
					fi.RelTable = tv
				} else if tv := tags["rel_through"]; tv != "" {
					fi.RelThrough = tv
				}
				break checkType
			default:
				err = fmt.Errorf("rel only allow these value: fk, one, m2m")
				goto wrongTag
			}
		}
		tag = "reverse"
		tagValue = tags[tag]
		if tagValue != "" {
			switch tagValue {
			case "one":
				fieldType = RelReverseOne
				break checkType
			case "many":
				fieldType = RelReverseMany
				if tv := tags["rel_table"]; tv != "" {
					fi.RelTable = tv
				} else if tv := tags["rel_through"]; tv != "" {
					fi.RelThrough = tv
				}
				break checkType
			default:
				err = fmt.Errorf("reverse only allow these value: one, many")
				goto wrongTag
			}
		}

		fieldType, err = getFieldType(addrField)
		if err != nil {
			goto end
		}
		if fieldType == TypeVarCharField {
			switch tags["type"] {
			case "char":
				fieldType = TypeCharField
			case "text":
				fieldType = TypeTextField
			case "json":
				fieldType = TypeJSONField
			case "jsonb":
				fieldType = TypeJsonbField
			}
		}
		if fieldType == TypeFloatField && (digits != "" || decimals != "") {
			fieldType = TypeDecimalField
		}
		if fieldType == TypeDateTimeField && tags["type"] == "date" {
			fieldType = TypeDateField
		}
		if fieldType == TypeTimeField && tags["type"] == "time" {
			fieldType = TypeTimeField
		}
	}

	// check the rel and reverse type
	// rel should Ptr
	// reverse should slice []*struct
	switch fieldType {
	case RelForeignKey, RelOneToOne, RelReverseOne:
		if field.Kind() != reflect.Ptr {
			err = fmt.Errorf("rel/reverse:one field must be *%s", field.Type().Name())
			goto end
		}
	case RelManyToMany, RelReverseMany:
		if field.Kind() != reflect.Slice {
			err = fmt.Errorf("rel/reverse:many field must be slice")
			goto end
		} else {
			if field.Type().Elem().Kind() != reflect.Ptr {
				err = fmt.Errorf("rel/reverse:many slice must be []*%s", field.Type().Elem().Name())
				goto end
			}
		}
	}

	if fieldType&IsFieldType == 0 {
		err = fmt.Errorf("wrong field type")
		goto end
	}

	fi.FieldType = fieldType
	fi.Name = sf.Name
	fi.Column = getColumnName(fieldType, addrField, sf, tags["column"])
	fi.AddrValue = addrField
	fi.Sf = sf
	fi.FullName = mi.FullName + mName + "." + sf.Name

	fi.Description = tags["description"]
	fi.Null = attrs["null"]
	fi.Index = attrs["index"]
	fi.Auto = attrs["auto"]
	fi.Pk = attrs["pk"]
	fi.Unique = attrs["unique"]

	// Mark object property if there is attribute "default" in the orm configuration
	if _, ok := tags["default"]; ok {
		fi.ColDefault = true
	}

	switch fieldType {
	case RelManyToMany, RelReverseMany, RelReverseOne:
		fi.Null = false
		fi.Index = false
		fi.Auto = false
		fi.Pk = false
		fi.Unique = false
	default:
		fi.DBcol = true
	}

	switch fieldType {
	case RelForeignKey, RelOneToOne, RelManyToMany:
		fi.Rel = true
		if fieldType == RelOneToOne {
			fi.Unique = true
		}
	case RelReverseMany, RelReverseOne:
		fi.Reverse = true
	}

	if fi.Rel && fi.DBcol {
		switch onDelete {
		case OdCascade, OdDoNothing:
		case OdSetDefault:
			if !initial.Exist() {
				err = errors.New("on_delete: set_default need set field a default value")
				goto end
			}
		case OdSetNULL:
			if !fi.Null {
				err = errors.New("on_delete: set_null need set field null")
				goto end
			}
		default:
			if onDelete == "" {
				onDelete = OdCascade
			} else {
				err = fmt.Errorf("on_delete value expected choice in `cascade,set_null,set_default,do_nothing`, unknown `%s`", onDelete)
				goto end
			}
		}

		fi.OnDelete = onDelete
	}

	switch fieldType {
	case TypeBooleanField:
	case TypeVarCharField, TypeCharField, TypeJSONField, TypeJsonbField:
		if size != "" {
			v, e := utils.StrTo(size).Int32()
			if e != nil {
				err = fmt.Errorf("wrong size value `%s`", size)
			} else {
				fi.Size = int(v)
			}
		} else {
			fi.Size = 255
			fi.ToText = true
		}
	case TypeTextField:
		fi.Index = false
		fi.Unique = false
	case TypeTimeField, TypeDateField, TypeDateTimeField:
		if fieldType == TypeDateTimeField {
			if precision != "" {
				v, e := utils.StrTo(precision).Int()
				if e != nil {
					err = fmt.Errorf("convert %s to int error:%v", precision, e)
				} else {
					fi.TimePrecision = &v
				}
			}
		}

		if attrs["auto_now"] {
			fi.AutoNow = true
		} else if attrs["auto_now_add"] {
			fi.AutoNowAdd = true
		}
	case TypeFloatField:
	case TypeDecimalField:
		d1 := digits
		d2 := decimals
		v1, er1 := utils.StrTo(d1).Int8()
		v2, er2 := utils.StrTo(d2).Int8()
		if er1 != nil || er2 != nil {
			err = fmt.Errorf("wrong digits/decimals value %s/%s", d2, d1)
			goto end
		}
		fi.Digits = int(v1)
		fi.Decimals = int(v2)
	default:
		switch {
		case fieldType&IsIntegerField > 0:
		case fieldType&IsRelField > 0:
		}
	}

	if fieldType&IsIntegerField == 0 {
		if fi.Auto {
			err = fmt.Errorf("non-integer type cannot set auto")
			goto end
		}
	}

	if fi.Auto || fi.Pk {
		if fi.Auto {
			switch addrField.Elem().Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
			default:
				err = fmt.Errorf("auto primary key only support int, int32, int64, uint, uint32, uint64 but found `%s`", addrField.Elem().Kind())
				goto end
			}
			fi.Pk = true
		}
		fi.Null = false
		fi.Index = false
		fi.Unique = false
	}

	if fi.Unique {
		fi.Index = false
	}

	// can not set default for these type
	if fi.Auto || fi.Pk || fi.Unique || fieldType == TypeTimeField || fieldType == TypeDateField || fieldType == TypeDateTimeField {
		initial.Clear()
	}

	if initial.Exist() {
		v := initial
		switch fieldType {
		case TypeBooleanField:
			_, err = v.Bool()
		case TypeFloatField, TypeDecimalField:
			_, err = v.Float64()
		case TypeBitField:
			_, err = v.Int8()
		case TypeSmallIntegerField:
			_, err = v.Int16()
		case TypeIntegerField:
			_, err = v.Int32()
		case TypeBigIntegerField:
			_, err = v.Int64()
		case TypePositiveBitField:
			_, err = v.Uint8()
		case TypePositiveSmallIntegerField:
			_, err = v.Uint16()
		case TypePositiveIntegerField:
			_, err = v.Uint32()
		case TypePositiveBigIntegerField:
			_, err = v.Uint64()
		}
		if err != nil {
			tag, tagValue = "default", tags["default"]
			goto wrongTag
		}
	}

	fi.Initial = initial
end:
	if err != nil {
		return nil, err
	}
	return
wrongTag:
	return nil, fmt.Errorf("wrong tag format: `%s:\"%s\"`, %s", tag, tagValue, err)
}
