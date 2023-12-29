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
	"github.com/valyala/bytebufferpool"

	"github.com/beego/beego/v2/client/orm/internal/models"
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
