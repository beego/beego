Beego
=======
Beego is an open source version of the scalable, non-blocking web server
and tools that power SNDA's CDN system. Documentation and downloads are
available at https://github.com/astaxie/beego

Beego is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).

Installation
============
To install:

    go install github.com/astaxie/beego

go version: go1 release

Quick Start
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
