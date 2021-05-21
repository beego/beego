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
	CommonOpts []BeegoHttpRequestOption

	Setting BeegoHTTPSettings
}

// HttpResponseCarrier If value implement HttpResponseCarrier. http.Response will pass to SetHttpResponse
type HttpResponseCarrier interface {
	SetHttpResponse(resp *http.Response)
}

// HttpBodyCarrier If value implement HttpBodyCarrier. http.Response.Body will pass to SetReader
type HttpBodyCarrier interface {
	SetReader(r *io.ReadCloser)
}

// HttpBytesCarrier If value implement HttpBytesCarrier.
// All the byte in http.Response.Body will pass to SetBytes
type HttpBytesCarrier interface {
	SetBytes(bytes []byte)
}

// HttpStatusCarrier If value implement HttpStatusCarrier. http.Response.StatusCode will pass to SetStatusCode
type HttpStatusCarrier interface {
	SetStatusCode(status int)
}

// HttpHeaderCarrier If value implement HttpHeaderCarrier. http.Response.Header will pass to SetHeader
type HttpHeadersCarrier interface {
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

	switch carrier := value.(type) {
	case HttpResponseCarrier:
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		resp.Body = ioutil.NopCloser(bytes.NewReader(b))
		carrier.SetHttpResponse(resp)
		fallthrough
	case HttpBodyCarrier:
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		reader := ioutil.NopCloser(bytes.NewReader(b))
		carrier.SetReader(&reader)
		fallthrough
	case HttpBytesCarrier:
		b, err := req.Bytes()
		if err != nil {
			return err
		}
		carrier.SetBytes(b)
		fallthrough
	case HttpStatusCarrier:
		resp, err := req.Response()
		if err != nil {
			return err
		}
		carrier.SetStatusCode(resp.StatusCode)
		fallthrough
	case HttpHeadersCarrier:
		resp, err := req.Response()
		if err != nil {
			return err
		}
		carrier.SetHeader(resp.Header)
	}

	return req.ToValue(value)
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
