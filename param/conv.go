package param

import (
	"fmt"
	"reflect"

	beecontext "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

func convertParam(param *MethodParam, paramType reflect.Type, ctx *beecontext.Context) (result reflect.Value) {
	var strValue string
	var reflectValue reflect.Value
	switch param.location {
	case body:
		strValue = string(ctx.Input.RequestBody)
	case header:
		strValue = ctx.Input.Header(param.name)
	default:
		strValue = ctx.Input.Query(param.name)
	}

	if strValue == "" {
		if param.required {
			ctx.Abort(400, fmt.Sprintf("Missing parameter %s", param.name))
		} else {
			strValue = param.defValue
		}
	}
	if strValue == "" {
		reflectValue = reflect.Zero(paramType)
	} else {
		value, err := param.parser.parse(strValue, paramType)
		if err != nil {
			logs.Debug(fmt.Sprintf("Error converting param %s to type %s. Value: %s, Error: %s", param.name, paramType, strValue, err))
			ctx.Abort(400, fmt.Sprintf("Invalid parameter %s. Can not convert %s to type %s", param.name, strValue, paramType))
		}

		reflectValue, err = safeConvert(reflect.ValueOf(value), paramType)
		if err != nil {
			panic(err)
		}
	}
	return reflectValue
}

func ConvertParams(methodParams []*MethodParam, methodType reflect.Type, ctx *beecontext.Context) (result []reflect.Value) {
	result = make([]reflect.Value, 0, len(methodParams))
	for i := 0; i < len(methodParams); i++ {
		reflectValue := convertParam(methodParams[i], methodType.In(i), ctx)
		result = append(result, reflectValue)
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
