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
	"net/http"
)

// Client provides an HTTP client supporting chain call
type Client struct {
	Name       string
	Endpoint   string
	CommonOpts []BeegoHttpRequestOption

	Setting BeegoHTTPSettings
	pointer responsePointer
}

type responsePointer struct {
	response      **http.Response
	statusCode    **int
	header        **http.Header
	headerValues  map[string]**string //用户传一个key，然后将key存在map的key里，header的value存在value里
	contentLength **int64
}

// NewClient return a new http client
func NewClient(name string, endpoint string, opts ...ClientOption) (*Client, error) {
	res := &Client{
		Name:     name,
		Endpoint: endpoint,
	}
	setting := GetDefaultSetting()
	res.Setting = setting
	for _, o := range opts {
		o(res)
	}
	return res, nil
}

// Response will set response to the pointer
func (c *Client) Response(resp **http.Response) *Client {
	newC := *c
	newC.pointer.response = resp
	return &newC
}

// StatusCode will set response StatusCode to the pointer
func (c *Client) StatusCode(code **int) *Client {
	newC := *c
	newC.pointer.statusCode = code
	return &newC
}

// Headers will set response Headers to the pointer
func (c *Client) Headers(headers **http.Header) *Client {
	newC := *c
	newC.pointer.header = headers
	return &newC
}

// HeaderValue will set response HeaderValue to the pointer
func (c *Client) HeaderValue(key string, value **string) *Client {
	newC := *c
	if newC.pointer.headerValues == nil {
		newC.pointer.headerValues = make(map[string]**string)
	}
	newC.pointer.headerValues[key] = value
	return &newC
}

// ContentType will set response ContentType to the pointer
func (c *Client) ContentType(contentType **string) *Client {
	return c.HeaderValue("Content-Type", contentType)
}

// ContentLength will set response ContentLength to the pointer
func (c *Client) ContentLength(contentLength **int64) *Client {
	newC := *c
	newC.pointer.contentLength = contentLength
	return &newC
}

// setPointers set the http response value to pointer
func (c *Client) setPointers(resp *http.Response) {
	if c.pointer.response != nil {
		*c.pointer.response = resp
	}
	if c.pointer.statusCode != nil {
		*c.pointer.statusCode = &resp.StatusCode
	}
	if c.pointer.header != nil {
		*c.pointer.header = &resp.Header
	}
	if c.pointer.headerValues != nil {
		for k, v := range c.pointer.headerValues {
			s := resp.Header.Get(k)
			*v = &s
		}
	}
	if c.pointer.contentLength != nil {
		*c.pointer.contentLength = &resp.ContentLength
	}
}

func (c *Client) customReq(req *BeegoHTTPRequest, opts []BeegoHttpRequestOption) {
	req.Setting(c.Setting)
	opts = append(c.CommonOpts, opts...)
	for _, o := range opts {
		o(req)
	}
}

// handleResponse try to parse body to meaningful value
func (c *Client) handleResponse(value interface{}, req *BeegoHTTPRequest) error {
	// send request
	resp, err := req.Response()
	if err != nil {
		return err
	}
	c.setPointers(resp)
	return req.ResponseForValue(value)
}

// Get Send a GET request and try to give its result value
func (c *Client) Get(value interface{}, path string, opts ...BeegoHttpRequestOption) error {
	req := Get(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}

// Post Send a POST request and try to give its result value
func (c *Client) Post(value interface{}, path string, body interface{}, opts ...BeegoHttpRequestOption) error {
	req := Post(c.Endpoint + path)
	c.customReq(req, opts)
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Put Send a Put request and try to give its result value
func (c *Client) Put(value interface{}, path string, body interface{}, opts ...BeegoHttpRequestOption) error {
	req := Put(c.Endpoint + path)
	c.customReq(req, opts)
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Delete Send a Delete request and try to give its result value
func (c *Client) Delete(value interface{}, path string, opts ...BeegoHttpRequestOption) error {
	req := Delete(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}

// Head Send a Head request and try to give its result value
func (c *Client) Head(value interface{}, path string, opts ...BeegoHttpRequestOption) error {
	req := Head(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}
