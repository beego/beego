package beego

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type handlerInfo struct {
	pattern        string
	regex          *regexp.Regexp
	params         map[int]string
	handlerType reflect.Type
}

type HandlerRegistor struct {
	routers    []*handlerInfo
	fixrouters []*handlerInfo
	filters    []http.HandlerFunc
}

func NewHandlerRegistor() *HandlerRegistor {
	return &HandlerRegistor{routers: make([]*handlerInfo, 0)}
}

func (p *HandlerRegistor) Add(pattern string, c HandlerInterface) {
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
		t := reflect.Indirect(reflect.ValueOf(c)).Type()
		route := &handlerInfo{}
		route.pattern = pattern
		route.handlerType = t

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
		t := reflect.Indirect(reflect.ValueOf(c)).Type()
		route := &handlerInfo{}
		route.regex = regex
		route.params = params
		route.pattern = pattern
		route.handlerType = t

		p.routers = append(p.routers, route)
	}

}

// Filter adds the middleware filter.
func (p *HandlerRegistor) Filter(filter http.HandlerFunc) {
	p.filters = append(p.filters, filter)
}

// FilterParam adds the middleware filter if the REST URL parameter exists.
func (p *HandlerRegistor) FilterParam(param string, filter http.HandlerFunc) {
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
func (p *HandlerRegistor) FilterPrefixPath(path string, filter http.HandlerFunc) {
	p.Filter(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, path) {
			filter(w, r)
		}
	})
}

// AutoRoute
func (p *HandlerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if !RecoverPanic {
				// go back to panic
				panic(err)
			} else {
				Critical("Handler crashed with error", err)
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					Critical(file, line)
				}
			}
		}
	}()
	w := &responseWriter{writer: rw}

	var runrouter *handlerInfo
	var findrouter bool

	params := make(map[string]string)

	//static file server
	for prefix, staticDir := range StaticDir {
		if strings.HasPrefix(r.URL.Path, prefix) {
			file := staticDir + r.URL.Path[len(prefix):]
			http.ServeFile(w, r, file)
			w.started = true
			return
		}
	}

	requestPath := r.URL.Path

	//first find path from the fixrouters to Improve Performance
	for _, route := range p.fixrouters {
		n := len(requestPath)
		//route like "/"
		if n == 1 {
			if requestPath == route.pattern {
				runrouter = route
				findrouter = true
				break
			} else {
				continue
			}
		}

		if (requestPath[n-1] != '/' && route.pattern == requestPath) ||
			(len(route.pattern) >= n-1 && requestPath[0:n-1] == route.pattern) {
			runrouter = route
			findrouter = true
			break
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
					params[route.params[i]] = match
				}
				//reassemble query params and add to RawQuery
				r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
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
		vc := reflect.New(runrouter.handlerType)

		//call the handler init function
		init := vc.MethodByName("Init")
		in := make([]reflect.Value, 2)
		ct := &Context{ResponseWriter: w, Request: r, Params: params}
		in[0] = reflect.ValueOf(ct)
		in[1] = reflect.ValueOf(runrouter.handlerType.Name())
		init.Call(in)
		//call prepare function
		in = make([]reflect.Value, 0)
		method := vc.MethodByName("Prepare")
		method.Call(in)

		//if response has written,yes don't run next
		if !w.started {
			if r.Method == "GET" {
				method = vc.MethodByName("Get")
				method.Call(in)
			} else if r.Method == "POST" {
				method = vc.MethodByName("Post")
				method.Call(in)
			} else if r.Method == "HEAD" {
				method = vc.MethodByName("Head")
				method.Call(in)
			} else if r.Method == "DELETE" {
				method = vc.MethodByName("Delete")
				method.Call(in)
			} else if r.Method == "PUT" {
				method = vc.MethodByName("Put")
				method.Call(in)
			} else if r.Method == "PATCH" {
				method = vc.MethodByName("Patch")
				method.Call(in)
			} else if r.Method == "OPTIONS" {
				method = vc.MethodByName("Options")
				method.Call(in)
			}
			if !w.started {
				if AutoRender {
					method = vc.MethodByName("Render")
					method.Call(in)
				}
				if !w.started {
					method = vc.MethodByName("Finish")
					method.Call(in)
				}
			}
		}
	}

	//if no matches to url, throw a not found exception
	if w.started == false {
		http.NotFound(w, r)
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
