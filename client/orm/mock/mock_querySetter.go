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
	"github.com/beego/beego/v2/client/orm"
)

// DoNothingQuerySetter do nothing
// usually you use this to build your mock QuerySetter
type DoNothingQuerySetter struct {
	
}

func (d *DoNothingQuerySetter) Filter(s string, i ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) FilterRaw(s string, s2 string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) Exclude(s string, i ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) SetCond(condition *orm.Condition) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) GetCond() *orm.Condition {
	return orm.NewCondition()
}

func (d *DoNothingQuerySetter) Limit(limit interface{}, args ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) Offset(offset interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) GroupBy(exprs ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) OrderBy(exprs ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) ForceIndex(indexes ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) UseIndex(indexes ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) IgnoreIndex(indexes ...string) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) RelatedSel(params ...interface{}) orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) Distinct() orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) ForUpdate() orm.QuerySeter {
	return d
}

func (d *DoNothingQuerySetter) Count() (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) Exist() bool {
	return true
}

func (d *DoNothingQuerySetter) Update(values orm.Params) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) Delete() (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) PrepareInsert() (orm.Inserter, error) {
	return nil, nil
}

func (d *DoNothingQuerySetter) All(container interface{}, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) One(container interface{}, cols ...string) error {
	return nil
}

func (d *DoNothingQuerySetter) Values(results *[]orm.Params, exprs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) ValuesList(results *[]orm.ParamsList, exprs ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) ValuesFlat(result *orm.ParamsList, expr string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) RowsToMap(result *orm.Params, keyCol, valueCol string) (int64, error) {
	return 0, nil
}

func (d *DoNothingQuerySetter) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	return 0, nil
}