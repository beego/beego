// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie
package beego

import (
	"net/http"
	"strings"

	beecontext "github.com/astaxie/beego/context"
)

type namespaceCond func(*beecontext.Context) bool

// Namespace is store all the info
type Namespace struct {
	prefix    string
	condition namespaceCond
	handlers  *ControllerRegistor
}

// get new Namespace
func NewNamespace(prefix string) *Namespace {
	cr := NewControllerRegistor()
	return &Namespace{
		prefix:   prefix,
		handlers: cr,
	}
}

// set condtion function
// if cond return true can run this namespace, else can't
// usage:
// ns.Cond(func (ctx *context.Context) bool{
//       if ctx.Input.Domain() == "api.beego.me" {
//         return true
//       }
//       return false
//   })
func (n *Namespace) Cond(cond namespaceCond) *Namespace {
	n.condition = cond
	return n
}

// add filter in the Namespace
// action has before & after
// FilterFunc
// usage:
// Filter("before", func (ctx *context.Context){
//       _, ok := ctx.Input.Session("uid").(int)
//       if !ok && ctx.Request.RequestURI != "/login" {
//          ctx.Redirect(302, "/login")
//        }
//   })
func (n *Namespace) Filter(action string, filter FilterFunc) *Namespace {
	if action == "before" {
		action = "BeforeRouter"
	} else if action == "after" {
		action = "FinishRouter"
	}
	n.handlers.AddFilter("*", action, filter)
	return n
}

// same as beego.Rourer
// refer: https://godoc.org/github.com/astaxie/beego#Router
func (n *Namespace) Router(rootpath string, c ControllerInterface, mappingMethods ...string) *Namespace {
	n.handlers.Add(rootpath, c, mappingMethods...)
	return n
}

// same as beego.AutoRouter
// refer: https://godoc.org/github.com/astaxie/beego#AutoRouter
func (n *Namespace) AutoRouter(c ControllerInterface) *Namespace {
	n.handlers.AddAuto(c)
	return n
}

// same as beego.AutoPrefix
// refer: https://godoc.org/github.com/astaxie/beego#AutoPrefix
func (n *Namespace) AutoPrefix(prefix string, c ControllerInterface) *Namespace {
	n.handlers.AddAutoPrefix(prefix, c)
	return n
}

// same as beego.Get
// refer: https://godoc.org/github.com/astaxie/beego#Get
func (n *Namespace) Get(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Get(rootpath, f)
	return n
}

// same as beego.Post
// refer: https://godoc.org/github.com/astaxie/beego#Post
func (n *Namespace) Post(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Post(rootpath, f)
	return n
}

// same as beego.Delete
// refer: https://godoc.org/github.com/astaxie/beego#Delete
func (n *Namespace) Delete(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Delete(rootpath, f)
	return n
}

// same as beego.Put
// refer: https://godoc.org/github.com/astaxie/beego#Put
func (n *Namespace) Put(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Put(rootpath, f)
	return n
}

// same as beego.Head
// refer: https://godoc.org/github.com/astaxie/beego#Head
func (n *Namespace) Head(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Head(rootpath, f)
	return n
}

// same as beego.Options
// refer: https://godoc.org/github.com/astaxie/beego#Options
func (n *Namespace) Options(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Options(rootpath, f)
	return n
}

// same as beego.Patch
// refer: https://godoc.org/github.com/astaxie/beego#Patch
func (n *Namespace) Patch(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Patch(rootpath, f)
	return n
}

// same as beego.Any
// refer: https://godoc.org/github.com/astaxie/beego#Any
func (n *Namespace) Any(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Any(rootpath, f)
	return n
}

// same as beego.Handler
// refer: https://godoc.org/github.com/astaxie/beego#Handler
func (n *Namespace) Handler(rootpath string, h http.Handler) *Namespace {
	n.handlers.Handler(rootpath, h)
	return n
}

// nest Namespace
// usage:
//ns := beego.NewNamespace(“/v1”).
//Namespace(
//    beego.NewNamespace("/shop").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("shopinfo"))
//    }),
//    beego.NewNamespace("/order").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("orderinfo"))
//    }),
//    beego.NewNamespace("/crm").
//        Get("/:id", func(ctx *context.Context) {
//            ctx.Output.Body([]byte("crminfo"))
//    }),
//)
func (n *Namespace) Namespace(ns ...*Namespace) *Namespace {
	for _, ni := range ns {
		n.handlers.Handler(ni.prefix, ni, true)
	}
	return n
}

// Namespace implement the http.Handler
func (n *Namespace) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//trim the preifix from URL.Path
	r.URL.Path = strings.TrimPrefix(r.URL.Path, n.prefix)
	// init context
	context := &beecontext.Context{
		ResponseWriter: rw,
		Request:        r,
		Input:          beecontext.NewInput(r),
		Output:         beecontext.NewOutput(),
	}
	context.Output.Context = context
	context.Output.EnableGzip = EnableGzip

	if context.Input.IsWebsocket() {
		context.ResponseWriter = rw
	}
	if n.condition != nil && !n.condition(context) {
		http.Error(rw, "Method Not Allowed", 405)
		return
	}
	n.handlers.ServeHTTP(rw, r)
}

// register Namespace into beego.Handler
// support multi Namespace
func AddNamespace(nl ...*Namespace) {
	for _, n := range nl {
		Handler(n.prefix, n, true)
	}
}
