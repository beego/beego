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
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// code, msg
const errFmt = "ERROR-%d, %s"

// Err returns an error representing c and msg.  If c is OK, returns nil.
func Error(c Code, msg string) error {
	return fmt.Errorf(errFmt, c.Code(), msg)
}

// Errorf returns error
func Errorf(c Code, format string, a ...interface{}) error {
	return Error(c, fmt.Sprintf(format, a...))
}

func Wrap(err error, c Code, msg string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, fmt.Sprintf(errFmt, c.Code(), msg))
}

func Wrapf(err error, c Code, format string, a ...interface{}) error {
	return Wrap(err, c, fmt.Sprintf(format, a...))
}

// FromError is very simple. It just parse error msg and check whether code has been register
// if code not being register, return unknown
// if err.Error() is not valid beego error code, return unknown
func FromError(err error) (Code, bool) {
	msg := err.Error()
	codeSeg := strings.SplitN(msg, ",", 2)
	if strings.HasPrefix(codeSeg[0], "ERROR-") {
		codeStr := strings.SplitN(codeSeg[0], "-", 2)
		if len(codeStr) < 2 {
			return Unknown, false
		}
		codeInt, e := strconv.ParseUint(codeStr[1], 10, 32)
		if e != nil {
			return Unknown, false
		}
		if code, ok := defaultCodeRegistry.Get(uint32(codeInt)); ok {
			return code, true
		}
	}
	return Unknown, false
}
