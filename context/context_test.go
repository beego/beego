package context

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupContext() (*Context, *httptest.ResponseRecorder) {

	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	rw := httptest.NewRecorder()
	ctx := &Context{
		ResponseWriter: rw,
		Request:        r,
		Input:          NewInput(r),
		Output:         NewOutput(),
	}
	ctx.Output.Context = ctx
	return ctx, rw
}

func TestRedirect(t *testing.T) {
	ctx, rw := setupContext()
	ctx.Redirect(302, "localhost")

	if e := 302; e != rw.Code {
		t.Errorf("got: %d want: %d\n", rw.Code, e)
	}
}

func TestWriteString(t *testing.T) {
	ctx, rw := setupContext()
	s := "This is only a test."

	ctx.WriteString(s)
	if e := 200; e != rw.Code {
		t.Errorf("got: %d want: %d\n", rw.Code, e)
	}
	if s != rw.Body.String() {
		t.Errorf("got: %s want: %s\n", rw.Body.String(), s)
	}
}
