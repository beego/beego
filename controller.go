package beego

import (
	"bytes"
	"crypto/hmac"
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

type Controller struct {
	Ctx           *context.Context
	Data          map[interface{}]interface{}
	ChildName     string
	TplNames      string
	Layout        string
	TplExt        string
	_xsrf_token   string
	gotofunc      string
	CruSession    session.SessionStore
	XSRFExpire    int
	AppController interface{}
}

type ControllerInterface interface {
	Init(ct *context.Context, childName string, app interface{})
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
}

func (c *Controller) Init(ctx *context.Context, childName string, app interface{}) {
	c.Data = make(map[interface{}]interface{})
	c.Layout = ""
	c.TplNames = ""
	c.ChildName = childName
	c.Ctx = ctx
	c.TplExt = "tpl"
	c.AppController = app
}

func (c *Controller) Prepare() {

}

func (c *Controller) Finish() {

}

func (c *Controller) Destructor() {
	if c.CruSession != nil {
		c.CruSession.SessionRelease()
	}
}

func (c *Controller) Get() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Head() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Patch() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Options() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

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

func (c *Controller) RenderString() (string, error) {
	b, e := c.RenderBytes()
	return string(b), e
}

func (c *Controller) RenderBytes() ([]byte, error) {
	//if the controller has set layout, then first get the tplname's content set the content to the layout
	if c.Layout != "" {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + strings.ToLower(c.Ctx.Request.Method) + "." + c.TplExt
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
		}
		tplcontent, _ := ioutil.ReadAll(newbytes)
		c.Data["LayoutContent"] = template.HTML(string(tplcontent))
		ibytes := bytes.NewBufferString("")
		err = BeeTemplates[c.Layout].ExecuteTemplate(ibytes, c.Layout, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
		}
		icontent, _ := ioutil.ReadAll(ibytes)
		return icontent, nil
	} else {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + strings.ToLower(c.Ctx.Request.Method) + "." + c.TplExt
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
		}
		icontent, _ := ioutil.ReadAll(ibytes)
		return icontent, nil
	}
	return []byte{}, nil
}

func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(code, url)
}

func (c *Controller) Abort(code string) {
	panic(code)
}

func (c *Controller) UrlFor(endpoint string, values ...string) string {
	if len(endpoint) <= 0 {
		return ""
	}
	if endpoint[0] == '.' {
		return UrlFor(reflect.Indirect(reflect.ValueOf(c.AppController)).Type().Name()+endpoint, values...)
	} else {
		return UrlFor(endpoint, values...)
	}
}

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

func (c *Controller) ServeJsonp() {
	var hasIndent bool
	if RunMode == "prod" {
		hasIndent = false
	} else {
		hasIndent = true
	}
	c.Ctx.Output.Jsonp(c.Data["jsonp"], hasIndent)
}

func (c *Controller) ServeXml() {
	var hasIndent bool
	if RunMode == "prod" {
		hasIndent = false
	} else {
		hasIndent = true
	}
	c.Ctx.Output.Xml(c.Data["xml"], hasIndent)
}

func (c *Controller) Input() url.Values {
	ct := c.Ctx.Request.Header.Get("Content-Type")
	if strings.Contains(ct, "multipart/form-data") {
		c.Ctx.Request.ParseMultipartForm(MaxMemory) //64MB
	} else {
		c.Ctx.Request.ParseForm()
	}
	return c.Ctx.Request.Form
}

func (c *Controller) ParseForm(obj interface{}) error {
	return ParseForm(c.Input(), obj)
}

func (c *Controller) GetString(key string) string {
	return c.Input().Get(key)
}

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

func (c *Controller) GetInt(key string) (int64, error) {
	return strconv.ParseInt(c.Input().Get(key), 10, 64)
}

func (c *Controller) GetBool(key string) (bool, error) {
	return strconv.ParseBool(c.Input().Get(key))
}

func (c *Controller) GetFloat(key string) (float64, error) {
	return strconv.ParseFloat(c.Input().Get(key), 64)
}

func (c *Controller) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Ctx.Request.FormFile(key)
}

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

func (c *Controller) StartSession() session.SessionStore {
	if c.CruSession == nil {
		c.CruSession = c.Ctx.Input.CruSession
	}
	return c.CruSession
}

func (c *Controller) SetSession(name interface{}, value interface{}) {
	if c.CruSession == nil {
		c.StartSession()
	}
	c.CruSession.Set(name, value)
}

func (c *Controller) GetSession(name interface{}) interface{} {
	if c.CruSession == nil {
		c.StartSession()
	}
	return c.CruSession.Get(name)
}

func (c *Controller) DelSession(name interface{}) {
	if c.CruSession == nil {
		c.StartSession()
	}
	c.CruSession.Delete(name)
}

func (c *Controller) SessionRegenerateID() {
	c.CruSession = GlobalSessions.SessionRegenerateId(c.Ctx.ResponseWriter, c.Ctx.Request)
	c.Ctx.Input.CruSession = c.CruSession
}

func (c *Controller) DestroySession() {
	GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
}

func (c *Controller) IsAjax() bool {
	return c.Ctx.Input.IsAjax()
}

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

func (c *Controller) SetSecureCookie(Secret, name, val string, age int64) {
	vs := base64.URLEncoding.EncodeToString([]byte(val))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	c.Ctx.SetCookie(name, cookie, age, "/")
}

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
			token = GetRandomString(15)
			c.SetSecureCookie(XSRFKEY, "_xsrf", token, expire)
		}
		c._xsrf_token = token
	}
	return c._xsrf_token
}

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

func (c *Controller) XsrfFormHtml() string {
	return "<input type=\"hidden\" name=\"_xsrf\" value=\"" +
		c._xsrf_token + "\"/>"
}

func (c *Controller) GoToFunc(funcname string) {
	if funcname[0] < 65 || funcname[0] > 90 {
		panic("GoToFunc should exported function")
	}
	c.gotofunc = funcname
}
