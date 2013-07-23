package validation

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	VALIDTAG = "valid"
)

var (
	// key: function name
	// value: the number of parameters
	funcs = make(map[string]int)

	// doesn't belong to validation functions
	unFuncs = map[string]bool{
		"Clear":     true,
		"HasErrors": true,
		"ErrorMap":  true,
		"Error":     true,
		"apply":     true,
		"Check":     true,
		"Valid":     true,
	}
)

func init() {
	v := &Validation{}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !unFuncs[m.Name] {
			funcs[m.Name] = m.Type.NumIn() - 3
		}
	}
}

type ValidFunc struct {
	Name   string
	Params []interface{}
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func getValidFuncs(f reflect.StructField) (vfs []ValidFunc, err error) {
	tag := f.Tag.Get(VALIDTAG)
	if len(tag) == 0 {
		return
	}
	fs := strings.Split(tag, ";")
	for _, vfunc := range fs {
		var vf ValidFunc
		vf, err = parseFunc(vfunc)
		if err != nil {
			return
		}
		vfs = append(vfs, vf)
	}
	return
}

func parseFunc(vfunc string) (v ValidFunc, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	vfunc = strings.TrimSpace(vfunc)
	start := strings.Index(vfunc, "(")
	var num int

	// doesn't need parameter valid function
	if start == -1 {
		if num, err = numIn(vfunc); err != nil {
			return
		}
		if num != 0 {
			err = fmt.Errorf("%s require %d parameters", vfunc, num)
			return
		}
		v = ValidFunc{Name: vfunc}
		return
	}

	end := strings.Index(vfunc, ")")
	if end == -1 {
		err = fmt.Errorf("invalid valid function")
		return
	}

	name := strings.TrimSpace(vfunc[:start])
	if num, err = numIn(name); err != nil {
		return
	}

	params := strings.Split(vfunc[start+1:end], ",")
	// the num of param must be equal
	if num != len(params) {
		err = fmt.Errorf("%s require %d parameters", name, num)
		return
	}

	v = ValidFunc{name, trim(params)}
	return
}

func numIn(name string) (num int, err error) {
	num, ok := funcs[name]
	if !ok {
		err = fmt.Errorf("doesn't exsits %s valid function", name)
	}
	return
}

func trim(s []string) []interface{} {
	ts := make([]interface{}, len(s))
	for i := 0; i < len(s); i++ {
		ts[i] = strings.TrimSpace(s[i])
	}
	return ts
}
