package beego

import (
	"html/template"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/session"
)

var (
	BeeApp        *App
	AppName       string
	AppPath       string
	AppConfigPath string
	StaticDir     map[string]string
	TemplateCache map[string]*template.Template
	HttpAddr      string
	HttpPort      int
	HttpTLS       bool
	HttpCertFile  string
	HttpKeyFile   string
	RecoverPanic  bool
	AutoRender    bool
	ViewsPath     string
	RunMode       string //"dev" or "prod"
	AppConfig     config.ConfigContainer
	//related to session
	GlobalSessions        *session.Manager //GlobalSessions
	SessionOn             bool             // whether auto start session,default is false
	SessionProvider       string           // default session provider  memory mysql redis
	SessionName           string           // sessionName cookie's name
	SessionGCMaxLifetime  int64            // session's gc maxlifetime
	SessionSavePath       string           // session savepath if use mysql/redis/file this set to the connectinfo
	SessionHashFunc       string
	SessionHashKey        string
	SessionCookieLifeTime int
	UseFcgi               bool
	MaxMemory             int64
	EnableGzip            bool   // enable gzip
	DirectoryIndex        bool   //enable DirectoryIndex default is false
	EnableHotUpdate       bool   //enable HotUpdate default is false
	HttpServerTimeOut     int64  //set httpserver timeout
	ErrorsShow            bool   //set weather show errors
	XSRFKEY               string //set XSRF
	EnableXSRF            bool
	XSRFExpire            int
	CopyRequestBody       bool //When in raw application, You want to the reqeustbody
	TemplateLeft          string
	TemplateRight         string
	BeegoServerName       string
	EnableAdmin           bool   //enable admin module to log api time
	AdminHttpAddr         string //admin module http addr
	AdminHttpPort         int
)

func init() {
	// create beeapp
	BeeApp = NewApp()

	// initialize default configurations
	os.Chdir(path.Dir(os.Args[0]))
	AppPath = path.Dir(os.Args[0])

	StaticDir = make(map[string]string)
	StaticDir["/static"] = "static"

	TemplateCache = make(map[string]*template.Template)

	// set this to 0.0.0.0 to make this app available to externally
	HttpAddr = "127.0.0.1"
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
	SessionCookieLifeTime = 3600

	UseFcgi = false

	MaxMemory = 1 << 26 //64MB

	EnableGzip = false

	AppConfigPath = path.Join(AppPath, "conf", "app.conf")

	HttpServerTimeOut = 0

	ErrorsShow = true

	XSRFKEY = "beegoxsrf"
	XSRFExpire = 0

	TemplateLeft = "{{"
	TemplateRight = "}}"

	BeegoServerName = "beegoServer"

	EnableAdmin = true
	AdminHttpAddr = "127.0.0.1"
	AdminHttpPort = 8088

	runtime.GOMAXPROCS(runtime.NumCPU())

	err := ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// panic unless the err is can not find default configuration file
		panic(err)
	}
}

//parse config now only support ini, next will support json
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
