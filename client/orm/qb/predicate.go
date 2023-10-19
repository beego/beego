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

type op string

const (
	opEqual op = "="
	opLT    op = "<"
	opGT    op = ">"
	opAnd   op = "AND"
	opOr    op = "OR"
	opNot   op = "NOT"
)

func (o op) String() string {
	return string(o)
}

// Predicate Represents a query condition
type Predicate struct {
	left  Expression
	op    op
	right Expression
}

// Expression Represents a statement
type Expression interface {
	expr()
}

func (Predicate) expr() {}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return valueOf(exp)
	}
}
func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAnd,
		right: r,
	}
}
func (p Predicate) Or(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOr,
		right: r,
	}
}
func Not(r Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: r,
	}
}
