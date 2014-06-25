// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package beego

import (
	"encoding/json"

	"github.com/astaxie/beego/context"
)

var GlobalDocApi map[string]interface{}

func init() {
	if EnableDocs {
		GlobalDocApi = make(map[string]interface{})
	}
}

func serverDocs(ctx *context.Context) {
	var obj interface{}
	if splat := ctx.Input.Param(":splat"); splat == "" {
		obj = GlobalDocApi["Root"]
	} else {
		if v, ok := GlobalDocApi[splat]; ok {
			obj = v
		}
	}
	if obj != nil {
		bt, err := json.Marshal(obj)
		if err != nil {
			ctx.Output.SetStatus(504)
			return
		}
		ctx.Output.Header("Content-Type", "application/json;charset=UTF-8")
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		ctx.Output.Body(bt)
		return
	}
	ctx.Output.SetStatus(404)
}
