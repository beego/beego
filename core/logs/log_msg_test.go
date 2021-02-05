// Copyright 2020
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

package logs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogMsg_OldStyleFormat(t *testing.T) {
	lg := &LogMsg{
		Level:      LevelDebug,
		Msg:        "Hello, world",
		When:       time.Date(2020, 9, 19, 20, 12, 37, 9, time.UTC),
		FilePath:   "/user/home/main.go",
		LineNumber: 13,
		Prefix:     "Cus",
	}
	res := lg.OldStyleFormat()
	assert.Equal(t, "[D] Cus Hello, world", res)

	lg.enableFuncCallDepth = true
	res = lg.OldStyleFormat()
	assert.Equal(t, "[D] [main.go:13] Cus Hello, world", res)

	lg.enableFullFilePath = true

	res = lg.OldStyleFormat()
	assert.Equal(t, "[D] [/user/home/main.go:13] Cus Hello, world", res)

	lg.Msg = "hello, %s"
	lg.Args = []interface{}{"world"}
	assert.Equal(t, "[D] [/user/home/main.go:13] Cus hello, world", lg.OldStyleFormat())
}
