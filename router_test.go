package beego

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestController struct {
	Controller
}

func (this *TestController) Get() {
	this.Data["Username"] = "astaxie"
	this.Ctx.Output.Body([]byte("ok"))
}

func (this *TestController) List() {
	this.Ctx.Output.Body([]byte("i am list"))
}

func (this *TestController) Myext() {
	this.Ctx.Output.Body([]byte(this.Ctx.Input.Params(":ext")))
}

func TestUserFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/api/list", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegistor()
	handler.Add("/api/list", &TestController{}, "*:List")
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("user define func can't run")
	}
}

func TestAutoFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/list", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegistor()
	handler.AddAuto(&TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("user define func can't run")
	}
}

func TestAutoExtFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test/myext.json", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegistor()
	handler.AddAuto(&TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "json" {
		t.Errorf("user define func can't run")
	}
}

func TestRouteOk(t *testing.T) {

	r, _ := http.NewRequest("GET", "/person/anderson/thomas?learn=kungfu", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegistor()
	handler.Add("/person/:last/:first", &TestController{})
	handler.ServeHTTP(w, r)

	lastNameParam := r.URL.Query().Get(":last")
	firstNameParam := r.URL.Query().Get(":first")
	learnParam := r.URL.Query().Get("learn")

	if lastNameParam != "anderson" {
		t.Errorf("url param set to [%s]; want [%s]", lastNameParam, "anderson")
	}
	if firstNameParam != "thomas" {
		t.Errorf("url param set to [%s]; want [%s]", firstNameParam, "thomas")
	}
	if learnParam != "kungfu" {
		t.Errorf("url param set to [%s]; want [%s]", learnParam, "kungfu")
	}
}

func TestManyRoute(t *testing.T) {

	r, _ := http.NewRequest("GET", "/beego32-12.html", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegistor()
	handler.Add("/beego:id([0-9]+)-:page([0-9]+).html", &TestController{})
	handler.ServeHTTP(w, r)

	id := r.URL.Query().Get(":id")
	page := r.URL.Query().Get(":page")

	if id != "32" {
		t.Errorf("url param set to [%s]; want [%s]", id, "32")
	}
	if page != "12" {
		t.Errorf("url param set to [%s]; want [%s]", page, "12")
	}
}

func TestNotFound(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegistor()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Code set to [%v]; want [%v]", w.Code, http.StatusNotFound)
	}
}

// TestStatic tests the ability to serve static
// content from the filesystem
func TestStatic(t *testing.T) {
	r, _ := http.NewRequest("GET", "/router_test.go", nil)
	w := httptest.NewRecorder()
	pwd, _ := os.Getwd()

	handler := NewControllerRegistor()
	SetStaticPath("/", pwd)
	handler.ServeHTTP(w, r)

	testFile, _ := ioutil.ReadFile(pwd + "/routes_test.go")
	if w.Body.String() != string(testFile) {
		t.Errorf("handler.Static failed to serve file")
	}
}
