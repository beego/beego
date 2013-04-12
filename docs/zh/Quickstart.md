# 快速入门
你对beego一无所知？没关系，这篇文档会很好的详细介绍beego的各个方面，看这个文档之前首先确认你已经安装了beego，如果你没有安装的话，请看这篇[安装指南](Install.md)

##导航

- [最小应用](#-1)
- [新建项目](#-2)
- [开发模式](#-3)
- [路由设置](#-4)
- [静态文件](#-5)
- [模板处理](#-6)
- [request处理](#-7)
- [跳转和错误](#-8)
- [response处理](#-9)
- [Sessions](#-10)
- [Cache设置](#-11)
- [安全的Map](#-12)
- [日志处理](#-13)
- [第三方应用集成](#-14)
- [部署编译应用](#-15)

## 最小应用
一个最小最简单的应用如下代码所示：

	package main

	import (
		"github.com/astaxie/beego"
	)
	
	type MainController struct {
		beego.Controller
	}
	
	func (this *MainController) Get() {
		this.Ctx.WriteString("hello world")
	}
	
	func main() {
		beego.Router("/", &MainController{})
		beego.Run()
	}

把上面的代码保存为hello.go，然后通过命令行进行编译并执行：

	$ go build main.go
	$ ./hello

这个时候你可以打开你的浏览器，通过这个地址浏览[http://127.0.0.1:8080](http://127.0.0.1:8080)返回“hello world”

那么上面的代码到底做了些什么呢？

1、首先我们引入了包`github.com/astaxie/beego`,我们知道Go语言里面引入包会深度优先的去执行引入包的初始化(变量和init函数，[更多](https://github.com/astaxie/build-web-application-with-golang/blob/master/ebook/02.3.md#maininit))，beego包中会初始化一个BeeAPP的应用，初始化一些参数。

2、定义Controller，这里我们定义了一个struct为`MainController`，充分利用了Go语言的组合的概念，匿名包含了`beego.Controller`，这样我们的`MainController`就拥有了`beego.Controller`的所有方法。

3、定义RESTFul方法，通过匿名组合之后，其实目前的`MainController`已经拥有了`Get`、`Post`、`Delete`、`Put`等方法，这些方法是分别用来对应用户请求的Method函数，如果用户发起的是`POST`请求，那么就执行`Post`函数。所以这里我们定义了`MainController`的`Get`方法用来重写继承的`Get`函数，这样当用户`GET`请求的时候就会执行该函数。
				
4、定义main函数，所有的Go应用程序和C语言一样都是Main函数作为入口，所以我们这里定义了我们应用的入口。

5、Router注册路由，路由就是告诉beego，当用户来请求的时候，该如何去调用相应的Controller，这里我们注册了请求`/`的时候，请求到`MainController`。这里我们需要知道，Router函数的两个参数函数，第一个是路径，第二个是Controller的指针。

6、Run应用，最后一步就是把在1中初始化的BeeApp开启起来，其实就是内部监听了8080端口:`Go默认情况会监听你本机所有的IP上面的8080端口`

停止服务的话，请按`ctrl+c`

## 新建项目

通过如下命令创建beego项目，首先进入gopath目录

	bee create hello
	
这样就建立了一个项目hello，目录结构如下所示

	.
	├── conf
	│   └── app.conf
	├── controllers
	│   └── default.go
	├── main.go
	├── models
	├── static
	│   ├── css
	│   ├── img
	│   └── js
	└── views
	    └── index.tpl

## 开发模式

## 路由设置

## 静态文件

## 模板处理

## request处理

## 跳转和错误

## response处理

## Sessions

## Cache设置

## 安全的Map

## 日志处理

## 第三方应用集成

## 部署编译应用