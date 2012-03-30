package beego

import {
    "net"
	"net/http"
	"net/http/fcgi"
    "log"
    "strconv"
    "./core"
}
type C struct {
	core.Content
}

type M struct{
	core.Model
}

type D struct{
	core.Config
}

type U struct{
	core.URL
}

type A struct{
	core.Controller
}

type V struct{
	core.View
}

type BeegoApp struct{
	Port int
}

func (app *BeegoApp) BeeListen(port int) {
	app.Port = port
	err := http.ListenAndServe(":"+strconv.Itoa(app.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (app *BeegoApp) BeeListenFcgi(port int) {
	app.Port = port
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fcgi.Serve(l, app.Handler)
}

func Run() {
	rootPath, _ := os.Getwd()
}