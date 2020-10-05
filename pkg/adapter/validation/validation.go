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

// Package validation for validations
//
//	import (
//		"github.com/astaxie/beego/validation"
//		"log"
//	)
//
//	type User struct {
//		Name string
//		Age int
//	}
//
//	func main() {
//		u := User{"man", 40}
//		valid := validation.Validation{}
//		valid.Required(u.Name, "name")
//		valid.MaxSize(u.Name, 15, "nameMax")
//		valid.Range(u.Age, 0, 140, "age")
//		if valid.HasErrors() {
//			// validation does not pass
//			// print invalid message
//			for _, err := range valid.Errors {
//				log.Println(err.Key, err.Message)
//			}
//		}
//		// or use like this
//		if v := valid.Max(u.Age, 140, "ageMax"); !v.Ok {
//			log.Println(v.Error.Key, v.Error.Message)
//		}
//	}
//
// more info: http://beego.me/docs/mvc/controller/validation.md
package validation

import (
	"fmt"
	"regexp"

	"github.com/astaxie/beego/pkg/core/validation"
)

// ValidFormer valid interface
type ValidFormer interface {
	Valid(*Validation)
}

// Error show the error
type Error validation.Error

// String Returns the Message.
func (e *Error) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Implement Error interface.
// Return e.String()
func (e *Error) Error() string { return e.String() }

// Result is returned from every validation method.
// It provides an indication of success, and a pointer to the Error (if any).
type Result validation.Result

// Key Get Result by given key string.
func (r *Result) Key(key string) *Result {
	if r.Error != nil {
		r.Error.Key = key
	}
	return r
}

// Message Set Result message by string or format string with args
func (r *Result) Message(message string, args ...interface{}) *Result {
	if r.Error != nil {
		if len(args) == 0 {
			r.Error.Message = message
		} else {
			r.Error.Message = fmt.Sprintf(message, args...)
		}
	}
	return r
}

// A Validation context manages data validation and error messages.
type Validation validation.Validation

// Clear Clean all ValidationError.
func (v *Validation) Clear() {
	(*validation.Validation)(v).Clear()
}

// HasErrors Has ValidationError nor not.
func (v *Validation) HasErrors() bool {
	return (*validation.Validation)(v).HasErrors()
}

// ErrorMap Return the errors mapped by key.
// If there are multiple validation errors associated with a single key, the
// first one "wins".  (Typically the first validation will be the more basic).
func (v *Validation) ErrorMap() map[string][]*Error {
	newErrors := (*validation.Validation)(v).ErrorMap()
	res := make(map[string][]*Error, len(newErrors))
	for n, es := range newErrors {
		errs := make([]*Error, 0, len(es))

		for _, e := range es {
			errs = append(errs, (*Error)(e))
		}

		res[n] = errs
	}
	return res
}

// Error Add an error to the validation context.
func (v *Validation) Error(message string, args ...interface{}) *Result {
	return (*Result)((*validation.Validation)(v).Error(message, args...))
}

// Required Test that the argument is non-nil and non-empty (if string or list)
func (v *Validation) Required(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Required(obj, key))
}

// Min Test that the obj is greater than min if obj's type is int
func (v *Validation) Min(obj interface{}, min int, key string) *Result {
	return (*Result)((*validation.Validation)(v).Min(obj, min, key))
}

// Max Test that the obj is less than max if obj's type is int
func (v *Validation) Max(obj interface{}, max int, key string) *Result {
	return (*Result)((*validation.Validation)(v).Max(obj, max, key))
}

// Range Test that the obj is between mni and max if obj's type is int
func (v *Validation) Range(obj interface{}, min, max int, key string) *Result {
	return (*Result)((*validation.Validation)(v).Range(obj, min, max, key))
}

// MinSize Test that the obj is longer than min size if type is string or slice
func (v *Validation) MinSize(obj interface{}, min int, key string) *Result {
	return (*Result)((*validation.Validation)(v).MinSize(obj, min, key))
}

// MaxSize Test that the obj is shorter than max size if type is string or slice
func (v *Validation) MaxSize(obj interface{}, max int, key string) *Result {
	return (*Result)((*validation.Validation)(v).MaxSize(obj, max, key))
}

// Length Test that the obj is same length to n if type is string or slice
func (v *Validation) Length(obj interface{}, n int, key string) *Result {
	return (*Result)((*validation.Validation)(v).Length(obj, n, key))
}

// Alpha Test that the obj is [a-zA-Z] if type is string
func (v *Validation) Alpha(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Alpha(obj, key))
}

// Numeric Test that the obj is [0-9] if type is string
func (v *Validation) Numeric(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Numeric(obj, key))
}

// AlphaNumeric Test that the obj is [0-9a-zA-Z] if type is string
func (v *Validation) AlphaNumeric(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).AlphaNumeric(obj, key))
}

// Match Test that the obj matches regexp if type is string
func (v *Validation) Match(obj interface{}, regex *regexp.Regexp, key string) *Result {
	return (*Result)((*validation.Validation)(v).Match(obj, regex, key))
}

// NoMatch Test that the obj doesn't match regexp if type is string
func (v *Validation) NoMatch(obj interface{}, regex *regexp.Regexp, key string) *Result {
	return (*Result)((*validation.Validation)(v).NoMatch(obj, regex, key))
}

// AlphaDash Test that the obj is [0-9a-zA-Z_-] if type is string
func (v *Validation) AlphaDash(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).AlphaDash(obj, key))
}

// Email Test that the obj is email address if type is string
func (v *Validation) Email(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Email(obj, key))
}

// IP Test that the obj is IP address if type is string
func (v *Validation) IP(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).IP(obj, key))
}

// Base64 Test that the obj is base64 encoded if type is string
func (v *Validation) Base64(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Base64(obj, key))
}

// Mobile Test that the obj is chinese mobile number if type is string
func (v *Validation) Mobile(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Mobile(obj, key))
}

// Tel Test that the obj is chinese telephone number if type is string
func (v *Validation) Tel(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Tel(obj, key))
}

// Phone Test that the obj is chinese mobile or telephone number if type is string
func (v *Validation) Phone(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).Phone(obj, key))
}

// ZipCode Test that the obj is chinese zip code if type is string
func (v *Validation) ZipCode(obj interface{}, key string) *Result {
	return (*Result)((*validation.Validation)(v).ZipCode(obj, key))
}

// key must like aa.bb.cc or aa.bb.
// AddError adds independent error message for the provided key
func (v *Validation) AddError(key, message string) {
	(*validation.Validation)(v).AddError(key, message)
}

// SetError Set error message for one field in ValidationError
func (v *Validation) SetError(fieldName string, errMsg string) *Error {
	return (*Error)((*validation.Validation)(v).SetError(fieldName, errMsg))
}

// Check Apply a group of validators to a field, in order, and return the
// ValidationResult from the first one that fails, or the last one that
// succeeds.
func (v *Validation) Check(obj interface{}, checks ...Validator) *Result {
	vldts := make([]validation.Validator, 0, len(checks))
	for _, v := range checks {
		vldts = append(vldts, validation.Validator(v))
	}
	return (*Result)((*validation.Validation)(v).Check(obj, vldts...))
}

// Valid Validate a struct.
// the obj parameter must be a struct or a struct pointer
func (v *Validation) Valid(obj interface{}) (b bool, err error) {
	return (*validation.Validation)(v).Valid(obj)
}

// RecursiveValid Recursively validate a struct.
// Step1: Validate by v.Valid
// Step2: If pass on step1, then reflect obj's fields
// Step3: Do the Recursively validation to all struct or struct pointer fields
func (v *Validation) RecursiveValid(objc interface{}) (bool, error) {
	return (*validation.Validation)(v).RecursiveValid(objc)
}

func (v *Validation) CanSkipAlso(skipFunc string) {
	(*validation.Validation)(v).CanSkipAlso(skipFunc)
}
