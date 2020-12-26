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

package httplib

import (
	"context"
	"fmt"
	"net/http"
)

// MockResponse will return mock response if find any suitable mock data
// if you want to test your code using httplib, you need this.
type MockResponseFilter struct {
	ms []*Mock
}

func NewMockResponseFilter() *MockResponseFilter {
	return &MockResponseFilter{
		ms: make([]*Mock, 0, 1),
	}
}

func (m *MockResponseFilter) FilterChain(next Filter) Filter {
	return func(ctx context.Context, req *BeegoHTTPRequest) (*http.Response, error) {

		ms := mockFromCtx(ctx)
		ms = append(ms, m.ms...)

		fmt.Printf("url: %s, mock: %d \n", req.url, len(ms))
		for _, mock := range ms {
			if mock.cond.Match(ctx, req) {
				return mock.resp, mock.err
			}
		}
		return next(ctx, req)
	}
}

func (m *MockResponseFilter) MockByPath(path string, resp *http.Response, err error) {
	m.Mock(NewSimpleCondition(path), resp, err)
}

func (m *MockResponseFilter) Clear() {
	m.ms = make([]*Mock, 0, 1)
}

// Mock add mock data
// If the cond.Match(...) = true, the resp and err will be returned
func (m *MockResponseFilter) Mock(cond RequestCondition, resp *http.Response, err error) {
	m.ms = append(m.ms, NewMock(cond, resp, err))
}
