package beego

import (
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
)

var (
	// beego application
	BeeApp                 *App
	// application configurations
	AppName                string
	AppPath                string
	AppConfigPath          string
	StaticDir              map[string]string
	// template caching map
	TemplateCache          map[string]*template.Template
	// files with should be compressed with gzip (.js,.css,etc)
	StaticExtensionsToGzip []string
	// http server configurations
	HttpAddr               string
	HttpPort               int
	HttpTLS                bool
	HttpCertFile           string
	HttpKeyFile            string
	// flag of auto recover panic
	RecoverPanic           bool
	// flag of render template automatically
	AutoRender             bool
	ViewsPath              string
	// run mode, "dev" or "prod"
	RunMode                string
	AppConfig              config.ConfigContainer
	// global session mananger
	GlobalSessions        *session.Manager
	// flag of starting session auto. default is false.
	SessionOn             bool
	// default session provider, memory, mysql , redis ,etc.
	SessionProvider       string
	// the cookie name when saving session id into cookie.
	SessionName           string
	// session gc time for auto cleaning expired session.
	SessionGCMaxLifetime  int64
	// if use mysql/redis/file provider, define save path to connection info.
	SessionSavePath       string
	// session hash generation func.
	SessionHashFunc       string
	// session hash salt string.
	SessionHashKey        string
	// the life time of session id in cookie.
	SessionCookieLifeTime int
	UseFcgi               bool
	MaxMemory             int64
	// flag of enable gzip
	EnableGzip            bool
	// flag of display directory index. default is false.
	DirectoryIndex        bool
	// flag of hot update checking in app self. default is false.
	EnableHotUpdate       bool
	HttpServerTimeOut     int64
	// flag of show errors in page. if true, show error and trace info in page rendered with error template.
	ErrorsShow            bool
	// xsrf hash salt string.
	XSRFKEY               string
	// flag of enable xsrf.
	EnableXSRF            bool
	// the expiry of xsrf value.
	XSRFExpire            int
	// flag of copy raw request body in context.
	CopyRequestBody       bool
	TemplateLeft          string
	TemplateRight         string
	// beego server name exported in response header.
	BeegoServerName       string
	// flag of enable admin module to log every request info.
	EnableAdmin           bool
	// http server configurations for admin module.
	AdminHttpAddr         string
	AdminHttpPort         int
)

func init() {
	// create beego application
	BeeApp = NewApp()

	// initialize default configurations
	AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(AppPath)

	StaticDir = make(map[string]string)
	StaticDir["/static"] = "static"

	StaticExtensionsToGzip = []string{".css", ".js"}

	TemplateCache = make(map[string]*template.Template)

	// set this to 0.0.0.0 to make this app available to externally
	HttpAddr = ""
	HttpPort = 8080

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

	UseFcgi = false

	MaxMemory = 1<<26 //64MB

	EnableGzip = false

	AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")

	HttpServerTimeOut = 0

	ErrorsShow = true

	XSRFKEY = "beegoxsrf"
	XSRFExpire = 0

	TemplateLeft = "{{"
	TemplateRight = "}}"

	BeegoServerName = "beegoServer"

	EnableAdmin = false
	AdminHttpAddr = "127.0.0.1"
	AdminHttpPort = 8088

	runtime.GOMAXPROCS(runtime.NumCPU())

	// init BeeLogger
	BeeLogger = logs.NewLogger(10000)
	BeeLogger.SetLogger("console", "")

	err := ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		Info(err)
	}
}

// ParseConfig parsed default config file.
// now only support ini, next will support json.
func ParseConfig() (err error) {
	AppConfig, err = config.NewConfig("ini", AppConfigPath)
	if err != nil {
		return err
	} else {
		HttpAddr = AppConfig.String("HttpAddr")

		if v, err := AppConfig.Int("HttpPort"); err == nil {
			HttpPort = v
		}

		if maxmemory, err := AppConfig.Int64("MaxMemory"); err == nil {
			MaxMemory = maxmemory
		}

		if appname := AppConfig.String("AppName"); appname != "" {
			AppName = appname
		}

		if runmode := AppConfig.String("RunMode"); runmode != "" {
			RunMode = runmode
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

		if sesshashfunc := AppConfig.String("SessionHashFunc"); sesshashfunc != "" {
			SessionHashFunc = sesshashfunc
		}

		if sesshashkey := AppConfig.String("SessionHashKey"); sesshashkey != "" {
			SessionHashKey = sesshashkey
		}

		if sessMaxLifeTime, err := AppConfig.Int("SessionGCMaxLifetime"); err == nil && sessMaxLifeTime != 0 {
			int64val, _ := strconv.ParseInt(strconv.Itoa(sessMaxLifeTime), 10, 64)
			SessionGCMaxLifetime = int64val
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

		if hotupdate, err := AppConfig.Bool("HotUpdate"); err == nil {
			EnableHotUpdate = hotupdate
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

		if httptls, err := AppConfig.Bool("HttpTLS"); err == nil {
			HttpTLS = httptls
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

		if sd := AppConfig.String("StaticDir"); sd != "" {
			for k := range StaticDir {
				delete(StaticDir, k)
			}
			sds := strings.Fields(sd)
			for _, v := range sds {
				if url2fsmap := strings.SplitN(v, ":", 2); len(url2fsmap) == 2 {
					StaticDir["/"+url2fsmap[0]] = url2fsmap[1]
				} else {
					StaticDir["/"+url2fsmap[0]] = url2fsmap[0]
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
	}
	return nil
}
