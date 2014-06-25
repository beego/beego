// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

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

	BeeLogger.Info("Running on %s", addr)

	var (
		err error
		l   net.Listener
	)
	endRunning := make(chan bool, 1)

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
		app.Server.Addr = addr
		app.Server.Handler = app.Handlers
		app.Server.ReadTimeout = time.Duration(HttpServerTimeOut) * time.Second
		app.Server.WriteTimeout = time.Duration(HttpServerTimeOut) * time.Second

		if EnableHttpTLS {
			go func() {
				if HttpsPort != 0 {
					app.Server.Addr = fmt.Sprintf("%s:%d", HttpAddr, HttpsPort)
				}
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
				err := app.Server.ListenAndServe()
				if err != nil {
					BeeLogger.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
				}
			}()
		}
	}

	<-endRunning
}
