package beego

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"strings"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	RequestBody    []byte
	Params         map[string]string
}

func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(body))
}

func (ctx *Context) Redirect(status int, url_ string) {
	ctx.ResponseWriter.Header().Set("Location", url_)
	ctx.ResponseWriter.WriteHeader(status)
}

func (ctx *Context) NotModified() {
	ctx.ResponseWriter.WriteHeader(304)
}

func (ctx *Context) NotFound(message string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(message))
}

//Sets the content type by extension, as defined in the mime package.
//For example, ctx.ContentType("json") sets the content-type to "application/json"
func (ctx *Context) ContentType(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		ctx.ResponseWriter.Header().Set("Content-Type", ctype)
	}
}

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.ResponseWriter.Header().Set(hdr, val)
	} else {
		ctx.ResponseWriter.Header().Add(hdr, val)
	}
}

//Sets a cookie -- duration is the amount of time in seconds. 0 = forever

//params:
//string name
//string value
//int64 expire = 0
//string $path
//string $domain
//bool $secure = false
//bool $httponly = false
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
	if len(others) > 0 {
		switch others[0].(type) {
		case int:
			fmt.Fprintf(&b, "; Max-Age=%d", others[0].(int))
		case int64:
			fmt.Fprintf(&b, "; Max-Age=%d", others[0].(int64))
		case int32:
			fmt.Fprintf(&b, "; Max-Age=%d", others[0].(int32))
		}
	} else {
		fmt.Fprintf(&b, "; Max-Age=0")
	}
	if len(others) > 1 {
		fmt.Fprintf(&b, "; Path=%s", sanitizeValue(others[1].(string)))
	}
	if len(others) > 2 {
		fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(others[2].(string)))
	}
	if len(others) > 3 {
		fmt.Fprintf(&b, "; Secure")
	}
	if len(others) > 4 {
		fmt.Fprintf(&b, "; HttpOnly")
	}
	ctx.SetHeader("Set-Cookie", b.String(), false)
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}

func (ctx *Context) GetCookie(key string) string {
	keycookie, err := ctx.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return keycookie.Value
}
