package beego

import (
	"github.com/astaxie/beego/session"
	"net/http"
	"path"
)

const VERSION = "0.9.0"

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

func DelStaticPath(url string) *App {
	delete(StaticDir, url)
	return BeeApp
}

//action has four values:
//BeforRouter
//AfterStatic
//BeforExec
//AfterExec
func AddFilter(pattern, action string, filter FilterFunc) *App {
	BeeApp.Filter(pattern, action, filter)
	return BeeApp
}

func Run() {
	//if AppConfigPath not In the conf/app.conf reParse config
	if AppConfigPath != path.Join(AppPath, "conf", "app.conf") {
		err := ParseConfig()
		if err != nil {
			if RunMode == "dev" {
				Warn(err)
			}
		}
	}

	if SessionOn {
		GlobalSessions, _ = session.NewManager(SessionProvider, SessionName, SessionGCMaxLifetime, SessionSavePath)
		go GlobalSessions.GC()
	}

	if AutoRender {
		err := BuildTemplate(ViewsPath)
		if err != nil {
			if RunMode == "dev" {
				Warn(err)
			}
		}
	}
	registerErrorHander()
	BeeApp.Run()
}
