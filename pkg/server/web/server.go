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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/astaxie/beego/pkg/infrastructure/logs"
	beecontext "github.com/astaxie/beego/pkg/server/web/context"

	"github.com/astaxie/beego/pkg/infrastructure/utils"
	"github.com/astaxie/beego/pkg/server/web/grace"
)

var (
	// BeeApp is an application instance
	// If you are using single server, you could use this
	// But if you need multiple servers, do not use this
	BeeApp *HttpServer
)

func init() {
	// create beego application
	BeeApp = NewHttpSever()
}

// HttpServer defines beego application with a new PatternServeMux.
type HttpServer struct {
	Handlers *ControllerRegister
	Server   *http.Server
	Cfg      *Config
}

// NewHttpSever returns a new beego application.
// this method will use the BConfig as the configure to create HttpServer
// Be careful that when you update BConfig, the server's Cfg will not be updated
func NewHttpSever() *HttpServer {
	return NewHttpServerWithCfg(*BConfig)
}

// NewHttpServerWithCfg will create an sever with specific cfg
func NewHttpServerWithCfg(cfg Config) *HttpServer {
	cfgPtr := &cfg
	cr := NewControllerRegisterWithCfg(cfgPtr)
	app := &HttpServer{
		Handlers: cr,
		Server:   &http.Server{},
		Cfg:      cfgPtr,
	}
	return app
}

// MiddleWare function for http.Handler
type MiddleWare func(http.Handler) http.Handler

// Run beego application.
func (app *HttpServer) Run(addr string, mws ...MiddleWare) {

	initBeforeHTTPRun()

	app.initAddr(addr)

	addr = app.Cfg.Listen.HTTPAddr

	if app.Cfg.Listen.HTTPPort != 0 {
		addr = fmt.Sprintf("%s:%d", app.Cfg.Listen.HTTPAddr, app.Cfg.Listen.HTTPPort)
	}

	var (
		err        error
		l          net.Listener
		endRunning = make(chan bool, 1)
	)

	// run cgi server
	if app.Cfg.Listen.EnableFcgi {
		if app.Cfg.Listen.EnableStdIo {
			if err = fcgi.Serve(nil, app.Handlers); err == nil { // standard I/O
				logs.Info("Use FCGI via standard I/O")
			} else {
				logs.Critical("Cannot use FCGI via standard I/O", err)
			}
			return
		}
		if app.Cfg.Listen.HTTPPort == 0 {
			// remove the Socket file before start
			if utils.FileExists(addr) {
				os.Remove(addr)
			}
			l, err = net.Listen("unix", addr)
		} else {
			l, err = net.Listen("tcp", addr)
		}
		if err != nil {
			logs.Critical("Listen: ", err)
		}
		if err = fcgi.Serve(l, app.Handlers); err != nil {
			logs.Critical("fcgi.Serve: ", err)
		}
		return
	}

	app.Server.Handler = app.Handlers
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] == nil {
			continue
		}
		app.Server.Handler = mws[i](app.Server.Handler)
	}
	app.Server.ReadTimeout = time.Duration(app.Cfg.Listen.ServerTimeOut) * time.Second
	app.Server.WriteTimeout = time.Duration(app.Cfg.Listen.ServerTimeOut) * time.Second
	app.Server.ErrorLog = logs.GetLogger("HTTP")

	// run graceful mode
	if app.Cfg.Listen.Graceful {
		httpsAddr := app.Cfg.Listen.HTTPSAddr
		app.Server.Addr = httpsAddr
		if app.Cfg.Listen.EnableHTTPS || app.Cfg.Listen.EnableMutualHTTPS {
			go func() {
				time.Sleep(1000 * time.Microsecond)
				if app.Cfg.Listen.HTTPSPort != 0 {
					httpsAddr = fmt.Sprintf("%s:%d", app.Cfg.Listen.HTTPSAddr, app.Cfg.Listen.HTTPSPort)
					app.Server.Addr = httpsAddr
				}
				server := grace.NewServer(httpsAddr, app.Server.Handler)
				server.Server.ReadTimeout = app.Server.ReadTimeout
				server.Server.WriteTimeout = app.Server.WriteTimeout
				if app.Cfg.Listen.EnableMutualHTTPS {
					if err := server.ListenAndServeMutualTLS(app.Cfg.Listen.HTTPSCertFile,
						app.Cfg.Listen.HTTPSKeyFile,
						app.Cfg.Listen.TrustCaFile); err != nil {
						logs.Critical("ListenAndServeTLS: ", err, fmt.Sprintf("%d", os.Getpid()))
						time.Sleep(100 * time.Microsecond)
					}
				} else {
					if app.Cfg.Listen.AutoTLS {
						m := autocert.Manager{
							Prompt:     autocert.AcceptTOS,
							HostPolicy: autocert.HostWhitelist(app.Cfg.Listen.Domains...),
							Cache:      autocert.DirCache(app.Cfg.Listen.TLSCacheDir),
						}
						app.Server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
						app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile = "", ""
					}
					if err := server.ListenAndServeTLS(app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile); err != nil {
						logs.Critical("ListenAndServeTLS: ", err, fmt.Sprintf("%d", os.Getpid()))
						time.Sleep(100 * time.Microsecond)
					}
				}
				endRunning <- true
			}()
		}
		if app.Cfg.Listen.EnableHTTP {
			go func() {
				server := grace.NewServer(addr, app.Server.Handler)
				server.Server.ReadTimeout = app.Server.ReadTimeout
				server.Server.WriteTimeout = app.Server.WriteTimeout
				if app.Cfg.Listen.ListenTCP4 {
					server.Network = "tcp4"
				}
				if err := server.ListenAndServe(); err != nil {
					logs.Critical("ListenAndServe: ", err, fmt.Sprintf("%d", os.Getpid()))
					time.Sleep(100 * time.Microsecond)
				}
				endRunning <- true
			}()
		}
		<-endRunning
		return
	}

	// run normal mode
	if app.Cfg.Listen.EnableHTTPS || app.Cfg.Listen.EnableMutualHTTPS {
		go func() {
			time.Sleep(1000 * time.Microsecond)
			if app.Cfg.Listen.HTTPSPort != 0 {
				app.Server.Addr = fmt.Sprintf("%s:%d", app.Cfg.Listen.HTTPSAddr, app.Cfg.Listen.HTTPSPort)
			} else if app.Cfg.Listen.EnableHTTP {
				logs.Info("Start https server error, conflict with http. Please reset https port")
				return
			}
			logs.Info("https server Running on https://%s", app.Server.Addr)
			if app.Cfg.Listen.AutoTLS {
				m := autocert.Manager{
					Prompt:     autocert.AcceptTOS,
					HostPolicy: autocert.HostWhitelist(app.Cfg.Listen.Domains...),
					Cache:      autocert.DirCache(app.Cfg.Listen.TLSCacheDir),
				}
				app.Server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
				app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile = "", ""
			} else if app.Cfg.Listen.EnableMutualHTTPS {
				pool := x509.NewCertPool()
				data, err := ioutil.ReadFile(app.Cfg.Listen.TrustCaFile)
				if err != nil {
					logs.Info("MutualHTTPS should provide TrustCaFile")
					return
				}
				pool.AppendCertsFromPEM(data)
				app.Server.TLSConfig = &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.ClientAuthType(app.Cfg.Listen.ClientAuth),
				}
			}
			if err := app.Server.ListenAndServeTLS(app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile); err != nil {
				logs.Critical("ListenAndServeTLS: ", err)
				time.Sleep(100 * time.Microsecond)
				endRunning <- true
			}
		}()

	}
	if app.Cfg.Listen.EnableHTTP {
		go func() {
			app.Server.Addr = addr
			logs.Info("http server Running on http://%s", app.Server.Addr)
			if app.Cfg.Listen.ListenTCP4 {
				ln, err := net.Listen("tcp4", app.Server.Addr)
				if err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
					return
				}
				if err = app.Server.Serve(ln); err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
					return
				}
			} else {
				if err := app.Server.ListenAndServe(); err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
				}
			}
		}()
	}
	<-endRunning
}

func (app *HttpServer) Start() {

}

// Router see HttpServer.Router
func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *HttpServer {
	return BeeApp.Router(rootpath, c, mappingMethods...)
}

// Router adds a patterned controller handler to BeeApp.
// it's an alias method of HttpServer.Router.
// usage:
//  simple router
//  beego.Router("/admin", &admin.UserController{})
//  beego.Router("/admin/index", &admin.ArticleController{})
//
//  regex router
//
//  beego.Router("/api/:id([0-9]+)", &controllers.RController{})
//
//  custom rules
//  beego.Router("/api/list",&RestController{},"*:ListFood")
//  beego.Router("/api/create",&RestController{},"post:CreateFood")
//  beego.Router("/api/update",&RestController{},"put:UpdateFood")
//  beego.Router("/api/delete",&RestController{},"delete:DeleteFood")
func (app *HttpServer) Router(rootPath string, c ControllerInterface, mappingMethods ...string) *HttpServer {
	app.Handlers.Add(rootPath, c, mappingMethods...)
	return app
}

// UnregisterFixedRoute see HttpServer.UnregisterFixedRoute
func UnregisterFixedRoute(fixedRoute string, method string) *HttpServer {
	return BeeApp.UnregisterFixedRoute(fixedRoute, method)
}

// UnregisterFixedRoute unregisters the route with the specified fixedRoute. It is particularly useful
// in web applications that inherit most routes from a base webapp via the underscore
// import, and aim to overwrite only certain paths.
// The method parameter can be empty or "*" for all HTTP methods, or a particular
// method type (e.g. "GET" or "POST") for selective removal.
//
// Usage (replace "GET" with "*" for all methods):
//  beego.UnregisterFixedRoute("/yourpreviouspath", "GET")
//  beego.Router("/yourpreviouspath", yourControllerAddress, "get:GetNewPage")
func (app *HttpServer) UnregisterFixedRoute(fixedRoute string, method string) *HttpServer {
	subPaths := splitPath(fixedRoute)
	if method == "" || method == "*" {
		for m := range HTTPMETHOD {
			if _, ok := app.Handlers.routers[m]; !ok {
				continue
			}
			if app.Handlers.routers[m].prefix == strings.Trim(fixedRoute, "/ ") {
				findAndRemoveSingleTree(app.Handlers.routers[m])
				continue
			}
			findAndRemoveTree(subPaths, app.Handlers.routers[m], m)
		}
		return app
	}
	// Single HTTP method
	um := strings.ToUpper(method)
	if _, ok := app.Handlers.routers[um]; ok {
		if app.Handlers.routers[um].prefix == strings.Trim(fixedRoute, "/ ") {
			findAndRemoveSingleTree(app.Handlers.routers[um])
			return app
		}
		findAndRemoveTree(subPaths, app.Handlers.routers[um], um)
	}
	return app
}

func findAndRemoveTree(paths []string, entryPointTree *Tree, method string) {
	for i := range entryPointTree.fixrouters {
		if entryPointTree.fixrouters[i].prefix == paths[0] {
			if len(paths) == 1 {
				if len(entryPointTree.fixrouters[i].fixrouters) > 0 {
					// If the route had children subtrees, remove just the functional leaf,
					// to allow children to function as before
					if len(entryPointTree.fixrouters[i].leaves) > 0 {
						entryPointTree.fixrouters[i].leaves[0] = nil
						entryPointTree.fixrouters[i].leaves = entryPointTree.fixrouters[i].leaves[1:]
					}
				} else {
					// Remove the *Tree from the fixrouters slice
					entryPointTree.fixrouters[i] = nil

					if i == len(entryPointTree.fixrouters)-1 {
						entryPointTree.fixrouters = entryPointTree.fixrouters[:i]
					} else {
						entryPointTree.fixrouters = append(entryPointTree.fixrouters[:i], entryPointTree.fixrouters[i+1:len(entryPointTree.fixrouters)]...)
					}
				}
				return
			}
			findAndRemoveTree(paths[1:], entryPointTree.fixrouters[i], method)
		}
	}
}

func findAndRemoveSingleTree(entryPointTree *Tree) {
	if entryPointTree == nil {
		return
	}
	if len(entryPointTree.fixrouters) > 0 {
		// If the route had children subtrees, remove just the functional leaf,
		// to allow children to function as before
		if len(entryPointTree.leaves) > 0 {
			entryPointTree.leaves[0] = nil
			entryPointTree.leaves = entryPointTree.leaves[1:]
		}
	}
}

// Include see HttpServer.Include
func Include(cList ...ControllerInterface) *HttpServer {
	return BeeApp.Include(cList...)
}

// Include will generate router file in the router/xxx.go from the controller's comments
// usage:
// beego.Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
// type BankAccount struct{
//   beego.Controller
// }
//
// register the function
// func (b *BankAccount)Mapping(){
//  b.Mapping("ShowAccount" , b.ShowAccount)
//  b.Mapping("ModifyAccount", b.ModifyAccount)
// }
//
// //@router /account/:id  [get]
// func (b *BankAccount) ShowAccount(){
//    //logic
// }
//
//
// //@router /account/:id  [post]
// func (b *BankAccount) ModifyAccount(){
//    //logic
// }
//
// the comments @router url methodlist
// url support all the function Router's pattern
// methodlist [get post head put delete options *]
func (app *HttpServer) Include(cList ...ControllerInterface) *HttpServer {
	app.Handlers.Include(cList...)
	return app
}

// RESTRouter see HttpServer.RESTRouter
func RESTRouter(rootpath string, c ControllerInterface) *HttpServer {
	return BeeApp.RESTRouter(rootpath, c)
}

// RESTRouter adds a restful controller handler to BeeApp.
// its' controller implements beego.ControllerInterface and
// defines a param "pattern/:objectId" to visit each resource.
func (app *HttpServer) RESTRouter(rootpath string, c ControllerInterface) *HttpServer {
	app.Router(rootpath, c)
	app.Router(path.Join(rootpath, ":objectId"), c)
	return app
}

// AutoRouter see HttpServer.AutoRouter
func AutoRouter(c ControllerInterface) *HttpServer {
	return BeeApp.AutoRouter(c)
}

// AutoRouter adds defined controller handler to BeeApp.
// it's same to HttpServer.AutoRouter.
// if beego.AddAuto(&MainContorlller{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func (app *HttpServer) AutoRouter(c ControllerInterface) *HttpServer {
	app.Handlers.AddAuto(c)
	return app
}

// AutoPrefix see HttpServer.AutoPrefix
func AutoPrefix(prefix string, c ControllerInterface) *HttpServer {
	return BeeApp.AutoPrefix(prefix, c)
}

// AutoPrefix adds controller handler to BeeApp with prefix.
// it's same to HttpServer.AutoRouterWithPrefix.
// if beego.AutoPrefix("/admin",&MainContorlller{}) and MainController has methods List and Page,
// visit the url /admin/main/list to exec List function or /admin/main/page to exec Page function.
func (app *HttpServer) AutoPrefix(prefix string, c ControllerInterface) *HttpServer {
	app.Handlers.AddAutoPrefix(prefix, c)
	return app
}

// Get see HttpServer.Get
func Get(rootpath string, f FilterFunc) *HttpServer {
	return BeeApp.Get(rootpath, f)
}

// Get used to register router for Get method
// usage:
//    beego.Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Get(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Get(rootpath, f)
	return app
}

// Post see HttpServer.Post
func Post(rootpath string, f FilterFunc) *HttpServer {
	return BeeApp.Post(rootpath, f)
}

// Post used to register router for Post method
// usage:
//    beego.Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Post(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Post(rootpath, f)
	return app
}

// Delete see HttpServer.Delete
func Delete(rootpath string, f FilterFunc) *HttpServer {
	return BeeApp.Delete(rootpath, f)
}

// Delete used to register router for Delete method
// usage:
//    beego.Delete("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Delete(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Delete(rootpath, f)
	return app
}

// Put see HttpServer.Put
func Put(rootpath string, f FilterFunc) *HttpServer {
	return BeeApp.Put(rootpath, f)
}

// Put used to register router for Put method
// usage:
//    beego.Put("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Put(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Put(rootpath, f)
	return app
}

// Head see HttpServer.Head
func Head(rootpath string, f FilterFunc) *HttpServer {
	return BeeApp.Head(rootpath, f)
}

// Head used to register router for Head method
// usage:
//    beego.Head("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Head(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Head(rootpath, f)
	return app
}

// Options see HttpServer.Options
func Options(rootpath string, f FilterFunc) *HttpServer {
	BeeApp.Handlers.Options(rootpath, f)
	return BeeApp
}

// Options used to register router for Options method
// usage:
//    beego.Options("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Options(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Options(rootpath, f)
	return app
}

// Patch see HttpServer.Patch
func Patch(rootpath string, f FilterFunc) *HttpServer {
	return BeeApp.Patch(rootpath, f)
}

// Patch used to register router for Patch method
// usage:
//    beego.Patch("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Patch(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Patch(rootpath, f)
	return app
}

// Any see HttpServer.Any
func Any(rootpath string, f FilterFunc) *HttpServer {
	return BeeApp.Any(rootpath, f)
}

// Any used to register router for all methods
// usage:
//    beego.Any("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Any(rootpath string, f FilterFunc) *HttpServer {
	app.Handlers.Any(rootpath, f)
	return app
}

// Handler see HttpServer.Handler
func Handler(rootpath string, h http.Handler, options ...interface{}) *HttpServer {
	return BeeApp.Handler(rootpath, h, options...)
}

// Handler used to register a Handler router
// usage:
//    beego.Handler("/api", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
//          fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//    }))
func (app *HttpServer) Handler(rootpath string, h http.Handler, options ...interface{}) *HttpServer {
	app.Handlers.Handler(rootpath, h, options...)
	return app
}

// InserFilter see HttpServer.InsertFilter
func InsertFilter(pattern string, pos int, filter FilterFunc, opts ...FilterOpt) *HttpServer {
	return BeeApp.InsertFilter(pattern, pos, filter, opts...)
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// beego.BeforeStatic, beego.BeforeRouter, beego.BeforeExec, beego.AfterExec and beego.FinishRouter.
// The bool params is for setting the returnOnOutput value (false allows multiple filters to execute)
func (app *HttpServer) InsertFilter(pattern string, pos int, filter FilterFunc, opts ...FilterOpt) *HttpServer {
	app.Handlers.InsertFilter(pattern, pos, filter, opts...)
	return app
}

// InsertFilterChain see HttpServer.InsertFilterChain
func InsertFilterChain(pattern string, filterChain FilterChain, opts ...FilterOpt) *HttpServer {
	return BeeApp.InsertFilterChain(pattern, filterChain, opts...)
}

// InsertFilterChain adds a FilterFunc built by filterChain.
// This filter will be executed before all filters.
// the filter's behavior like stack's behavior
// and the last filter is serving the http request
func (app *HttpServer) InsertFilterChain(pattern string, filterChain FilterChain, opts ...FilterOpt) *HttpServer {
	app.Handlers.InsertFilterChain(pattern, filterChain, opts...)
	return app
}

func (app *HttpServer) initAddr(addr string) {
	strs := strings.Split(addr, ":")
	if len(strs) > 0 && strs[0] != "" {
		app.Cfg.Listen.HTTPAddr = strs[0]
		app.Cfg.Listen.Domains = []string{strs[0]}
	}
	if len(strs) > 1 && strs[1] != "" {
		app.Cfg.Listen.HTTPPort, _ = strconv.Atoi(strs[1])
	}
}

func (app *HttpServer) LogAccess(ctx *beecontext.Context, startTime *time.Time, statusCode int) {
	// Skip logging if AccessLogs config is false
	if !app.Cfg.Log.AccessLogs {
		return
	}
	// Skip logging static requests unless EnableStaticLogs config is true
	if !app.Cfg.Log.EnableStaticLogs && DefaultAccessLogFilter.Filter(ctx) {
		return
	}
	var (
		requestTime time.Time
		elapsedTime time.Duration
		r           = ctx.Request
	)
	if startTime != nil {
		requestTime = *startTime
		elapsedTime = time.Since(*startTime)
	}
	record := &logs.AccessLogRecord{
		RemoteAddr:     ctx.Input.IP(),
		RequestTime:    requestTime,
		RequestMethod:  r.Method,
		Request:        fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		ServerProtocol: r.Proto,
		Host:           r.Host,
		Status:         statusCode,
		ElapsedTime:    elapsedTime,
		HTTPReferrer:   r.Header.Get("Referer"),
		HTTPUserAgent:  r.Header.Get("User-Agent"),
		RemoteUser:     r.Header.Get("Remote-User"),
		BodyBytesSent:  r.ContentLength,
	}
	logs.AccessLog(record, app.Cfg.Log.AccessLogsFormat)
}
