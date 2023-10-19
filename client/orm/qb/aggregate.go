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

// Aggregate represents the aggregation function, such as AVG, MAX, MIN, etc
type Aggregate struct {
	fn    string
	arg   string
	alias string
}

func (a Aggregate) selectable() {}

func (a Aggregate) expr() {}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: alias,
	}
}

func Avg(c string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: c,
	}
}
