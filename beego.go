package beego

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

var (
	BeeApp       *App
	AppName      string
	AppPath      string
	StaticDir    map[string]string
	HttpAddr     string
	HttpPort     int
	RecoverPanic bool
	AutoRender   bool
	ViewsPath    string
	RunMode      string //"dev" or "prod"
	AppConfig    *Config
)

func init() {
	BeeApp = NewApp()
	AppPath, _ = os.Getwd()
	StaticDir = make(map[string]string)
	var err error
	AppConfig, err = LoadConfig(path.Join(AppPath, "conf", "app.conf"))
	if err != nil {
		Trace("open Config err:", err)
		HttpAddr = ""
		HttpPort = 8080
		AppName = "beego"
		RunMode = "prod"
		AutoRender = true
		RecoverPanic = true
		ViewsPath = "views"
	} else {
		HttpAddr = AppConfig.String("httpaddr")
		if v, err := AppConfig.Int("httpport"); err != nil {
			HttpPort = 8080
		} else {
			HttpPort = v
		}
		AppName = AppConfig.String("appname")
		if runmode := AppConfig.String("runmode"); runmode != "" {
			RunMode = runmode
		} else {
			RunMode = "prod"
		}
		if ar, err := AppConfig.Bool("autorender"); err != nil {
			AutoRender = true
		} else {
			AutoRender = ar
		}
		if ar, err := AppConfig.Bool("autorecover"); err != nil {
			RecoverPanic = true
		} else {
			RecoverPanic = ar
		}
		if views := AppConfig.String("viewspath"); views == "" {
			ViewsPath = "views"
		} else {
			ViewsPath = views
		}
	}
	StaticDir["/static"] = "static"

}

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
	addr := fmt.Sprintf("%s:%d", HttpAddr, HttpPort)
	err := http.ListenAndServe(addr, app.Handlers)
	if err != nil {
		BeeLogger.Fatal("ListenAndServe: ", err)
	}
}

func (app *App) RegisterController(path string, c ControllerInterface) *App {
	app.Handlers.Add(path, c)
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

func (app *App) ErrorLog(ctx *Context) {
	BeeLogger.Printf("[ERR] host: '%s', request: '%s %s', proto: '%s', ua: '%s', remote: '%s'\n", ctx.Request.Host, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.Proto, ctx.Request.UserAgent(), ctx.Request.RemoteAddr)
}

func (app *App) AccessLog(ctx *Context) {
	BeeLogger.Printf("[ACC] host: '%s', request: '%s %s', proto: '%s', ua: %s'', remote: '%s'\n", ctx.Request.Host, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.Proto, ctx.Request.UserAgent(), ctx.Request.RemoteAddr)
}
