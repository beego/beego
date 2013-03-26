package beego

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/astaxie/session"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Handler struct {
	Ctx       *Context
	Data      map[interface{}]interface{}
	ChildName string
	TplNames  string
	Layout    string
	TplExt    string
}

type HandlerInterface interface {
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

func (c *Handler) Init(ctx *Context, cn string) {
	c.Data = make(map[interface{}]interface{})
	c.Layout = ""
	c.TplNames = ""
	c.ChildName = cn
	c.Ctx = ctx
	c.TplExt = "html"

}

func (c *Handler) Prepare() {

}

func (c *Handler) Finish() {

}

func (c *Handler) Get() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Handler) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Handler) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Handler) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Handler) Head() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Handler) Patch() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Handler) Options() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Handler) SetSession(name string, value interface{}) {
	ss := c.StartSession()
	ss.Set(name, value)
}

func (c *Handler) GetSession(name string) interface{} {
	ss := c.StartSession()
	return ss.Get(name)
}

func (c *Handler) DelSession(name string) {
	ss := c.StartSession()
	ss.Delete(name)
}

func (c *Handler) Render() error {
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

func (c *Handler) RenderString() (string, error) {
	b, e := c.RenderBytes()
	return string(b), e
}

func (c *Handler) RenderBytes() ([]byte, error) {
	//if the handler has set layout, then first get the tplname's content set the content to the layout
	if c.Layout != "" {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt
		}
		_, file := path.Split(c.TplNames)
		subdir := path.Dir(c.TplNames)
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
		_, file := path.Split(c.TplNames)
		subdir := path.Dir(c.TplNames)
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

func (c *Handler) Redirect(url string, code int) {
	c.Ctx.Redirect(code, url)
}

func (c *Handler) ServeJson() {
	content, err := json.MarshalIndent(c.Data["json"], "", "  ")
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("json")
	c.Ctx.ResponseWriter.Write(content)
}

func (c *Handler) ServeXml() {
	content, err := xml.Marshal(c.Data["xml"])
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("xml")
	c.Ctx.ResponseWriter.Write(content)
}

func (c *Handler) Input() url.Values {
	c.Ctx.Request.ParseForm()
	return c.Ctx.Request.Form
}

func (c *Handler) StartSession() (sess session.Session) {
	sess = GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	return
}
