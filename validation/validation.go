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

// package for validations
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
//		if v := valid.Max(u.Age, 140); !v.Ok {
//			log.Println(v.Error.Key, v.Error.Message)
//		}
//	}
//
// more info: http://beego.me/docs/mvc/controller/validation.md
package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ValidFormer interface {
	Valid(*Validation)
}

type ValidationError struct {
	Message, Key, Name, Field, Tmpl string
	Value                           interface{}
	LimitValue                      interface{}
}

// Returns the Message.
func (e *ValidationError) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// A ValidationResult is returned from every validation method.
// It provides an indication of success, and a pointer to the Error (if any).
type ValidationResult struct {
	Error *ValidationError
	Ok    bool
}

// Get ValidationResult by given key string.
func (r *ValidationResult) Key(key string) *ValidationResult {
	if r.Error != nil {
		r.Error.Key = key
	}
	return r
}

// Set ValidationResult message by string or format string with args
func (r *ValidationResult) Message(message string, args ...interface{}) *ValidationResult {
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
type Validation struct {
	Errors    []*ValidationError
	ErrorsMap map[string]*ValidationError
}

// Clean all ValidationError.
func (v *Validation) Clear() {
	v.Errors = []*ValidationError{}
}

// Has ValidationError nor not.
func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

// Return the errors mapped by key.
// If there are multiple validation errors associated with a single key, the
// first one "wins".  (Typically the first validation will be the more basic).
func (v *Validation) ErrorMap() map[string]*ValidationError {
	return v.ErrorsMap
}

// Add an error to the validation context.
func (v *Validation) Error(message string, args ...interface{}) *ValidationResult {
	result := (&ValidationResult{
		Ok:    false,
		Error: &ValidationError{},
	}).Message(message, args...)
	v.Errors = append(v.Errors, result.Error)
	return result
}

// Test that the argument is non-nil and non-empty (if string or list)
func (v *Validation) Required(obj interface{}, key string) *ValidationResult {
	return v.apply(Required{key}, obj)
}

// Test that the obj is greater than min if obj's type is int
func (v *Validation) Min(obj interface{}, min int, key string) *ValidationResult {
	return v.apply(Min{min, key}, obj)
}

// Test that the obj is less than max if obj's type is int
func (v *Validation) Max(obj interface{}, max int, key string) *ValidationResult {
	return v.apply(Max{max, key}, obj)
}

// Test that the obj is between mni and max if obj's type is int
func (v *Validation) Range(obj interface{}, min, max int, key string) *ValidationResult {
	return v.apply(Range{Min{Min: min}, Max{Max: max}, key}, obj)
}

// Test that the obj is longer than min size if type is string or slice
func (v *Validation) MinSize(obj interface{}, min int, key string) *ValidationResult {
	return v.apply(MinSize{min, key}, obj)
}

// Test that the obj is shorter than max size if type is string or slice
func (v *Validation) MaxSize(obj interface{}, max int, key string) *ValidationResult {
	return v.apply(MaxSize{max, key}, obj)
}

// Test that the obj is same length to n if type is string or slice
func (v *Validation) Length(obj interface{}, n int, key string) *ValidationResult {
	return v.apply(Length{n, key}, obj)
}

// Test that the obj is [a-zA-Z] if type is string
func (v *Validation) Alpha(obj interface{}, key string) *ValidationResult {
	return v.apply(Alpha{key}, obj)
}

// Test that the obj is [0-9] if type is string
func (v *Validation) Numeric(obj interface{}, key string) *ValidationResult {
	return v.apply(Numeric{key}, obj)
}

// Test that the obj is [0-9a-zA-Z] if type is string
func (v *Validation) AlphaNumeric(obj interface{}, key string) *ValidationResult {
	return v.apply(AlphaNumeric{key}, obj)
}

// Test that the obj matches regexp if type is string
func (v *Validation) Match(obj interface{}, regex *regexp.Regexp, key string) *ValidationResult {
	return v.apply(Match{regex, key}, obj)
}

// Test that the obj doesn't match regexp if type is string
func (v *Validation) NoMatch(obj interface{}, regex *regexp.Regexp, key string) *ValidationResult {
	return v.apply(NoMatch{Match{Regexp: regex}, key}, obj)
}

// Test that the obj is [0-9a-zA-Z_-] if type is string
func (v *Validation) AlphaDash(obj interface{}, key string) *ValidationResult {
	return v.apply(AlphaDash{NoMatch{Match: Match{Regexp: alphaDashPattern}}, key}, obj)
}

// Test that the obj is email address if type is string
func (v *Validation) Email(obj interface{}, key string) *ValidationResult {
	return v.apply(Email{Match{Regexp: emailPattern}, key}, obj)
}

// Test that the obj is IP address if type is string
func (v *Validation) IP(obj interface{}, key string) *ValidationResult {
	return v.apply(IP{Match{Regexp: ipPattern}, key}, obj)
}

// Test that the obj is base64 encoded if type is string
func (v *Validation) Base64(obj interface{}, key string) *ValidationResult {
	return v.apply(Base64{Match{Regexp: base64Pattern}, key}, obj)
}

// Test that the obj is chinese mobile number if type is string
func (v *Validation) Mobile(obj interface{}, key string) *ValidationResult {
	return v.apply(Mobile{Match{Regexp: mobilePattern}, key}, obj)
}

// Test that the obj is chinese telephone number if type is string
func (v *Validation) Tel(obj interface{}, key string) *ValidationResult {
	return v.apply(Tel{Match{Regexp: telPattern}, key}, obj)
}

// Test that the obj is chinese mobile or telephone number if type is string
func (v *Validation) Phone(obj interface{}, key string) *ValidationResult {
	return v.apply(Phone{Mobile{Match: Match{Regexp: mobilePattern}},
		Tel{Match: Match{Regexp: telPattern}}, key}, obj)
}

// Test that the obj is chinese zip code if type is string
func (v *Validation) ZipCode(obj interface{}, key string) *ValidationResult {
	return v.apply(ZipCode{Match{Regexp: zipCodePattern}, key}, obj)
}

func (v *Validation) apply(chk Validator, obj interface{}) *ValidationResult {
	if chk.IsSatisfied(obj) {
		return &ValidationResult{Ok: true}
	}

	// Add the error to the validation context.
	key := chk.GetKey()
	Name := key
	Field := ""

	parts := strings.Split(key, ".")
	if len(parts) == 2 {
		Field = parts[0]
		Name = parts[1]
	}

	err := &ValidationError{
		Message:    chk.DefaultMessage(),
		Key:        key,
		Name:       Name,
		Field:      Field,
		Value:      obj,
		Tmpl:       MessageTmpls[Name],
		LimitValue: chk.GetLimitValue(),
	}
	v.setError(err)

	// Also return it in the result.
	return &ValidationResult{
		Ok:    false,
		Error: err,
	}
}

func (v *Validation) setError(err *ValidationError) {
	v.Errors = append(v.Errors, err)
	if v.ErrorsMap == nil {
		v.ErrorsMap = make(map[string]*ValidationError)
	}
	if _, ok := v.ErrorsMap[err.Field]; !ok {
		v.ErrorsMap[err.Field] = err
	}
}

// Set error message for one field in ValidationError
func (v *Validation) SetError(fieldName string, errMsg string) *ValidationError {
	err := &ValidationError{Key: fieldName, Field: fieldName, Tmpl: errMsg, Message: errMsg}
	v.setError(err)
	return err
}

// Apply a group of validators to a field, in order, and return the
// ValidationResult from the first one that fails, or the last one that
// succeeds.
func (v *Validation) Check(obj interface{}, checks ...Validator) *ValidationResult {
	var result *ValidationResult
	for _, check := range checks {
		result = v.apply(check, obj)
		if !result.Ok {
			return result
		}
	}
	return result
}

// Validate a struct.
// the obj parameter must be a struct or a struct pointer
func (v *Validation) Valid(obj interface{}) (b bool, err error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		err = fmt.Errorf("%v must be a struct or a struct pointer", obj)
		return
	}

	for i := 0; i < objT.NumField(); i++ {
		var vfs []ValidFunc
		if vfs, err = getValidFuncs(objT.Field(i)); err != nil {
			return
		}
		for _, vf := range vfs {
			if _, err = funcs.Call(vf.Name,
				mergeParam(v, objV.Field(i).Interface(), vf.Params)...); err != nil {
				return
			}
		}
	}

	if !v.HasErrors() {
		if form, ok := obj.(ValidFormer); ok {
			form.Valid(v)
		}
	}

	return !v.HasErrors(), nil
}
