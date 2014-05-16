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

type Namespace struct {
	prefix    string
	condition namespaceCond
	handlers  *ControllerRegistor
}

func NewNamespace(prefix string) *Namespace {
	cr := NewControllerRegistor()
	return &Namespace{
		prefix:   prefix,
		handlers: cr,
	}
}

func (n *Namespace) Cond(cond namespaceCond) *Namespace {
	n.condition = cond
	return n
}

func (n *Namespace) Filter(action string, filter FilterFunc) *Namespace {
	if action == "before" {
		action = "BeforeRouter"
	} else if action == "after" {
		action = "FinishRouter"
	}
	n.handlers.AddFilter("*", action, filter)
	return n
}

func (n *Namespace) Router(rootpath string, c ControllerInterface, mappingMethods ...string) *Namespace {
	n.handlers.Add(rootpath, c, mappingMethods...)
	return n
}

func (n *Namespace) AutoRouter(c ControllerInterface) *Namespace {
	n.handlers.AddAuto(c)
	return n
}

func (n *Namespace) AutoPrefix(prefix string, c ControllerInterface) *Namespace {
	n.handlers.AddAutoPrefix(prefix, c)
	return n
}

func (n *Namespace) Get(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Get(rootpath, f)
	return n
}

func (n *Namespace) Post(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Post(rootpath, f)
	return n
}

func (n *Namespace) Delete(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Delete(rootpath, f)
	return n
}

func (n *Namespace) Put(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Put(rootpath, f)
	return n
}

func (n *Namespace) Head(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Head(rootpath, f)
	return n
}

func (n *Namespace) Options(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Options(rootpath, f)
	return n
}

func (n *Namespace) Patch(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Patch(rootpath, f)
	return n
}

func (n *Namespace) Any(rootpath string, f FilterFunc) *Namespace {
	n.handlers.Any(rootpath, f)
	return n
}

func (n *Namespace) Handler(rootpath string, h http.Handler) *Namespace {
	n.handlers.Handler(rootpath, h)
	return n
}

func (n *Namespace) Namespace(ns *Namespace) *Namespace {
	n.handlers.Handler(ns.prefix, ns, true)
	return n
}

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
	}
	n.handlers.ServeHTTP(rw, r)
}

func AddNamespace(nl ...*Namespace) {
	for _, n := range nl {
		Handler(n.prefix, n, true)
	}
}
