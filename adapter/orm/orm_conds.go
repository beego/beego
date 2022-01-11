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

package orm

import (
	"github.com/beego/beego/v2/client/orm"
)

// ExprSep define the expression separation
const (
	ExprSep = "__"
)

// Condition struct.
// work for WHERE conditions.
type Condition orm.Condition

// NewCondition return new condition struct
func NewCondition() *Condition {
	return (*Condition)(orm.NewCondition())
}

// Raw add raw sql to condition
func (c Condition) Raw(expr string, sql string) *Condition {
	return (*Condition)((orm.Condition)(c).Raw(expr, sql))
}

// And add expression to condition
func (c Condition) And(expr string, args ...interface{}) *Condition {
	return (*Condition)((orm.Condition)(c).And(expr, args...))
}

// AndNot add NOT expression to condition
func (c Condition) AndNot(expr string, args ...interface{}) *Condition {
	return (*Condition)((orm.Condition)(c).AndNot(expr, args...))
}

// AndCond combine a condition to current condition
func (c *Condition) AndCond(cond *Condition) *Condition {
	return (*Condition)((*orm.Condition)(c).AndCond((*orm.Condition)(cond)))
}

// AndNotCond combine a AND NOT condition to current condition
func (c *Condition) AndNotCond(cond *Condition) *Condition {
	return (*Condition)((*orm.Condition)(c).AndNotCond((*orm.Condition)(cond)))
}

// Or add OR expression to condition
func (c Condition) Or(expr string, args ...interface{}) *Condition {
	return (*Condition)((orm.Condition)(c).Or(expr, args...))
}

// OrNot add OR NOT expression to condition
func (c Condition) OrNot(expr string, args ...interface{}) *Condition {
	return (*Condition)((orm.Condition)(c).OrNot(expr, args...))
}

// OrCond combine a OR condition to current condition
func (c *Condition) OrCond(cond *Condition) *Condition {
	return (*Condition)((*orm.Condition)(c).OrCond((*orm.Condition)(cond)))
}

// OrNotCond combine a OR NOT condition to current condition
func (c *Condition) OrNotCond(cond *Condition) *Condition {
	return (*Condition)((*orm.Condition)(c).OrNotCond((*orm.Condition)(cond)))
}

// IsEmpty check the condition arguments are empty or not.
func (c *Condition) IsEmpty() bool {
	return (*orm.Condition)(c).IsEmpty()
}
