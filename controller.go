package beego

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/session"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	Ctx         *Context
	Data        map[interface{}]interface{}
	ChildName   string
	TplNames    string
	Layout      string
	TplExt      string
	_xsrf_token string
	gotofunc    string
	CruSession  session.SessionStore
	XSRFExpire  int
}

type ControllerInterface interface {
	Init(ct *Context, cn string)
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

func (c *Controller) Init(ctx *Context, cn string) {
	c.Data = make(map[interface{}]interface{})
	c.Layout = ""
	c.TplNames = ""
	c.ChildName = cn
	c.Ctx = ctx
	c.TplExt = "tpl"
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
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.writeToWriter(rb)
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
		subdir := path.Dir(c.TplNames)
		_, file := path.Split(c.TplNames)
		newbytes := bytes.NewBufferString("")
		if _, ok := BeeTemplates[subdir]; !ok {
			panic("can't find templatefile in the path:" + c.TplNames)
			return []byte{}, errors.New("can't find templatefile in the path:" + c.TplNames)
		}
		BeeTemplates[subdir].ExecuteTemplate(newbytes, file, c.Data)
		tplcontent, _ := ioutil.ReadAll(newbytes)
		c.Data["LayoutContent"] = template.HTML(string(tplcontent))
		subdir = path.Dir(c.Layout)
		_, file = path.Split(c.Layout)
		ibytes := bytes.NewBufferString("")
		err := BeeTemplates[subdir].ExecuteTemplate(ibytes, file, c.Data)
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
		subdir := path.Dir(c.TplNames)
		_, file := path.Split(c.TplNames)
		ibytes := bytes.NewBufferString("")
		if _, ok := BeeTemplates[subdir]; !ok {
			panic("can't find templatefile in the path:" + c.TplNames)
			return []byte{}, errors.New("can't find templatefile in the path:" + c.TplNames)
		}
		err := BeeTemplates[subdir].ExecuteTemplate(ibytes, file, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
		}
		icontent, _ := ioutil.ReadAll(ibytes)
		return icontent, nil
	}
	return []byte{}, nil
}

func (c *Controller) writeToWriter(rb []byte) {
	output_writer := c.Ctx.ResponseWriter.(io.Writer)
	if EnableGzip == true && c.Ctx.Request.Header.Get("Accept-Encoding") != "" {
		splitted := strings.SplitN(c.Ctx.Request.Header.Get("Accept-Encoding"), ",", -1)
		encodings := make([]string, len(splitted))

		for i, val := range splitted {
			encodings[i] = strings.TrimSpace(val)
		}
		for _, val := range encodings {
			if val == "gzip" {
				c.Ctx.ResponseWriter.Header().Set("Content-Encoding", "gzip")
				output_writer, _ = gzip.NewWriterLevel(c.Ctx.ResponseWriter, gzip.BestSpeed)

				break
			} else if val == "deflate" {
				c.Ctx.ResponseWriter.Header().Set("Content-Encoding", "deflate")
				output_writer, _ = flate.NewWriter(c.Ctx.ResponseWriter, flate.BestSpeed)
				break
			}
		}
	} else {
		c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(rb)), true)
	}
	output_writer.Write(rb)
	switch output_writer.(type) {
	case *gzip.Writer:
		output_writer.(*gzip.Writer).Close()
	case *flate.Writer:
		output_writer.(*flate.Writer).Close()
	case io.WriteCloser:
		output_writer.(io.WriteCloser).Close()
	}
}

func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(code, url)
}

func (c *Controller) Abort(code string) {
	panic(code)
}

func (c *Controller) ServeJson(encoding ...bool) {
	var content []byte
	var err error
	if RunMode == "prod" {
		content, err = json.Marshal(c.Data["json"])
	} else {
		content, err = json.MarshalIndent(c.Data["json"], "", "  ")
	}
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if len(encoding) > 0 && encoding[0] == true {
		content = []byte(stringsToJson(string(content)))
	}
	c.writeToWriter(content)
}

func (c *Controller) ServeJsonp() {
	var content []byte
	var err error
	if RunMode == "prod" {
		content, err = json.Marshal(c.Data["jsonp"])
	} else {
		content, err = json.MarshalIndent(c.Data["jsonp"], "", "  ")
	}
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	callback := c.Ctx.Request.Form.Get("callback")
	if callback == "" {
		http.Error(c.Ctx.ResponseWriter, `"callback" parameter required`, http.StatusInternalServerError)
		return
	}
	callback_content := bytes.NewBufferString(callback)
	callback_content.WriteString("(")
	callback_content.Write(content)
	callback_content.WriteString(");\r\n")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	c.writeToWriter(callback_content.Bytes())
}

func (c *Controller) ServeXml() {
	var content []byte
	var err error
	if RunMode == "prod" {
		content, err = xml.Marshal(c.Data["xml"])
	} else {
		content, err = xml.MarshalIndent(c.Data["xml"], "", "  ")
	}
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	c.writeToWriter(content)
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
		c.CruSession = GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
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

func (c *Controller) DestroySession() {
	GlobalSessions.SessionDestroy(c.Ctx.ResponseWriter, c.Ctx.Request)
}

func (c *Controller) IsAjax() bool {
	return (c.Ctx.Request.Header.Get("HTTP_X_REQUESTED_WITH") == "XMLHttpRequest")
}

func (c *Controller) XsrfToken() string {
	if c._xsrf_token == "" {
		token := c.Ctx.GetCookie("_xsrf")
		if token == "" {
			h := hmac.New(sha1.New, []byte(XSRFKEY))
			fmt.Fprintf(h, "%s:%d", c.Ctx.Request.RemoteAddr, time.Now().UnixNano())
			tok := fmt.Sprintf("%s:%d", h.Sum(nil), time.Now().UnixNano())
			token = base64.URLEncoding.EncodeToString([]byte(tok))
			expire := 0
			if c.XSRFExpire > 0 {
				expire = c.XSRFExpire
			} else {
				expire = XSRFExpire
			}
			c.Ctx.SetCookie("_xsrf", token, expire)
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
	}

	if c._xsrf_token != token {
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
