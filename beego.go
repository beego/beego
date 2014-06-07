// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package beego

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego/middleware"
	"github.com/astaxie/beego/session"
)

// beego web framework version.
const VERSION = "1.2.0"

type hookfunc func() error //hook function to run
var hooks []hookfunc       //hook function slice to store the hookfunc

type groupRouter struct {
	pattern        string
	controller     ControllerInterface
	mappingMethods string
}

// RouterGroups which will store routers
type GroupRouters []groupRouter

// Get a new GroupRouters
func NewGroupRouters() GroupRouters {
	return make(GroupRouters, 0)
}

// Add Router in the GroupRouters
// it is for plugin or module to register router
func (gr *GroupRouters) AddRouter(pattern string, c ControllerInterface, mappingMethod ...string) {
	var newRG groupRouter
	if len(mappingMethod) > 0 {
		newRG = groupRouter{
			pattern,
			c,
			mappingMethod[0],
		}
	} else {
		newRG = groupRouter{
			pattern,
			c,
			"",
		}
	}
	*gr = append(*gr, newRG)
}

func (gr *GroupRouters) AddAuto(c ControllerInterface) {
	newRG := groupRouter{
		"",
		c,
		"",
	}
	*gr = append(*gr, newRG)
}

// AddGroupRouter with the prefix
// it will register the router in BeeApp
// the follow code is write in modules:
// GR:=NewGroupRouters()
// GR.AddRouter("/login",&UserController,"get:Login")
// GR.AddRouter("/logout",&UserController,"get:Logout")
// GR.AddRouter("/register",&UserController,"get:Reg")
// the follow code is write in app:
// import "github.com/beego/modules/auth"
// AddRouterGroup("/admin", auth.GR)
func AddGroupRouter(prefix string, groups GroupRouters) *App {
	for _, v := range groups {
		if v.pattern == "" {
			BeeApp.AutoRouterWithPrefix(prefix, v.controller)
		} else if v.mappingMethods != "" {
			BeeApp.Router(prefix+v.pattern, v.controller, v.mappingMethods)
		} else {
			BeeApp.Router(prefix+v.pattern, v.controller)
		}

	}
	return BeeApp
}

// Router adds a patterned controller handler to BeeApp.
// it's an alias method of App.Router.
func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *App {
	BeeApp.Router(rootpath, c, mappingMethods...)
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
func AutoRouter(c ControllerInterface) *App {
	BeeApp.AutoRouter(c)
	return BeeApp
}

// AutoPrefix adds controller handler to BeeApp with prefix.
// it's same to App.AutoRouterWithPrefix.
func AutoPrefix(prefix string, c ControllerInterface) *App {
	BeeApp.AutoRouterWithPrefix(prefix, c)
	return BeeApp
}

// register router for Get method
func Get(rootpath string, f FilterFunc) *App {
	BeeApp.Get(rootpath, f)
	return BeeApp
}

// register router for Post method
func Post(rootpath string, f FilterFunc) *App {
	BeeApp.Post(rootpath, f)
	return BeeApp
}

// register router for Delete method
func Delete(rootpath string, f FilterFunc) *App {
	BeeApp.Delete(rootpath, f)
	return BeeApp
}

// register router for Put method
func Put(rootpath string, f FilterFunc) *App {
	BeeApp.Put(rootpath, f)
	return BeeApp
}

// register router for Head method
func Head(rootpath string, f FilterFunc) *App {
	BeeApp.Head(rootpath, f)
	return BeeApp
}

// register router for Options method
func Options(rootpath string, f FilterFunc) *App {
	BeeApp.Options(rootpath, f)
	return BeeApp
}

// register router for Patch method
func Patch(rootpath string, f FilterFunc) *App {
	BeeApp.Patch(rootpath, f)
	return BeeApp
}

// register router for all method
func Any(rootpath string, f FilterFunc) *App {
	BeeApp.Any(rootpath, f)
	return BeeApp
}

// register router for own Handler
func Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	BeeApp.Handler(rootpath, h, options...)
	return BeeApp
}

// ErrorHandler registers http.HandlerFunc to each http err code string.
// usage:
// 	beego.ErrorHandler("404",NotFound)
//	beego.ErrorHandler("500",InternalServerError)
func Errorhandler(err string, h http.HandlerFunc) *App {
	middleware.Errorhandler(err, h)
	return BeeApp
}

// SetViewsPath sets view directory to BeeApp.
// it's alias of App.SetViewsPath.
func SetViewsPath(path string) *App {
	BeeApp.SetViewsPath(path)
	return BeeApp
}

// SetStaticPath sets static directory and url prefix to BeeApp.
// it's alias of App.SetStaticPath.
func SetStaticPath(url string, path string) *App {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	url = strings.TrimRight(url, "/")
	StaticDir[url] = path
	return BeeApp
}

// DelStaticPath removes the static folder setting in this url pattern in beego application.
// it's alias of App.DelStaticPath.
func DelStaticPath(url string) *App {
	delete(StaticDir, url)
	return BeeApp
}

// [Deprecated] use InsertFilter.
// Filter adds a FilterFunc under pattern condition and named action.
// The actions contains BeforeRouter,AfterStatic,BeforeExec,AfterExec and FinishRouter.
// it's alias of App.Filter.
func AddFilter(pattern, action string, filter FilterFunc) *App {
	BeeApp.Filter(pattern, action, filter)
	return BeeApp
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// beego.BeforeRouter, beego.AfterStatic, beego.BeforeExec, beego.AfterExec and beego.FinishRouter.
// it's alias of App.InsertFilter.
func InsertFilter(pattern string, pos int, filter FilterFunc) *App {
	BeeApp.InsertFilter(pattern, pos, filter)
	return BeeApp
}

// The hookfunc will run in beego.Run()
// such as sessionInit, middlerware start, buildtemplate, admin start
func AddAPPStartHook(hf hookfunc) {
	hooks = append(hooks, hf)
}

// Run beego application.
// it's alias of App.Run.
func Run() {
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
				`"providerConfig":"` + SessionSavePath + `",` +
				`"secure":` + strconv.FormatBool(HttpTLS) + `,` +
				`"sessionIDHashFunc":"` + SessionHashFunc + `",` +
				`"sessionIDHashKey":"` + SessionHashKey + `",` +
				`"enableSetCookie":` + strconv.FormatBool(SessionAutoSetCookie) + `,` +
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

	middleware.VERSION = VERSION
	middleware.AppName = AppName
	middleware.RegisterErrorHandler()
}

func TestBeegoInit(apppath string) {
	AppPath = apppath
	AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")
	err := ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		Info(err)
	}
	os.Chdir(AppPath)
	initBeforeHttpRun()
}

//添加远程配置文件
func Config(config string) {
	if len(config) != 0 {

		resp, err := http.Get(config)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("Download remote config error:", err)
		}
		if resp.StatusCode == 200 {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Read remote config error:", err)
			} else {
				f, _ := ioutil.TempFile("", "beego")
				ioutil.WriteFile(f.Name(), data, os.ModeAppend)
				AppConfigPath = f.Name()
				AppConfigRemote = config
			}

		}

	}
}
func init() {
	hooks = make([]hookfunc, 0)
	//init mime
	AddAPPStartHook(initMime)
}
