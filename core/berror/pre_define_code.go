// Copyright 2021 beego
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
)

// pre define code

// Unknown indicates got some error which is not defined
var Unknown = DefineCode(5000001, "error", "Unknown", fmt.Sprintf(`
Unknown error code. Usually you will see this code in three cases:
1. You forget to define Code or function DefineCode not being executed;
2. This is not Beego's error but you call FromError();
3. Beego got unexpected error and don't know how to handle it, and then return Unknown error

A common practice to DefineCode looks like:
%s

In this way, you may forget to import this package, and got Unknown error. 

Sometimes, you believe you got Beego error, but actually you don't, and then you call FromError(err)

`, goCodeBlock(`
import your_package

func init() {
    DefineCode(5100100, "your_module", "detail")
    // ...
}
`)))

func goCodeBlock(code string) string {
	return codeBlock("go", code)
}
func codeBlock(lan string, code string) string {
	return fmt.Sprintf("```%s\n%s\n```", lan, code)
}
