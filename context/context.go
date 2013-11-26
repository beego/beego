package context

import (
	"github.com/astaxie/beego/middleware"
	"net/http"
)

type Context struct {
	Input          *BeegoInput
	Output         *BeegoOutput
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func (ctx *Context) Redirect(status int, localurl string) {
	ctx.Output.Header("Location", localurl)
	ctx.Output.SetStatus(status)
}

func (ctx *Context) Abort(status int, body string) {
	ctx.Output.SetStatus(status)
	ctx.Output.Body([]byte(body))

	if e, ok := middleware.HTTPExceptionMaps[status]; ok {
		if len(body) >= 1 {
			e.Description = body
		}
		panic(e)
	}
}

func (ctx *Context) WriteString(content string) {
	ctx.Output.Body([]byte(content))
}

func (ctx *Context) GetCookie(key string) string {
	return ctx.Input.Cookie(key)
}

func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	ctx.Output.Cookie(name, value, others...)
}
