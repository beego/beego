package param

import (
	"github.com/beego/beego/v2/server/web/context/param"
)

// MethodParamOption defines a func which apply options on a MethodParam
type MethodParamOption func(*MethodParam)

// IsRequired indicates that this param is required and can not be omitted from the http request
var IsRequired MethodParamOption = func(p *MethodParam) {
	param.IsRequired((*param.MethodParam)(p))
}

// InHeader indicates that this param is passed via an http header
var InHeader MethodParamOption = func(p *MethodParam) {
	param.InHeader((*param.MethodParam)(p))
}

// InPath indicates that this param is part of the URL path
var InPath MethodParamOption = func(p *MethodParam) {
	param.InPath((*param.MethodParam)(p))
}

// InBody indicates that this param is passed as an http request body
var InBody MethodParamOption = func(p *MethodParam) {
	param.InBody((*param.MethodParam)(p))
}

// Default provides a default value for the http param
func Default(defaultValue interface{}) MethodParamOption {
	return newMpoToOld(param.Default(defaultValue))
}

func newMpoToOld(n param.MethodParamOption) MethodParamOption {
	return func(methodParam *MethodParam) {
		n((*param.MethodParam)(methodParam))
	}
}

func oldMpoToNew(old MethodParamOption) param.MethodParamOption {
	return func(methodParam *param.MethodParam) {
		old((*MethodParam)(methodParam))
	}
}
