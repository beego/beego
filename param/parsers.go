package param

import "strconv"

type paramParser interface {
	parse(value string) (interface{}, error)
	zeroValue() interface{}
}

type boolParser struct {
}

func (p boolParser) parse(value string) (interface{}, error) {
	return strconv.ParseBool(value)
}

func (p boolParser) zeroValue() interface{} {
	return false
}

type stringParser struct {
}

func (p stringParser) parse(value string) (interface{}, error) {
	return value, nil
}

func (p stringParser) zeroValue() interface{} {
	return ""
}

type intParser struct {
}

func (p intParser) parse(value string) (interface{}, error) {
	return strconv.Atoi(value)
}

func (p intParser) zeroValue() interface{} {
	return 0
}
