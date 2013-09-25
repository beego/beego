package context

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	Input          *BeegoInput
	Output         *BeegoOutput
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	SecureKey      string
}

func (ctx *Context) Redirect(status int, localurl string) {
	ctx.Output.Header("Location", localurl)
	ctx.Output.SetStatus(status)
}

func (ctx *Context) Abort(status int, body string) {
	ctx.Output.SetStatus(status)
	ctx.Output.Body([]byte(body))
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

//code is migrated from [Go web](https://github.com/hoisie/web/blob/master/web.go)
//all rights are reserved by @hoisie
func getCookieSig(key string, val []byte, timestamp string) string {
	hm := hmac.New(sha1.New, []byte(key))

	hm.Write(val)
	hm.Write([]byte(timestamp))

	hex := fmt.Sprintf("%02x", hm.Sum(nil))
	return hex
}

//code is migrated from [Go web](https://github.com/hoisie/web/blob/master/web.go)
//all rights are reserved by @hoisie
func (ctx *Context) GetSecureCookie(key string) string {
	if cookie := ctx.GetCookie(key); cookie != "" {
		parts := strings.SplitN(cookie, "|", 3)

		val := parts[0]
		timestamp := parts[1]
		sig := parts[2]

		if getCookieSig(ctx.SecureKey, []byte(val), timestamp) != sig {
			return ""
		}

		ts, _ := strconv.ParseInt(timestamp, 0, 64)

		if time.Now().Unix()-31*86400 > ts {
			return ""
		}

		buf := bytes.NewBufferString(val)
		encoder := base64.NewDecoder(base64.StdEncoding, buf)

		res, _ := ioutil.ReadAll(encoder)
		return string(res)
	}

	return ""
}

//code is migrated from [Go web](https://github.com/hoisie/web/blob/master/web.go)
//all rights are reserved by @hoisie
func (ctx *Context) SetSecureCookie(name string, val string, others ...interface{}) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(val))
	encoder.Close()
	vs := buf.String()
	vb := buf.Bytes()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := getCookieSig(ctx.SecureKey, vb, timestamp)
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")

	ctx.Output.Cookie(name, cookie, others...)
}
