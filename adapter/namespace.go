// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapter

import (
	"net/http"

	adtContext "github.com/beego/beego/v2/adapter/context"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

type namespaceCond func(*adtContext.Context) bool

// LinkNamespace used as link action
type LinkNamespace func(*Namespace)

// Namespace is store all the info
type Namespace web.Namespace

// NewNamespace get new Namespace
func NewNamespace(prefix string, params ...LinkNamespace) *Namespace {
	nps := oldToNewLinkNs(params)
	return (*Namespace)(web.NewNamespace(prefix, nps...))
}

func oldToNewLinkNs(params []LinkNamespace) []web.LinkNamespace {
	nps := make([]web.LinkNamespace, 0, len(params))
	for i := 0; i < len(params); i++ {
		p := params[i]
		nps = append(nps, func(namespace *web.Namespace) {
			p((*Namespace)(namespace))
		})
	}
	return nps
}

// Cond set condition function
// if cond return true can run this namespace, else can't
// usage:
//
//	ns.Cond(func (ctx *context.Context) bool{
//	      if ctx.Input.Domain() == "api.beego.vip" {
//	        return true
//	      }
//	      return false
//	  })
//
// Cond as the first filter
func (n *Namespace) Cond(cond namespaceCond) *Namespace {
	(*web.Namespace)(n).Cond(func(context *context.Context) bool {
		return cond((*adtContext.Context)(context))
	})
	return n
}

// Filter add filter in the Namespace
// action has before & after
// FilterFunc
// usage:
//
//	Filter("before", func (ctx *context.Context){
//	      _, ok := ctx.Input.Session("uid").(int)
//	      if !ok && ctx.Request.RequestURI != "/login" {
//	         ctx.Redirect(302, "/login")
//	       }
//	  })
func (n *Namespace) Filter(action string, filter ...FilterFunc) *Namespace {
	nfs := oldToNewFilter(filter)
	(*web.Namespace)(n).Filter(action, nfs...)
	return n
}

func oldToNewFilter(filter []FilterFunc) []web.FilterFunc {
	nfs := make([]web.FilterFunc, 0, len(filter))
	for i := 0; i < len(filter); i++ {
		f := filter[i]
		nfs = append(nfs, func(ctx *context.Context) {
			f((*adtContext.Context)(ctx))
		})
	}
	return nfs
}

// Router same as beego.Rourer
// refer: https://godoc.org/github.com/beego/beego/v2#Router
func (n *Namespace) Router(rootpath string, c ControllerInterface, mappingMethods ...string) *Namespace {
	(*web.Namespace)(n).Router(rootpath, c, mappingMethods...)
	return n
}

// AutoRouter same as beego.AutoRouter
// refer: https://godoc.org/github.com/beego/beego/v2#AutoRouter
func (n *Namespace) AutoRouter(c ControllerInterface) *Namespace {
	(*web.Namespace)(n).AutoRouter(c)
	return n
}

// AutoPrefix same as beego.AutoPrefix
// refer: https://godoc.org/github.com/beego/beego/v2#AutoPrefix
func (n *Namespace) AutoPrefix(prefix string, c ControllerInterface) *Namespace {
	(*web.Namespace)(n).AutoPrefix(prefix, c)
	return n
}

// Get same as beego.Get
// refer: https://godoc.org/github.com/beego/beego/v2#Get
func (n *Namespace) Get(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Get(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Post same as beego.Post
// refer: https://godoc.org/github.com/beego/beego/v2#Post
func (n *Namespace) Post(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Post(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Delete same as beego.Delete
// refer: https://godoc.org/github.com/beego/beego/v2#Delete
func (n *Namespace) Delete(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Delete(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Put same as beego.Put
// refer: https://godoc.org/github.com/beego/beego/v2#Put
func (n *Namespace) Put(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Put(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Head same as beego.Head
// refer: https://godoc.org/github.com/beego/beego/v2#Head
func (n *Namespace) Head(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Head(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Options same as beego.Options
// refer: https://godoc.org/github.com/beego/beego/v2#Options
func (n *Namespace) Options(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Options(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Patch same as beego.Patch
// refer: https://godoc.org/github.com/beego/beego/v2#Patch
func (n *Namespace) Patch(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Patch(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Any same as beego.Any
// refer: https://godoc.org/github.com/beego/beego/v2#Any
func (n *Namespace) Any(rootpath string, f FilterFunc) *Namespace {
	(*web.Namespace)(n).Any(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return n
}

// Handler same as beego.Handler
// refer: https://godoc.org/github.com/beego/beego/v2#Handler
func (n *Namespace) Handler(rootpath string, h http.Handler) *Namespace {
	(*web.Namespace)(n).Handler(rootpath, h)
	return n
}

// Include add include class
// refer: https://godoc.org/github.com/beego/beego/v2#Include
func (n *Namespace) Include(cList ...ControllerInterface) *Namespace {
	nL := oldToNewCtrlIntfs(cList)
	(*web.Namespace)(n).Include(nL...)
	return n
}

// Namespace add nest Namespace
// usage:
// ns := beego.NewNamespace(“/v1”).
// Namespace(
//
//	beego.NewNamespace("/shop").
//	    Get("/:id", func(ctx *context.Context) {
//	        ctx.Output.Body([]byte("shopinfo"))
//	}),
//	beego.NewNamespace("/order").
//	    Get("/:id", func(ctx *context.Context) {
//	        ctx.Output.Body([]byte("orderinfo"))
//	}),
//	beego.NewNamespace("/crm").
//	    Get("/:id", func(ctx *context.Context) {
//	        ctx.Output.Body([]byte("crminfo"))
//	}),
//
// )
func (n *Namespace) Namespace(ns ...*Namespace) *Namespace {
	nns := oldToNewNs(ns)
	(*web.Namespace)(n).Namespace(nns...)
	return n
}

func oldToNewNs(ns []*Namespace) []*web.Namespace {
	nns := make([]*web.Namespace, 0, len(ns))
	for _, n := range ns {
		nns = append(nns, (*web.Namespace)(n))
	}
	return nns
}

// AddNamespace register Namespace into beego.Handler
// support multi Namespace
func AddNamespace(nl ...*Namespace) {
	nnl := oldToNewNs(nl)
	web.AddNamespace(nnl...)
}

// NSCond is Namespace Condition
func NSCond(cond namespaceCond) LinkNamespace {
	wc := web.NSCond(func(b *context.Context) bool {
		return cond((*adtContext.Context)(b))
	})
	return func(namespace *Namespace) {
		wc((*web.Namespace)(namespace))
	}
}

// NSBefore Namespace BeforeRouter filter
func NSBefore(filterList ...FilterFunc) LinkNamespace {
	nfs := oldToNewFilter(filterList)
	wf := web.NSBefore(nfs...)
	return func(namespace *Namespace) {
		wf((*web.Namespace)(namespace))
	}
}

// NSAfter add Namespace FinishRouter filter
func NSAfter(filterList ...FilterFunc) LinkNamespace {
	nfs := oldToNewFilter(filterList)
	wf := web.NSAfter(nfs...)
	return func(namespace *Namespace) {
		wf((*web.Namespace)(namespace))
	}
}

// NSInclude Namespace Include ControllerInterface
func NSInclude(cList ...ControllerInterface) LinkNamespace {
	nfs := oldToNewCtrlIntfs(cList)
	wi := web.NSInclude(nfs...)
	return func(namespace *Namespace) {
		wi((*web.Namespace)(namespace))
	}
}

// NSRouter call Namespace Router
func NSRouter(rootpath string, c ControllerInterface, mappingMethods ...string) LinkNamespace {
	wn := web.NSRouter(rootpath, c, mappingMethods...)
	return func(namespace *Namespace) {
		wn((*web.Namespace)(namespace))
	}
}

// NSGet call Namespace Get
func NSGet(rootpath string, f FilterFunc) LinkNamespace {
	ln := web.NSGet(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		ln((*web.Namespace)(ns))
	}
}

// NSPost call Namespace Post
func NSPost(rootpath string, f FilterFunc) LinkNamespace {
	wp := web.NSPost(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		wp((*web.Namespace)(ns))
	}
}

// NSHead call Namespace Head
func NSHead(rootpath string, f FilterFunc) LinkNamespace {
	wb := web.NSHead(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		wb((*web.Namespace)(ns))
	}
}

// NSPut call Namespace Put
func NSPut(rootpath string, f FilterFunc) LinkNamespace {
	wn := web.NSPut(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}

// NSDelete call Namespace Delete
func NSDelete(rootpath string, f FilterFunc) LinkNamespace {
	wn := web.NSDelete(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}

// NSAny call Namespace Any
func NSAny(rootpath string, f FilterFunc) LinkNamespace {
	wn := web.NSAny(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}

// NSOptions call Namespace Options
func NSOptions(rootpath string, f FilterFunc) LinkNamespace {
	wo := web.NSOptions(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		wo((*web.Namespace)(ns))
	}
}

// NSPatch call Namespace Patch
func NSPatch(rootpath string, f FilterFunc) LinkNamespace {
	wn := web.NSPatch(rootpath, func(ctx *context.Context) {
		f((*adtContext.Context)(ctx))
	})
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}

// NSAutoRouter call Namespace AutoRouter
func NSAutoRouter(c ControllerInterface) LinkNamespace {
	wn := web.NSAutoRouter(c)
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}

// NSAutoPrefix call Namespace AutoPrefix
func NSAutoPrefix(prefix string, c ControllerInterface) LinkNamespace {
	wn := web.NSAutoPrefix(prefix, c)
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}

// NSNamespace add sub Namespace
func NSNamespace(prefix string, params ...LinkNamespace) LinkNamespace {
	nps := oldToNewLinkNs(params)
	wn := web.NSNamespace(prefix, nps...)
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}

// NSHandler add handler
func NSHandler(rootpath string, h http.Handler) LinkNamespace {
	wn := web.NSHandler(rootpath, h)
	return func(ns *Namespace) {
		wn((*web.Namespace)(ns))
	}
}
