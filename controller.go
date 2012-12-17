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
	Ct        *Context
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

func (c *Controller) Init(ct *Context, cn string) {
	c.Data = make(map[interface{}]interface{})
	c.Tpl = template.New(cn + ct.Request.Method)
	c.Tpl = c.Tpl.Funcs(beegoTplFuncMap)
	c.Layout = ""
	c.TplNames = ""
	c.ChildName = cn
	c.Ct = ct
	c.TplExt = "tpl"

}

func (c *Controller) Prepare() {

}

func (c *Controller) Finish() {

}

func (c *Controller) Get() {
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Delete() {
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Put() {
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Head() {
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Patch() {
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Options() {
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Render() error {
	//if the controller has set layout, then first get the tplname's content set the content to the layout
	if c.Layout != "" {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + c.Ct.Request.Method + "." + c.TplExt
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
		err = t.ExecuteTemplate(c.Ct.ResponseWriter, file, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
		}
	} else {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + c.Ct.Request.Method + "." + c.TplExt
		}
		t, err := c.Tpl.ParseFiles(path.Join(ViewsPath, c.TplNames))
		if err != nil {
			Trace("template ParseFiles err:", err)
		}
		_, file := path.Split(c.TplNames)
		err = t.ExecuteTemplate(c.Ct.ResponseWriter, file, c.Data)
		if err != nil {
			Trace("template Execute err:", err)
		}
	}
	return nil
}

func (c *Controller) Redirect(url string, code int) {
	c.Ct.Redirect(code, url)
}

func (c *Controller) ServeJson() {
	content, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		http.Error(c.Ct.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ct.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ct.ContentType("json")
	c.Ct.ResponseWriter.Write(content)
}

func (c *Controller) ServeXml() {
	content, err := xml.Marshal(c.Data)
	if err != nil {
		http.Error(c.Ct.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ct.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ct.ContentType("xml")
	c.Ct.ResponseWriter.Write(content)
}

func (c *Controller) Input() url.Values {
	c.Ct.Request.ParseForm()
	return c.Ct.Request.Form
}
