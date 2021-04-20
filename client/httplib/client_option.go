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
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type ClientOption func(client *Client) error
type BeegoHttpRequestOption func(request *BeegoHTTPRequest) error

// WithEnableCookie will enable cookie in all subsequent request
func WithEnableCookie(enable bool) ClientOption {
	return func(client *Client) error {
		client.Setting.EnableCookie = enable
		return nil
	}
}

// WithEnableCookie will adds UA in all subsequent request
func WithUserAgent(userAgent string) ClientOption {
	return func(client *Client) error {
		client.Setting.UserAgent = userAgent
		return nil
	}
}

// WithTLSClientConfig will adds tls config in all subsequent request
func WithTLSClientConfig(config *tls.Config) ClientOption {
	return func(client *Client) error {
		client.Setting.TLSClientConfig = config
		return nil
	}
}

// WithTransport will set transport field in all subsequent request
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(client *Client) error {
		client.Setting.Transport = transport
		return nil
	}
}

// WithProxy will set http proxy field in all subsequent request
func WithProxy(proxy func(*http.Request) (*url.URL, error)) ClientOption {
	return func(client *Client) error {
		client.Setting.Proxy = proxy
		return nil
	}
}

// WithCheckRedirect will specifies the policy for handling redirects in all subsequent request
func WithCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(client *Client) error {
		client.Setting.CheckRedirect = redirect
		return nil
	}
}

// WithHTTPSetting can replace beegoHTTPSeting
func WithHTTPSetting(setting BeegoHTTPSettings) ClientOption {
	return func(client *Client) error {
		client.Setting = &setting
		return nil
	}
}

// WithEnableGzip will enable gzip in all subsequent request
func WithEnableGzip(enable bool) ClientOption {
	return func(client *Client) error {
		client.Setting.Gzip = enable
		return nil
	}
}

// BeegoHttpRequestOption

// WithTimeout sets connect time out and read-write time out for BeegoRequest.
func WithTimeout(connectTimeout, readWriteTimeout time.Duration) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		request.SetTimeout(connectTimeout, readWriteTimeout)
		return nil
	}
}

// WithHeader adds header item string in request.
func WithHeader(key, value string) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		request.Header(key, value)
		return nil
	}
}

// WithCookie adds a cookie to the request.
func WithCookie(cookie *http.Cookie) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		request.Header("Cookie", cookie.String())
		return nil
	}
}

// Withtokenfactory adds a custom function to set Authorization
func WithTokenFactory(tokenFactory func() (string, error)) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		t, err := tokenFactory()
		if err != nil {
			return err
		}
		request.Header("Authorization", t)
		return nil
	}
}

// WithBasicAuth adds a custom function to set basic auth
func WithBasicAuth(basicAuth func() (string, string, error)) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		username, password, err := basicAuth()
		if err != nil {
			return err
		}
		request.SetBasicAuth(username, password)
		return nil
	}
}

// WithFilters will use the filter as the invocation filters
func WithFilters(fcs ...FilterChain) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		request.SetFilters(fcs...)
		return nil
	}
}

// WithContentType adds ContentType in header
func WithContentType(contentType string) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		request.Header("Content-Type", contentType)
		return nil
	}
}

// WithParam adds query param in to request.
func WithParam(key, value string) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		request.Param(key, value)
		return nil
	}
}

// WithRetry set retry times and delay for the request
// default is 0 (never retry)
// -1 retry indefinitely (forever)
// Other numbers specify the exact retry amount
func WithRetry(times int, delay time.Duration) BeegoHttpRequestOption {
	return func(request *BeegoHTTPRequest) error {
		request.Retries(times)
		request.RetryDelay(delay)
		return nil
	}
}
