package beego

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type errorTestController struct {
	Controller
}

func (ec *errorTestController) Get() {
	errorCode, err := ec.GetInt("code")
	if err != nil {
		ec.Abort("parse code error")
	}
	if errorCode != 0 {
		ec.CustomAbort(errorCode, ec.GetString("code"))
	}
	ec.Abort("404")
}

func TestErrorCode_01(t *testing.T) {
	registerDefaultErrorHandler()
	for k, _ := range ErrorMaps {
		r, _ := http.NewRequest("GET", "/error?code="+k, nil)
		w := httptest.NewRecorder()

		handler := NewControllerRegister()
		handler.Add("/error", &errorTestController{})
		handler.ServeHTTP(w, r)
		code, _ := strconv.Atoi(k)
		if w.Code != code {
			t.Fail()
		}
		if !strings.Contains(string(w.Body.Bytes()), http.StatusText(code)) {
			t.Fail()
		}
	}
}

func TestErrorCode_02(t *testing.T) {
	registerDefaultErrorHandler()
	r, _ := http.NewRequest("GET", "/error?code=0", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/error", &errorTestController{})
	handler.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Fail()

	}
}

func TestErrorCode_03(t *testing.T) {
	registerDefaultErrorHandler()
	r, _ := http.NewRequest("GET", "/error?code=crash", nil)
	w := httptest.NewRecorder()

	handler := NewControllerRegister()
	handler.Add("/error", &errorTestController{})
	handler.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fail()
	}
	if string(w.Body.Bytes()) != "parse code error" {
		t.Fail()
	}
}
