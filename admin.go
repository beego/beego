package beego

import (
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego/toolbox"
)

var BeeAdminApp *AdminApp

//func MyFilterMonitor(method, requestPath string, t time.Duration) bool {
//	if method == "POST" {
//		return false
//	}
//	if t.Nanoseconds() < 100 {
//		return false
//	}
//	if strings.HasPrefix(requestPath, "/astaxie") {
//		return false
//	}
//	return true
//}

//beego.FilterMonitorFunc = MyFilterMonitor
var FilterMonitorFunc func(string, string, time.Duration) bool

func init() {
	BeeAdminApp = &AdminApp{
		routers: make(map[string]http.HandlerFunc),
	}
	BeeAdminApp.Route("/", AdminIndex)
	BeeAdminApp.Route("/qps", QpsIndex)
	BeeAdminApp.Route("/prof", ProfIndex)
	BeeAdminApp.Route("/healthcheck", toolbox.Healthcheck)
	BeeAdminApp.Route("/task", toolbox.TaskStatus)
	BeeAdminApp.Route("/runtask", toolbox.RunTask)
	FilterMonitorFunc = func(string, string, time.Duration) bool { return true }
}

func AdminIndex(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Welcome to Admin Dashboard\n"))
	rw.Write([]byte("There are servral functions:\n"))
	rw.Write([]byte("1. Record all request and request time, http://localhost:8088/qps\n"))
	rw.Write([]byte("2. Get runtime profiling data by the pprof, http://localhost:8088/prof\n"))
	rw.Write([]byte("3. Get healthcheck result from http://localhost:8088/prof\n"))
	rw.Write([]byte("4. Get current task infomation from taskhttp://localhost:8088/task \n"))
	rw.Write([]byte("5. To run a task passed a param http://localhost:8088/runtask\n"))

}

func QpsIndex(rw http.ResponseWriter, r *http.Request) {
	toolbox.StatisticsMap.GetMap(rw)
}

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

type AdminApp struct {
	routers map[string]http.HandlerFunc
}

func (admin *AdminApp) Route(pattern string, f http.HandlerFunc) {
	admin.routers[pattern] = f
}

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
