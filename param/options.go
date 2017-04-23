package param

import (
	"fmt"
)

var IsRequired MethodParamOption = func(p *MethodParam) {
	p.required = true
}

var InHeader MethodParamOption = func(p *MethodParam) {
	p.location = header
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
