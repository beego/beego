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

package beego

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

//commonly used mime-types
const (
	applicationJson = "application/json"
	applicationXml  = "application/xml"
	textXml         = "text/xml"
)

var (
	// custom error when user stop request handler manually.
	USERSTOPRUN                                            = errors.New("User stop run")
	GlobalControllerRouter map[string][]ControllerComments = make(map[string][]ControllerComments) //pkgpath+controller:comments
)

// store the comment for the controller method
type ControllerComments struct {
	Method           string
	Router           string
	AllowHTTPMethods []string
	Params           []map[string]string
}

// Controller defines some basic http request handler operations, such as
// http context, template and view, session and xsrf.
type Controller struct {
	Ctx            *context.Context
	Data           map[interface{}]interface{}
	controllerName string
	actionName     string
	TplNames       string
	Layout         string
	LayoutSections map[string]string // the key is the section name and the value is the template name
	TplExt         string
	_xsrf_token    string
	gotofunc       string
	CruSession     session.SessionStore
	XSRFExpire     int
	AppController  interface{}
	EnableRender   bool
	EnableXSRF     bool
	methodMapping  map[string]func() //method:routertree
}

// ControllerInterface is an interface to uniform all controller handler.
type ControllerInterface interface {
	Init(ct *context.Context, controllerName, actionName string, app interface{})
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Finish()
	Render() error
	XsrfToken() string
	CheckXsrfCookie() bool
	HandlerFunc(fn string) bool
	URLMapping()
}

// Init generates default values of controller operations.
func (c *Controller) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Layout = ""
	c.TplNames = ""
	c.controllerName = controllerName
	c.actionName = actionName
	c.Ctx = ctx
	c.TplExt = "tpl"
	c.AppController = app
	c.EnableRender = true
	c.EnableXSRF = true
	c.Data = ctx.Input.Data
	c.methodMapping = make(map[string]func())
}

// Prepare runs after Init before request function execution.
func (c *Controller) Prepare() {

}

// Finish runs after request function execution.
func (c *Controller) Finish() {

}

// Get adds a request function to handle GET request.
func (c *Controller) Get() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Post adds a request function to handle POST request.
func (c *Controller) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Delete adds a request function to handle DELETE request.
func (c *Controller) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Put adds a request function to handle PUT request.
func (c *Controller) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Head adds a request function to handle HEAD request.
func (c *Controller) Head() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Patch adds a request function to handle PATCH request.
func (c *Controller) Patch() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// Options adds a request function to handle OPTIONS request.
func (c *Controller) Options() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

// call function fn
func (c *Controller) HandlerFunc(fnname string) bool {
	if v, ok := c.methodMapping[fnname]; ok {
		v()
		return true
	} else {
		return false
	}
}

// URLMapping register the internal Controller router.
func (c *Controller) URLMapping() {
}

func (c *Controller) Mapping(method string, fn func()) {
	c.methodMapping[method] = fn
}

// Render sends the response with rendered template bytes as text/html type.
func (c *Controller) Render() error {
	if !c.EnableRender {
		return nil
	}
	rb, err := c.RenderBytes()

	if err != nil {
		return err
	} else {
		c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
		c.Ctx.Output.Body(rb)
	}
	return nil
}

// RenderString returns the rendered template string. Do not send out response.
func (c *Controller) RenderString() (string, error) {
	b, e := c.RenderBytes()
	return string(b), e
}

// RenderBytes returns the bytes of rendered template string. Do not send out response.
func (c *Controller) RenderBytes() ([]byte, error) {
	//if the controller has set layout, then first get the tplname's content set the content to the layout
	if c.Layout != "" {
		if c.TplNames == "" {
			c.TplNames = strings.ToLower(c.controllerName) + "/" + strings.ToLower(c.actionName) + "." + c.TplExt
		}
		if RunMode == "dev" {
			BuildTemplate(ViewsPath)
		}
		newbytes := bytes.NewBufferString("")
		if _, ok := BeeTemplates[c.TplNames]; !ok {
			panic("can't find templatefile in the path:" + c.TplNames)
		}
		err := BeeTemplates[c.TplNames].ExecuteTemplate(newbytes, c.TplNames, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
			return nil, err
		}
		tplcontent, _ := ioutil.ReadAll(newbytes)
		c.Data["LayoutContent"] = template.HTML(string(tplcontent))

		if c.LayoutSections != nil {
			for sectionName, sectionTpl := range c.LayoutSections {
				if sectionTpl == "" {
					c.Data[sectionName] = ""
					continue
				}

				sectionBytes := bytes.NewBufferString("")
				err = BeeTemplates[sectionTpl].ExecuteTemplate(sectionBytes, sectionTpl, c.Data)
				if err != nil {
					Trace("template Execute err:", err)
					return nil, err
				}
				sectionContent, _ := ioutil.ReadAll(sectionBytes)
				c.Data[sectionName] = template.HTML(string(sectionContent))
			}
		}

		ibytes := bytes.NewBufferString("")
		err = BeeTemplates[c.Layout].ExecuteTemplate(ibytes, c.Layout, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
			return nil, err
		}
		icontent, _ := ioutil.ReadAll(ibytes)
		return icontent, nil
	} else {
		if c.TplNames == "" {
			c.TplNames = strings.ToLower(c.controllerName) + "/" + strings.ToLower(c.actionName) + "." + c.TplExt
		}
		if RunMode == "dev" {
			BuildTemplate(ViewsPath)
		}
		ibytes := bytes.NewBufferString("")
		if _, ok := BeeTemplates[c.TplNames]; !ok {
			panic("can't find templatefile in the path:" + c.TplNames)
		}
		err := BeeTemplates[c.TplNames].ExecuteTemplate(ibytes, c.TplNames, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
			return nil, err
		}
		icontent, _ := ioutil.ReadAll(ibytes)
		return icontent, nil
	}
}

// Redirect sends the redirection response to url with status code.
func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(code, url)
}

// Aborts stops controller handler and show the error data if code is defined in ErrorMap or code string.
func (c *Controller) Abort(code string) {
	status, err := strconv.Atoi(code)
	if err != nil {
		status = 200
	}
	c.CustomAbort(status, code)
}

// CustomAbort stops controller handler and show the error data, it's similar Aborts, but support status code and body.
func (c *Controller) CustomAbort(status int, body string) {
	c.Ctx.ResponseWriter.WriteHeader(status)
	// first panic from ErrorMaps, is is user defined error functions.
	if _, ok := ErrorMaps[body]; ok {
		panic(body)
	}
	// last panic user string
	c.Ctx.ResponseWriter.Write([]byte(body))
	panic(USERSTOPRUN)
}

// StopRun makes panic of USERSTOPRUN error and go to recover function if defined.
func (c *Controller) StopRun() {
	panic(USERSTOPRUN)
}

// UrlFor does another controller handler in this request function.
// it goes to this controller method if endpoint is not clear.
func (c *Controller) UrlFor(endpoint string, values ...interface{}) string {
	if len(endpoint) <= 0 {
		return ""
	}
	if endpoint[0] == '.' {
		return UrlFor(reflect.Indirect(reflect.ValueOf(c.AppController)).Type().Name()+endpoint, values...)
	} else {
		return UrlFor(endpoint, values...)
	}
}

// ServeJson sends a json response with encoding charset.
func (c *Controller) ServeJson(encoding ...bool) {
	var hasIndent bool
	var hasencoding bool
	if RunMode == "prod" {
		hasIndent = false
	} else {
		hasIndent = true
	}
	if len(encoding) > 0 && encoding[0] == true {
		hasencoding = true
	}
	c.Ctx.Output.Json(c.Data["json"], hasIndent, hasencoding)
}

// ServeJsonp sends a jsonp response.
func (c *Controller) ServeJsonp() {
	var hasIndent bool
	if RunMode == "prod" {
		hasIndent = false
	} else {
		hasIndent = true
	}
	c.Ctx.Output.Jsonp(c.Data["jsonp"], hasIndent)
}

// ServeXml sends xml response.
func (c *Controller) ServeXml() {
	var hasIndent bool
	if RunMode == "prod" {
		hasIndent = false
	} else {
		hasIndent = true
	}
	c.Ctx.Output.Xml(c.Data["xml"], hasIndent)
}

// ServeFormatted serve Xml OR Json, depending on the value of the Accept header
func (c *Controller) ServeFormatted() {
	accept := c.Ctx.Input.Header("Accept")
	switch accept {
	case applicationJson:
		c.ServeJson()
	case applicationXml, textXml:
		c.ServeXml()
	default:
		c.ServeJson()
	}
}

// Input returns the input data map from POST or PUT request body and query string.
func (c *Controller) Input() url.Values {
	if c.Ctx.Request.Form == nil {
		c.Ctx.Request.ParseForm()
	}
	return c.Ctx.Request.Form
}

// ParseForm maps input data map to obj struct.
func (c *Controller) ParseForm(obj interface{}) error {
	return ParseForm(c.Input(), obj)
}

// GetString returns the input value by key string or the default value while it's present and input is blank
func (c *Controller) GetString(key string, def ...string) string {
	var defv string
	if len(def) > 0 {
		defv = def[0]
	}

	if v := c.Ctx.Input.Query(key); v != "" {
		return v
	} else {
		return defv
	}
}

// GetStrings returns the input string slice by key string or the default value while it's present and input is blank
// it's designed for multi-value input field such as checkbox(input[type=checkbox]), multi-selection.
func (c *Controller) GetStrings(key string, def ...[]string) []string {
	var defv []string
	if len(def) > 0 {
		defv = def[0]
	}

	f := c.Input()
	if f == nil {
		return defv
	}

	vs := f[key]
	if len(vs) > 0 {
		return vs
	} else {
		return defv
	}
}

// GetInt returns input as an int or the default value while it's present and input is blank
func (c *Controller) GetInt(key string, def ...int) (int, error) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.Atoi(strv)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.Atoi(strv)
	}
}

// GetInt8 return input as an int8 or the default value while it's present and input is blank
func (c *Controller) GetInt8(key string, def ...int8) (int8, error) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		i64, err := strconv.ParseInt(strv, 10, 8)
		i8 := int8(i64)
		return i8, err
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		i64, err := strconv.ParseInt(strv, 10, 8)
		i8 := int8(i64)
		return i8, err
	}
}

// GetInt16 returns input as an int16 or the default value while it's present and input is blank
func (c *Controller) GetInt16(key string, def ...int16) (int16, error) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		i64, err := strconv.ParseInt(strv, 10, 16)
		i16 := int16(i64)
		return i16, err
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		i64, err := strconv.ParseInt(strv, 10, 16)
		i16 := int16(i64)
		return i16, err
	}
}

// GetInt32 returns input as an int32 or the default value while it's present and input is blank
func (c *Controller) GetInt32(key string, def ...int32) (int32, error) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		i64, err := strconv.ParseInt(c.Ctx.Input.Query(key), 10, 32)
		i32 := int32(i64)
		return i32, err
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		i64, err := strconv.ParseInt(c.Ctx.Input.Query(key), 10, 32)
		i32 := int32(i64)
		return i32, err
	}
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func (c *Controller) GetInt64(key string, def ...int64) (int64, error) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.ParseInt(strv, 10, 64)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.ParseInt(strv, 10, 64)
	}
}

// GetBool returns input value as bool or the default value while it's present and input is blank.
func (c *Controller) GetBool(key string, def ...bool) (bool, error) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.ParseBool(strv)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.ParseBool(strv)
	}
}

// GetFloat returns input value as float64 or the default value while it's present and input is blank.
func (c *Controller) GetFloat(key string, def ...float64) (float64, error) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		return strconv.ParseFloat(strv, 64)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.ParseFloat(strv, 64)
	}
}

// GetFile returns the file data in file upload field named as key.
// it returns the first one of multi-uploaded files.
func (c *Controller) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Ctx.Request.FormFile(key)
}

// GetFiles return multi-upload files
// files, err:=c.Getfiles("myfiles")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNoContent)
//		return
//	}
// for i, _ := range files {
//	//for each fileheader, get a handle to the actual file
//	file, err := files[i].Open()
//	defer file.Close()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	//create destination file making sure the path is writeable.
//	dst, err := os.Create("upload/" + files[i].Filename)
//	defer dst.Close()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	//copy the uploaded file to the destination file
//	if _, err := io.Copy(dst, file); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
// }
func (c *Controller) GetFiles(key string) ([]*multipart.FileHeader, error) {
	files, ok := c.Ctx.Request.MultipartForm.File[key]
	if ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

// SaveToFile saves uploaded file to new path.
// it only operates the first one of mutil-upload form file field.
func (c *Controller) SaveToFile(fromfile, tofile string) error {
	file, _, err := c.Ctx.Request.FormFile(fromfile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

// StartSession starts session and load old session data info this controller.
func (c *Controller) StartSession() session.SessionStore {
	if c.CruSession == nil {
		c.CruSession = c.Ctx.Input.CruSession
	}
	return c.CruSession
}

// SetSession puts value into session.
func (c *Controller) SetSession(name interface{}, value interface{}) {
	if c.CruSession == nil {
		c.StartSession()
	}
	c.CruSession.Set(name, value)
}

// GetSession gets value from session.
func (c *Controller) GetSession(name interface{}) interface{} {
	if c.CruSession == nil {
		c.StartSession()
	}
	return c.CruSession.Get(name)
}

// SetSession removes value from session.
func (c *Controller) DelSession(name interface{}) {
	if c.CruSession == nil {
		c.StartSession()
	}
	c.CruSession.Delete(name)
}

// SessionRegenerateID regenerates session id for this session.
// the session data have no changes.
func (c *Controller) SessionRegenerateID() {
	if c.CruSession != nil {
		c.CruSession.SessionRelease(c.Ctx.ResponseWriter)
	}
	c.CruSession = GlobalSessions.SessionRegenerateId(c.Ctx.ResponseWriter, c.Ctx.Request)
	c.Ctx.Input.CruSession = c.CruSession
}

// DestroySession cleans session data and session cookie.
func (c *Controller) DestroySession() {
	c.Ctx.Input.CruSession.Flush()
	GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
}

// IsAjax returns this request is ajax or not.
func (c *Controller) IsAjax() bool {
	return c.Ctx.Input.IsAjax()
}

// GetSecureCookie returns decoded cookie value from encoded browser cookie values.
func (c *Controller) GetSecureCookie(Secret, key string) (string, bool) {
	return c.Ctx.GetSecureCookie(Secret, key)
}

// SetSecureCookie puts value into cookie after encoded the value.
func (c *Controller) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	c.Ctx.SetSecureCookie(Secret, name, value, others...)
}

// XsrfToken creates a xsrf token string and returns.
func (c *Controller) XsrfToken() string {
	if c._xsrf_token == "" {
		var expire int64
		if c.XSRFExpire > 0 {
			expire = int64(c.XSRFExpire)
		} else {
			expire = int64(XSRFExpire)
		}
		c._xsrf_token = c.Ctx.XsrfToken(XSRFKEY, expire)
	}
	return c._xsrf_token
}

// CheckXsrfCookie checks xsrf token in this request is valid or not.
// the token can provided in request header "X-Xsrftoken" and "X-CsrfToken"
// or in form field value named as "_xsrf".
func (c *Controller) CheckXsrfCookie() bool {
	if !c.EnableXSRF {
		return true
	}
	return c.Ctx.CheckXsrfCookie()
}

// XsrfFormHtml writes an input field contains xsrf token value.
func (c *Controller) XsrfFormHtml() string {
	return "<input type=\"hidden\" name=\"_xsrf\" value=\"" +
		c._xsrf_token + "\"/>"
}

// GetControllerAndAction gets the executing controller name and action name.
func (c *Controller) GetControllerAndAction() (controllerName, actionName string) {
	return c.controllerName, c.actionName
}
