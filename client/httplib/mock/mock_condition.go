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

package mock

import (
	"context"
	"encoding/json"
	"net/textproto"
	"regexp"

	"github.com/beego/beego/v2/client/httplib"
)

type RequestCondition interface {
	Match(ctx context.Context, req *httplib.BeegoHTTPRequest) bool
}

// reqCondition create condition
// - path: same path
// - pathReg: request path match pathReg
// - method: same method
// - Query parameters (key, value)
// - header (key, value)
// - Body json format, contains specific (key, value).
type SimpleCondition struct {
	pathReg string
	path    string
	method  string
	query   map[string]string
	header  map[string]string
	body    map[string]interface{}
}

func NewSimpleCondition(path string, opts ...simpleConditionOption) *SimpleCondition {
	sc := &SimpleCondition{
		path:   path,
		query:  make(map[string]string),
		header: make(map[string]string),
		body:   map[string]interface{}{},
	}
	for _, o := range opts {
		o(sc)
	}
	return sc
}

func (sc *SimpleCondition) Match(ctx context.Context, req *httplib.BeegoHTTPRequest) bool {
	res := true
	if len(sc.path) > 0 {
		res = sc.matchPath(ctx, req)
	} else if len(sc.pathReg) > 0 {
		res = sc.matchPathReg(ctx, req)
	} else {
		return false
	}
	return res &&
		sc.matchMethod(ctx, req) &&
		sc.matchQuery(ctx, req) &&
		sc.matchHeader(ctx, req) &&
		sc.matchBodyFields(ctx, req)
}

func (sc *SimpleCondition) matchPath(ctx context.Context, req *httplib.BeegoHTTPRequest) bool {
	path := req.GetRequest().URL.Path
	return path == sc.path
}

func (sc *SimpleCondition) matchPathReg(ctx context.Context, req *httplib.BeegoHTTPRequest) bool {
	path := req.GetRequest().URL.Path
	if b, err := regexp.Match(sc.pathReg, []byte(path)); err == nil {
		return b
	}
	return false
}

func (sc *SimpleCondition) matchQuery(ctx context.Context, req *httplib.BeegoHTTPRequest) bool {
	qs := req.GetRequest().URL.Query()
	for k, v := range sc.query {
		if uv, ok := qs[k]; !ok || uv[0] != v {
			return false
		}
	}
	return true
}

func (sc *SimpleCondition) matchHeader(ctx context.Context, req *httplib.BeegoHTTPRequest) bool {
	headers := req.GetRequest().Header
	for k, v := range sc.header {
		if uv, ok := headers[k]; !ok || uv[0] != v {
			return false
		}
	}
	return true
}

func (sc *SimpleCondition) matchBodyFields(ctx context.Context, req *httplib.BeegoHTTPRequest) bool {
	if len(sc.body) == 0 {
		return true
	}
	getBody := req.GetRequest().GetBody
	body, err := getBody()
	if err != nil {
		return false
	}
	bytes := make([]byte, req.GetRequest().ContentLength)
	_, err = body.Read(bytes)
	if err != nil {
		return false
	}

	m := make(map[string]interface{})

	err = json.Unmarshal(bytes, &m)

	if err != nil {
		return false
	}

	for k, v := range sc.body {
		if uv, ok := m[k]; !ok || uv != v {
			return false
		}
	}
	return true
}

func (sc *SimpleCondition) matchMethod(ctx context.Context, req *httplib.BeegoHTTPRequest) bool {
	if len(sc.method) > 0 {
		return sc.method == req.GetRequest().Method
	}
	return true
}

type simpleConditionOption func(sc *SimpleCondition)

func WithPathReg(pathReg string) simpleConditionOption {
	return func(sc *SimpleCondition) {
		sc.pathReg = pathReg
	}
}

func WithQuery(key, value string) simpleConditionOption {
	return func(sc *SimpleCondition) {
		sc.query[key] = value
	}
}

func WithHeader(key, value string) simpleConditionOption {
	return func(sc *SimpleCondition) {
		sc.header[textproto.CanonicalMIMEHeaderKey(key)] = value
	}
}

func WithJsonBodyFields(field string, value interface{}) simpleConditionOption {
	return func(sc *SimpleCondition) {
		sc.body[field] = value
	}
}

func WithMethod(method string) simpleConditionOption {
	return func(sc *SimpleCondition) {
		sc.method = method
	}
}
