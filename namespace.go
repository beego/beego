// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie
package beego

import (
	"net/http"

	beecontext "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/middleware"
)

type namespaceCond func(*beecontext.Context) bool

// Namespace is store all the info
type Namespace struct {
	prefix   string
	handlers *ControllerRegistor
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
// Cond as the first filter
func (n *Namespace) Cond(cond namespaceCond) *Namespace {
	fn := func(ctx *beecontext.Context) {
		if !cond(ctx) {
			middleware.Exception("405", ctx.ResponseWriter, ctx.Request, "Method not allowed")
		}
	}
	if v, ok := n.handlers.filters[BeforeRouter]; ok {
		mr := new(FilterRouter)
		mr.tree = NewTree()
		mr.pattern = "*"
		mr.filterFunc = fn
		mr.tree.AddRouter("*", true)
		n.handlers.filters[BeforeRouter] = append([]*FilterRouter{mr}, v...)
	} else {
		n.handlers.InsertFilter("*", BeforeRouter, fn)
	}
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
	var a int
	if action == "before" {
		a = BeforeRouter
	} else if action == "after" {
		a = FinishRouter
	}
	n.handlers.InsertFilter("*", a, filter)
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
		for k, v := range ni.handlers.routers {
			if t, ok := n.handlers.routers[k]; ok {
				n.handlers.routers[k].AddTree(ni.prefix, v)
			} else {
				t = NewTree()
				t.AddTree(ni.prefix, v)
				n.handlers.routers[k] = t
			}
		}
		if n.handlers.enableFilter {
			for pos, filterList := range ni.handlers.filters {
				for _, mr := range filterList {
					t := NewTree()
					t.AddTree(ni.prefix, mr.tree)
					mr.tree = t
					n.handlers.insertFilterRouter(pos, mr)
				}
			}
		}
	}
	return n
}

// register Namespace into beego.Handler
// support multi Namespace
func AddNamespace(nl ...*Namespace) {
	for _, n := range nl {
		for k, v := range n.handlers.routers {
			if t, ok := BeeApp.Handlers.routers[k]; ok {
				BeeApp.Handlers.routers[k].AddTree(n.prefix, v)
			} else {
				t = NewTree()
				t.AddTree(n.prefix, v)
				BeeApp.Handlers.routers[k] = t
			}
		}
		if n.handlers.enableFilter {
			for pos, filterList := range n.handlers.filters {
				for _, mr := range filterList {
					t := NewTree()
					t.AddTree(n.prefix, mr.tree)
					mr.tree = t
					BeeApp.Handlers.insertFilterRouter(pos, mr)
				}
			}
		}
	}
}
