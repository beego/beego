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
