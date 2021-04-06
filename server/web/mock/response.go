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
	"encoding/json"
	"net/http"
)

// HttpResponse mock response, which should be used in tests
type HttpResponse struct {
	body       []byte
	header     http.Header
	StatusCode int
}

// NewMockHttpResponse you should only use this in your test code
func NewMockHttpResponse() *HttpResponse {
	return &HttpResponse{
		body: make([]byte, 0),
		header: make(http.Header),
	}
}

// Header return headers
func (m *HttpResponse) Header() http.Header {
	return m.header
}

// Write append the body
func (m *HttpResponse) Write(bytes []byte) (int, error) {
	m.body = append(m.body, bytes...)
	return len(bytes), nil
}

// WriteHeader set the status code
func (m *HttpResponse) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

// JsonUnmarshal convert the body to object
func (m *HttpResponse) JsonUnmarshal(value interface{}) error {
	return json.Unmarshal(m.body, value)
}

// BodyToString return the body as the string
func (m *HttpResponse) BodyToString() string {
	return string(m.body)
}
