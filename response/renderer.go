package response

import (
	"strconv"

	beecontext "github.com/astaxie/beego/context"
)

type Renderer interface {
	Render(ctx *beecontext.Context)
}

type rendererFunc func(ctx *beecontext.Context)

func (f rendererFunc) Render(ctx *beecontext.Context) {
	f(ctx)
}

type StatusCode int

func (s StatusCode) Error() string {
	return strconv.Itoa(int(s))
}

func (s StatusCode) Render(ctx *beecontext.Context) {
	ctx.Output.SetStatus(int(s))
}
