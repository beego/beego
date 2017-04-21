package param

import (
	"fmt"
	"reflect"

	beecontext "github.com/astaxie/beego/context"
)

//Keeps param information to be auto passed to controller methods
type MethodParam struct {
	name     string
	parser   paramParser
	location paramLocation
	required bool
	defValue interface{}
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

func newParam(name string, parser paramParser, opts []MethodParamOption) (param *MethodParam) {
	param = &MethodParam{name: name, parser: parser}
	for _, option := range opts {
		option(param)
	}
	return
}

func ConvertParams(methodParams []*MethodParam, methodType reflect.Type, ctx *beecontext.Context) (result []reflect.Value) {
	result = make([]reflect.Value, 0, len(methodParams))
	i := 0
	for _, p := range methodParams {
		var strValue string
		var value interface{}
		switch p.location {
		case body:
			strValue = string(ctx.Input.RequestBody)
		case header:
			strValue = ctx.Input.Header(p.name)
		default:
			strValue = ctx.Input.Query(p.name)
		}

		if strValue == "" {
			if p.required {
				ctx.Abort(400, "Missing argument "+p.name)
			} else if p.defValue != nil {
				value = p.defValue
			} else {
				value = p.parser.zeroValue()
			}
		} else {
			var err error
			value, err = p.parser.parse(strValue)
			if err != nil {
				//TODO: handle err
			}
		}
		reflectValue, err := safeConvert(reflect.ValueOf(value), methodType.In(i))
		if err != nil {
			//TODO: handle err
		}
		result = append(result, reflectValue)
		i++
	}
	return
}

func safeConvert(value reflect.Value, t reflect.Type) (result reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	result = value.Convert(t)
	return
}
