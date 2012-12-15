## Beego
=======
Beego is an open source version of the scalable, non-blocking web server
and tools that power SNDA's CDN system. Documentation and downloads are
available at https://github.com/astaxie/beego

Beego is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).

## Installation
============
To install:

    go get github.com/astaxie/beego

go version: go1 release

## Quick Start
============
Here is the canonical "Hello, world" example app for beego:

	package main
	
	import (
		"github.com/astaxie/beego"
	)
	
	type MainController struct {
		beego.Controller
	}
	
	func (this *MainController) Get() {
		this.Ct.WriteString("hello world")
	}
	
	func main() {
		beego.BeeApp.RegisterController("/", &MainController{})
		beego.BeeApp.Run()
	}
	
default port:8080

	http get http://localhost:8080
	HTTP/1.1 200 OK
	Content-Type: text/plain; charset=utf-8
	Date: Sat, 15 Dec 2012 16:03:00 GMT
	Transfer-Encoding: chunked
	
	hello world
	
	
## Router
============
In beego, a route is a struct paired with a URL-matching pattern. The strcut has many method with the same name of http method to server the http request. Each route is associated with a block:

	beego.BeeApp.RegisterController("/", &controllers.MainController{})
	beego.BeeApp.RegisterController("/admin", &admin.UserController{})
	beego.BeeApp.RegisterController("/admin/index", &admin.ArticleController{})
	beego.BeeApp.RegisterController("/admin/addpkg", &admin.AddController{})

You can specify custom regular expressions for routes:

	beego.BeeApp.RegisterController("/admin/editpkg/:id([0-9]+)", &admin.EditController{})
	beego.BeeApp.RegisterController("/admin/delpkg/:id([0-9]+)", &admin.DelController{})
	beego.BeeApp.RegisterController("/:pkg(.*)", &controllers.MainController{})
	
You can also create routes for static files:

	beego.BeeApp.SetStaticPath("/static","/public")
	
this will serve any files in /static, including files in subdirectories. For example request `/static/logo.gif` or `/static/style/main.css` will server with the file in the path `/pulic/logo.gif` or `/public/style/main.css`

## Filters / Middleware
============
You can apply filters to routes, which is useful for enforcing security, redirects, etc.

You can, for example, filter all request to enforce some type of security:

	var FilterUser = func(w http.ResponseWriter, r *http.Request) {
	    if r.URL.User == nil || r.URL.User.Username() != "admin" {
	        http.Error(w, "", http.StatusUnauthorized)
	    }
	}
	
	beego.BeeApp.Filter(FilterUser)
	
You can also apply filters only when certain REST URL Parameters exist:

	beego.BeeApp.RegisterController("/:id([0-9]+)", &admin.EditController{})
	beego.BeeApp.FilterParam("id", func(rw http.ResponseWriter, r *http.Request) {
	    ...
	})
	
also You can apply filters only when certain prefix URL path exist:

	beego.BeeApp.FilterPrefixPath("/admin", func(rw http.ResponseWriter, r *http.Request) {
	    … auth 
	})
 		