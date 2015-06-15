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

// beego is an open-source, high-performance, modularity, full-stack web framework
//
// package main
//
// import "github.com/astaxie/beego"
//
// func main() {
//  	beego.Run()
// }
//
// more infomation: http://beego.me
package beego

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego/session"
)

// beego web framework version.
const VERSION = "1.5.0"

type hookfunc func() error //hook function to run
var hooks []hookfunc       //hook function slice to store the hookfunc

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

// Router add list from
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

// register router for Get method
// usage:
//    beego.Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Get(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Get(rootpath, f)
	return BeeApp
}

// register router for Post method
// usage:
//    beego.Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Post(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Post(rootpath, f)
	return BeeApp
}

// register router for Delete method
// usage:
//    beego.Delete("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Delete(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Delete(rootpath, f)
	return BeeApp
}

// register router for Put method
// usage:
//    beego.Put("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Put(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Put(rootpath, f)
	return BeeApp
}

// register router for Head method
// usage:
//    beego.Head("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Head(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Head(rootpath, f)
	return BeeApp
}

// register router for Options method
// usage:
//    beego.Options("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Options(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Options(rootpath, f)
	return BeeApp
}

// register router for Patch method
// usage:
//    beego.Patch("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Patch(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Patch(rootpath, f)
	return BeeApp
}

// register router for all method
// usage:
//    beego.Any("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Any(rootpath string, f FilterFunc) *App {
	BeeApp.Handlers.Any(rootpath, f)
	return BeeApp
}

// register router for own Handler
// usage:
//    beego.Handler("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	BeeApp.Handlers.Handler(rootpath, h, options...)
	return BeeApp
}

// SetViewsPath sets view directory path in beego application.
func SetViewsPath(path string) *App {
	ViewsPath = path
	return BeeApp
}

// SetStaticPath sets static directory path and proper url pattern in beego application.
// if beego.SetStaticPath("static","public"), visit /static/* to load static file in folder "public".
func SetStaticPath(url string, path string) *App {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	url = strings.TrimRight(url, "/")
	StaticDir[url] = path
	return BeeApp
}

// DelStaticPath removes the static folder setting in this url pattern in beego application.
func DelStaticPath(url string) *App {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	url = strings.TrimRight(url, "/")
	delete(StaticDir, url)
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

// The hookfunc will run in beego.Run()
// such as sessionInit, middlerware start, buildtemplate, admin start
func AddAPPStartHook(hf hookfunc) {
	hooks = append(hooks, hf)
}

// Run beego application.
// beego.Run() default run on HttpPort
// beego.Run(":8089")
// beego.Run("127.0.0.1:8089")
func Run(params ...string) {
	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			HttpAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			HttpPort, _ = strconv.Atoi(strs[1])
		}
	}
	initBeforeHttpRun()

	if EnableAdmin {
		go beeAdminApp.Run()
	}

	BeeApp.Run()
}

func initBeforeHttpRun() {
	// if AppConfigPath not In the conf/app.conf reParse config
	if AppConfigPath != filepath.Join(AppPath, "conf", "app.conf") {
		err := ParseConfig()
		if err != nil && AppConfigPath != filepath.Join(workPath, "conf", "app.conf") {
			// configuration is critical to app, panic here if parse failed
			panic(err)
		}
	}

	//init mime
	AddAPPStartHook(initMime)

	// do hooks function
	for _, hk := range hooks {
		err := hk()
		if err != nil {
			panic(err)
		}
	}

	if SessionOn {
		var err error
		sessionConfig := AppConfig.String("sessionConfig")
		if sessionConfig == "" {
			sessionConfig = `{"cookieName":"` + SessionName + `",` +
				`"gclifetime":` + strconv.FormatInt(SessionGCMaxLifetime, 10) + `,` +
				`"providerConfig":"` + filepath.ToSlash(SessionSavePath) + `",` +
				`"secure":` + strconv.FormatBool(EnableHttpTLS) + `,` +
				`"enableSetCookie":` + strconv.FormatBool(SessionAutoSetCookie) + `,` +
				`"domain":"` + SessionDomain + `",` +
				`"cookieLifeTime":` + strconv.Itoa(SessionCookieLifeTime) + `}`
		}
		GlobalSessions, err = session.NewManager(SessionProvider,
			sessionConfig)
		if err != nil {
			panic(err)
		}
		go GlobalSessions.GC()
	}

	err := BuildTemplate(ViewsPath)
	if err != nil {
		if RunMode == "dev" {
			Warn(err)
		}
	}

	registerDefaultErrorHandler()

	if EnableDocs {
		Get("/docs", serverDocs)
		Get("/docs/*", serverDocs)
	}
}

// this function is for test package init
func TestBeegoInit(apppath string) {
	AppPath = apppath
	os.Setenv("BEEGO_RUNMODE", "test")
	AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")
	err := ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		Info(err)
	}
	os.Chdir(AppPath)
	initBeforeHttpRun()
}

func init() {
	hooks = make([]hookfunc, 0)
}
