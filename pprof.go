package beego

import (
	"net/http/pprof"
)

type ProfController struct {
	Controller
}

func (this *ProfController) Get() {
	ptype := this.Ctx.Params[":pp"]
	if ptype == "" {
		pprof.Index(this.Ctx.ResponseWriter, this.Ctx.Request)
	} else if ptype == "cmdline" {
		pprof.Cmdline(this.Ctx.ResponseWriter, this.Ctx.Request)
	} else if ptype == "profile" {
		pprof.Profile(this.Ctx.ResponseWriter, this.Ctx.Request)
	} else if ptype == "symbol" {
		pprof.Symbol(this.Ctx.ResponseWriter, this.Ctx.Request)
	}
}
