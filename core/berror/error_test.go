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

package berror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCode1 = DefineCode(1, "unit_test", "TestError", "Hello, test code1")

var testErr = errors.New("hello, this is error")

func TestErrorf(t *testing.T) {
	msg := Errorf(testCode1, "errorf %s", "aaaa")
	assert.NotNil(t, msg)
	assert.Equal(t, "ERROR-1, errorf aaaa", msg.Error())
}

func TestWrapf(t *testing.T) {
	err := Wrapf(testErr, testCode1, "Wrapf %s", "aaaa")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, testErr))
}

func TestFromError(t *testing.T) {
	err := errors.New("ERROR-1, errorf aaaa")
	code, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, testCode1, code)
	assert.Equal(t, "unit_test", code.Module())
	assert.Equal(t, "Hello, test code1", code.Desc())

	err = errors.New("not beego error")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR-2, not register")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR-aaa, invalid code")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("aaaaaaaaaaaaaa")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR-2-3, invalid error")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR, invalid error")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)
}
