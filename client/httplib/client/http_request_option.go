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

package client

import "github.com/beego/beego/v2/client/httplib"

type BeegoHttpRequestOption func(request *httplib.BeegoHTTPRequest) error

// Req设置
func WithTokenFactory(tokenFactory func() (string, error)) BeegoHttpRequestOption
func WithBasicAuth(basicAuth func() (string, error)) BeegoHttpRequestOption
func WithFilters(fcs ...httplib.FilterChain) BeegoHttpRequestOption
func WithContentType(contentType string) BeegoHttpRequestOption
