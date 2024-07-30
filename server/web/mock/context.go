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
	"net/http"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/mock/common"
)

func NewMockContext(req *http.Request) (*beegoCtx.Context, *common.HttpResponse) {
	ctx := beegoCtx.NewContext()
	resp := common.NewMockHttpResponse()
	ctx.Reset(resp, req)
	return ctx, resp
}
