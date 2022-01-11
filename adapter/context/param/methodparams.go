package param

import (
	"github.com/beego/beego/v2/server/web/context/param"
)

// MethodParam keeps param information to be auto passed to controller methods
type MethodParam param.MethodParam

// New creates a new MethodParam with name and specific options
func New(name string, opts ...MethodParamOption) *MethodParam {
	newOps := make([]param.MethodParamOption, 0, len(opts))
	for _, o := range opts {
		newOps = append(newOps, oldMpoToNew(o))
	}
	return (*MethodParam)(param.New(name, newOps...))
}

// Make creates an array of MethodParmas or an empty array
func Make(list ...*MethodParam) []*MethodParam {
	if len(list) > 0 {
		return list
	}
	return nil
}

func (mp *MethodParam) String() string {
	return (*param.MethodParam)(mp).String()
}
