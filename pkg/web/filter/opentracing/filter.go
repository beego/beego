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
	"github.com/opentracing/opentracing-go"

	beego "github.com/astaxie/beego/pkg"
	"github.com/astaxie/beego/pkg/context"
)

// FilterChainBuilder provides an extension point that we can support more configurations if necessary
type FilterChainBuilder struct {
	// CustomSpanFunc makes users to custom the span.
	CustomSpanFunc func(span opentracing.Span, ctx *context.Context)
}


func (builder *FilterChainBuilder) FilterChain(next beego.FilterFunc) beego.FilterFunc {
	return func(ctx *context.Context) {
		span := opentracing.SpanFromContext(ctx.Request.Context())
		spanCtx := ctx.Request.Context()
		if span == nil {
			operationName := ctx.Input.URL()
			// it means that there is not any span, so we create a span as the root span.
			// TODO, if we support multiple servers, this need to be changed
			route, found := beego.BeeApp.Handlers.FindRouter(ctx)
			if found {
				operationName = route.GetPattern()
			}
			span, spanCtx = opentracing.StartSpanFromContext(spanCtx, operationName)
			newReq := ctx.Request.Clone(spanCtx)
			ctx.Reset(ctx.ResponseWriter.ResponseWriter, newReq)
		}

		defer span.Finish()
		next(ctx)
		// if you think we need to do more things, feel free to create an issue to tell us
		span.SetTag("status", ctx.Output.Status)
		span.SetTag("method", ctx.Input.Method())
		span.SetTag("route", ctx.Input.GetData("RouterPattern"))
		if builder.CustomSpanFunc != nil {
			builder.CustomSpanFunc(span, ctx)
		}
	}
}
