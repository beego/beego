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
	"reflect"

	"github.com/beego/beego/v2/core/validation"
)

const (
	// ValidTag struct tag
	ValidTag = validation.ValidTag

	LabelTag = validation.LabelTag
)

var ErrInt64On32 = validation.ErrInt64On32

// CustomFunc is for custom validate function
type CustomFunc func(v *Validation, obj interface{}, key string)

// AddCustomFunc Add a custom function to validation
// The name can not be:
//
//	Clear
//	HasErrors
//	ErrorMap
//	Error
//	Check
//	Valid
//	NoMatch
//
// If the name is same with exists function, it will replace the origin valid function
func AddCustomFunc(name string, f CustomFunc) error {
	return validation.AddCustomFunc(name, func(v *validation.Validation, obj interface{}, key string) {
		f((*Validation)(v), obj, key)
	})
}

// ValidFunc Valid function type
type ValidFunc validation.ValidFunc

// Funcs Validate function map
type Funcs validation.Funcs

// Call validate values with named type string
func (f Funcs) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	return (validation.Funcs(f)).Call(name, params...)
}
