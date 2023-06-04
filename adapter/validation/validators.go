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

package validation

import (
	"sync"

	"github.com/beego/beego/v2/core/validation"
)

// CanSkipFuncs will skip valid if RequiredFirst is true and the struct field's value is empty
var CanSkipFuncs = validation.CanSkipFuncs

// MessageTmpls store commond validate template
var MessageTmpls = map[string]string{
	"Required":     "Can not be empty",
	"Min":          "Minimum is %d",
	"Max":          "Maximum is %d",
	"Range":        "Range is %d to %d",
	"MinSize":      "Minimum size is %d",
	"MaxSize":      "Maximum size is %d",
	"Length":       "Required length is %d",
	"Alpha":        "Must be valid alpha characters",
	"Numeric":      "Must be valid numeric characters",
	"AlphaNumeric": "Must be valid alpha or numeric characters",
	"Match":        "Must match %s",
	"NoMatch":      "Must not match %s",
	"AlphaDash":    "Must be valid alpha or numeric or dash(-_) characters",
	"Email":        "Must be a valid email address",
	"IP":           "Must be a valid ip address",
	"Base64":       "Must be valid base64 characters",
	"Mobile":       "Must be valid mobile number",
	"Tel":          "Must be valid telephone number",
	"Phone":        "Must be valid telephone or mobile phone number",
	"ZipCode":      "Must be valid zipcode",
}

var once sync.Once

// SetDefaultMessage set default messages
// if not set, the default messages are
//
//	"Required":     "Can not be empty",
//	"Min":          "Minimum is %d",
//	"Max":          "Maximum is %d",
//	"Range":        "Range is %d to %d",
//	"MinSize":      "Minimum size is %d",
//	"MaxSize":      "Maximum size is %d",
//	"Length":       "Required length is %d",
//	"Alpha":        "Must be valid alpha characters",
//	"Numeric":      "Must be valid numeric characters",
//	"AlphaNumeric": "Must be valid alpha or numeric characters",
//	"Match":        "Must match %s",
//	"NoMatch":      "Must not match %s",
//	"AlphaDash":    "Must be valid alpha or numeric or dash(-_) characters",
//	"Email":        "Must be a valid email address",
//	"IP":           "Must be a valid ip address",
//	"Base64":       "Must be valid base64 characters",
//	"Mobile":       "Must be valid mobile number",
//	"Tel":          "Must be valid telephone number",
//	"Phone":        "Must be valid telephone or mobile phone number",
//	"ZipCode":      "Must be valid zipcode",
func SetDefaultMessage(msg map[string]string) {
	validation.SetDefaultMessage(msg)
}

// Validator interface
type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
	GetKey() string
	GetLimitValue() interface{}
}

// Required struct
type Required validation.Required

// IsSatisfied judge whether obj has value
func (r Required) IsSatisfied(obj interface{}) bool {
	return validation.Required(r).IsSatisfied(obj)
}

// DefaultMessage return the default error message
func (r Required) DefaultMessage() string {
	return validation.Required(r).DefaultMessage()
}

// GetKey return the r.Key
func (r Required) GetKey() string {
	return validation.Required(r).GetKey()
}

// GetLimitValue return nil now
func (r Required) GetLimitValue() interface{} {
	return validation.Required(r).GetLimitValue()
}

// Min check struct
type Min validation.Min

// IsSatisfied judge whether obj is valid
// not support int64 on 32-bit platform
func (m Min) IsSatisfied(obj interface{}) bool {
	return validation.Min(m).IsSatisfied(obj)
}

// DefaultMessage return the default min error message
func (m Min) DefaultMessage() string {
	return validation.Min(m).DefaultMessage()
}

// GetKey return the m.Key
func (m Min) GetKey() string {
	return validation.Min(m).GetKey()
}

// GetLimitValue return the limit value, Min
func (m Min) GetLimitValue() interface{} {
	return validation.Min(m).GetLimitValue()
}

// Max validate struct
type Max validation.Max

// IsSatisfied judge whether obj is valid
// not support int64 on 32-bit platform
func (m Max) IsSatisfied(obj interface{}) bool {
	return validation.Max(m).IsSatisfied(obj)
}

// DefaultMessage return the default max error message
func (m Max) DefaultMessage() string {
	return validation.Max(m).DefaultMessage()
}

// GetKey return the m.Key
func (m Max) GetKey() string {
	return validation.Max(m).GetKey()
}

// GetLimitValue return the limit value, Max
func (m Max) GetLimitValue() interface{} {
	return validation.Max(m).GetLimitValue()
}

// Range Requires an integer to be within Min, Max inclusive.
type Range validation.Range

// IsSatisfied judge whether obj is valid
// not support int64 on 32-bit platform
func (r Range) IsSatisfied(obj interface{}) bool {
	return validation.Range(r).IsSatisfied(obj)
}

// DefaultMessage return the default Range error message
func (r Range) DefaultMessage() string {
	return validation.Range(r).DefaultMessage()
}

// GetKey return the m.Key
func (r Range) GetKey() string {
	return validation.Range(r).GetKey()
}

// GetLimitValue return the limit value, Max
func (r Range) GetLimitValue() interface{} {
	return validation.Range(r).GetLimitValue()
}

// MinSize Requires an array or string to be at least a given length.
type MinSize validation.MinSize

// IsSatisfied judge whether obj is valid
func (m MinSize) IsSatisfied(obj interface{}) bool {
	return validation.MinSize(m).IsSatisfied(obj)
}

// DefaultMessage return the default MinSize error message
func (m MinSize) DefaultMessage() string {
	return validation.MinSize(m).DefaultMessage()
}

// GetKey return the m.Key
func (m MinSize) GetKey() string {
	return validation.MinSize(m).GetKey()
}

// GetLimitValue return the limit value
func (m MinSize) GetLimitValue() interface{} {
	return validation.MinSize(m).GetLimitValue()
}

// MaxSize Requires an array or string to be at most a given length.
type MaxSize validation.MaxSize

// IsSatisfied judge whether obj is valid
func (m MaxSize) IsSatisfied(obj interface{}) bool {
	return validation.MaxSize(m).IsSatisfied(obj)
}

// DefaultMessage return the default MaxSize error message
func (m MaxSize) DefaultMessage() string {
	return validation.MaxSize(m).DefaultMessage()
}

// GetKey return the m.Key
func (m MaxSize) GetKey() string {
	return validation.MaxSize(m).GetKey()
}

// GetLimitValue return the limit value
func (m MaxSize) GetLimitValue() interface{} {
	return validation.MaxSize(m).GetLimitValue()
}

// Length Requires an array or string to be exactly a given length.
type Length validation.Length

// IsSatisfied judge whether obj is valid
func (l Length) IsSatisfied(obj interface{}) bool {
	return validation.Length(l).IsSatisfied(obj)
}

// DefaultMessage return the default Length error message
func (l Length) DefaultMessage() string {
	return validation.Length(l).DefaultMessage()
}

// GetKey return the m.Key
func (l Length) GetKey() string {
	return validation.Length(l).GetKey()
}

// GetLimitValue return the limit value
func (l Length) GetLimitValue() interface{} {
	return validation.Length(l).GetLimitValue()
}

// Alpha check the alpha
type Alpha validation.Alpha

// IsSatisfied judge whether obj is valid
func (a Alpha) IsSatisfied(obj interface{}) bool {
	return validation.Alpha(a).IsSatisfied(obj)
}

// DefaultMessage return the default Length error message
func (a Alpha) DefaultMessage() string {
	return validation.Alpha(a).DefaultMessage()
}

// GetKey return the m.Key
func (a Alpha) GetKey() string {
	return validation.Alpha(a).GetKey()
}

// GetLimitValue return the limit value
func (a Alpha) GetLimitValue() interface{} {
	return validation.Alpha(a).GetLimitValue()
}

// Numeric check number
type Numeric validation.Numeric

// IsSatisfied judge whether obj is valid
func (n Numeric) IsSatisfied(obj interface{}) bool {
	return validation.Numeric(n).IsSatisfied(obj)
}

// DefaultMessage return the default Length error message
func (n Numeric) DefaultMessage() string {
	return validation.Numeric(n).DefaultMessage()
}

// GetKey return the n.Key
func (n Numeric) GetKey() string {
	return validation.Numeric(n).GetKey()
}

// GetLimitValue return the limit value
func (n Numeric) GetLimitValue() interface{} {
	return validation.Numeric(n).GetLimitValue()
}

// AlphaNumeric check alpha and number
type AlphaNumeric validation.AlphaNumeric

// IsSatisfied judge whether obj is valid
func (a AlphaNumeric) IsSatisfied(obj interface{}) bool {
	return validation.AlphaNumeric(a).IsSatisfied(obj)
}

// DefaultMessage return the default Length error message
func (a AlphaNumeric) DefaultMessage() string {
	return validation.AlphaNumeric(a).DefaultMessage()
}

// GetKey return the a.Key
func (a AlphaNumeric) GetKey() string {
	return validation.AlphaNumeric(a).GetKey()
}

// GetLimitValue return the limit value
func (a AlphaNumeric) GetLimitValue() interface{} {
	return validation.AlphaNumeric(a).GetLimitValue()
}

// Match Requires a string to match a given regex.
type Match validation.Match

// IsSatisfied judge whether obj is valid
func (m Match) IsSatisfied(obj interface{}) bool {
	return validation.Match(m).IsSatisfied(obj)
}

// DefaultMessage return the default Match error message
func (m Match) DefaultMessage() string {
	return validation.Match(m).DefaultMessage()
}

// GetKey return the m.Key
func (m Match) GetKey() string {
	return validation.Match(m).GetKey()
}

// GetLimitValue return the limit value
func (m Match) GetLimitValue() interface{} {
	return validation.Match(m).GetLimitValue()
}

// NoMatch Requires a string to not match a given regex.
type NoMatch validation.NoMatch

// IsSatisfied judge whether obj is valid
func (n NoMatch) IsSatisfied(obj interface{}) bool {
	return validation.NoMatch(n).IsSatisfied(obj)
}

// DefaultMessage return the default NoMatch error message
func (n NoMatch) DefaultMessage() string {
	return validation.NoMatch(n).DefaultMessage()
}

// GetKey return the n.Key
func (n NoMatch) GetKey() string {
	return validation.NoMatch(n).GetKey()
}

// GetLimitValue return the limit value
func (n NoMatch) GetLimitValue() interface{} {
	return validation.NoMatch(n).GetLimitValue()
}

// AlphaDash check not Alpha
type AlphaDash validation.AlphaDash

// DefaultMessage return the default AlphaDash error message
func (a AlphaDash) DefaultMessage() string {
	return validation.AlphaDash(a).DefaultMessage()
}

// GetKey return the n.Key
func (a AlphaDash) GetKey() string {
	return validation.AlphaDash(a).GetKey()
}

// GetLimitValue return the limit value
func (a AlphaDash) GetLimitValue() interface{} {
	return validation.AlphaDash(a).GetLimitValue()
}

// Email check struct
type Email validation.Email

// DefaultMessage return the default Email error message
func (e Email) DefaultMessage() string {
	return validation.Email(e).DefaultMessage()
}

// GetKey return the n.Key
func (e Email) GetKey() string {
	return validation.Email(e).GetKey()
}

// GetLimitValue return the limit value
func (e Email) GetLimitValue() interface{} {
	return validation.Email(e).GetLimitValue()
}

// IP check struct
type IP validation.IP

// DefaultMessage return the default IP error message
func (i IP) DefaultMessage() string {
	return validation.IP(i).DefaultMessage()
}

// GetKey return the i.Key
func (i IP) GetKey() string {
	return validation.IP(i).GetKey()
}

// GetLimitValue return the limit value
func (i IP) GetLimitValue() interface{} {
	return validation.IP(i).GetLimitValue()
}

// Base64 check struct
type Base64 validation.Base64

// DefaultMessage return the default Base64 error message
func (b Base64) DefaultMessage() string {
	return validation.Base64(b).DefaultMessage()
}

// GetKey return the b.Key
func (b Base64) GetKey() string {
	return validation.Base64(b).GetKey()
}

// GetLimitValue return the limit value
func (b Base64) GetLimitValue() interface{} {
	return validation.Base64(b).GetLimitValue()
}

// Mobile check struct
type Mobile validation.Mobile

// DefaultMessage return the default Mobile error message
func (m Mobile) DefaultMessage() string {
	return validation.Mobile(m).DefaultMessage()
}

// GetKey return the m.Key
func (m Mobile) GetKey() string {
	return validation.Mobile(m).GetKey()
}

// GetLimitValue return the limit value
func (m Mobile) GetLimitValue() interface{} {
	return validation.Mobile(m).GetLimitValue()
}

// Tel check telephone struct
type Tel validation.Tel

// DefaultMessage return the default Tel error message
func (t Tel) DefaultMessage() string {
	return validation.Tel(t).DefaultMessage()
}

// GetKey return the t.Key
func (t Tel) GetKey() string {
	return validation.Tel(t).GetKey()
}

// GetLimitValue return the limit value
func (t Tel) GetLimitValue() interface{} {
	return validation.Tel(t).GetLimitValue()
}

// Phone just for chinese telephone or mobile phone number
type Phone validation.Phone

// IsSatisfied judge whether obj is valid
func (p Phone) IsSatisfied(obj interface{}) bool {
	return validation.Phone(p).IsSatisfied(obj)
}

// DefaultMessage return the default Phone error message
func (p Phone) DefaultMessage() string {
	return validation.Phone(p).DefaultMessage()
}

// GetKey return the p.Key
func (p Phone) GetKey() string {
	return validation.Phone(p).GetKey()
}

// GetLimitValue return the limit value
func (p Phone) GetLimitValue() interface{} {
	return validation.Phone(p).GetLimitValue()
}

// ZipCode check the zip struct
type ZipCode validation.ZipCode

// DefaultMessage return the default Zip error message
func (z ZipCode) DefaultMessage() string {
	return validation.ZipCode(z).DefaultMessage()
}

// GetKey return the z.Key
func (z ZipCode) GetKey() string {
	return validation.ZipCode(z).GetKey()
}

// GetLimitValue return the limit value
func (z ZipCode) GetLimitValue() interface{} {
	return validation.ZipCode(z).GetLimitValue()
}
