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

// DoNothingQueryM2Mer do nothing
// use it to build mock orm.QueryM2Mer
type DoNothingQueryM2Mer struct{}

func (d *DoNothingQueryM2Mer) AddWithCtx(ctx context.Context, i ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingQueryM2Mer) RemoveWithCtx(ctx context.Context, i ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingQueryM2Mer) ExistWithCtx(ctx context.Context, i interface{}) bool {
	return true
}

func (d *DoNothingQueryM2Mer) ClearWithCtx(ctx context.Context) (int64, error) {
	return 0, nil
}

func (d *DoNothingQueryM2Mer) CountWithCtx(ctx context.Context) (int64, error) {
	return 0, nil
}

func (d *DoNothingQueryM2Mer) Add(i ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingQueryM2Mer) Remove(i ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingQueryM2Mer) Exist(i interface{}) bool {
	return true
}

func (d *DoNothingQueryM2Mer) Clear() (int64, error) {
	return 0, nil
}

func (d *DoNothingQueryM2Mer) Count() (int64, error) {
	return 0, nil
}

type QueryM2MerCondition struct {
	tableName string
	name      string
}

func NewQueryM2MerCondition(tableName string, name string) *QueryM2MerCondition {
	return &QueryM2MerCondition{
		tableName: tableName,
		name:      name,
	}
}

func (q *QueryM2MerCondition) Match(ctx context.Context, inv *orm.Invocation) bool {
	res := true
	if len(q.tableName) > 0 {
		res = res && (q.tableName == inv.GetTableName())
	}
	if len(q.name) > 0 {
		res = res && (len(inv.Args) > 1) && (q.name == inv.Args[1].(string))
	}
	return res
}
