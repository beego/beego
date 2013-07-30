package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type fieldChoices []StrTo

func (f *fieldChoices) Add(s StrTo) {
	if f.Have(s) == false {
		*f = append(*f, s)
	}
}

func (f *fieldChoices) Clear() {
	*f = fieldChoices([]StrTo{})
}

func (f *fieldChoices) Have(s StrTo) bool {
	for _, v := range *f {
		if v == s {
			return true
		}
	}
	return false
}

func (f *fieldChoices) Clone() fieldChoices {
	return *f
}

type primaryKeys []*fieldInfo

func (p *primaryKeys) Add(fi *fieldInfo) {
	*p = append(*p, fi)
}

func (p primaryKeys) Exist(fi *fieldInfo) (int, bool) {
	for i, v := range p {
		if v == fi {
			return i, true
		}
	}
	return -1, false
}

func (p primaryKeys) IsMulti() bool {
	return len(p) > 1
}

func (p primaryKeys) IsEmpty() bool {
	return len(p) == 0
}

type fields struct {
	pk            primaryKeys
	auto          *fieldInfo
	columns       map[string]*fieldInfo
	fields        map[string]*fieldInfo
	fieldsLow     map[string]*fieldInfo
	fieldsByType  map[int][]*fieldInfo
	fieldsRel     []*fieldInfo
	fieldsReverse []*fieldInfo
	fieldsDB      []*fieldInfo
	rels          []*fieldInfo
	orders        []string
	dbcols        []string
}

func (f *fields) Add(fi *fieldInfo) (added bool) {
	if f.fields[fi.name] == nil && f.columns[fi.column] == nil {
		f.columns[fi.column] = fi
		f.fields[fi.name] = fi
		f.fieldsLow[strings.ToLower(fi.name)] = fi
	} else {
		return
	}
	if _, ok := f.fieldsByType[fi.fieldType]; ok == false {
		f.fieldsByType[fi.fieldType] = make([]*fieldInfo, 0)
	}
	f.fieldsByType[fi.fieldType] = append(f.fieldsByType[fi.fieldType], fi)
	f.orders = append(f.orders, fi.column)
	if fi.dbcol {
		f.dbcols = append(f.dbcols, fi.column)
		f.fieldsDB = append(f.fieldsDB, fi)
	}
	if fi.rel {
		f.fieldsRel = append(f.fieldsRel, fi)
	}
	if fi.reverse {
		f.fieldsReverse = append(f.fieldsReverse, fi)
	}
	return true
}

func (f *fields) GetByName(name string) *fieldInfo {
	return f.fields[name]
}

func (f *fields) GetByColumn(column string) *fieldInfo {
	return f.columns[column]
}

func (f *fields) GetByAny(name string) (*fieldInfo, bool) {
	if fi, ok := f.fields[name]; ok {
		return fi, ok
	}
	if fi, ok := f.fieldsLow[strings.ToLower(name)]; ok {
		return fi, ok
	}
	if fi, ok := f.columns[name]; ok {
		return fi, ok
	}
	return nil, false
}

func newFields() *fields {
	f := new(fields)
	f.fields = make(map[string]*fieldInfo)
	f.fieldsLow = make(map[string]*fieldInfo)
	f.columns = make(map[string]*fieldInfo)
	f.fieldsByType = make(map[int][]*fieldInfo)
	return f
}

type fieldInfo struct {
	mi                  *modelInfo
	fieldIndex          int
	fieldType           int
	dbcol               bool
	inModel             bool
	name                string
	fullName            string
	column              string
	addrValue           *reflect.Value
	sf                  *reflect.StructField
	auto                bool
	pk                  bool
	null                bool
	blank               bool
	index               bool
	unique              bool
	initial             StrTo
	choices             fieldChoices
	maxLength           int
	auto_now            bool
	auto_now_add        bool
	rel                 bool
	reverse             bool
	reverseField        string
	reverseFieldInfo    *fieldInfo
	relTable            string
	relThrough          string
	relThroughModelInfo *modelInfo
	relModelInfo        *modelInfo
	digits              int
	decimals            int
	isFielder           bool
	onDelete            string
}

func newFieldInfo(mi *modelInfo, field reflect.Value, sf reflect.StructField) (fi *fieldInfo, err error) {
	var (
		tag       string
		tagValue  string
		choices   fieldChoices
		values    fieldChoices
		initial   StrTo
		fieldType int
		attrs     map[string]bool
		tags      map[string]string
		parts     []string
		addrField reflect.Value
	)

	fi = new(fieldInfo)

	if field.Kind() != reflect.Ptr && field.Kind() != reflect.Slice && field.CanAddr() {
		addrField = field.Addr()
	} else {
		addrField = field
	}

	parseStructTag(sf.Tag.Get(defaultStructTagName), &attrs, &tags)

	digits := tags["digits"]
	decimals := tags["decimals"]
	maxLength := tags["max_length"]
	onDelete := tags["on_delete"]

checkType:
	switch f := addrField.Interface().(type) {
	case Fielder:
		fi.isFielder = true
		if field.Kind() == reflect.Ptr {
			err = fmt.Errorf("the model Fielder can not be use ptr")
			goto end
		}
		fieldType = f.FieldType()
		if fieldType&IsRelField > 0 {
			err = fmt.Errorf("unsupport rel type custom field")
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
					fi.relTable = tv
				} else if tv := tags["rel_through"]; tv != "" {
					fi.relThrough = tv
				}
				break checkType
			default:
				err = fmt.Errorf("error")
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
				break checkType
			default:
				err = fmt.Errorf("error")
				goto wrongTag
			}
		}

		fieldType, err = getFieldType(addrField)
		if err != nil {
			goto end
		}
		if fieldType == TypeTextField && maxLength != "" {
			fieldType = TypeCharField
		}
		if fieldType == TypeFloatField && (digits != "" || decimals != "") {
			fieldType = TypeDecimalField
		}
		if fieldType == TypeDateTimeField && attrs["date"] {
			fieldType = TypeDateField
		}
	}

	switch fieldType {
	case RelForeignKey, RelOneToOne, RelReverseOne:
		if _, ok := addrField.Interface().(Modeler); ok == false {
			err = fmt.Errorf("rel/reverse:one field must be implements Modeler")
			goto end
		}
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
			if _, ok := reflect.New(field.Type().Elem()).Elem().Interface().(Modeler); ok == false {
				err = fmt.Errorf("rel/reverse:many slice element must be implements Modeler")
				goto end
			}
		}
	}

	if fieldType&IsFieldType == 0 {
		err = fmt.Errorf("wrong field type")
		goto end
	}

	fi.fieldType = fieldType
	fi.name = sf.Name
	fi.column = getColumnName(fieldType, addrField, sf, tags["column"])
	fi.addrValue = &addrField
	fi.sf = &sf
	fi.fullName = mi.fullName + "." + sf.Name

	fi.null = attrs["null"]
	fi.blank = attrs["blank"]
	fi.index = attrs["index"]
	fi.auto = attrs["auto"]
	fi.pk = attrs["pk"]
	fi.unique = attrs["unique"]

	switch fieldType {
	case RelManyToMany, RelReverseMany, RelReverseOne:
		fi.null = false
		fi.blank = false
		fi.index = false
		fi.auto = false
		fi.pk = false
		fi.unique = false
	default:
		fi.dbcol = true
	}

	switch fieldType {
	case RelForeignKey, RelOneToOne, RelManyToMany:
		fi.rel = true
		if fieldType == RelOneToOne {
			fi.unique = true
		}
	case RelReverseMany, RelReverseOne:
		fi.reverse = true
	}

	if fi.rel && fi.dbcol {
		switch onDelete {
		case od_CASCADE, od_DO_NOTHING:
		case od_SET_DEFAULT:
			if tags["default"] == "" {
				err = errors.New("on_delete: set_default need set field a default value")
				goto end
			}
		case od_SET_NULL:
			if fi.null == false {
				err = errors.New("on_delete: set_null need set field null")
				goto end
			}
		default:
			if onDelete == "" {
				onDelete = od_CASCADE
			} else {
				err = fmt.Errorf("on_delete value expected choice in `cascade,set_null,set_default,do_nothing`, unknown `%s`", onDelete)
				goto end
			}
		}

		fi.onDelete = onDelete
	}

	switch fieldType {
	case TypeBooleanField:
	case TypeCharField:
		if maxLength != "" {
			v, e := StrTo(maxLength).Int32()
			if e != nil {
				err = fmt.Errorf("wrong maxLength value `%s`", maxLength)
			} else {
				fi.maxLength = int(v)
			}
		} else {
			err = fmt.Errorf("maxLength must be specify")
		}
	case TypeTextField:
		fi.index = false
		fi.unique = false
	case TypeDateField, TypeDateTimeField:
		if attrs["auto_now"] {
			fi.auto_now = true
		} else if attrs["auto_now_add"] {
			fi.auto_now_add = true
		}
	case TypeFloatField:
	case TypeDecimalField:
		d1 := digits
		d2 := decimals
		v1, er1 := StrTo(d1).Int16()
		v2, er2 := StrTo(d2).Int16()
		if er1 != nil || er2 != nil {
			err = fmt.Errorf("wrong digits/decimals value %s/%s", d2, d1)
			goto end
		}
		fi.digits = int(v1)
		fi.decimals = int(v2)
	default:
		switch {
		case fieldType&IsIntegerField > 0:
		case fieldType&IsRelField > 0:
		}
	}

	if fieldType&IsIntegerField == 0 {
		if fi.auto {
			err = fmt.Errorf("non-integer type cannot set auto")
			goto end
		}

		if fi.pk || fi.index || fi.unique {
			if fieldType != TypeCharField && fieldType != RelOneToOne {
				err = fmt.Errorf("cannot set pk/index/unique")
				goto end
			}
		}
	}

	if fi.auto || fi.pk {
		if fi.auto {
			fi.pk = true
		}
		fi.null = false
		fi.blank = false
		fi.index = false
		fi.unique = false
	}

	if fi.unique {
		fi.null = false
		fi.blank = false
		fi.index = false
	}

	parts = strings.Split(tags["choices"], ",")
	if len(parts) > 1 {
		for _, v := range parts {
			choices.Add(StrTo(strings.TrimSpace(v)))
		}
	}

	initial.Clear()
	if v, ok := tags["default"]; ok {
		initial.Set(v)
	}

	if fi.auto || fi.pk || fi.unique || fieldType == TypeDateField || fieldType == TypeDateTimeField {
		// can not set default
		choices.Clear()
		initial.Clear()
	}

	values = choices.Clone()

	if initial.Exist() {
		values.Add(initial)
	}

	for i, v := range values {
		switch fieldType {
		case TypeBooleanField:
			_, err = v.Bool()
		case TypeFloatField, TypeDecimalField:
			_, err = v.Float64()
		case TypeSmallIntegerField:
			_, err = v.Int16()
		case TypeIntegerField:
			_, err = v.Int32()
		case TypeBigIntegerField:
			_, err = v.Int64()
		case TypePositiveSmallIntegerField:
			_, err = v.Uint16()
		case TypePositiveIntegerField:
			_, err = v.Uint32()
		case TypePositiveBigIntegerField:
			_, err = v.Uint64()
		}
		if err != nil {
			if initial.Exist() && len(values) == i {
				tag, tagValue = "default", tags["default"]
			} else {
				tag, tagValue = "choices", tags["choices"]
			}
			goto wrongTag
		}
	}

	if len(choices) > 0 && initial.Exist() {
		if choices.Have(initial) == false {
			err = fmt.Errorf("default value `%s` not in choices `%s`", tags["default"], tags["choices"])
			goto end
		}
	}

	fi.choices = choices
	fi.initial = initial

end:
	if err != nil {
		return nil, err
	}
	return
wrongTag:
	return nil, fmt.Errorf("wrong tag format: `%s:\"%s\"`, %s", tag, tagValue, err)
}
