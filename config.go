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
	"reflect"
	"runtime"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/config/base"
	_ "github.com/astaxie/beego/config/file"
	"github.com/astaxie/beego/context"
	_ "github.com/astaxie/beego/encoder/yaml"
	"github.com/astaxie/beego/logs"
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
	RecoverFunc         func(*context.Context)
	CopyRequestBody     bool
	EnableGzip          bool
	MaxMemory           int64
	EnableErrorsShow    bool
	EnableErrorsRender  bool
	Listen              Listen
	WebConfig           WebConfig
	Log                 LogConfig
}

// Listen holds for http and https related config
type Listen struct {
	Graceful          bool // Graceful means use graceful module to start the server
	ServerTimeOut     int64
	ListenTCP4        bool
	EnableHTTP        bool
	HTTPAddr          string
	HTTPPort          int
	AutoTLS           bool
	Domains           []string
	TLSCacheDir       string
	EnableHTTPS       bool
	EnableMutualHTTPS bool
	HTTPSAddr         string
	HTTPSPort         int
	HTTPSCertFile     string
	HTTPSKeyFile      string
	TrustCaFile       string
	EnableAdmin       bool
	AdminAddr         string
	AdminPort         int
	EnableFcgi        bool
	EnableStdIo       bool // EnableStdIo works with EnableFcgi Use FCGI via standard I/O
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
	SessionOn                    bool
	SessionProvider              string
	SessionName                  string
	SessionGCMaxLifetime         int64
	SessionProviderConfig        string
	SessionCookieLifeTime        int
	SessionAutoSetCookie         bool
	SessionDomain                string
	SessionDisableHTTPOnly       bool // used to allow for cross domain cookies/javascript cookies.
	SessionEnableSidInHTTPHeader bool // enable store/get the sessionId into/from http headers
	SessionNameInHTTPHeader      string
	SessionEnableSidInURLQuery   bool // enable get the sessionId from Url Query params
}

// LogConfig holds Log related config
type LogConfig struct {
	AccessLogs       bool
	EnableStaticLogs bool   //log static files requests default: false
	AccessLogsFormat string //access log format: JSON_FORMAT, APACHE_FORMAT or empty string
	FileLineNum      bool
	Outputs          map[string]string // Store Adaptor : config
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
	// appConfigProvider is the provider for the config, default is file
	appConfigProvider = "file"
)

func init() {
	BConfig = newBConfig()
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var filename = "app.yaml"
	if os.Getenv("BEEGO_RUNMODE") != "" {
		filename = os.Getenv("BEEGO_RUNMODE") + ".app.yaml"
	}
	appConfigPath = filepath.Join(workPath, "conf", filename)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "conf", filename)
		if !utils.FileExists(appConfigPath) {
			AppConfig = &beegoAppConfig{Configer: base.NewBaseConfig()}
			return
		}
	}
	if err = parseConfig(appConfigPath); err != nil {
		panic(err)
	}
}

func recoverPanic(ctx *context.Context) {
	if err := recover(); err != nil {
		if err == ErrAbort {
			return
		}
		if !BConfig.RecoverPanic {
			panic(err)
		}
		if BConfig.EnableErrorsShow {
			if _, ok := ErrorMaps[fmt.Sprint(err)]; ok {
				exception(fmt.Sprint(err), ctx)
				return
			}
		}
		var stack string
		logs.Critical("the request url is ", ctx.Input.URL())
		logs.Critical("Handler crashed with error", err)
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		}
		if BConfig.RunMode == DEV && BConfig.EnableErrorsRender {
			showErr(err, ctx, stack)
		}
		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}
	}
}

func newBConfig() *Config {
	return &Config{
		AppName:             "beego",
		RunMode:             PROD,
		RouterCaseSensitive: true,
		ServerName:          "beegoServer:" + VERSION,
		RecoverPanic:        true,
		RecoverFunc:         recoverPanic,
		CopyRequestBody:     false,
		EnableGzip:          false,
		MaxMemory:           1 << 26, //64MB
		EnableErrorsShow:    true,
		EnableErrorsRender:  true,
		Listen: Listen{
			Graceful:      false,
			ServerTimeOut: 0,
			ListenTCP4:    false,
			EnableHTTP:    true,
			AutoTLS:       false,
			Domains:       []string{},
			TLSCacheDir:   ".",
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
				SessionOn:                    false,
				SessionProvider:              "memory",
				SessionName:                  "beegosessionID",
				SessionGCMaxLifetime:         3600,
				SessionProviderConfig:        "",
				SessionDisableHTTPOnly:       false,
				SessionCookieLifeTime:        0, //set cookie default is the browser life
				SessionAutoSetCookie:         true,
				SessionDomain:                "",
				SessionEnableSidInHTTPHeader: false, // enable store/get the sessionId into/from http headers
				SessionNameInHTTPHeader:      "Beegosessionid",
				SessionEnableSidInURLQuery:   false, // enable get the sessionId from Url Query params
			},
		},
		Log: LogConfig{
			AccessLogs:       false,
			EnableStaticLogs: false,
			AccessLogsFormat: "APACHE_FORMAT",
			FileLineNum:      true,
			Outputs:          map[string]string{"console": ""},
		},
	}
}

// now only support ini, next will support json.
func parseConfig(appConfigPath string) (err error) {
	AppConfig, err = newAppConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	return assignConfig(AppConfig)
}

func assignConfig(ac config.Configer) error {
	for _, i := range []interface{}{BConfig, &BConfig.Listen, &BConfig.WebConfig, &BConfig.Log, &BConfig.WebConfig.Session} {
		assignSingleConfig(i, ac)
	}
	// set the run mode first
	if envRunMode := os.Getenv("BEEGO_RUNMODE"); envRunMode != "" {
		BConfig.RunMode = envRunMode
	} else if runMode := ac.String("RunMode"); runMode != "" {
		BConfig.RunMode = runMode
	}

	if sd := ac.String("StaticDir"); sd != "" {
		BConfig.WebConfig.StaticDir = map[string]string{}
		sds := strings.Fields(sd)
		for _, v := range sds {
			if url2fsmap := strings.SplitN(v, ":", 2); len(url2fsmap) == 2 {
				BConfig.WebConfig.StaticDir["/"+strings.Trim(url2fsmap[0], "/")] = url2fsmap[1]
			} else {
				BConfig.WebConfig.StaticDir["/"+strings.Trim(url2fsmap[0], "/")] = url2fsmap[0]
			}
		}
	}

	if sgz := ac.String("StaticExtensionsToGzip"); sgz != "" {
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

	if lo := ac.String("LogOutputs"); lo != "" {
		// if lo is not nil or empty
		// means user has set his own LogOutputs
		// clear the default setting to BConfig.Log.Outputs
		BConfig.Log.Outputs = make(map[string]string)
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
	logs.Reset()
	for adaptor, conf := range BConfig.Log.Outputs {
		err := logs.SetLogger(adaptor, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s with the config %q got err:%s", adaptor, conf, err.Error()))
		}
	}
	logs.SetLogFuncCall(BConfig.Log.FileLineNum)

	return nil
}

func assignSingleConfig(p interface{}, ac config.Configer) {
	pt := reflect.TypeOf(p)
	if pt.Kind() != reflect.Ptr {
		return
	}
	pt = pt.Elem()
	if pt.Kind() != reflect.Struct {
		return
	}
	pv := reflect.ValueOf(p).Elem()

	for i := 0; i < pt.NumField(); i++ {
		pf := pv.Field(i)
		if !pf.CanSet() {
			continue
		}
		name := pt.Field(i).Name
		switch pf.Kind() {
		case reflect.String:
			pf.SetString(ac.DefaultString(name, pf.String()))
		case reflect.Int, reflect.Int64:
			pf.SetInt(ac.DefaultInt64(pf.Int(), name))
		case reflect.Bool:
			valDef :=  pf.Bool()
			pf.SetBool(ac.DefaultBool(valDef, name))
		case reflect.Struct:
		default:
			//do nothing here
		}
	}
}

// LoadAppConfig allow developer to apply a config
func LoadAppConfig(adapterName, configName string) error {
	appConfigProvider = adapterName
	return parseConfig(configName)
}

type beegoAppConfig struct {
	config.Configer
}

func newAppConfig(appConfigProvider, appConfigPath string) (*beegoAppConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)

	if err != nil {
		return nil, err
	}
	return &beegoAppConfig{ac}, nil
}
