# 快速入门
你对beego一无所知？没关系，这篇文档会很好的详细介绍beego的各个方面，看这个文档之前首先确认你已经安装了beego，如果你没有安装的话，请看这篇[安装指南](Install.md)

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

1、首先我们引入了包`github.com/astaxie/beego`,beego包中会初始化一个BeeAPP的应用，

2、定义了Controller

3、定义了RESTFul方法
				
4、定义了main函数

5、Route注册路由

6、Run应用

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