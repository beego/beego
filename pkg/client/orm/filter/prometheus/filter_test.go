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

package prometheus

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/astaxie/beego/pkg/client/orm"
)

func TestFilterChainBuilder_FilterChain1(t *testing.T) {
	next := func(ctx context.Context, inv *orm.Invocation) []interface{} {
		inv.Method = "coming"
		return []interface{}{}
	}
	builder := &FilterChainBuilder{}
	filter := builder.FilterChain(next)

	assert.NotNil(t, builder.summaryVec)
	assert.NotNil(t, filter)

	inv := &orm.Invocation{}
	filter(context.Background(), inv)
	assert.Equal(t, "coming", inv.Method)

	inv = &orm.Invocation{
		Method:      "Hello",
		TxStartTime: time.Now(),
	}
	builder.reportTxn(context.Background(), inv)

	inv = &orm.Invocation{
		Method: "Begin",
	}

	ctx := context.Background()
	// it will be ignored
	builder.report(ctx, inv, time.Second)

	inv.Method = "Commit"
	builder.report(ctx, inv, time.Second)

	inv.Method = "Update"
	builder.report(ctx, inv, time.Second)

}
