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

package web

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/beego/beego/v2"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

// Config is the main struct for BConfig
// TODO after supporting multiple servers, remove common config to somewhere else
type Config struct {
	// AppName
	// @Description Application's name. You'd better set it because we use it to do some logging and tracing
	// @Default beego
	AppName string // Application name
	// RunMode
	// @Description it's the same as environment. In general, we have different run modes.
	// For example, the most common case is using dev, test, prod three environments
	// when you are developing the application, you should set it as dev
	// when you completed coding and want QA to test your code, you should deploy your application to test environment
	// and the RunMode should be set as test
	// when you completed all tests, you want to deploy it to prod, you should set it to prod
	// You should never set RunMode="dev" when you deploy the application to prod
	// because Beego will do more things which need Go SDK and other tools when it found out the RunMode="dev"
	// @Default dev
	RunMode string // Running Mode: dev | prod

	// RouterCaseSensitive
	// @Description If it was true, it means that the router is case sensitive.
	// For example, when you register a router with pattern "/hello",
	// 1. If this is true, and the request URL is "/Hello", it won't match this pattern
	// 2. If this is false and the request URL is "/Hello", it will match this pattern
	// @Default true
	RouterCaseSensitive bool
	// RecoverPanic
	// @Description if it was true, Beego will try to recover from panic when it serves your http request
	// So you should notice that it doesn't mean that Beego will recover all panic cases.
	// @Default true
	RecoverPanic bool
	// CopyRequestBody
	// @Description if it's true, Beego will copy the request body. But if the request body's size > MaxMemory,
	// Beego will return 413 as http status
	// If you are building RESTful API, please set it to true.
	// And if you want to read data from request Body multiple times, please set it to true
	// In general, if you don't meet any performance issue, you could set it to true
	// @Default false
	CopyRequestBody bool
	// EnableGzip
	// @Description If it was true, Beego will try to compress data by using zip algorithm.
	// But there are two points:
	// 1. Only static resources will be compressed
	// 2. Only those static resource which has the extension specified by StaticExtensionsToGzip will be compressed
	// @Default false
	EnableGzip bool
	// EnableErrorsShow
	// @Description If it's true, Beego will show error message to page
	// it will work with ErrorMaps which allows you register some error handler
	// You may want to set it to false when application was deploy to prod environment
	// because you may not want to expose your internal error msg to your users
	// it's a little bit unsafe
	// @Default true
	EnableErrorsShow bool
	// EnableErrorsRender
	// @Description If it's true, it will output the error msg as a page. It's similar to EnableErrorsShow
	// And this configure item only work in dev run mode (see RunMode)
	// @Default true
	EnableErrorsRender bool
	// ServerName
	// @Description server name. For example, in large scale system,
	// you may want to deploy your application to several machines, so that each of them has a server name
	// we suggest you'd better set value because Beego use this to output some DEBUG msg,
	// or integrated with other tools such as tracing, metrics
	// @Default
	ServerName string

	// RecoverFunc
	// @Description when Beego want to recover from panic, it will use this func as callback
	// see RecoverPanic
	// @Default defaultRecoverPanic
	RecoverFunc func(*context.Context, *Config)
	// @Description MaxMemory and MaxUploadSize are used to limit the request body
	// if the request is not uploading file, MaxMemory is the max size of request body
	// if the request is uploading file, MaxUploadSize is the max size of request body
	// if CopyRequestBody is true, this value will be used as the threshold of request body
	// see CopyRequestBody
	// the default value is 1 << 26 (64MB)
	// @Default 67108864
	MaxMemory int64
	// MaxUploadSize
	// @Description  MaxMemory and MaxUploadSize are used to limit the request body
	// if the request is not uploading file, MaxMemory is the max size of request body
	// if the request is uploading file, MaxUploadSize is the max size of request body
	// the default value is 1 << 30 (1GB)
	// @Default 1073741824
	MaxUploadSize int64
	// Listen
	// @Description the configuration about socket or http protocol
	Listen Listen
	// WebConfig
	// @Description the configuration about Web
	WebConfig WebConfig
	// LogConfig
	// @Description log configuration
	Log LogConfig
}

// Listen holds for http and https related config
type Listen struct {
	// Graceful
	// @Description means use graceful module to start the server
	// @Default false
	Graceful bool
	// ListenTCP4
	// @Description if it's true, means that Beego only work for TCP4
	// please check net.Listen function
	// In general, you should not set it to true
	// @Default false
	ListenTCP4 bool
	// EnableHTTP
	// @Description if it's true, Beego will accept HTTP request.
	// But if you want to use HTTPS only, please set it to false
	// see EnableHTTPS
	// @Default true
	EnableHTTP bool
	// AutoTLS
	// @Description If it's true, Beego will use default value to initialize the TLS configure
	// But those values could be override if you have custom value.
	// see Domains, TLSCacheDir
	// @Default false
	AutoTLS bool
	// EnableHTTPS
	// @Description If it's true, Beego will accept HTTPS request.
	// Now, you'd better use HTTPS protocol on prod environment to get better security
	// In prod, the best option is EnableHTTPS=true and EnableHTTP=false
	// see EnableHTTP
	// @Default false
	EnableHTTPS bool
	// EnableMutualHTTPS
	// @Description if it's true, Beego will handle requests on incoming mutual TLS connections
	// see Server.ListenAndServeMutualTLS
	// @Default false
	EnableMutualHTTPS bool
	// EnableAdmin
	// @Description if it's true, Beego will provide admin service.
	// You can visit the admin service via browser.
	// The default port is 8088
	// see AdminPort
	// @Default false
	EnableAdmin bool
	// EnableFcgi
	// @Description
	// @Default false
	EnableFcgi bool
	// EnableStdIo
	// @Description EnableStdIo works with EnableFcgi Use FCGI via standard I/O
	// @Default false
	EnableStdIo bool
	// ServerTimeOut
	// @Description Beego use this as ReadTimeout and WriteTimeout
	// The unit is second.
	// see http.Server.ReadTimeout, WriteTimeout
	// @Default 0
	ServerTimeOut int64
	// HTTPAddr
	// @Description Beego listen to this address when the application start up.
	// @Default ""
	HTTPAddr string
	// HTTPPort
	// @Description Beego listen to this port
	// you'd better change this value when you deploy to prod environment
	// @Default 8080
	HTTPPort int
	// Domains
	// @Description Beego use this to configure TLS. Those domains are "white list" domain
	// @Default []
	Domains []string
	// TLSCacheDir
	// @Description Beego use this as cache dir to store TLS cert data
	// @Default ""
	TLSCacheDir string
	// HTTPSAddr
	// @Description Beego will listen to this address to accept HTTPS request
	// see EnableHTTPS
	// @Default ""
	HTTPSAddr string
	// HTTPSPort
	// @Description  Beego will listen to this port to accept HTTPS request
	// @Default 10443
	HTTPSPort int
	// HTTPSCertFile
	// @Description Beego read this file as cert file
	// When you are using HTTPS protocol, please configure it
	// see HTTPSKeyFile
	// @Default ""
	HTTPSCertFile string
	// HTTPSKeyFile
	// @Description Beego read this file as key file
	// When you are using HTTPS protocol, please configure it
	// see HTTPSCertFile
	// @Default ""
	HTTPSKeyFile string
	// TrustCaFile
	// @Description Beego read this file as CA file
	// @Default ""
	TrustCaFile string
	// AdminAddr
	// @Description Beego will listen to this address to provide admin service
	// In general, it should be the same with your application address, HTTPAddr or HTTPSAddr
	// @Default ""
	AdminAddr string
	// AdminPort
	// @Description  Beego will listen to this port to provide admin service
	// @Default 8088
	AdminPort int
	// @Description Beego use this tls.ClientAuthType to initialize TLS connection
	// The default value is tls.RequireAndVerifyClientCert
	// @Default 4
	ClientAuth int
}

// WebConfig holds web related config
type WebConfig struct {
	// AutoRender
	// @Description If it's true, Beego will render the page based on your template and data
	// In general, keep it as true.
	// But if you are building RESTFul API and you don't have any page,
	// you can set it to false
	// @Default true
	AutoRender bool
	// Deprecated: Beego didn't use it anymore
	EnableDocs bool
	// EnableXSRF
	// @Description If it's true, Beego will help to provide XSRF support
	// But you should notice that, now Beego only work for HTTPS protocol with XSRF
	// because it's not safe if using HTTP protocol
	// And, the cookie storing XSRF token has two more flags HttpOnly and Secure
	// It means that you must use HTTPS protocol and you can not read the token from JS script
	// This is completed different from Beego 1.x because we got many security reports
	// And if you are in dev environment, you could set it to false
	// @Default false
	EnableXSRF bool
	// DirectoryIndex
	// @Description When Beego serves static resources request, it will look up the file.
	// If the file is directory, Beego will try to find the index.html as the response
	// But if the index.html is not exist or it's a directory,
	// Beego will return 403 response if DirectoryIndex is **false**
	// @Default false
	DirectoryIndex bool
	// FlashName
	// @Description the cookie's name when Beego try to store the flash data into cookie
	// @Default BEEGO_FLASH
	FlashName string
	// FlashSeparator
	// @Description When Beego read flash data from request, it uses this as the separator
	// @Default BEEGOFLASH
	FlashSeparator string
	// StaticDir
	// @Description Beego uses this as static resources' root directory.
	// It means that Beego will try to search static resource from this start point
	// It's a map, the key is the path and the value is the directory
	// For example, the default value is /static => static,
	// which means that when Beego got a request with path /static/xxx
	// Beego will try to find the resource from static directory
	// @Default /static => static
	StaticDir map[string]string
	// StaticExtensionsToGzip
	// @Description The static resources with those extension will be compressed if EnableGzip is true
	// @Default [".css", ".js" ]
	StaticExtensionsToGzip []string
	// StaticCacheFileSize
	// @Description If the size of static resource < StaticCacheFileSize, Beego will try to handle it by itself,
	// it means that Beego will compressed the file data (if enable) and cache this file.
	// But if the file size > StaticCacheFileSize, Beego just simply delegate the request to http.ServeFile
	// the default value is 100KB.
	// the max memory size of caching static files is StaticCacheFileSize * StaticCacheFileNum
	// see StaticCacheFileNum
	// @Default 102400
	StaticCacheFileSize int
	// StaticCacheFileNum
	// @Description Beego use it to control the memory usage of caching static resource file
	// If the caching files > StaticCacheFileNum, Beego use LRU algorithm to remove caching file
	// the max memory size of caching static files is StaticCacheFileSize * StaticCacheFileNum
	// see StaticCacheFileSize
	// @Default 1000
	StaticCacheFileNum int
	// TemplateLeft
	// @Description Beego use this to render page
	// see TemplateRight
	// @Default {{
	TemplateLeft string
	// TemplateRight
	// @Description Beego use this to render page
	// see TemplateLeft
	// @Default }}
	TemplateRight string
	// ViewsPath
	// @Description The directory of Beego application storing template
	// @Default views
	ViewsPath string
	// CommentRouterPath
	// @Description Beego scans this directory and its sub directory to generate router
	// Beego only scans this directory when it's in dev environment
	// @Default controllers
	CommentRouterPath string
	// XSRFKey
	// @Description the name of cookie storing XSRF token
	// see EnableXSRF
	// @Default beegoxsrf
	XSRFKey string
	// XSRFExpire
	// @Description the expiration time of XSRF token cookie
	// second
	// @Default 0
	XSRFExpire int
	// @Description session related config
	Session SessionConfig
}

// SessionConfig holds session related config
type SessionConfig struct {
	// SessionOn
	// @Description if it's true, Beego will auto manage session
	// @Default false
	SessionOn bool
	// SessionAutoSetCookie
	// @Description if it's true, Beego will put the session token into cookie too
	// @Default true
	SessionAutoSetCookie bool
	// SessionDisableHTTPOnly
	// @Description used to allow for cross domain cookies/javascript cookies
	// In general, you should not set it to true unless you understand the risk
	// @Default false
	SessionDisableHTTPOnly bool
	// SessionEnableSidInHTTPHeader
	// @Description enable store/get the sessionId into/from http headers
	// @Default false
	SessionEnableSidInHTTPHeader bool
	// SessionEnableSidInURLQuery
	// @Description enable get the sessionId from Url Query params
	// @Default false
	SessionEnableSidInURLQuery bool
	// SessionProvider
	// @Description session provider's name.
	// You should confirm that this provider has been register via session.Register method
	// the default value is memory. This is not suitable for distributed system
	// @Default memory
	SessionProvider string
	// SessionName
	// @Description If SessionAutoSetCookie is true, we use this value as the cookie's name
	// @Default beegosessionID
	SessionName string
	// SessionGCMaxLifetime
	// @Description Beego will GC session to clean useless session.
	// unit: second
	// @Default 3600
	SessionGCMaxLifetime int64
	// SessionProviderConfig
	// @Description the config of session provider
	// see SessionProvider
	// you should read the document of session provider to learn how to set this value
	// @Default ""
	SessionProviderConfig string
	// SessionCookieLifeTime
	// @Description If SessionAutoSetCookie is true,
	// we use this value as the expiration time and max age of the cookie
	// unit second
	// @Default 0
	SessionCookieLifeTime int
	// SessionDomain
	// @Description If SessionAutoSetCookie is true, we use this value as the cookie's domain
	// @Default ""
	SessionDomain string
	// SessionNameInHTTPHeader
	// @Description if SessionEnableSidInHTTPHeader is true, this value will be used as the http header
	// @Default Beegosessionid
	SessionNameInHTTPHeader string
	// SessionCookieSameSite
	// @Description If SessionAutoSetCookie is true, we use this value as the cookie's same site policy
	// the default value is http.SameSiteDefaultMode
	// @Default 1
	SessionCookieSameSite http.SameSite

	// SessionIDPrefix
	// @Description session id's prefix
	// @Default ""
	SessionIDPrefix string
}

// LogConfig holds Log related config
type LogConfig struct {
	// AccessLogs
	// @Description If it's true, Beego will log the HTTP request info
	// @Default false
	AccessLogs bool
	// EnableStaticLogs
	// @Description log static files requests
	// @Default false
	EnableStaticLogs bool
	// FileLineNum
	// @Description if it's true, it will log the line number
	// @Default true
	FileLineNum bool
	// AccessLogsFormat
	// @Description access log format: JSON_FORMAT, APACHE_FORMAT or empty string
	// @Default APACHE_FORMAT
	AccessLogsFormat string
	// Outputs
	// @Description the destination of access log
	// the key is log adapter and the value is adapter's configure
	// @Default "console" => ""
	Outputs map[string]string // Store Adaptor : config
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
	// WorkPath is the absolute path to project root directory
	WorkPath string
)

func init() {
	BConfig = newBConfig()
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	WorkPath, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	filename := "app.conf"
	if os.Getenv("BEEGO_RUNMODE") != "" {
		filename = os.Getenv("BEEGO_RUNMODE") + ".app.conf"
	}
	appConfigPath = filepath.Join(WorkPath, "conf", filename)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "conf", filename)
		if !utils.FileExists(appConfigPath) {
			AppConfig = &beegoAppConfig{innerConfig: config.NewFakeConfig()}
			return
		}
	}
	if err = parseConfig(appConfigPath); err != nil {
		panic(err)
	}
}

func defaultRecoverPanic(ctx *context.Context, cfg *Config) {
	if err := recover(); err != nil {
		if err == ErrAbort {
			return
		}
		if !cfg.RecoverPanic {
			panic(err)
		}
		if cfg.EnableErrorsShow {
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

		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}

		if cfg.RunMode == DEV && cfg.EnableErrorsRender {
			showErr(err, ctx, stack)
		}
	}
}

func newBConfig() *Config {
	res := &Config{
		AppName:             "beego",
		RunMode:             PROD,
		RouterCaseSensitive: true,
		ServerName:          "beegoServer:" + beego.VERSION,
		RecoverPanic:        true,

		CopyRequestBody:    false,
		EnableGzip:         false,
		MaxMemory:          1 << 26, // 64MB
		MaxUploadSize:      1 << 30, // 1GB
		EnableErrorsShow:   true,
		EnableErrorsRender: true,
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
			ClientAuth:    int(tls.RequireAndVerifyClientCert),
		},
		WebConfig: WebConfig{
			AutoRender:             true,
			EnableDocs:             false,
			FlashName:              "BEEGO_FLASH",
			FlashSeparator:         "BEEGOFLASH",
			DirectoryIndex:         false,
			StaticDir:              map[string]string{"/static": "static"},
			StaticExtensionsToGzip: []string{".css", ".js"},
			StaticCacheFileSize:    1024 * 100,
			StaticCacheFileNum:     1000,
			TemplateLeft:           "{{",
			TemplateRight:          "}}",
			ViewsPath:              "views",
			CommentRouterPath:      "controllers",
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
				SessionCookieLifeTime:        0, // set cookie default is the browser life
				SessionAutoSetCookie:         true,
				SessionDomain:                "",
				SessionEnableSidInHTTPHeader: false, // enable store/get the sessionId into/from http headers
				SessionNameInHTTPHeader:      "Beegosessionid",
				SessionEnableSidInURLQuery:   false, // enable get the sessionId from Url Query params
				SessionCookieSameSite:        http.SameSiteDefaultMode,
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

	res.RecoverFunc = defaultRecoverPanic
	return res
}

// now only support ini, next will support json.
func parseConfig(appConfigPath string) (err error) {
	AppConfig, err = newAppConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	return assignConfig(AppConfig)
}

// assignConfig is tricky.
// For 1.x, it use assignSingleConfig to parse the file
// but for 2.x, we use Unmarshaler method
func assignConfig(ac config.Configer) error {
	parseConfigForV1(ac)

	err := ac.Unmarshaler("", BConfig)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("Unmarshaler config file to BConfig failed. "+
			"And if you are working on v1.x config file, please ignore this, err: %s", err))
		return err
	}

	// init log
	logs.Reset()
	for adaptor, cfg := range BConfig.Log.Outputs {
		err := logs.SetLogger(adaptor, cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s with the config %q got err:%s", adaptor, cfg, err.Error()))
			return err
		}
	}
	logs.SetLogFuncCall(BConfig.Log.FileLineNum)
	return nil
}

func parseConfigForV1(ac config.Configer) {
	for _, i := range []interface{}{BConfig, &BConfig.Listen, &BConfig.WebConfig, &BConfig.Log, &BConfig.WebConfig.Session} {
		assignSingleConfig(i, ac)
	}

	// set the run mode first
	if envRunMode := os.Getenv("BEEGO_RUNMODE"); envRunMode != "" {
		BConfig.RunMode = envRunMode
	} else if runMode, err := ac.String("RunMode"); runMode != "" && err == nil {
		BConfig.RunMode = runMode
	}

	if sd, err := ac.String("StaticDir"); sd != "" && err == nil {
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

	if sgz, err := ac.String("StaticExtensionsToGzip"); sgz != "" && err == nil {
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

	if sfs, err := ac.Int("StaticCacheFileSize"); err == nil {
		BConfig.WebConfig.StaticCacheFileSize = sfs
	}

	if sfn, err := ac.Int("StaticCacheFileNum"); err == nil {
		BConfig.WebConfig.StaticCacheFileNum = sfn
	}

	if lo, err := ac.String("LogOutputs"); lo != "" && err == nil {
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
			pf.SetInt(ac.DefaultInt64(name, pf.Int()))
		case reflect.Bool:
			pf.SetBool(ac.DefaultBool(name, pf.Bool()))
		case reflect.Struct:
		default:
			// do nothing here
		}
	}
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

	appConfigPath = absConfigPath
	appConfigProvider = adapterName

	return parseConfig(appConfigPath)
}

type beegoAppConfig struct {
	config.BaseConfiger
	innerConfig config.Configer
}

func newAppConfig(appConfigProvider, appConfigPath string) (*beegoAppConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &beegoAppConfig{innerConfig: ac}, nil
}

func (b *beegoAppConfig) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	return b.innerConfig.Unmarshaler(prefix, obj, opt...)
}

func (b *beegoAppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(BConfig.RunMode+"::"+key, val); err != nil {
		return b.innerConfig.Set(key, val)
	}
	return nil
}

func (b *beegoAppConfig) String(key string) (string, error) {
	if v, err := b.innerConfig.String(BConfig.RunMode + "::" + key); v != "" && err == nil {
		return v, nil
	}
	return b.innerConfig.String(key)
}

func (b *beegoAppConfig) Strings(key string) ([]string, error) {
	if v, err := b.innerConfig.Strings(BConfig.RunMode + "::" + key); len(v) > 0 && err == nil {
		return v, nil
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
	if v, err := b.String(key); v != "" && err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v, err := b.Strings(key); len(v) != 0 && err == nil {
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
