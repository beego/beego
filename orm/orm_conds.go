// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie, slene

package orm

import (
	"fmt"
	"strings"
)

const (
	ExprSep = "__"
)

type condValue struct {
	exprs  []string
	args   []interface{}
	cond   *Condition
	isOr   bool
	isNot  bool
	isCond bool
}

// condition struct.
// work for WHERE conditions.
type Condition struct {
	params []condValue
}

// return new condition struct
func NewCondition() *Condition {
	c := &Condition{}
	return c
}

// add expression to condition
func (c Condition) And(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.And> args cannot empty"))
	}
	c.params = append(c.params, condValue{exprs: strings.Split(expr, ExprSep), args: args})
	return &c
}

// add NOT expression to condition
func (c Condition) AndNot(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.AndNot> args cannot empty"))
	}
	c.params = append(c.params, condValue{exprs: strings.Split(expr, ExprSep), args: args, isNot: true})
	return &c
}

// combine a condition to current condition
func (c *Condition) AndCond(cond *Condition) *Condition {
	c = c.clone()
	if c == cond {
		panic(fmt.Errorf("<Condition.AndCond> cannot use self as sub cond"))
	}
	if cond != nil {
		c.params = append(c.params, condValue{cond: cond, isCond: true})
	}
	return c
}

// add OR expression to condition
func (c Condition) Or(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.Or> args cannot empty"))
	}
	c.params = append(c.params, condValue{exprs: strings.Split(expr, ExprSep), args: args, isOr: true})
	return &c
}

// add OR NOT expression to condition
func (c Condition) OrNot(expr string, args ...interface{}) *Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition.OrNot> args cannot empty"))
	}
	c.params = append(c.params, condValue{exprs: strings.Split(expr, ExprSep), args: args, isNot: true, isOr: true})
	return &c
}

// combine a OR condition to current condition
func (c *Condition) OrCond(cond *Condition) *Condition {
	c = c.clone()
	if c == cond {
		panic(fmt.Errorf("<Condition.OrCond> cannot use self as sub cond"))
	}
	if cond != nil {
		c.params = append(c.params, condValue{cond: cond, isCond: true, isOr: true})
	}
	return c
}

// check the condition arguments are empty or not.
func (c *Condition) IsEmpty() bool {
	return len(c.params) == 0
}

// clone a condition
func (c Condition) clone() *Condition {
	return &c
}
