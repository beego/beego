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

package mock

import (
	"context"
	"net/http"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
)

const mockCtxKey = "beego-httplib-mock"

func init() {
	InitMockSetting()
}

type Stub interface {
	Mock(cond RequestCondition, resp *http.Response, err error)
	Clear()
	MockByPath(path string, resp *http.Response, err error)
}

var mockFilter = &MockResponseFilter{}

func InitMockSetting() {
	httplib.AddDefaultFilter(mockFilter.FilterChain)
}

func StartMock() Stub {
	return mockFilter
}

func CtxWithMock(ctx context.Context, mock... *Mock) context.Context {
	return context.WithValue(ctx, mockCtxKey, mock)
}

func mockFromCtx(ctx context.Context) []*Mock {
	ms := ctx.Value(mockCtxKey)
	if ms != nil {
		if res, ok := ms.([]*Mock); ok {
			return res
		}
		logs.Error("mockCtxKey found in context, but value is not type []*Mock")
	}
	return nil
}

type Mock struct {
	cond RequestCondition
	resp *http.Response
	err  error
}

func NewMockByPath(path string, resp *http.Response, err error) *Mock {
	return NewMock(NewSimpleCondition(path), resp, err)
}

func NewMock(con RequestCondition, resp *http.Response, err error) *Mock {
	return &Mock{
		cond: con,
		resp: resp,
		err: err,
	}
}


