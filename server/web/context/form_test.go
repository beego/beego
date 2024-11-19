// Copyright 2024 beego
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

package context

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormValue(t *testing.T) {
	typ := reflect.TypeOf(TestStruct{})
	defField, _ := typ.FieldByName("DefaultField")
	noDefField, _ := typ.FieldByName("NoDefaultField")
	testCases := []struct {
		name string
		tag  string

		form  url.Values
		field reflect.StructField

		wantRes string
		wantOk  bool
	}{
		{
			name:  "use value",
			tag:   "defaultField",
			field: defField,
			form: map[string][]string{
				"defaultField": {"abc", "bcd"},
			},
			wantRes: "abc",
			wantOk:  true,
		},
		{
			name:  "empty value",
			tag:   "defaultField",
			field: defField,
			form: map[string][]string{
				"defaultField": {"", "bcd"},
			},
			wantRes: "",
			wantOk:  true,
		},
		{
			name:    "use default value",
			tag:     "defaultField",
			field:   defField,
			form:    map[string][]string{},
			wantRes: "Tom",
			wantOk:  true,
		},
		{
			name:   "no value",
			tag:    "no",
			field:  noDefField,
			form:   map[string][]string{},
			wantOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := formValue(tc.tag, tc.form, tc.field)
			assert.Equal(t, tc.wantRes, val)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}

type TestStruct struct {
	DefaultField   string `form:"defaultField" default:"Tom"`
	NoDefaultField string `form:"noDefaultField"`
}
