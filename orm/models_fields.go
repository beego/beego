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
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	// bool
	TypeBooleanField = 1 << iota

	// string
	TypeCharField

	// string
	TypeTextField

	// time.Time
	TypeDateField
	// time.Time
	TypeDateTimeField

	// int8
	TypeBitField
	// int16
	TypeSmallIntegerField
	// int32
	TypeIntegerField
	// int64
	TypeBigIntegerField
	// uint8
	TypePositiveBitField
	// uint16
	TypePositiveSmallIntegerField
	// uint32
	TypePositiveIntegerField
	// uint64
	TypePositiveBigIntegerField

	// float64
	TypeFloatField
	// float64
	TypeDecimalField

	RelForeignKey
	RelOneToOne
	RelManyToMany
	RelReverseOne
	RelReverseMany
)

const (
	IsIntegerField        = ^-TypePositiveBigIntegerField >> 4 << 5
	IsPostiveIntegerField = ^-TypePositiveBigIntegerField >> 8 << 9
	IsRelField            = ^-RelReverseMany >> 14 << 15
	IsFieldType           = ^-RelReverseMany<<1 + 1
)

// A true/false field.
type BooleanField bool

func (e BooleanField) Value() bool {
	return bool(e)
}

func (e *BooleanField) Set(d bool) {
	*e = BooleanField(d)
}

func (e *BooleanField) String() string {
	return strconv.FormatBool(e.Value())
}

func (e *BooleanField) FieldType() int {
	return TypeBooleanField
}

func (e *BooleanField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case bool:
		e.Set(d)
	case string:
		v, err := StrTo(d).Bool()
		if err != nil {
			e.Set(v)
		}
		return err
	default:
		return errors.New(fmt.Sprintf("<BooleanField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *BooleanField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(BooleanField)

// A string field
// required values tag: size
// The size is enforced at the database level and in models’s validation.
// eg: `orm:"size(120)"`
type CharField string

func (e CharField) Value() string {
	return string(e)
}

func (e *CharField) Set(d string) {
	*e = CharField(d)
}

func (e *CharField) String() string {
	return e.Value()
}

func (e *CharField) FieldType() int {
	return TypeCharField
}

func (e *CharField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		e.Set(d)
	default:
		return errors.New(fmt.Sprintf("<CharField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *CharField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(CharField)

// A date, represented in go by a time.Time instance.
// only date values like 2006-01-02
// Has a few extra, optional attr tag:
//
// auto_now:
// Automatically set the field to now every time the object is saved. Useful for “last-modified” timestamps.
// Note that the current date is always used; it’s not just a default value that you can override.
//
// auto_now_add:
// Automatically set the field to now when the object is first created. Useful for creation of timestamps.
// Note that the current date is always used; it’s not just a default value that you can override.
//
// eg: `orm:"auto_now"` or `orm:"auto_now_add"`
type DateField time.Time

func (e DateField) Value() time.Time {
	return time.Time(e)
}

func (e *DateField) Set(d time.Time) {
	*e = DateField(d)
}

func (e *DateField) String() string {
	return e.Value().String()
}

func (e *DateField) FieldType() int {
	return TypeDateField
}

func (e *DateField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case time.Time:
		e.Set(d)
	case string:
		v, err := timeParse(d, format_Date)
		if err != nil {
			e.Set(v)
		}
		return err
	default:
		return errors.New(fmt.Sprintf("<DateField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *DateField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(DateField)

// A date, represented in go by a time.Time instance.
// datetime values like 2006-01-02 15:04:05
// Takes the same extra arguments as DateField.
type DateTimeField time.Time

func (e DateTimeField) Value() time.Time {
	return time.Time(e)
}

func (e *DateTimeField) Set(d time.Time) {
	*e = DateTimeField(d)
}

func (e *DateTimeField) String() string {
	return e.Value().String()
}

func (e *DateTimeField) FieldType() int {
	return TypeDateTimeField
}

func (e *DateTimeField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case time.Time:
		e.Set(d)
	case string:
		v, err := timeParse(d, format_DateTime)
		if err != nil {
			e.Set(v)
		}
		return err
	default:
		return errors.New(fmt.Sprintf("<DateTimeField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *DateTimeField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(DateTimeField)

// A floating-point number represented in go by a float32 value.
type FloatField float64

func (e FloatField) Value() float64 {
	return float64(e)
}

func (e *FloatField) Set(d float64) {
	*e = FloatField(d)
}

func (e *FloatField) String() string {
	return ToStr(e.Value(), -1, 32)
}

func (e *FloatField) FieldType() int {
	return TypeFloatField
}

func (e *FloatField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case float32:
		e.Set(float64(d))
	case float64:
		e.Set(d)
	case string:
		v, err := StrTo(d).Float64()
		if err != nil {
			e.Set(v)
		}
	default:
		return errors.New(fmt.Sprintf("<FloatField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *FloatField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(FloatField)

// -32768 to 32767
type SmallIntegerField int16

func (e SmallIntegerField) Value() int16 {
	return int16(e)
}

func (e *SmallIntegerField) Set(d int16) {
	*e = SmallIntegerField(d)
}

func (e *SmallIntegerField) String() string {
	return ToStr(e.Value())
}

func (e *SmallIntegerField) FieldType() int {
	return TypeSmallIntegerField
}

func (e *SmallIntegerField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case int16:
		e.Set(d)
	case string:
		v, err := StrTo(d).Int16()
		if err != nil {
			e.Set(v)
		}
	default:
		return errors.New(fmt.Sprintf("<SmallIntegerField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *SmallIntegerField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(SmallIntegerField)

// -2147483648 to 2147483647
type IntegerField int32

func (e IntegerField) Value() int32 {
	return int32(e)
}

func (e *IntegerField) Set(d int32) {
	*e = IntegerField(d)
}

func (e *IntegerField) String() string {
	return ToStr(e.Value())
}

func (e *IntegerField) FieldType() int {
	return TypeIntegerField
}

func (e *IntegerField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case int32:
		e.Set(d)
	case string:
		v, err := StrTo(d).Int32()
		if err != nil {
			e.Set(v)
		}
	default:
		return errors.New(fmt.Sprintf("<IntegerField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *IntegerField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(IntegerField)

// -9223372036854775808 to 9223372036854775807.
type BigIntegerField int64

func (e BigIntegerField) Value() int64 {
	return int64(e)
}

func (e *BigIntegerField) Set(d int64) {
	*e = BigIntegerField(d)
}

func (e *BigIntegerField) String() string {
	return ToStr(e.Value())
}

func (e *BigIntegerField) FieldType() int {
	return TypeBigIntegerField
}

func (e *BigIntegerField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case int64:
		e.Set(d)
	case string:
		v, err := StrTo(d).Int64()
		if err != nil {
			e.Set(v)
		}
	default:
		return errors.New(fmt.Sprintf("<BigIntegerField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *BigIntegerField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(BigIntegerField)

// 0 to 65535
type PositiveSmallIntegerField uint16

func (e PositiveSmallIntegerField) Value() uint16 {
	return uint16(e)
}

func (e *PositiveSmallIntegerField) Set(d uint16) {
	*e = PositiveSmallIntegerField(d)
}

func (e *PositiveSmallIntegerField) String() string {
	return ToStr(e.Value())
}

func (e *PositiveSmallIntegerField) FieldType() int {
	return TypePositiveSmallIntegerField
}

func (e *PositiveSmallIntegerField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case uint16:
		e.Set(d)
	case string:
		v, err := StrTo(d).Uint16()
		if err != nil {
			e.Set(v)
		}
	default:
		return errors.New(fmt.Sprintf("<PositiveSmallIntegerField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *PositiveSmallIntegerField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(PositiveSmallIntegerField)

// 0 to 4294967295
type PositiveIntegerField uint32

func (e PositiveIntegerField) Value() uint32 {
	return uint32(e)
}

func (e *PositiveIntegerField) Set(d uint32) {
	*e = PositiveIntegerField(d)
}

func (e *PositiveIntegerField) String() string {
	return ToStr(e.Value())
}

func (e *PositiveIntegerField) FieldType() int {
	return TypePositiveIntegerField
}

func (e *PositiveIntegerField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case uint32:
		e.Set(d)
	case string:
		v, err := StrTo(d).Uint32()
		if err != nil {
			e.Set(v)
		}
	default:
		return errors.New(fmt.Sprintf("<PositiveIntegerField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *PositiveIntegerField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(PositiveIntegerField)

// 0 to 18446744073709551615
type PositiveBigIntegerField uint64

func (e PositiveBigIntegerField) Value() uint64 {
	return uint64(e)
}

func (e *PositiveBigIntegerField) Set(d uint64) {
	*e = PositiveBigIntegerField(d)
}

func (e *PositiveBigIntegerField) String() string {
	return ToStr(e.Value())
}

func (e *PositiveBigIntegerField) FieldType() int {
	return TypePositiveIntegerField
}

func (e *PositiveBigIntegerField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case uint64:
		e.Set(d)
	case string:
		v, err := StrTo(d).Uint64()
		if err != nil {
			e.Set(v)
		}
	default:
		return errors.New(fmt.Sprintf("<PositiveBigIntegerField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *PositiveBigIntegerField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(PositiveBigIntegerField)

// A large text field.
type TextField string

func (e TextField) Value() string {
	return string(e)
}

func (e *TextField) Set(d string) {
	*e = TextField(d)
}

func (e *TextField) String() string {
	return e.Value()
}

func (e *TextField) FieldType() int {
	return TypeTextField
}

func (e *TextField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		e.Set(d)
	default:
		return errors.New(fmt.Sprintf("<TextField.SetRaw> unknown value `%s`", value))
	}
	return nil
}

func (e *TextField) RawValue() interface{} {
	return e.Value()
}

var _ Fielder = new(TextField)
