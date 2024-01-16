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

package qb

import (
	"github.com/beego/beego/v2/client/orm/internal/models"
	"github.com/beego/beego/v2/client/orm/qb/errs"
	"github.com/valyala/bytebufferpool"
)

type builder struct {
	buffer *bytebufferpool.ByteBuffer
	model  *models.ModelInfo
	args   []any
}

func (b *builder) space() {
	b.writeByte(' ')
}

func (b *builder) writeString(val string) {
	_, _ = b.buffer.WriteString(val)
}

func (b *builder) writeByte(c byte) {
	_ = b.buffer.WriteByte(c)
}

func (b *builder) end() {
	b.writeByte(';')
}

func (b *builder) comma() {
	b.writeByte(',')
}

func (b *builder) buildPredicates(predicates []Predicate) error {
	p := predicates[0]
	for i := 1; i < len(predicates); i++ {
		p = p.And(predicates[i])
	}
	return b.buildExpression(p)
}

func (b *builder) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column:
		b.writeByte('`')
		b.writeString(exp.name)
		b.writeByte('`')
	case value:
		b.writeByte('?')
		b.args = append(b.args, exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			b.writeByte('(')
		}
		if err := b.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			b.writeByte(')')
		}
		b.space()
		b.writeString(exp.op.String())
		b.space()

		_, rp := exp.right.(Predicate)
		if rp {
			b.writeByte('(')
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			b.writeByte(')')
		}
	default:
		return errs.NewErrUnsupportedExpressionType(exp)
	}
	return nil
}
