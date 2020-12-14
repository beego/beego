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
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Define the Type enum
const (
	TypeBooleanField              = orm.TypeBooleanField
	TypeVarCharField              = orm.TypeVarCharField
	TypeCharField                 = orm.TypeCharField
	TypeTextField                 = orm.TypeTextField
	TypeTimeField                 = orm.TypeTimeField
	TypeDateField                 = orm.TypeDateField
	TypeDateTimeField             = orm.TypeDateTimeField
	TypeBitField                  = orm.TypeBitField
	TypeSmallIntegerField         = orm.TypeSmallIntegerField
	TypeIntegerField              = orm.TypeIntegerField
	TypeBigIntegerField           = orm.TypeBigIntegerField
	TypePositiveBitField          = orm.TypePositiveBitField
	TypePositiveSmallIntegerField = orm.TypePositiveSmallIntegerField
	TypePositiveIntegerField      = orm.TypePositiveIntegerField
	TypePositiveBigIntegerField   = orm.TypePositiveBigIntegerField
	TypeFloatField                = orm.TypeFloatField
	TypeDecimalField              = orm.TypeDecimalField
	TypeJSONField                 = orm.TypeJSONField
	TypeJsonbField                = orm.TypeJsonbField
	RelForeignKey                 = orm.RelForeignKey
	RelOneToOne                   = orm.RelOneToOne
	RelManyToMany                 = orm.RelManyToMany
	RelReverseOne                 = orm.RelReverseOne
	RelReverseMany                = orm.RelReverseMany
)

// Define some logic enum
const (
	IsIntegerField         = orm.IsIntegerField
	IsPositiveIntegerField = orm.IsPositiveIntegerField
	IsRelField             = orm.IsRelField
	IsFieldType            = orm.IsFieldType
)

// BooleanField A true/false field.
type BooleanField orm.BooleanField

// Value return the BooleanField
func (e BooleanField) Value() bool {
	return orm.BooleanField(e).Value()
}

// Set will set the BooleanField
func (e *BooleanField) Set(d bool) {
	(*orm.BooleanField)(e).Set(d)
}

// String format the Bool to string
func (e *BooleanField) String() string {
	return (*orm.BooleanField)(e).String()
}

// FieldType return BooleanField the type
func (e *BooleanField) FieldType() int {
	return (*orm.BooleanField)(e).FieldType()
}

// SetRaw set the interface to bool
func (e *BooleanField) SetRaw(value interface{}) error {
	return (*orm.BooleanField)(e).SetRaw(value)
}

// RawValue return the current value
func (e *BooleanField) RawValue() interface{} {
	return (*orm.BooleanField)(e).RawValue()
}

// verify the BooleanField implement the Fielder interface
var _ Fielder = new(BooleanField)

// CharField A string field
// required values tag: size
// The size is enforced at the database level and in models’s validation.
// eg: `orm:"size(120)"`
type CharField orm.CharField

// Value return the CharField's Value
func (e CharField) Value() string {
	return orm.CharField(e).Value()
}

// Set CharField value
func (e *CharField) Set(d string) {
	(*orm.CharField)(e).Set(d)
}

// String return the CharField
func (e *CharField) String() string {
	return (*orm.CharField)(e).String()
}

// FieldType return the enum type
func (e *CharField) FieldType() int {
	return (*orm.CharField)(e).FieldType()
}

// SetRaw set the interface to string
func (e *CharField) SetRaw(value interface{}) error {
	return (*orm.CharField)(e).SetRaw(value)
}

// RawValue return the CharField value
func (e *CharField) RawValue() interface{} {
	return (*orm.CharField)(e).RawValue()
}

// verify CharField implement Fielder
var _ Fielder = new(CharField)

// TimeField A time, represented in go by a time.Time instance.
// only time values like 10:00:00
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
type TimeField orm.TimeField

// Value return the time.Time
func (e TimeField) Value() time.Time {
	return orm.TimeField(e).Value()
}

// Set set the TimeField's value
func (e *TimeField) Set(d time.Time) {
	(*orm.TimeField)(e).Set(d)
}

// String convert time to string
func (e *TimeField) String() string {
	return (*orm.TimeField)(e).String()
}

// FieldType return enum type Date
func (e *TimeField) FieldType() int {
	return (*orm.TimeField)(e).FieldType()
}

// SetRaw convert the interface to time.Time. Allow string and time.Time
func (e *TimeField) SetRaw(value interface{}) error {
	return (*orm.TimeField)(e).SetRaw(value)
}

// RawValue return time value
func (e *TimeField) RawValue() interface{} {
	return (*orm.TimeField)(e).RawValue()
}

var _ Fielder = new(TimeField)

// DateField A date, represented in go by a time.Time instance.
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
type DateField orm.DateField

// Value return the time.Time
func (e DateField) Value() time.Time {
	return orm.DateField(e).Value()
}

// Set set the DateField's value
func (e *DateField) Set(d time.Time) {
	(*orm.DateField)(e).Set(d)
}

// String convert datetime to string
func (e *DateField) String() string {
	return (*orm.DateField)(e).String()
}

// FieldType return enum type Date
func (e *DateField) FieldType() int {
	return (*orm.DateField)(e).FieldType()
}

// SetRaw convert the interface to time.Time. Allow string and time.Time
func (e *DateField) SetRaw(value interface{}) error {
	return (*orm.DateField)(e).SetRaw(value)
}

// RawValue return Date value
func (e *DateField) RawValue() interface{} {
	return (*orm.DateField)(e).RawValue()
}

// verify DateField implement fielder interface
var _ Fielder = new(DateField)

// DateTimeField A date, represented in go by a time.Time instance.
// datetime values like 2006-01-02 15:04:05
// Takes the same extra arguments as DateField.
type DateTimeField orm.DateTimeField

// Value return the datetime value
func (e DateTimeField) Value() time.Time {
	return orm.DateTimeField(e).Value()
}

// Set set the time.Time to datetime
func (e *DateTimeField) Set(d time.Time) {
	(*orm.DateTimeField)(e).Set(d)
}

// String return the time's String
func (e *DateTimeField) String() string {
	return (*orm.DateTimeField)(e).String()
}

// FieldType return the enum TypeDateTimeField
func (e *DateTimeField) FieldType() int {
	return (*orm.DateTimeField)(e).FieldType()
}

// SetRaw convert the string or time.Time to DateTimeField
func (e *DateTimeField) SetRaw(value interface{}) error {
	return (*orm.DateTimeField)(e).SetRaw(value)
}

// RawValue return the datetime value
func (e *DateTimeField) RawValue() interface{} {
	return (*orm.DateTimeField)(e).RawValue()
}

// verify datetime implement fielder
var _ Fielder = new(DateTimeField)

// FloatField A floating-point number represented in go by a float32 value.
type FloatField orm.FloatField

// Value return the FloatField value
func (e FloatField) Value() float64 {
	return orm.FloatField(e).Value()
}

// Set the Float64
func (e *FloatField) Set(d float64) {
	(*orm.FloatField)(e).Set(d)
}

// String return the string
func (e *FloatField) String() string {
	return (*orm.FloatField)(e).String()
}

// FieldType return the enum type
func (e *FloatField) FieldType() int {
	return (*orm.FloatField)(e).FieldType()
}

// SetRaw converter interface Float64 float32 or string to FloatField
func (e *FloatField) SetRaw(value interface{}) error {
	return (*orm.FloatField)(e).SetRaw(value)
}

// RawValue return the FloatField value
func (e *FloatField) RawValue() interface{} {
	return (*orm.FloatField)(e).RawValue()
}

// verify FloatField implement Fielder
var _ Fielder = new(FloatField)

// SmallIntegerField -32768 to 32767
type SmallIntegerField orm.SmallIntegerField

// Value return int16 value
func (e SmallIntegerField) Value() int16 {
	return orm.SmallIntegerField(e).Value()
}

// Set the SmallIntegerField value
func (e *SmallIntegerField) Set(d int16) {
	(*orm.SmallIntegerField)(e).Set(d)
}

// String convert smallint to string
func (e *SmallIntegerField) String() string {
	return (*orm.SmallIntegerField)(e).String()
}

// FieldType return enum type SmallIntegerField
func (e *SmallIntegerField) FieldType() int {
	return (*orm.SmallIntegerField)(e).FieldType()
}

// SetRaw convert interface int16/string to int16
func (e *SmallIntegerField) SetRaw(value interface{}) error {
	return (*orm.SmallIntegerField)(e).SetRaw(value)
}

// RawValue return smallint value
func (e *SmallIntegerField) RawValue() interface{} {
	return (*orm.SmallIntegerField)(e).RawValue()
}

// verify SmallIntegerField implement Fielder
var _ Fielder = new(SmallIntegerField)

// IntegerField -2147483648 to 2147483647
type IntegerField orm.IntegerField

// Value return the int32
func (e IntegerField) Value() int32 {
	return orm.IntegerField(e).Value()
}

// Set IntegerField value
func (e *IntegerField) Set(d int32) {
	(*orm.IntegerField)(e).Set(d)
}

// String convert Int32 to string
func (e *IntegerField) String() string {
	return (*orm.IntegerField)(e).String()
}

// FieldType return the enum type
func (e *IntegerField) FieldType() int {
	return (*orm.IntegerField)(e).FieldType()
}

// SetRaw convert interface int32/string to int32
func (e *IntegerField) SetRaw(value interface{}) error {
	return (*orm.IntegerField)(e).SetRaw(value)
}

// RawValue return IntegerField value
func (e *IntegerField) RawValue() interface{} {
	return (*orm.IntegerField)(e).RawValue()
}

// verify IntegerField implement Fielder
var _ Fielder = new(IntegerField)

// BigIntegerField -9223372036854775808 to 9223372036854775807.
type BigIntegerField orm.BigIntegerField

// Value return int64
func (e BigIntegerField) Value() int64 {
	return orm.BigIntegerField(e).Value()
}

// Set the BigIntegerField value
func (e *BigIntegerField) Set(d int64) {
	(*orm.BigIntegerField)(e).Set(d)
}

// String convert BigIntegerField to string
func (e *BigIntegerField) String() string {
	return (*orm.BigIntegerField)(e).String()
}

// FieldType return enum type
func (e *BigIntegerField) FieldType() int {
	return (*orm.BigIntegerField)(e).FieldType()
}

// SetRaw convert interface int64/string to int64
func (e *BigIntegerField) SetRaw(value interface{}) error {
	return (*orm.BigIntegerField)(e).SetRaw(value)
}

// RawValue return BigIntegerField value
func (e *BigIntegerField) RawValue() interface{} {
	return (*orm.BigIntegerField)(e).RawValue()
}

// verify BigIntegerField implement Fielder
var _ Fielder = new(BigIntegerField)

// PositiveSmallIntegerField 0 to 65535
type PositiveSmallIntegerField orm.PositiveSmallIntegerField

// Value return uint16
func (e PositiveSmallIntegerField) Value() uint16 {
	return orm.PositiveSmallIntegerField(e).Value()
}

// Set PositiveSmallIntegerField value
func (e *PositiveSmallIntegerField) Set(d uint16) {
	(*orm.PositiveSmallIntegerField)(e).Set(d)
}

// String convert uint16 to string
func (e *PositiveSmallIntegerField) String() string {
	return (*orm.PositiveSmallIntegerField)(e).String()
}

// FieldType return enum type
func (e *PositiveSmallIntegerField) FieldType() int {
	return (*orm.PositiveSmallIntegerField)(e).FieldType()
}

// SetRaw convert Interface uint16/string to uint16
func (e *PositiveSmallIntegerField) SetRaw(value interface{}) error {
	return (*orm.PositiveSmallIntegerField)(e).SetRaw(value)
}

// RawValue returns PositiveSmallIntegerField value
func (e *PositiveSmallIntegerField) RawValue() interface{} {
	return (*orm.PositiveSmallIntegerField)(e).RawValue()
}

// verify PositiveSmallIntegerField implement Fielder
var _ Fielder = new(PositiveSmallIntegerField)

// PositiveIntegerField 0 to 4294967295
type PositiveIntegerField orm.PositiveIntegerField

// Value return PositiveIntegerField value. Uint32
func (e PositiveIntegerField) Value() uint32 {
	return orm.PositiveIntegerField(e).Value()
}

// Set the PositiveIntegerField value
func (e *PositiveIntegerField) Set(d uint32) {
	(*orm.PositiveIntegerField)(e).Set(d)
}

// String convert PositiveIntegerField to string
func (e *PositiveIntegerField) String() string {
	return (*orm.PositiveIntegerField)(e).String()
}

// FieldType return enum type
func (e *PositiveIntegerField) FieldType() int {
	return (*orm.PositiveIntegerField)(e).FieldType()
}

// SetRaw convert interface uint32/string to Uint32
func (e *PositiveIntegerField) SetRaw(value interface{}) error {
	return (*orm.PositiveIntegerField)(e).SetRaw(value)
}

// RawValue return the PositiveIntegerField Value
func (e *PositiveIntegerField) RawValue() interface{} {
	return (*orm.PositiveIntegerField)(e).RawValue()
}

// verify PositiveIntegerField implement Fielder
var _ Fielder = new(PositiveIntegerField)

// PositiveBigIntegerField 0 to 18446744073709551615
type PositiveBigIntegerField orm.PositiveBigIntegerField

// Value return uint64
func (e PositiveBigIntegerField) Value() uint64 {
	return orm.PositiveBigIntegerField(e).Value()
}

// Set PositiveBigIntegerField value
func (e *PositiveBigIntegerField) Set(d uint64) {
	(*orm.PositiveBigIntegerField)(e).Set(d)
}

// String convert PositiveBigIntegerField to string
func (e *PositiveBigIntegerField) String() string {
	return (*orm.PositiveBigIntegerField)(e).String()
}

// FieldType return enum type
func (e *PositiveBigIntegerField) FieldType() int {
	return (*orm.PositiveBigIntegerField)(e).FieldType()
}

// SetRaw convert interface uint64/string to Uint64
func (e *PositiveBigIntegerField) SetRaw(value interface{}) error {
	return (*orm.PositiveBigIntegerField)(e).SetRaw(value)
}

// RawValue return PositiveBigIntegerField value
func (e *PositiveBigIntegerField) RawValue() interface{} {
	return (*orm.PositiveBigIntegerField)(e).RawValue()
}

// verify PositiveBigIntegerField implement Fielder
var _ Fielder = new(PositiveBigIntegerField)

// TextField A large text field.
type TextField orm.TextField

// Value return TextField value
func (e TextField) Value() string {
	return orm.TextField(e).Value()
}

// Set the TextField value
func (e *TextField) Set(d string) {
	(*orm.TextField)(e).Set(d)
}

// String convert TextField to string
func (e *TextField) String() string {
	return (*orm.TextField)(e).String()
}

// FieldType return enum type
func (e *TextField) FieldType() int {
	return (*orm.TextField)(e).FieldType()
}

// SetRaw convert interface string to string
func (e *TextField) SetRaw(value interface{}) error {
	return (*orm.TextField)(e).SetRaw(value)
}

// RawValue return TextField value
func (e *TextField) RawValue() interface{} {
	return (*orm.TextField)(e).RawValue()
}

// verify TextField implement Fielder
var _ Fielder = new(TextField)

// JSONField postgres json field.
type JSONField orm.JSONField

// Value return JSONField value
func (j JSONField) Value() string {
	return orm.JSONField(j).Value()
}

// Set the JSONField value
func (j *JSONField) Set(d string) {
	(*orm.JSONField)(j).Set(d)
}

// String convert JSONField to string
func (j *JSONField) String() string {
	return (*orm.JSONField)(j).String()
}

// FieldType return enum type
func (j *JSONField) FieldType() int {
	return (*orm.JSONField)(j).FieldType()
}

// SetRaw convert interface string to string
func (j *JSONField) SetRaw(value interface{}) error {
	return (*orm.JSONField)(j).SetRaw(value)
}

// RawValue return JSONField value
func (j *JSONField) RawValue() interface{} {
	return (*orm.JSONField)(j).RawValue()
}

// verify JSONField implement Fielder
var _ Fielder = new(JSONField)

// JsonbField postgres json field.
type JsonbField orm.JsonbField

// Value return JsonbField value
func (j JsonbField) Value() string {
	return orm.JsonbField(j).Value()
}

// Set the JsonbField value
func (j *JsonbField) Set(d string) {
	(*orm.JsonbField)(j).Set(d)
}

// String convert JsonbField to string
func (j *JsonbField) String() string {
	return (*orm.JsonbField)(j).String()
}

// FieldType return enum type
func (j *JsonbField) FieldType() int {
	return (*orm.JsonbField)(j).FieldType()
}

// SetRaw convert interface string to string
func (j *JsonbField) SetRaw(value interface{}) error {
	return (*orm.JsonbField)(j).SetRaw(value)
}

// RawValue return JsonbField value
func (j *JsonbField) RawValue() interface{} {
	return (*orm.JsonbField)(j).RawValue()
}

// verify JsonbField implement Fielder
var _ Fielder = new(JsonbField)
