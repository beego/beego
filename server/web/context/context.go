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

// Package context provide the context utils
// Usage:
//
//	import "github.com/beego/beego/v2/server/web/context"
//
//	ctx := context.Context{Request:req,ResponseWriter:rw}
//
package context

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	"github.com/beego/beego/v2/core/utils"
	"github.com/beego/beego/v2/server/web/session"
)

// Commonly used mime-types
const (
	ApplicationJSON  = "application/json"
	ApplicationXML   = "application/xml"
	ApplicationForm  = "application/x-www-form-urlencoded"
	ApplicationProto = "application/x-protobuf"
	ApplicationYAML  = "application/x-yaml"
	TextXML          = "text/xml"

	formatTime      = "15:04:05"
	formatDate      = "2006-01-02"
	formatDateTime  = "2006-01-02 15:04:05"
	formatDateTimeT = "2006-01-02T15:04:05"
)

// NewContext return the Context with Input and Output
func NewContext() *Context {
	return &Context{
		Input:  NewInput(),
		Output: NewOutput(),
	}
}

// Context Http request context struct including BeegoInput, BeegoOutput, http.Request and http.ResponseWriter.
// BeegoInput and BeegoOutput provides an api to operate request and response more easily.
type Context struct {
	Input          *BeegoInput
	Output         *BeegoOutput
	Request        *http.Request
	ResponseWriter *Response
	_xsrfToken     string
}

func (ctx *Context) Bind(obj interface{}) error {
	ct, exist := ctx.Request.Header["Content-Type"]
	if !exist || len(ct) == 0 {
		return ctx.BindJSON(obj)
	}
	i, l := 0, len(ct[0])
	for i < l && ct[0][i] != ';' {
		i++
	}
	switch ct[0][0:i] {
	case ApplicationJSON:
		return ctx.BindJSON(obj)
	case ApplicationXML, TextXML:
		return ctx.BindXML(obj)
	case ApplicationForm:
		return ctx.BindForm(obj)
	case ApplicationProto:
		return ctx.BindProtobuf(obj.(proto.Message))
	case ApplicationYAML:
		return ctx.BindYAML(obj)
	default:
		return errors.New("Unsupported Content-Type:" + ct[0])
	}
}

// Resp sends response based on the Accept Header
// By default response will be in JSON
func (ctx *Context) Resp(data interface{}) error {
	accept := ctx.Input.Header("Accept")
	switch accept {
	case ApplicationYAML:
		return ctx.YamlResp(data)
	case ApplicationXML, TextXML:
		return ctx.XMLResp(data)
	case ApplicationProto:
		return ctx.ProtoResp(data.(proto.Message))
	default:
		return ctx.JSONResp(data)
	}
}

func (ctx *Context) JSONResp(data interface{}) error {
	return ctx.Output.JSON(data, false, false)
}

func (ctx *Context) XMLResp(data interface{}) error {
	return ctx.Output.XML(data, false)
}

func (ctx *Context) YamlResp(data interface{}) error {
	return ctx.Output.YAML(data)
}

func (ctx *Context) ProtoResp(data proto.Message) error {
	return ctx.Output.Proto(data)
}

// BindYAML only read data from http request body
func (ctx *Context) BindYAML(obj interface{}) error {
	return yaml.Unmarshal(ctx.Input.RequestBody, obj)
}

// BindForm will parse form values to struct via tag.
func (ctx *Context) BindForm(obj interface{}) error {
	err := ctx.Request.ParseForm()
	if err != nil {
		return err
	}
	return ParseForm(ctx.Request.Form, obj)
}

// BindJSON only read data from http request body
func (ctx *Context) BindJSON(obj interface{}) error {
	return json.Unmarshal(ctx.Input.RequestBody, obj)
}

// BindProtobuf only read data from http request body
func (ctx *Context) BindProtobuf(obj proto.Message) error {
	return proto.Unmarshal(ctx.Input.RequestBody, obj)
}

// BindXML only read data from http request body
func (ctx *Context) BindXML(obj interface{}) error {
	return xml.Unmarshal(ctx.Input.RequestBody, obj)
}

// ParseForm will parse form values to struct via tag.
func ParseForm(form url.Values, obj interface{}) error {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if !isStructPtr(objT) {
		return fmt.Errorf("%v must be  a struct pointer", obj)
	}
	objT = objT.Elem()
	objV = objV.Elem()
	return parseFormToStruct(form, objT, objV)
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// Reset initializes Context, BeegoInput and BeegoOutput
func (ctx *Context) Reset(rw http.ResponseWriter, r *http.Request) {
	ctx.Request = r
	if ctx.ResponseWriter == nil {
		ctx.ResponseWriter = &Response{}
	}
	ctx.ResponseWriter.reset(rw)
	ctx.Input.Reset(ctx)
	ctx.Output.Reset(ctx)
	ctx._xsrfToken = ""
}

// Redirect redirects to localurl with http header status code.
func (ctx *Context) Redirect(status int, localurl string) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, localurl, status)
}

// Abort stops the request.
// If beego.ErrorMaps exists, panic body.
func (ctx *Context) Abort(status int, body string) {
	ctx.Output.SetStatus(status)
	panic(body)
}

// WriteString writes a string to response body.
func (ctx *Context) WriteString(content string) {
	_, _ = ctx.ResponseWriter.Write([]byte(content))
}

// GetCookie gets a cookie from a request for a given key.
// (Alias of BeegoInput.Cookie)
func (ctx *Context) GetCookie(key string) string {
	return ctx.Input.Cookie(key)
}

// SetCookie sets a cookie for a response.
// (Alias of BeegoOutput.Cookie)
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	ctx.Output.Cookie(name, value, others...)
}

// GetSecureCookie gets a secure cookie from a request for a given key.
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

	h := hmac.New(sha256.New, []byte(Secret))
	_, err := fmt.Fprintf(h, "%s%s", vs, timestamp)
	if err != nil {
		return "", false
	}

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// SetSecureCookie sets a secure cookie for a response.
func (ctx *Context) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha256.New, []byte(Secret))
	_, _ = fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.Output.Cookie(name, cookie, others...)
}

// XSRFToken creates and returns an xsrf token string
func (ctx *Context) XSRFToken(key string, expire int64) string {
	if ctx._xsrfToken == "" {
		token, ok := ctx.GetSecureCookie(key, "_xsrf")
		if !ok {
			token = string(utils.RandomCreateBytes(32))
			// TODO make it configurable
			ctx.SetSecureCookie(key, "_xsrf", token, expire, "/", "", true, true)
		}
		ctx._xsrfToken = token
	}
	return ctx._xsrfToken
}

// CheckXSRFCookie checks if the XSRF token in this request is valid or not.
// The token can be provided in the request header in the form "X-Xsrftoken" or "X-CsrfToken"
// or in form field value named as "_xsrf".
func (ctx *Context) CheckXSRFCookie() bool {
	token := ctx.Input.Query("_xsrf")
	if token == "" {
		token = ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		ctx.Abort(422, "422")
		return false
	}
	if ctx._xsrfToken != token {
		ctx.Abort(417, "417")
		return false
	}
	return true
}

// RenderMethodResult renders the return value of a controller method to the output
func (ctx *Context) RenderMethodResult(result interface{}) {
	if result != nil {
		renderer, ok := result.(Renderer)
		if !ok {
			err, ok := result.(error)
			if ok {
				renderer = errorRenderer(err)
			} else {
				renderer = jsonRenderer(result)
			}
		}
		renderer.Render(ctx)
	}
}

// Session return session store of this context of request
func (ctx *Context) Session() (store session.Store, err error) {
	if ctx.Input != nil {
		if ctx.Input.CruSession != nil {
			store = ctx.Input.CruSession
			return
		} else {
			err = errors.New(`no valid session store(please initialize session)`)
			return
		}
	} else {
		err = errors.New(`no valid input`)
		return
	}
}

// Response is a wrapper for the http.ResponseWriter
// Started:  if true, response was already written to so the other handler will not be executed
type Response struct {
	http.ResponseWriter
	Started bool
	Status  int
	Elapsed time.Duration
}

func (r *Response) reset(rw http.ResponseWriter) {
	r.ResponseWriter = rw
	r.Status = 0
	r.Started = false
}

// Write writes the data to the connection as part of a HTTP reply,
// and sets `Started` to true.
// Started:  if true, the response was already sent
func (r *Response) Write(p []byte) (int, error) {
	r.Started = true
	return r.ResponseWriter.Write(p)
}

// WriteHeader sends a HTTP response header with status code,
// and sets `Started` to true.
func (r *Response) WriteHeader(code int) {
	if r.Status > 0 {
		// prevent multiple response.WriteHeader calls
		return
	}
	r.Status = code
	r.Started = true
	r.ResponseWriter.WriteHeader(code)
}

// Hijack hijacker for http
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

// Flush http.Flusher
func (r *Response) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// CloseNotify http.CloseNotifier
func (r *Response) CloseNotify() <-chan bool {
	if cn, ok := r.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

// Pusher http.Pusher
func (r *Response) Pusher() (pusher http.Pusher) {
	if pusher, ok := r.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
