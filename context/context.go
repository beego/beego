package context

import (
	"net/http"

	"github.com/astaxie/beego/middleware"
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
	// first panic from ErrorMaps, is is user defined error functions.
	if _, ok := middleware.ErrorMaps[body]; ok {
		panic(body)
	}
	// second panic from HTTPExceptionMaps, it is system defined functions.
	if e, ok := middleware.HTTPExceptionMaps[status]; ok {
		if len(body) >= 1 {
			e.Description = body
		}
		panic(e)
	}
	// last panic user string
	panic(body)
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
