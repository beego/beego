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
	"reflect"
	"strconv"
	"text/template"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/beego/beego/grace"
	"github.com/beego/beego/logs"
	"github.com/beego/beego/toolbox"
	"github.com/beego/beego/utils"
)

// BeeAdminApp is the default adminApp used by admin module.
var beeAdminApp *adminApp

// FilterMonitorFunc is default monitor filter when admin module is enable.
// if this func returns, admin module records qps for this request by condition of this function logic.
// usage:
// 	func MyFilterMonitor(method, requestPath string, t time.Duration, pattern string, statusCode int) bool {
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
var FilterMonitorFunc func(string, string, time.Duration, string, int) bool

func init() {
	beeAdminApp = &adminApp{
		routers: make(map[string]http.HandlerFunc),
	}
	// keep in mind that all data should be html escaped to avoid XSS attack
	beeAdminApp.Route("/", adminIndex)
	beeAdminApp.Route("/qps", qpsIndex)
	beeAdminApp.Route("/prof", profIndex)
	beeAdminApp.Route("/healthcheck", healthcheck)
	beeAdminApp.Route("/task", taskStatus)
	beeAdminApp.Route("/listconf", listConf)
	beeAdminApp.Route("/metrics", promhttp.Handler().ServeHTTP)
	FilterMonitorFunc = func(string, string, time.Duration, string, int) bool { return true }
}

// AdminIndex is the default http.Handler for admin module.
// it matches url pattern "/".
func adminIndex(rw http.ResponseWriter, _ *http.Request) {
	writeTemplate(rw, map[interface{}]interface{}{}, indexTpl, defaultScriptsTpl)
}

// QpsIndex is the http.Handler for writing qps statistics map result info in http.ResponseWriter.
// it's registered with url pattern "/qps" in admin module.
func qpsIndex(rw http.ResponseWriter, _ *http.Request) {
	data := make(map[interface{}]interface{})
	data["Content"] = toolbox.StatisticsMap.GetMap()

	// do html escape before display path, avoid xss
	if content, ok := (data["Content"]).(M); ok {
		if resultLists, ok := (content["Data"]).([][]string); ok {
			for i := range resultLists {
				if len(resultLists[i]) > 0 {
					resultLists[i][0] = template.HTMLEscapeString(resultLists[i][0])
				}
			}
		}
	}

	writeTemplate(rw, data, qpsTpl, defaultScriptsTpl)
}

// ListConf is the http.Handler of displaying all beego configuration values as key/value pair.
// it's registered with url pattern "/listconf" in admin module.
func listConf(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")
	if command == "" {
		rw.Write([]byte("command not support"))
		return
	}

	data := make(map[interface{}]interface{})
	switch command {
	case "conf":
		m := make(M)
		list("BConfig", BConfig, m)
		m["AppConfigPath"] = template.HTMLEscapeString(appConfigPath)
		m["AppConfigProvider"] = template.HTMLEscapeString(appConfigProvider)
		tmpl := template.Must(template.New("dashboard").Parse(dashboardTpl))
		tmpl = template.Must(tmpl.Parse(configTpl))
		tmpl = template.Must(tmpl.Parse(defaultScriptsTpl))

		data["Content"] = m

		tmpl.Execute(rw, data)

	case "router":
		content := PrintTree()
		content["Fields"] = []string{
			"Router Pattern",
			"Methods",
			"Controller",
		}
		data["Content"] = content
		data["Title"] = "Routers"
		writeTemplate(rw, data, routerAndFilterTpl, defaultScriptsTpl)
	case "filter":
		var (
			content = M{
				"Fields": []string{
					"Router Pattern",
					"Filter Function",
				},
			}
			filterTypes    = []string{}
			filterTypeData = make(M)
		)

		if BeeApp.Handlers.enableFilter {
			var filterType string
			for k, fr := range map[int]string{
				BeforeStatic: "Before Static",
				BeforeRouter: "Before Router",
				BeforeExec:   "Before Exec",
				AfterExec:    "After Exec",
				FinishRouter: "Finish Router"} {
				if bf := BeeApp.Handlers.filters[k]; len(bf) > 0 {
					filterType = fr
					filterTypes = append(filterTypes, filterType)
					resultList := new([][]string)
					for _, f := range bf {
						var result = []string{
							// void xss
							template.HTMLEscapeString(f.pattern),
							template.HTMLEscapeString(utils.GetFuncName(f.filterFunc)),
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
		writeTemplate(rw, data, routerAndFilterTpl, defaultScriptsTpl)
	default:
		rw.Write([]byte("command not support"))
	}
}

func list(root string, p interface{}, m M) {
	pt := reflect.TypeOf(p)
	pv := reflect.ValueOf(p)
	if pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
		pv = pv.Elem()
	}
	for i := 0; i < pv.NumField(); i++ {
		var key string
		if root == "" {
			key = pt.Field(i).Name
		} else {
			key = root + "." + pt.Field(i).Name
		}
		if pv.Field(i).Kind() == reflect.Struct {
			list(key, pv.Field(i).Interface(), m)
		} else {
			m[key] = pv.Field(i).Interface()
		}
	}
}

// PrintTree prints all registered routers.
func PrintTree() M {
	var (
		content     = M{}
		methods     = []string{}
		methodsData = make(M)
	)
	for method, t := range BeeApp.Handlers.routers {

		resultList := new([][]string)

		printTree(resultList, t)

		methods = append(methods, template.HTMLEscapeString(method))
		methodsData[template.HTMLEscapeString(method)] = resultList
	}

	content["Data"] = methodsData
	content["Methods"] = methods
	return content
}

func printTree(resultList *[][]string, t *Tree) {
	for _, tr := range t.fixrouters {
		printTree(resultList, tr)
	}
	if t.wildcard != nil {
		printTree(resultList, t.wildcard)
	}
	for _, l := range t.leaves {
		if v, ok := l.runObject.(*ControllerInfo); ok {
			if v.routerType == routerTypeBeego {
				var result = []string{
					template.HTMLEscapeString(v.pattern),
					template.HTMLEscapeString(fmt.Sprintf("%s", v.methods)),
					template.HTMLEscapeString(v.controllerType.String()),
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeRESTFul {
				var result = []string{
					template.HTMLEscapeString(v.pattern),
					template.HTMLEscapeString(fmt.Sprintf("%s", v.methods)),
					"",
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeHandler {
				var result = []string{
					template.HTMLEscapeString(v.pattern),
					"",
					"",
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
	if command == "" {
		return
	}

	var (
		format = r.Form.Get("format")
		data   = make(map[interface{}]interface{})
		result bytes.Buffer
	)
	toolbox.ProcessInput(command, &result)
	data["Content"] = template.HTMLEscapeString(result.String())

	if format == "json" && command == "gc summary" {
		dataJSON, err := json.Marshal(data)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(rw, dataJSON)
		return
	}

	data["Title"] = template.HTMLEscapeString(command)
	defaultTpl := defaultScriptsTpl
	if command == "gc summary" {
		defaultTpl = gcAjaxTpl
	}
	writeTemplate(rw, data, profillingTpl, defaultTpl)
}

// Healthcheck is a http.Handler calling health checking and showing the result.
// it's in "/healthcheck" pattern in admin module.
func healthcheck(rw http.ResponseWriter, r *http.Request) {
	var (
		result     []string
		data       = make(map[interface{}]interface{})
		resultList = new([][]string)
		content    = M{
			"Fields": []string{"Name", "Message", "Status"},
		}
	)

	for name, h := range toolbox.AdminCheckList {
		if err := h.Check(); err != nil {
			result = []string{
				"error",
				template.HTMLEscapeString(name),
				template.HTMLEscapeString(err.Error()),
			}
		} else {
			result = []string{
				"success",
				template.HTMLEscapeString(name),
				"OK",
			}
		}
		*resultList = append(*resultList, result)
	}

	queryParams := r.URL.Query()
	jsonFlag := queryParams.Get("json")
	shouldReturnJSON, _ := strconv.ParseBool(jsonFlag)

	if shouldReturnJSON {
		response := buildHealthCheckResponseList(resultList)
		jsonResponse, err := json.Marshal(response)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		} else {
			writeJSON(rw, jsonResponse)
		}
		return
	}

	content["Data"] = resultList
	data["Content"] = content
	data["Title"] = "Health Check"

	writeTemplate(rw, data, healthCheckTpl, defaultScriptsTpl)
}

func buildHealthCheckResponseList(healthCheckResults *[][]string) []map[string]interface{} {
	response := make([]map[string]interface{}, len(*healthCheckResults))

	for i, healthCheckResult := range *healthCheckResults {
		currentResultMap := make(map[string]interface{})

		currentResultMap["name"] = healthCheckResult[0]
		currentResultMap["message"] = healthCheckResult[1]
		currentResultMap["status"] = healthCheckResult[2]

		response[i] = currentResultMap
	}

	return response

}

func writeJSON(rw http.ResponseWriter, jsonData []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonData)
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
			if err := t.Run(); err != nil {
				data["Message"] = []string{"error", template.HTMLEscapeString(fmt.Sprintf("%s", err))}
			}
			data["Message"] = []string{"success", template.HTMLEscapeString(fmt.Sprintf("%s run success,Now the Status is <br>%s", taskname, t.GetStatus()))}
		} else {
			data["Message"] = []string{"warning", template.HTMLEscapeString(fmt.Sprintf("there's no task which named: %s", taskname))}
		}
	}

	// List Tasks
	content := make(M)
	resultList := new([][]string)
	var fields = []string{
		"Task Name",
		"Task Spec",
		"Task Status",
		"Last Time",
		"",
	}
	for tname, tk := range toolbox.AdminTaskList {
		result := []string{
			template.HTMLEscapeString(tname),
			template.HTMLEscapeString(tk.GetSpec()),
			template.HTMLEscapeString(tk.GetStatus()),
			template.HTMLEscapeString(tk.GetPrev().String()),
		}
		*resultList = append(*resultList, result)
	}

	content["Fields"] = fields
	content["Data"] = resultList
	data["Content"] = content
	data["Title"] = "Tasks"
	writeTemplate(rw, data, tasksTpl, defaultScriptsTpl)
}

func writeTemplate(rw http.ResponseWriter, data map[interface{}]interface{}, tpls ...string) {
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
	addr := BConfig.Listen.AdminAddr

	if BConfig.Listen.AdminPort != 0 {
		addr = fmt.Sprintf("%s:%d", BConfig.Listen.AdminAddr, BConfig.Listen.AdminPort)
	}
	for p, f := range admin.routers {
		http.Handle(p, f)
	}
	logs.Info("Admin server Running on %s", addr)

	var err error
	if BConfig.Listen.Graceful {
		err = grace.ListenAndServe(addr, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		logs.Critical("Admin ListenAndServe: ", err, fmt.Sprintf("%d", os.Getpid()))
	}
}
