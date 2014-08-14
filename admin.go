// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie
package beego

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
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
	beeAdminApp.Route("/listconf", listConf)
	FilterMonitorFunc = func(string, string, time.Duration) bool { return true }
}

// AdminIndex is the default http.Handler for admin module.
// it matches url pattern "/".
func adminIndex(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
	tmpl = template.Must(tmpl.Parse(indexTpl))
	data := make(map[interface{}]interface{})
	tmpl.Execute(rw, data)
}

// QpsIndex is the http.Handler for writing qbs statistics map result info in http.ResponseWriter.
// it's registered with url pattern "/qbs" in admin module.
func qpsIndex(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
	tmpl = template.Must(tmpl.Parse(qpsTpl))
	data := make(map[interface{}]interface{})
	data["Content"] = toolbox.StatisticsMap.GetMap()

	tmpl.Execute(rw, data)

}

// ListConf is the http.Handler of displaying all beego configuration values as key/value pair.
// it's registered with url pattern "/listconf" in admin module.
func listConf(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")
	if command != "" {
		data := make(map[interface{}]interface{})
		switch command {
		case "conf":
			m := make(map[string]interface{})

			m["AppName"] = AppName
			m["AppPath"] = AppPath
			m["AppConfigPath"] = AppConfigPath
			m["StaticDir"] = StaticDir
			m["StaticExtensionsToGzip"] = StaticExtensionsToGzip
			m["HttpAddr"] = HttpAddr
			m["HttpPort"] = HttpPort
			m["HttpTLS"] = EnableHttpTLS
			m["HttpCertFile"] = HttpCertFile
			m["HttpKeyFile"] = HttpKeyFile
			m["RecoverPanic"] = RecoverPanic
			m["AutoRender"] = AutoRender
			m["ViewsPath"] = ViewsPath
			m["RunMode"] = RunMode
			m["SessionOn"] = SessionOn
			m["SessionProvider"] = SessionProvider
			m["SessionName"] = SessionName
			m["SessionGCMaxLifetime"] = SessionGCMaxLifetime
			m["SessionSavePath"] = SessionSavePath
			m["SessionHashFunc"] = SessionHashFunc
			m["SessionHashKey"] = SessionHashKey
			m["SessionCookieLifeTime"] = SessionCookieLifeTime
			m["UseFcgi"] = UseFcgi
			m["MaxMemory"] = MaxMemory
			m["EnableGzip"] = EnableGzip
			m["DirectoryIndex"] = DirectoryIndex
			m["HttpServerTimeOut"] = HttpServerTimeOut
			m["ErrorsShow"] = ErrorsShow
			m["XSRFKEY"] = XSRFKEY
			m["EnableXSRF"] = EnableXSRF
			m["XSRFExpire"] = XSRFExpire
			m["CopyRequestBody"] = CopyRequestBody
			m["TemplateLeft"] = TemplateLeft
			m["TemplateRight"] = TemplateRight
			m["BeegoServerName"] = BeegoServerName
			m["EnableAdmin"] = EnableAdmin
			m["AdminHttpAddr"] = AdminHttpAddr
			m["AdminHttpPort"] = AdminHttpPort

			tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
			tmpl = template.Must(tmpl.Parse(configTpl))

			data["Content"] = m

			tmpl.Execute(rw, data)

		case "router":
			resultList := new([][]string)

			var result = []string{
				fmt.Sprintf("header"),
				fmt.Sprintf("Router Pattern"),
				fmt.Sprintf("Methods"),
				fmt.Sprintf("Controller"),
			}
			*resultList = append(*resultList, result)

			for method, t := range BeeApp.Handlers.routers {
				var result = []string{
					fmt.Sprintf("success"),
					fmt.Sprintf("Method: %s", method),
					fmt.Sprintf(""),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)

				printTree(resultList, t)
			}
			data["Content"] = resultList
			data["Title"] = "Routers"
			tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
			tmpl = template.Must(tmpl.Parse(routerAndFilterTpl))
			tmpl.Execute(rw, data)
		case "filter":
			resultList := new([][]string)

			var result = []string{
				fmt.Sprintf("header"),
				fmt.Sprintf("Router Pattern"),
				fmt.Sprintf("Filter Function"),
			}
			*resultList = append(*resultList, result)

			if BeeApp.Handlers.enableFilter {
				var result = []string{
					fmt.Sprintf("success"),
					fmt.Sprintf("Before Router"),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)

				if bf, ok := BeeApp.Handlers.filters[BeforeRouter]; ok {
					for _, f := range bf {

						var result = []string{
							fmt.Sprintf(""),
							fmt.Sprintf("%s", f.pattern),
							fmt.Sprintf("%s", utils.GetFuncName(f.filterFunc)),
						}
						*resultList = append(*resultList, result)

					}
				}
				result = []string{
					fmt.Sprintf("success"),
					fmt.Sprintf("Before Exec"),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)
				if bf, ok := BeeApp.Handlers.filters[BeforeExec]; ok {
					for _, f := range bf {

						var result = []string{
							fmt.Sprintf(""),
							fmt.Sprintf("%s", f.pattern),
							fmt.Sprintf("%s", utils.GetFuncName(f.filterFunc)),
						}
						*resultList = append(*resultList, result)

					}
				}
				result = []string{
					fmt.Sprintf("success"),
					fmt.Sprintf("AfterExec Exec"),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)

				if bf, ok := BeeApp.Handlers.filters[AfterExec]; ok {
					for _, f := range bf {

						var result = []string{
							fmt.Sprintf(""),
							fmt.Sprintf("%s", f.pattern),
							fmt.Sprintf("%s", utils.GetFuncName(f.filterFunc)),
						}
						*resultList = append(*resultList, result)

					}
				}
				result = []string{
					fmt.Sprintf("success"),
					fmt.Sprintf("Finish Router"),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)

				if bf, ok := BeeApp.Handlers.filters[FinishRouter]; ok {
					for _, f := range bf {

						var result = []string{
							fmt.Sprintf(""),
							fmt.Sprintf("%s", f.pattern),
							fmt.Sprintf("%s", utils.GetFuncName(f.filterFunc)),
						}
						*resultList = append(*resultList, result)

					}
				}
			}
			data["Content"] = resultList
			data["Title"] = "Filters"
			tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
			tmpl = template.Must(tmpl.Parse(routerAndFilterTpl))
			tmpl.Execute(rw, data)

		default:
			rw.Write([]byte("command not support"))
		}
	} else {
	}
}

func printTree(resultList *[][]string, t *Tree) {
	for _, tr := range t.fixrouters {
		printTree(resultList, tr)
	}
	if t.wildcard != nil {
		printTree(resultList, t.wildcard)
	}
	for _, l := range t.leaves {
		if v, ok := l.runObject.(*controllerInfo); ok {
			if v.routerType == routerTypeBeego {
				var result = []string{
					fmt.Sprintf(""),
					fmt.Sprintf("%s", v.pattern),
					fmt.Sprintf("%s", v.methods),
					fmt.Sprintf("%s", v.controllerType),
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeRESTFul {
				var result = []string{
					fmt.Sprintf(""),
					fmt.Sprintf("%s", v.pattern),
					fmt.Sprintf("%s", v.methods),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeHandler {
				var result = []string{
					fmt.Sprintf(""),
					fmt.Sprintf("%s", v.pattern),
					fmt.Sprintf(""),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)
			}
		}
	}
}

// ProfIndex is a http.Handler for showing profile command.
// it's in url pattern "/prof" in admin module.
func profIndex(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")
	data := make(map[interface{}]interface{})

	var result bytes.Buffer
	if command != "" {
		toolbox.ProcessInput(command, &result)
		data["Content"] = result.String()
		data["Title"] = command

		tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
		tmpl = template.Must(tmpl.Parse(profillingTpl))
		tmpl.Execute(rw, data)
	} else {
	}
}

// Healthcheck is a http.Handler calling health checking and showing the result.
// it's in "/healthcheck" pattern in admin module.
func healthcheck(rw http.ResponseWriter, req *http.Request) {
	data := make(map[interface{}]interface{})

	resultList := new([][]string)
	var result = []string{
		fmt.Sprintf("header"),
		fmt.Sprintf("Name"),
		fmt.Sprintf("Status"),
	}
	*resultList = append(*resultList, result)

	for name, h := range toolbox.AdminCheckList {
		if err := h.Check(); err != nil {
			result = []string{
				fmt.Sprintf("error"),
				fmt.Sprintf("%s", name),
				fmt.Sprintf("%s", err.Error()),
			}

		} else {
			result = []string{
				fmt.Sprintf("success"),
				fmt.Sprintf("%s", name),
				fmt.Sprintf("OK"),
			}

		}
		*resultList = append(*resultList, result)
	}

	data["Content"] = resultList
	data["Title"] = "Health Check"
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
	tmpl = template.Must(tmpl.Parse(healthCheckTpl))
	tmpl.Execute(rw, data)

}

// TaskStatus is a http.Handler with running task status (task name, status and the last execution).
// it's in "/task" pattern in admin module.
func taskStatus(rw http.ResponseWriter, req *http.Request) {
	data := make(map[interface{}]interface{})

	// Run Task
	req.ParseForm()
	taskname := req.Form.Get("taskname")
	if taskname != "" {

		if t, ok := toolbox.AdminTaskList[taskname]; ok {
			err := t.Run()
			if err != nil {
				data["Message"] = []string{"error", fmt.Sprintf("%s", err)}
			}
			data["Message"] = []string{"success", fmt.Sprintf("%s run success,Now the Status is %s", taskname, t.GetStatus())}
		} else {
			data["Message"] = []string{"warning", fmt.Sprintf("there's no task which named: %s", taskname)}
		}
	}

	// List Tasks
	resultList := new([][]string)
	var result = []string{
		fmt.Sprintf("header"),
		fmt.Sprintf("Task Name"),
		fmt.Sprintf("Task Spec"),
		fmt.Sprintf("Task Function"),
	}
	*resultList = append(*resultList, result)
	for tname, tk := range toolbox.AdminTaskList {
		result = []string{
			fmt.Sprintf(""),
			fmt.Sprintf("%s", tname),
			fmt.Sprintf("%s", tk.GetStatus()),
			fmt.Sprintf("%s", tk.GetPrev().String()),
		}
		*resultList = append(*resultList, result)
	}

	data["Content"] = resultList
	data["Title"] = "Tasks"
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
	tmpl = template.Must(tmpl.Parse(tasksTpl))
	tmpl.Execute(rw, data)
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
