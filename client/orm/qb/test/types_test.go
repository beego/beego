// Copyright 2021 ecodeclub
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

package test

import (
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

func TestJsonColumn_Scan(t *testing.T) {
	type User struct {
		Name string
	}
	testCases := []struct {
		name    string
		input   any
		wantVal User
		wantErr string
	}{
		{
			name:  "empty string",
			input: ``,
		},
		{
			name:    "no fields",
			input:   `{}`,
			wantVal: User{},
		},
		{
			name:    "string",
			input:   `{"name":"Tom"}`,
			wantVal: User{Name: "Tom"},
		},
		{
			name:  "nil bytes",
			input: []byte(nil),
		},
		{
			name:  "empty bytes",
			input: []byte(""),
		},
		{
			name:    "bytes",
			input:   []byte(`{"name":"Tom"}`),
			wantVal: User{Name: "Tom"},
		},
		{
			name: "nil",
		},
		{
			name:  "empty bytes ptr",
			input: ekit.ToPtr[[]byte]([]byte("")),
		},
		{
			name:    "bytes ptr",
			input:   ekit.ToPtr[[]byte]([]byte(`{"name":"Tom"}`)),
			wantVal: User{Name: "Tom"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := &JsonColumn{}
			err := js.Scan(tc.input)
			if tc.wantErr != "" {
				assert.EqualError(t, err, tc.wantErr)
				return
			} else {
				assert.Nil(t, err)
			}
			_, err = js.Value()
			assert.Nil(t, err)
			assert.EqualValues(t, tc.wantVal, js.Val)
		})
	}
}
