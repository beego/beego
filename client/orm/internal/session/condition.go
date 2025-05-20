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

package session

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm/clauses"
)

// ExprSep define the expression separation
const (
	ExprSep = clauses.ExprSep
)

type CondValue struct {
	Exprs  []string
	Args   []interface{}
	Cond   *Condition
	IsOr   bool
	IsNot  bool
	IsCond bool
	IsRaw  bool
	Sql    string
}

// Condition struct.
// work for WHERE conditions.
type Condition struct {
	params []CondValue
}

// NewCondition return new condition struct
func NewCondition() *Condition {
	c := &Condition{}
	return c
}

// Raw add raw sql to condition
func (c Condition) Raw(expr string, sql string) *Condition {
	if len(sql) == 0 {
		panic(fmt.Errorf("<Condition.Raw> sql cannot empty"))
	}
	c.params = append(c.params, CondValue{Exprs: strings.Split(expr, ExprSep), Sql: sql, IsRaw: true})
	return &c
}

// And add expression to condition
func (c Condition) And(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.And> args cannot empty"))
	}
	c.params = append(c.params, CondValue{Exprs: strings.Split(expr, ExprSep), Args: args})
	return &c
}

// AndNot add NOT expression to condition
func (c Condition) AndNot(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.AndNot> args cannot empty"))
	}
	c.params = append(c.params, CondValue{Exprs: strings.Split(expr, ExprSep), Args: args, IsNot: true})
	return &c
}

// AndCond combine a condition to current condition
func (c *Condition) AndCond(cond *Condition) *Condition {
	if c == cond {
		panic(fmt.Errorf("<Condition.AndCond> cannot use self as sub cond"))
	}

	c = c.clone()

	if cond != nil {
		c.params = append(c.params, CondValue{Cond: cond, IsCond: true})
	}
	return c
}

// AndNotCond combine an AND NOT condition to current condition
func (c *Condition) AndNotCond(cond *Condition) *Condition {
	c = c.clone()
	if c == cond {
		panic(fmt.Errorf("<Condition.AndNotCond> cannot use self as sub cond"))
	}

	if cond != nil {
		c.params = append(c.params, CondValue{Cond: cond, IsCond: true, IsNot: true})
	}
	return c
}

// Or add OR expression to condition
func (c Condition) Or(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.Or> args cannot empty"))
	}
	c.params = append(c.params, CondValue{Exprs: strings.Split(expr, ExprSep), Args: args, IsOr: true})
	return &c
}

// OrNot add OR NOT expression to condition
func (c Condition) OrNot(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.OrNot> args cannot empty"))
	}
	c.params = append(c.params, CondValue{Exprs: strings.Split(expr, ExprSep), Args: args, IsNot: true, IsOr: true})
	return &c
}

// OrCond combine an OR condition to current condition
func (c *Condition) OrCond(cond *Condition) *Condition {
	c = c.clone()
	if c == cond {
		panic(fmt.Errorf("<Condition.OrCond> cannot use self as sub cond"))
	}
	if cond != nil {
		c.params = append(c.params, CondValue{Cond: cond, IsCond: true, IsOr: true})
	}
	return c
}

// OrNotCond combine an OR NOT condition to current condition
func (c *Condition) OrNotCond(cond *Condition) *Condition {
	c = c.clone()
	if c == cond {
		panic(fmt.Errorf("<Condition.OrNotCond> cannot use self as sub cond"))
	}

	if cond != nil {
		c.params = append(c.params, CondValue{Cond: cond, IsCond: true, IsNot: true, IsOr: true})
	}
	return c
}

// IsEmpty check the condition arguments are empty or not.
func (c *Condition) IsEmpty() bool {
	return len(c.params) == 0
}

// clone clone a condition
func (c Condition) clone() *Condition {
	params := make([]CondValue, len(c.params))
	copy(params, c.params)
	c.params = params
	return &c
}
