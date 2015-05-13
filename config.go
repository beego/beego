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
	ListenTCP4             bool
	EnableHttpTLS          bool
	HttpsPort              int
	HttpCertFile           string
	HttpKeyFile            string
	RecoverPanic           bool // flag of auto recover panic
	AutoRender             bool // flag of render template automatically
	ViewsPath              string
	AppConfig              *beegoAppConfig
	RunMode                string           // run mode, "dev" or "prod"
	GlobalSessions         *session.Manager // global session mananger
	SessionOn              bool             // flag of starting session auto. default is false.
	SessionProvider        string           // default session provider, memory, mysql , redis ,etc.
	SessionName            string           // the cookie name when saving session id into cookie.
	SessionGCMaxLifetime   int64            // session gc time for auto cleaning expired session.
	SessionSavePath        string           // if use mysql/redis/file provider, define save path to connection info.
	SessionCookieLifeTime  int              // the life time of session id in cookie.
	SessionAutoSetCookie   bool             // auto setcookie
	SessionDomain          string           // the cookie domain default is empty
	UseFcgi                bool
	UseStdIo               bool
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
	RouterCaseSensitive    bool   // router case sensitive default is true
	AccessLogs             bool   // print access logs, default is false
	Graceful               bool   // use graceful start the server
)

type beegoAppConfig struct {
	innerConfig config.ConfigContainer
}

func newAppConfig(AppConfigProvider, AppConfigPath string) (*beegoAppConfig, error) {
	ac, err := config.NewConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		return nil, err
	}
	rac := &beegoAppConfig{ac}
	return rac, nil
}

func (b *beegoAppConfig) Set(key, val string) error {
	err := b.innerConfig.Set(RunMode+"::"+key, val)
	if err == nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *beegoAppConfig) String(key string) string {
	v := b.innerConfig.String(RunMode + "::" + key)
	if v == "" {
		return b.innerConfig.String(key)
	}
	return v
}

func (b *beegoAppConfig) Strings(key string) []string {
	v := b.innerConfig.Strings(RunMode + "::" + key)
	if v[0] == "" {
		return b.innerConfig.Strings(key)
	}
	return v
}

func (b *beegoAppConfig) Int(key string) (int, error) {
	v, err := b.innerConfig.Int(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int(key)
	}
	return v, nil
}

func (b *beegoAppConfig) Int64(key string) (int64, error) {
	v, err := b.innerConfig.Int64(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int64(key)
	}
	return v, nil
}

func (b *beegoAppConfig) Bool(key string) (bool, error) {
	v, err := b.innerConfig.Bool(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Bool(key)
	}
	return v, nil
}

func (b *beegoAppConfig) Float(key string) (float64, error) {
	v, err := b.innerConfig.Float(RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Float(key)
	}
	return v, nil
}

func (b *beegoAppConfig) DefaultString(key string, defaultval string) string {
	v := b.String(key)
	if v != "" {
		return v
	}
	return defaultval
}

func (b *beegoAppConfig) DefaultStrings(key string, defaultval []string) []string {
	v := b.Strings(key)
	if len(v) != 0 {
		return v
	}
	return defaultval
}

func (b *beegoAppConfig) DefaultInt(key string, defaultval int) int {
	v, err := b.Int(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *beegoAppConfig) DefaultInt64(key string, defaultval int64) int64 {
	v, err := b.Int64(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *beegoAppConfig) DefaultBool(key string, defaultval bool) bool {
	v, err := b.Bool(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *beegoAppConfig) DefaultFloat(key string, defaultval float64) float64 {
	v, err := b.Float(key)
	if err == nil {
		return v
	}
	return defaultval
}

func (b *beegoAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *beegoAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *beegoAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}

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
	SessionCookieLifeTime = 0 //set cookie default is the brower life
	SessionAutoSetCookie = true

	UseFcgi = false
	UseStdIo = false

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

	RouterCaseSensitive = true

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init BeeLogger
	BeeLogger = logs.NewLogger(10000)
	err := BeeLogger.SetLogger("console", "")
	if err != nil {
		fmt.Println("init console log error:", err)
	}
	SetLogFuncCall(true)

	err = ParseConfig()
	if err != nil && os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		ac := config.NewFakeConfig()
		AppConfig = &beegoAppConfig{ac}
		Warning(err)
	}
}

// ParseConfig parsed default config file.
// now only support ini, next will support json.
func ParseConfig() (err error) {
	AppConfig, err = newAppConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		return err
	}
	envRunMode := os.Getenv("BEEGO_RUNMODE")
	// set the runmode first
	if envRunMode != "" {
		RunMode = envRunMode
	} else if runmode := AppConfig.String("RunMode"); runmode != "" {
		RunMode = runmode
	}

	HttpAddr = AppConfig.String("HttpAddr")

	if v, err := AppConfig.Int("HttpPort"); err == nil {
		HttpPort = v
	}

	if v, err := AppConfig.Bool("ListenTCP4"); err == nil {
		ListenTCP4 = v
	}

	if v, err := AppConfig.Bool("EnableHttpListen"); err == nil {
		EnableHttpListen = v
	}

	if maxmemory, err := AppConfig.Int64("MaxMemory"); err == nil {
		MaxMemory = maxmemory
	}

	if appname := AppConfig.String("AppName"); appname != "" {
		AppName = appname
	}

	if autorender, err := AppConfig.Bool("AutoRender"); err == nil {
		AutoRender = autorender
	}

	if autorecover, err := AppConfig.Bool("RecoverPanic"); err == nil {
		RecoverPanic = autorecover
	}

	if views := AppConfig.String("ViewsPath"); views != "" {
		ViewsPath = views
	}

	if sessionon, err := AppConfig.Bool("SessionOn"); err == nil {
		SessionOn = sessionon
	}

	if sessProvider := AppConfig.String("SessionProvider"); sessProvider != "" {
		SessionProvider = sessProvider
	}

	if sessName := AppConfig.String("SessionName"); sessName != "" {
		SessionName = sessName
	}

	if sesssavepath := AppConfig.String("SessionSavePath"); sesssavepath != "" {
		SessionSavePath = sesssavepath
	}

	if sessMaxLifeTime, err := AppConfig.Int64("SessionGCMaxLifetime"); err == nil && sessMaxLifeTime != 0 {
		SessionGCMaxLifetime = sessMaxLifeTime
	}

	if sesscookielifetime, err := AppConfig.Int("SessionCookieLifeTime"); err == nil && sesscookielifetime != 0 {
		SessionCookieLifeTime = sesscookielifetime
	}

	if usefcgi, err := AppConfig.Bool("UseFcgi"); err == nil {
		UseFcgi = usefcgi
	}

	if enablegzip, err := AppConfig.Bool("EnableGzip"); err == nil {
		EnableGzip = enablegzip
	}

	if directoryindex, err := AppConfig.Bool("DirectoryIndex"); err == nil {
		DirectoryIndex = directoryindex
	}

	if timeout, err := AppConfig.Int64("HttpServerTimeOut"); err == nil {
		HttpServerTimeOut = timeout
	}

	if errorsshow, err := AppConfig.Bool("ErrorsShow"); err == nil {
		ErrorsShow = errorsshow
	}

	if copyrequestbody, err := AppConfig.Bool("CopyRequestBody"); err == nil {
		CopyRequestBody = copyrequestbody
	}

	if xsrfkey := AppConfig.String("XSRFKEY"); xsrfkey != "" {
		XSRFKEY = xsrfkey
	}

	if enablexsrf, err := AppConfig.Bool("EnableXSRF"); err == nil {
		EnableXSRF = enablexsrf
	}

	if expire, err := AppConfig.Int("XSRFExpire"); err == nil {
		XSRFExpire = expire
	}

	if tplleft := AppConfig.String("TemplateLeft"); tplleft != "" {
		TemplateLeft = tplleft
	}

	if tplright := AppConfig.String("TemplateRight"); tplright != "" {
		TemplateRight = tplright
	}

	if httptls, err := AppConfig.Bool("EnableHttpTLS"); err == nil {
		EnableHttpTLS = httptls
	}

	if httpsport, err := AppConfig.Int("HttpsPort"); err == nil {
		HttpsPort = httpsport
	}

	if certfile := AppConfig.String("HttpCertFile"); certfile != "" {
		HttpCertFile = certfile
	}

	if keyfile := AppConfig.String("HttpKeyFile"); keyfile != "" {
		HttpKeyFile = keyfile
	}

	if serverName := AppConfig.String("BeegoServerName"); serverName != "" {
		BeegoServerName = serverName
	}

	if flashname := AppConfig.String("FlashName"); flashname != "" {
		FlashName = flashname
	}

	if flashseperator := AppConfig.String("FlashSeperator"); flashseperator != "" {
		FlashSeperator = flashseperator
	}

	if sd := AppConfig.String("StaticDir"); sd != "" {
		for k := range StaticDir {
			delete(StaticDir, k)
		}
		sds := strings.Fields(sd)
		for _, v := range sds {
			if url2fsmap := strings.SplitN(v, ":", 2); len(url2fsmap) == 2 {
				StaticDir["/"+strings.TrimRight(url2fsmap[0], "/")] = url2fsmap[1]
			} else {
				StaticDir["/"+strings.TrimRight(url2fsmap[0], "/")] = url2fsmap[0]
			}
		}
	}

	if sgz := AppConfig.String("StaticExtensionsToGzip"); sgz != "" {
		extensions := strings.Split(sgz, ",")
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

	if enableadmin, err := AppConfig.Bool("EnableAdmin"); err == nil {
		EnableAdmin = enableadmin
	}

	if adminhttpaddr := AppConfig.String("AdminHttpAddr"); adminhttpaddr != "" {
		AdminHttpAddr = adminhttpaddr
	}

	if adminhttpport, err := AppConfig.Int("AdminHttpPort"); err == nil {
		AdminHttpPort = adminhttpport
	}

	if enabledocs, err := AppConfig.Bool("EnableDocs"); err == nil {
		EnableDocs = enabledocs
	}

	if casesensitive, err := AppConfig.Bool("RouterCaseSensitive"); err == nil {
		RouterCaseSensitive = casesensitive
	}
	if graceful, err := AppConfig.Bool("Graceful"); err == nil {
		Graceful = graceful
	}
	return nil
}
