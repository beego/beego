// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package context

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// Get secure cookie from request by a given key.
func (ctx *Context) GetSecureCookie(Secret, key string) (string, bool) {
	val := ctx.Input.Cookie(key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// Set Secure cookie for response.
func (ctx *Context) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.Output.Cookie(name, cookie, others...)
}
