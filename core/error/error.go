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

package error

import (
	"fmt"

	"github.com/pkg/errors"
)

type Code int32

func (c Code) ToInt32() int32 {
	return int32(c)
}


type Error struct {
	Code Code
	Msg string
	Cause error
}

func (be *Error) String() string {
	return fmt.Sprintf("code: %d, msg: %s", be.Code.ToInt32(), be.Msg)
}

func New(code Code, msg string) *Error {
	return &Error{
		Code: code,
		Msg: msg,
	}
}

func Wrap(cause error, code Code, msg string) {
	errors.Wrap()
}

func Convert(err error) *Error {

}

