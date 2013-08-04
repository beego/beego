package beego

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
}

type userHandler struct {
	pattern string
	regex   *regexp.Regexp
	params  map[int]string
	h       http.Handler
}

type ControllerRegistor struct {
	routers      []*controllerInfo
	fixrouters   []*controllerInfo
	filters      []http.HandlerFunc
	userHandlers map[string]*userHandler
	autoRouter   map[string]map[string]reflect.Type //key:controller key:method value:reflect.type
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{
		routers:      make([]*controllerInfo, 0),
		userHandlers: make(map[string]*userHandler),
		autoRouter:   make(map[string]map[string]reflect.Type),
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
		route.controllerType = t
		p.routers = append(p.routers, route)
	}
}

func (p *ControllerRegistor) AddAuto(c ControllerInterface) {
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

func (p *ControllerRegistor) AddHandler(pattern string, c http.Handler) {
	parts := strings.Split(pattern, "/")

	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}
	if j == 0 {
		//now create the Route
		uh := &userHandler{}
		uh.pattern = pattern
		uh.h = c
		p.userHandlers[pattern] = uh
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
		uh := &userHandler{}
		uh.regex = regex
		uh.params = params
		uh.pattern = pattern
		uh.h = c
		p.userHandlers[pattern] = uh
	}
}

// Filter adds the middleware filter.
func (p *ControllerRegistor) Filter(filter http.HandlerFunc) {
	p.filters = append(p.filters, filter)
}

// FilterParam adds the middleware filter if the REST URL parameter exists.
func (p *ControllerRegistor) FilterParam(param string, filter http.HandlerFunc) {
	if !strings.HasPrefix(param, ":") {
		param = ":" + param
	}

	p.Filter(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get(param)
		if len(p) > 0 {
			filter(w, r)
		}
	})
}

// FilterPrefixPath adds the middleware filter if the prefix path exists.
func (p *ControllerRegistor) FilterPrefixPath(path string, filter http.HandlerFunc) {
	p.Filter(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, path) {
			filter(w, r)
		}
	})
}

// AutoRoute
func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errstr := fmt.Sprint(err)
			if handler, ok := ErrorMaps[errstr]; ok && ErrorsShow {
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
						ShowErr(err, rw, r, stack)
					}
				}
			}
		}
	}()

	w := &responseWriter{writer: rw}

	w.Header().Set("Server", "beegoServer")
	var runrouter *controllerInfo
	var findrouter bool

	params := make(map[string]string)

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
				if h, ok := ErrorMaps["403"]; ok {
					h(w, r)
					return
				} else {
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.WriteHeader(403)
					fmt.Fprintln(w, "403 Forbidden")
					return
				}
			}
			http.ServeFile(w, r, file)
			w.started = true
			return
		}
	}

	requestPath := r.URL.Path

	var requestbody []byte

	if CopyRequestBody {
		requestbody, _ = ioutil.ReadAll(r.Body)

		r.Body.Close()

		bf := bytes.NewBuffer(requestbody)

		r.Body = ioutil.NopCloser(bf)
	}

	r.ParseMultipartForm(MaxMemory)

	//user defined Handler
	for pattern, c := range p.userHandlers {
		if c.regex == nil && pattern == requestPath {
			c.h.ServeHTTP(rw, r)
			return
		} else if c.regex == nil {
			continue
		}

		//check if Route pattern matches url
		if !c.regex.MatchString(requestPath) {
			continue
		}

		//get submatches (params)
		matches := c.regex.FindStringSubmatch(requestPath)

		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(requestPath) {
			continue
		}

		if len(c.params) > 0 {
			//add url parameters to the query param map
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(c.params[i], match)
				r.Form.Add(c.params[i], match)
				params[c.params[i]] = match
			}
			//reassemble query params and add to RawQuery
			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
			//r.URL.RawQuery = url.Values(values).Encode()
		}
		c.h.ServeHTTP(rw, r)
		return
	}

	//first find path from the fixrouters to Improve Performance
	for _, route := range p.fixrouters {
		n := len(requestPath)
		//route like "/"
		//if n == 1 {
		//	else {
		//		continue
		//	}
		//}
		if requestPath == route.pattern {
			runrouter = route
			findrouter = true
			break
		}
		// pattern /admin   url /admin 200  /admin/ 404
		// pattern /admin/  url /admin 301  /admin/ 200
		if requestPath[n-1] != '/' && len(route.pattern) == n+1 &&
			route.pattern[n] == '/' && route.pattern[:n-1] == requestPath {
			http.Redirect(w, r, requestPath+"/", 301)
			return
		}
	}

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
					r.Form.Add(route.params[i], match)
					params[route.params[i]] = match
				}
				//reassemble query params and add to RawQuery
				r.URL.RawQuery = url.Values(values).Encode()
				//r.URL.RawQuery = url.Values(values).Encode()
			}
			runrouter = route
			findrouter = true
			break
		}
	}

	if runrouter != nil {
		//execute middleware filters
		for _, filter := range p.filters {
			filter(w, r)
			if w.started {
				return
			}
		}

		//Invoke the request handler
		vc := reflect.New(runrouter.controllerType)

		//call the controller init function
		init := vc.MethodByName("Init")
		in := make([]reflect.Value, 2)
		ct := &Context{ResponseWriter: w, Request: r, Params: params, RequestBody: requestbody}

		in[0] = reflect.ValueOf(ct)
		in[1] = reflect.ValueOf(runrouter.controllerType.Name())
		init.Call(in)
		//call prepare function
		in = make([]reflect.Value, 0)
		method := vc.MethodByName("Prepare")
		method.Call(in)

		//if response has written,yes don't run next
		if !w.started {
			if r.Method == "GET" {
				if m, ok := runrouter.methods["get"]; ok {
					method = vc.MethodByName(m)
				} else if m, ok = runrouter.methods["*"]; ok {
					method = vc.MethodByName(m)
				} else {
					method = vc.MethodByName("Get")
				}
				method.Call(in)
			} else if r.Method == "HEAD" {
				if m, ok := runrouter.methods["head"]; ok {
					method = vc.MethodByName(m)
				} else if m, ok = runrouter.methods["*"]; ok {
					method = vc.MethodByName(m)
				} else {
					method = vc.MethodByName("Head")
				}
				method.Call(in)
			} else if r.Method == "DELETE" || (r.Method == "POST" && r.Form.Get("_method") == "delete") {
				if m, ok := runrouter.methods["delete"]; ok {
					method = vc.MethodByName(m)
				} else if m, ok = runrouter.methods["*"]; ok {
					method = vc.MethodByName(m)
				} else {
					method = vc.MethodByName("Delete")
				}
				method.Call(in)
			} else if r.Method == "PUT" || (r.Method == "POST" && r.Form.Get("_method") == "put") {
				if m, ok := runrouter.methods["put"]; ok {
					method = vc.MethodByName(m)
				} else if m, ok = runrouter.methods["*"]; ok {
					method = vc.MethodByName(m)
				} else {
					method = vc.MethodByName("Put")
				}
				method.Call(in)
			} else if r.Method == "POST" {
				if m, ok := runrouter.methods["post"]; ok {
					method = vc.MethodByName(m)
				} else if m, ok = runrouter.methods["*"]; ok {
					method = vc.MethodByName(m)
				} else {
					method = vc.MethodByName("Post")
				}
				method.Call(in)
			} else if r.Method == "PATCH" {
				if m, ok := runrouter.methods["patch"]; ok {
					method = vc.MethodByName(m)
				} else if m, ok = runrouter.methods["*"]; ok {
					method = vc.MethodByName(m)
				} else {
					method = vc.MethodByName("Patch")
				}
				method.Call(in)
			} else if r.Method == "OPTIONS" {
				if m, ok := runrouter.methods["options"]; ok {
					method = vc.MethodByName(m)
				} else if m, ok = runrouter.methods["*"]; ok {
					method = vc.MethodByName(m)
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
			if !w.started {
				if AutoRender {
					method = vc.MethodByName("Render")
					method.Call(in)
				}
				method = vc.MethodByName("Finish")
				method.Call(in)
			}
		}
		method = vc.MethodByName("Destructor")
		method.Call(in)
	}

	//start autorouter

	if !findrouter {
		for cName, methodmap := range p.autoRouter {
			if strings.ToLower(requestPath) == "/"+cName+"/" {
				requestPath = requestPath + "index"
			}
			if strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/") {
				for mName, controllerType := range methodmap {
					if strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/"+strings.ToLower(mName)) {
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
						in := make([]reflect.Value, 2)
						ct := &Context{ResponseWriter: w, Request: r, Params: params, RequestBody: requestbody}

						in[0] = reflect.ValueOf(ct)
						in[1] = reflect.ValueOf(controllerType.Name())
						init.Call(in)
						//call prepare function
						in = make([]reflect.Value, 0)
						method := vc.MethodByName("Prepare")
						method.Call(in)
						method = vc.MethodByName(mName)
						method.Call(in)
						if !w.started {
							if AutoRender {
								method = vc.MethodByName("Render")
								method.Call(in)
							}
							method = vc.MethodByName("Finish")
							method.Call(in)
						}
						method = vc.MethodByName("Destructor")
						method.Call(in)
						// set find
						findrouter = true
					}
				}
			}
		}
	}
	//if no matches to url, throw a not found exception
	if !findrouter {
		if h, ok := ErrorMaps["404"]; ok {
			w.status = 404
			h(w, r)
		} else {
			http.NotFound(w, r)
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
