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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/astaxie/beego/context"
)

type TestController struct {
	Controller
}

func (tc *TestController) Get() {
	tc.Data["Username"] = "astaxie"
	tc.Ctx.Output.Body([]byte("ok"))
}

func (tc *TestController) Post() {
	tc.Ctx.Output.Body([]byte(tc.Ctx.Input.Query(":name")))
}

func (tc *TestController) Param() {
	tc.Ctx.Output.Body([]byte(tc.Ctx.Input.Query(":name")))
}

func (tc *TestController) List() {
	tc.Ctx.Output.Body([]byte("i am list"))
}

func (tc *TestController) Params() {
	tc.Ctx.Output.Body([]byte(tc.Ctx.Input.Params["0"] + tc.Ctx.Input.Params["1"] + tc.Ctx.Input.Params["2"]))
}

func (tc *TestController) Myext() {
	tc.Ctx.Output.Body([]byte(tc.Ctx.Input.Param(":ext")))
}

func (tc *TestController) GetUrl() {
	tc.Ctx.Output.Body([]byte(tc.UrlFor(".Myext")))
}

func (t *TestController) GetParams() {
	t.Ctx.WriteString(t.Ctx.Input.Query(":last") + "+" +
		t.Ctx.Input.Query(":first") + "+" + t.Ctx.Input.Query("learn"))
}

func (t *TestController) GetManyRouter() {
	t.Ctx.WriteString(t.Ctx.Input.Query(":id") + t.Ctx.Input.Query(":page"))
}

type ResStatus struct {
	Code int
	Msg  string
}

type JsonController struct {
	Controller
}

func (this *JsonController) Prepare() {
	this.Data["json"] = "prepare"
	this.ServeJson(true)
}

func (this *JsonController) Get() {
	this.Data["Username"] = "astaxie"
	this.Ctx.Output.Body([]byte("ok"))
}

func TestUrlFor(t *testing.T) {
	handler := NewControllerRegister()
	handler.Add("/api/list", &TestController{}, "*:List")
	handler.Add("/person/:last/:first", &TestController{}, "*:Param")
	if a := handler.UrlFor("TestController.List"); a != "/api/list" {
		Info(a)
		t.Errorf("TestController.List must equal to /api/list")
	}
	if a := handler.UrlFor("TestController.Param", ":last", "xie", ":first", "asta"); a != "/person/xie/asta" {
		t.Errorf("TestController.Param must equal to /person/xie/asta, but get " + a)
	}
}

func TestUrlFor3(t *testing.T) {
	handler := NewControllerRegister()
	handler.AddAuto(&TestController{})
	if a := handler.UrlFor("TestController.Myext"); a != "/test/myext" && a != "/Test/Myext" {
		t.Errorf("TestController.Myext must equal to /test/myext, but get " + a)
	}
	if a := handler.UrlFor("TestController.GetUrl"); a != "/test/geturl" && a != "/Test/GetUrl" {
		t.Errorf("TestController.GetUrl must equal to /test/geturl, but get " + a)
	}
}

func TestUrlFor2(t *testing.T) {
	handler := NewControllerRegister()
	handler.Add("/v1/:v/cms_:id(.+)_:page(.+).html", &TestController{}, "*:List")
	handler.Add("/v1/:username/edit", &TestController{}, "get:GetUrl")
	handler.Add("/v1/:v(.+)_cms/ttt_:id(.+)_:page(.+).html", &TestController{}, "*:Param")
	handler.Add("/:year:int/:month:int/:title/:entid", &TestController{})
	if handler.UrlFor("TestController.GetUrl", ":username", "astaxie") != "/v1/astaxie/edit" {
		Info(handler.UrlFor("TestController.GetUrl"))
		t.Errorf("TestController.List must equal to /v1/astaxie/edit")
	}

	if handler.UrlFor("TestController.List", ":v", "za", ":id", "12", ":page", "123") !=
		"/v1/za/cms_12_123.html" {
		Info(handler.UrlFor("TestController.List"))
		t.Errorf("TestController.List must equal to /v1/za/cms_12_123.html")
	}
	if handler.UrlFor("TestController.Param", ":v", "za", ":id", "12", ":page", "123") !=
		"/v1/za_cms/ttt_12_123.html" {
		Info(handler.UrlFor("TestController.Param"))
		t.Errorf("TestController.List must equal to /v1/za_cms/ttt_12_123.html")
	}
	if handler.UrlFor("TestController.Get", ":year", "1111", ":month", "11",
		":title", "aaaa", ":entid", "aaaa") !=
		"/1111/11/aaaa/aaaa" {
		Info(handler.UrlFor("TestController.Get"))
		t.Errorf("TestController.Get must equal to /1111/11/aaaa/aaaa")
	}
}

func TestUserFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/list", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/api/list", &TestController{}, "*:List")
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("user define func can't run")
	}
}

func TestPostFunc(t *testing.T) {
	r, _ := http.NewRequest("POST", "/astaxie", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/:name", &TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "astaxie" {
		t.Errorf("post func should astaxie")
	}
}

func TestAutoFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/list", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.AddAuto(&TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("user define func can't run")
	}
}

func TestAutoFunc2(t *testing.T) {
	r, _ := http.NewRequest("GET", "/Test/List", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.AddAuto(&TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("user define func can't run")
	}
}

func TestAutoFuncParams(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/params/2009/11/12", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.AddAuto(&TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "20091112" {
		t.Errorf("user define func can't run")
	}
}

func TestAutoExtFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/myext.json", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.AddAuto(&TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "json" {
		t.Errorf("user define func can't run")
	}
}

func TestRouteOk(t *testing.T) {

	r, _ := http.NewRequest("GET", "/person/anderson/thomas?learn=kungfu", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/person/:last/:first", &TestController{}, "get:GetParams")
	handler.ServeHTTP(w, r)
	body := w.Body.String()
	if body != "anderson+thomas+kungfu" {
		t.Errorf("url param set to [%s];", body)
	}
}

func TestManyRoute(t *testing.T) {

	r, _ := http.NewRequest("GET", "/beego32-12.html", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/beego:id([0-9]+)-:page([0-9]+).html", &TestController{}, "get:GetManyRouter")
	handler.ServeHTTP(w, r)

	body := w.Body.String()

	if body != "3212" {
		t.Errorf("url param set to [%s];", body)
	}
}

func TestNotFound(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Code set to [%v]; want [%v]", w.Code, http.StatusNotFound)
	}
}

// TestStatic tests the ability to serve static
// content from the filesystem
func TestStatic(t *testing.T) {
	r, _ := http.NewRequest("GET", "/static/js/jquery.js", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("handler.Static failed to serve file")
	}
}

func TestPrepare(t *testing.T) {
	r, _ := http.NewRequest("GET", "/json/list", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/json/list", &JsonController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != `"prepare"` {
		t.Errorf(w.Body.String() + "user define func can't run")
	}
}

func TestAutoPrefix(t *testing.T) {
	r, _ := http.NewRequest("GET", "/admin/test/list", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.AddAutoPrefix("/admin", &TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("TestAutoPrefix can't run")
	}
}

func TestRouterGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Get("/user", func(ctx *context.Context) {
		ctx.Output.Body([]byte("Get userlist"))
	})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "Get userlist" {
		t.Errorf("TestRouterGet can't run")
	}
}

func TestRouterPost(t *testing.T) {
	r, _ := http.NewRequest("POST", "/user/123", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Post("/user/:id", func(ctx *context.Context) {
		ctx.Output.Body([]byte(ctx.Input.Param(":id")))
	})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "123" {
		t.Errorf("TestRouterPost can't run")
	}
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sayhello"))
}

func TestRouterHandler(t *testing.T) {
	r, _ := http.NewRequest("POST", "/sayhi", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Handler("/sayhi", http.HandlerFunc(sayhello))
	handler.ServeHTTP(w, r)
	if w.Body.String() != "sayhello" {
		t.Errorf("TestRouterHandler can't run")
	}
}

//
// Benchmarks NewApp:
//

func beegoFilterFunc(ctx *context.Context) {
	ctx.WriteString("hello")
}

type AdminController struct {
	Controller
}

func (a *AdminController) Get() {
	a.Ctx.WriteString("hello")
}

func TestRouterFunc(t *testing.T) {
	mux := NewControllerRegister()
	mux.Get("/action", beegoFilterFunc)
	mux.Post("/action", beegoFilterFunc)
	rw, r := testRequest("GET", "/action")
	mux.ServeHTTP(rw, r)
	if rw.Body.String() != "hello" {
		t.Errorf("TestRouterFunc can't run")
	}
}

func BenchmarkFunc(b *testing.B) {
	mux := NewControllerRegister()
	mux.Get("/action", beegoFilterFunc)
	rw, r := testRequest("GET", "/action")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(rw, r)
	}
}

func BenchmarkController(b *testing.B) {
	mux := NewControllerRegister()
	mux.Add("/action", &AdminController{})
	rw, r := testRequest("GET", "/action")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(rw, r)
	}
}

func testRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()

	return recorder, request
}

// Execution point: BeforeRouter
// expectation: only BeforeRouter function is executed, notmatch output as router doesn't handle
func TestFilterBeforeRouter(t *testing.T) {
	testName := "TestFilterBeforeRouter"
	url := "/beforeRouter"

	mux := NewControllerRegister()
	mux.InsertFilter(url, BeforeRouter, beegoBeforeRouter1)

	mux.Get(url, beegoFilterFunc)

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if strings.Contains(rw.Body.String(), "BeforeRouter1") == false {
		t.Errorf(testName + " BeforeRouter did not run")
	}
	if strings.Contains(rw.Body.String(), "hello") == true {
		t.Errorf(testName + " BeforeRouter did not return properly")
	}
}

// Execution point: BeforeExec
// expectation: only BeforeExec function is executed, match as router determines route only
func TestFilterBeforeExec(t *testing.T) {
	testName := "TestFilterBeforeExec"
	url := "/beforeExec"

	mux := NewControllerRegister()
	mux.InsertFilter(url, BeforeRouter, beegoFilterNoOutput)
	mux.InsertFilter(url, BeforeExec, beegoBeforeExec1)

	mux.Get(url, beegoFilterFunc)

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if strings.Contains(rw.Body.String(), "BeforeExec1") == false {
		t.Errorf(testName + " BeforeExec did not run")
	}
	if strings.Contains(rw.Body.String(), "hello") == true {
		t.Errorf(testName + " BeforeExec did not return properly")
	}
	if strings.Contains(rw.Body.String(), "BeforeRouter") == true {
		t.Errorf(testName + " BeforeRouter ran in error")
	}
}

// Execution point: AfterExec
// expectation: only AfterExec function is executed, match as router handles
func TestFilterAfterExec(t *testing.T) {
	testName := "TestFilterAfterExec"
	url := "/afterExec"

	mux := NewControllerRegister()
	mux.InsertFilter(url, BeforeRouter, beegoFilterNoOutput)
	mux.InsertFilter(url, BeforeExec, beegoFilterNoOutput)
	mux.InsertFilter(url, AfterExec, beegoAfterExec1, false)

	mux.Get(url, beegoFilterFunc)

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if strings.Contains(rw.Body.String(), "AfterExec1") == false {
		t.Errorf(testName + " AfterExec did not run")
	}
	if strings.Contains(rw.Body.String(), "hello") == false {
		t.Errorf(testName + " handler did not run properly")
	}
	if strings.Contains(rw.Body.String(), "BeforeRouter") == true {
		t.Errorf(testName + " BeforeRouter ran in error")
	}
	if strings.Contains(rw.Body.String(), "BeforeExec") == true {
		t.Errorf(testName + " BeforeExec ran in error")
	}
}

// Execution point: FinishRouter
// expectation: only FinishRouter function is executed, match as router handles
func TestFilterFinishRouter(t *testing.T) {
	testName := "TestFilterFinishRouter"
	url := "/finishRouter"

	mux := NewControllerRegister()
	mux.InsertFilter(url, BeforeRouter, beegoFilterNoOutput)
	mux.InsertFilter(url, BeforeExec, beegoFilterNoOutput)
	mux.InsertFilter(url, AfterExec, beegoFilterNoOutput)
	mux.InsertFilter(url, FinishRouter, beegoFinishRouter1)

	mux.Get(url, beegoFilterFunc)

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if strings.Contains(rw.Body.String(), "FinishRouter1") == true {
		t.Errorf(testName + " FinishRouter did not run")
	}
	if strings.Contains(rw.Body.String(), "hello") == false {
		t.Errorf(testName + " handler did not run properly")
	}
	if strings.Contains(rw.Body.String(), "AfterExec1") == true {
		t.Errorf(testName + " AfterExec ran in error")
	}
	if strings.Contains(rw.Body.String(), "BeforeRouter") == true {
		t.Errorf(testName + " BeforeRouter ran in error")
	}
	if strings.Contains(rw.Body.String(), "BeforeExec") == true {
		t.Errorf(testName + " BeforeExec ran in error")
	}
}

// Execution point: FinishRouter
// expectation: only first FinishRouter function is executed, match as router handles
func TestFilterFinishRouterMultiFirstOnly(t *testing.T) {
	testName := "TestFilterFinishRouterMultiFirstOnly"
	url := "/finishRouterMultiFirstOnly"

	mux := NewControllerRegister()
	mux.InsertFilter(url, FinishRouter, beegoFinishRouter1, false)
	mux.InsertFilter(url, FinishRouter, beegoFinishRouter2)

	mux.Get(url, beegoFilterFunc)

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if strings.Contains(rw.Body.String(), "FinishRouter1") == false {
		t.Errorf(testName + " FinishRouter1 did not run")
	}
	if strings.Contains(rw.Body.String(), "hello") == false {
		t.Errorf(testName + " handler did not run properly")
	}
	// not expected in body
	if strings.Contains(rw.Body.String(), "FinishRouter2") == true {
		t.Errorf(testName + " FinishRouter2 did run")
	}
}

// Execution point: FinishRouter
// expectation: both FinishRouter functions execute, match as router handles
func TestFilterFinishRouterMulti(t *testing.T) {
	testName := "TestFilterFinishRouterMulti"
	url := "/finishRouterMulti"

	mux := NewControllerRegister()
	mux.InsertFilter(url, FinishRouter, beegoFinishRouter1, false)
	mux.InsertFilter(url, FinishRouter, beegoFinishRouter2, false)

	mux.Get(url, beegoFilterFunc)

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if strings.Contains(rw.Body.String(), "FinishRouter1") == false {
		t.Errorf(testName + " FinishRouter1 did not run")
	}
	if strings.Contains(rw.Body.String(), "hello") == false {
		t.Errorf(testName + " handler did not run properly")
	}
	if strings.Contains(rw.Body.String(), "FinishRouter2") == false {
		t.Errorf(testName + " FinishRouter2 did not run properly")
	}
}

func beegoFilterNoOutput(ctx *context.Context) {
	return
}
func beegoBeforeRouter1(ctx *context.Context) {
	ctx.WriteString("|BeforeRouter1")
}
func beegoBeforeRouter2(ctx *context.Context) {
	ctx.WriteString("|BeforeRouter2")
}
func beegoBeforeExec1(ctx *context.Context) {
	ctx.WriteString("|BeforeExec1")
}
func beegoBeforeExec2(ctx *context.Context) {
	ctx.WriteString("|BeforeExec2")
}
func beegoAfterExec1(ctx *context.Context) {
	ctx.WriteString("|AfterExec1")
}
func beegoAfterExec2(ctx *context.Context) {
	ctx.WriteString("|AfterExec2")
}
func beegoFinishRouter1(ctx *context.Context) {
	ctx.WriteString("|FinishRouter1")
}
func beegoFinishRouter2(ctx *context.Context) {
	ctx.WriteString("|FinishRouter2")
}
