// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package beego

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego/context"
)

type TestController struct {
	Controller
}

func (this *TestController) Get() {
	this.Data["Username"] = "astaxie"
	this.Ctx.Output.Body([]byte("ok"))
}

func (this *TestController) Post() {
	this.Ctx.Output.Body([]byte(this.Ctx.Input.Query(":name")))
}

func (this *TestController) Param() {
	this.Ctx.Output.Body([]byte(this.Ctx.Input.Query(":name")))
}

func (this *TestController) List() {
	this.Ctx.Output.Body([]byte("i am list"))
}

func (this *TestController) Params() {
	this.Ctx.Output.Body([]byte(this.Ctx.Input.Params["0"] + this.Ctx.Input.Params["1"] + this.Ctx.Input.Params["2"]))
}

func (this *TestController) Myext() {
	this.Ctx.Output.Body([]byte(this.Ctx.Input.Param(":ext")))
}

func (this *TestController) GetUrl() {
	this.Ctx.Output.Body([]byte(this.UrlFor(".Myext")))
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
	if a := handler.UrlFor("TestController.Myext"); a != "/test/myext" {
		t.Errorf("TestController.Myext must equal to /test/myext, but get " + a)
	}
	if a := handler.UrlFor("TestController.GetUrl"); a != "/test/geturl" {
		t.Errorf("TestController.GetUrl must equal to /test/geturl, but get " + a)
	}
}

func TestUrlFor2(t *testing.T) {
	handler := NewControllerRegister()
	handler.Add("/v1/:v/cms_:id(.+)_:page(.+).html", &TestController{}, "*:List")
	handler.Add("/v1/:v(.+)_cms/ttt_:id(.+)_:page(.+).html", &TestController{}, "*:Param")
	handler.Add("/:year:int/:month:int/:title/:entid", &TestController{})
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
