package controller

import (
	"github.com/astaxie/beego/beego"
	"../model"
)

func UserIndex(w beego.A) {
	userinfo :=model.User.getAll()
	beego.V.Render(w,"users/index")
}