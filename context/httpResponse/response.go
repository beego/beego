package httpResponse

import (
	"strconv"

	"net/http"

	beecontext "github.com/astaxie/beego/context"
)

const (
	//BadRequest indicates http error 400
	BadRequest StatusCode = http.StatusBadRequest

	//NotFound indicates http error 404
	NotFound StatusCode = http.StatusNotFound
)

// Redirect renders http 302 with a URL
func Redirect(localurl string) error {
	return statusCodeWithRender{302, func(ctx *beecontext.Context) {
		ctx.Redirect(302, localurl)
	}}
}

// StatusCode sets the http response status code
type StatusCode int

func (s StatusCode) Error() string {
	return strconv.Itoa(int(s))
}

// Render sets the http status code
func (s StatusCode) Render(ctx *beecontext.Context) {
	ctx.Output.SetStatus(int(s))
}

type statusCodeWithRender struct {
	statusCode int
	f          func(ctx *beecontext.Context)
}

//assert that statusCodeWithRender implements Renderer interface
var _r beecontext.Renderer = (*statusCodeWithRender)(nil)

func (s statusCodeWithRender) Error() string {
	return strconv.Itoa(s.statusCode)
}

func (s statusCodeWithRender) Render(ctx *beecontext.Context) {
	s.f(ctx)
}
