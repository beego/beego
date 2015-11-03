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
	// AccessLogs represent whether output the access logs, default is false
	AccessLogs bool
	// AdminHTTPAddr is address for admin
	AdminHTTPAddr string
	// AdminHTTPPort is listens port for admin
	AdminHTTPPort int
	// AppConfig is the instance of Config, store the config information from file
	AppConfig *beegoAppConfig
	// AppName represent Application name, always the project folder name
	AppName string
	// AppPath is the path to the application
	AppPath string
	// AppConfigPath is the path to the config files
	AppConfigPath string
	// AppConfigProvider is the provider for the config, default is ini
	AppConfigProvider string
	// AutoRender is a flag of render template automatically. It's always turn off in API application
	// default is true
	AutoRender bool
	// BeegoServerName exported in response header.
	BeegoServerName string
	// CopyRequestBody is just useful for raw request body in context. default is false
	CopyRequestBody bool
	// DirectoryIndex wheather display directory index. default is false.
	DirectoryIndex bool
	// EnableAdmin means turn on admin module to log every request info.
	EnableAdmin bool
	// EnableDocs enable generate docs & server docs API Swagger
	EnableDocs bool
	// EnableErrorsShow wheather show errors in page. if true, show error and trace info in page rendered with error template.
	EnableErrorsShow bool
	// EnableFcgi turn on the fcgi Listen, default is false
	EnableFcgi bool
	// EnableGzip means gzip the response
	EnableGzip bool
	// EnableHTTPListen represent whether turn on the HTTP, default is true
	EnableHTTPListen bool
	// EnableHTTPTLS represent whether turn on the HTTPS, default is true
	EnableHTTPTLS bool
	// EnableStdIo works with EnableFcgi Use FCGI via standard I/O
	EnableStdIo bool
	// EnableXSRF whether turn on xsrf. default is false
	EnableXSRF bool
	// FlashName is the name of the flash variable found in response header and cookie
	FlashName string
	// FlashSeperator used to seperate flash key:value, default is BEEGOFLASH
	FlashSeperator string
	// GlobalSessions is the instance for the session manager
	GlobalSessions *session.Manager
	// Graceful means use graceful module to start the server
	Graceful bool
	// workPath is always the same as AppPath, but sometime when it started with other
	// program, like supervisor
	workPath string
	// ListenTCP4 represent only Listen in TCP4, default is false
	ListenTCP4 bool
	// MaxMemory The whole request body is parsed and up to a total of maxMemory
	// bytes of its file parts are stored in memory, with the remainder stored on disk in temporary files
	MaxMemory int64
	// HTTPAddr is the TCP network address addr for HTTP
	HTTPAddr string
	// HTTPPort is listens port for HTTP
	HTTPPort int
	// HTTPSPort is listens port for HTTPS
	HTTPSPort int
	// HTTPCertFile is the path to certificate file
	HTTPCertFile string
	// HTTPKeyFile is the path to private key file
	HTTPKeyFile string
	// HTTPServerTimeOut HTTP server timeout. default is 0, no timeout
	HTTPServerTimeOut int64
	// RecoverPanic is a flag for auto recover panic, default is true
	RecoverPanic bool
	// RouterCaseSensitive means whether router case sensitive, default is true
	RouterCaseSensitive bool
	// RunMode represent the staging, "dev" or "prod"
	RunMode string
	// SessionOn means whether turn on the session auto when application started. default is false.
	SessionOn bool
	// SessionProvider means session provider, e.q memory, mysql, redis,etc.
	SessionProvider string
	// SessionName is the cookie name when saving session id into cookie.
	SessionName string
	// SessionGCMaxLifetime for auto cleaning expired session.
	SessionGCMaxLifetime int64
	// SessionProviderConfig is for the provider config, define save path or connection info.
	SessionProviderConfig string
	// SessionCookieLifeTime means the life time of session id in cookie.
	SessionCookieLifeTime int
	// SessionAutoSetCookie auto setcookie
	SessionAutoSetCookie bool
	// SessionDomain means the cookie domain default is empty
	SessionDomain string
	// StaticDir store the static path, key is path, value is the folder
	StaticDir map[string]string
	// StaticExtensionsToGzip stores the extensions which need to gzip(.js,.css,etc)
	StaticExtensionsToGzip []string
	// TemplateCache store the caching template
	TemplateCache map[string]*template.Template
	// TemplateLeft left delimiter
	TemplateLeft string
	// TemplateRight right delimiter
	TemplateRight string
	// ViewsPath means the template folder
	ViewsPath string
	// XSRFKEY xsrf hash salt string.
	XSRFKEY string
	// XSRFExpire is the expiry of xsrf value.
	XSRFExpire int
)

type beegoAppConfig struct {
	innerConfig config.Configer
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
	EnableHTTPListen = true //default enable http Listen

	HTTPAddr = ""
	HTTPPort = 8080

	HTTPSPort = 10443

	AppName = "beego"

	RunMode = "dev" //default runmod

	AutoRender = true

	RecoverPanic = true

	ViewsPath = "views"

	SessionOn = false
	SessionProvider = "memory"
	SessionName = "beegosessionID"
	SessionGCMaxLifetime = 3600
	SessionProviderConfig = ""
	SessionCookieLifeTime = 0 //set cookie default is the brower life
	SessionAutoSetCookie = true

	MaxMemory = 1 << 26 //64MB

	HTTPServerTimeOut = 0

	EnableErrorsShow = true

	XSRFKEY = "beegoxsrf"
	XSRFExpire = 0

	TemplateLeft = "{{"
	TemplateRight = "}}"

	BeegoServerName = "beegoServer:" + VERSION

	EnableAdmin = false
	AdminHTTPAddr = "127.0.0.1"
	AdminHTTPPort = 8088

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

	HTTPAddr = AppConfig.String("HTTPAddr")

	if v, err := AppConfig.Int("HTTPPort"); err == nil {
		HTTPPort = v
	}

	if v, err := AppConfig.Bool("ListenTCP4"); err == nil {
		ListenTCP4 = v
	}

	if v, err := AppConfig.Bool("EnableHTTPListen"); err == nil {
		EnableHTTPListen = v
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

	if sessProvConfig := AppConfig.String("SessionProviderConfig"); sessProvConfig != "" {
		SessionProviderConfig = sessProvConfig
	}

	if sessMaxLifeTime, err := AppConfig.Int64("SessionGCMaxLifetime"); err == nil && sessMaxLifeTime != 0 {
		SessionGCMaxLifetime = sessMaxLifeTime
	}

	if sesscookielifetime, err := AppConfig.Int("SessionCookieLifeTime"); err == nil && sesscookielifetime != 0 {
		SessionCookieLifeTime = sesscookielifetime
	}

	if enableFcgi, err := AppConfig.Bool("EnableFcgi"); err == nil {
		EnableFcgi = enableFcgi
	}

	if enablegzip, err := AppConfig.Bool("EnableGzip"); err == nil {
		EnableGzip = enablegzip
	}

	if directoryindex, err := AppConfig.Bool("DirectoryIndex"); err == nil {
		DirectoryIndex = directoryindex
	}

	if timeout, err := AppConfig.Int64("HTTPServerTimeOut"); err == nil {
		HTTPServerTimeOut = timeout
	}

	if errorsshow, err := AppConfig.Bool("EnableErrorsShow"); err == nil {
		EnableErrorsShow = errorsshow
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

	if httptls, err := AppConfig.Bool("EnableHTTPTLS"); err == nil {
		EnableHTTPTLS = httptls
	}

	if httpsport, err := AppConfig.Int("HTTPSPort"); err == nil {
		HTTPSPort = httpsport
	}

	if certfile := AppConfig.String("HTTPCertFile"); certfile != "" {
		HTTPCertFile = certfile
	}

	if keyfile := AppConfig.String("HTTPKeyFile"); keyfile != "" {
		HTTPKeyFile = keyfile
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

	if adminhttpaddr := AppConfig.String("AdminHTTPAddr"); adminhttpaddr != "" {
		AdminHTTPAddr = adminhttpaddr
	}

	if adminhttpport, err := AppConfig.Int("AdminHTTPPort"); err == nil {
		AdminHTTPPort = adminhttpport
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
