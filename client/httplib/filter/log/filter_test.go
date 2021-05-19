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

package log

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/httplib"
)

func TestFilterChain(t *testing.T) {
	next := func(ctx context.Context, req *httplib.BeegoHTTPRequest) (*http.Response, error) {
		time.Sleep(100 * time.Millisecond)
		return &http.Response{
			StatusCode: 404,
		}, nil
	}
	builder := NewFilterChainBuilder()
	filter := builder.FilterChain(next)
	req := httplib.Get("https://github.com/notifications?query=repo%3Aastaxie%2Fbeego")
	resp, err := filter(context.Background(), req)
	defer resp.Body.Close()
	assert.NotNil(t, resp)
	assert.Nil(t, err)
}

func TestContains(t *testing.T) {
	jsonType := "application/json"
	cases := []struct {
		Name        string
		Types       []string
		ContentType string
		Expected    bool
	}{
		{"case1", []string{jsonType}, jsonType, true},
		{"case2", []string{"text/plain"}, jsonType, false},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := contains(c.Types, c.ContentType); ans != c.Expected {
				t.Fatalf("Types: %v, ContentType: %v, expected %v, but %v got",
					c.Types, c.ContentType, c.Expected, ans)
			}
		})
	}
}
