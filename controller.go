package beego

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
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
	"time"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

var (
	// custom error when user stop request handler manually.
	USERSTOPRUN = errors.New("User stop run")
)

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
}

// Init generates default values of controller operations.
func (c *Controller) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Data = make(map[interface{}]interface{})
	c.Layout = ""
	c.TplNames = ""
	c.controllerName = controllerName
	c.actionName = actionName
	c.Ctx = ctx
	c.TplExt = "tpl"
	c.AppController = app
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

// Render sends the response with rendered template bytes as text/html type.
func (c *Controller) Render() error {
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

// RenderBytes returns the bytes of renderd tempate string. Do not send out response.
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
			return []byte{}, errors.New("can't find templatefile in the path:" + c.TplNames)
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
				if (sectionTpl == "") {
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
			return []byte{}, errors.New("can't find templatefile in the path:" + c.TplNames)
		}
		err := BeeTemplates[c.TplNames].ExecuteTemplate(ibytes, c.TplNames, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
			return nil, err
		}
		icontent, _ := ioutil.ReadAll(ibytes)
		return icontent, nil
	}
	return []byte{}, nil
}

// Redirect sends the redirection response to url with status code.
func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(code, url)
}

// Aborts stops controller handler and show the error data if code is defined in ErrorMap or code string.
func (c *Controller) Abort(code string) {
	status, err := strconv.Atoi(code)
	if err == nil {
		c.Ctx.Abort(status, code)
	} else {
		c.Ctx.Abort(200, code)
	}
}

// StopRun makes panic of USERSTOPRUN error and go to recover function if defined.
func (c *Controller) StopRun() {
	panic(USERSTOPRUN)
}

// UrlFor does another controller handler in this request function.
// it goes to this controller method if endpoint is not clear.
func (c *Controller) UrlFor(endpoint string, values ...string) string {
	if len(endpoint) <= 0 {
		return ""
	}
	if endpoint[0] == '.' {
		return UrlFor(reflect.Indirect(reflect.ValueOf(c.AppController)).Type().Name()+endpoint, values...)
	} else {
		return UrlFor(endpoint, values...)
	}
	return ""
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

// Input returns the input data map from POST or PUT request body and query string.
func (c *Controller) Input() url.Values {
	ct := c.Ctx.Request.Header.Get("Content-Type")
	if strings.Contains(ct, "multipart/form-data") {
		c.Ctx.Request.ParseMultipartForm(MaxMemory) //64MB
	} else {
		c.Ctx.Request.ParseForm()
	}
	return c.Ctx.Request.Form
}

// ParseForm maps input data map to obj struct.
func (c *Controller) ParseForm(obj interface{}) error {
	return ParseForm(c.Input(), obj)
}

// GetString returns the input value by key string.
func (c *Controller) GetString(key string) string {
	return c.Input().Get(key)
}

// GetStrings returns the input string slice by key string.
// it's designed for multi-value input field such as checkbox(input[type=checkbox]), multi-selection.
func (c *Controller) GetStrings(key string) []string {
	r := c.Ctx.Request
	if r.Form == nil {
		return []string{}
	}
	vs := r.Form[key]
	if len(vs) > 0 {
		return vs
	}
	return []string{}
}

// GetInt returns input value as int64.
func (c *Controller) GetInt(key string) (int64, error) {
	return strconv.ParseInt(c.Input().Get(key), 10, 64)
}

// GetBool returns input value as bool.
func (c *Controller) GetBool(key string) (bool, error) {
	return strconv.ParseBool(c.Input().Get(key))
}

// GetFloat returns input value as float64.
func (c *Controller) GetFloat(key string) (float64, error) {
	return strconv.ParseFloat(c.Input().Get(key), 64)
}

// GetFile returns the file data in file upload field named as key.
// it returns the first one of multi-uploaded files.
func (c *Controller) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Ctx.Request.FormFile(key)
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
	c.CruSession = GlobalSessions.SessionRegenerateId(c.Ctx.ResponseWriter, c.Ctx.Request)
	c.Ctx.Input.CruSession = c.CruSession
}

// DestroySession cleans session data and session cookie.
func (c *Controller) DestroySession() {
	GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
}

// IsAjax returns this request is ajax or not.
func (c *Controller) IsAjax() bool {
	return c.Ctx.Input.IsAjax()
}

// GetSecureCookie returns decoded cookie value from encoded browser cookie values.
func (c *Controller) GetSecureCookie(Secret, key string) (string, bool) {
	val := c.Ctx.GetCookie(key)
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

// SetSecureCookie puts value into cookie after encoded the value.
func (c *Controller) SetSecureCookie(Secret, name, val string, age int64) {
	vs := base64.URLEncoding.EncodeToString([]byte(val))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	c.Ctx.SetCookie(name, cookie, age, "/")
}

// XsrfToken creates a xsrf token string and returns.
func (c *Controller) XsrfToken() string {
	if c._xsrf_token == "" {
		token, ok := c.GetSecureCookie(XSRFKEY, "_xsrf")
		if !ok {
			var expire int64
			if c.XSRFExpire > 0 {
				expire = int64(c.XSRFExpire)
			} else {
				expire = int64(XSRFExpire)
			}
			token = getRandomString(15)
			c.SetSecureCookie(XSRFKEY, "_xsrf", token, expire)
		}
		c._xsrf_token = token
	}
	return c._xsrf_token
}

// CheckXsrfCookie checks xsrf token in this request is valid or not.
// the token can provided in request header "X-Xsrftoken" and "X-CsrfToken"
// or in form field value named as "_xsrf".
func (c *Controller) CheckXsrfCookie() bool {
	token := c.GetString("_xsrf")
	if token == "" {
		token = c.Ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = c.Ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		c.Ctx.Abort(403, "'_xsrf' argument missing from POST")
	} else if c._xsrf_token != token {
		c.Ctx.Abort(403, "XSRF cookie does not match POST argument")
	}
	return true
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

// getRandomString returns random string.
func getRandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
