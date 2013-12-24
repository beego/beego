package beego

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/middleware"
	"github.com/astaxie/beego/session"
)

// beego web framework version.
const VERSION = "1.0.1"

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

// Run beego application.
// it's alias of App.Run.
func Run() {
	// if AppConfigPath not In the conf/app.conf reParse config
	if AppConfigPath != filepath.Join(AppPath, "conf", "app.conf") {
		err := ParseConfig()
		if err != nil {
			// configuration is critical to app, panic here if parse failed
			panic(err)
		}
	}

	//init mime
	initMime()

	if SessionOn {
		GlobalSessions, _ = session.NewManager(SessionProvider,
			SessionName,
			SessionGCMaxLifetime,
			SessionSavePath,
			HttpTLS,
			SessionHashFunc,
			SessionHashKey,
			SessionCookieLifeTime)
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
	middleware.RegisterErrorHander()

	if EnableAdmin {
		go BeeAdminApp.Run()
	}

	BeeApp.Run()
}
