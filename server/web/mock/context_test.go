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

package mock

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/server/web"
)

type TestController struct {
	web.Controller
}

func TestMockContext(t *testing.T) {
	req, err := http.NewRequest("GET", "https://localhost:8080/hello?name=tom", bytes.NewReader([]byte{}))
	assert.Nil(t, err)
	ctx, resp := NewMockContext(req)
	ctrl := &TestController{
		Controller: web.Controller{
			Ctx: ctx,
		},
	}
	ctrl.HelloWorld()
	result := resp.Body.String()
	assert.Equal(t, "name=tom", result)
}

// GET hello?name=XXX
func (c *TestController) HelloWorld() {
	name := c.GetString("name")
	c.Ctx.WriteString(fmt.Sprintf("name=%s", name))
}

func (c *TestController) HelloSession() {
	err := c.SessionRegenerateID()
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	_ = c.SetSession("name", "Tom")
	c.Ctx.WriteString("set")
}

func (c *TestController) HelloSessionName() {
	name := c.CruSession.Get(context.Background(), "name")
	c.Ctx.WriteString(name.(string))
}
