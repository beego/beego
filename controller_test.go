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
	"testing"

	"github.com/astaxie/beego/context"
)

func TestGetInt(t *testing.T) {
	i := context.NewInput()
	i.SetParam("age", "40")
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt("age")
	if val != 40 {
		t.Errorf("TestGetInt expect 40,get %T,%v", val, val)
	}
}

func TestGetInt8(t *testing.T) {
	i := context.NewInput()
	i.SetParam("age", "40")
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt8("age")
	if val != 40 {
		t.Errorf("TestGetInt8 expect 40,get %T,%v", val, val)
	}
	//Output: int8
}

func TestGetInt16(t *testing.T) {
	i := context.NewInput()
	i.SetParam("age", "40")
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt16("age")
	if val != 40 {
		t.Errorf("TestGetInt16 expect 40,get %T,%v", val, val)
	}
}

func TestGetInt32(t *testing.T) {
	i := context.NewInput()
	i.SetParam("age", "40")
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt32("age")
	if val != 40 {
		t.Errorf("TestGetInt32 expect 40,get %T,%v", val, val)
	}
}

func TestGetInt64(t *testing.T) {
	i := context.NewInput()
	i.SetParam("age", "40")
	ctx := &context.Context{Input: i}
	ctrlr := Controller{Ctx: ctx}
	val, _ := ctrlr.GetInt64("age")
	if val != 40 {
		t.Errorf("TestGeetInt64 expect 40,get %T,%v", val, val)
	}
}
