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

package beego

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego/session"
)

// beego web framework version.
const VERSION = "1.5.0"

//hook function to run
type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) //hook function slice to store the hookfunc
)

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
	initBeforeHttpRun()

	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			HttpAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			HttpPort, _ = strconv.Atoi(strs[1])
		}
	}

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

	if AutoRender {
		err := BuildTemplate(ViewsPath)
		if err != nil && RunMode == "dev" {
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
