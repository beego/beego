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

package adapter

import (
	"mime/multipart"
	"net/url"

	"github.com/astaxie/beego/pkg/adapter/session"
	webContext "github.com/astaxie/beego/pkg/server/web/context"

	"github.com/astaxie/beego/pkg/server/web"
)

var (
	// ErrAbort custom error when user stop request handler manually.
	ErrAbort = web.ErrAbort
	// GlobalControllerRouter store comments with controller. pkgpath+controller:comments
	GlobalControllerRouter = web.GlobalControllerRouter
)

// ControllerFilter store the filter for controller
type ControllerFilter web.ControllerFilter

// ControllerFilterComments store the comment for controller level filter
type ControllerFilterComments web.ControllerFilterComments

// ControllerImportComments store the import comment for controller needed
type ControllerImportComments web.ControllerImportComments

// ControllerComments store the comment for the controller method
type ControllerComments web.ControllerComments

// ControllerCommentsSlice implements the sort interface
type ControllerCommentsSlice web.ControllerCommentsSlice

func (p ControllerCommentsSlice) Len() int {
	return (web.ControllerCommentsSlice)(p).Len()
}
func (p ControllerCommentsSlice) Less(i, j int) bool {
	return (web.ControllerCommentsSlice)(p).Less(i, j)
}
func (p ControllerCommentsSlice) Swap(i, j int) {
	(web.ControllerCommentsSlice)(p).Swap(i, j)
}

// Controller defines some basic http request handler operations, such as
// http context, template and view, session and xsrf.
type Controller web.Controller

func (c *Controller) Init(ctx *webContext.Context, controllerName, actionName string, app interface{}) {
	(*web.Controller)(c).Init(ctx, controllerName, actionName, app)
}

// ControllerInterface is an interface to uniform all controller handler.
type ControllerInterface web.ControllerInterface

// Prepare runs after Init before request function execution.
func (c *Controller) Prepare() {
	(*web.Controller)(c).Prepare()
}

// Finish runs after request function execution.
func (c *Controller) Finish() {
	(*web.Controller)(c).Finish()
}

// Get adds a request function to handle GET request.
func (c *Controller) Get() {
	(*web.Controller)(c).Get()
}

// Post adds a request function to handle POST request.
func (c *Controller) Post() {
	(*web.Controller)(c).Post()
}

// Delete adds a request function to handle DELETE request.
func (c *Controller) Delete() {
	(*web.Controller)(c).Delete()
}

// Put adds a request function to handle PUT request.
func (c *Controller) Put() {
	(*web.Controller)(c).Put()
}

// Head adds a request function to handle HEAD request.
func (c *Controller) Head() {
	(*web.Controller)(c).Head()
}

// Patch adds a request function to handle PATCH request.
func (c *Controller) Patch() {
	(*web.Controller)(c).Patch()
}

// Options adds a request function to handle OPTIONS request.
func (c *Controller) Options() {
	(*web.Controller)(c).Options()
}

// Trace adds a request function to handle Trace request.
// this method SHOULD NOT be overridden.
// https://tools.ietf.org/html/rfc7231#section-4.3.8
// The TRACE method requests a remote, application-level loop-back of
// the request message.  The final recipient of the request SHOULD
// reflect the message received, excluding some fields described below,
// back to the client as the message body of a 200 (OK) response with a
// Content-Type of "message/http" (Section 8.3.1 of [RFC7230]).
func (c *Controller) Trace() {
	(*web.Controller)(c).Trace()
}

// HandlerFunc call function with the name
func (c *Controller) HandlerFunc(fnname string) bool {
	return (*web.Controller)(c).HandlerFunc(fnname)
}

// URLMapping register the internal Controller router.
func (c *Controller) URLMapping() {
	(*web.Controller)(c).URLMapping()
}

// Mapping the method to function
func (c *Controller) Mapping(method string, fn func()) {
	(*web.Controller)(c).Mapping(method, fn)
}

// Render sends the response with rendered template bytes as text/html type.
func (c *Controller) Render() error {
	return (*web.Controller)(c).Render()
}

// RenderString returns the rendered template string. Do not send out response.
func (c *Controller) RenderString() (string, error) {
	return (*web.Controller)(c).RenderString()
}

// RenderBytes returns the bytes of rendered template string. Do not send out response.
func (c *Controller) RenderBytes() ([]byte, error) {
	return (*web.Controller)(c).RenderBytes()
}

// Redirect sends the redirection response to url with status code.
func (c *Controller) Redirect(url string, code int) {
	(*web.Controller)(c).Redirect(url, code)
}

// SetData set the data depending on the accepted
func (c *Controller) SetData(data interface{}) {
	(*web.Controller)(c).SetData(data)
}

// Abort stops controller handler and show the error data if code is defined in ErrorMap or code string.
func (c *Controller) Abort(code string) {
	(*web.Controller)(c).Abort(code)
}

// CustomAbort stops controller handler and show the error data, it's similar Aborts, but support status code and body.
func (c *Controller) CustomAbort(status int, body string) {
	(*web.Controller)(c).CustomAbort(status, body)
}

// StopRun makes panic of USERSTOPRUN error and go to recover function if defined.
func (c *Controller) StopRun() {
	(*web.Controller)(c).StopRun()
}

// URLFor does another controller handler in this request function.
// it goes to this controller method if endpoint is not clear.
func (c *Controller) URLFor(endpoint string, values ...interface{}) string {
	return (*web.Controller)(c).URLFor(endpoint, values...)
}

// ServeJSON sends a json response with encoding charset.
func (c *Controller) ServeJSON(encoding ...bool) {
	(*web.Controller)(c).ServeJSON(encoding...)
}

// ServeJSONP sends a jsonp response.
func (c *Controller) ServeJSONP() {
	(*web.Controller)(c).ServeJSONP()
}

// ServeXML sends xml response.
func (c *Controller) ServeXML() {
	(*web.Controller)(c).ServeXML()
}

// ServeYAML sends yaml response.
func (c *Controller) ServeYAML() {
	(*web.Controller)(c).ServeYAML()
}

// ServeFormatted serve YAML, XML OR JSON, depending on the value of the Accept header
func (c *Controller) ServeFormatted(encoding ...bool) {
	(*web.Controller)(c).ServeFormatted(encoding...)
}

// Input returns the input data map from POST or PUT request body and query string.
func (c *Controller) Input() url.Values {
	return (*web.Controller)(c).Input()
}

// ParseForm maps input data map to obj struct.
func (c *Controller) ParseForm(obj interface{}) error {
	return (*web.Controller)(c).ParseForm(obj)
}

// GetString returns the input value by key string or the default value while it's present and input is blank
func (c *Controller) GetString(key string, def ...string) string {
	return (*web.Controller)(c).GetString(key, def...)
}

// GetStrings returns the input string slice by key string or the default value while it's present and input is blank
// it's designed for multi-value input field such as checkbox(input[type=checkbox]), multi-selection.
func (c *Controller) GetStrings(key string, def ...[]string) []string {
	return (*web.Controller)(c).GetStrings(key, def...)
}

// GetInt returns input as an int or the default value while it's present and input is blank
func (c *Controller) GetInt(key string, def ...int) (int, error) {
	return (*web.Controller)(c).GetInt(key, def...)
}

// GetInt8 return input as an int8 or the default value while it's present and input is blank
func (c *Controller) GetInt8(key string, def ...int8) (int8, error) {
	return (*web.Controller)(c).GetInt8(key, def...)
}

// GetUint8 return input as an uint8 or the default value while it's present and input is blank
func (c *Controller) GetUint8(key string, def ...uint8) (uint8, error) {
	return (*web.Controller)(c).GetUint8(key, def...)
}

// GetInt16 returns input as an int16 or the default value while it's present and input is blank
func (c *Controller) GetInt16(key string, def ...int16) (int16, error) {
	return (*web.Controller)(c).GetInt16(key, def...)
}

// GetUint16 returns input as an uint16 or the default value while it's present and input is blank
func (c *Controller) GetUint16(key string, def ...uint16) (uint16, error) {
	return (*web.Controller)(c).GetUint16(key, def...)
}

// GetInt32 returns input as an int32 or the default value while it's present and input is blank
func (c *Controller) GetInt32(key string, def ...int32) (int32, error) {
	return (*web.Controller)(c).GetInt32(key, def...)
}

// GetUint32 returns input as an uint32 or the default value while it's present and input is blank
func (c *Controller) GetUint32(key string, def ...uint32) (uint32, error) {
	return (*web.Controller)(c).GetUint32(key, def...)
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func (c *Controller) GetInt64(key string, def ...int64) (int64, error) {
	return (*web.Controller)(c).GetInt64(key, def...)
}

// GetUint64 returns input value as uint64 or the default value while it's present and input is blank.
func (c *Controller) GetUint64(key string, def ...uint64) (uint64, error) {
	return (*web.Controller)(c).GetUint64(key, def...)
}

// GetBool returns input value as bool or the default value while it's present and input is blank.
func (c *Controller) GetBool(key string, def ...bool) (bool, error) {
	return (*web.Controller)(c).GetBool(key, def...)
}

// GetFloat returns input value as float64 or the default value while it's present and input is blank.
func (c *Controller) GetFloat(key string, def ...float64) (float64, error) {
	return (*web.Controller)(c).GetFloat(key, def...)
}

// GetFile returns the file data in file upload field named as key.
// it returns the first one of multi-uploaded files.
func (c *Controller) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return (*web.Controller)(c).GetFile(key)
}

// GetFiles return multi-upload files
// files, err:=c.GetFiles("myfiles")
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
	return (*web.Controller)(c).GetFiles(key)
}

// SaveToFile saves uploaded file to new path.
// it only operates the first one of mutil-upload form file field.
func (c *Controller) SaveToFile(fromfile, tofile string) error {
	return (*web.Controller)(c).SaveToFile(fromfile, tofile)
}

// StartSession starts session and load old session data info this controller.
func (c *Controller) StartSession() session.Store {
	s := (*web.Controller)(c).StartSession()
	return session.CreateNewToOldStoreAdapter(s)
}

// SetSession puts value into session.
func (c *Controller) SetSession(name interface{}, value interface{}) {
	(*web.Controller)(c).SetSession(name, value)
}

// GetSession gets value from session.
func (c *Controller) GetSession(name interface{}) interface{} {
	return (*web.Controller)(c).GetSession(name)
}

// DelSession removes value from session.
func (c *Controller) DelSession(name interface{}) {
	(*web.Controller)(c).DelSession(name)
}

// SessionRegenerateID regenerates session id for this session.
// the session data have no changes.
func (c *Controller) SessionRegenerateID() {
	(*web.Controller)(c).SessionRegenerateID()
}

// DestroySession cleans session data and session cookie.
func (c *Controller) DestroySession() {
	(*web.Controller)(c).DestroySession()
}

// IsAjax returns this request is ajax or not.
func (c *Controller) IsAjax() bool {
	return (*web.Controller)(c).IsAjax()
}

// GetSecureCookie returns decoded cookie value from encoded browser cookie values.
func (c *Controller) GetSecureCookie(Secret, key string) (string, bool) {
	return (*web.Controller)(c).GetSecureCookie(Secret, key)
}

// SetSecureCookie puts value into cookie after encoded the value.
func (c *Controller) SetSecureCookie(Secret, name, value string, others ...interface{}) {
	(*web.Controller)(c).SetSecureCookie(Secret, name, value, others...)
}

// XSRFToken creates a CSRF token string and returns.
func (c *Controller) XSRFToken() string {
	return (*web.Controller)(c).XSRFToken()
}

// CheckXSRFCookie checks xsrf token in this request is valid or not.
// the token can provided in request header "X-Xsrftoken" and "X-CsrfToken"
// or in form field value named as "_xsrf".
func (c *Controller) CheckXSRFCookie() bool {
	return (*web.Controller)(c).CheckXSRFCookie()
}

// XSRFFormHTML writes an input field contains xsrf token value.
func (c *Controller) XSRFFormHTML() string {
	return (*web.Controller)(c).XSRFFormHTML()
}

// GetControllerAndAction gets the executing controller name and action name.
func (c *Controller) GetControllerAndAction() (string, string) {
	return (*web.Controller)(c).GetControllerAndAction()
}
