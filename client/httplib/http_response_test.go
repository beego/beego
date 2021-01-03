// Copyright 2020 beego
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

package httplib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNewHttpResponseWithJsonBody(t *testing.T) {
	// string
	resp := NewHttpResponseWithJsonBody("{}")
	assert.Equal(t, int64(2), resp.ContentLength)

	resp = NewHttpResponseWithJsonBody([]byte("{}"))
	assert.Equal(t, int64(2), resp.ContentLength)

	resp = NewHttpResponseWithJsonBody(&user{
		Name: "Tom",
	})
	assert.True(t, resp.ContentLength > 0)
}
