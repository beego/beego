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

package opentracing

import (
	"context"
	"net/http"
	"strconv"

	logKit "github.com/go-kit/kit/log"
	opentracingKit "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/astaxie/beego/pkg/httplib"
)

type FilterChainBuilder struct {
	// CustomSpanFunc users are able to custom their span
	CustomSpanFunc func(span opentracing.Span, ctx context.Context,
		req *httplib.BeegoHTTPRequest, resp *http.Response, err error)
}

func (builder *FilterChainBuilder) FilterChain(next httplib.Filter) httplib.Filter {

	return func(ctx context.Context, req *httplib.BeegoHTTPRequest) (*http.Response, error) {

		method := req.GetRequest().Method
		host := req.GetRequest().URL.Host
		path := req.GetRequest().URL.Path

		proto := req.GetRequest().Proto

		scheme := req.GetRequest().URL.Scheme

		operationName := host + path + "#" + method
		span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		inject := opentracingKit.ContextToHTTP(opentracing.GlobalTracer(), logKit.NewNopLogger())
		inject(spanCtx, req.GetRequest())
		resp, err := next(spanCtx, req)

		if resp != nil {
			span.SetTag("status", strconv.Itoa(resp.StatusCode))
		}

		span.SetTag("method", method)
		span.SetTag("host", host)
		span.SetTag("path", path)
		span.SetTag("proto", proto)
		span.SetTag("scheme", scheme)

		span.LogFields(log.String("url", req.GetRequest().URL.String()))

		if err != nil {
			span.LogFields(log.String("error", err.Error()))
		}

		if builder.CustomSpanFunc != nil {
			builder.CustomSpanFunc(span, ctx, req, resp, err)
		}
		return resp, err
	}
}
