package param

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type paramParser interface {
	parse(value string, toType reflect.Type) (interface{}, error)
}

type parserFunc func(value string, toType reflect.Type) (interface{}, error)

func (f parserFunc) parse(value string, toType reflect.Type) (interface{}, error) {
	return f(value, toType)
}

type boolParser struct {
}

func (p boolParser) parse(value string, toType reflect.Type) (interface{}, error) {
	return strconv.ParseBool(value)
}

type stringParser struct {
}

func (p stringParser) parse(value string, toType reflect.Type) (interface{}, error) {
	return value, nil
}

type intParser struct {
}

func (p intParser) parse(value string, toType reflect.Type) (interface{}, error) {
	return strconv.Atoi(value)
}

type floatParser struct {
}

func (p floatParser) parse(value string, toType reflect.Type) (interface{}, error) {
	if toType.Kind() == reflect.Float32 {
		res, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		return float32(res), nil
	}
	return strconv.ParseFloat(value, 64)
}

type timeParser struct {
}

func (p timeParser) parse(value string, toType reflect.Type) (result interface{}, err error) {
	result, err = time.Parse(time.RFC3339, value)
	if err != nil {
		result, err = time.Parse("2006-01-02", value)
	}
	return
}

type jsonParser struct {
}

func (p jsonParser) parse(value string, toType reflect.Type) (interface{}, error) {
	pResult := reflect.New(toType)
	v := pResult.Interface()
	err := json.Unmarshal([]byte(value), v)
	if err != nil {
		return nil, err
	}
	return pResult.Elem().Interface(), nil
}

func sliceParser(elemParser paramParser) paramParser {
	return parserFunc(func(value string, toType reflect.Type) (interface{}, error) {
		values := strings.Split(value, ",")
		result := reflect.MakeSlice(toType, 0, len(values))
		elemType := toType.Elem()
		for _, v := range values {
			parsedValue, err := elemParser.parse(v, elemType)
			if err != nil {
				return nil, err
			}
			result = reflect.Append(result, reflect.ValueOf(parsedValue))
		}
		return result.Interface(), nil
	})
}
