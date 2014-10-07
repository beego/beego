/*

The forms package provides utilities to define forms as structs,
render them in the view templates and process / validate the incoming requests.

An example CreateUser action

Define the action form (actions/user/new.go)

  package user

  import (
    "github.com/beego/i18n"
  )

  type NewUser struct {
    Locale   i18n.Locale `form:"-"`
    Username string      `form:"type(text)" valid:"Required;MinSize(3)"`
    Password string      `form:"type(password)" valid:"Required;MinSize(3)"`
  }

  func (this *NewUser) Placeholders() map[string]string {
    return map[string]string{
      "Username": "new_user.username_placeholder",
      "Password": "new_user.password_placeholder",
    }
  }

  func (this *NewUser) Save() (err error) {
    // use this.Username and this.Password to create the user and return
    return
  }

Define the routes (routes/users.go)

  package routers

  import "controllers"

  func init() {
    beego.Router("/users/new", &controllers.UsersController{}, "get:New")
    beego.Router("/users", &controllers.UsersController{}, "post:Create")
  }

Setup the controller (controllers/users.go)

  package controllers

  import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/utils/forms"

    "actions/user"
  )

  type UsersController struct {
    beego.Controller
  }

  func (this *UsersController) Prepare() {
    this.Layout = "layouts/default.html.tpl"
  }

  func (this *UsersController) New() {
    this.TplNames = "users/new.html.tpl"
    form := user.NewUser{}
    forms.SetFormSets(this, &form)
  }

  func (this *UsersController) Create() {
    this.TplNames = "users/new.html.tpl"
    form := user.NewUser{}
    forms.SetFormSets(this, &form)
    if !forms.ValidFormSets(this, &form) {
      return
    }
    err := form.Save()
    if err != nil {
      return
    }
    this.Redirect("/users/"+form.Username, 302)
  }

Design the form template (views/users/new.html.tpl)

  <form action="{{urlfor "UsersController.Create"}}" method="post" class="form-horizontal" role="form">
    {{.xsrf_html | str2html}}
    <div class="form-group">
      <h1>Create User</h1>
    </div>
    {{template "shared/horizontal_form/fields.html.tpl" .CreateUserFormSets}}
    <div class="form-group">
      <button type="submit" class="btn btn-default">Create User</button>
    </div>
  </form>

The required Twitter Bootstrap horizontal form template

This template is meant to be reused within a project across different actions.

views/shared/horizontal_form/fields.html.tpl

  {{range $field := .FieldList}}
      {{template "shared/horizontal_form/group.html.tpl" $field}}
  {{end}}

view/shared/horizontal_form/group.html.tpl

  {{if eq .Type "hidden"}}
      {{call .Field}}
  {{else}}
      <div class="form-group{{if .Error}} has-error{{end}}">
        {{template "shared/horizontal_form/field.html.tpl" .}}
      </div>
  {{end}}

view/shared/horizontal_form/field.html.tpl

  {{if eq .Type "hidden"}}
      {{call .Field}}
  {{else}}
    <label class="col-md-3 control-label" for="{{.Id}}">{{.LabelText}}</label>
    <div class="col-md-6">
      {{call .Field}}
      {{if .Error}}<p class="error-block">{{.Error}}</p>{{end}}
      {{if .Help}}<p class="help-block">{{.Help}}</p>{{end}}
    </div>
  {{end}}

*/
package forms
