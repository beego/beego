package beego

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"time"

	"github.com/astaxie/beego/context"
)

// FilterFunc defines filter function type.
type FilterFunc func(*context.Context)

// App defines beego application with a new PatternServeMux.
type App struct {
	Handlers *ControllerRegistor
}

// NewApp returns a new beego application.
func NewApp() *App {
	cr := NewControllerRegistor()
	app := &App{Handlers: cr}
	return app
}

// Run beego application.
func (app *App) Run() {
	addr := HttpAddr

	if HttpPort != 0 {
		addr = fmt.Sprintf("%s:%d", HttpAddr, HttpPort)
	}

	BeeLogger.Info("Running on %s", addr)

	var (
		err error
		l   net.Listener
	)

	if UseFcgi {
		if HttpPort == 0 {
			l, err = net.Listen("unix", addr)
		} else {
			l, err = net.Listen("tcp", addr)
		}
		if err != nil {
			BeeLogger.Critical("Listen: ", err)
		}
		err = fcgi.Serve(l, app.Handlers)
	} else {
		if EnableHotUpdate {
			server := &http.Server{
				Handler:      app.Handlers,
				ReadTimeout:  time.Duration(HttpServerTimeOut) * time.Second,
				WriteTimeout: time.Duration(HttpServerTimeOut) * time.Second,
			}
			laddr, err := net.ResolveTCPAddr("tcp", addr)
			if nil != err {
				BeeLogger.Critical("ResolveTCPAddr:", err)
			}
			l, err = GetInitListener(laddr)
			theStoppable = newStoppable(l)
			err = server.Serve(theStoppable)
			theStoppable.wg.Wait()
			CloseSelf()
		} else {
			s := &http.Server{
				Addr:         addr,
				Handler:      app.Handlers,
				ReadTimeout:  time.Duration(HttpServerTimeOut) * time.Second,
				WriteTimeout: time.Duration(HttpServerTimeOut) * time.Second,
			}
			if HttpTLS {
				err = s.ListenAndServeTLS(HttpCertFile, HttpKeyFile)
			} else {
				err = s.ListenAndServe()
			}
		}
	}

	if err != nil {
		BeeLogger.Critical("ListenAndServe: ", err)
		time.Sleep(100 * time.Microsecond)
	}
}

// Router adds a url-patterned controller handler.
// The path argument supports regex rules and specific placeholders.
// The c argument needs a controller handler implemented beego.ControllerInterface.
// The mapping methods argument only need one string to define custom router rules.
// usage:
//  simple router
//  beego.Router("/admin", &admin.UserController{})
//  beego.Router("/admin/index", &admin.ArticleController{})
//
//  regex router
//
//  beego.Router(“/api/:id([0-9]+)“, &controllers.RController{})
//
//  custom rules
//  beego.Router("/api/list",&RestController{},"*:ListFood")
//  beego.Router("/api/create",&RestController{},"post:CreateFood")
//  beego.Router("/api/update",&RestController{},"put:UpdateFood")
//  beego.Router("/api/delete",&RestController{},"delete:DeleteFood")
func (app *App) Router(path string, c ControllerInterface, mappingMethods ...string) *App {
	app.Handlers.Add(path, c, mappingMethods...)
	return app
}

// AutoRouter adds beego-defined controller handler.
// if beego.AddAuto(&MainContorlller{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func (app *App) AutoRouter(c ControllerInterface) *App {
	app.Handlers.AddAuto(c)
	return app
}

// UrlFor creates a url with another registered controller handler with params.
// The endpoint is formed as path.controller.name to defined the controller method which will run.
// The values need key-pair data to assign into controller method.
func (app *App) UrlFor(endpoint string, values ...string) string {
	return app.Handlers.UrlFor(endpoint, values...)
}

// [Deprecated] use InsertFilter.
// Filter adds a FilterFunc under pattern condition and named action.
// The actions contains BeforeRouter,AfterStatic,BeforeExec,AfterExec and FinishRouter.
func (app *App) Filter(pattern, action string, filter FilterFunc) *App {
	app.Handlers.AddFilter(pattern, action, filter)
	return app
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// beego.BeforeRouter, beego.AfterStatic, beego.BeforeExec, beego.AfterExec and beego.FinishRouter.
func (app *App) InsertFilter(pattern string, pos int, filter FilterFunc) *App {
	app.Handlers.InsertFilter(pattern, pos, filter)
	return app
}

// SetViewsPath sets view directory path in beego application.
// it returns beego application self.
func (app *App) SetViewsPath(path string) *App {
	ViewsPath = path
	return app
}

// SetStaticPath sets static directory path and proper url pattern in beego application.
// if beego.SetStaticPath("static","public"), visit /static/* to load static file in folder "public".
// it returns beego application self.
func (app *App) SetStaticPath(url string, path string) *App {
	StaticDir[url] = path
	return app
}

// DelStaticPath removes the static folder setting in this url pattern in beego application.
// it returns beego application self.
func (app *App) DelStaticPath(url string) *App {
	delete(StaticDir, url)
	return app
}
