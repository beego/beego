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
	"html/template"
	"net/url"
	"reflect"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
)

type FormController interface {
	GetSession(name interface{}) interface{}
	SetSession(name interface{}, value interface{})
	GetCtx() *context.Context
	Input() url.Values
	GetLocale() FormLocaler
	Tr(format string, args ...interface{}) string // from i18n
}

// check form once, void re-submit
func FormOnceNotMatch(controller FormController) bool {
	notMatch := false
	recreat := false

	// get token from request param / header
	var value string
	if vus, ok := controller.Input()["_once"]; ok && len(vus) > 0 {
		value = vus[0]
	} else {
		value = controller.GetCtx().Input.Header("X-Form-Once")
	}

	// exist in session
	if v, ok := controller.GetSession("form_once").(string); ok && v != "" {
		// not match
		if value != v {
			notMatch = true
		} else {
			// if matched then re-creat once
			recreat = true
		}
	}

	FormOnceCreate(controller, recreat)
	return notMatch
}

// create form once html
func FormOnceCreate(controller FormController, args ...bool) {
	var value string
	var creat bool
	creat = len(args) > 0 && args[0]
	if !creat {
		if v, ok := controller.GetSession("form_once").(string); ok && v != "" {
			value = v
		} else {
			creat = true
		}
	}
	if creat {
		value = GetRandomString(10)
		controller.SetSession("form_once", value)
	}
	controller.GetCtx().Input.Data["once_token"] = value
	controller.GetCtx().Input.Data["once_html"] = template.HTML(`<input type="hidden" name="_once" value="` + value + `">`)
}

func validForm(controller FormController, form interface{}, names ...string) (bool, map[string]*validation.ValidationError) {
	// parse request params to form ptr struct
	ParseForm(form, controller.Input())

	// Put data back in case users input invalid data for any section.
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	controller.GetCtx().Input.Data[name] = form

	errName := name + "Error"

	// check form once
	if FormOnceNotMatch(controller) {
		return false, nil
	}

	// Verify basic input.
	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := valid.ErrorMap()
		controller.GetCtx().Input.Data[errName] = &valid
		return false, errs
	}
	return true, nil
}

// valid form and put errors to tempalte context
func ValidForm(controller FormController, form interface{}, names ...string) bool {
	valid, _ := validForm(controller, form, names...)
	return valid
}

// valid form and put errors to tempalte context
func ValidFormSets(controller FormController, form interface{}, names ...string) bool {
	valid, errs := validForm(controller, form, names...)
	setFormSets(controller, form, errs, names...)
	return valid
}

func SetFormSets(controller FormController, form interface{}, names ...string) *FormSets {
	return setFormSets(controller, form, nil, names...)
}

func setFormSets(controller FormController, form interface{}, errs map[string]*validation.ValidationError, names ...string) *FormSets {
	formSets := NewFormSets(form, errs, controller.GetLocale())
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	name += "Sets"
	controller.GetCtx().Input.Data[name] = formSets

	return formSets
}

// add valid error to FormError
func SetFormError(controller FormController, form interface{}, fieldName, errMsg string, names ...string) {
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	errName := name + "Error"
	setsName := name + "Sets"

	if valid, ok := controller.GetCtx().Input.Data[errName].(*validation.Validation); ok {
		valid.SetError(fieldName, controller.Tr(errMsg))
	}

	if fSets, ok := controller.GetCtx().Input.Data[setsName].(*FormSets); ok {
		fSets.SetError(fieldName, errMsg)
	}
}
