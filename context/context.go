// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Usage:
//
//	import "github.com/astaxie/beego/context"
//
//	ctx := context.Context{Request:req,ResponseWriter:rw}
//
//  more docs http://beego.me/docs/module/context.md
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
	"github.com/astaxie/beego/utils"
)

// Http request context struct including BeegoInput, BeegoOutput, http.Request and http.ResponseWriter.
// BeegoInput and BeegoOutput provides some api to operate request and response more easily.
type Context struct {
	Input          *BeegoInput
	Output         *BeegoOutput
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	_xsrf_token    string
}

// Redirect does redirection to localurl with http header status code.
// It sends http response header directly.
func (ctx *Context) Redirect(status int, localurl string) {
	ctx.Output.Header("Location", localurl)
	ctx.ResponseWriter.WriteHeader(status)
}

// Abort stops this request.
// if middleware.ErrorMaps exists, panic body.
// if middleware.HTTPExceptionMaps exists, panic HTTPException struct with status and body string.
func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
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
	ctx.ResponseWriter.Write([]byte(body))
	panic("User stop run")
}

// Write string to response body.
// it sends response body.
func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
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

// XsrfToken creates a xsrf token string and returns.
func (ctx *Context) XsrfToken(key string, expire int64) string {
	if ctx._xsrf_token == "" {
		token, ok := ctx.GetSecureCookie(key, "_xsrf")
		if !ok {
			token = string(utils.RandomCreateBytes(32))
			ctx.SetSecureCookie(key, "_xsrf", token, expire)
		}
		ctx._xsrf_token = token
	}
	return ctx._xsrf_token
}

// CheckXsrfCookie checks xsrf token in this request is valid or not.
// the token can provided in request header "X-Xsrftoken" and "X-CsrfToken"
// or in form field value named as "_xsrf".
func (ctx *Context) CheckXsrfCookie() bool {
	token := ctx.Input.Query("_xsrf")
	if token == "" {
		token = ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		ctx.Abort(403, "'_xsrf' argument missing from POST")
	} else if ctx._xsrf_token != token {
		ctx.Abort(403, "XSRF cookie does not match POST argument")
	}
	return true
}
