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
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
	"time"

	"github.com/astaxie/beego/grace"
	"github.com/astaxie/beego/utils"
)

var (
	// BeeApp is an application instance
	BeeApp *App
)

func init() {
	// create beego application
	BeeApp = NewApp()
}

// App defines beego application with a new PatternServeMux.
type App struct {
	Handlers *ControllerRegister
	Server   *http.Server
}

// NewApp returns a new beego application.
func NewApp() *App {
	cr := NewControllerRegister()
	app := &App{Handlers: cr, Server: &http.Server{}}
	return app
}

// Run beego application.
func (app *App) Run() {
	addr := HttpAddr

	if HttpPort != 0 {
		addr = fmt.Sprintf("%s:%d", HttpAddr, HttpPort)
	}

	var (
		err error
		l   net.Listener
	)
	endRunning := make(chan bool, 1)

	if UseFcgi {
		if UseStdIo {
			err = fcgi.Serve(nil, app.Handlers) // standard I/O
			if err == nil {
				BeeLogger.Info("Use FCGI via standard I/O")
			} else {
				BeeLogger.Info("Cannot use FCGI via standard I/O", err)
			}
		} else {
			if HttpPort == 0 {
				// remove the Socket file before start
				if utils.FileExists(addr) {
					os.Remove(addr)
				}
				l, err = net.Listen("unix", addr)
			} else {
				l, err = net.Listen("tcp", addr)
			}
			if err != nil {
				BeeLogger.Critical("Listen: ", err)
			}
			err = fcgi.Serve(l, app.Handlers)
		}
	} else {
		if Graceful {
			app.Server.Addr = addr
			app.Server.Handler = app.Handlers
			app.Server.ReadTimeout = time.Duration(HttpServerTimeOut) * time.Second
			app.Server.WriteTimeout = time.Duration(HttpServerTimeOut) * time.Second
			if EnableHttpTLS {
				go func() {
					time.Sleep(20 * time.Microsecond)
					if HttpsPort != 0 {
						addr = fmt.Sprintf("%s:%d", HttpAddr, HttpsPort)
						app.Server.Addr = addr
					}
					server := grace.NewServer(addr, app.Handlers)
					server.Server = app.Server
					err := server.ListenAndServeTLS(HttpCertFile, HttpKeyFile)
					if err != nil {
						BeeLogger.Critical("ListenAndServeTLS: ", err, fmt.Sprintf("%d", os.Getpid()))
						time.Sleep(100 * time.Microsecond)
						endRunning <- true
					}
				}()
			}
			if EnableHttpListen {
				go func() {
					server := grace.NewServer(addr, app.Handlers)
					server.Server = app.Server
					if ListenTCP4 && HttpAddr == "" {
						server.Network = "tcp4"
					}
					err := server.ListenAndServe()
					if err != nil {
						BeeLogger.Critical("ListenAndServe: ", err, fmt.Sprintf("%d", os.Getpid()))
						time.Sleep(100 * time.Microsecond)
						endRunning <- true
					}
				}()
			}
		} else {
			app.Server.Addr = addr
			app.Server.Handler = app.Handlers
			app.Server.ReadTimeout = time.Duration(HttpServerTimeOut) * time.Second
			app.Server.WriteTimeout = time.Duration(HttpServerTimeOut) * time.Second

			if EnableHttpTLS {
				go func() {
					time.Sleep(20 * time.Microsecond)
					if HttpsPort != 0 {
						app.Server.Addr = fmt.Sprintf("%s:%d", HttpAddr, HttpsPort)
					}
					BeeLogger.Info("https server Running on %s", app.Server.Addr)
					err := app.Server.ListenAndServeTLS(HttpCertFile, HttpKeyFile)
					if err != nil {
						BeeLogger.Critical("ListenAndServeTLS: ", err)
						time.Sleep(100 * time.Microsecond)
						endRunning <- true
					}
				}()
			}

			if EnableHttpListen {
				go func() {
					app.Server.Addr = addr
					BeeLogger.Info("http server Running on %s", app.Server.Addr)
					if ListenTCP4 && HttpAddr == "" {
						ln, err := net.Listen("tcp4", app.Server.Addr)
						if err != nil {
							BeeLogger.Critical("ListenAndServe: ", err)
							time.Sleep(100 * time.Microsecond)
							endRunning <- true
							return
						}
						err = app.Server.Serve(ln)
						if err != nil {
							BeeLogger.Critical("ListenAndServe: ", err)
							time.Sleep(100 * time.Microsecond)
							endRunning <- true
							return
						}
					} else {
						err := app.Server.ListenAndServe()
						if err != nil {
							BeeLogger.Critical("ListenAndServe: ", err)
							time.Sleep(100 * time.Microsecond)
							endRunning <- true
						}
					}
				}()
			}
		}

	}
	<-endRunning
}

// Router adds a patterned controller handler to BeeApp.
// it's an alias method of App.Router.
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
func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *App {
	BeeApp.Handlers.Add(rootpath, c, mappingMethods...)
	return BeeApp
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
//}
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
func Include(cList ...ControllerInterface) *App {
	BeeApp.Handlers.Include(cList...)
	return BeeApp
}

// RESTRouter adds a restful controller handler to BeeApp.
// its' controller implements beego.ControllerInterface and
// defines a param "pattern/:objectId" to visit each resource.
func RESTRouter(rootpath string, c ControllerInterface) *App {
	Router(rootpath, c)
	Router(path.Join(rootpath, ":objectId"), c)
	return BeeApp
}

// AutoRouter adds defined controller handler to BeeApp.
// it's same to App.AutoRouter.
// if beego.AddAuto(&MainContorlller{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func AutoRouter(c ControllerInterface) *App {
	BeeApp.Handlers.AddAuto(c)
	return BeeApp
}

// AutoPrefix adds controller handler to BeeApp with prefix.
// it's same to App.AutoRouterWithPrefix.
// if beego.AutoPrefix("/admin",&MainContorlller{}) and MainController has methods List and Page,
// visit the url /admin/main/list to exec List function or /admin/main/page to exec Page function.
func AutoPrefix(prefix string, c ControllerInterface) *App {
	BeeApp.Handlers.AddAutoPrefix(prefix, c)
	return BeeApp
}

// Get used to register router for Get method
// usage:
//    beego.Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Get(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Get(rootpath, f)
	return BeeApp
}

// Post used to register router for Post method
// usage:
//    beego.Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Post(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Post(rootpath, f)
	return BeeApp
}

// Delete used to register router for Delete method
// usage:
//    beego.Delete("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Delete(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Delete(rootpath, f)
	return BeeApp
}

// Put used to register router for Put method
// usage:
//    beego.Put("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Put(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Put(rootpath, f)
	return BeeApp
}

// Head used to register router for Head method
// usage:
//    beego.Head("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Head(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Head(rootpath, f)
	return BeeApp
}

// Options used to register router for Options method
// usage:
//    beego.Options("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Options(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Options(rootpath, f)
	return BeeApp
}

// Patch used to register router for Patch method
// usage:
//    beego.Patch("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Patch(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Patch(rootpath, f)
	return BeeApp
}

// Any used to register router for all methods
// usage:
//    beego.Any("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Any(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Any(rootpath, f)
	return BeeApp
}

// Handler used to register a Handler router
// usage:
//    beego.Handler("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	BeeApp.Handlers.Handler(rootpath, h, options...)
	return BeeApp
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// beego.BeforeStatic, beego.BeforeRouter, beego.BeforeExec, beego.AfterExec and beego.FinishRouter.
// The bool params is for setting the returnOnOutput value (false allows multiple filters to execute)
func InsertFilter(pattern string, pos int, filter FilterFunc, params ...bool) *App {
	BeeApp.Handlers.InsertFilter(pattern, pos, filter, params...)
	return BeeApp
}
