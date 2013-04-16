# 快速入门
你对beego一无所知？没关系，这篇文档会很好的详细介绍beego的各个方面，看这个文档之前首先确认你已经安装了beego，如果你没有安装的话，请看这篇[安装指南](Install.md)

**导航**

- [最小应用](#-1)
- [新建项目](#-2)
- [开发模式](#-3)
- [路由设置](#-4)
- [静态文件](#-5)
- [过滤和中间件](#-6)
- [模板处理](#-7)
- [request处理](#-8)
- [跳转和错误](#-9)
- [response处理](#-10)
- [Sessions](#-11)
- [Cache设置](#-12)
- [安全的Map](#-13)
- [日志处理](#-14)
- [第三方应用集成](#-15)
- [部署编译应用](#-16)

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

通过bee创建的项目，beego默认情况下是开发模式。
	
我们可以通过如下的方式改变我们的模式：

	beego.RunMode = "pro"

或者我们在conf/app.conf下面设置如下：

	runmode = pro

以上两种效果一样。

开发模式中

- 开发模式下，如果你的目录不存在views目录，那么会出现类似下面的错误提示：

		2013/04/13 19:36:17 [W] [stat views: no such file or directory]

- 模板会自动重新加载不缓存。
- 如果服务端出错，那么就会在浏览器端显示如下类似的截图：

![](images/dev.png)

## 路由设置

路由的主要功能是实现从请求地址到实现方法，beego中封装了`Controller`，所以路由是从路径到`ControllerInterface`的过程，`ControllerInterface`的方法有如下：

	type ControllerInterface interface {
		Init(ct *Context, cn string)
		Prepare()
		Get()
		Post()
		Delete()
		Put()
		Head()
		Patch()
		Options()
		Finish()
		Render() error
	}

这些方法`beego.Controller`都已经实现了，所以只要用户定义struct的时候匿名包含就可以了。当然更灵活的方法就是用户可以去自定义类似的方法，然后实现自己的逻辑。

用户可以通过如下的方式进行路由设置：

	beego.Router("/", &controllers.MainController{})
	beego.Router("/admin", &admin.UserController{})
	beego.Router("/admin/index", &admin.ArticleController{})
	beego.Router("/admin/addpkg", &admin.AddController{})

为了用户更加方便的路由设置，beego参考了sinatra的路由实现，支持多种方式的路由：

- beego.Router("/api/:id([0-9]+)", &controllers.RController{})    
	自定义正则匹配	//匹配 /api/123 :id= 123 

- beego.Router("/news/:all", &controllers.RController{})    
	全匹配方式 //匹配 /news/path/to/123.html :all= path/to/123.html
	
- beego.Router("/user/:username([\w]+)", &controllers.RController{})    
	正则字符串匹配 //匹配 /user/astaxie    :username = astaxie
	
- beego.Router("/download/*.*", &controllers.RController{})    
	*匹配方式 //匹配 /download/file/api.xml     :path= file/api   :ext=xml
	
- beego.Router("/download/ceshi/*", &controllers.RController{})   
	*全匹配方式 //匹配  /download/ceshi/file/api.json  :splat=file/api.json
	
- beego.Router("/:id:int", &controllers.RController{})    
	int类型设置方式  //匹配 :id为int类型，框架帮你实现了正则([0-9]+)
	
- beego.Router("/:hi:string", &controllers.RController{})   
	string类型设置方式 //匹配 :hi为string类型。框架帮你实现了正则([\w]+)

## 静态文件
Go语言内部其实已经提供了`http.ServeFile`，通过这个函数可以实现静态文件的服务。beego针对这个功能进行了一层封装，通过下面的方式进行静态文件注册：

	beego.SetStaticPath("/static","public")
	
- 第一个参数是路径，url路径信息
- 第二个参数是静态文件目录（相对应用所在的目录）

beego支持多个目录的静态文件注册，用户可以注册如下的静态文件目录：

	beego.SetStaticPath("/images","images")
	beego.SetStaticPath("/css","css")
	beego.SetStaticPath("/js","js")

设置了如上的静态目录之后，用户访问`/images/login/login.png`，那么就会访问应用对应的目录下面的`images/login/login.png`文件。如果是访问`/static/img/logo.png`，那么就访问`public/img/logo.png`文件。

## 过滤和中间件
beego支持自定义过滤中间件，例如安全验证，强制跳转等

如下例子所示，验证用户名是否是admin，应用于全部的请求：

	var FilterUser = func(w http.ResponseWriter, r *http.Request) {
	    if r.URL.User == nil || r.URL.User.Username() != "admin" {
	        http.Error(w, "", http.StatusUnauthorized)
	    }
	}

	beego.Filter(FilterUser)
	
还可以通过参数进行过滤，如果匹配参数就执行

	beego.Router("/:id([0-9]+)", &admin.EditController{})
	beego.FilterParam("id", func(rw http.ResponseWriter, r *http.Request) {
	    dosomething()
	})
	
当然你还可以通过前缀过滤

	beego.FilterPrefixPath("/admin", func(rw http.ResponseWriter, r *http.Request) {
	    dosomething()
	})

## 模板处理
### 模板目录
beego中默认的模板目录是`views`，用户可以把你的模板文件放到该目录下，beego会自动在该目录下的所有模板文件进行解析并缓存，开发模式下会每次重新解析，不做缓存。当然用户可以通过如下的方式改变模板的目录：

	beego.ViewsPath = "/myviewpath"
### 自动渲染
beego中用户无需手动的调用渲染输出模板，beego会自动的在调用玩相应的method方法之后调用Render函数，当然如果你的应用是不需要模板输出的，那么你可以在配置文件或者在main.go中设置关闭自动渲染。

配置文件配置如下：
	
	autorender = false

main.go文件中设置如下：

	beego.AutoRender = false
### 模板名称
beego采用了Go语言内置的模板引擎，所有模板的语法和Go的一模一样，至于如何写模板文件，详细的请参考[模板教程](https://github.com/astaxie/build-web-application-with-golang/blob/master/ebook/07.4.md)。

用户通过在Controller的对应方法中设置相应的模板名称，beego会自动的在viewpath目录下查询该文件并渲染，例如下面的设置，beego会在admin下面找add.tpl文件进行渲染：

	this.TplNames = "admin/add.tpl"

我们看到上面的模板后缀名是tpl，beego默认情况下支持tpl和html后缀名的模板文件，如果你的后缀名不是这两种，请进行如下设置：

	beego.AddTemplateExt("你文件的后缀名")

当你设置了自动渲染，然后在你的Controller中没有设置任何的TplNames，那么beego会自动设置你的模板文件如下：

	c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt

也就是你对应的Controller名字+请求方法名.模板后缀，也就是如果你的Controller名是`AddController`，请求方法是`POST`，默认的文件后缀是`tpl`，那么就会默认请求`/viewpath/AddController/POST.tpl`文件。

### lauout设计
beego支持layout设计，例如你在管理系统中，其实整个的管理界面是固定的，支会变化中间的部分，那么你可以通过如下的设置：
	
	this.Layout = "admin/layout.html"
	this.TplNames = "admin/add.tpl" 

在layout.html中你必须设置如下的变量：

	{{.LayoutContent}}
	
beego就会首先解析TplNames指定的文件，获取内容赋值给LayoutContent，然后最后渲染layout.html文件。

目前采用首先把目录下所有的文件进行缓存，所以用户还可以通过类似这样的方式实现layout：

	{{template "header.html"}}
	处理逻辑
	{{template "footer.html"}}

### 模板函数
beego支持用户定义模板函数，但是必须在`beego.Run()`调用之前，设置如下：

	func hello(in string)(out string){
	    out = in + "world"
	    return
	}
	
	beego.AddFuncMap("hi",hello)

定义之后你就可以在模板中这样使用了：

	{{.Content | hi}}

目前beego内置的模板函数有如下：

* markdown 
	
	实现了把markdown文本转化为html信息，使用方法{{markdown .Content}}
* dateformat 

	实现了时间的格式化，返回字符串，使用方法{{dateformat .Time "2006-01-02T15:04:05Z07:00"}}
* date 

	实现了类似PHP的date函数，可以很方便的根据字符串返回时间，使用方法{{date .T "Y-m-d H:i:s"}}
* compare 

	实现了比较两个对象的比较，如果相同返回true，否者false，使用方法{{compare .A .B}}
* substr 

	实现了字符串的截取，支持中文截取的完美截取，使用方法{{substr .Str 0 30}}
* html2str 

	实现了把html转化为字符串，剔除一些script、css之类的元素，返回纯文本信息，使用方法{{html2str .Htmlinfo}}
* str2html 

	实现了把相应的字符串当作HTML来输出，不转义，使用方法{{str2html .Strhtml}}
* htmlquote 

	实现了基本的html字符转义，使用方法{{htmlquote .quote}}
* htmlunquote 	

	实现了基本的反转移字符，使用方法{{htmlunquote .unquote}}
	
## request处理
我们经常需要获取用户传递的数据，包括Get、POST等方式的请求，beego里面会自动解析这些数据，你可以通过如下方式获取数据

- GetString
- GetInt
- GetBool


### 文件上传

- GetFile
- SaveToFile

## 跳转和错误

## response处理

## Sessions

## Cache设置

## 安全的Map

## 日志处理

## 第三方应用集成

## 部署编译应用

- [beego介绍](README.md)
- [一步一步开发应用](Tutorial.md)