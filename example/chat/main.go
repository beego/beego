package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/example/chat/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WSController{})
	beego.Run()
}
