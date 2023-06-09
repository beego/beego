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
	"github.com/beego/beego/v2/client/orm/internal/models"
)

// Define the Type enum
const (
	TypeBooleanField              = models.TypeBooleanField
	TypeVarCharField              = models.TypeVarCharField
	TypeCharField                 = models.TypeCharField
	TypeTextField                 = models.TypeTextField
	TypeTimeField                 = models.TypeTimeField
	TypeDateField                 = models.TypeDateField
	TypeDateTimeField             = models.TypeDateTimeField
	TypeBitField                  = models.TypeBitField
	TypeSmallIntegerField         = models.TypeSmallIntegerField
	TypeIntegerField              = models.TypeIntegerField
	TypeBigIntegerField           = models.TypeBigIntegerField
	TypePositiveBitField          = models.TypePositiveBitField
	TypePositiveSmallIntegerField = models.TypePositiveSmallIntegerField
	TypePositiveIntegerField      = models.TypePositiveIntegerField
	TypePositiveBigIntegerField   = models.TypePositiveBigIntegerField
	TypeFloatField                = models.TypeFloatField
	TypeDecimalField              = models.TypeDecimalField
	TypeJSONField                 = models.TypeJSONField
	TypeJsonbField                = models.TypeJsonbField
	RelForeignKey                 = models.RelForeignKey
	RelOneToOne                   = models.RelOneToOne
	RelManyToMany                 = models.RelManyToMany
	RelReverseOne                 = models.RelReverseOne
	RelReverseMany                = models.RelReverseMany
)

// Define some logic enum
const (
	IsIntegerField         = models.IsIntegerField
	IsPositiveIntegerField = models.IsPositiveIntegerField
	IsRelField             = models.IsRelField
	IsFieldType            = models.IsFieldType
)

// BooleanField A true/false field.
type BooleanField = models.BooleanField

// verify the BooleanField implement the Fielder interface
var _ Fielder = new(BooleanField)

// CharField A string field
// required values tag: size
// The size is enforced at the database level and in models’s validation.
// eg: `orm:"size(120)"`
type CharField = models.CharField

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
type TimeField = models.TimeField

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
type DateField = models.DateField

// verify DateField implement fielder interface
var _ Fielder = new(DateField)

// DateTimeField A date, represented in go by a time.Time instance.
// datetime values like 2006-01-02 15:04:05
// Takes the same extra arguments as DateField.
type DateTimeField = models.DateTimeField

// verify datetime implement fielder
var _ models.Fielder = new(DateTimeField)

// FloatField A floating-point number represented in go by a float32 value.
type FloatField = models.FloatField

// verify FloatField implement Fielder
var _ Fielder = new(FloatField)

// SmallIntegerField -32768 to 32767
type SmallIntegerField = models.SmallIntegerField

// verify SmallIntegerField implement Fielder
var _ Fielder = new(SmallIntegerField)

// IntegerField -2147483648 to 2147483647
type IntegerField = models.IntegerField

// verify IntegerField implement Fielder
var _ Fielder = new(IntegerField)

// BigIntegerField -9223372036854775808 to 9223372036854775807.
type BigIntegerField = models.BigIntegerField

// verify BigIntegerField implement Fielder
var _ Fielder = new(BigIntegerField)

// PositiveSmallIntegerField 0 to 65535
type PositiveSmallIntegerField = models.PositiveSmallIntegerField

// verify PositiveSmallIntegerField implement Fielder
var _ Fielder = new(PositiveSmallIntegerField)

// PositiveIntegerField 0 to 4294967295
type PositiveIntegerField = models.PositiveIntegerField

// verify PositiveIntegerField implement Fielder
var _ Fielder = new(PositiveIntegerField)

// PositiveBigIntegerField 0 to 18446744073709551615
type PositiveBigIntegerField = models.PositiveBigIntegerField

// verify PositiveBigIntegerField implement Fielder
var _ Fielder = new(PositiveBigIntegerField)

// TextField A large text field.
type TextField = models.TextField

// verify TextField implement Fielder
var _ Fielder = new(TextField)

// JSONField postgres json field.
type JSONField = models.JSONField

// verify JSONField implement Fielder
var _ models.Fielder = new(JSONField)

// JsonbField postgres json field.
type JsonbField = models.JsonbField

// verify JsonbField implement Fielder
var _ models.Fielder = new(JsonbField)
