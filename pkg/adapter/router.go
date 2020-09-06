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
	"time"

	beecontext "github.com/astaxie/beego/pkg/adapter/context"
	"github.com/astaxie/beego/pkg/server/web/context"

	"github.com/astaxie/beego/pkg/server/web"
)

// default filter execution points
const (
	BeforeStatic = web.BeforeStatic
	BeforeRouter = web.BeforeRouter
	BeforeExec   = web.BeforeExec
	AfterExec    = web.AfterExec
	FinishRouter = web.FinishRouter
)

var (
	// HTTPMETHOD list the supported http methods.
	HTTPMETHOD = web.HTTPMETHOD

	// DefaultAccessLogFilter will skip the accesslog if return true
	DefaultAccessLogFilter FilterHandler = &newToOldFtHdlAdapter{
		delegate: web.DefaultAccessLogFilter,
	}
)

// FilterHandler is an interface for
type FilterHandler interface {
	Filter(*beecontext.Context) bool
}

type newToOldFtHdlAdapter struct {
	delegate web.FilterHandler
}

func (n *newToOldFtHdlAdapter) Filter(ctx *beecontext.Context) bool {
	return n.delegate.Filter((*context.Context)(ctx))
}

// ExceptMethodAppend to append a slice's value into "exceptMethod", for controller's methods shouldn't reflect to AutoRouter
func ExceptMethodAppend(action string) {
	web.ExceptMethodAppend(action)
}

// ControllerInfo holds information about the controller.
type ControllerInfo web.ControllerInfo

func (c *ControllerInfo) GetPattern() string {
	return (*web.ControllerInfo)(c).GetPattern()
}

// ControllerRegister containers registered router rules, controller handlers and filters.
type ControllerRegister web.ControllerRegister

// NewControllerRegister returns a new ControllerRegister.
func NewControllerRegister() *ControllerRegister {
	return (*ControllerRegister)(web.NewControllerRegister())
}

// Add controller handler and pattern rules to ControllerRegister.
// usage:
//	default methods is the same name as method
//	Add("/user",&UserController{})
//	Add("/api/list",&RestController{},"*:ListFood")
//	Add("/api/create",&RestController{},"post:CreateFood")
//	Add("/api/update",&RestController{},"put:UpdateFood")
//	Add("/api/delete",&RestController{},"delete:DeleteFood")
//	Add("/api",&RestController{},"get,post:ApiFunc"
//	Add("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
func (p *ControllerRegister) Add(pattern string, c ControllerInterface, mappingMethods ...string) {
	(*web.ControllerRegister)(p).Add(pattern, c, mappingMethods...)
}

// Include only when the Runmode is dev will generate router file in the router/auto.go from the controller
// Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
func (p *ControllerRegister) Include(cList ...ControllerInterface) {
	nls := oldToNewCtrlIntfs(cList)
	(*web.ControllerRegister)(p).Include(nls...)
}

// GetContext returns a context from pool, so usually you should remember to call Reset function to clean the context
// And don't forget to give back context to pool
// example:
//  ctx := p.GetContext()
//  ctx.Reset(w, q)
//  defer p.GiveBackContext(ctx)
func (p *ControllerRegister) GetContext() *beecontext.Context {
	return (*beecontext.Context)((*web.ControllerRegister)(p).GetContext())
}

// GiveBackContext put the ctx into pool so that it could be reuse
func (p *ControllerRegister) GiveBackContext(ctx *beecontext.Context) {
	(*web.ControllerRegister)(p).GiveBackContext((*context.Context)(ctx))
}

// Get add get method
// usage:
//    Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Get(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Get(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Post add post method
// usage:
//    Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Post(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Post(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Put add put method
// usage:
//    Put("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Put(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Put(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Delete add delete method
// usage:
//    Delete("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Delete(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Delete(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Head add head method
// usage:
//    Head("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Head(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Head(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Patch add patch method
// usage:
//    Patch("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Patch(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Patch(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Options add options method
// usage:
//    Options("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Options(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Options(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Any add all method
// usage:
//    Any("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Any(pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).Any(pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// AddMethod add http method router
// usage:
//    AddMethod("get","/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) AddMethod(method, pattern string, f FilterFunc) {
	(*web.ControllerRegister)(p).AddMethod(method, pattern, func(ctx *context.Context) {
		f((*beecontext.Context)(ctx))
	})
}

// Handler add user defined Handler
func (p *ControllerRegister) Handler(pattern string, h http.Handler, options ...interface{}) {
	(*web.ControllerRegister)(p).Handler(pattern, h, options)
}

// AddAuto router to ControllerRegister.
// example beego.AddAuto(&MainContorlller{}),
// MainController has method List and Page.
// visit the url /main/list to execute List function
// /main/page to execute Page function.
func (p *ControllerRegister) AddAuto(c ControllerInterface) {
	(*web.ControllerRegister)(p).AddAuto(c)
}

// AddAutoPrefix Add auto router to ControllerRegister with prefix.
// example beego.AddAutoPrefix("/admin",&MainContorlller{}),
// MainController has method List and Page.
// visit the url /admin/main/list to execute List function
// /admin/main/page to execute Page function.
func (p *ControllerRegister) AddAutoPrefix(prefix string, c ControllerInterface) {
	(*web.ControllerRegister)(p).AddAutoPrefix(prefix, c)
}

// InsertFilter Add a FilterFunc with pattern rule and action constant.
// params is for:
//   1. setting the returnOnOutput value (false allows multiple filters to execute)
//   2. determining whether or not params need to be reset.
func (p *ControllerRegister) InsertFilter(pattern string, pos int, filter FilterFunc, params ...bool) error {
	opts := oldToNewFilterOpts(params)
	return (*web.ControllerRegister)(p).InsertFilter(pattern, pos, func(ctx *context.Context) {
		filter((*beecontext.Context)(ctx))
	}, opts...)
}

func oldToNewFilterOpts(params []bool) []web.FilterOpt {
	opts := make([]web.FilterOpt, 0, 4)
	if len(params) > 0 {
		opts = append(opts, web.WithReturnOnOutput(params[0]))
	} else {
		// the default value should be true
		opts = append(opts, web.WithReturnOnOutput(true))
	}
	if len(params) > 1 {
		opts = append(opts, web.WithResetParams(params[1]))
	}
	return opts
}

// URLFor does another controller handler in this request function.
// it can access any controller method.
func (p *ControllerRegister) URLFor(endpoint string, values ...interface{}) string {
	return (*web.ControllerRegister)(p).URLFor(endpoint, values...)
}

// Implement http.Handler interface.
func (p *ControllerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	(*web.ControllerRegister)(p).ServeHTTP(rw, r)
}

// FindRouter Find Router info for URL
func (p *ControllerRegister) FindRouter(ctx *beecontext.Context) (routerInfo *ControllerInfo, isFind bool) {
	r, ok := (*web.ControllerRegister)(p).FindRouter((*context.Context)(ctx))
	return (*ControllerInfo)(r), ok
}

// LogAccess logging info HTTP Access
func LogAccess(ctx *beecontext.Context, startTime *time.Time, statusCode int) {
	web.LogAccess((*context.Context)(ctx), startTime, statusCode)
}
