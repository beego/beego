// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/example/beeapi/models"
)

type ObjectController struct {
	beego.Controller
}

func (o *ObjectController) Post() {
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	objectid := models.AddOne(ob)
	o.Data["json"] = map[string]string{"ObjectId": objectid}
	o.ServeJson()
}

func (o *ObjectController) Get() {
	objectId := o.Ctx.Input.Params[":objectId"]
	if objectId != "" {
		ob, err := models.GetOne(objectId)
		if err != nil {
			o.Data["json"] = err
		} else {
			o.Data["json"] = ob
		}
	} else {
		obs := models.GetAll()
		o.Data["json"] = obs
	}
	o.ServeJson()
}

func (o *ObjectController) Put() {
	objectId := o.Ctx.Input.Params[":objectId"]
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)

	err := models.Update(objectId, ob.Score)
	if err != nil {
		o.Data["json"] = err
	} else {
		o.Data["json"] = "update success!"
	}
	o.ServeJson()
}

func (o *ObjectController) Delete() {
	objectId := o.Ctx.Input.Params[":objectId"]
	models.Delete(objectId)
	o.Data["json"] = "delete success!"
	o.ServeJson()
}
