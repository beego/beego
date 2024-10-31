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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// NewHttpResponseWithJsonBody will try to convert the data to json format
// usually you only use this when you want to mock http Response
func NewHttpResponseWithJsonBody(data interface{}) *http.Response {
	var body []byte
	if str, ok := data.(string); ok {
		body = []byte(str)
	} else if bts, ok := data.([]byte); ok {
		body = bts
	} else {
		body, _ = json.Marshal(data)
	}
	return &http.Response{
		ContentLength: int64(len(body)),
		Body:          io.NopCloser(bytes.NewReader(body)),
	}
}
