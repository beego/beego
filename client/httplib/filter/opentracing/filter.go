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

	logKit "github.com/go-kit/kit/log"
	opentracingKit "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/beego/beego/v2/client/httplib"
)

type FilterChainBuilder struct {
	// TagURL true will tag span with url
	TagURL bool
	// CustomSpanFunc users are able to custom their span
	CustomSpanFunc func(span opentracing.Span, ctx context.Context,
		req *httplib.BeegoHTTPRequest, resp *http.Response, err error)
}

func (builder *FilterChainBuilder) FilterChain(next httplib.Filter) httplib.Filter {

	return func(ctx context.Context, req *httplib.BeegoHTTPRequest) (*http.Response, error) {

		method := req.GetRequest().Method

		operationName := method + "#" + req.GetRequest().URL.String()
		span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		inject := opentracingKit.ContextToHTTP(opentracing.GlobalTracer(), logKit.NewNopLogger())
		inject(spanCtx, req.GetRequest())
		resp, err := next(spanCtx, req)

		if resp != nil {
			span.SetTag("http.status_code", resp.StatusCode)
		}
		span.SetTag("http.method", method)
		span.SetTag("peer.hostname", req.GetRequest().URL.Host)

		span.SetTag("http.scheme", req.GetRequest().URL.Scheme)
		span.SetTag("span.kind", "client")
		span.SetTag("component", "beego")

		if builder.TagURL {
			span.SetTag("http.url", req.GetRequest().URL.String())
		}
		span.LogFields(log.String("http.url", req.GetRequest().URL.String()))

		if err != nil {
			span.SetTag("error", true)
			span.LogFields(log.String("message", err.Error()))
		} else if resp != nil && !(resp.StatusCode < 300 && resp.StatusCode >= 200) {
			span.SetTag("error", true)
		}

		span.SetTag("peer.address", req.GetRequest().RemoteAddr)
		span.SetTag("http.proto", req.GetRequest().Proto)

		if builder.CustomSpanFunc != nil {
			builder.CustomSpanFunc(span, ctx, req, resp, err)
		}
		return resp, err
	}
}
