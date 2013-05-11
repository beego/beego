package beego

import (
	"fmt"
	"github.com/astaxie/beego/session"
	"html/template"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
	"runtime"
)

const VERSION = "0.5.0"

var (
	BeeApp        *App
	AppName       string
	AppPath       string
	AppConfigPath string
	StaticDir     map[string]string
	TemplateCache map[string]*template.Template
	HttpAddr      string
	HttpPort      int
	RecoverPanic  bool
	AutoRender    bool
	PprofOn       bool
	ViewsPath     string
	RunMode       string //"dev" or "prod"
	AppConfig     *Config
	//related to session 
	GlobalSessions       *session.Manager //GlobalSessions
	SessionOn            bool             // wheather auto start session,default is false
	SessionProvider      string           // default session provider  memory mysql redis 
	SessionName          string           // sessionName cookie's name
	SessionGCMaxLifetime int64            // session's gc maxlifetime
	SessionSavePath      string           // session savepath if use mysql/redis/file this set to the connectinfo
	UseFcgi              bool
	MaxMemory            int64
	EnableGzip           bool // enable gzip
)

func init() {
	os.Chdir(path.Dir(os.Args[0]))
	BeeApp = NewApp()
	AppPath, _ = os.Getwd()
	StaticDir = make(map[string]string)
	TemplateCache = make(map[string]*template.Template)
	HttpAddr = ""
	HttpPort = 8080
	AppName = "beego"
	RunMode = "dev" //default runmod
	AutoRender = true
	RecoverPanic = true
	PprofOn = false
	ViewsPath = "views"
	SessionOn = false
	SessionProvider = "memory"
	SessionName = "beegosessionID"
	SessionGCMaxLifetime = 3600
	SessionSavePath = ""
	UseFcgi = false
	MaxMemory = 1 << 26 //64MB
	EnableGzip = false
	StaticDir["/static"] = "static"
	AppConfigPath = path.Join(AppPath, "conf", "app.conf")
	ParseConfig()
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
	var err error
	if UseFcgi {
		l, e := net.Listen("tcp", addr)
		if e != nil {
			BeeLogger.Fatal("Listen: ", e)
		}
		err = fcgi.Serve(l, app.Handlers)
	} else {
		err = http.ListenAndServe(addr, app.Handlers)
	}
	if err != nil {
		BeeLogger.Fatal("ListenAndServe: ", err)
	}
}

func (app *App) Router(path string, c ControllerInterface) *App {
	app.Handlers.Add(path, c)
	return app
}

func (app *App) Filter(filter http.HandlerFunc) *App {
	app.Handlers.Filter(filter)
	return app
}

func (app *App) FilterParam(param string, filter http.HandlerFunc) *App {
	app.Handlers.FilterParam(param, filter)
	return app
}

func (app *App) FilterPrefixPath(path string, filter http.HandlerFunc) *App {
	app.Handlers.FilterPrefixPath(path, filter)
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
	BeeLogger.Printf("[ACC] host: '%s', request: '%s %s', proto: '%s', ua: '%s', remote: '%s'\n", ctx.Request.Host, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.Proto, ctx.Request.UserAgent(), ctx.Request.RemoteAddr)
}

func RegisterController(path string, c ControllerInterface) *App {
	BeeApp.Router(path, c)
	return BeeApp
}

func Router(path string, c ControllerInterface) *App {
	BeeApp.Router(path, c)
	return BeeApp
}

func RouterHandler(path string, c http.Handler) *App {
	BeeApp.Handlers.AddHandler(path, c)
	return BeeApp
}

func Errorhandler(err string, h http.HandlerFunc) *App {
	ErrorMaps[err] = h
	return BeeApp
}

func SetViewsPath(path string) *App {
	BeeApp.SetViewsPath(path)
	return BeeApp
}

func SetStaticPath(url string, path string) *App {
	StaticDir[url] = path
	return BeeApp
}

func Filter(filter http.HandlerFunc) *App {
	BeeApp.Filter(filter)
	return BeeApp
}

func FilterParam(param string, filter http.HandlerFunc) *App {
	BeeApp.FilterParam(param, filter)
	return BeeApp
}

func FilterPrefixPath(path string, filter http.HandlerFunc) *App {
	BeeApp.FilterPrefixPath(path, filter)
	return BeeApp
}

func Run() {
	if AppConfigPath != path.Join(AppPath, "conf", "app.conf") {
		err := ParseConfig()
		if err != nil {
			if RunMode == "dev" {
				Warn(err)
			}
		}
	}
	if PprofOn {
		BeeApp.Router(`/debug/pprof`, &ProfController{})
		BeeApp.Router(`/debug/pprof/:pp([\w]+)`, &ProfController{})
	}
	if SessionOn {
		GlobalSessions, _ = session.NewManager(SessionProvider, SessionName, SessionGCMaxLifetime, SessionSavePath)
		go GlobalSessions.GC()
	}
	err := BuildTemplate(ViewsPath)
	if err != nil {
		if RunMode == "dev" {
			Warn(err)
		}
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	registerErrorHander()
	BeeApp.Run()
}
