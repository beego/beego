// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package forms

import (
	"fmt"
	"html/template"
	"reflect"
	"time"

	"github.com/astaxie/beego"
)

func init() {
	initCommonFields()
}

func initCommonFields() {
	RegisterFieldCreater("text", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="text" value="%v" class="form-control"%s%s>`,
				fSet.Id, fSet.Name, fSet.Value, fSet.Placeholder, fSet.Attrs))
		}
	})

	RegisterFieldCreater("textarea", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<textarea id="%s" name="%s" rows="5" class="form-control"%s%s>%v</textarea>`,
				fSet.Id, fSet.Name, fSet.Placeholder, fSet.Attrs, fSet.Value))
		}
	})

	RegisterFieldCreater("password", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="password" value="%v" class="form-control"%s%s>`,
				fSet.Id, fSet.Name, fSet.Value, fSet.Placeholder, fSet.Attrs))
		}
	})

	RegisterFieldCreater("hidden", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="hidden" value="%v"%s>`, fSet.Id, fSet.Name, fSet.Value, fSet.Attrs))
		}
	})

	datetimeFunc := func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			t := fSet.Value.(time.Time)
			tval := beego.Date(t, DateTimeFormat)
			if tval == "0001-01-01 00:00:00" {
				t = time.Now()
			}
			if fSet.Type == "date" {
				tval = beego.Date(t, DateOnlyFormat)
			}
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="%s" value="%s" class="form-control"%s%s>`,
				fSet.Id, fSet.Name, fSet.Type, tval, fSet.Placeholder, fSet.Attrs))
		}
	}

	RegisterFieldCreater("date", datetimeFunc)
	RegisterFieldCreater("datetime", datetimeFunc)

	RegisterFieldCreater("checkbox", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			var checked string
			if b, ok := fSet.Value.(bool); ok && b {
				checked = "checked"
			}
			return template.HTML(fmt.Sprintf(
				`<label for="%s" class="checkbox">%s<input id="%s" name="%s" type="checkbox" %s></label>`,
				fSet.Id, fSet.LabelText, fSet.Id, fSet.Name, checked))
		}
	})

	RegisterFieldCreater("select", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			var options string
			str := fmt.Sprintf(`<select id="%s" name="%s" class="form-control"%s%s>%s</select>`,
				fSet.Id, fSet.Name, fSet.Placeholder, fSet.Attrs)

			fun := fSet.FormElm.Addr().MethodByName(fSet.Name + "SelectData")

			if fun.IsValid() {
				results := fun.Call([]reflect.Value{})
				if len(results) > 0 {
					v := results[0]
					if v.CanInterface() {
						if vu, ok := v.Interface().([][]string); ok {

							var vs []string
							val := reflect.ValueOf(fSet.Value)
							if val.Kind() == reflect.Slice {
								vs = make([]string, 0, val.Len())
								for i := 0; i < val.Len(); i++ {
									vs = append(vs, ToStr(val.Index(i).Interface()))
								}
							}

							isMulti := len(vs) > 0
							for _, parts := range vu {
								var n, v string
								switch {
								case len(parts) > 1:
									n, v = fSet.Locale.Tr(parts[0]), parts[1]
								case len(parts) == 1:
									n, v = fSet.Locale.Tr(parts[0]), parts[0]
								}
								var selected string
								if isMulti {
									for _, e := range vs {
										if e == v {
											selected = ` selected="selected"`
											break
										}
									}
								} else if ToStr(fSet.Value) == v {
									selected = ` selected="selected"`
								}
								options += fmt.Sprintf(`<option value="%s"%s>%s</option>`, v, selected, n)
							}
						}
					}
				}
			}

			if len(options) == 0 {
				options = fmt.Sprintf(`<option value="%v">%v</option>`, fSet.Value, fSet.Value)
			}

			return template.HTML(fmt.Sprintf(str, options))
		}
	})
}
