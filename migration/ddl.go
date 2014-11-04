// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package migration

type Table struct {
	TableName string
	Columns   []*Column
}

func (t *Table) Create() string {
	return ""
}

func (t *Table) Drop() string {
	return ""
}

type Column struct {
	Name    string
	Type    string
	Default interface{}
}

func Create(tbname string, columns ...Column) string {
	return ""
}

func Drop(tbname string, columns ...Column) string {
	return ""
}

func TableDDL(tbname string, columns ...Column) string {
	return ""
}
