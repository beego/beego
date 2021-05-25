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
	"io"
	"io/ioutil"
	"net/http"
)

// Client provides an HTTP client supporting chain call
type Client struct {
	Name       string
	Endpoint   string
	CommonOpts []BeegoHTTPRequestOption

	Setting BeegoHTTPSettings
}

// HTTPResponseCarrier If value implement HTTPResponseCarrier. http.Response will pass to SetHTTPResponse
type HTTPResponseCarrier interface {
	SetHTTPResponse(resp *http.Response)
}

// HTTPBodyCarrier If value implement HTTPBodyCarrier. http.Response.Body will pass to SetReader
type HTTPBodyCarrier interface {
	SetReader(r io.ReadCloser)
}

// HTTPBytesCarrier If value implement HTTPBytesCarrier.
// All the byte in http.Response.Body will pass to SetBytes
type HTTPBytesCarrier interface {
	SetBytes(bytes []byte)
}

// HTTPStatusCarrier If value implement HTTPStatusCarrier. http.Response.StatusCode will pass to SetStatusCode
type HTTPStatusCarrier interface {
	SetStatusCode(status int)
}

// HttpHeaderCarrier If value implement HttpHeaderCarrier. http.Response.Header will pass to SetHeader
type HTTPHeadersCarrier interface {
	SetHeader(header map[string][]string)
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

func (c *Client) customReq(req *BeegoHTTPRequest, opts []BeegoHTTPRequestOption) {
	req.Setting(c.Setting)
	opts = append(c.CommonOpts, opts...)
	for _, o := range opts {
		o(req)
	}
}

// handleResponse try to parse body to meaningful value
func (c *Client) handleResponse(value interface{}, req *BeegoHTTPRequest) error {
	err := c.handleCarrier(value, req)
	if err != nil {
		return err
	}

	return req.ToValue(value)
}

// handleCarrier set http data to value
func (c *Client) handleCarrier(value interface{}, req *BeegoHTTPRequest) error {
	resp, err := req.Response()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if value == nil {
		return err
	}
	if carrier, ok := value.(HTTPResponseCarrier); ok {
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		resp.Body = ioutil.NopCloser(bytes.NewReader(b))
		carrier.SetHTTPResponse(resp)
	}
	if carrier, ok := value.(HTTPBodyCarrier); ok {
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		reader := ioutil.NopCloser(bytes.NewReader(b))
		carrier.SetReader(reader)
	}
	if carrier, ok := value.(HTTPBytesCarrier); ok {
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		carrier.SetBytes(b)
	}
	if carrier, ok := value.(HTTPStatusCarrier); ok {
		carrier.SetStatusCode(resp.StatusCode)
	}
	if carrier, ok := value.(HTTPHeadersCarrier); ok {
		carrier.SetHeader(resp.Header)
	}
	return nil
}

// Get Send a GET request and try to give its result value
func (c *Client) Get(value interface{}, path string, opts ...BeegoHTTPRequestOption) error {
	req := Get(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}

// Post Send a POST request and try to give its result value
func (c *Client) Post(value interface{}, path string, body interface{}, opts ...BeegoHTTPRequestOption) error {
	req := Post(c.Endpoint + path)
	c.customReq(req, opts)
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Put Send a Put request and try to give its result value
func (c *Client) Put(value interface{}, path string, body interface{}, opts ...BeegoHTTPRequestOption) error {
	req := Put(c.Endpoint + path)
	c.customReq(req, opts)
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Delete Send a Delete request and try to give its result value
func (c *Client) Delete(value interface{}, path string, opts ...BeegoHTTPRequestOption) error {
	req := Delete(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}

// Head Send a Head request and try to give its result value
func (c *Client) Head(value interface{}, path string, opts ...BeegoHTTPRequestOption) error {
	req := Head(c.Endpoint + path)
	c.customReq(req, opts)
	return c.handleResponse(value, req)
}
