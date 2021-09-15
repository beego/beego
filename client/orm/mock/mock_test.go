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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/orm"
)

func TestOrmStub_FilterChain(t *testing.T) {
	os := newOrmStub()
	inv := &orm.Invocation{
		Args: []interface{}{10},
	}
	i := 1
	os.FilterChain(func(ctx context.Context, inv *orm.Invocation) []interface{} {
		i++
		return nil
	})(context.Background(), inv)

	assert.Equal(t, 2, i)

	m := NewMock(NewSimpleCondition("", ""), nil, func(inv *orm.Invocation) {
		arg := inv.Args[0]
		j := arg.(int)
		inv.Args[0] = j + 1
	})
	os.Mock(m)

	os.FilterChain(nil)(context.Background(), inv)
	assert.Equal(t, 11, inv.Args[0])

	inv.Args[0] = 10
	ctxMock := NewMock(NewSimpleCondition("", ""), nil, func(inv *orm.Invocation) {
		arg := inv.Args[0]
		j := arg.(int)
		inv.Args[0] = j + 3
	})

	os.FilterChain(nil)(CtxWithMock(context.Background(), ctxMock), inv)
	assert.Equal(t, 13, inv.Args[0])
}
