// Copyright 2023 beego. All Rights Reserved.
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

package errs

import (
	"errors"
	"fmt"
)

var ErrGetByMd = errors.New("orm: Unknown error in get model")

// NewErrUnknownField returns an error representing an unknown field
// Generally, it means that you may have entered a column name or an incorrect field name
func NewErrUnknownField(fd string) error {
	return fmt.Errorf("orm: Unknown field %s", fd)
}

// NewErrUnsupportedExpressionType returns an error message that does not support the expression
func NewErrUnsupportedExpressionType(exp any) error {
	return fmt.Errorf("orm: Unsupported expression %v", exp)
}
