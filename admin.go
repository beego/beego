package beego

import (
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego/toolbox"
	"github.com/astaxie/beego/utils"
)

// BeeAdminApp is the default AdminApp used by admin module.
var BeeAdminApp *AdminApp

// FilterMonitorFunc is default monitor filter when admin module is enable.
// if this func returns, admin module records qbs for this request by condition of this function logic.
// usage:
// 	func MyFilterMonitor(method, requestPath string, t time.Duration) bool {
//	 	if method == "POST" {
//			return false
//	 	}
//	 	if t.Nanoseconds() < 100 {
//			return false
//	 	}
//	 	if strings.HasPrefix(requestPath, "/astaxie") {
//			return false
//	 	}
//	 	return true
// 	}
// 	beego.FilterMonitorFunc = MyFilterMonitor.
var FilterMonitorFunc func(string, string, time.Duration) bool

func init() {
	BeeAdminApp = &AdminApp{
		routers: make(map[string]http.HandlerFunc),
	}
	BeeAdminApp.Route("/", AdminIndex)
	BeeAdminApp.Route("/qps", QpsIndex)
	BeeAdminApp.Route("/prof", ProfIndex)
	BeeAdminApp.Route("/healthcheck", Healthcheck)
	BeeAdminApp.Route("/task", TaskStatus)
	BeeAdminApp.Route("/runtask", RunTask)
	BeeAdminApp.Route("/listconf", ListConf)
	FilterMonitorFunc = func(string, string, time.Duration) bool { return true }
}

// AdminIndex is the default http.Handler for admin module.
// it matches url pattern "/".
func AdminIndex(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Welcome to Admin Dashboard\n"))
	rw.Write([]byte("There are servral functions:\n"))
	rw.Write([]byte("1. Record all request and request time, http://localhost:8088/qps\n"))
	rw.Write([]byte("2. Get runtime profiling data by the pprof, http://localhost:8088/prof\n"))
	rw.Write([]byte("3. Get healthcheck result from http://localhost:8088/healthcheck\n"))
	rw.Write([]byte("4. Get current task infomation from taskhttp://localhost:8088/task \n"))
	rw.Write([]byte("5. To run a task passed a param http://localhost:8088/runtask\n"))
	rw.Write([]byte("6. Get all confige & router infomation http://localhost:8088/listconf\n"))

}

// QpsIndex is the http.Handler for writing qbs statistics map result info in http.ResponseWriter.
// it's registered with url pattern "/qbs" in admin module.
func QpsIndex(rw http.ResponseWriter, r *http.Request) {
	toolbox.StatisticsMap.GetMap(rw)
}

// ListConf is the http.Handler of displaying all beego configuration values as key/value pair.
// it's registered with url pattern "/listconf" in admin module.
func ListConf(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")
	if command != "" {
		switch command {
		case "conf":
			fmt.Fprintln(rw, "list all beego's conf:")
			fmt.Fprintln(rw, "AppName:", AppName)
			fmt.Fprintln(rw, "AppPath:", AppPath)
			fmt.Fprintln(rw, "AppConfigPath:", AppConfigPath)
			fmt.Fprintln(rw, "StaticDir:", StaticDir)
			fmt.Fprintln(rw, "StaticExtensionsToGzip:", StaticExtensionsToGzip)
			fmt.Fprintln(rw, "HttpAddr:", HttpAddr)
			fmt.Fprintln(rw, "HttpPort:", HttpPort)
			fmt.Fprintln(rw, "HttpTLS:", HttpTLS)
			fmt.Fprintln(rw, "HttpCertFile:", HttpCertFile)
			fmt.Fprintln(rw, "HttpKeyFile:", HttpKeyFile)
			fmt.Fprintln(rw, "RecoverPanic:", RecoverPanic)
			fmt.Fprintln(rw, "AutoRender:", AutoRender)
			fmt.Fprintln(rw, "ViewsPath:", ViewsPath)
			fmt.Fprintln(rw, "RunMode:", RunMode)
			fmt.Fprintln(rw, "SessionOn:", SessionOn)
			fmt.Fprintln(rw, "SessionProvider:", SessionProvider)
			fmt.Fprintln(rw, "SessionName:", SessionName)
			fmt.Fprintln(rw, "SessionGCMaxLifetime:", SessionGCMaxLifetime)
			fmt.Fprintln(rw, "SessionSavePath:", SessionSavePath)
			fmt.Fprintln(rw, "SessionHashFunc:", SessionHashFunc)
			fmt.Fprintln(rw, "SessionHashKey:", SessionHashKey)
			fmt.Fprintln(rw, "SessionCookieLifeTime:", SessionCookieLifeTime)
			fmt.Fprintln(rw, "UseFcgi:", UseFcgi)
			fmt.Fprintln(rw, "MaxMemory:", MaxMemory)
			fmt.Fprintln(rw, "EnableGzip:", EnableGzip)
			fmt.Fprintln(rw, "DirectoryIndex:", DirectoryIndex)
			fmt.Fprintln(rw, "EnableHotUpdate:", EnableHotUpdate)
			fmt.Fprintln(rw, "HttpServerTimeOut:", HttpServerTimeOut)
			fmt.Fprintln(rw, "ErrorsShow:", ErrorsShow)
			fmt.Fprintln(rw, "XSRFKEY:", XSRFKEY)
			fmt.Fprintln(rw, "EnableXSRF:", EnableXSRF)
			fmt.Fprintln(rw, "XSRFExpire:", XSRFExpire)
			fmt.Fprintln(rw, "CopyRequestBody:", CopyRequestBody)
			fmt.Fprintln(rw, "TemplateLeft:", TemplateLeft)
			fmt.Fprintln(rw, "TemplateRight:", TemplateRight)
			fmt.Fprintln(rw, "BeegoServerName:", BeegoServerName)
			fmt.Fprintln(rw, "EnableAdmin:", EnableAdmin)
			fmt.Fprintln(rw, "AdminHttpAddr:", AdminHttpAddr)
			fmt.Fprintln(rw, "AdminHttpPort:", AdminHttpPort)
		case "router":
			fmt.Fprintln(rw, "Print all router infomation:")
			for _, router := range BeeApp.Handlers.fixrouters {
				if router.hasMethod {
					fmt.Fprintln(rw, router.pattern, "----", router.methods, "----", router.controllerType.Name())
				} else {
					fmt.Fprintln(rw, router.pattern, "----", router.controllerType.Name())
				}
			}
			for _, router := range BeeApp.Handlers.routers {
				if router.hasMethod {
					fmt.Fprintln(rw, router.pattern, "----", router.methods, "----", router.controllerType.Name())
				} else {
					fmt.Fprintln(rw, router.pattern, "----", router.controllerType.Name())
				}
			}
			if BeeApp.Handlers.enableAuto {
				for controllerName, methodObj := range BeeApp.Handlers.autoRouter {
					fmt.Fprintln(rw, controllerName, "----")
					for methodName, obj := range methodObj {
						fmt.Fprintln(rw, "        ", methodName, "-----", obj.Name())
					}
				}
			}
		case "filter":
			fmt.Fprintln(rw, "Print all filter infomation:")
			if BeeApp.Handlers.enableFilter {
				fmt.Fprintln(rw, "BeforeRouter:")
				if bf, ok := BeeApp.Handlers.filters[BeforeRouter]; ok {
					for _, f := range bf {
						fmt.Fprintln(rw, f.pattern, utils.GetFuncName(f.filterFunc))
					}
				}
				fmt.Fprintln(rw, "AfterStatic:")
				if bf, ok := BeeApp.Handlers.filters[AfterStatic]; ok {
					for _, f := range bf {
						fmt.Fprintln(rw, f.pattern, utils.GetFuncName(f.filterFunc))
					}
				}
				fmt.Fprintln(rw, "BeforeExec:")
				if bf, ok := BeeApp.Handlers.filters[BeforeExec]; ok {
					for _, f := range bf {
						fmt.Fprintln(rw, f.pattern, utils.GetFuncName(f.filterFunc))
					}
				}
				fmt.Fprintln(rw, "AfterExec:")
				if bf, ok := BeeApp.Handlers.filters[AfterExec]; ok {
					for _, f := range bf {
						fmt.Fprintln(rw, f.pattern, utils.GetFuncName(f.filterFunc))
					}
				}
				fmt.Fprintln(rw, "FinishRouter:")
				if bf, ok := BeeApp.Handlers.filters[FinishRouter]; ok {
					for _, f := range bf {
						fmt.Fprintln(rw, f.pattern, utils.GetFuncName(f.filterFunc))
					}
				}
			}
		default:
			rw.Write([]byte("command not support"))
		}
	} else {
		rw.Write([]byte("ListConf support this command:\n"))
		rw.Write([]byte("1. command=conf\n"))
		rw.Write([]byte("2. command=router\n"))
		rw.Write([]byte("3. command=filter\n"))
	}
}

// ProfIndex is a http.Handler for showing profile command.
// it's in url pattern "/prof" in admin module.
func ProfIndex(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")
	if command != "" {
		toolbox.ProcessInput(command, rw)
	} else {
		rw.Write([]byte("request url like '/prof?command=lookup goroutine'\n"))
		rw.Write([]byte("the command have below types:\n"))
		rw.Write([]byte("1. lookup goroutine\n"))
		rw.Write([]byte("2. lookup heap\n"))
		rw.Write([]byte("3. lookup threadcreate\n"))
		rw.Write([]byte("4. lookup block\n"))
		rw.Write([]byte("5. start cpuprof\n"))
		rw.Write([]byte("6. stop cpuprof\n"))
		rw.Write([]byte("7. get memprof\n"))
		rw.Write([]byte("8. gc summary\n"))
	}
}

// Healthcheck is a http.Handler calling health checking and showing the result.
// it's in "/healthcheck" pattern in admin module.
func Healthcheck(rw http.ResponseWriter, req *http.Request) {
	for name, h := range toolbox.AdminCheckList {
		if err := h.Check(); err != nil {
			fmt.Fprintf(rw, "%s : ok\n", name)
		} else {
			fmt.Fprintf(rw, "%s : %s\n", name, err.Error())
		}
	}
}

// TaskStatus is a http.Handler with running task status (task name, status and the last execution).
// it's in "/task" pattern in admin module.
func TaskStatus(rw http.ResponseWriter, req *http.Request) {
	for tname, tk := range toolbox.AdminTaskList {
		fmt.Fprintf(rw, "%s:%s:%s", tname, tk.GetStatus(), tk.GetPrev().String())
	}
}

// RunTask is a http.Handler to run a Task from the "query string.
// the request url likes /runtask?taskname=sendmail.
func RunTask(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	taskname := req.Form.Get("taskname")
	if t, ok := toolbox.AdminTaskList[taskname]; ok {
		err := t.Run()
		if err != nil {
			fmt.Fprintf(rw, "%v", err)
		}
		fmt.Fprintf(rw, "%s run success,Now the Status is %s", t.GetStatus())
	} else {
		fmt.Fprintf(rw, "there's no task which named:%s", taskname)
	}
}

// AdminApp is an http.HandlerFunc map used as BeeAdminApp.
type AdminApp struct {
	routers map[string]http.HandlerFunc
}

// Route adds http.HandlerFunc to AdminApp with url pattern.
func (admin *AdminApp) Route(pattern string, f http.HandlerFunc) {
	admin.routers[pattern] = f
}

// Run AdminApp http server.
// Its addr is defined in configuration file as adminhttpaddr and adminhttpport.
func (admin *AdminApp) Run() {
	if len(toolbox.AdminTaskList) > 0 {
		toolbox.StartTask()
	}
	addr := AdminHttpAddr

	if AdminHttpPort != 0 {
		addr = fmt.Sprintf("%s:%d", AdminHttpAddr, AdminHttpPort)
	}
	for p, f := range admin.routers {
		http.Handle(p, f)
	}
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		BeeLogger.Critical("Admin ListenAndServe: ", err)
	}
}
