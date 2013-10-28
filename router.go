package beego

import (
	"fmt"
	beecontext "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/middleware"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
	filters      map[string][]*FilterRouter
	enableAuto   bool
	autoRouter   map[string]map[string]reflect.Type //key:controller key:method value:reflect.type
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{
		routers:    make([]*controllerInfo, 0),
		autoRouter: make(map[string]map[string]reflect.Type),
		filters:    make(map[string][]*FilterRouter),
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
func (p *ControllerRegistor) AddFilter(pattern, action string, filter FilterFunc) {
	p.enableFilter = true
	mr := new(FilterRouter)
	mr.filterFunc = filter

	parts := strings.Split(pattern, "/")
	j := 0
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
			return
		}
		mr.regex = regex
		mr.hasregex = true
	}
	mr.pattern = pattern
	p.filters[action] = append(p.filters[action], mr)
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

	if SessionOn {
		context.Input.CruSession = GlobalSessions.SessionStart(w, r)
	}

	var runrouter *controllerInfo
	var findrouter bool

	params := make(map[string]string)

	if p.enableFilter {
		if l, ok := p.filters["BeforRouter"]; ok {
			for _, filterR := range l {
				if filterR.ValidRouter(r.URL.Path) {
					filterR.filterFunc(context)
					if w.started {
						return
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
			return
		}
		if strings.HasPrefix(r.URL.Path, prefix) {
			file := staticDir + r.URL.Path[len(prefix):]
			finfo, err := os.Stat(file)
			if err != nil {
				return
			}
			//if the request is dir and DirectoryIndex is false then
			if finfo.IsDir() && !DirectoryIndex {
				middleware.Exception("403", rw, r, "403 Forbidden")
				return
			}
			http.ServeFile(w, r, file)
			w.started = true
			return
		}
	}

	if p.enableFilter {
		if l, ok := p.filters["AfterStatic"]; ok {
			for _, filterR := range l {
				if filterR.ValidRouter(r.URL.Path) {
					filterR.filterFunc(context)
					if w.started {
						return
					}
				}
			}
		}
	}
	requestPath := r.URL.Path

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
			return
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
			if l, ok := p.filters["BeforExec"]; ok {
				for _, filterR := range l {
					if filterR.ValidRouter(r.URL.Path) {
						filterR.filterFunc(context)
						if w.started {
							return
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
			if l, ok := p.filters["AfterExec"]; ok {
				for _, filterR := range l {
					if filterR.ValidRouter(r.URL.Path) {
						filterR.filterFunc(context)
						if w.started {
							return
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
					return
				}

				if strings.ToLower(requestPath) == "/"+cName+"/" {
					requestPath = requestPath + "index"
				}
				if strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/") {
					for mName, controllerType := range methodmap {
						if strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/"+strings.ToLower(mName)) {
							if r.Method == "POST" {
								r.ParseMultipartForm(MaxMemory)
							}
							//execute middleware filters
							if p.enableFilter {
								if l, ok := p.filters["BeforExec"]; ok {
									for _, filterR := range l {
										if filterR.ValidRouter(r.URL.Path) {
											filterR.filterFunc(context)
											if w.started {
												return
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
								if l, ok := p.filters["AfterExec"]; ok {
									for _, filterR := range l {
										if filterR.ValidRouter(r.URL.Path) {
											filterR.filterFunc(context)
											if w.started {
												return
											}
										}
									}
								}
							}
							method = vc.MethodByName("Destructor")
							method.Call(in)
							// set find
							findrouter = true
							goto Last
						}
					}
				}
			}
		}
	}

Last:
	//if no matches to url, throw a not found exception
	if !findrouter {
		middleware.Exception("404", rw, r, "")
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
