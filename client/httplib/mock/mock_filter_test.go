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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/httplib"
)

func TestMockResponseFilter_FilterChain(t *testing.T) {
	req := httplib.Get("http://localhost:8080/abc/s")
	ft := NewMockResponseFilter()

	expectedResp := httplib.NewHttpResponseWithJsonBody(`{}`)
	expectedErr := errors.New("expected error")
	ft.Mock(NewSimpleCondition("/abc/s"), expectedResp, expectedErr)

	req.AddFilters(ft.FilterChain)

	resp, err := req.DoRequest()
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedResp, resp)

	req = httplib.Get("http://localhost:8080/abcd/s")
	req.AddFilters(ft.FilterChain)

	resp, err = req.DoRequest()
	assert.NotEqual(t, expectedErr, err)
	assert.NotEqual(t, expectedResp, resp)

	req = httplib.Get("http://localhost:8080/abc/s")
	req.AddFilters(ft.FilterChain)
	expectedResp1 := httplib.NewHttpResponseWithJsonBody(map[string]string{})
	expectedErr1 := errors.New("expected error")
	ft.Mock(NewSimpleCondition("/abc/abs/bbc"), expectedResp1, expectedErr1)

	resp, err = req.DoRequest()
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedResp, resp)

	req = httplib.Get("http://localhost:8080/abc/abs/bbc")
	req.AddFilters(ft.FilterChain)
	ft.Mock(NewSimpleCondition("/abc/abs/bbc"), expectedResp1, expectedErr1)
	resp, err = req.DoRequest()
	assert.Equal(t, expectedErr1, err)
	assert.Equal(t, expectedResp1, resp)
}
