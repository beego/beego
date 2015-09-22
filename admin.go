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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/astaxie/beego/grace"
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
	execTpl(rw, map[interface{}]interface{}{}, indexTpl, defaultScriptsTpl)
}

// QpsIndex is the http.Handler for writing qbs statistics map result info in http.ResponseWriter.
// it's registered with url pattern "/qbs" in admin module.
func qpsIndex(rw http.ResponseWriter, r *http.Request) {
	data := make(map[interface{}]interface{})
	data["Content"] = toolbox.StatisticsMap.GetMap()
	execTpl(rw, data, qpsTpl, defaultScriptsTpl)
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
			m["HTTPAddr"] = HTTPAddr
			m["HTTPPort"] = HTTPPort
			m["HTTPTLS"] = EnableHTTPTLS
			m["HTTPCertFile"] = HTTPCertFile
			m["HTTPKeyFile"] = HTTPKeyFile
			m["RecoverPanic"] = RecoverPanic
			m["AutoRender"] = AutoRender
			m["ViewsPath"] = ViewsPath
			m["RunMode"] = RunMode
			m["SessionOn"] = SessionOn
			m["SessionProvider"] = SessionProvider
			m["SessionName"] = SessionName
			m["SessionGCMaxLifetime"] = SessionGCMaxLifetime
			m["SessionProviderConfig"] = SessionProviderConfig
			m["SessionCookieLifeTime"] = SessionCookieLifeTime
			m["EnabelFcgi"] = EnabelFcgi
			m["MaxMemory"] = MaxMemory
			m["EnableGzip"] = EnableGzip
			m["DirectoryIndex"] = DirectoryIndex
			m["HTTPServerTimeOut"] = HTTPServerTimeOut
			m["EnableErrorsShow"] = EnableErrorsShow
			m["XSRFKEY"] = XSRFKEY
			m["EnableXSRF"] = EnableXSRF
			m["XSRFExpire"] = XSRFExpire
			m["CopyRequestBody"] = CopyRequestBody
			m["TemplateLeft"] = TemplateLeft
			m["TemplateRight"] = TemplateRight
			m["BeegoServerName"] = BeegoServerName
			m["EnableAdmin"] = EnableAdmin
			m["AdminHTTPAddr"] = AdminHTTPAddr
			m["AdminHTTPPort"] = AdminHTTPPort

			tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
			tmpl = template.Must(tmpl.Parse(configTpl))
			tmpl = template.Must(tmpl.Parse(defaultScriptsTpl))

			data["Content"] = m

			tmpl.Execute(rw, data)

		case "router":
			content := make(map[string]interface{})

			var fields = []string{
				fmt.Sprintf("Router Pattern"),
				fmt.Sprintf("Methods"),
				fmt.Sprintf("Controller"),
			}
			content["Fields"] = fields

			methods := []string{}
			methodsData := make(map[string]interface{})
			for method, t := range BeeApp.Handlers.routers {

				resultList := new([][]string)

				printTree(resultList, t)

				methods = append(methods, method)
				methodsData[method] = resultList
			}

			content["Data"] = methodsData
			content["Methods"] = methods
			data["Content"] = content
			data["Title"] = "Routers"
			execTpl(rw, data, routerAndFilterTpl, defaultScriptsTpl)
		case "filter":
			content := make(map[string]interface{})

			var fields = []string{
				fmt.Sprintf("Router Pattern"),
				fmt.Sprintf("Filter Function"),
			}
			content["Fields"] = fields

			filterTypes := []string{}
			filterTypeData := make(map[string]interface{})

			if BeeApp.Handlers.enableFilter {
				var filterType string
				for k, fr := range map[int]string{
					BeforeStatic: "Before Static",
					BeforeRouter: "Before Router",
					BeforeExec:   "Before Exec",
					AfterExec:    "After Exec",
					FinishRouter: "Finish Router"} {
					if bf, ok := BeeApp.Handlers.filters[k]; ok {
						filterType = fr
						filterTypes = append(filterTypes, filterType)
						resultList := new([][]string)
						for _, f := range bf {
							var result = []string{
								fmt.Sprintf("%s", f.pattern),
								fmt.Sprintf("%s", utils.GetFuncName(f.filterFunc)),
							}
							*resultList = append(*resultList, result)
						}
						filterTypeData[filterType] = resultList
					}
				}
			}

			content["Data"] = filterTypeData
			content["Methods"] = filterTypes

			data["Content"] = content
			data["Title"] = "Filters"
			execTpl(rw, data, routerAndFilterTpl, defaultScriptsTpl)
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
					fmt.Sprintf("%s", v.pattern),
					fmt.Sprintf("%s", v.methods),
					fmt.Sprintf("%s", v.controllerType),
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeRESTFul {
				var result = []string{
					fmt.Sprintf("%s", v.pattern),
					fmt.Sprintf("%s", v.methods),
					fmt.Sprintf(""),
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeHandler {
				var result = []string{
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
	format := r.Form.Get("format")
	data := make(map[interface{}]interface{})

	var result bytes.Buffer
	if command != "" {
		toolbox.ProcessInput(command, &result)
		data["Content"] = result.String()

		if format == "json" && command == "gc summary" {
			dataJSON, err := json.Marshal(data)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.Write(dataJSON)
			return
		}

		data["Title"] = command
		defaultTpl := defaultScriptsTpl
		if command == "gc summary" {
			defaultTpl = gcAjaxTpl
		}
		execTpl(rw, data, profillingTpl, defaultTpl)
	}
}

// Healthcheck is a http.Handler calling health checking and showing the result.
// it's in "/healthcheck" pattern in admin module.
func healthcheck(rw http.ResponseWriter, req *http.Request) {
	data := make(map[interface{}]interface{})

	var result = []string{}
	fields := []string{
		fmt.Sprintf("Name"),
		fmt.Sprintf("Message"),
		fmt.Sprintf("Status"),
	}
	resultList := new([][]string)

	content := make(map[string]interface{})

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

	content["Fields"] = fields
	content["Data"] = resultList
	data["Content"] = content
	data["Title"] = "Health Check"
	execTpl(rw, data, healthCheckTpl, defaultScriptsTpl)
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
			data["Message"] = []string{"success", fmt.Sprintf("%s run success,Now the Status is <br>%s", taskname, t.GetStatus())}
		} else {
			data["Message"] = []string{"warning", fmt.Sprintf("there's no task which named: %s", taskname)}
		}
	}

	// List Tasks
	content := make(map[string]interface{})
	resultList := new([][]string)
	var result = []string{}
	var fields = []string{
		fmt.Sprintf("Task Name"),
		fmt.Sprintf("Task Spec"),
		fmt.Sprintf("Task Status"),
		fmt.Sprintf("Last Time"),
		fmt.Sprintf(""),
	}
	for tname, tk := range toolbox.AdminTaskList {
		result = []string{
			fmt.Sprintf("%s", tname),
			fmt.Sprintf("%s", tk.GetSpec()),
			fmt.Sprintf("%s", tk.GetStatus()),
			fmt.Sprintf("%s", tk.GetPrev().String()),
		}
		*resultList = append(*resultList, result)
	}

	content["Fields"] = fields
	content["Data"] = resultList
	data["Content"] = content
	data["Title"] = "Tasks"
	execTpl(rw, data, tasksTpl, defaultScriptsTpl)
}

func execTpl(rw http.ResponseWriter, data map[interface{}]interface{}, tpls ...string) {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
	for _, tpl := range tpls {
		tmpl = template.Must(tmpl.Parse(tpl))
	}
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
	addr := AdminHTTPAddr

	if AdminHTTPPort != 0 {
		addr = fmt.Sprintf("%s:%d", AdminHTTPAddr, AdminHTTPPort)
	}
	for p, f := range admin.routers {
		http.Handle(p, f)
	}
	BeeLogger.Info("Admin server Running on %s", addr)

	var err error
	if Graceful {
		err = grace.ListenAndServe(addr, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		BeeLogger.Critical("Admin ListenAndServe: ", err, fmt.Sprintf("%d", os.Getpid()))
	}
}
