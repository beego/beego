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
	handler := NewControllerRegistor()
	handler.AddFilter("/person/:last/:first", "AfterStatic", FilterUser)
	handler.Add("/person/:last/:first", &TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am astaXie" {
		t.Errorf("user define func can't run")
	}
}
