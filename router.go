package beego

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
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
	BeforeRouter = iota
	AfterStatic
	BeforeExec
	AfterExec
	FinishRouter
)

var (
	// supported http methods.
	HTTPMETHOD = []string{"get", "post", "put", "delete", "patch", "options", "head"}
)

type controllerInfo struct {
	pattern        string
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
	methods        map[string]string
	hasMethod      bool
}

// ControllerRegistor containers registered router rules, controller handlers and filters.
type ControllerRegistor struct {
	routers      []*controllerInfo // regexp router storage
	fixrouters   []*controllerInfo // fixed router storage
	enableFilter bool
	filters      map[int][]*FilterRouter
	enableAuto   bool
	autoRouter   map[string]map[string]reflect.Type //key:controller key:method value:reflect.type
}

// NewControllerRegistor returns a new ControllerRegistor.
func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{
		routers:    make([]*controllerInfo, 0),
		autoRouter: make(map[string]map[string]reflect.Type),
		filters:    make(map[int][]*FilterRouter),
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
				panic("method mapping format is invalid")
			}
			comma := strings.Split(colon[0], ",")
			for _, m := range comma {
				if m == "*" || utils.InSlice(strings.ToLower(m), HTTPMETHOD) {
					if val := reflectVal.MethodByName(colon[1]); val.IsValid() {
						methods[strings.ToLower(m)] = colon[1]
					} else {
						panic(colon[1] + " method doesn't exist in the controller " + t.Name())
					}
				} else {
					panic(v + " is an invalid method mapping. Method doesn't exist " + m)
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

// Add auto router to ControllerRegistor.
// example beego.AddAuto(&MainContorlller{}),
// MainController has method List and Page.
// visit the url /main/list to exec List function
// /main/page to exec Page function.
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

// [Deprecated] use InsertFilter.
// Add FilterFunc with pattern for action.
func (p *ControllerRegistor) AddFilter(pattern, action string, filter FilterFunc) {
	mr := buildFilter(pattern, filter)
	switch action {
	case "BeforeRouter":
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

// Add a FilterFunc with pattern rule and action constant.
func (p *ControllerRegistor) InsertFilter(pattern string, pos int, filter FilterFunc) {
	mr := buildFilter(pattern, filter)
	p.filters[pos] = append(p.filters[pos], mr)
	p.enableFilter = true
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
			if utils.InSlice(strings.ToLower(methodName), HTTPMETHOD) {
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
			if utils.InSlice(strings.ToLower(methodName), HTTPMETHOD) {
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

// Implement http.Handler interface.
func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if err == USERSTOPRUN {
				return
			}
			if _, ok := err.(middleware.HTTPException); ok {
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
							Critical(file, line)
							stack = stack + fmt.Sprintln(file, line)
						}
						middleware.ShowErr(err, rw, r, stack)
					}
				} else {
					if !RecoverPanic {
						panic(err)
					} else {
						// in production model show all infomation
						if ErrorsShow {
							handler := p.getErrorHandler(fmt.Sprint(err))
							handler(rw, r)
							return
						} else {
							Critical("the request url is ", r.URL.Path)
							Critical("Handler crashed with error", err)
							for i := 1; ; i++ {
								_, file, line, ok := runtime.Caller(i)
								if !ok {
									break
								}
								Critical(file, line)
							}
						}
					}
				}

			}
		}
	}()

	starttime := time.Now()
	requestPath := r.URL.Path
	var runrouter reflect.Type
	var findrouter bool
	var runMethod string
	params := make(map[string]string)

	w := &responseWriter{writer: rw}
	w.Header().Set("Server", BeegoServerName)

	// init context
	context := &beecontext.Context{
		ResponseWriter: w,
		Request:        r,
		Input:          beecontext.NewInput(r),
		Output:         beecontext.NewOutput(),
	}
	context.Output.Context = context
	context.Output.EnableGzip = EnableGzip

	if context.Input.IsWebsocket() {
		context.ResponseWriter = rw
	}

	// defined filter function
	do_filter := func(pos int) (started bool) {
		if p.enableFilter {
			if l, ok := p.filters[pos]; ok {
				for _, filterR := range l {
					if ok, p := filterR.ValidRouter(r.URL.Path); ok {
						context.Input.Params = p
						filterR.filterFunc(context)
						if w.started {
							return true
						}
					}
				}
			}
		}

		return false
	}

	// session init
	if SessionOn {
		context.Input.CruSession = GlobalSessions.SessionStart(w, r)
		defer context.Input.CruSession.SessionRelease()
	}

	if !utils.InSlice(strings.ToLower(r.Method), HTTPMETHOD) {
		http.Error(w, "Method Not Allowed", 405)
		goto Admin
	}

	if do_filter(BeforeRouter) {
		goto Admin
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

			//This block obtained from (https://github.com/smithfox/beego) - it should probably get merged into astaxie/beego after a pull request
			isStaticFileToCompress := false
			if StaticExtensionsToGzip != nil && len(StaticExtensionsToGzip) > 0 {
				for _, statExtension := range StaticExtensionsToGzip {
					if strings.HasSuffix(strings.ToLower(file), strings.ToLower(statExtension)) {
						isStaticFileToCompress = true
						break
					}
				}
			}

			if isStaticFileToCompress {
				if EnableGzip {
					w.contentEncoding = GetAcceptEncodingZip(r)
				}

				memzipfile, err := OpenMemZipFile(file, w.contentEncoding)
				if err != nil {
					return
				}

				w.InitHeadContent(finfo.Size())

				http.ServeContent(w, r, file, finfo.ModTime(), memzipfile)
			} else {
				http.ServeFile(w, r, file)
			}

			w.started = true
			goto Admin
		}
	}

	if do_filter(AfterStatic) {
		goto Admin
	}

	if CopyRequestBody {
		context.Input.Body()
	}

	//first find path from the fixrouters to Improve Performance
	for _, route := range p.fixrouters {
		n := len(requestPath)
		if requestPath == route.pattern {
			runMethod = p.getRunMethod(r.Method, context, route)
			if runMethod != "" {
				runrouter = route.controllerType
				findrouter = true
				break
			}
		}
		// pattern /admin   url /admin 200  /admin/ 200
		// pattern /admin/  url /admin 301  /admin/ 200
		if requestPath[n-1] != '/' && len(route.pattern) == n+1 &&
			route.pattern[n] == '/' && route.pattern[:n] == requestPath {
			http.Redirect(w, r, requestPath+"/", 301)
			goto Admin
		}
		if requestPath[n-1] == '/' && n >= 2 && requestPath[:n-2] == route.pattern {
			runMethod = p.getRunMethod(r.Method, context, route)
			if runMethod != "" {
				runrouter = route.controllerType
				findrouter = true
				break
			}
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
			runMethod = p.getRunMethod(r.Method, context, route)
			if runMethod != "" {
				runrouter = route.controllerType
				context.Input.Params = params
				findrouter = true
				break
			}
		}
	}

	if !findrouter && p.enableAuto {
		// deal with url with diffirent ext
		// /controller/simple
		// /controller/simple.html
		// /controller/simple.json
		// /controller/simple.rss
		lastindex := strings.LastIndex(requestPath, "/")
		lastsub := requestPath[lastindex+1:]
		if subindex := strings.LastIndex(lastsub, "."); subindex != -1 {
			context.Input.Params[":ext"] = lastsub[subindex+1:]
			r.URL.Query().Add(":ext", lastsub[subindex+1:])
			r.URL.RawQuery = r.URL.Query().Encode()
			requestPath = requestPath[:len(requestPath)-len(lastsub[subindex:])]
		}
		for cName, methodmap := range p.autoRouter {
			// if prev already find the router break
			if findrouter {
				break
			}
			if strings.ToLower(requestPath) == "/"+cName {
				http.Redirect(w, r, requestPath+"/", 301)
				goto Admin
			}
			// if there's no action, set the default action to index
			if strings.ToLower(requestPath) == "/"+cName+"/" {
				requestPath = requestPath + "index"
			}
			// if the request path start with controllerName
			if strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/") {
				for mName, controllerType := range methodmap {
					if strings.ToLower(requestPath) == "/"+cName+"/"+strings.ToLower(mName) ||
						(strings.HasPrefix(strings.ToLower(requestPath), "/"+cName+"/"+strings.ToLower(mName)) &&
							requestPath[len("/"+cName+"/"+strings.ToLower(mName)):len("/"+cName+"/"+strings.ToLower(mName))+1] == "/") {
						runrouter = controllerType
						runMethod = mName
						findrouter = true
						//parse params
						otherurl := requestPath[len("/"+cName+"/"+strings.ToLower(mName)):]
						if len(otherurl) > 1 {
							plist := strings.Split(otherurl, "/")
							for k, v := range plist[1:] {
								context.Input.Params[strconv.Itoa(k)] = v
							}
						}
						break
					}
				}
			}
		}
	}

	//if no matches to url, throw a not found exception
	if !findrouter {
		middleware.Exception("404", rw, r, "")
		goto Admin
	}

	if findrouter {
		if r.Method == "POST" {
			r.ParseMultipartForm(MaxMemory)
		}
		//execute middleware filters
		if do_filter(BeforeExec) {
			goto Admin
		}

		//Invoke the request handler
		vc := reflect.New(runrouter)
		execController, ok := vc.Interface().(ControllerInterface)
		if !ok {
			panic("controller is not ControllerInterface")
		}

		//call the controller init function
		execController.Init(context, runrouter.Name(), runMethod, vc.Interface())

		//if XSRF is Enable then check cookie where there has any cookie in the  request's cookie _csrf
		if EnableXSRF {
			execController.XsrfToken()
			if r.Method == "POST" || r.Method == "DELETE" || r.Method == "PUT" ||
				(r.Method == "POST" && (r.Form.Get("_method") == "delete" || r.Form.Get("_method") == "put")) {
				execController.CheckXsrfCookie()
			}
		}

		//call prepare function
		execController.Prepare()

		if !w.started {
			//exec main logic
			switch runMethod {
			case "Get":
				execController.Get()
			case "Post":
				execController.Post()
			case "Delete":
				execController.Delete()
			case "Put":
				execController.Put()
			case "Head":
				execController.Head()
			case "Patch":
				execController.Patch()
			case "Options":
				execController.Options()
			default:
				in := make([]reflect.Value, 0)
				method := vc.MethodByName(runMethod)
				method.Call(in)
			}

			//render template
			if !w.started && !context.Input.IsWebsocket() {
				if AutoRender {
					if err := execController.Render(); err != nil {
						panic(err)
					}

				}
			}
		}

		// finish all runrouter. release resource
		execController.Finish()

		//execute middleware filters
		if do_filter(AfterExec) {
			goto Admin
		}
	}

Admin:
	do_filter(FinishRouter)

	//admin module record QPS
	if EnableAdmin {
		timeend := time.Since(starttime)
		if FilterMonitorFunc(r.Method, requestPath, timeend) {
			if runrouter != nil {
				go toolbox.StatisticsMap.AddStatistics(r.Method, requestPath, runrouter.Name(), timeend)
			} else {
				go toolbox.StatisticsMap.AddStatistics(r.Method, requestPath, "", timeend)
			}
		}
	}
}

// there always should be error handler that sets error code accordingly for all unhandled errors.
// in order to have custom UI for error page it's necessary to override "500" error.
func (p *ControllerRegistor) getErrorHandler(errorCode string) func(rw http.ResponseWriter, r *http.Request) {
	handler := middleware.SimpleServerError
	ok := true
	if errorCode != "" {
		handler, ok = middleware.ErrorMaps[errorCode]
		if !ok {
			handler, ok = middleware.ErrorMaps["500"]
		}
		if !ok || handler == nil {
			handler = middleware.SimpleServerError
		}
	}

	return handler
}

// returns method name from request header or form field.
// sometimes browsers can't create PUT and DELETE request.
// set a form field "_method" instead.
func (p *ControllerRegistor) getRunMethod(method string, context *beecontext.Context, router *controllerInfo) string {
	method = strings.ToLower(method)
	if method == "post" && strings.ToLower(context.Input.Query("_method")) == "put" {
		method = "put"
	}
	if method == "post" && strings.ToLower(context.Input.Query("_method")) == "delete" {
		method = "delete"
	}
	if router.hasMethod {
		if m, ok := router.methods[method]; ok {
			return m
		} else if m, ok = router.methods["*"]; ok {
			return m
		} else {
			return ""
		}
	} else {
		return strings.Title(method)
	}
}

//responseWriter is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type responseWriter struct {
	writer          http.ResponseWriter
	started         bool
	status          int
	contentEncoding string
}

// Header returns the header map that will be sent by WriteHeader.
func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

// Init content-length header.
func (w *responseWriter) InitHeadContent(contentlength int64) {
	if w.contentEncoding == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
	} else if w.contentEncoding == "deflate" {
		w.Header().Set("Content-Encoding", "deflate")
	} else {
		w.Header().Set("Content-Length", strconv.FormatInt(contentlength, 10))
	}
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
