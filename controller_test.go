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

package beego

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

func ExampleGetInt() {

	i := &context.BeegoInput{Params: map[string]string{"age": "40"}}
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}

	val, _ := ctrlr.GetInt("age")
	fmt.Printf("%T", val)
	//Output: int
}

func ExampleGetInt8() {

	i := &context.BeegoInput{Params: map[string]string{"age": "40"}}
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}

	val, _ := ctrlr.GetInt8("age")
	fmt.Printf("%T", val)
	//Output: int8
}

func ExampleGetInt16() {

	i := &context.BeegoInput{Params: map[string]string{"age": "40"}}
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}

	val, _ := ctrlr.GetInt16("age")
	fmt.Printf("%T", val)
	//Output: int16
}

func ExampleGetInt32() {

	i := &context.BeegoInput{Params: map[string]string{"age": "40"}}
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}

	val, _ := ctrlr.GetInt32("age")
	fmt.Printf("%T", val)
	//Output: int32
}

func ExampleGetInt64() {

	i := &context.BeegoInput{Params: map[string]string{"age": "40"}}
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}

	val, _ := ctrlr.GetInt64("age")
	fmt.Printf("%T", val)
	//Output: int64
}
