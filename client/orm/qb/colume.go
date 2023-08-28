// Copyright 2023 beego. All Rights Reserved.
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

package qb

// Column contains the info of single column
type Column struct {
	name  string
	alias string
	order string
}

func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
		order: c.order,
	}
}
func (c Column) selectable() {}

func (c Column) expr() {}
func (c Column) Asc() Column {
	return Column{
		name:  c.name,
		alias: c.alias,
		order: " ASC",
	}
}
func (c Column) Desc() Column {
	return Column{
		name:  c.name,
		alias: c.alias,
		order: " DESC",
	}
}

type value struct {
	val any
}

func (c value) expr() {}

func valueOf(val any) value {
	return value{
		val: val,
	}
}

func C(name string) Column {
	return Column{
		name: name,
	}
}
func (c Column) EQ(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEqual,
		right: exprOf(arg),
	}
}
func (c Column) GT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: exprOf(arg),
	}
}
func (c Column) LT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: exprOf(arg),
	}
}
