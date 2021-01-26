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
	"database/sql"

	"github.com/beego/beego/v2/client/orm"
)

type DoNothingRawSetter struct {
}

func (d *DoNothingRawSetter) Exec() (sql.Result, error) {
	return nil, nil
}

func (d *DoNothingRawSetter) QueryRow(containers ...interface{}) error {
	return nil
}

func (d *DoNothingRawSetter) QueryRows(containers ...interface{}) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSetter) SetArgs(i ...interface{}) orm.RawSeter {
	return d
}

func (d *DoNothingRawSetter) Values(container *[]orm.Params, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSetter) ValuesList(container *[]orm.ParamsList, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSetter) ValuesFlat(container *orm.ParamsList, cols ...string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSetter) RowsToMap(result *orm.Params, keyCol, valueCol string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSetter) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	return 0, nil
}

func (d *DoNothingRawSetter) Prepare() (orm.RawPreparer, error) {
	return nil, nil
}