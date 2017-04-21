package param

var InHeader MethodParamOption = func(p *MethodParam) {
	p.location = header
}

var IsRequired MethodParamOption = func(p *MethodParam) {
	p.required = true
}

func Default(defValue interface{}) MethodParamOption {
	return func(p *MethodParam) {
		p.defValue = defValue
	}
}
