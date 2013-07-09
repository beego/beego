package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/example/beeapi/models"
)

type ResponseInfo struct {
}

type ObejctController struct {
	beego.Controller
}

func (this *ObejctController) Post() {
	var ob models.Object
	json.Unmarshal(this.Ctx.RequestBody, &ob)
	objectid := models.AddOne(ob)
	this.Data["json"] = "{\"ObjectId\":\"" + objectid + "\"}"
	this.ServeJson()
}

func (this *ObejctController) Get() {
	objectId := this.Ctx.Params[":objectId"]
	if objectId != "" {
		ob, err := models.GetOne(objectId)
		if err != nil {
			this.Data["json"] = err
		} else {
			this.Data["json"] = ob
		}
	} else {
		obs := models.GetAll()
		this.Data["json"] = obs
	}
	this.ServeJson()
}

func (this *ObejctController) Put() {
	objectId := this.Ctx.Params[":objectId"]
	var ob models.Object
	json.Unmarshal(this.Ctx.RequestBody, &ob)

	err := models.Update(objectId, ob.Score)
	if err != nil {
		this.Data["json"] = err
	} else {
		this.Data["json"] = "update success!"
	}
	this.ServeJson()
}

func (this *ObejctController) Delete() {
	objectId := this.Ctx.Params[":objectId"]
	models.Delete(objectId)
	this.Data["json"] = "delete success!"
	this.ServeJson()
}
