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

import (
	"net/http"
	"net/url"
	"time"
)

type ClientOption func(client *Client) error

// client设置
func WithTimeout(connectTimeout, readWriteTimeout time.Duration) ClientOption
func WithEnableCookie(enable bool) ClientOption
func WithUserAgent(userAgent string) ClientOption
func WithCookie(cookie *http.Cookie) ClientOption
func WithTransport(transport http.RoundTripper) ClientOption
func WithProxy(proxy func(*http.Request) (*url.URL, error)) ClientOption
func WithCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) ClientOption
func WithAccept(accept string) ClientOption
