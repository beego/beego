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

	"github.com/beego/beego/v2/client/orm"
)

type Mock struct {
	cond Condition
	resp []interface{}
	cb   func(inv *orm.Invocation)
}

func NewMock(cond Condition, resp []interface{}, cb func(inv *orm.Invocation)) *Mock {
	return &Mock{
		cond: cond,
		resp: resp,
		cb:   cb,
	}
}

type Condition interface {
	Match(ctx context.Context, inv *orm.Invocation) bool
}

type SimpleCondition struct {
	tableName string
	method    string
}

func NewSimpleCondition(tableName string, methodName string) Condition {
	return &SimpleCondition{
		tableName: tableName,
		method:    methodName,
	}
}

func (s *SimpleCondition) Match(ctx context.Context, inv *orm.Invocation) bool {
	res := true
	if len(s.tableName) != 0 {
		res = res && (s.tableName == inv.GetTableName())
	}

	if len(s.method) != 0 {
		res = res && (s.method == inv.Method)
	}
	return res
}
