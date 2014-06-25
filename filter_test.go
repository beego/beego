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

var FilterUser = func(ctx *context.Context) {
	ctx.Output.Body([]byte("i am " + ctx.Input.Params[":last"] + ctx.Input.Params[":first"]))
}

func TestFilter(t *testing.T) {
	r, _ := http.NewRequest("GET", "/person/asta/Xie", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/person/:last/:first", BeforeRouter, FilterUser)
	handler.Add("/person/:last/:first", &TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am astaXie" {
		t.Errorf("user define func can't run")
	}
}

var FilterAdminUser = func(ctx *context.Context) {
	ctx.Output.Body([]byte("i am admin"))
}

// Filter pattern /admin/:all
// all url like    /admin/    /admin/xie    will all get filter

func TestPatternTwo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/admin/", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/admin/?:all", BeforeRouter, FilterAdminUser)
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am admin" {
		t.Errorf("filter /admin/ can't run")
	}
}

func TestPatternThree(t *testing.T) {
	r, _ := http.NewRequest("GET", "/admin/astaxie", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/admin/:all", BeforeRouter, FilterAdminUser)
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am admin" {
		t.Errorf("filter /admin/astaxie can't run")
	}
}
