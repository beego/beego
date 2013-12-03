package beego

import (
	"net/http"
	"path"
	"strings"

	"github.com/astaxie/beego/middleware"
	"github.com/astaxie/beego/session"
)

const VERSION = "0.9.9"

func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *App {
	BeeApp.Router(rootpath, c, mappingMethods...)
	return BeeApp
}

func RESTRouter(rootpath string, c ControllerInterface) *App {
	Router(rootpath, c)
	Router(path.Join(rootpath, ":objectId"), c)
	return BeeApp
}

func AutoRouter(c ControllerInterface) *App {
	BeeApp.AutoRouter(c)
	return BeeApp
}

func Errorhandler(err string, h http.HandlerFunc) *App {
	middleware.Errorhandler(err, h)
	return BeeApp
}

func SetViewsPath(path string) *App {
	BeeApp.SetViewsPath(path)
	return BeeApp
}

func SetStaticPath(url string, path string) *App {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	StaticDir[url] = path
	return BeeApp
}

func DelStaticPath(url string) *App {
	delete(StaticDir, url)
	return BeeApp
}

//!!DEPRECATED!! use InsertFilter
//action has four values:
//BeforRouter
//AfterStatic
//BeforExec
//AfterExec
func AddFilter(pattern, action string, filter FilterFunc) *App {
	BeeApp.Filter(pattern, action, filter)
	return BeeApp
}

func InsertFilter(pattern string, pos int, filter FilterFunc) *App {
	BeeApp.InsertFilter(pattern, pos, filter)
	return BeeApp
}

func Run() {
	InitConfig()

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
