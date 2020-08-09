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

	"github.com/opentracing/opentracing-go"

	"github.com/astaxie/beego/pkg/orm"
)

// FilterChainBuilder provides an extension point
// this Filter's behavior looks a little bit strange
// for example:
// if we want to trace QuerySetter
// actually we trace invoking "QueryTable" and "QueryTableWithCtx"
type FilterChainBuilder struct {
	// CustomSpanFunc users are able to custom their span
	CustomSpanFunc func(span opentracing.Span, ctx context.Context, inv *orm.Invocation)
}

func (builder *FilterChainBuilder) FilterChain(next orm.Filter) orm.Filter {
	return func(ctx context.Context, inv *orm.Invocation) {
		operationName := builder.operationName(ctx, inv)
		span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		next(spanCtx, inv)
		span.SetTag("Method", inv.Method)
		span.SetTag("Table", inv.GetTableName())
		span.SetTag("InsideTx", inv.InsideTx)
		span.SetTag("TxName", spanCtx.Value(orm.TxNameKey))

		if builder.CustomSpanFunc != nil {
			builder.CustomSpanFunc(span, spanCtx, inv)
		}

	}
}

func (builder *FilterChainBuilder) operationName(ctx context.Context, inv *orm.Invocation) string {
	if n, ok := ctx.Value(orm.TxNameKey).(string); ok {
		return inv.Method + "#" + n
	}
	return inv.Method + "#" + inv.GetTableName()
}
