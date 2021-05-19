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
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/httplib"
)

func TestStartMock(t *testing.T) {
	// httplib.defaultSetting.FilterChains = []httplib.FilterChain{mockFilter.FilterChain}

	stub := StartMock()
	// defer stub.Clear()

	expectedResp := httplib.NewHttpResponseWithJsonBody([]byte(`{}`))
	expectedErr := errors.New("expected err")

	stub.Mock(NewSimpleCondition("/abc"), expectedResp, expectedErr)

	resp, err := OriginalCodeUsingHttplib()

	defer expectedResp.Body.Close()
	defer resp.Body.Close()

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedResp, resp)
}

// TestStartMock_Isolation Test StartMock that
// mock only work for this request
func TestStartMock_Isolation(t *testing.T) {
	// httplib.defaultSetting.FilterChains = []httplib.FilterChain{mockFilter.FilterChain}
	// setup global stub
	stub := StartMock()
	globalMockResp := httplib.NewHttpResponseWithJsonBody([]byte(`{}`))
	globalMockErr := errors.New("expected err")
	stub.Mock(NewSimpleCondition("/abc"), globalMockResp, globalMockErr)

	expectedResp := httplib.NewHttpResponseWithJsonBody(struct {
		A string `json:"a"`
	}{
		A: "aaa",
	})
	expectedErr := errors.New("expected err aa")
	m := NewMockByPath("/abc", expectedResp, expectedErr)
	ctx := CtxWithMock(context.Background(), m)

	resp, err := OriginnalCodeUsingHttplibPassCtx(ctx)

	defer globalMockResp.Body.Close()
	defer resp.Body.Close()

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expectedResp, resp)
}

func OriginnalCodeUsingHttplibPassCtx(ctx context.Context) (*http.Response, error) {
	return httplib.Get("http://localhost:7777/abc").DoRequestWithCtx(ctx)
}

func OriginalCodeUsingHttplib() (*http.Response, error) {
	return httplib.Get("http://localhost:7777/abc").DoRequest()
}
