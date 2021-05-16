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
}

// If value implement this interface. http.response will saved by SetHttpResponse
type HttpResponseCarrier interface {
	SetHttpResponse(resp *http.Response)
}

// If value implement this interface. bytes of http.response will saved by SetHttpResponse
type ResponseBytesCarrier interface {
	// Cause of when user get http.response, the body stream is closed. So need to pass bytes by
	SetBytes(bytes []byte)
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
	if carrier, ok := (value).(HttpResponseCarrier); ok {
		(carrier).SetHttpResponse(resp)
	}
	if carrier, ok := (value).(ResponseBytesCarrier); ok {
		bytes, err := req.Bytes()
		if err != nil {
			return err
		}
		(carrier).SetBytes(bytes)
	}

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
