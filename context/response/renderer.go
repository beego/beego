package response

import (
	"strconv"

	beecontext "github.com/astaxie/beego/context"
)

// Renderer defines an http response renderer
type Renderer interface {
	Render(ctx *beecontext.Context)
}

type rendererFunc func(ctx *beecontext.Context)

func (f rendererFunc) Render(ctx *beecontext.Context) {
	f(ctx)
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
	rendererFunc
}

func (s statusCodeWithRender) Error() string {
	return strconv.Itoa(s.statusCode)
}
