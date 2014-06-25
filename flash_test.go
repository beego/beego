// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package beego

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestFlashController struct {
	Controller
}

func (this *TestFlashController) TestWriteFlash() {
	flash := NewFlash()
	flash.Notice("TestFlashString")
	flash.Store(&this.Controller)
	// we choose to serve json because we don't want to load a template html file
	this.ServeJson(true)
}

func TestFlashHeader(t *testing.T) {
	// create fake GET request
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// setup the handler
	handler := NewControllerRegister()
	handler.Add("/", &TestFlashController{}, "get:TestWriteFlash")
	handler.ServeHTTP(w, r)

	// get the Set-Cookie value
	sc := w.Header().Get("Set-Cookie")
	// match for the expected header
	res := strings.Contains(sc, "BEEGO_FLASH=%00notice%23BEEGOFLASH%23TestFlashString%00")
	// validate the assertion
	if res != true {
		t.Errorf("TestFlashHeader() unable to validate flash message")
	}
}
