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

package web

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"github.com/beego/beego/v2/core/utils"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/context/param"
)

// default filter execution points
const (
	BeforeStatic = iota
	BeforeRouter
	BeforeExec
	AfterExec
	FinishRouter
)

const (
	routerTypeBeego = iota
	routerTypeRESTFul
	routerTypeHandler
)

var (
	// HTTPMETHOD list the supported http methods.
	HTTPMETHOD = map[string]bool{
		"GET":       true,
		"POST":      true,
		"PUT":       true,
		"DELETE":    true,
		"PATCH":     true,
		"OPTIONS":   true,
		"HEAD":      true,
		"TRACE":     true,
		"CONNECT":   true,
		"MKCOL":     true,
		"COPY":      true,
		"MOVE":      true,
		"PROPFIND":  true,
		"PROPPATCH": true,
		"LOCK":      true,
		"UNLOCK":    true,
	}
	// these web.Controller's methods shouldn't reflect to AutoRouter
	exceptMethod = []string{"Abort", "CheckXSRFCookie", "CustomAbort", "DelSession",
		"DestroySession", "Finish", "GetBool", "GetControllerAndAction",
		"GetFile", "GetFiles", "GetFloat", "GetInt", "GetInt16",
		"GetInt32", "GetInt64", "GetInt8", "GetSecureCookie", "GetSession",
		"GetString", "GetStrings", "GetUint16", "GetUint32", "GetUint64",
		"GetUint8", "HandlerFunc", "Init", "Input",
		"IsAjax", "Mapping", "ParseForm",
		"Prepare", "Redirect", "Render", "RenderBytes",
		"RenderString", "SaveToFile", "SaveToFileWithBuffer", "SaveToFileWithBuffer",
		"ServeFormatted", "ServeJSON", "ServeJSONP", "ServeXML", "ServeYAML",
		"SessionRegenerateID", "SetData", "SetSecureCookie", "SetSession", "StartSession",
		"StopRun", "URLFor", "URLMapping", "XSRFFormHTML",
		"XSRFToken",
	}

	urlPlaceholder = "{{placeholder}}"
	// DefaultAccessLogFilter will skip the accesslog if return true
	DefaultAccessLogFilter FilterHandler = &logFilter{}
)

// FilterHandler is an interface for
type FilterHandler interface {
	Filter(*beecontext.Context) bool
}

// default log filter static file will not show
type logFilter struct {
}

func (l *logFilter) Filter(ctx *beecontext.Context) bool {
	requestPath := path.Clean(ctx.Request.URL.Path)
	if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
		return true
	}
	for prefix := range BConfig.WebConfig.StaticDir {
		if strings.HasPrefix(requestPath, prefix) {
			return true
		}
	}
	return false
}

// ExceptMethodAppend to append a slice's value into "exceptMethod", for controller's methods shouldn't reflect to AutoRouter
func ExceptMethodAppend(action string) {
	exceptMethod = append(exceptMethod, action)
}

// ControllerInfo holds information about the controller.
type ControllerInfo struct {
	pattern        string
	controllerType reflect.Type
	methods        map[string]string
	handler        http.Handler
	runFunction    FilterFunc
	routerType     int
	initialize     func() ControllerInterface
	methodParams   []*param.MethodParam
}

func (c *ControllerInfo) GetPattern() string {
	return c.pattern
}

// ControllerRegister containers registered router rules, controller handlers and filters.
type ControllerRegister struct {
	routers      map[string]*Tree
	enablePolicy bool
	policies     map[string]*Tree
	enableFilter bool
	filters      [FinishRouter + 1][]*FilterRouter
	pool         sync.Pool

	// the filter created by FilterChain
	chainRoot *FilterRouter

	cfg *Config
}

// NewControllerRegister returns a new ControllerRegister.
// Usually you should not use this method
// please use NewControllerRegisterWithCfg
func NewControllerRegister() *ControllerRegister {
	return NewControllerRegisterWithCfg(BeeApp.Cfg)
}

func NewControllerRegisterWithCfg(cfg *Config) *ControllerRegister {
	res := &ControllerRegister{
		routers:  make(map[string]*Tree),
		policies: make(map[string]*Tree),
		pool: sync.Pool{
			New: func() interface{} {
				return beecontext.NewContext()
			},
		},
		cfg: cfg,
	}
	res.chainRoot = newFilterRouter("/*", res.serveHttp, WithCaseSensitive(false))
	return res
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
	p.addWithMethodParams(pattern, c, nil, mappingMethods...)
}

func (p *ControllerRegister) addWithMethodParams(pattern string, c ControllerInterface, methodParams []*param.MethodParam, mappingMethods ...string) {
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	methods := make(map[string]string)
	if len(mappingMethods) > 0 {
		semi := strings.Split(mappingMethods[0], ";")
		for _, v := range semi {
			colon := strings.Split(v, ":")
			if len(colon) != 2 {
				panic("method mapping format is invalid")
			}
			comma := strings.Split(colon[0], ",")
			for _, m := range comma {
				if m == "*" || HTTPMETHOD[strings.ToUpper(m)] {
					if val := reflectVal.MethodByName(colon[1]); val.IsValid() {
						methods[strings.ToUpper(m)] = colon[1]
					} else {
						panic("'" + colon[1] + "' method doesn't exist in the controller " + t.Name())
					}
				} else {
					panic(v + " is an invalid method mapping. Method doesn't exist " + m)
				}
			}
		}
	}

	route := &ControllerInfo{}
	route.pattern = pattern
	route.methods = methods
	route.routerType = routerTypeBeego
	route.controllerType = t
	route.initialize = func() ControllerInterface {
		vc := reflect.New(route.controllerType)
		execController, ok := vc.Interface().(ControllerInterface)
		if !ok {
			panic("controller is not ControllerInterface")
		}

		elemVal := reflect.ValueOf(c).Elem()
		elemType := reflect.TypeOf(c).Elem()
		execElem := reflect.ValueOf(execController).Elem()

		numOfFields := elemVal.NumField()
		for i := 0; i < numOfFields; i++ {
			fieldType := elemType.Field(i)
			elemField := execElem.FieldByName(fieldType.Name)
			if elemField.CanSet() {
				fieldVal := elemVal.Field(i)
				elemField.Set(fieldVal)
			}
		}

		return execController
	}

	route.methodParams = methodParams
	if len(methods) == 0 {
		for m := range HTTPMETHOD {
			p.addToRouter(m, pattern, route)
		}
	} else {
		for k := range methods {
			if k == "*" {
				for m := range HTTPMETHOD {
					p.addToRouter(m, pattern, route)
				}
			} else {
				p.addToRouter(k, pattern, route)
			}
		}
	}
}

func (p *ControllerRegister) addToRouter(method, pattern string, r *ControllerInfo) {
	if !p.cfg.RouterCaseSensitive {
		pattern = strings.ToLower(pattern)
	}
	if t, ok := p.routers[method]; ok {
		t.AddRouter(pattern, r)
	} else {
		t := NewTree()
		t.AddRouter(pattern, r)
		p.routers[method] = t
	}
}

// Include only when the Runmode is dev will generate router file in the router/auto.go from the controller
// Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
func (p *ControllerRegister) Include(cList ...ControllerInterface) {
	for _, c := range cList {
		reflectVal := reflect.ValueOf(c)
		t := reflect.Indirect(reflectVal).Type()
		key := t.PkgPath() + ":" + t.Name()
		if comm, ok := GlobalControllerRouter[key]; ok {
			for _, a := range comm {
				for _, f := range a.Filters {
					p.InsertFilter(f.Pattern, f.Pos, f.Filter, WithReturnOnOutput(f.ReturnOnOutput), WithResetParams(f.ResetParams))
				}
				p.addWithMethodParams(a.Router, c, a.MethodParams, strings.Join(a.AllowHTTPMethods, ",")+":"+a.Method)
			}
		}
	}
}

// GetContext returns a context from pool, so usually you should remember to call Reset function to clean the context
// And don't forget to give back context to pool
// example:
//  ctx := p.GetContext()
//  ctx.Reset(w, q)
//  defer p.GiveBackContext(ctx)
func (p *ControllerRegister) GetContext() *beecontext.Context {
	return p.pool.Get().(*beecontext.Context)
}

// GiveBackContext put the ctx into pool so that it could be reuse
func (p *ControllerRegister) GiveBackContext(ctx *beecontext.Context) {
	p.pool.Put(ctx)
}

// Get add get method
// usage:
//    Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Get(pattern string, f FilterFunc) {
	p.AddMethod("get", pattern, f)
}

// Post add post method
// usage:
//    Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Post(pattern string, f FilterFunc) {
	p.AddMethod("post", pattern, f)
}

// Put add put method
// usage:
//    Put("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Put(pattern string, f FilterFunc) {
	p.AddMethod("put", pattern, f)
}

// Delete add delete method
// usage:
//    Delete("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Delete(pattern string, f FilterFunc) {
	p.AddMethod("delete", pattern, f)
}

// Head add head method
// usage:
//    Head("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Head(pattern string, f FilterFunc) {
	p.AddMethod("head", pattern, f)
}

// Patch add patch method
// usage:
//    Patch("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Patch(pattern string, f FilterFunc) {
	p.AddMethod("patch", pattern, f)
}

// Options add options method
// usage:
//    Options("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Options(pattern string, f FilterFunc) {
	p.AddMethod("options", pattern, f)
}

// Any add all method
// usage:
//    Any("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Any(pattern string, f FilterFunc) {
	p.AddMethod("*", pattern, f)
}

// AddMethod add http method router
// usage:
//    AddMethod("get","/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) AddMethod(method, pattern string, f FilterFunc) {
	method = strings.ToUpper(method)
	if method != "*" && !HTTPMETHOD[method] {
		panic("not support http method: " + method)
	}
	route := &ControllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeRESTFul
	route.runFunction = f
	methods := make(map[string]string)
	if method == "*" {
		for val := range HTTPMETHOD {
			methods[val] = val
		}
	} else {
		methods[method] = method
	}
	route.methods = methods
	for k := range methods {
		if k == "*" {
			for m := range HTTPMETHOD {
				p.addToRouter(m, pattern, route)
			}
		} else {
			p.addToRouter(k, pattern, route)
		}
	}
}

// Handler add user defined Handler
func (p *ControllerRegister) Handler(pattern string, h http.Handler, options ...interface{}) {
	route := &ControllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeHandler
	route.handler = h
	if len(options) > 0 {
		if _, ok := options[0].(bool); ok {
			pattern = path.Join(pattern, "?:all(.*)")
		}
	}
	for m := range HTTPMETHOD {
		p.addToRouter(m, pattern, route)
	}
}

// AddAuto router to ControllerRegister.
// example beego.AddAuto(&MainContorlller{}),
// MainController has method List and Page.
// visit the url /main/list to execute List function
// /main/page to execute Page function.
func (p *ControllerRegister) AddAuto(c ControllerInterface) {
	p.AddAutoPrefix("/", c)
}

// AddAutoPrefix Add auto router to ControllerRegister with prefix.
// example beego.AddAutoPrefix("/admin",&MainContorlller{}),
// MainController has method List and Page.
// visit the url /admin/main/list to execute List function
// /admin/main/page to execute Page function.
func (p *ControllerRegister) AddAutoPrefix(prefix string, c ControllerInterface) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	for i := 0; i < rt.NumMethod(); i++ {
		if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
			route := &ControllerInfo{}
			route.routerType = routerTypeBeego
			route.methods = map[string]string{"*": rt.Method(i).Name}
			route.controllerType = ct
			pattern := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name), "*")
			patternInit := path.Join(prefix, controllerName, rt.Method(i).Name, "*")
			patternFix := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
			patternFixInit := path.Join(prefix, controllerName, rt.Method(i).Name)
			route.pattern = pattern
			for m := range HTTPMETHOD {
				p.addToRouter(m, pattern, route)
				p.addToRouter(m, patternInit, route)
				p.addToRouter(m, patternFix, route)
				p.addToRouter(m, patternFixInit, route)
			}
		}
	}
}

// InsertFilter Add a FilterFunc with pattern rule and action constant.
// params is for:
//   1. setting the returnOnOutput value (false allows multiple filters to execute)
//   2. determining whether or not params need to be reset.
func (p *ControllerRegister) InsertFilter(pattern string, pos int, filter FilterFunc, opts ...FilterOpt) error {
	opts = append(opts, WithCaseSensitive(p.cfg.RouterCaseSensitive))
	mr := newFilterRouter(pattern, filter, opts...)
	return p.insertFilterRouter(pos, mr)
}

// InsertFilterChain is similar to InsertFilter,
// but it will using chainRoot.filterFunc as input to build a new filterFunc
// for example, assume that chainRoot is funcA
// and we add new FilterChain
// fc := func(next) {
//     return func(ctx) {
//           // do something
//           next(ctx)
//           // do something
//     }
// }
func (p *ControllerRegister) InsertFilterChain(pattern string, chain FilterChain, opts ...FilterOpt) {
	root := p.chainRoot
	filterFunc := chain(root.filterFunc)
	opts = append(opts, WithCaseSensitive(p.cfg.RouterCaseSensitive))
	p.chainRoot = newFilterRouter(pattern, filterFunc, opts...)
	p.chainRoot.next = root

}

// add Filter into
func (p *ControllerRegister) insertFilterRouter(pos int, mr *FilterRouter) (err error) {
	if pos < BeforeStatic || pos > FinishRouter {
		return errors.New("can not find your filter position")
	}
	p.enableFilter = true
	p.filters[pos] = append(p.filters[pos], mr)
	return nil
}

// URLFor does another controller handler in this request function.
// it can access any controller method.
func (p *ControllerRegister) URLFor(endpoint string, values ...interface{}) string {
	paths := strings.Split(endpoint, ".")
	if len(paths) <= 1 {
		logs.Warn("urlfor endpoint must like path.controller.method")
		return ""
	}
	if len(values)%2 != 0 {
		logs.Warn("urlfor params must key-value pair")
		return ""
	}
	params := make(map[string]string)
	if len(values) > 0 {
		key := ""
		for k, v := range values {
			if k%2 == 0 {
				key = fmt.Sprint(v)
			} else {
				params[key] = fmt.Sprint(v)
			}
		}
	}
	controllerName := strings.Join(paths[:len(paths)-1], "/")
	methodName := paths[len(paths)-1]
	for m, t := range p.routers {
		ok, url := p.getURL(t, "/", controllerName, methodName, params, m)
		if ok {
			return url
		}
	}
	return ""
}

func (p *ControllerRegister) getURL(t *Tree, url, controllerName, methodName string, params map[string]string, httpMethod string) (bool, string) {
	for _, subtree := range t.fixrouters {
		u := path.Join(url, subtree.prefix)
		ok, u := p.getURL(subtree, u, controllerName, methodName, params, httpMethod)
		if ok {
			return ok, u
		}
	}
	if t.wildcard != nil {
		u := path.Join(url, urlPlaceholder)
		ok, u := p.getURL(t.wildcard, u, controllerName, methodName, params, httpMethod)
		if ok {
			return ok, u
		}
	}
	for _, l := range t.leaves {
		if c, ok := l.runObject.(*ControllerInfo); ok {
			if c.routerType == routerTypeBeego &&
				strings.HasSuffix(path.Join(c.controllerType.PkgPath(), c.controllerType.Name()), controllerName) {
				find := false
				if HTTPMETHOD[strings.ToUpper(methodName)] {
					if len(c.methods) == 0 {
						find = true
					} else if m, ok := c.methods[strings.ToUpper(methodName)]; ok && m == strings.ToUpper(methodName) {
						find = true
					} else if m, ok = c.methods["*"]; ok && m == methodName {
						find = true
					}
				}
				if !find {
					for m, md := range c.methods {
						if (m == "*" || m == httpMethod) && md == methodName {
							find = true
						}
					}
				}
				if find {
					if l.regexps == nil {
						if len(l.wildcards) == 0 {
							return true, strings.Replace(url, "/"+urlPlaceholder, "", 1) + toURL(params)
						}
						if len(l.wildcards) == 1 {
							if v, ok := params[l.wildcards[0]]; ok {
								delete(params, l.wildcards[0])
								return true, strings.Replace(url, urlPlaceholder, v, 1) + toURL(params)
							}
							return false, ""
						}
						if len(l.wildcards) == 3 && l.wildcards[0] == "." {
							if p, ok := params[":path"]; ok {
								if e, isok := params[":ext"]; isok {
									delete(params, ":path")
									delete(params, ":ext")
									return true, strings.Replace(url, urlPlaceholder, p+"."+e, -1) + toURL(params)
								}
							}
						}
						canSkip := false
						for _, v := range l.wildcards {
							if v == ":" {
								canSkip = true
								continue
							}
							if u, ok := params[v]; ok {
								delete(params, v)
								url = strings.Replace(url, urlPlaceholder, u, 1)
							} else {
								if canSkip {
									canSkip = false
									continue
								}
								return false, ""
							}
						}
						return true, url + toURL(params)
					}
					var i int
					var startReg bool
					regURL := ""
					for _, v := range strings.Trim(l.regexps.String(), "^$") {
						if v == '(' {
							startReg = true
							continue
						} else if v == ')' {
							startReg = false
							if v, ok := params[l.wildcards[i]]; ok {
								delete(params, l.wildcards[i])
								regURL = regURL + v
								i++
							} else {
								break
							}
						} else if !startReg {
							regURL = string(append([]rune(regURL), v))
						}
					}
					if l.regexps.MatchString(regURL) {
						ps := strings.Split(regURL, "/")
						for _, p := range ps {
							url = strings.Replace(url, urlPlaceholder, p, 1)
						}
						return true, url + toURL(params)
					}
				}
			}
		}
	}

	return false, ""
}

func (p *ControllerRegister) execFilter(context *beecontext.Context, urlPath string, pos int) (started bool) {
	var preFilterParams map[string]string
	for _, filterR := range p.filters[pos] {
		b, done := filterR.filter(context, urlPath, preFilterParams)
		if done {
			return b
		}
	}
	return false
}

// Implement http.Handler interface.
func (p *ControllerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	ctx := p.GetContext()

	ctx.Reset(rw, r)
	defer p.GiveBackContext(ctx)

	var preFilterParams map[string]string
	p.chainRoot.filter(ctx, p.getUrlPath(ctx), preFilterParams)
}

func (p *ControllerRegister) serveHttp(ctx *beecontext.Context) {
	var err error
	startTime := time.Now()
	r := ctx.Request
	rw := ctx.ResponseWriter.ResponseWriter
	var (
		runRouter    reflect.Type
		findRouter   bool
		runMethod    string
		methodParams []*param.MethodParam
		routerInfo   *ControllerInfo
		isRunnable   bool
	)

	if p.cfg.RecoverFunc != nil {
		defer p.cfg.RecoverFunc(ctx, p.cfg)
	}

	ctx.Output.EnableGzip = p.cfg.EnableGzip

	if p.cfg.RunMode == DEV {
		ctx.Output.Header("Server", p.cfg.ServerName)
	}

	urlPath := p.getUrlPath(ctx)

	// filter wrong http method
	if !HTTPMETHOD[r.Method] {
		exception("405", ctx)
		goto Admin
	}

	// filter for static file
	if len(p.filters[BeforeStatic]) > 0 && p.execFilter(ctx, urlPath, BeforeStatic) {
		goto Admin
	}

	serverStaticRouter(ctx)

	if ctx.ResponseWriter.Started {
		findRouter = true
		goto Admin
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {

		if ctx.Input.IsUpload() {
			ctx.Input.Context.Request.Body = http.MaxBytesReader(ctx.Input.Context.ResponseWriter,
				ctx.Input.Context.Request.Body,
				p.cfg.MaxUploadSize)
		} else if p.cfg.CopyRequestBody {
			// connection will close if the incoming data are larger (RFC 7231, 6.5.11)
			if r.ContentLength > p.cfg.MaxMemory {
				logs.Error(errors.New("payload too large"))
				exception("413", ctx)
				goto Admin
			}
			ctx.Input.CopyBody(p.cfg.MaxMemory)
		} else {
			ctx.Input.Context.Request.Body = http.MaxBytesReader(ctx.Input.Context.ResponseWriter,
				ctx.Input.Context.Request.Body,
				p.cfg.MaxMemory)
		}

		err = ctx.Input.ParseFormOrMultiForm(p.cfg.MaxMemory)
		if err != nil {
			logs.Error(err)
			if strings.Contains(err.Error(), `http: request body too large`) {
				exception("413", ctx)
			} else {
				exception("500", ctx)
			}
			goto Admin
		}
	}

	// session init
	if p.cfg.WebConfig.Session.SessionOn {
		ctx.Input.CruSession, err = GlobalSessions.SessionStart(rw, r)
		if err != nil {
			logs.Error(err)
			exception("503", ctx)
			goto Admin
		}
		defer func() {
			if ctx.Input.CruSession != nil {
				ctx.Input.CruSession.SessionRelease(nil, rw)
			}
		}()
	}
	if len(p.filters[BeforeRouter]) > 0 && p.execFilter(ctx, urlPath, BeforeRouter) {
		goto Admin
	}
	// User can define RunController and RunMethod in filter
	if ctx.Input.RunController != nil && ctx.Input.RunMethod != "" {
		findRouter = true
		runMethod = ctx.Input.RunMethod
		runRouter = ctx.Input.RunController
	} else {
		routerInfo, findRouter = p.FindRouter(ctx)
	}

	// if no matches to url, throw a not found exception
	if !findRouter {
		exception("404", ctx)
		goto Admin
	}
	if splat := ctx.Input.Param(":splat"); splat != "" {
		for k, v := range strings.Split(splat, "/") {
			ctx.Input.SetParam(strconv.Itoa(k), v)
		}
	}

	if routerInfo != nil {
		// store router pattern into context
		ctx.Input.SetData("RouterPattern", routerInfo.pattern)
	}

	// execute middleware filters
	if len(p.filters[BeforeExec]) > 0 && p.execFilter(ctx, urlPath, BeforeExec) {
		goto Admin
	}

	// check policies
	if p.execPolicy(ctx, urlPath) {
		goto Admin
	}

	if routerInfo != nil {
		if routerInfo.routerType == routerTypeRESTFul {
			if _, ok := routerInfo.methods[r.Method]; ok {
				isRunnable = true
				routerInfo.runFunction(ctx)
			} else {
				exception("405", ctx)
				goto Admin
			}
		} else if routerInfo.routerType == routerTypeHandler {
			isRunnable = true
			routerInfo.handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
		} else {
			runRouter = routerInfo.controllerType
			methodParams = routerInfo.methodParams
			method := r.Method
			if r.Method == http.MethodPost && ctx.Input.Query("_method") == http.MethodPut {
				method = http.MethodPut
			}
			if r.Method == http.MethodPost && ctx.Input.Query("_method") == http.MethodDelete {
				method = http.MethodDelete
			}
			if m, ok := routerInfo.methods[method]; ok {
				runMethod = m
			} else if m, ok = routerInfo.methods["*"]; ok {
				runMethod = m
			} else {
				runMethod = method
			}
		}
	}

	// also defined runRouter & runMethod from filter
	if !isRunnable {
		// Invoke the request handler
		var execController ControllerInterface
		if routerInfo != nil && routerInfo.initialize != nil {
			execController = routerInfo.initialize()
		} else {
			vc := reflect.New(runRouter)
			var ok bool
			execController, ok = vc.Interface().(ControllerInterface)
			if !ok {
				panic("controller is not ControllerInterface")
			}
		}

		// call the controller init function
		execController.Init(ctx, runRouter.Name(), runMethod, execController)

		// call prepare function
		execController.Prepare()

		// if XSRF is Enable then check cookie where there has any cookie in the  request's cookie _csrf
		if p.cfg.WebConfig.EnableXSRF {
			execController.XSRFToken()
			if r.Method == http.MethodPost || r.Method == http.MethodDelete || r.Method == http.MethodPut ||
				(r.Method == http.MethodPost && (ctx.Input.Query("_method") == http.MethodDelete || ctx.Input.Query("_method") == http.MethodPut)) {
				execController.CheckXSRFCookie()
			}
		}

		execController.URLMapping()

		if !ctx.ResponseWriter.Started {
			// exec main logic
			switch runMethod {
			case http.MethodGet:
				execController.Get()
			case http.MethodPost:
				execController.Post()
			case http.MethodDelete:
				execController.Delete()
			case http.MethodPut:
				execController.Put()
			case http.MethodHead:
				execController.Head()
			case http.MethodPatch:
				execController.Patch()
			case http.MethodOptions:
				execController.Options()
			case http.MethodTrace:
				execController.Trace()
			default:
				if !execController.HandlerFunc(runMethod) {
					vc := reflect.ValueOf(execController)
					method := vc.MethodByName(runMethod)
					in := param.ConvertParams(methodParams, method.Type(), ctx)
					out := method.Call(in)

					// For backward compatibility we only handle response if we had incoming methodParams
					if methodParams != nil {
						p.handleParamResponse(ctx, execController, out)
					}
				}
			}

			// render template
			if !ctx.ResponseWriter.Started && ctx.Output.Status == 0 {
				if p.cfg.WebConfig.AutoRender {
					if err := execController.Render(); err != nil {
						logs.Error(err)
					}
				}
			}
		}

		// finish all runRouter. release resource
		execController.Finish()
	}

	// execute middleware filters
	if len(p.filters[AfterExec]) > 0 && p.execFilter(ctx, urlPath, AfterExec) {
		goto Admin
	}

	if len(p.filters[FinishRouter]) > 0 && p.execFilter(ctx, urlPath, FinishRouter) {
		goto Admin
	}

Admin:
	// admin module record QPS

	statusCode := ctx.ResponseWriter.Status
	if statusCode == 0 {
		statusCode = 200
	}

	LogAccess(ctx, &startTime, statusCode)

	timeDur := time.Since(startTime)
	ctx.ResponseWriter.Elapsed = timeDur
	if p.cfg.Listen.EnableAdmin {
		pattern := ""
		if routerInfo != nil {
			pattern = routerInfo.pattern
		}

		if FilterMonitorFunc(r.Method, r.URL.Path, timeDur, pattern, statusCode) {
			routerName := ""
			if runRouter != nil {
				routerName = runRouter.Name()
			}
			go StatisticsMap.AddStatistics(r.Method, r.URL.Path, routerName, timeDur)
		}
	}

	if p.cfg.RunMode == DEV && !p.cfg.Log.AccessLogs {
		match := map[bool]string{true: "match", false: "nomatch"}
		devInfo := fmt.Sprintf("|%15s|%s %3d %s|%13s|%8s|%s %-7s %s %-3s",
			ctx.Input.IP(),
			logs.ColorByStatus(statusCode), statusCode, logs.ResetColor(),
			timeDur.String(),
			match[findRouter],
			logs.ColorByMethod(r.Method), r.Method, logs.ResetColor(),
			r.URL.Path)
		if routerInfo != nil {
			devInfo += fmt.Sprintf("   r:%s", routerInfo.pattern)
		}

		logs.Debug(devInfo)
	}
	// Call WriteHeader if status code has been set changed
	if ctx.Output.Status != 0 {
		ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
	}
}

func (p *ControllerRegister) getUrlPath(ctx *beecontext.Context) string {
	urlPath := ctx.Request.URL.Path
	if !p.cfg.RouterCaseSensitive {
		urlPath = strings.ToLower(urlPath)
	}
	return urlPath
}

func (p *ControllerRegister) handleParamResponse(context *beecontext.Context, execController ControllerInterface, results []reflect.Value) {
	// looping in reverse order for the case when both error and value are returned and error sets the response status code
	for i := len(results) - 1; i >= 0; i-- {
		result := results[i]
		if result.Kind() != reflect.Interface || !result.IsNil() {
			resultValue := result.Interface()
			context.RenderMethodResult(resultValue)
		}
	}
	if !context.ResponseWriter.Started && len(results) > 0 && context.Output.Status == 0 {
		context.Output.SetStatus(200)
	}
}

// FindRouter Find Router info for URL
func (p *ControllerRegister) FindRouter(context *beecontext.Context) (routerInfo *ControllerInfo, isFind bool) {
	var urlPath = context.Input.URL()
	if !p.cfg.RouterCaseSensitive {
		urlPath = strings.ToLower(urlPath)
	}
	httpMethod := context.Input.Method()
	if t, ok := p.routers[httpMethod]; ok {
		runObject := t.Match(urlPath, context)
		if r, ok := runObject.(*ControllerInfo); ok {
			return r, true
		}
	}
	return
}

func toURL(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	u := "?"
	for k, v := range params {
		u += k + "=" + v + "&"
	}
	return strings.TrimRight(u, "&")
}

// LogAccess logging info HTTP Access
func LogAccess(ctx *beecontext.Context, startTime *time.Time, statusCode int) {
	BeeApp.LogAccess(ctx, startTime, statusCode)
}
