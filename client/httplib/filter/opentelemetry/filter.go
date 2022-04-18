// Copyright 2021 beego
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

package opentelemetry

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/beego/beego/v2/client/httplib"
)

type CustomSpanFunc func(span trace.Span, ctx context.Context, req *httplib.BeegoHTTPRequest, resp *http.Response, err error)

type OtelFilterChainBuilder struct {
	// TagURL true will tag span with url
	tagURL bool
	// CustomSpanFunc users are able to custom their span
	customSpanFunc CustomSpanFunc
}

func NewOpenTelemetryFilter(tagURL bool, spanFunc CustomSpanFunc) *OtelFilterChainBuilder {
	return &OtelFilterChainBuilder{
		tagURL:         tagURL,
		customSpanFunc: spanFunc,
	}
}

func (builder *OtelFilterChainBuilder) FilterChain(next httplib.Filter) httplib.Filter {
	return func(ctx context.Context, req *httplib.BeegoHTTPRequest) (*http.Response, error) {
		method := req.GetRequest().Method

		operationName := method + "#" + req.GetRequest().URL.Path
		spanCtx, span := otel.Tracer("beego").Start(ctx, operationName)
		defer span.End()

		otel.GetTextMapPropagator().Inject(spanCtx, propagation.HeaderCarrier(req.GetRequest().Header))

		resp, err := next(spanCtx, req)

		if resp != nil {
			span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		}
		span.SetAttributes(attribute.String("http.method", method))
		span.SetAttributes(attribute.String("peer.hostname", req.GetRequest().URL.Host))

		span.SetAttributes(attribute.String("http.scheme", req.GetRequest().URL.Scheme))
		span.SetAttributes(attribute.String("span.kind", "client"))
		span.SetAttributes(attribute.String("component", "beego"))

		if builder.tagURL {
			span.SetAttributes(attribute.String("http.url", req.GetRequest().URL.String()))
		}

		if err != nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(err)
		} else if resp != nil && !(resp.StatusCode < 300 && resp.StatusCode >= 200) {
			span.SetAttributes(attribute.Bool("error", true))
		}

		span.SetAttributes(attribute.String("peer.address", req.GetRequest().RemoteAddr))
		span.SetAttributes(attribute.String("http.proto", req.GetRequest().Proto))

		if builder.customSpanFunc != nil {
			builder.customSpanFunc(span, ctx, req, resp, err)
		}
		return resp, err
	}
}
