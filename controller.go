package beego

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/astaxie/beego/session"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type Controller struct {
	Ctx       *Context
	Data      map[interface{}]interface{}
	ChildName string
	TplNames  string
	Layout    string
	TplExt    string
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
		c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(rb)), true)
		c.Ctx.ContentType("text/html")
		c.Ctx.ResponseWriter.Write(rb)
		return nil
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
			c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt
		}
		if RunMode == "dev" {
			BuildTemplate(ViewsPath)
		}
		subdir := path.Dir(c.TplNames)
		_, file := path.Split(c.TplNames)
		newbytes := bytes.NewBufferString("")
		BeeTemplates[subdir].ExecuteTemplate(newbytes, file, c.Data)
		tplcontent, _ := ioutil.ReadAll(newbytes)
		c.Data["LayoutContent"] = template.HTML(string(tplcontent))
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
			c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt
		}
		if RunMode == "dev" {
			BuildTemplate(ViewsPath)
		}
		subdir := path.Dir(c.TplNames)
		_, file := path.Split(c.TplNames)
		ibytes := bytes.NewBufferString("")
		err := BeeTemplates[subdir].ExecuteTemplate(ibytes, file, c.Data)
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

func (c *Controller) ServeJson() {
	content, err := json.MarshalIndent(c.Data["json"], "", "  ")
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("json")
	c.Ctx.ResponseWriter.Write(content)
}

func (c *Controller) ServeXml() {
	content, err := xml.Marshal(c.Data["xml"])
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("xml")
	c.Ctx.ResponseWriter.Write(content)
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

func (c *Controller) GetString(key string) string {
	return c.Input().Get(key)
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

func (c *Controller) StartSession() (sess session.SessionStore) {
	sess = GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	return
}

func (c *Controller) SetSession(name string, value interface{}) {
	ss := c.StartSession()
	defer ss.SessionRelease()
	ss.Set(name, value)
}

func (c *Controller) GetSession(name string) interface{} {
	ss := c.StartSession()
	defer ss.SessionRelease()
	return ss.Get(name)
}

func (c *Controller) DelSession(name string) {
	ss := c.StartSession()
	defer ss.SessionRelease()
	ss.Delete(name)
}
