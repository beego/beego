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
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestHttpHandler struct {
}

func (t *TestHttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("coming"))
}

func TestJLWriter_WriteMsg(t *testing.T) {
	// start sever

	http.Handle("/", &TestHttpHandler{})
	go http.ListenAndServe(":12124", nil)

	jl := newJLWriter()
	jl.Init(`{
"webhookurl":"http://localhost:12124/hello",
"redirecturl":"nil",
"imageurl":"a"
}`)
	err := jl.WriteMsg(&LogMsg{
		Msg: "world",
	})

	jl.Flush()
	jl.Destroy()
	assert.Nil(t, err)
}

func TestJLWriter_Format(t *testing.T) {
	lg := &LogMsg{
		Level:      LevelDebug,
		Msg:        "Hello, world",
		When:       time.Date(2020, 9, 19, 20, 12, 37, 9, time.UTC),
		FilePath:   "/user/home/main.go",
		LineNumber: 13,
		Prefix:     "Cus",
	}
	jl := newJLWriter().(*JLWriter)
	res := jl.Format(lg)
	assert.Equal(t, "2020-09-19 20:12:37 [D] Cus Hello, world", res)
}
