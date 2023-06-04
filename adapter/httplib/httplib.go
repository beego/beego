// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package httplib is used as http.Client
// Usage:
//
// import "github.com/beego/beego/v2/client/httplib"
//
//	b := httplib.Post("http://beego.vip/")
//	b.Param("username","astaxie")
//	b.Param("password","123456")
//	b.PostFile("uploadfile1", "httplib.pdf")
//	b.PostFile("uploadfile2", "httplib.txt")
//	str, err := b.String()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(str)
package httplib

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/beego/beego/v2/client/httplib"
)

// SetDefaultSetting Overwrite default settings
func SetDefaultSetting(setting BeegoHTTPSettings) {
	httplib.SetDefaultSetting(httplib.BeegoHTTPSettings(setting))
}

// NewBeegoRequest return *BeegoHttpRequest with specific method
func NewBeegoRequest(rawurl, method string) *BeegoHTTPRequest {
	return &BeegoHTTPRequest{
		delegate: httplib.NewBeegoRequest(rawurl, method),
	}
}

// Get returns *BeegoHttpRequest with GET method.
func Get(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "GET")
}

// Post returns *BeegoHttpRequest with POST method.
func Post(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "POST")
}

// Put returns *BeegoHttpRequest with PUT method.
func Put(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "PUT")
}

// Delete returns *BeegoHttpRequest DELETE method.
func Delete(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "DELETE")
}

// Head returns *BeegoHttpRequest with HEAD method.
func Head(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "HEAD")
}

// BeegoHTTPSettings is the http.Client setting
type BeegoHTTPSettings httplib.BeegoHTTPSettings

// BeegoHTTPRequest provides more useful methods for requesting one url than http.Request.
type BeegoHTTPRequest struct {
	delegate *httplib.BeegoHTTPRequest
}

// GetRequest return the request object
func (b *BeegoHTTPRequest) GetRequest() *http.Request {
	return b.delegate.GetRequest()
}

// Setting Change request settings
func (b *BeegoHTTPRequest) Setting(setting BeegoHTTPSettings) *BeegoHTTPRequest {
	b.delegate.Setting(httplib.BeegoHTTPSettings(setting))
	return b
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (b *BeegoHTTPRequest) SetBasicAuth(username, password string) *BeegoHTTPRequest {
	b.delegate.SetBasicAuth(username, password)
	return b
}

// SetEnableCookie sets enable/disable cookiejar
func (b *BeegoHTTPRequest) SetEnableCookie(enable bool) *BeegoHTTPRequest {
	b.delegate.SetEnableCookie(enable)
	return b
}

// SetUserAgent sets User-Agent header field
func (b *BeegoHTTPRequest) SetUserAgent(useragent string) *BeegoHTTPRequest {
	b.delegate.SetUserAgent(useragent)
	return b
}

// Retries sets Retries times.
// default is 0 means no retried.
// -1 means retried forever.
// others means retried times.
func (b *BeegoHTTPRequest) Retries(times int) *BeegoHTTPRequest {
	b.delegate.Retries(times)
	return b
}

func (b *BeegoHTTPRequest) RetryDelay(delay time.Duration) *BeegoHTTPRequest {
	b.delegate.RetryDelay(delay)
	return b
}

// SetTimeout sets connect time out and read-write time out for BeegoRequest.
func (b *BeegoHTTPRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *BeegoHTTPRequest {
	b.delegate.SetTimeout(connectTimeout, readWriteTimeout)
	return b
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (b *BeegoHTTPRequest) SetTLSClientConfig(config *tls.Config) *BeegoHTTPRequest {
	b.delegate.SetTLSClientConfig(config)
	return b
}

// Header add header item string in request.
func (b *BeegoHTTPRequest) Header(key, value string) *BeegoHTTPRequest {
	b.delegate.Header(key, value)
	return b
}

// SetHost set the request host
func (b *BeegoHTTPRequest) SetHost(host string) *BeegoHTTPRequest {
	b.delegate.SetHost(host)
	return b
}

// SetProtocolVersion Set the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (b *BeegoHTTPRequest) SetProtocolVersion(vers string) *BeegoHTTPRequest {
	b.delegate.SetProtocolVersion(vers)
	return b
}

// SetCookie add cookie into request.
func (b *BeegoHTTPRequest) SetCookie(cookie *http.Cookie) *BeegoHTTPRequest {
	b.delegate.SetCookie(cookie)
	return b
}

// SetTransport set the setting transport
func (b *BeegoHTTPRequest) SetTransport(transport http.RoundTripper) *BeegoHTTPRequest {
	b.delegate.SetTransport(transport)
	return b
}

// SetProxy set the http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
//		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
//		return u, nil
//	}
func (b *BeegoHTTPRequest) SetProxy(proxy func(*http.Request) (*url.URL, error)) *BeegoHTTPRequest {
	b.delegate.SetProxy(proxy)
	return b
}

// SetCheckRedirect specifies the policy for handling redirects.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func (b *BeegoHTTPRequest) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *BeegoHTTPRequest {
	b.delegate.SetCheckRedirect(redirect)
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *BeegoHTTPRequest) Param(key, value string) *BeegoHTTPRequest {
	b.delegate.Param(key, value)
	return b
}

// PostFile add a post file to the request
func (b *BeegoHTTPRequest) PostFile(formname, filename string) *BeegoHTTPRequest {
	b.delegate.PostFile(formname, filename)
	return b
}

// Body adds request raw body.
// it supports string and []byte.
func (b *BeegoHTTPRequest) Body(data interface{}) *BeegoHTTPRequest {
	b.delegate.Body(data)
	return b
}

// XMLBody adds request raw body encoding by XML.
func (b *BeegoHTTPRequest) XMLBody(obj interface{}) (*BeegoHTTPRequest, error) {
	_, err := b.delegate.XMLBody(obj)
	return b, err
}

// YAMLBody adds request raw body encoding by YAML.
func (b *BeegoHTTPRequest) YAMLBody(obj interface{}) (*BeegoHTTPRequest, error) {
	_, err := b.delegate.YAMLBody(obj)
	return b, err
}

// JSONBody adds request raw body encoding by JSON.
func (b *BeegoHTTPRequest) JSONBody(obj interface{}) (*BeegoHTTPRequest, error) {
	_, err := b.delegate.JSONBody(obj)
	return b, err
}

// DoRequest will do the client.Do
func (b *BeegoHTTPRequest) DoRequest() (resp *http.Response, err error) {
	return b.delegate.DoRequest()
}

// String returns the body string in response.
// it calls Response inner.
func (b *BeegoHTTPRequest) String() (string, error) {
	return b.delegate.String()
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (b *BeegoHTTPRequest) Bytes() ([]byte, error) {
	return b.delegate.Bytes()
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (b *BeegoHTTPRequest) ToFile(filename string) error {
	return b.delegate.ToFile(filename)
}

// ToJSON returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *BeegoHTTPRequest) ToJSON(v interface{}) error {
	return b.delegate.ToJSON(v)
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *BeegoHTTPRequest) ToXML(v interface{}) error {
	return b.delegate.ToXML(v)
}

// ToYAML returns the map that marshals from the body bytes as yaml in response .
// it calls Response inner.
func (b *BeegoHTTPRequest) ToYAML(v interface{}) error {
	return b.delegate.ToYAML(v)
}

// Response executes request client gets response mannually.
func (b *BeegoHTTPRequest) Response() (*http.Response, error) {
	return b.delegate.Response()
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return httplib.TimeoutDialer(cTimeout, rwTimeout)
}
