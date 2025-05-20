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
	"github.com/beego/beego/v2/client/orm/internal/session"
)

var stub = newOrmStub()

func init() {
	session.AddGlobalFilterChain(stub.FilterChain)
}

type Stub interface {
	Mock(m *Mock)
	Clear()
}

type OrmStub struct {
	ms []*Mock
}

func StartMock() Stub {
	return stub
}

func newOrmStub() *OrmStub {
	return &OrmStub{
		ms: make([]*Mock, 0, 4),
	}
}

func (o *OrmStub) Mock(m *Mock) {
	o.ms = append(o.ms, m)
}

func (o *OrmStub) Clear() {
	o.ms = make([]*Mock, 0, 4)
}

func (o *OrmStub) FilterChain(next session.Filter) session.Filter {
	return func(ctx context.Context, inv *session.Invocation) []interface{} {
		ms := mockFromCtx(ctx)
		ms = append(ms, o.ms...)

		for _, mock := range ms {
			if mock.cond.Match(ctx, inv) {
				if mock.cb != nil {
					mock.cb(inv)
				}
				return mock.resp
			}
		}
		return next(ctx, inv)
	}
}
