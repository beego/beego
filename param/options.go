package param

import (
	"fmt"
)

type MethodParamOption func(*MethodParam)

var IsRequired MethodParamOption = func(p *MethodParam) {
	p.required = true
}

var InHeader MethodParamOption = func(p *MethodParam) {
	p.location = header
}

var InPath MethodParamOption = func(p *MethodParam) {
	p.location = path
}

var InBody MethodParamOption = func(p *MethodParam) {
	p.location = body
}

func Default(defValue interface{}) MethodParamOption {
	return func(p *MethodParam) {
		if defValue != nil {
			p.defValue = fmt.Sprint(defValue)
		}
	}
}
