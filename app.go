package beego

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"time"

	"github.com/astaxie/beego/context"
)

type FilterFunc func(*context.Context)

type App struct {
	Handlers *ControllerRegistor
}

// New returns a new PatternServeMux.
func NewApp() *App {
	cr := NewControllerRegistor()
	app := &App{Handlers: cr}
	return app
}

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
			l, err = GetInitListner(laddr)
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

func (app *App) Router(path string, c ControllerInterface, mappingMethods ...string) *App {
	app.Handlers.Add(path, c, mappingMethods...)
	return app
}

func (app *App) AutoRouter(c ControllerInterface) *App {
	app.Handlers.AddAuto(c)
	return app
}

func (app *App) UrlFor(endpoint string, values ...string) string {
	return app.Handlers.UrlFor(endpoint, values...)
}
func (app *App) Filter(pattern, action string, filter FilterFunc) *App {
	app.Handlers.AddFilter(pattern, action, filter)
	return app
}

func (app *App) InsertFilter(pattern string, pos int, filter FilterFunc) *App {
	app.Handlers.InsertFilter(pattern, pos, filter)
	return app
}

func (app *App) SetViewsPath(path string) *App {
	ViewsPath = path
	return app
}

func (app *App) SetStaticPath(url string, path string) *App {
	StaticDir[url] = path
	return app
}

func (app *App) DelStaticPath(url string) *App {
	delete(StaticDir, url)
	return app
}
