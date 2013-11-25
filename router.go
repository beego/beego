package beego

import (
	"fmt"
	beecontext "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/middleware"
	"github.com/astaxie/beego/toolbox"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	BeforeRouter = iota
	AfterStatic
	BeforeExec
	AfterExec
	FinishRouter
)

var HTTPMETHOD = []string{"get", "post", "put", "delete", "patch", "options", "head"}

type controllerInfo struct {
	pattern        string
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
	methods        map[string]string
	hasMethod      bool
}

type ControllerRegistor struct {
	routers      []*controllerInfo
	fixrouters   []*controllerInfo
	enableFilter bool
	filters      map[int][]*FilterRouter
	enableAuto   bool
	autoRouter   map[string]map[string]reflect.Type //key:controller key:method value:reflect.type
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{
		routers:    make([]*controllerInfo, 0),
		autoRouter: make(map[string]map[string]reflect.Type),
		filters:    make(map[int][]*FilterRouter),
	}
}

//methods support like this:
//default methods is the same name as method
//Add("/user",&UserController{})
//Add("/api/list",&RestController{},"*:ListFood")
//Add("/api/create",&RestController{},"post:CreateFood")
//Add("/api/update",&RestController{},"put:UpdateFood")
//Add("/api/delete",&RestController{},"delete:DeleteFood")
//Add("/api",&RestController{},"get,post:ApiFunc")
//Add("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
func (p *ControllerRegistor) Add(pattern string, c ControllerInterface, mappingMethods ...string) {
	parts := strings.Split(pattern, "/")

	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "(.+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
				//match /user/:id:int ([0-9]+)
				//match /post/:username:string	([\w]+)
			} else if lindex := strings.LastIndex(part, ":"); lindex != 0 {
				switch part[lindex:] {
				case ":int":
					expr = "([0-9]+)"
					part = part[:lindex]
				case ":string":
					expr = `([\w]+)`
					part = part[:lindex]
				}
			}
			params[j] = part
			parts[i] = expr
			j++
		}
		if strings.HasPrefix(part, "*") {
			expr := "(.+)"
			if part == "*.*" {
				params[j] = ":path"
				parts[i] = "([^.]+).([^.]+)"
				j++
				params[j] = ":ext"
				j++
			} else {
				params[j] = ":splat"
				parts[i] = expr
				j++
			}
		}
		//url like someprefix:id(xxx).html
		if strings.Contains(part, ":") && strings.Contains(part, "(") && strings.Contains(part, ")") {
			var out []rune
			var start bool
			var startexp bool
			var param []rune
			var expt []rune
			for _, v := range part {
				if start {
					if v != '(' {
						param = append(param, v)
						continue
					}
				}
				if startexp {
					if v != ')' {
						expt = append(expt, v)
						continue
					}
				}
				if v == ':' {
					param = make([]rune, 0)
					param = append(param, ':')
					start = true
				} else if v == '(' {
					startexp = true
					start = false
					params[j] = string(param)
					j++
					expt = make([]rune, 0)
					expt = append(expt, '(')
				} else if v == ')' {
					startexp = false
					expt = append(expt, ')')
					out = append(out, expt...)
				} else {
					out = append(out, v)
				}
			}
			parts[i] = string(out)
		}
	}
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	methods := make(map[string]string)
	if len(mappingMethods) > 0 {
		semi := strings.Split(mappingMethods[0], ";")
		for _, v := range semi {
			colon := strings.Split(v, ":")
			if len(colon) != 2 {
				panic("method mapping fomate is error")
			}
			comma := strings.Split(colon[0], ",")
			for _, m := range comma {
				if m == "*" || inSlice(strings.ToLower(m), HTTPMETHOD) {
					if val := reflectVal.MethodByName(colon[1]); val.IsValid() {
						methods[strings.ToLower(m)] = colon[1]
					} else {
						panic(colon[1] + " method don't exist in the controller " + t.Name())
					}
				} else {
					panic(v + " is an error method mapping,Don't exist method named " + m)
				}
			}
		}
	}
	if j == 0 {
		//now create the Route
		route := &controllerInfo{}
		route.pattern = pattern
		route.controllerType = t
		route.methods = methods
		if len(methods) > 0 {
			route.hasMethod = true
		}
		p.fixrouters = append(p.fixrouters, route)
	} else { // add regexp routers
		//recreate the url pattern, with parameters replaced
		//by regular expressions. then compile the regex
		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {
			//TODO add error handling here to avoid panic
			panic(regexErr)
			return
		}

		//now create the Route

		route := &controllerInfo{}
		route.regex = regex
		route.params = params
		route.pattern = pattern
		route.methods = methods
		if len(methods) > 0 {
			route.hasMethod = true
		}
		route.controllerType = t
		p.routers = append(p.routers, route)
	}
}

func (p *ControllerRegistor) AddAuto(c ControllerInterface) {
	p.enableAuto = true
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	firstParam := strings.ToLower(strings.TrimSuffix(ct.Name(), "Controller"))
	if _, ok := p.autoRouter[firstParam]; ok {
		return
	} else {
		p.autoRouter[firstParam] = make(map[string]reflect.Type)
	}
	for i := 0; i < rt.NumMethod(); i++ {
		p.autoRouter[firstParam][rt.Method(i).Name] = ct
	}
}

// Filter adds the middleware filter.
func buildFilter(pattern string, filter FilterFunc) *FilterRouter {
	mr := new(FilterRouter)
	mr.filterFunc = filter
	parts := strings.Split(pattern, "/")
	j := 0
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "(.+)"
			//a user may choose to override the default expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
				//match /user/:id:int ([0-9]+)
				//match /post/:username:string	([\w]+)
			} else if lindex := strings.LastIndex(part, ":"); lindex != 0 {
				switch part[lindex:] {
				case ":int":
					expr = "([0-9]+)"
					part = part[:lindex]
				case ":string":
					expr = `([\w]+)`
					part = part[:lindex]
				}
			}
			parts[i] = expr
			j++
		}
	}
	if j != 0 {
		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {
			//TODO add error handling here to avoid panic
			panic(regexErr)
		}
		mr.regex = regex
		mr.hasregex = true
	}
	mr.pattern = pattern
	return mr
}

//p.filters[action] = append(p.filters[action], mr)
func (p *ControllerRegistor) AddFilter(pattern, action string, filter FilterFunc) {
	mr := buildFilter(pattern, filter)
	switch action {
	case "BeforRouter":
		p.filters[BeforeRouter] = append(p.filters[BeforeRouter], mr)
	case "AfterStatic":
		p.filters[AfterStatic] = append(p.filters[AfterStatic], mr)
	case "BeforeExec":
		p.filters[BeforeExec] = append(p.filters[BeforeExec], mr)
	case "AfterExec":
		p.filters[AfterExec] = append(p.filters[AfterExec], mr)
	case "FinishRouter":
		p.filters[FinishRouter] = append(p.filters[FinishRouter], mr)
	}
	p.enableFilter = true
}

func (p *ControllerRegistor) InsertFilter(pattern string, filterPos int, filter FilterFunc) {
	mr := buildFilter(pattern, filter)
	p.filters[filterPos] = append(p.filters[filterPos], mr)
	p.enableFilter = true
}

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
	urlv := url.Values{}
	if len(values) > 0 {
		key := ""
		for k, v := range values {
			if k%2 == 0 {
				key = v
			} else {
				urlv.Set(key, v)
			}
		}
	}
	controllName := strings.Join(paths[:len(paths)-1], ".")
	methodName := paths[len(paths)-1]
	for _, route := range p.fixrouters {
		if route.controllerType.Name() == controllName {
			var finded bool
			if inSlice(strings.ToLower(methodName), HTTPMETHOD) {
				if route.hasMethod {
					if m, ok := route.methods[strings.ToLower(methodName)]; ok && m != methodName {
						finded = false
					} else if m, ok = route.methods["*"]; ok && m != methodName {
						finded = false
					} else {
						finded = true
					}
				} else {
					finded = true
				}
			} else if route.hasMethod {
				for _, md := range route.methods {
					if md == methodName {
						finded = true
					}
				}
			}
			if !finded {
				continue
			}
			if len(values) > 0 {
				return route.pattern + "?" + urlv.Encode()
			}
			return route.pattern
		}
	}
	for _, route := range p.routers {
		if route.controllerType.Name() == controllName {
			var finded bool
			if inSlice(strings.ToLower(methodName), HTTPMETHOD) {
				if route.hasMethod {
					if m, ok := route.methods[strings.ToLower(methodName)]; ok && m != methodName {
						finded = false
					} else if m, ok = route.methods["*"]; ok && m != methodName {
						finded = false
					} else {
						finded = true
					}
				} else {
					finded = true
				}
			} else if route.hasMethod {
				for _, md := range route.methods {
					if md == methodName {
						finded = true
					}
				}
			}
			if !finded {
				continue
			}
			var returnurl string
			var i int
			var startreg bool
			for _, v := range route.regex.String() {
				if v == '(' {
					startreg = true
					continue
				} else if v == ')' {
					startreg = false
					returnurl = returnurl + urlv.Get(route.params[i])
					i++
				} else if !startreg {
					returnurl = string(append([]rune(returnurl), v))
				}
			}
			if route.regex.MatchString(returnurl) {
				return returnurl
			}
		}
	}
	if p.enableAuto {
		for cName, methodList := range p.autoRouter {
			if strings.ToLower(strings.TrimSuffix(paths[len(paths)-2], "Controller")) == cName {
				if _, ok := methodList[methodName]; ok {
					if len(values) > 0 {
						return "/" + strings.TrimSuffix(paths[len(paths)-2], "Controller") + "/" + methodName + "?" + urlv.Encode()
					} else {
						return "/" + strings.TrimSuffix(paths[len(paths)-2], "Controller") + "/" + methodName
					}
				}
			}
		}
	}
	return ""
}

// AutoRoute
func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errstr := fmt.Sprint(err)
			if handler, ok := middleware.ErrorMaps[errstr]; ok && ErrorsShow {
				handler(rw, r)
			} else {
				if !RecoverPanic {
					// go back to panic
					panic(err)
				} else {
					var stack string
					Critical("Handler crashed with error", err)
					for i := 1; ; i++ {
						_, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}
						Critical(file, line)
						if RunMode == "dev" {
							stack = stack + fmt.Sprintln(file, line)
						}
					}
					if RunMode == "dev" {
						middleware.ShowErr(err, rw, r, stack)
					}
				}
			}
		}
	}()

	starttime := time.Now()
	requestPath := r.URL.Path
	var runrouter *controllerInfo
	var findrouter bool
	params := make(map[string]string)

	w := &responseWriter{writer: rw}
	w.Header().Set("Server", BeegoServerName)
	context := &beecontext.Context{
		ResponseWriter: w,
		Request:        r,
		Input:          beecontext.NewInput(r),
		Output:         beecontext.NewOutput(w),
	}
	context.Output.Context = context
	context.Output.EnableGzip = EnableGzip

	if context.Input.IsWebsocket() {
		context.ResponseWriter = rw
		context.Output = beecontext.NewOutput(rw)
	}

	if !inSlice(strings.ToLower(r.Method), HTTPMETHOD) {
		http.Error(w, "Method Not Allowed", 405)
		goto Admin
	}

	if p.enableFilter {
		if l, ok := p.filters[BeforeRouter]; ok {
			for _, filterR := range l {
				if filterR.ValidRouter(r.URL.Path) {
					filterR.filterFunc(context)
					if w.started {
						goto Admin
					}
				}
			}
		}
	}

	//static file server
	for prefix, staticDir := range StaticDir {
		if r.URL.Path == "/favicon.ico" {
			file := staticDir + r.URL.Path
			http.ServeFile(w, r, file)
			w.started = true
			goto Admin
		}
		if strings.HasPrefix(r.URL.Path, prefix) {
			file := staticDir + r.URL.Path[len(prefix):]
			finfo, err := os.Stat(file)
			if err != nil {
				if RunMode == "dev" {
					Warn(err)
				}
				http.NotFound(w, r)
				goto Admin
			}
			//if the request is dir and DirectoryIndex is false then
			if finfo.IsDir() && !DirectoryIndex {
				middleware.Exception("403", rw, r, "403 Forbidden")
				goto Admin
			}
			http.ServeFile(w, r, file)
			w.started = true
			goto Admin
		}
	}

	// session init after static file
	if SessionOn {
		context.Input.CruSession = GlobalSessions.SessionStart(w, r)
	}

	if p.enableFilter {
		if l, ok := p.filters[AfterStatic]; ok {
			for _, filterR := range l {
				if filterR.ValidRouter(r.URL.Path) {
					filterR.filterFunc(context)
					if w.started {
						goto Admin
					}
				}
			}
		}
	}

	if CopyRequestBody {
		context.Input.Body()
	}

	//first find path from the fixrouters to Improve Performance
	for _, route := range p.fixrouters {
		n := len(requestPath)
		if requestPath == route.pattern {
			runrouter = route
			findrouter = true
			break
		}
		// pattern /admin   url /admin 200  /admin/ 404
		// pattern /admin/  url /admin 301  /admin/ 200
		if requestPath[n-1] != '/' && len(route.pattern) == n+1 &&
			route.pattern[n] == '/' && route.pattern[:n] == requestPath {
			http.Redirect(w, r, requestPath+"/", 301)
			goto Admin
		}
	}

	//find regex's router
	if !findrouter {
		//find a matching Route
		for _, route := range p.routers {

			//check if Route pattern matches url
			if !route.regex.MatchString(requestPath) {
				continue
			}

			//get submatches (params)
			matches := route.regex.FindStringSubmatch(requestPath)

			//double check that the Route matches the URL pattern.
			if len(matches[0]) != len(requestPath) {
				continue
			}

			if len(route.params) > 0 {
				//add url parameters to the query param map
				values := r.URL.Query()
				for i, match := range matches[1:] {
					values.Add(route.params[i], match)
					params[route.params[i]] = match
				}
				//reassemble query params and add to RawQuery
				r.URL.RawQuery = url.Values(values).Encode()
			}
			runrouter = route
			findrouter = true
			break
		}
	}

	context.Input.Param = params

	if runrouter != nil {
		if r.Method == "POST" {
			r.ParseMultipartForm(MaxMemory)
		}
		//execute middleware filters
		if p.enableFilter {
			if l, ok := p.filters[BeforeExec]; ok {
				for _, filterR := range l {
					if filterR.ValidRouter(r.URL.Path) {
						filterR.filterFunc(context)
						if w.started {
							goto Admin
						}
					}
				}
			}
		}
		//Invoke the request handler
		vc := reflect.New(runrouter.controllerType)

		//call the controller init function
		method := vc.MethodByName("Init")
		in := make([]reflect.Value, 3)
		in[0] = reflect.ValueOf(context)
		in[1] = reflect.ValueOf(runrouter.controllerType.Name())
		in[2] = reflect.ValueOf(vc.Interface())
		method.Call(in)

		//if XSRF is Enable then check cookie where there has any cookie in the  request's cookie _csrf
		if EnableXSRF {
			in = make([]reflect.Value, 0)
			method = vc.MethodByName("XsrfToken")
			method.Call(in)
			if r.Method == "POST" || r.Method == "DELETE" || r.Method == "PUT" ||
				(r.Method == "POST" && (r.Form.Get("_method") == "delete" || r.Form.Get("_method") == "put")) {
				method = vc.MethodByName("CheckXsrfCookie")
				method.Call(in)
			}
		}

		//call prepare function
		in = make([]reflect.Value, 0)
		method = vc.MethodByName("Prepare")
		method.Call(in)

		//if response has written,yes don't run next
		if !w.started {
			if r.Method == "GET" {
				if runrouter.hasMethod {
					if m, ok := runrouter.methods["get"]; ok {
						method = vc.MethodByName(m)
					} else if m, ok = runrouter.methods["*"]; ok {
						method = vc.MethodByName(m)
					} else {
						method = vc.MethodByName("Get")
					}
				} else {
					method = vc.MethodByName("Get")
				}
				method.Call(in)
			} else if r.Method == "HEAD" {
				if runrouter.hasMethod {
					if m, ok := runrouter.methods["head"]; ok {
						method = vc.MethodByName(m)
					} else if m, ok = runrouter.methods["*"]; ok {
						method = vc.MethodByName(m)
					} else {
						method = vc.MethodByName("Head")
					}
				} else {
					method = vc.MethodByName("Head")
				}

				method.Call(in)
			} else if r.Method == "DELETE" || (r.Method == "POST" && r.Form.Get("_method") == "delete") {
				if runrouter.hasMethod {
					if m, ok := runrouter.methods["delete"]; ok {
						method = vc.MethodByName(m)
					} else if m, ok = runrouter.methods["*"]; ok {
						method = vc.MethodByName(m)
					} else {
						method = vc.MethodByName("Delete")
					}
				} else {
					method = vc.MethodByName("Delete")
				}
				method.Call(in)
			} else if r.Method == "PUT" || (r.Method == "POST" && r.Form.Get("_method") == "put") {
				if runrouter.hasMethod {
					if m, ok := runrouter.methods["put"]; ok {
						method = vc.MethodByName(m)
					} else if m, ok = runrouter.methods["*"]; ok {
						method = vc.MethodByName(m)
					} else {
						method = vc.MethodByName("Put")
					}
				} else {
					method = vc.MethodByName("Put")
				}
				method.Call(in)
			} else if r.Method == "POST" {
				if runrouter.hasMethod {
					if m, ok := runrouter.methods["post"]; ok {
						method = vc.MethodByName(m)
					} else if m, ok = runrouter.methods["*"]; ok {
						method = vc.MethodByName(m)
					} else {
						method = vc.MethodByName("Post")
					}
				} else {
					method = vc.MethodByName("Post")
				}
				method.Call(in)
			} else if r.Method == "PATCH" {
				if runrouter.hasMethod {
					if m, ok := runrouter.methods["patch"]; ok {
						method = vc.MethodByName(m)
					} else if m, ok = runrouter.methods["*"]; ok {
						method = vc.MethodByName(m)
					} else {
						method = vc.MethodByName("Patch")
					}
				} else {
					method = vc.MethodByName("Patch")
				}
				method.Call(in)
			} else if r.Method == "OPTIONS" {
				if runrouter.hasMethod {
					if m, ok := runrouter.methods["options"]; ok {
						method = vc.MethodByName(m)
					} else if m, ok = runrouter.methods["*"]; ok {
						method = vc.MethodByName(m)
					} else {
						method = vc.MethodByName("Options")
					}
				} else {
					method = vc.MethodByName("Options")
				}
				method.Call(in)
			}
			gotofunc := vc.Elem().FieldByName("gotofunc").String()
			if gotofunc != "" {
				method = vc.MethodByName(gotofunc)
				if method.IsValid() {
					method.Call(in)
				} else {
					panic("gotofunc is exists:" + gotofunc)
				}
			}
			if !w.started && !context.Input.IsWebsocket() {
				if AutoRender {
					method = vc.MethodByName("Render")
					method.Call(in)
				}
			}
		}

		method = vc.MethodByName("Finish")
		method.Call(in)
		//execute middleware filters
		if p.enableFilter {
			if l, ok := p.filters[AfterExec]; ok {
				for _, filterR := range l {
					if filterR.ValidRouter(r.URL.Path) {
						filterR.filterFunc(context)
						if w.started {
							goto Admin
						}
					}
				}
			}
		}
		method = vc.MethodByName("Destructor")
		method.Call(in)
	}

	//start autorouter

	if p.enableAuto {
		if !findrouter {
			lastindex := strings.LastIndex(requestPath, "/")
			lastsub := requestPath[lastindex+1:]
			if subindex := strings.LastIndex(lastsub, "."); subindex != -1 {
				context.Input.Param[":ext"] = lastsub[subindex+1:]
				r.URL.Query().Add(":ext", lastsub[subindex+1:])
				r.URL.RawQuery = r.URL.Query().Encode()
				requestPath = requestPath[:len(requestPath)-len(lastsub[subindex:])]
			}
			for cName, methodmap := range p.autoRouter {

				if strings.ToLower(requestPath) == "/"+cName {
					http.Redirect(w, r, requestPath+"/", 301)
					goto Admin
				}

				if strings.ToLower(requestPath) == "/"+cName+"/" {
					requestPath = requestPath + "index"
				}
				if strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/") {
					for mName, controllerType := range methodmap {
						if strings.ToLower(requestPath) == "/"+cName+"/"+strings.ToLower(mName) ||
							(strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/"+strings.ToLower(mName)) &&
								requestPath[len("/"+cName+"/"+strings.ToLower(mName)):len("/"+cName+"/"+strings.ToLower(mName))+1] == "/") {
							if r.Method == "POST" {
								r.ParseMultipartForm(MaxMemory)
							}
							// set find
							findrouter = true
							//execute middleware filters
							if p.enableFilter {
								if l, ok := p.filters[BeforeExec]; ok {
									for _, filterR := range l {
										if filterR.ValidRouter(r.URL.Path) {
											filterR.filterFunc(context)
											if w.started {
												goto Admin
											}
										}
									}
								}
							}
							//parse params
							otherurl := requestPath[len("/"+cName+"/"+strings.ToLower(mName)):]
							if len(otherurl) > 1 {
								plist := strings.Split(otherurl, "/")
								for k, v := range plist[1:] {
									params[strconv.Itoa(k)] = v
								}
							}
							//Invoke the request handler
							vc := reflect.New(controllerType)

							//call the controller init function
							init := vc.MethodByName("Init")
							in := make([]reflect.Value, 3)
							in[0] = reflect.ValueOf(context)
							in[1] = reflect.ValueOf(controllerType.Name())
							in[2] = reflect.ValueOf(vc.Interface())
							init.Call(in)
							//call prepare function
							in = make([]reflect.Value, 0)
							method := vc.MethodByName("Prepare")
							method.Call(in)
							method = vc.MethodByName(mName)
							method.Call(in)
							//if XSRF is Enable then check cookie where there has any cookie in the  request's cookie _csrf
							if EnableXSRF {
								method = vc.MethodByName("XsrfToken")
								method.Call(in)
								if r.Method == "POST" || r.Method == "DELETE" || r.Method == "PUT" ||
									(r.Method == "POST" && (r.Form.Get("_method") == "delete" || r.Form.Get("_method") == "put")) {
									method = vc.MethodByName("CheckXsrfCookie")
									method.Call(in)
								}
							}
							if !w.started && !context.Input.IsWebsocket() {
								if AutoRender {
									method = vc.MethodByName("Render")
									method.Call(in)
								}
							}
							method = vc.MethodByName("Finish")
							method.Call(in)
							//execute middleware filters
							if p.enableFilter {
								if l, ok := p.filters[AfterExec]; ok {
									for _, filterR := range l {
										if filterR.ValidRouter(r.URL.Path) {
											filterR.filterFunc(context)
											if w.started {
												goto Admin
											}
										}
									}
								}
							}
							method = vc.MethodByName("Destructor")
							method.Call(in)
							goto Admin
						}
					}
				}
			}
		}
	}

	//if no matches to url, throw a not found exception
	if !findrouter {
		middleware.Exception("404", rw, r, "")
	}

Admin:
	if p.enableFilter {
		if l, ok := p.filters[FinishRouter]; ok {
			for _, filterR := range l {
				if filterR.ValidRouter(r.URL.Path) {
					filterR.filterFunc(context)
					if w.started {
						break
					}
				}
			}
		}
	}
	//admin module record QPS
	if EnableAdmin {
		timeend := time.Since(starttime)
		if FilterMonitorFunc(r.Method, requestPath, timeend) {
			if runrouter != nil {
				go toolbox.StatisticsMap.AddStatistics(r.Method, requestPath, runrouter.controllerType.Name(), timeend)
			} else {
				go toolbox.StatisticsMap.AddStatistics(r.Method, requestPath, "", timeend)
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
// and sets `started` to true
func (w *responseWriter) Write(p []byte) (int, error) {
	w.started = true
	return w.writer.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.started = true
	w.writer.WriteHeader(code)
}
