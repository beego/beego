package beego

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/session"
	"html/template"
	"os"
	"path"
	"runtime"
	"strconv"
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
	PprofOn       bool
	ViewsPath     string
	RunMode       string //"dev" or "prod"
	AppConfig     config.ConfigContainer
	//related to session
	GlobalSessions       *session.Manager //GlobalSessions
	SessionOn            bool             // whether auto start session,default is false
	SessionProvider      string           // default session provider  memory mysql redis
	SessionName          string           // sessionName cookie's name
	SessionGCMaxLifetime int64            // session's gc maxlifetime
	SessionSavePath      string           // session savepath if use mysql/redis/file this set to the connectinfo
	UseFcgi              bool
	MaxMemory            int64
	EnableGzip           bool   // enable gzip
	DirectoryIndex       bool   //enable DirectoryIndex default is false
	EnableHotUpdate      bool   //enable HotUpdate default is false
	HttpServerTimeOut    int64  //set httpserver timeout
	ErrorsShow           bool   //set weather show errors
	XSRFKEY              string //set XSRF
	EnableXSRF           bool
	XSRFExpire           int
	CopyRequestBody      bool //When in raw application, You want to the reqeustbody
	TemplateLeft         string
	TemplateRight        string
)

func init() {
	os.Chdir(path.Dir(os.Args[0]))
	BeeApp = NewApp()
	AppPath, _ = path.Dir(os.Args[0])
	StaticDir = make(map[string]string)
	TemplateCache = make(map[string]*template.Template)
	HttpAddr = ""
	HttpPort = 8080
	AppName = "beego"
	RunMode = "dev" //default runmod
	AutoRender = true
	RecoverPanic = true
	PprofOn = false
	ViewsPath = "views"
	SessionOn = false
	SessionProvider = "memory"
	SessionName = "beegosessionID"
	SessionGCMaxLifetime = 3600
	SessionSavePath = ""
	UseFcgi = false
	MaxMemory = 1 << 26 //64MB
	EnableGzip = false
	StaticDir["/static"] = "static"
	AppConfigPath = path.Join(AppPath, "conf", "app.conf")
	HttpServerTimeOut = 0
	ErrorsShow = true
	XSRFKEY = "beegoxsrf"
	XSRFExpire = 0
	TemplateLeft = "{{"
	TemplateRight = "}}"
	ParseConfig()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func ParseConfig() (err error) {
	AppConfig, err = config.NewConfig("ini", AppConfigPath)
	if err != nil {
		return err
	} else {
		HttpAddr = AppConfig.String("httpaddr")
		if v, err := AppConfig.Int("httpport"); err == nil {
			HttpPort = v
		}
		if maxmemory, err := AppConfig.Int64("maxmemory"); err == nil {
			MaxMemory = maxmemory
		}
		AppName = AppConfig.String("appname")
		if runmode := AppConfig.String("runmode"); runmode != "" {
			RunMode = runmode
		}
		if autorender, err := AppConfig.Bool("autorender"); err == nil {
			AutoRender = autorender
		}
		if autorecover, err := AppConfig.Bool("autorecover"); err == nil {
			RecoverPanic = autorecover
		}
		if pprofon, err := AppConfig.Bool("pprofon"); err == nil {
			PprofOn = pprofon
		}
		if views := AppConfig.String("viewspath"); views != "" {
			ViewsPath = views
		}
		if sessionon, err := AppConfig.Bool("sessionon"); err == nil {
			SessionOn = sessionon
		}
		if sessProvider := AppConfig.String("sessionprovider"); sessProvider != "" {
			SessionProvider = sessProvider
		}
		if sessName := AppConfig.String("sessionname"); sessName != "" {
			SessionName = sessName
		}
		if sesssavepath := AppConfig.String("sessionsavepath"); sesssavepath != "" {
			SessionSavePath = sesssavepath
		}
		if sessMaxLifeTime, err := AppConfig.Int("sessiongcmaxlifetime"); err == nil && sessMaxLifeTime != 0 {
			int64val, _ := strconv.ParseInt(strconv.Itoa(sessMaxLifeTime), 10, 64)
			SessionGCMaxLifetime = int64val
		}
		if usefcgi, err := AppConfig.Bool("usefcgi"); err == nil {
			UseFcgi = usefcgi
		}
		if enablegzip, err := AppConfig.Bool("enablegzip"); err == nil {
			EnableGzip = enablegzip
		}
		if directoryindex, err := AppConfig.Bool("directoryindex"); err == nil {
			DirectoryIndex = directoryindex
		}
		if hotupdate, err := AppConfig.Bool("hotupdate"); err == nil {
			EnableHotUpdate = hotupdate
		}
		if timeout, err := AppConfig.Int64("httpservertimeout"); err == nil {
			HttpServerTimeOut = timeout
		}
		if errorsshow, err := AppConfig.Bool("errorsshow"); err == nil {
			ErrorsShow = errorsshow
		}
		if copyrequestbody, err := AppConfig.Bool("copyrequestbody"); err == nil {
			CopyRequestBody = copyrequestbody
		}
		if xsrfkey := AppConfig.String("xsrfkey"); xsrfkey != "" {
			XSRFKEY = xsrfkey
		}
		if enablexsrf, err := AppConfig.Bool("enablexsrf"); err == nil {
			EnableXSRF = enablexsrf
		}
		if expire, err := AppConfig.Int("xsrfexpire"); err == nil {
			XSRFExpire = expire
		}
		if tplleft := AppConfig.String("templateleft"); tplleft != "" {
			TemplateLeft = tplleft
		}
		if tplright := AppConfig.String("templateright"); tplright != "" {
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
	}
	return nil
}
