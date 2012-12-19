package beego

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Controller struct {
	Ctx       *Context
	Tpl       *template.Template
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
	c.Tpl = template.New(cn + ctx.Request.Method)
	c.Tpl = c.Tpl.Funcs(beegoTplFuncMap)
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
	//if the controller has set layout, then first get the tplname's content set the content to the layout
	if c.Layout != "" {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt
		}
		t, err := c.Tpl.ParseFiles(path.Join(ViewsPath, c.TplNames), path.Join(ViewsPath, c.Layout))
		if err != nil {
			Trace("template ParseFiles err:", err)
		}
		_, file := path.Split(c.TplNames)
		newbytes := bytes.NewBufferString("")
		t.ExecuteTemplate(newbytes, file, c.Data)
		tplcontent, _ := ioutil.ReadAll(newbytes)
		c.Data["LayoutContent"] = template.HTML(string(tplcontent))
		_, file = path.Split(c.Layout)
		err = t.ExecuteTemplate(c.Ctx.ResponseWriter, file, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
		}
	} else {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt
		}
		t, err := c.Tpl.ParseFiles(path.Join(ViewsPath, c.TplNames))
		if err != nil {
			Trace("template ParseFiles err:", err)
		}
		_, file := path.Split(c.TplNames)
		err = t.ExecuteTemplate(c.Ctx.ResponseWriter, file, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
		}
	}
	return nil
}

func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(code, url)
}

func (c *Controller) ServeJson() {
	content, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("json")
	c.Ctx.ResponseWriter.Write(content)
}

func (c *Controller) ServeXml() {
	content, err := xml.Marshal(c.Data)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("xml")
	c.Ctx.ResponseWriter.Write(content)
}

func (c *Controller) Input() url.Values {
	c.Ctx.Request.ParseForm()
	return c.Ctx.Request.Form
}
