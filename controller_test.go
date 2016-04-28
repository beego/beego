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

	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego/context"
	"fmt"
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

type testController struct {
	Controller
}

func (t *testController) Get() {
	typ := t.GetString("type")
	code, _ := t.GetInt("code")
	switch typ {
	case "abort":
		t.Abort(fmt.Sprint(code))
	case "abort-content":
		t.Ctx.Abort(code,t.GetString("body"))
	default:
		t.CustomAbort(code, t.GetString("body"))
	}
}

func TestController_Abort_01(t *testing.T) {
	mux := NewControllerRegister()
	mux.Add("/test", &testController{})
	hrw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/test?type=abort&code=401", nil)
	mux.ServeHTTP(hrw, r)
	if hrw.Code != 401 {
		t.Log(hrw.Code)
		t.FailNow()
	}
}

func TestController_Abort_02(t *testing.T) {
	mux := NewControllerRegister()
	mux.Add("/test", &testController{})
	hrw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.com/test?type=default&code=501&body=testController", nil)
	mux.ServeHTTP(hrw, r)
	if hrw.Code != 501{
		t.Log(hrw.Code)
		t.FailNow()
	}
	if string(hrw.Body.Bytes())!="testController"{
		t.Log(string(hrw.Body.Bytes()))
		t.FailNow()
	}
}

func TestController_Abort_03(t *testing.T) {
	mux := NewControllerRegister()
	mux.Add("/test", &testController{})
	hrw := httptest.NewRecorder()
	registerDefaultErrorHandler()
	r, _ := http.NewRequest("GET", "http://example.com/test?type=default&code=501&body=testController", nil)
	mux.ServeHTTP(hrw, r)
	if hrw.Code != 501{
		t.Log(hrw.Code)
		t.FailNow()
	}
	if string(hrw.Body.Bytes())!="testController"{
		t.Log(string(hrw.Body.Bytes()))
		t.FailNow()
	}
}

func TestController_Abort_04(t *testing.T) {
	mux := NewControllerRegister()
	mux.Add("/test", &testController{})
	hrw := httptest.NewRecorder()
	registerDefaultErrorHandler()
	r, _ := http.NewRequest("GET", "http://example.com/test?type=abort-content&code=501&body=501", nil)
	mux.ServeHTTP(hrw, r)
	if hrw.Code != 501{
		t.Log(hrw.Code)
		t.FailNow()
	}
	//if string(hrw.Body.Bytes())!="testController"{
	//	t.Log(string(hrw.Body.Bytes()))
	//	t.FailNow()
	//}
}
