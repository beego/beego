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
	"strings"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/berror"
)

type Client struct {
	Name       string
	Endpoint   string
	CommonOpts []BeegoHttpRequestOption

	Setting *httplib.BeegoHTTPSettings
	pointer *ResponsePointer
}

type ResponsePointer struct {
	response      **http.Response
	statusCode    **int
	header        **http.Header
	headerValues  map[string]**string //用户传一个key，然后将key存在map的key里，header的value存在value里
	contentLength **int64
}

// NewClient
func NewClient(name string, endpoint string, opts ...ClientOption) (*Client, error) {
	res := &Client{
		Name:     name,
		Endpoint: endpoint,
	}
	for _, o := range opts {
		err := o(res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// Response will set response to the pointer
func (c *Client) Response(resp **http.Response) *Client {
	if c.pointer == nil {
		newC := *c
		newC.pointer = &ResponsePointer{
			response: resp,
		}
		return &newC
	}
	c.pointer.response = resp
	return c
}

// StatusCode will set response StatusCode to the pointer
func (c *Client) StatusCode(code **int) *Client {
	if c.pointer == nil {
		newC := *c
		newC.pointer = &ResponsePointer{
			statusCode: code,
		}
		return &newC
	}
	c.pointer.statusCode = code
	return c
}

// Headers will set response Headers to the pointer
func (c *Client) Headers(headers **http.Header) *Client {
	if c.pointer == nil {
		newC := *c
		newC.pointer = &ResponsePointer{
			header: headers,
		}
		return &newC
	}
	c.pointer.header = headers
	return c
}

// HeaderValue will set response HeaderValue to the pointer
func (c *Client) HeaderValue(key string, value **string) *Client {
	if c.pointer == nil {
		newC := *c
		newC.pointer = &ResponsePointer{
			headerValues: map[string]**string{
				key: value,
			},
		}
		return &newC
	}
	if c.pointer.headerValues == nil {
		c.pointer.headerValues = map[string]**string{}
	}
	c.pointer.headerValues[key] = value
	return c
}

// ContentType will set response ContentType to the pointer
func (c *Client) ContentType(contentType **string) *Client {
	return c.HeaderValue("Content-Type", contentType)
}

// ContentLength will set response ContentLength to the pointer
func (c *Client) ContentLength(contentLength **int64) *Client {
	if c.pointer == nil {
		newC := *c
		newC.pointer = &ResponsePointer{
			contentLength: contentLength,
		}
		return &newC
	}
	c.pointer.contentLength = contentLength
	return c
}

// setPointers set the http response value to pointer
func (c *Client) setPointers(resp *http.Response) {
	if c.pointer == nil {
		return
	}
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

// initRequest will apply all the client setting, common option and request option
func (c *Client) newRequest(method, path string, opts []BeegoHttpRequestOption) (*httplib.BeegoHTTPRequest, error) {
	var req *httplib.BeegoHTTPRequest
	switch method {
	case http.MethodGet:
		req = httplib.Get(c.Endpoint + path)
	case http.MethodPost:
		req = httplib.Post(c.Endpoint + path)
	case http.MethodPut:
		req = httplib.Put(c.Endpoint + path)
	case http.MethodDelete:
		req = httplib.Delete(c.Endpoint + path)
	case http.MethodHead:
		req = httplib.Head(c.Endpoint + path)
	}

	req = req.Setting(*c.Setting)
	for _, o := range c.CommonOpts {
		err := o(req)
		if err != nil {
			return nil, err
		}
	}
	for _, o := range opts {
		err := o(req)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

// handleResponse try to parse body to meaningful value
func (c *Client) handleResponse(value interface{}, req *httplib.BeegoHTTPRequest) error {
	// send request
	resp, err := req.Response()
	if err != nil {
		return err
	}
	c.setPointers(resp)

	// handle basic type
	switch v := value.(type) {
	case **string:
		s, err := req.String()
		if err != nil {
			return nil
		}
		*v = &s
		return nil
	case **[]byte:
		bs, err := req.Bytes()
		if err != nil {
			return nil
		}
		*v = &bs
		return nil
	}

	// try to parse it as content type
	switch strings.Split(resp.Header.Get("Content-Type"), ";")[0] {
	case "application/json":
		return req.ToJSON(value)
	case "text/xml":
		return req.ToXML(value)
	case "text/yaml", "application/x-yaml":
		return req.ToYAML(value)
	}

	// try to parse it anyway
	if err := req.ToJSON(value); err == nil {
		return nil
	}
	if err := req.ToYAML(value); err == nil {
		return nil
	}
	if err := req.ToXML(value); err == nil {
		return nil
	}

	// TODO add new error type about can't parse body
	return berror.Error(httplib.UnsupportedBodyType, "unsupported body data")
}

// Get Send a GET request and try to give its result value
func (c *Client) Get(value interface{}, path string, opts ...BeegoHttpRequestOption) error {
	req, err := c.newRequest(http.MethodGet, path, opts)
	if err != nil {
		return err
	}
	return c.handleResponse(value, req)
}

// Post Send a POST request and try to give its result value
func (c *Client) Post(value interface{}, path string, body interface{}, opts ...BeegoHttpRequestOption) error {
	req, err := c.newRequest(http.MethodPost, path, opts)
	if err != nil {
		return err
	}
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Put Send a Put request and try to give its result value
func (c *Client) Put(value interface{}, path string, body interface{}, opts ...BeegoHttpRequestOption) error {
	req, err := c.newRequest(http.MethodPut, path, opts)
	if err != nil {
		return err
	}
	if body != nil {
		req = req.Body(body)
	}
	return c.handleResponse(value, req)
}

// Delete Send a Delete request and try to give its result value
func (c *Client) Delete(value interface{}, path string, opts ...BeegoHttpRequestOption) error {
	req, err := c.newRequest(http.MethodDelete, path, opts)
	if err != nil {
		return err
	}
	return c.handleResponse(value, req)
}

// Head Send a Head request and try to give its result value
func (c *Client) Head(value interface{}, path string, opts ...BeegoHttpRequestOption) error {
	req, err := c.newRequest(http.MethodHead, path, opts)
	if err != nil {
		return err
	}
	return c.handleResponse(value, req)
}
