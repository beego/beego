// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package beego

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/toolbox"
	"github.com/astaxie/beego/utils"
)

// BeeAdminApp is the default adminApp used by admin module.
var beeAdminApp *adminApp

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
	beeAdminApp = &adminApp{
		routers: make(map[string]http.HandlerFunc),
	}
	beeAdminApp.Route("/", adminIndex)
	beeAdminApp.Route("/qps", qpsIndex)
	beeAdminApp.Route("/prof", profIndex)
	beeAdminApp.Route("/healthcheck", healthcheck)
	beeAdminApp.Route("/task", taskStatus)
	beeAdminApp.Route("/runtask", runTask)
	beeAdminApp.Route("/listconf", listConf)
	FilterMonitorFunc = func(string, string, time.Duration) bool { return true }
}

// AdminIndex is the default http.Handler for admin module.
// it matches url pattern "/".
func adminIndex(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("<html><head><title>beego admin dashboard</title></head><body>"))
	rw.Write([]byte("Welcome to Admin Dashboard<br>\n"))
	rw.Write([]byte("There are servral functions:<br>\n"))
	rw.Write([]byte("1. Record all request and request time, <a href='/qps'>http://localhost:" + strconv.Itoa(AdminHttpPort) + "/qps</a><br>\n"))
	rw.Write([]byte("2. Get runtime profiling data by the pprof, <a href='/prof'>http://localhost:" + strconv.Itoa(AdminHttpPort) + "/prof</a><br>\n"))
	rw.Write([]byte("3. Get healthcheck result from <a href='/healthcheck'>http://localhost:" + strconv.Itoa(AdminHttpPort) + "/healthcheck</a><br>\n"))
	rw.Write([]byte("4. Get current task infomation from task <a href='/task'>http://localhost:" + strconv.Itoa(AdminHttpPort) + "/task</a><br> \n"))
	rw.Write([]byte("5. To run a task passed a param <a href='/runtask'>http://localhost:" + strconv.Itoa(AdminHttpPort) + "/runtask</a><br>\n"))
	rw.Write([]byte("6. Get all confige & router infomation <a href='/listconf'>http://localhost:" + strconv.Itoa(AdminHttpPort) + "/listconf</a><br>\n"))
	rw.Write([]byte("</body></html>"))
}

// QpsIndex is the http.Handler for writing qbs statistics map result info in http.ResponseWriter.
// it's registered with url pattern "/qbs" in admin module.
func qpsIndex(rw http.ResponseWriter, r *http.Request) {
	toolbox.StatisticsMap.GetMap(rw)
}

// ListConf is the http.Handler of displaying all beego configuration values as key/value pair.
// it's registered with url pattern "/listconf" in admin module.
func listConf(rw http.ResponseWriter, r *http.Request) {
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
			fmt.Fprintln(rw, "HttpTLS:", EnableHttpTLS)
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
			for method, t := range BeeApp.Handlers.routers {
				fmt.Fprintln(rw)
				fmt.Fprintln(rw)
				fmt.Fprintln(rw, "		Method:", method)
				printTree(rw, t)
			}
			// @todo print routers
		case "filter":
			fmt.Fprintln(rw, "Print all filter infomation:")
			if BeeApp.Handlers.enableFilter {
				fmt.Fprintln(rw, "BeforeRouter:")
				if bf, ok := BeeApp.Handlers.filters[BeforeRouter]; ok {
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
		rw.Write([]byte("<html><head><title>beego admin dashboard</title></head><body>"))
		rw.Write([]byte("ListConf support this command:<br>\n"))
		rw.Write([]byte("1. <a href='?command=conf'>command=conf</a><br>\n"))
		rw.Write([]byte("2. <a href='?command=router'>command=router</a><br>\n"))
		rw.Write([]byte("3. <a href='?command=filter'>command=filter</a><br>\n"))
		rw.Write([]byte("</body></html>"))
	}
}

func printTree(rw http.ResponseWriter, t *Tree) {
	for _, tr := range t.fixrouters {
		printTree(rw, tr)
	}
	if t.wildcard != nil {
		printTree(rw, t.wildcard)
	}
	for _, l := range t.leaves {
		if v, ok := l.runObject.(*controllerInfo); ok {
			if v.routerType == routerTypeBeego {
				fmt.Fprintln(rw, v.pattern, v.methods, v.controllerType.Name())
			} else if v.routerType == routerTypeRESTFul {
				fmt.Fprintln(rw, v.pattern, v.methods)
			} else if v.routerType == routerTypeHandler {
				fmt.Fprintln(rw, v.pattern, "handler")
			}
		}
	}
}

// ProfIndex is a http.Handler for showing profile command.
// it's in url pattern "/prof" in admin module.
func profIndex(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")
	if command != "" {
		toolbox.ProcessInput(command, rw)
	} else {
		rw.Write([]byte("<html><head><title>beego admin dashboard</title></head><body>"))
		rw.Write([]byte("request url like '/prof?command=lookup goroutine'<br>\n"))
		rw.Write([]byte("the command have below types:<br>\n"))
		rw.Write([]byte("1. <a href='?command=lookup goroutine'>lookup goroutine</a><br>\n"))
		rw.Write([]byte("2. <a href='?command=lookup heap'>lookup heap</a><br>\n"))
		rw.Write([]byte("3. <a href='?command=lookup threadcreate'>lookup threadcreate</a><br>\n"))
		rw.Write([]byte("4. <a href='?command=lookup block'>lookup block</a><br>\n"))
		rw.Write([]byte("5. <a href='?command=start cpuprof'>start cpuprof</a><br>\n"))
		rw.Write([]byte("6. <a href='?command=stop cpuprof'>stop cpuprof</a><br>\n"))
		rw.Write([]byte("7. <a href='?command=get memprof'>get memprof</a><br>\n"))
		rw.Write([]byte("8. <a href='?command=gc summary'>gc summary</a><br>\n"))
		rw.Write([]byte("</body></html>"))
	}
}

// Healthcheck is a http.Handler calling health checking and showing the result.
// it's in "/healthcheck" pattern in admin module.
func healthcheck(rw http.ResponseWriter, req *http.Request) {
	for name, h := range toolbox.AdminCheckList {
		if err := h.Check(); err != nil {
			fmt.Fprintf(rw, "%s : %s\n", name, err.Error())
		} else {
			fmt.Fprintf(rw, "%s : ok\n", name)
		}
	}
}

// TaskStatus is a http.Handler with running task status (task name, status and the last execution).
// it's in "/task" pattern in admin module.
func taskStatus(rw http.ResponseWriter, req *http.Request) {
	for tname, tk := range toolbox.AdminTaskList {
		fmt.Fprintf(rw, "%s:%s:%s", tname, tk.GetStatus(), tk.GetPrev().String())
	}
}

// RunTask is a http.Handler to run a Task from the "query string.
// the request url likes /runtask?taskname=sendmail.
func runTask(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	taskname := req.Form.Get("taskname")
	if t, ok := toolbox.AdminTaskList[taskname]; ok {
		err := t.Run()
		if err != nil {
			fmt.Fprintf(rw, "%v", err)
		}
		fmt.Fprintf(rw, "%s run success,Now the Status is %s", taskname, t.GetStatus())
	} else {
		fmt.Fprintf(rw, "there's no task which named:%s", taskname)
	}
}

// adminApp is an http.HandlerFunc map used as beeAdminApp.
type adminApp struct {
	routers map[string]http.HandlerFunc
}

// Route adds http.HandlerFunc to adminApp with url pattern.
func (admin *adminApp) Route(pattern string, f http.HandlerFunc) {
	admin.routers[pattern] = f
}

// Run adminApp http server.
// Its addr is defined in configuration file as adminhttpaddr and adminhttpport.
func (admin *adminApp) Run() {
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
