// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package forms

import (
	"github.com/astaxie/beego/forms/elements"
)

type Form struct {
	elements map[string]elements.ElementInterface
}

func NewForm() *Form {
	f := &Form{}
	f.elements = make(map[string]elements.ElementInterface)
	return f
}

func (f *Form) Init() *Form {
	f.elements = make(map[string]elements.ElementInterface)
}

func (f *Form) Valid(data map[string]interface{}) bool {
	for k, val := range data {
		if e, ok := f.elements[k]; ok {
			if !e.Valid(val) {
				return false
			}
		}
	}
	return true
}

func (f *Form) SetData() {

}

func (f *Form) Bind() {

}

func (f *Form) SaveData() {

}

func (f *Form) AddElement(e elements.ElementInterface) {
	name := e.GetName()
	f.elements[name] = e
}

func (f *Form) Render() string {

}
