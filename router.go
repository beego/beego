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

package beego

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	beecontext "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/middleware"
	"github.com/astaxie/beego/toolbox"
	"github.com/astaxie/beego/utils"
)

const (
	// default filter execution points
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
	// supported http methods.
	HTTPMETHOD = map[string]string{
		"GET":     "GET",
		"POST":    "POST",
		"PUT":     "PUT",
		"DELETE":  "DELETE",
		"PATCH":   "PATCH",
		"OPTIONS": "OPTIONS",
		"HEAD":    "HEAD",
		"TRACE":   "TRACE",
		"CONNECT": "CONNECT",
	}
	// these beego.Controller's methods shouldn't reflect to AutoRouter
	exceptMethod = []string{"Init", "Prepare", "Finish", "Render", "RenderString",
		"RenderBytes", "Redirect", "Abort", "StopRun", "UrlFor", "ServeJson", "ServeJsonp",
		"ServeXml", "Input", "ParseForm", "GetString", "GetStrings", "GetInt", "GetBool",
		"GetFloat", "GetFile", "SaveToFile", "StartSession", "SetSession", "GetSession",
		"DelSession", "SessionRegenerateID", "DestroySession", "IsAjax", "GetSecureCookie",
		"SetSecureCookie", "XsrfToken", "CheckXsrfCookie", "XsrfFormHtml",
		"GetControllerAndAction"}

	url_placeholder                = "{{placeholder}}"
	DefaultLogFilter FilterHandler = &logFilter{}
)

type FilterHandler interface {
	Filter(*beecontext.Context) bool
}

// default log filter static file will not show
type logFilter struct {
}

func (l *logFilter) Filter(ctx *beecontext.Context) bool {
	requestPath := path.Clean(ctx.Input.Request.URL.Path)
	if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
		return true
	}
	for prefix, _ := range StaticDir {
		if strings.HasPrefix(requestPath, prefix) {
			return true
		}
	}
	return false
}

// To append a slice's value into "exceptMethod", for controller's methods shouldn't reflect to AutoRouter
func ExceptMethodAppend(action string) {
	exceptMethod = append(exceptMethod, action)
}

type controllerInfo struct {
	pattern        string
	controllerType reflect.Type
	methods        map[string]string
	handler        http.Handler
	runfunction    FilterFunc
	routerType     int
}

// ControllerRegistor containers registered router rules, controller handlers and filters.
type ControllerRegistor struct {
	routers      map[string]*Tree
	enableFilter bool
	filters      map[int][]*FilterRouter
}

// NewControllerRegister returns a new ControllerRegistor.
func NewControllerRegister() *ControllerRegistor {
	return &ControllerRegistor{
		routers: make(map[string]*Tree),
		filters: make(map[int][]*FilterRouter),
	}
}

// Add controller handler and pattern rules to ControllerRegistor.
// usage:
//	default methods is the same name as method
//	Add("/user",&UserController{})
//	Add("/api/list",&RestController{},"*:ListFood")
//	Add("/api/create",&RestController{},"post:CreateFood")
//	Add("/api/update",&RestController{},"put:UpdateFood")
//	Add("/api/delete",&RestController{},"delete:DeleteFood")
//	Add("/api",&RestController{},"get,post:ApiFunc")
//	Add("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
func (p *ControllerRegistor) Add(pattern string, c ControllerInterface, mappingMethods ...string) {
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
				if _, ok := HTTPMETHOD[strings.ToUpper(m)]; m == "*" || ok {
					if val := reflectVal.MethodByName(colon[1]); val.IsValid() {
						methods[strings.ToUpper(m)] = colon[1]
					} else {
						panic(colon[1] + " method doesn't exist in the controller " + t.Name())
					}
				} else {
					panic(v + " is an invalid method mapping. Method doesn't exist " + m)
				}
			}
		}
	}

	route := &controllerInfo{}
	route.pattern = pattern
	route.methods = methods
	route.routerType = routerTypeBeego
	route.controllerType = t
	if len(methods) == 0 {
		for _, m := range HTTPMETHOD {
			p.addToRouter(m, pattern, route)
		}
	} else {
		for k, _ := range methods {
			if k == "*" {
				for _, m := range HTTPMETHOD {
					p.addToRouter(m, pattern, route)
				}
			} else {
				p.addToRouter(k, pattern, route)
			}
		}
	}
}

func (p *ControllerRegistor) addToRouter(method, pattern string, r *controllerInfo) {
	if !RouterCaseSensitive {
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

// only when the Runmode is dev will generate router file in the router/auto.go from the controller
// Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
func (p *ControllerRegistor) Include(cList ...ControllerInterface) {
	if RunMode == "dev" {
		skip := make(map[string]bool, 10)
		for _, c := range cList {
			reflectVal := reflect.ValueOf(c)
			t := reflect.Indirect(reflectVal).Type()
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				panic("you are in dev mode. So please set gopath")
			}
			pkgpath := ""

			wgopath := filepath.SplitList(gopath)
			for _, wg := range wgopath {
				wg, _ = filepath.EvalSymlinks(filepath.Join(wg, "src", t.PkgPath()))
				if utils.FileExists(wg) {
					pkgpath = wg
					break
				}
			}
			if pkgpath != "" {
				if _, ok := skip[pkgpath]; !ok {
					skip[pkgpath] = true
					parserPkg(pkgpath, t.PkgPath())
				}
			}
		}
	}
	for _, c := range cList {
		reflectVal := reflect.ValueOf(c)
		t := reflect.Indirect(reflectVal).Type()
		key := t.PkgPath() + ":" + t.Name()
		if comm, ok := GlobalControllerRouter[key]; ok {
			for _, a := range comm {
				p.Add(a.Router, c, strings.Join(a.AllowHTTPMethods, ",")+":"+a.Method)
			}
		}
	}
}

// add get method
// usage:
//    Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Get(pattern string, f FilterFunc) {
	p.AddMethod("get", pattern, f)
}

// add post method
// usage:
//    Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Post(pattern string, f FilterFunc) {
	p.AddMethod("post", pattern, f)
}

// add put method
// usage:
//    Put("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Put(pattern string, f FilterFunc) {
	p.AddMethod("put", pattern, f)
}

// add delete method
// usage:
//    Delete("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Delete(pattern string, f FilterFunc) {
	p.AddMethod("delete", pattern, f)
}

// add head method
// usage:
//    Head("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Head(pattern string, f FilterFunc) {
	p.AddMethod("head", pattern, f)
}

// add patch method
// usage:
//    Patch("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Patch(pattern string, f FilterFunc) {
	p.AddMethod("patch", pattern, f)
}

// add options method
// usage:
//    Options("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Options(pattern string, f FilterFunc) {
	p.AddMethod("options", pattern, f)
}

// add all method
// usage:
//    Any("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Any(pattern string, f FilterFunc) {
	p.AddMethod("*", pattern, f)
}

// add http method router
// usage:
//    AddMethod("get","/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) AddMethod(method, pattern string, f FilterFunc) {
	if _, ok := HTTPMETHOD[strings.ToUpper(method)]; method != "*" && !ok {
		panic("not support http method: " + method)
	}
	route := &controllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeRESTFul
	route.runfunction = f
	methods := make(map[string]string)
	if method == "*" {
		for _, val := range HTTPMETHOD {
			methods[val] = val
		}
	} else {
		methods[strings.ToUpper(method)] = strings.ToUpper(method)
	}
	route.methods = methods
	for k, _ := range methods {
		if k == "*" {
			for _, m := range HTTPMETHOD {
				p.addToRouter(m, pattern, route)
			}
		} else {
			p.addToRouter(k, pattern, route)
		}
	}
}

// add user defined Handler
func (p *ControllerRegistor) Handler(pattern string, h http.Handler, options ...interface{}) {
	route := &controllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeHandler
	route.handler = h
	if len(options) > 0 {
		if _, ok := options[0].(bool); ok {
			pattern = path.Join(pattern, "?:all")
		}
	}
	for _, m := range HTTPMETHOD {
		p.addToRouter(m, pattern, route)
	}
}

// Add auto router to ControllerRegistor.
// example beego.AddAuto(&MainContorlller{}),
// MainController has method List and Page.
// visit the url /main/list to execute List function
// /main/page to execute Page function.
func (p *ControllerRegistor) AddAuto(c ControllerInterface) {
	p.AddAutoPrefix("/", c)
}

// Add auto router to ControllerRegistor with prefix.
// example beego.AddAutoPrefix("/admin",&MainContorlller{}),
// MainController has method List and Page.
// visit the url /admin/main/list to execute List function
// /admin/main/page to execute Page function.
func (p *ControllerRegistor) AddAutoPrefix(prefix string, c ControllerInterface) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	for i := 0; i < rt.NumMethod(); i++ {
		if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
			route := &controllerInfo{}
			route.routerType = routerTypeBeego
			route.methods = map[string]string{"*": rt.Method(i).Name}
			route.controllerType = ct
			pattern := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name), "*")
			patternInit := path.Join(prefix, controllerName, rt.Method(i).Name, "*")
			patternfix := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
			patternfixInit := path.Join(prefix, controllerName, rt.Method(i).Name)
			route.pattern = pattern
			for _, m := range HTTPMETHOD {
				p.addToRouter(m, pattern, route)
				p.addToRouter(m, patternInit, route)
				p.addToRouter(m, patternfix, route)
				p.addToRouter(m, patternfixInit, route)
			}
		}
	}
}

// Add a FilterFunc with pattern rule and action constant.
// The bool params is for setting the returnOnOutput value (false allows multiple filters to execute)
func (p *ControllerRegistor) InsertFilter(pattern string, pos int, filter FilterFunc, params ...bool) error {

	mr := new(FilterRouter)
	mr.tree = NewTree()
	mr.pattern = pattern
	mr.filterFunc = filter
	if !RouterCaseSensitive {
		pattern = strings.ToLower(pattern)
	}
	if len(params) == 0 {
		mr.returnOnOutput = true
	} else {
		mr.returnOnOutput = params[0]
	}
	mr.tree.AddRouter(pattern, true)
	return p.insertFilterRouter(pos, mr)
}

// add Filter into
func (p *ControllerRegistor) insertFilterRouter(pos int, mr *FilterRouter) error {
	p.filters[pos] = append(p.filters[pos], mr)
	p.enableFilter = true
	return nil
}

// UrlFor does another controller handler in this request function.
// it can access any controller method.
func (p *ControllerRegistor) UrlFor(endpoint string, values ...string) string {
	paths := strings.Split(endpoint, ".")
	if len(paths) <= 1 {
		Warn("urlfor endpoint must like path.controller.method")
		return ""
	}
	if len(values)%2 != 0 {
		Warn("urlfor params must key-value pair")
		return ""
	}
	params := make(map[string]string)
	if len(values) > 0 {
		key := ""
		for k, v := range values {
			if k%2 == 0 {
				key = v
			} else {
				params[key] = v
			}
		}
	}
	controllName := strings.Join(paths[:len(paths)-1], "/")
	methodName := paths[len(paths)-1]
	for m, t := range p.routers {
		ok, url := p.geturl(t, "/", controllName, methodName, params, m)
		if ok {
			return url
		}
	}
	return ""
}

func (p *ControllerRegistor) geturl(t *Tree, url, controllName, methodName string, params map[string]string, httpMethod string) (bool, string) {
	for k, subtree := range t.fixrouters {
		u := path.Join(url, k)
		ok, u := p.geturl(subtree, u, controllName, methodName, params, httpMethod)
		if ok {
			return ok, u
		}
	}
	if t.wildcard != nil {
		u := path.Join(url, url_placeholder)
		ok, u := p.geturl(t.wildcard, u, controllName, methodName, params, httpMethod)
		if ok {
			return ok, u
		}
	}
	for _, l := range t.leaves {
		if c, ok := l.runObject.(*controllerInfo); ok {
			if c.routerType == routerTypeBeego &&
				strings.HasSuffix(path.Join(c.controllerType.PkgPath(), c.controllerType.Name()), controllName) {
				find := false
				if _, ok := HTTPMETHOD[strings.ToUpper(methodName)]; ok {
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
							return true, strings.Replace(url, "/"+url_placeholder, "", 1) + tourl(params)
						}
						if len(l.wildcards) == 1 {
							if v, ok := params[l.wildcards[0]]; ok {
								delete(params, l.wildcards[0])
								return true, strings.Replace(url, url_placeholder, v, 1) + tourl(params)
							} else {
								return false, ""
							}
						}
						if len(l.wildcards) == 3 && l.wildcards[0] == "." {
							if p, ok := params[":path"]; ok {
								if e, isok := params[":ext"]; isok {
									delete(params, ":path")
									delete(params, ":ext")
									return true, strings.Replace(url, url_placeholder, p+"."+e, -1) + tourl(params)
								}
							}
						}
						canskip := false
						for _, v := range l.wildcards {
							if v == ":" {
								canskip = true
								continue
							}
							if u, ok := params[v]; ok {
								delete(params, v)
								url = strings.Replace(url, url_placeholder, u, 1)
							} else {
								if canskip {
									canskip = false
									continue
								} else {
									return false, ""
								}
							}
						}
						return true, url + tourl(params)
					} else {
						var i int
						var startreg bool
						regurl := ""
						for _, v := range strings.Trim(l.regexps.String(), "^$") {
							if v == '(' {
								startreg = true
								continue
							} else if v == ')' {
								startreg = false
								if v, ok := params[l.wildcards[i]]; ok {
									delete(params, l.wildcards[i])
									regurl = regurl + v
									i++
								} else {
									break
								}
							} else if !startreg {
								regurl = string(append([]rune(regurl), v))
							}
						}
						if l.regexps.MatchString(regurl) {
							ps := strings.Split(regurl, "/")
							for _, p := range ps {
								url = strings.Replace(url, url_placeholder, p, 1)
							}
							return true, url + tourl(params)
						}
					}
				}
			}
		}
	}

	return false, ""
}

// Implement http.Handler interface.
func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer p.recoverPanic(rw, r)
	starttime := time.Now()
	var runrouter reflect.Type
	var findrouter bool
	var runMethod string
	var routerInfo *controllerInfo

	w := &responseWriter{writer: rw}

	if RunMode == "dev" {
		w.Header().Set("Server", BeegoServerName)
	}

	// init context
	context := &beecontext.Context{
		ResponseWriter: w,
		Request:        r,
		Input:          beecontext.NewInput(r),
		Output:         beecontext.NewOutput(),
	}
	context.Output.Context = context
	context.Output.EnableGzip = EnableGzip

	var urlPath string
	if !RouterCaseSensitive {
		urlPath = strings.ToLower(r.URL.Path)
	} else {
		urlPath = r.URL.Path
	}
	// defined filter function
	do_filter := func(pos int) (started bool) {
		if p.enableFilter {
			if l, ok := p.filters[pos]; ok {
				for _, filterR := range l {
					if ok, p := filterR.ValidRouter(urlPath); ok {
						context.Input.Params = p
						filterR.filterFunc(context)
						if filterR.returnOnOutput && w.started {
							return true
						}
					}
				}
			}
		}

		return false
	}

	// filter wrong httpmethod
	if _, ok := HTTPMETHOD[r.Method]; !ok {
		http.Error(w, "Method Not Allowed", 405)
		goto Admin
	}

	// filter for static file
	if do_filter(BeforeStatic) {
		goto Admin
	}

	serverStaticRouter(context)
	if w.started {
		findrouter = true
		goto Admin
	}

	// session init
	if SessionOn {
		var err error
		context.Input.CruSession, err = GlobalSessions.SessionStart(w, r)
		if err != nil {
			Error(err)
			middleware.Exception("503", rw, r, "")
			return
		}
		defer func() {
			context.Input.CruSession.SessionRelease(w)
		}()
	}

	if r.Method != "GET" && r.Method != "HEAD" {
		if CopyRequestBody && !context.Input.IsUpload() {
			context.Input.CopyBody()
		}
		context.Input.ParseFormOrMulitForm(MaxMemory)
	}

	if do_filter(BeforeRouter) {
		goto Admin
	}

	if context.Input.RunController != nil && context.Input.RunMethod != "" {
		findrouter = true
		runMethod = context.Input.RunMethod
		runrouter = context.Input.RunController
	}

	if !findrouter {
		http_method := r.Method

		if http_method == "POST" && context.Input.Query("_method") == "PUT" {
			http_method = "PUT"
		}

		if http_method == "POST" && context.Input.Query("_method") == "DELETE" {
			http_method = "DELETE"
		}

		if t, ok := p.routers[http_method]; ok {
			runObject, p := t.Match(urlPath)
			if r, ok := runObject.(*controllerInfo); ok {
				routerInfo = r
				findrouter = true
				if splat, ok := p[":splat"]; ok {
					splatlist := strings.Split(splat, "/")
					for k, v := range splatlist {
						p[strconv.Itoa(k)] = v
					}
				}
				context.Input.Params = p
			}
		}

	}

	//if no matches to url, throw a not found exception
	if !findrouter {
		middleware.Exception("404", rw, r, "")
		goto Admin
	}

	if findrouter {
		//execute middleware filters
		if do_filter(BeforeExec) {
			goto Admin
		}
		isRunable := false
		if routerInfo != nil {
			if routerInfo.routerType == routerTypeRESTFul {
				if _, ok := routerInfo.methods[r.Method]; ok {
					isRunable = true
					routerInfo.runfunction(context)
				} else {
					middleware.Exception("405", rw, r, "Method Not Allowed")
					goto Admin
				}
			} else if routerInfo.routerType == routerTypeHandler {
				isRunable = true
				routerInfo.handler.ServeHTTP(rw, r)
			} else {
				runrouter = routerInfo.controllerType
				method := r.Method
				if r.Method == "POST" && context.Input.Query("_method") == "PUT" {
					method = "PUT"
				}
				if r.Method == "POST" && context.Input.Query("_method") == "DELETE" {
					method = "DELETE"
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

		// also defined runrouter & runMethod from filter
		if !isRunable {
			//Invoke the request handler
			vc := reflect.New(runrouter)
			execController, ok := vc.Interface().(ControllerInterface)
			if !ok {
				panic("controller is not ControllerInterface")
			}

			//call the controller init function
			execController.Init(context, runrouter.Name(), runMethod, vc.Interface())

			//call prepare function
			execController.Prepare()

			//if XSRF is Enable then check cookie where there has any cookie in the  request's cookie _csrf
			if EnableXSRF {
				execController.XsrfToken()
				if r.Method == "POST" || r.Method == "DELETE" || r.Method == "PUT" ||
					(r.Method == "POST" && (context.Input.Query("_method") == "DELETE" || context.Input.Query("_method") == "PUT")) {
					execController.CheckXsrfCookie()
				}
			}

			execController.URLMapping()

			if !w.started {
				//exec main logic
				switch runMethod {
				case "GET":
					execController.Get()
				case "POST":
					execController.Post()
				case "DELETE":
					execController.Delete()
				case "PUT":
					execController.Put()
				case "HEAD":
					execController.Head()
				case "PATCH":
					execController.Patch()
				case "OPTIONS":
					execController.Options()
				default:
					if !execController.HandlerFunc(runMethod) {
						in := make([]reflect.Value, 0)
						method := vc.MethodByName(runMethod)
						method.Call(in)
					}
				}

				//render template
				if !w.started && context.Output.Status == 0 {
					if AutoRender {
						if err := execController.Render(); err != nil {
							panic(err)
						}
					}
				}
			}

			// finish all runrouter. release resource
			execController.Finish()
		}

		//execute middleware filters
		if do_filter(AfterExec) {
			goto Admin
		}
	}

	do_filter(FinishRouter)

Admin:
	timeend := time.Since(starttime)
	//admin module record QPS
	if EnableAdmin {
		if FilterMonitorFunc(r.Method, r.URL.Path, timeend) {
			if runrouter != nil {
				go toolbox.StatisticsMap.AddStatistics(r.Method, r.URL.Path, runrouter.Name(), timeend)
			} else {
				go toolbox.StatisticsMap.AddStatistics(r.Method, r.URL.Path, "", timeend)
			}
		}
	}

	if RunMode == "dev" {
		var devinfo string
		if findrouter {
			if routerInfo != nil {
				devinfo = fmt.Sprintf("| % -10s | % -40s | % -16s | % -10s | % -40s |", r.Method, r.URL.Path, timeend.String(), "match", routerInfo.pattern)
			} else {
				devinfo = fmt.Sprintf("| % -10s | % -40s | % -16s | % -10s |", r.Method, r.URL.Path, timeend.String(), "match")
			}
		} else {
			devinfo = fmt.Sprintf("| % -10s | % -40s | % -16s | % -10s |", r.Method, r.URL.Path, timeend.String(), "notmatch")
		}
		if DefaultLogFilter == nil || !DefaultLogFilter.Filter(context) {
			Debug(devinfo)
		}
	}

	// Call WriteHeader if status code has been set changed
	if context.Output.Status != 0 {
		w.writer.WriteHeader(context.Output.Status)
	}
}

func (p *ControllerRegistor) recoverPanic(rw http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		if err == USERSTOPRUN {
			return
		}
		if he, ok := err.(middleware.HTTPException); ok {
			rw.Write([]byte(he.Description))
			// catch intented errors, only for HTTP 4XX and 5XX
		} else {
			if RunMode == "dev" {
				if !RecoverPanic {
					panic(err)
				} else {
					if ErrorsShow {
						if handler, ok := middleware.ErrorMaps[fmt.Sprint(err)]; ok {
							handler(rw, r)
							return
						}
					}
					var stack string
					Critical("the request url is ", r.URL.Path)
					Critical("Handler crashed with error", err)
					for i := 1; ; i++ {
						_, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}
						Critical(fmt.Sprintf("%s:%d", file, line))
						stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
					}
					middleware.ShowErr(err, rw, r, stack)
				}
			} else {
				if !RecoverPanic {
					panic(err)
				} else {
					// in production model show all infomation
					if ErrorsShow {
						if handler, ok := middleware.ErrorMaps[fmt.Sprint(err)]; ok {
							handler(rw, r)
							return
						} else if handler, ok := middleware.ErrorMaps["503"]; ok {
							handler(rw, r)
							return
						} else {
							rw.Write([]byte(fmt.Sprint(err)))
						}
					} else {
						Critical("the request url is ", r.URL.Path)
						Critical("Handler crashed with error", err)
						for i := 1; ; i++ {
							_, file, line, ok := runtime.Caller(i)
							if !ok {
								break
							}
							Critical(fmt.Sprintf("%s:%d", file, line))
						}
					}
				}
			}

		}
	}
}

//responseWriter is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type responseWriter struct {
	writer  http.ResponseWriter
	started bool
	status  int
}

// Header returns the header map that will be sent by WriteHeader.
func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true.
// started means the response has sent out.
func (w *responseWriter) Write(p []byte) (int, error) {
	w.started = true
	return w.writer.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.started = true
	w.writer.WriteHeader(code)
}

// hijacker for http
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.writer.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

func tourl(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	u := "?"
	for k, v := range params {
		u += k + "=" + v + "&"
	}
	return strings.TrimRight(u, "&")
}
