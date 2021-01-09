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
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/beego/beego/v2/client/orm"
)

func TestFilterChainBuilder_FilterChain(t *testing.T) {
	next := func(ctx context.Context, inv *orm.Invocation) []interface{} {
		inv.TxName = "Hello"
		return []interface{}{}
	}

	builder := &FilterChainBuilder{
		CustomSpanFunc: func(span opentracing.Span, ctx context.Context, inv *orm.Invocation) {
			span.SetTag("hello", "hell")
		},
	}

	inv := &orm.Invocation{
		Method:      "Hello",
		TxStartTime: time.Now(),
	}
	builder.FilterChain(next)(context.Background(), inv)
}
