package beego

import (
	"fmt"
	"net/http"
)

var BeeAdminApp *AdminApp

func init() {
	BeeAdminApp = &AdminApp{
		routers: make(map[string]http.HandlerFunc),
	}
	BeeAdminApp.Route("/", AdminIndex)
}

func AdminIndex(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Welcome to Admin Dashboard"))
}

type AdminApp struct {
	routers map[string]http.HandlerFunc
}

func (admin *AdminApp) Route(pattern string, f http.HandlerFunc) {
	admin.routers[pattern] = f
}

func (admin *AdminApp) Run() {
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
