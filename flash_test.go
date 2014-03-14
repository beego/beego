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
	handler := NewControllerRegistor()
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
