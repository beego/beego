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

package bytebuffer

import "github.com/valyala/bytebufferpool"

var (
	// Get returns an empty byte buffer from the pool, exported from beego/bytebuffer.
	Get = bytebufferpool.Get
	// GetLen returns byte buffer with fixed length from the pool, exported from beego/bytebuffer
	GetLen = func(n int) *bytebufferpool.ByteBuffer {
		b := bytebufferpool.Get()
		if cap(b.B) < n {
			b.B = make([]byte, n)
			return b
		}
		b.B = b.B[:n]
		return b
	}
	// Put returns byte buffer to the pool, exported from beego/bytebuffer
	Put = bytebufferpool.Put
)
