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
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/session"
	"github.com/astaxie/beego/utils"
)

type BeegoConfig struct {
	AppName             string //Application name
	RunMode             string //Running Mode: dev | prod
	RouterCaseSensitive bool
	ServerName          string
	RecoverPanic        bool
	CopyRequestBody     bool
	EnableGzip          bool
	MaxMemory           int64
	EnableErrorsShow    bool
	Listen              Listen
	WebConfig           WebConfig
	Log                 LogConfig
}

type Listen struct {
	Graceful      bool // Graceful means use graceful module to start the server
	ServerTimeOut int64
	ListenTCP4    bool
	HTTPEnable    bool
	HTTPAddr      string
	HTTPPort      int
	HTTPSEnable   bool
	HTTPSAddr     string
	HTTPSPort     int
	HTTPSCertFile string
	HTTPSKeyFile  string
	AdminEnable   bool
	AdminAddr     string
	AdminPort     int
	EnableFcgi    bool
	EnableStdIo   bool // EnableStdIo works with EnableFcgi Use FCGI via standard I/O
}

type WebConfig struct {
	AutoRender             bool
	EnableDocs             bool
	FlashName              string
	FlashSeperator         string
	DirectoryIndex         bool
	StaticDir              map[string]string
	StaticExtensionsToGzip []string
	TemplateLeft           string
	TemplateRight          string
	ViewsPath              string
	EnableXSRF             bool
	XSRFKEY                string
	XSRFExpire             int
	Session                SessionConfig
}

type SessionConfig struct {
	SessionOn             bool
	SessionProvider       string
	SessionName           string
	SessionGCMaxLifetime  int64
	SessionProviderConfig string
	SessionCookieLifeTime int
	SessionAutoSetCookie  bool
	SessionDomain         string
}

type LogConfig struct {
	AccessLogs  bool
	FileLineNum bool
	Output      map[string]string // Store Adaptor : config
}

var (
	// BConfig is the default config for Application
	BConfig *BeegoConfig
	// AppConfig is the instance of Config, store the config information from file
	AppConfig *beegoAppConfig
	// AppConfigPath is the path to the config files
	AppConfigPath string
	// AppConfigProvider is the provider for the config, default is ini
	AppConfigProvider = "ini"
	// TemplateCache stores template caching
	TemplateCache map[string]*template.Template
	// GlobalSessions is the instance for the session manager
	GlobalSessions *session.Manager
)

func init() {
	BConfig = &BeegoConfig{
		AppName:             "beego",
		RunMode:             "dev",
		RouterCaseSensitive: true,
		ServerName:          "beegoServer:" + VERSION,
		RecoverPanic:        true,
		CopyRequestBody:     false,
		EnableGzip:          false,
		MaxMemory:           1 << 26, //64MB
		EnableErrorsShow:    true,
		Listen: Listen{
			Graceful:      false,
			ServerTimeOut: 0,
			ListenTCP4:    false,
			HTTPEnable:    true,
			HTTPAddr:      "",
			HTTPPort:      8080,
			HTTPSEnable:   false,
			HTTPSAddr:     "",
			HTTPSPort:     10443,
			HTTPSCertFile: "",
			HTTPSKeyFile:  "",
			AdminEnable:   false,
			AdminAddr:     "",
			AdminPort:     8088,
			EnableFcgi:    false,
			EnableStdIo:   false,
		},
		WebConfig: WebConfig{
			AutoRender:             true,
			EnableDocs:             false,
			FlashName:              "BEEGO_FLASH",
			FlashSeperator:         "BEEGOFLASH",
			DirectoryIndex:         false,
			StaticDir:              map[string]string{"/static": "static"},
			StaticExtensionsToGzip: []string{".css", ".js"},
			TemplateLeft:           "{{",
			TemplateRight:          "}}",
			ViewsPath:              "views",
			EnableXSRF:             false,
			XSRFKEY:                "beegoxsrf",
			XSRFExpire:             0,
			Session: SessionConfig{
				SessionOn:             false,
				SessionProvider:       "memory",
				SessionName:           "beegosessionID",
				SessionGCMaxLifetime:  3600,
				SessionProviderConfig: "",
				SessionCookieLifeTime: 0, //set cookie default is the brower life
				SessionAutoSetCookie:  true,
				SessionDomain:         "",
			},
		},
		Log: LogConfig{
			AccessLogs:  false,
			FileLineNum: true,
			Output:      map[string]string{"console": ""},
		},
	}
}

// ParseConfig parsed default config file.
// now only support ini, next will support json.
func ParseConfig() (err error) {
	if AppConfigPath == "" {
		if utils.FileExists(filepath.Join("conf", "app.conf")) {
			AppConfigPath = filepath.Join("conf", "app.conf")
		} else {
			ac := config.NewFakeConfig()
			AppConfig = &beegoAppConfig{ac}
			return
		}
	}
	AppConfig, err = newAppConfig(AppConfigProvider, AppConfigPath)
	if err != nil {
		return err
	}
	envRunMode := os.Getenv("BEEGO_RUNMODE")
	// set the runmode first
	if envRunMode != "" {
		BConfig.RunMode = envRunMode
	} else if runmode := AppConfig.String("RunMode"); runmode != "" {
		BConfig.RunMode = runmode
	}

	BConfig.Listen.HTTPAddr = AppConfig.String("HTTPAddr")

	if v, err := AppConfig.Int("HTTPPort"); err == nil {
		BConfig.Listen.HTTPPort = v
	}

	if v, err := AppConfig.Bool("ListenTCP4"); err == nil {
		BConfig.Listen.ListenTCP4 = v
	}

	if v, err := AppConfig.Bool("EnableHTTPListen"); err == nil {
		BConfig.Listen.HTTPEnable = v
	}

	if maxmemory, err := AppConfig.Int64("MaxMemory"); err == nil {
		BConfig.MaxMemory = maxmemory
	}

	if appname := AppConfig.String("AppName"); appname != "" {
		BConfig.AppName = appname
	}

	if autorender, err := AppConfig.Bool("AutoRender"); err == nil {
		BConfig.WebConfig.AutoRender = autorender
	}

	if autorecover, err := AppConfig.Bool("RecoverPanic"); err == nil {
		BConfig.RecoverPanic = autorecover
	}

	if views := AppConfig.String("ViewsPath"); views != "" {
		BConfig.WebConfig.ViewsPath = views
	}

	if sessionon, err := AppConfig.Bool("SessionOn"); err == nil {
		BConfig.WebConfig.Session.SessionOn = sessionon
	}

	if sessProvider := AppConfig.String("SessionProvider"); sessProvider != "" {
		BConfig.WebConfig.Session.SessionProvider = sessProvider
	}

	if sessName := AppConfig.String("SessionName"); sessName != "" {
		BConfig.WebConfig.Session.SessionName = sessName
	}

	if sessProvConfig := AppConfig.String("SessionProviderConfig"); sessProvConfig != "" {
		BConfig.WebConfig.Session.SessionProviderConfig = sessProvConfig
	}

	if sessMaxLifeTime, err := AppConfig.Int64("SessionGCMaxLifetime"); err == nil && sessMaxLifeTime != 0 {
		BConfig.WebConfig.Session.SessionGCMaxLifetime = sessMaxLifeTime
	}

	if sesscookielifetime, err := AppConfig.Int("SessionCookieLifeTime"); err == nil && sesscookielifetime != 0 {
		BConfig.WebConfig.Session.SessionCookieLifeTime = sesscookielifetime
	}

	if enableFcgi, err := AppConfig.Bool("EnableFcgi"); err == nil {
		BConfig.Listen.EnableFcgi = enableFcgi
	}

	if enablegzip, err := AppConfig.Bool("EnableGzip"); err == nil {
		BConfig.EnableGzip = enablegzip
	}

	if directoryindex, err := AppConfig.Bool("DirectoryIndex"); err == nil {
		BConfig.WebConfig.DirectoryIndex = directoryindex
	}

	if timeout, err := AppConfig.Int64("HTTPServerTimeOut"); err == nil {
		BConfig.Listen.ServerTimeOut = timeout
	}

	if errorsshow, err := AppConfig.Bool("EnableErrorsShow"); err == nil {
		BConfig.EnableErrorsShow = errorsshow
	}

	if copyrequestbody, err := AppConfig.Bool("CopyRequestBody"); err == nil {
		BConfig.CopyRequestBody = copyrequestbody
	}

	if xsrfkey := AppConfig.String("XSRFKEY"); xsrfkey != "" {
		BConfig.WebConfig.XSRFKEY = xsrfkey
	}

	if enablexsrf, err := AppConfig.Bool("EnableXSRF"); err == nil {
		BConfig.WebConfig.EnableXSRF = enablexsrf
	}

	if expire, err := AppConfig.Int("XSRFExpire"); err == nil {
		BConfig.WebConfig.XSRFExpire = expire
	}

	if tplleft := AppConfig.String("TemplateLeft"); tplleft != "" {
		BConfig.WebConfig.TemplateLeft = tplleft
	}

	if tplright := AppConfig.String("TemplateRight"); tplright != "" {
		BConfig.WebConfig.TemplateRight = tplright
	}

	if httptls, err := AppConfig.Bool("EnableHTTPTLS"); err == nil {
		BConfig.Listen.HTTPSEnable = httptls
	}

	if httpsport, err := AppConfig.Int("HTTPSPort"); err == nil {
		BConfig.Listen.HTTPSPort = httpsport
	}

	if certfile := AppConfig.String("HTTPCertFile"); certfile != "" {
		BConfig.Listen.HTTPSCertFile = certfile
	}

	if keyfile := AppConfig.String("HTTPKeyFile"); keyfile != "" {
		BConfig.Listen.HTTPSKeyFile = keyfile
	}

	if serverName := AppConfig.String("BeegoServerName"); serverName != "" {
		BConfig.ServerName = serverName
	}

	if flashname := AppConfig.String("FlashName"); flashname != "" {
		BConfig.WebConfig.FlashName = flashname
	}

	if flashseperator := AppConfig.String("FlashSeperator"); flashseperator != "" {
		BConfig.WebConfig.FlashSeperator = flashseperator
	}

	if sd := AppConfig.String("StaticDir"); sd != "" {
		for k := range BConfig.WebConfig.StaticDir {
			delete(BConfig.WebConfig.StaticDir, k)
		}
		sds := strings.Fields(sd)
		for _, v := range sds {
			if url2fsmap := strings.SplitN(v, ":", 2); len(url2fsmap) == 2 {
				BConfig.WebConfig.StaticDir["/"+strings.TrimRight(url2fsmap[0], "/")] = url2fsmap[1]
			} else {
				BConfig.WebConfig.StaticDir["/"+strings.TrimRight(url2fsmap[0], "/")] = url2fsmap[0]
			}
		}
	}

	if sgz := AppConfig.String("StaticExtensionsToGzip"); sgz != "" {
		extensions := strings.Split(sgz, ",")
		fileExts := []string{}
		for _, ext := range extensions {
			ext = strings.TrimSpace(ext)
			if ext == "" {
				continue
			}
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			fileExts = append(fileExts, ext)
		}
		if len(fileExts) > 0 {
			BConfig.WebConfig.StaticExtensionsToGzip = fileExts
		}
	}

	if enableadmin, err := AppConfig.Bool("EnableAdmin"); err == nil {
		BConfig.Listen.AdminEnable = enableadmin
	}

	if adminhttpaddr := AppConfig.String("AdminHTTPAddr"); adminhttpaddr != "" {
		BConfig.Listen.AdminAddr = adminhttpaddr
	}

	if adminhttpport, err := AppConfig.Int("AdminHTTPPort"); err == nil {
		BConfig.Listen.AdminPort = adminhttpport
	}

	if enabledocs, err := AppConfig.Bool("EnableDocs"); err == nil {
		BConfig.WebConfig.EnableDocs = enabledocs
	}

	if casesensitive, err := AppConfig.Bool("RouterCaseSensitive"); err == nil {
		BConfig.RouterCaseSensitive = casesensitive
	}
	if graceful, err := AppConfig.Bool("Graceful"); err == nil {
		BConfig.Listen.Graceful = graceful
	}
	return nil
}

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
	err := b.innerConfig.Set(BConfig.RunMode+"::"+key, val)
	if err == nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *beegoAppConfig) String(key string) string {
	v := b.innerConfig.String(BConfig.RunMode + "::" + key)
	if v == "" {
		return b.innerConfig.String(key)
	}
	return v
}

func (b *beegoAppConfig) Strings(key string) []string {
	v := b.innerConfig.Strings(BConfig.RunMode + "::" + key)
	if v[0] == "" {
		return b.innerConfig.Strings(key)
	}
	return v
}

func (b *beegoAppConfig) Int(key string) (int, error) {
	v, err := b.innerConfig.Int(BConfig.RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int(key)
	}
	return v, nil
}

func (b *beegoAppConfig) Int64(key string) (int64, error) {
	v, err := b.innerConfig.Int64(BConfig.RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Int64(key)
	}
	return v, nil
}

func (b *beegoAppConfig) Bool(key string) (bool, error) {
	v, err := b.innerConfig.Bool(BConfig.RunMode + "::" + key)
	if err != nil {
		return b.innerConfig.Bool(key)
	}
	return v, nil
}

func (b *beegoAppConfig) Float(key string) (float64, error) {
	v, err := b.innerConfig.Float(BConfig.RunMode + "::" + key)
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
