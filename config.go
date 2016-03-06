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
	"os"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/session"
	"github.com/astaxie/beego/utils"
)

// Config is the main struct for BConfig
type Config struct {
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

// Listen holds for http and https related config
type Listen struct {
	Graceful      bool // Graceful means use graceful module to start the server
	ServerTimeOut int64
	ListenTCP4    bool
	EnableHTTP    bool
	HTTPAddr      string
	HTTPPort      int
	EnableHTTPS   bool
	HTTPSAddr     string
	HTTPSPort     int
	HTTPSCertFile string
	HTTPSKeyFile  string
	EnableAdmin   bool
	AdminAddr     string
	AdminPort     int
	EnableFcgi    bool
	EnableStdIo   bool // EnableStdIo works with EnableFcgi Use FCGI via standard I/O
}

// WebConfig holds web related config
type WebConfig struct {
	AutoRender             bool
	EnableDocs             bool
	FlashName              string
	FlashSeparator         string
	DirectoryIndex         bool
	StaticDir              map[string]string
	StaticExtensionsToGzip []string
	TemplateLeft           string
	TemplateRight          string
	ViewsPath              string
	EnableXSRF             bool
	XSRFKey                string
	XSRFExpire             int
	Session                SessionConfig
}

// SessionConfig holds session related config
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

// LogConfig holds Log related config
type LogConfig struct {
	AccessLogs  bool
	FileLineNum bool
	Outputs     map[string]string // Store Adaptor : config
}

var (
	// BConfig is the default config for Application
	BConfig *Config
	// AppConfig is the instance of Config, store the config information from file
	AppConfig *beegoAppConfig
	// AppPath is the absolute path to the app
	AppPath string
	// GlobalSessions is the instance for the session manager
	GlobalSessions *session.Manager

	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "ini"
)

func init() {
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	os.Chdir(AppPath)

	BConfig = &Config{
		AppName:             "beego",
		RunMode:             DEV,
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
			EnableHTTP:    true,
			HTTPAddr:      "",
			HTTPPort:      8080,
			EnableHTTPS:   false,
			HTTPSAddr:     "",
			HTTPSPort:     10443,
			HTTPSCertFile: "",
			HTTPSKeyFile:  "",
			EnableAdmin:   false,
			AdminAddr:     "",
			AdminPort:     8088,
			EnableFcgi:    false,
			EnableStdIo:   false,
		},
		WebConfig: WebConfig{
			AutoRender:             true,
			EnableDocs:             false,
			FlashName:              "BEEGO_FLASH",
			FlashSeparator:         "BEEGOFLASH",
			DirectoryIndex:         false,
			StaticDir:              map[string]string{"/static": "static"},
			StaticExtensionsToGzip: []string{".css", ".js"},
			TemplateLeft:           "{{",
			TemplateRight:          "}}",
			ViewsPath:              "views",
			EnableXSRF:             false,
			XSRFKey:                "beegoxsrf",
			XSRFExpire:             0,
			Session: SessionConfig{
				SessionOn:             false,
				SessionProvider:       "memory",
				SessionName:           "beegosessionID",
				SessionGCMaxLifetime:  3600,
				SessionProviderConfig: "",
				SessionCookieLifeTime: 0, //set cookie default is the browser life
				SessionAutoSetCookie:  true,
				SessionDomain:         "",
			},
		},
		Log: LogConfig{
			AccessLogs:  false,
			FileLineNum: true,
			Outputs:     map[string]string{"console": ""},
		},
	}

	appConfigPath = filepath.Join(AppPath, "conf", "app.conf")
	if !utils.FileExists(appConfigPath) {
		AppConfig = &beegoAppConfig{innerConfig: config.NewFakeConfig()}
		return
	}

	if err := parseConfig(appConfigPath); err != nil {
		panic(err)
	}
}

// now only support ini, next will support json.
func parseConfig(appConfigPath string) (err error) {
	AppConfig, err = newAppConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	// set the run mode first
	if envRunMode := os.Getenv("BEEGO_RUNMODE"); envRunMode != "" {
		BConfig.RunMode = envRunMode
	} else if runMode := AppConfig.String("RunMode"); runMode != "" {
		BConfig.RunMode = runMode
	}

	BConfig.AppName = AppConfig.DefaultString("AppName", BConfig.AppName)
	BConfig.RecoverPanic = AppConfig.DefaultBool("RecoverPanic", BConfig.RecoverPanic)
	BConfig.RouterCaseSensitive = AppConfig.DefaultBool("RouterCaseSensitive", BConfig.RouterCaseSensitive)
	BConfig.ServerName = AppConfig.DefaultString("ServerName", BConfig.ServerName)
	BConfig.EnableGzip = AppConfig.DefaultBool("EnableGzip", BConfig.EnableGzip)
	BConfig.EnableErrorsShow = AppConfig.DefaultBool("EnableErrorsShow", BConfig.EnableErrorsShow)
	BConfig.CopyRequestBody = AppConfig.DefaultBool("CopyRequestBody", BConfig.CopyRequestBody)
	BConfig.MaxMemory = AppConfig.DefaultInt64("MaxMemory", BConfig.MaxMemory)
	BConfig.Listen.Graceful = AppConfig.DefaultBool("Graceful", BConfig.Listen.Graceful)
	BConfig.Listen.HTTPAddr = AppConfig.String("HTTPAddr")
	BConfig.Listen.HTTPPort = AppConfig.DefaultInt("HTTPPort", BConfig.Listen.HTTPPort)
	BConfig.Listen.ListenTCP4 = AppConfig.DefaultBool("ListenTCP4", BConfig.Listen.ListenTCP4)
	BConfig.Listen.EnableHTTP = AppConfig.DefaultBool("EnableHTTP", BConfig.Listen.EnableHTTP)
	BConfig.Listen.EnableHTTPS = AppConfig.DefaultBool("EnableHTTPS", BConfig.Listen.EnableHTTPS)
	BConfig.Listen.HTTPSAddr = AppConfig.DefaultString("HTTPSAddr", BConfig.Listen.HTTPSAddr)
	BConfig.Listen.HTTPSPort = AppConfig.DefaultInt("HTTPSPort", BConfig.Listen.HTTPSPort)
	BConfig.Listen.HTTPSCertFile = AppConfig.DefaultString("HTTPSCertFile", BConfig.Listen.HTTPSCertFile)
	BConfig.Listen.HTTPSKeyFile = AppConfig.DefaultString("HTTPSKeyFile", BConfig.Listen.HTTPSKeyFile)
	BConfig.Listen.EnableAdmin = AppConfig.DefaultBool("EnableAdmin", BConfig.Listen.EnableAdmin)
	BConfig.Listen.AdminAddr = AppConfig.DefaultString("AdminAddr", BConfig.Listen.AdminAddr)
	BConfig.Listen.AdminPort = AppConfig.DefaultInt("AdminPort", BConfig.Listen.AdminPort)
	BConfig.Listen.EnableFcgi = AppConfig.DefaultBool("EnableFcgi", BConfig.Listen.EnableFcgi)
	BConfig.Listen.EnableStdIo = AppConfig.DefaultBool("EnableStdIo", BConfig.Listen.EnableStdIo)
	BConfig.Listen.ServerTimeOut = AppConfig.DefaultInt64("ServerTimeOut", BConfig.Listen.ServerTimeOut)
	BConfig.WebConfig.AutoRender = AppConfig.DefaultBool("AutoRender", BConfig.WebConfig.AutoRender)
	BConfig.WebConfig.ViewsPath = AppConfig.DefaultString("ViewsPath", BConfig.WebConfig.ViewsPath)
	BConfig.WebConfig.DirectoryIndex = AppConfig.DefaultBool("DirectoryIndex", BConfig.WebConfig.DirectoryIndex)
	BConfig.WebConfig.FlashName = AppConfig.DefaultString("FlashName", BConfig.WebConfig.FlashName)
	BConfig.WebConfig.FlashSeparator = AppConfig.DefaultString("FlashSeparator", BConfig.WebConfig.FlashSeparator)
	BConfig.WebConfig.EnableDocs = AppConfig.DefaultBool("EnableDocs", BConfig.WebConfig.EnableDocs)
	BConfig.WebConfig.XSRFKey = AppConfig.DefaultString("XSRFKEY", BConfig.WebConfig.XSRFKey)
	BConfig.WebConfig.EnableXSRF = AppConfig.DefaultBool("EnableXSRF", BConfig.WebConfig.EnableXSRF)
	BConfig.WebConfig.XSRFExpire = AppConfig.DefaultInt("XSRFExpire", BConfig.WebConfig.XSRFExpire)
	BConfig.WebConfig.TemplateLeft = AppConfig.DefaultString("TemplateLeft", BConfig.WebConfig.TemplateLeft)
	BConfig.WebConfig.TemplateRight = AppConfig.DefaultString("TemplateRight", BConfig.WebConfig.TemplateRight)
	BConfig.WebConfig.Session.SessionOn = AppConfig.DefaultBool("SessionOn", BConfig.WebConfig.Session.SessionOn)
	BConfig.WebConfig.Session.SessionProvider = AppConfig.DefaultString("SessionProvider", BConfig.WebConfig.Session.SessionProvider)
	BConfig.WebConfig.Session.SessionName = AppConfig.DefaultString("SessionName", BConfig.WebConfig.Session.SessionName)
	BConfig.WebConfig.Session.SessionProviderConfig = AppConfig.DefaultString("SessionProviderConfig", BConfig.WebConfig.Session.SessionProviderConfig)
	BConfig.WebConfig.Session.SessionGCMaxLifetime = AppConfig.DefaultInt64("SessionGCMaxLifetime", BConfig.WebConfig.Session.SessionGCMaxLifetime)
	BConfig.WebConfig.Session.SessionCookieLifeTime = AppConfig.DefaultInt("SessionCookieLifeTime", BConfig.WebConfig.Session.SessionCookieLifeTime)
	BConfig.WebConfig.Session.SessionAutoSetCookie = AppConfig.DefaultBool("SessionAutoSetCookie", BConfig.WebConfig.Session.SessionAutoSetCookie)
	BConfig.WebConfig.Session.SessionDomain = AppConfig.DefaultString("SessionDomain", BConfig.WebConfig.Session.SessionDomain)
	BConfig.Log.AccessLogs = AppConfig.DefaultBool("LogAccessLogs", BConfig.Log.AccessLogs)
	BConfig.Log.FileLineNum = AppConfig.DefaultBool("LogFileLineNum", BConfig.Log.FileLineNum)

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

	if lo := AppConfig.String("LogOutputs"); lo != "" {
		los := strings.Split(lo, ";")
		for _, v := range los {
			if logType2Config := strings.SplitN(v, ",", 2); len(logType2Config) == 2 {
				BConfig.Log.Outputs[logType2Config[0]] = logType2Config[1]
			} else {
				continue
			}
		}
	}

	//init log
	BeeLogger.Reset()
	for adaptor, config := range BConfig.Log.Outputs {
		err = BeeLogger.SetLogger(adaptor, config)
		if err != nil {
			fmt.Printf("%s with the config `%s` got err:%s\n", adaptor, config, err)
		}
	}
	SetLogFuncCall(BConfig.Log.FileLineNum)

	return nil
}

// LoadAppConfig allow developer to apply a config file
func LoadAppConfig(adapterName, configPath string) error {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	if !utils.FileExists(absConfigPath) {
		return fmt.Errorf("the target config file: %s don't exist", configPath)
	}

	if absConfigPath == appConfigPath {
		return nil
	}

	appConfigPath = absConfigPath
	appConfigProvider = adapterName

	return parseConfig(appConfigPath)
}

type beegoAppConfig struct {
	innerConfig config.Configer
}

func newAppConfig(appConfigProvider, appConfigPath string) (*beegoAppConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &beegoAppConfig{ac}, nil
}

func (b *beegoAppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(BConfig.RunMode+"::"+key, val); err != nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *beegoAppConfig) String(key string) string {
	if v := b.innerConfig.String(BConfig.RunMode + "::" + key); v != "" {
		return v
	}
	return b.innerConfig.String(key)
}

func (b *beegoAppConfig) Strings(key string) []string {
	if v := b.innerConfig.Strings(BConfig.RunMode + "::" + key); v[0] != "" {
		return v
	}
	return b.innerConfig.Strings(key)
}

func (b *beegoAppConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

func (b *beegoAppConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

func (b *beegoAppConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

func (b *beegoAppConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

func (b *beegoAppConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
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
