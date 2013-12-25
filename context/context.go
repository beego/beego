package context

import (
	"net/http"

	"github.com/astaxie/beego/middleware"
)

// Http request context struct including BeegoInput, BeegoOutput, http.Request and http.ResponseWriter.
// BeegoInput and BeegoOutput provides some api to operate request and response more easily.
type Context struct {
	Input          *BeegoInput
	Output         *BeegoOutput
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// Redirect does redirection to localurl with http header status code.
// It sends http response header directly.
func (ctx *Context) Redirect(status int, localurl string) {
	ctx.Output.Header("Location", localurl)
	ctx.Output.SetStatus(status)
}

// Abort stops this request.
// if middleware.ErrorMaps exists, panic body.
// if middleware.HTTPExceptionMaps exists, panic HTTPException struct with status and body string.
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

// Write string to response body.
// it sends response body.
func (ctx *Context) WriteString(content string) {
	ctx.Output.Body([]byte(content))
}

// Get cookie from request by a given key.
// It's alias of BeegoInput.Cookie.
func (ctx *Context) GetCookie(key string) string {
	return ctx.Input.Cookie(key)
}

// Set cookie for response.
// It's alias of BeegoOutput.Cookie.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	ctx.Output.Cookie(name, value, others...)
}
