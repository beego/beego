// Copyright 2021 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"github.com/beego/beego/v2/core/berror"
)

var NilCacheAdapter = berror.DefineCode(4002001, moduleName, "NilCacheAdapter", `
It means that you register cache adapter by pass nil.
A cache adapter is an instance of Cache interface. 
`)

var DuplicateAdapter = berror.DefineCode(4002002, moduleName, "DuplicateAdapter", `
You register two adapter with same name. In beego cache module, one name one adapter.
Once you got this error, please check the error stack, search adapter 
`)

var UnknownAdapter = berror.DefineCode(4002003, moduleName, "UnknownAdapter", `
Unknown adapter, do you forget to register the adapter?
You must register adapter before use it. For example, if you want to use redis implementation, 
you must import the cache/redis package.
`)

var IncrementOverflow = berror.DefineCode(4002004, moduleName, "IncrementOverflow", `
The increment operation will overflow.
`)

var DecrementOverflow = berror.DefineCode(4002005, moduleName, "DecrementOverflow", `
The decrement operation will overflow.
`)

var NotIntegerType = berror.DefineCode(4002006, moduleName, "NotIntegerType", `
The type of value is not (u)int (u)int32 (u)int64. 
When you want to call Incr or Decr function of Cache API, you must confirm that the value's type is one of (u)int (u)int32 (u)int64.
`)