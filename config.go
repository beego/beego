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
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
	"github.com/astaxie/beego/utils"
)

var (
	BeeApp                 *App // beego application
	AppName                string
	AppPath                string
	workPath               string
	AppConfigPath          string
	StaticDir              map[string]string
	TemplateCache          map[string]*template.Template // template caching map
	StaticExtensionsToGzip []string                      // files with should be compressed with gzip (.js,.css,etc)
	EnableHttpListen       bool
	HttpAddr               string
	HttpPort               int
	EnableHttpTLS          bool
	HttpsPort              int
	HttpCertFile           string
	HttpKeyFile            string
	RecoverPanic           bool // flag of auto recover panic
	AutoRender             bool // flag of render template automatically
	ViewsPath              string
	RunMode                string // run mode, "dev" or "prod"
	AppConfig              config.ConfigContainer
	GlobalSessions         *session.Manager // global session mananger
	SessionOn              bool             // flag of starting session auto. default is false.
	SessionProvider        string           // default session provider, memory, mysql , redis ,etc.
	SessionName            string           // the cookie name when saving session id into cookie.
	SessionGCMaxLifetime   int64            // session gc time for auto cleaning expired session.
	SessionSavePath        string           // if use mysql/redis/file provider, define save path to connection info.
	SessionHashFunc        string           // session hash generation func.
	SessionHashKey         string           // session hash salt string.
	SessionCookieLifeTime  int              // the life time of session id in cookie.
	SessionAutoSetCookie   bool             // auto setcookie
	SessionDomain          string           // the cookie domain default is empty
	UseFcgi                bool
	MaxMemory              int64
	EnableGzip             bool // flag of enable gzip
	DirectoryIndex         bool // flag of display directory index. default is false.
	HttpServerTimeOut      int64
	ErrorsShow             bool   // flag of show errors in page. if true, show error and trace info in page rendered with error template.
	XSRFKEY                string // xsrf hash salt string.
	EnableXSRF             bool   // flag of enable xsrf.
	XSRFExpire             int    // the expiry of xsrf value.
	CopyRequestBody        bool   // flag of copy raw request body in context.
	TemplateLeft           string
	TemplateRight          string
	BeegoServerName        string // beego server name exported in response header.
	EnableAdmin            bool   // flag of enable admin module to log every request info.
	AdminHttpAddr          string // http server configurations for admin module.
	AdminHttpPort          int
	FlashName              string // name of the flash variable found in response header and cookie
	FlashSeperator         string // used to seperate flash key:value
	AppConfigProvider      string // config provider
	EnableDocs             bool   // enable generate docs & server docs API Swagger
)

func init() {
	// create beego application
	BeeApp = NewApp()

	workPath, _ = os.Getwd()
	workPath, _ = filepath.Abs(workPath)
	// initialize default configurations
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")

	if workPath != AppPath {
		if utils.FileExists(AppConfigPath) {
			os.Chdir(AppPath)
		} else {
			AppConfigPath = filepath.Join(workPath, "conf", "app.conf")
		}
	}

	AppConfigProvider = "ini"

	StaticDir = make(map[string]string)
	StaticDir["/static"] = "static"

	StaticExtensionsToGzip = []string{".css", ".js"}

	TemplateCache = make(map[string]*template.Template)

	// set this to 0.0.0.0 to make this app available to externally
	EnableHttpListen = true //default enable http Listen

	HttpAddr = ""
	HttpPort = 8080

	HttpsPort = 10443

	AppName = "beego"

	RunMode = "dev" //default runmod

	AutoRender = true

	RecoverPanic = true

	ViewsPath = "views"

	SessionOn = false
	SessionProvider = "memory"
	SessionName = "beegosessionID"
	SessionGCMaxLifetime = 3600
	SessionSavePath = ""
	SessionHashFunc = "sha1"
	SessionHashKey = "beegoserversessionkey"
	SessionCookieLifeTime = 0 //set cookie default is the brower life
	SessionAutoSetCookie = true

	UseFcgi = false

	MaxMemory = 1 << 26 //64MB

	EnableGzip = false

	HttpServerTimeOut = 0

	ErrorsShow = true

	XSRFKEY = "beegoxsrf"
	XSRFExpire = 0

	TemplateLeft = "{{"
	TemplateRight = "}}"

	BeegoServerName = "beegoServer:" + VERSION

	EnableAdmin = false
	AdminHttpAddr = "127.0.0.1"
	AdminHttpPort = 8088

	FlashName = "BEEGO_FLASH"
	FlashSeperator = "BEEGOFLASH"

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init BeeLogger
	BeeLogger = logs.NewLogger(10000)
	err := BeeLogger.SetLogger("console", "")
	if err != nil {
		fmt.Println("init console log error:", err)
	}

	err = ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		Info(err)
	}
}

// ParseConfig parsed default config file.
// now only support ini, next will support json.
func ParseConfig() (err error) {
	AppConfig, err = config.NewConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		AppConfig = config.NewFakeConfig()
		return err
	} else {

		if v, err := GetConfig("string", "HttpAddr"); err == nil {
			HttpAddr = v.(string)
		}

		if v, err := GetConfig("int", "HttpPort"); err == nil {
			HttpPort = v.(int)
		}

		if v, err := GetConfig("bool", "EnableHttpListen"); err == nil {
			EnableHttpListen = v.(bool)
		}

		if maxmemory, err := GetConfig("int64", "MaxMemory"); err == nil {
			MaxMemory = maxmemory.(int64)
		}

		if appname, _ := GetConfig("string", "AppName"); appname != "" {
			AppName = appname.(string)
		}

		if runmode, _ := GetConfig("string", "RunMode"); runmode != "" {
			RunMode = runmode.(string)
		}

		if autorender, err := GetConfig("bool", "AutoRender"); err == nil {
			AutoRender = autorender.(bool)
		}

		if autorecover, err := GetConfig("bool", "RecoverPanic"); err == nil {
			RecoverPanic = autorecover.(bool)
		}

		if views, _ := GetConfig("string", "ViewsPath"); views != "" {
			ViewsPath = views.(string)
		}

		if sessionon, err := GetConfig("bool", "SessionOn"); err == nil {
			SessionOn = sessionon.(bool)
		}

		if sessProvider, _ := GetConfig("string", "SessionProvider"); sessProvider != "" {
			SessionProvider = sessProvider.(string)
		}

		if sessName, _ := GetConfig("string", "SessionName"); sessName != "" {
			SessionName = sessName.(string)
		}

		if sesssavepath, _ := GetConfig("string", "SessionSavePath"); sesssavepath != "" {
			SessionSavePath = sesssavepath.(string)
		}

		if sesshashfunc, _ := GetConfig("string", "SessionHashFunc"); sesshashfunc != "" {
			SessionHashFunc = sesshashfunc.(string)
		}

		if sesshashkey, _ := GetConfig("string", "SessionHashKey"); sesshashkey != "" {
			SessionHashKey = sesshashkey.(string)
		}

		if sessMaxLifeTime, err := GetConfig("int64", "SessionGCMaxLifetime"); err == nil && sessMaxLifeTime != 0 {
			SessionGCMaxLifetime = sessMaxLifeTime.(int64)
		}

		if sesscookielifetime, err := GetConfig("int", "SessionCookieLifeTime"); err == nil && sesscookielifetime != 0 {
			SessionCookieLifeTime = sesscookielifetime.(int)
		}

		if usefcgi, err := GetConfig("bool", "UseFcgi"); err == nil {
			UseFcgi = usefcgi.(bool)
		}

		if enablegzip, err := GetConfig("bool", "EnableGzip"); err == nil {
			EnableGzip = enablegzip.(bool)
		}

		if directoryindex, err := GetConfig("bool", "DirectoryIndex"); err == nil {
			DirectoryIndex = directoryindex.(bool)
		}

		if timeout, err := GetConfig("int64", "HttpServerTimeOut"); err == nil {
			HttpServerTimeOut = timeout.(int64)
		}

		if errorsshow, err := GetConfig("bool", "ErrorsShow"); err == nil {
			ErrorsShow = errorsshow.(bool)
		}

		if copyrequestbody, err := GetConfig("bool", "CopyRequestBody"); err == nil {
			CopyRequestBody = copyrequestbody.(bool)
		}

		if xsrfkey, _ := GetConfig("string", "XSRFKEY"); xsrfkey != "" {
			XSRFKEY = xsrfkey.(string)
		}

		if enablexsrf, err := GetConfig("bool", "EnableXSRF"); err == nil {
			EnableXSRF = enablexsrf.(bool)
		}

		if expire, err := GetConfig("int", "XSRFExpire"); err == nil {
			XSRFExpire = expire.(int)
		}

		if tplleft, _ := GetConfig("string", "TemplateLeft"); tplleft != "" {
			TemplateLeft = tplleft.(string)
		}

		if tplright, _ := GetConfig("string", "TemplateRight"); tplright != "" {
			TemplateRight = tplright.(string)
		}

		if httptls, err := GetConfig("bool", "EnableHttpTLS"); err == nil {
			EnableHttpTLS = httptls.(bool)
		}

		if httpsport, err := GetConfig("int", "HttpsPort"); err == nil {
			HttpsPort = httpsport.(int)
		}

		if certfile, _ := GetConfig("string", "HttpCertFile"); certfile != "" {
			HttpCertFile = certfile.(string)
		}

		if keyfile, _ := GetConfig("string", "HttpKeyFile"); keyfile != "" {
			HttpKeyFile = keyfile.(string)
		}

		if serverName, _ := GetConfig("string", "BeegoServerName"); serverName != "" {
			BeegoServerName = serverName.(string)
		}

		if flashname, _ := GetConfig("string", "FlashName"); flashname != "" {
			FlashName = flashname.(string)
		}

		if flashseperator, _ := GetConfig("string", "FlashSeperator"); flashseperator != "" {
			FlashSeperator = flashseperator.(string)
		}

		if sd, _ := GetConfig("string", "StaticDir"); sd != "" {
			for k := range StaticDir {
				delete(StaticDir, k)
			}
			sds := strings.Fields(sd.(string))
			for _, v := range sds {
				if url2fsmap := strings.SplitN(v, ":", 2); len(url2fsmap) == 2 {
					StaticDir["/"+strings.TrimRight(url2fsmap[0], "/")] = url2fsmap[1]
				} else {
					StaticDir["/"+strings.TrimRight(url2fsmap[0], "/")] = url2fsmap[0]
				}
			}
		}

		if sgz, _ := GetConfig("string", "StaticExtensionsToGzip"); sgz != "" {
			extensions := strings.Split(sgz.(string), ",")
			if len(extensions) > 0 {
				StaticExtensionsToGzip = []string{}
				for _, ext := range extensions {
					if len(ext) == 0 {
						continue
					}
					extWithDot := ext
					if extWithDot[:1] != "." {
						extWithDot = "." + extWithDot
					}
					StaticExtensionsToGzip = append(StaticExtensionsToGzip, extWithDot)
				}
			}
		}

		if enableadmin, err := GetConfig("bool", "EnableAdmin"); err == nil {
			EnableAdmin = enableadmin.(bool)
		}

		if adminhttpaddr, _ := GetConfig("string", "AdminHttpAddr"); adminhttpaddr != "" {
			AdminHttpAddr = adminhttpaddr.(string)
		}

		if adminhttpport, err := GetConfig("int", "AdminHttpPort"); err == nil {
			AdminHttpPort = adminhttpport.(int)
		}

		if enabledocs, err := GetConfig("bool", "EnableDocs"); err == nil {
			EnableDocs = enabledocs.(bool)
		}
	}
	return nil
}

// Getconfig throw the Runmode
// [dev]
// name = astaixe
// IsEnable = false
// [prod]
// name = slene
// IsEnable = true
//
// usage:
// GetConfig("string", "name")
// GetConfig("bool", "IsEnable")
func GetConfig(typ, key string) (interface{}, error) {
	switch typ {
	case "string":
		v := AppConfig.String(RunMode + "::" + key)
		if v == "" {
			v = AppConfig.String(key)
		}
		return v, nil
	case "strings":
		v := AppConfig.Strings(RunMode + "::" + key)
		if len(v) == 0 {
			v = AppConfig.Strings(key)
		}
		return v, nil
	case "int":
		v, err := AppConfig.Int(RunMode + "::" + key)
		if err != nil || v == 0 {
			return AppConfig.Int(key)
		}
		return v, nil
	case "bool":
		v, err := AppConfig.Bool(RunMode + "::" + key)
		if err != nil {
			return AppConfig.Bool(key)
		}
		return v, nil
	case "int64":
		v, err := AppConfig.Int64(RunMode + "::" + key)
		if err != nil || v == 0 {
			return AppConfig.Int64(key)
		}
		return v, nil
	case "float":
		v, err := AppConfig.Float(RunMode + "::" + key)
		if err != nil || v == 0 {
			return AppConfig.Float(key)
		}
		return v, nil
	}
	return "", errors.New("not support type")
}
