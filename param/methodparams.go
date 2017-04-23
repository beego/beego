package param

//Keeps param information to be auto passed to controller methods
type MethodParam struct {
	name     string
	parser   paramParser
	location paramLocation
	required bool
	defValue string
}

type paramLocation byte

const (
	param paramLocation = iota
	body
	header
)

type MethodParamOption func(*MethodParam)

func Bool(name string, opts ...MethodParamOption) *MethodParam {
	return newParam(name, boolParser{}, opts)
}

func String(name string, opts ...MethodParamOption) *MethodParam {
	return newParam(name, stringParser{}, opts)
}

func Int(name string, opts ...MethodParamOption) *MethodParam {
	return newParam(name, intParser{}, opts)
}

func Float(name string, opts ...MethodParamOption) *MethodParam {
	return newParam(name, floatParser{}, opts)
}

func Time(name string, opts ...MethodParamOption) *MethodParam {
	return newParam(name, timeParser{}, opts)
}

func Json(name string, opts ...MethodParamOption) *MethodParam {
	return newParam(name, jsonParser{}, opts)
}

func AsSlice(param *MethodParam) *MethodParam {
	param.parser = sliceParser(param.parser)
	return param
}

func newParam(name string, parser paramParser, opts []MethodParamOption) (param *MethodParam) {
	param = &MethodParam{name: name, parser: parser}
	for _, option := range opts {
		option(param)
	}
	return
}
