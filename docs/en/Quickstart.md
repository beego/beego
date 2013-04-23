# Quick start
Hey, you say you've never heard about Beego and don't know how to use it? Don't worry, after you read this section, you will know a lot about Beego. Before you start reading, make sure you installed Beego in your computer, if not, check this tutorial: [Installation](Install.md)

**Navigation**

- [Hello world](#hello-world)
- [New project](#new-project)
- [Development mode](#development-mode)
- [Router](#router)
- [Static files](#static-files)
- [Filter and middleware](#filter-and-middleware)
- [Controller](#controller)
- [Template](#template)
- [Handle request](#handle-request)
- [Redirect and error](#-15)
- [Handle response](#response)
- [Sessions](#sessions)
- [Cache](#cache)
- [Safe map](#map)
- [Log](#-16)
- [Configuration](#-17)
- [Beego arguments](#-18)
- [Integrated third-party applications](#-19)
- [Deployment](#-20)

## Hello world
This is an example of "Hello world" in Beego:

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

Save file as "hello.go", build and run it:

	$ go build main.go
	$ ./hello

Open address [http://127.0.0.1:8080](http://127.0.0.1:8080) in your browser and you will see "hello world".

What happened in behind above example?

1. We import package `github.com/astaxie/beego`. As we know that Go initialize packages and runs init() function in every package(more detail [here](https://github.com/Unknwon/build-web-application-with-golang_EN/blob/master/eBook/02.3.md#main-function-and-init-function)), so Beego initializes the BeeApp application at this time.
2. Define controller. We define a struct called `MainController` with a anonymous field `beego.Controller`, so the `MainController` has all methods that `beego.Controller` has.
3. Define RESTful methods. Once we use anonymous combination, `MainController` has already had `Get`, `Post`, `Delete`, `Put` and other methods, these methods will be called when user sends corresponding request, like `Post` method for requests that are using POST method. Therefore, after we overloaded `Get` method in `MainController`, all GET requests will use `Get` method in `MainController` instead of in `beego.Controller`.
4. Define main function. All applications in Go use main function as entry point as C does.
5. Register routers, it tells Beego which controller is responsibility for specific requests. Here we register `/` for `MainController`, so all requests in `/` will be handed to `MainController`. Be aware that the first argument is the path and the second one is pointer of controller that you want to register.
6. Run application in port 8080 as default, press `Ctrl+c` to exit.

## New project
Get into your $GOPATH, then use following command to setup Beego project:

	bee create hello
	
It generates folders and files for your project, directory structure as follows:

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

## Development mode
Beego uses development mode as default, you can use following code to change mode in your application:

	beego.RunMode = "pro"

Or use configuration file in `conf/app.conf`, and input following content:

	runmode = pro

No differences between two ways.

In development mode, you have following effects:

- If you don't have directory `views`, it prints following error prompt: 

		2013/04/13 19:36:17 [W] [stat views: no such file or directory]

- It doesn't cache template and reload every time.
- If panic occurs in your server, it prints information like following screen shot:

![](images/dev.png)

## Router
The main function of router is to connect request URL and handler. Beego wrapped `Controller`, so it connects request URL and `ControllerInterface`. The `ControllerInterface` has following methods:

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

`beego.Controller` implemented all of them, so you just use this struct as anonymous field in your controller struct. Of course you have to overload corresponding methods for more specific usages.

Users can use following ways to register route rules:

	beego.Router("/", &controllers.MainController{})
	beego.Router("/admin", &admin.UserController{})
	beego.Router("/admin/index", &admin.ArticleController{})
	beego.Router("/admin/addpkg", &admin.AddController{})

For more convenient configure route rules, Beego references the idea from sinatra, so it supports more kinds of route rules as follows:

- beego.Router("/api/:id([0-9]+)", &controllers.RController{})    

		Customized regular expression match 	// match /api/123 :id= 123 

- beego.Router("/news/:all", &controllers.RController{})   
 
		Match rest of all // match /news/path/to/123.html :all= path/to/123.html
	
- beego.Router("/user/:username([\w]+)", &controllers.RController{})   
 
		Regular expression // match /user/astaxie    :username = astaxie
	
- beego.Router("/download/`*`.`*`", &controllers.RController{})   

		Wildcard character // match /download/file/api.xml     :path= file/api   :ext=xml
	
- beego.Router("/download/ceshi/`*`", &controllers.RController{})   

		wildcard character match rest of all // match  /download/ceshi/file/api.json  :splat=file/api.json
	
- beego.Router("/:id:int", &controllers.RController{})   
 
		Match type int  // match :id is int type, Beego uses regular expression ([0-9]+) automatically
	
- beego.Router("/:hi:string", &controllers.RController{})   

		Match type string // match :hi is string type, Beego uses regular expression ([\w]+) automatically

##Static files
Go provides `http.ServeFile` for static files, Beego wrapped this function and use following way to register static file folder:

	beego.SetStaticPath("/static","public")
	
- The first argument is the path of your URL.
- The second argument is the directory in your application path.

Beego supports multiple static file directories as follows:

	beego.SetStaticPath("/images","images")
	beego.SetStaticPath("/css","css")
	beego.SetStaticPath("/js","js")

After you setting static directory, when users visit `/images/login/login.png`，Beego accesses `images/login/login.png` in related to your application directory. One more example, if users visit `/static/img/logo.png`, Beego accesses file `public/img/logo.png`.

##Filter and middleware
Beego supports customized filter and middleware, such as security verification, force redirect, etc.

Here is an example of verify user name of all requests, check if it's admin.

	var FilterUser = func(w http.ResponseWriter, r *http.Request) {
	    if r.URL.User == nil || r.URL.User.Username() != "admin" {
	        http.Error(w, "", http.StatusUnauthorized)
	    }
	}

	beego.Filter(FilterUser)
	
You can also filter by arguments:

	beego.Router("/:id([0-9]+)", &admin.EditController{})
	beego.FilterParam("id", func(rw http.ResponseWriter, r *http.Request) {
	    dosomething()
	})
	
Filter by prefix is also available:

	beego.FilterPrefixPath("/admin", func(rw http.ResponseWriter, r *http.Request) {
	    dosomething()
	})

##Controller
Use `beego.controller` as anonymous in your controller struct to implement the interface in Beego:

	type xxxController struct {
	    beego.Controller
	}

`beego.Controller` implemented`beego.ControllerInterface`, `beego.ControllerInterface` defined following methods:

- Init(ct `*`Context, cn string)

	Initialize context, controller's name, template's name, and container of template arguments
	
- Prepare()

	This is for expend usages, it executes before all the following methods. Users can overload this method for verification for example.
   
- Get()

	This method executes when client sends request as GET method, 403 as default status code. Users overload this method for customized handle process of GET method.
	
- Post()

	This method executes when client sends request as POST method, 403 as default status code. Users overload this method for customized handle process of POST method.

- Delete()

	This method executes when client sends request as DELETE method, 403 as default status code. Users overload this method for customized handle process of DELETE method.

- Put()

	This method executes when client sends request as PUT method, 403 as default status code. Users overload this method for customized handle process of PUT method.

- Head()

	This method executes when client sends request as HEAD method, 403 as default status code. Users overload this method for customized handle process of HEAD method.

- Patch()

	This method executes when client sends request as PATCH method, 403 as default status code. Users overload this method for customized handle process of PATCH method.

- Options()

	This method executes when client sends request as OPTIONS method, 403 as default status code. Users overload this method for customized handle process of OPTIONS method.

- Finish()

	This method executes after corresponding method finished, empty as default. User overload this method for more usages like close database, clean data, etc.

- Render() error

	This method is for rendering template, it executes automatically when you set beego.AutoRender to true.

Overload all methods for all customized logic processes, let's see an example:

	type AddController struct {
	    beego.Controller
	}
	
	func (this *AddController) Prepare() {
	
	}
	
	func (this *AddController) Get() {
		this.Data["content"] ="value"
	    this.Layout = "admin/layout.html"
	    this.TplNames = "admin/add.tpl"
	}
	
	func (this *AddController) Post() {
	    pkgname := this.GetString("pkgname")
	    content := this.GetString("content")
	    pk := models.GetCruPkg(pkgname)
	    if pk.Id == 0 {
	        var pp models.PkgEntity
	        pp.Pid = 0
	        pp.Pathname = pkgname
	        pp.Intro = pkgname
	        models.InsertPkg(pp)
	        pk = models.GetCruPkg(pkgname)
	    }
	    var at models.Article
	    at.Pkgid = pk.Id
	    at.Content = content
	    models.InsertArticle(at)
	    this.Ctx.Redirect(302, "/admin/index")
	}	

##Template
###Template directory
Beego uses `views` as the default directory for template files, parses and caches them as needed(cache is not enable in develop mode), but you can **change**(because only one directory can be used for template files) its directory using following code:

	beego.ViewsPath = "/myviewpath"

###Auto-render
You don't need to call render function manually, Beego calls it automatically after corresponding methods executed. If your application is somehow doesn't need templates, you can disable this feature either in code of `main.go` or configuration file.

To disable auto-render in configuration file:
	
	autorender = false

To disable auto-render in `main.go`(before you call `beego.Run()` to run the application):

	beego.AutoRender = false
	
###Template data
You can use `this.Data` in controller methods to access the data in templates. Suppose you want to get content of `{{.Content}}`, you can use following code to do this:
	
	this.Data["Context"] = "value"
		
###Template name
Beego uses built-in template engine of Go, so there is no different in syntax. As for how to write template file, please visit [Template tutorial](https://github.com/Unknwon/build-web-application-with-golang_EN/blob/master/eBook/07.4.md)。

Beego parses template files in `viewpath` and render it after you set the name of the template file in controller methods. For example, Beego finds the file `add.tpl` in directory `admin` in following code:

	this.TplNames = "admin/add.tpl"

Beego supports two kinds of extensions for template files, which are `tpl` and `html`, if you want to use other extensions, you have to use following code to let Beego know:

	beego.AddTemplateExt("<your template file extension>")

If you enabled auto-render and you don't tell Beego which template file you are going to use in controller methods, Beego uses following format to find the template file if it exists:

	c.TplNames = c.ChildName + "/" + c.Ctx.Request.Method + "." + c.TplExt

Which is `<corresponding controller name>/<request method name>.<template extension>`. For example, your controller name is `AddController` and the request method is POST, and the default file extension is `tpl`, so Beego will try to find file `/<viewpath>/AddController/POST.tpl`.

###Layout design
Beego supports layout design, which means if you are working on an administration application, and some part of its user interface is exactly same all the time, then you can make this part as a layout.
	
	this.Layout = "admin/layout.html"
	this.TplNames = "admin/add.tpl" 

You have to set following variable in order to make Beego possible to insert your dynamic content:

	{{.LayoutContent}}
	
Beego parses template file and assign content to `LayoutContent`, and render them together.

Right now, Beego caches all template files, so you can use following way to implement another kind of layout:

	{{template "header.html"}}
	Handle logic
	{{template "footer.html"}}

###Template function
Beego supports customized template functions that are registered before you call `beego.Run()`.

	func hello(in string)(out string){
	    out = in + "world"
	    return
	}
	
	beego.AddFuncMap("hi",hello)

Then you can use this function in your template files:

	{{.Content | hi}}

There are some built-in template functions:

* markdown 
	
	This function converts markdown content to HTML format, use {{markdown .Content}} in template files.
* dateformat 

	This function converts time to formatted string, use {{dateformat .Time "2006-01-02T15:04:05Z07:00"}} in template files.
* date 

	This function implements date function like in PHP, use formatted string to get corresponding time, use {{date .T "Y-m-d H:i:s"}} in template files.
* compare 

	This functions compares two objects, returns true if they are same, false otherwise, use {{compare .A .B}} in template files.
* substr 

	This function cuts out string from another string by index, it supports UTF-8 characters, use {{substr .Str 0 30}} in template files.
* html2str 

	This function escapes HTML to raw string, use {{html2str .Htmlinfo}} in template files.
* str2html 

	This function outputs string in HTML format without escaping, use {{str2html .Strhtml}} in template files.
* htmlquote 

	This functions implements basic HTML escape, use {{htmlquote .quote}} in template files.
* htmlunquote 	

	This functions implements basic invert-escape of HTML, use {{htmlunquote .unquote}} in template files.
	
##Handle request
We always need to get data from users, including methods like GET, POST, etc. Beego parses these data automatically, and you can access them by following code:

- GetString(key string) string
- GetInt(key string) (int64, error)
- GetBool(key string) (bool, error)

Usage example:

	func (this *MainController) Post() {
		jsoninfo := this.GetString("jsoninfo")
		if jsoninfo == "" {
			this.Ctx.WriteString("jsoninfo is empty")
			return
		}
	}

If you need other types that are not included above, like you need int64 instead of int, then you need to do following way:

	func (this *MainController) Post() {
		id := this.Input().Get("id")
		intid, err := strconv.Atoi(id)
	}		

To use `this.Ctx.Request` for more information about request, and object properties and method please read [Request](http://golang.org/pkg/net/http/#Request)

###File upload
It's very easy to upload file through Beego, but don't forget to add `enctype="multipart/form-data"` in your form, otherwise the browser will not upload anything.

Files will be saved in memory, if the size is greater than cache memory, the rest part will be saved as temporary file. The default cache memory is 64 MB, and you can using following ways to change this size.

In code:

	beego.MaxMemory = 1<<22 

In configuration file:

	maxmemory = 1<<22

Beego provides two convenient functions to upload files:

- GetFile(key string) (multipart.File, `*`multipart.FileHeader, error)

	This function is mainly used to read file name element `the_file` in form and returns corresponding information. You can use this information either filter or save files.
	
- SaveToFile(fromfile, tofile string) error

	This function a wrapper of GetFile and gives ability to save file.
	
This is an example to save file that is uploaded:
	
	func (this *MainController) Post() {
		this.SaveToFile("the_file","/var/www/uploads/uploaded_file.txt"")
	}

###Output Json and XML
Beego considered API function design at the beginning, and we often use Json or XML format data as output. Therefore, it's no reason that Beego doesn't support it:

Set `content-type` to `application/json` for output raw Json format data:

	func (this *AddController) Get() {
	    mystruct := { ... }
	    this.Data["json"] = &mystruct
	    this.ServeJson()
	}	

Set `content-type` to `application/xml` for output raw XML format data:

	func (this *AddController) Get() {
	    mystruct := { ... }
	    this.Data["xml"]=&mystruct
	    this.ServeXml()
	}	
	
## 跳转和错误
我们在做Web开发的时候，经常会遇到页面调整和错误处理，beego这这方面也进行了考虑，通过`Redirect`方法来进行跳转：

	func (this *AddController) Get() {
	   this.Redirect("/", 302)
	}	

@todo 错误处理还需要后期改进

## response处理
response可能会有集中情况：

1. 模板输出

	模板输出上面模板介绍里面已经介绍，beego会在执行完相应的Controller里面的对应的Method之后输出到模板。
	
2. 跳转

	上一节介绍的跳转就是我们经常用到的页面之间的跳转
	
3. 字符串输出

	有些时候我们只是想输出相应的一个字符串，那么我们可以通过如下的代码实现
	
		this.Ctx.WriteString("ok")

## Sessions
beego内置了session模块，目前session模块支持的后端引擎包括memory、file、mysql、redis四中，用户也可以根据相应的interface实现自己的引擎。

beego中使用session相当方便，只要在main入口函数中设置如下：

	beego.SessionOn = true

或者通过配置文件配置如下：

	sessionon = true

通过这种方式就可以开启session，如何使用session，请看下面的例子：

	func (this *MainController) Get() {
		v := this.GetSession("asta")
		if v == nil {
			this.SetSession("asta", int(1))
			this.Data["num"] = 0
		} else {
			this.SetSession("asta", v.(int)+1)
			this.Data["num"] = v.(int)
		}
		this.TplNames = "index.tpl"
	}

上面的例子中我们知道session有几个方便的方法：

- SetSession(name string, value interface{})
- GetSession(name string) interface{}
- DelSession(name string)

session操作主要有设置session、获取session、删除session

当然你要可以通过下面的方式自己控制相应的逻辑这些逻辑：

	sess:=this.StartSession()
	defer sess.SessionRelease()

sess对象具有如下方法：

* sess.Set()
* sess.Get()
* sess.Delete()
* sess.SessionID()

但是我还是建议大家采用SetSession、GetSession、DelSession三个方法来操作，避免自己在操作的过程中资源没释放的问题。

关于Session模块使用中的一些参数设置：

- SessionOn

	设置是否开启Session，默认是false，配置文件对应的参数名：sessionon
	
- SessionProvider

	设置Session的引擎，默认是memory，目前支持还有file、mysql、redis等，配置文件对应的参数名：sessionprovider
	
- SessionName

	设置cookies的名字，Session默认是保存在用户的浏览器cookies里面的，默认名是beegosessionID，配置文件对应的参数名是：sessionname
	
- SessionGCMaxLifetime

	设置Session过期的时间，默认值是3600秒，配置文件对应的参数：sessiongcmaxlifetime
	
- SessionSavePath
	
	设置对应file、mysql、redis引擎的保存路径或者链接地址，默认值是空，配置文件对应的参数：sessionsavepath


当SessionProvider为file时，SessionSavePath是只保存文件的目录，如下所示：

	beego.SessionProvider = "file"
	beego.SessionSavePath = "./tmp"

当SessionProvider为mysql时，SessionSavePath是链接地址，采用[go-sql-driver](https://github.com/go-sql-driver/mysql)，如下所示：

	beego.SessionProvider = "mysql"
	beego.SessionSavePath = "username:password@protocol(address)/dbname?param=value"
	
当SessionProvider为redis时，SessionSavePath是redis的链接地址，采用了[redigo](https://github.com/garyburd/redigo)，如下所示：

	beego.SessionProvider = "redis"
	beego.SessionSavePath = "127.0.0.1:6379"	

## Cache设置
beego内置了一个cache模块，实现了类似memcache的功能，缓存数据在内存中，主要的使用方法如下：

	var (
		urllist *beego.BeeCache
	)
	
	func init() {
		urllist = beego.NewBeeCache()
		urllist.Every = 0 //不过期
		urllist.Start()
	}

	func (this *ShortController) Post() {
		var result ShortResult
		longurl := this.Input().Get("longurl")
		beego.Info(longurl)
		result.UrlLong = longurl
		urlmd5 := models.GetMD5(longurl)
		beego.Info(urlmd5)
		if urllist.IsExist(urlmd5) {
			result.UrlShort = urllist.Get(urlmd5).(string)
		} else {
			result.UrlShort = models.Generate()
			err := urllist.Put(urlmd5, result.UrlShort, 0)
			if err != nil {
				beego.Info(err)
			}
			err = urllist.Put(result.UrlShort, longurl, 0)
			if err != nil {
				beego.Info(err)
			}
		}
		this.Data["json"] = result
		this.ServeJson()
	}	
	
上面这个例子演示了如何使用beego的Cache模块，主要是通过`beego.NewBeeCache`初始化一个对象，然后设置过期时间，开启过期检测，在业务逻辑中就可以通过如下的接口进行增删改的操作：

- Get(name string) interface{}
- Put(name string, value interface{}, expired int) error
- Delete(name string) (ok bool, err error)
- IsExist(name string) bool

## 安全的Map
我们知道在Go语言里面map是非线程安全的，详细的[atomic_maps](http://golang.org/doc/faq#atomic_maps)。但是我们在平常的业务中经常需要用到线程安全的map，特别是在goroutine的情况下，所以beego内置了一个简单的线程安全的map：

	bm := NewBeeMap()
	if !bm.Set("astaxie", 1) {
		t.Error("set Error")
	}
	if !bm.Check("astaxie") {
		t.Error("check err")
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	bm.Delete("astaxie")
	if bm.Check("astaxie") {
		t.Error("delete err")
	}
	
上面演示了如何使用线程安全的Map，主要的接口有：

- Get(k interface{}) interface{}
- Set(k interface{}, v interface{}) bool
- Check(k interface{}) bool
- Delete(k interface{})

## 日志处理
beego默认有一个初始化的BeeLogger对象输出内容到stdout中，你可以通过如下的方式设置自己的输出：

	beego.SetLogger(*log.Logger)

只要你的输出符合`*log.Logger`就可以，例如输出到文件：

	fd,err := os.OpenFile("/var/log/beeapp/beeapp.log", os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
	    beego.Critical("openfile beeapp.log:", err)
	    return
	}
	lg := log.New(fd, "", log.Ldate|log.Ltime)
	beego.SetLogger(lg)
### 不同级别的log日志函数

* Trace(v ...interface{})
* Debug(v ...interface{})
* Info(v ...interface{})
* Warn(v ...interface{})
* Error(v ...interface{})
* Critical(v ...interface{})

你可以通过下面的方式设置不同的日志分级：

	beego.SetLevel(beego.LevelError)
	
当你代码中有很多日志输出之后，如果想上线，但是你不想输出Trace、Debug、Info等信息，那么你可以设置如下：

	beego.SetLevel(beego.LevelWarning)

这样的话就不会输出小于这个level的日志，日志的排序如下：

LevelTrace、LevelDebug、LevelInfo、LevelWarning、	LevelError、LevelCritical	

用户可以根据不同的级别输出不同的错误信息，如下例子所示：

### Examples of log messages
- Trace

	* "Entered parse function validation block"
	* "Validation: entered second 'if'"
	* "Dictionary 'Dict' is empty. Using default value"
- Debug

	* "Web page requested: http://somesite.com Params='...'"
	* "Response generated. Response size: 10000. Sending."
	* "New file received. Type:PNG Size:20000"
- Info
	* "Web server restarted"
	* "Hourly statistics: Requested pages: 12345 Errors: 123 ..."
	* "Service paused. Waiting for 'resume' call"
- Warn
	* "Cache corrupted for file='test.file'. Reading from back-end"
	* "Database 192.168.0.7/DB not responding. Using backup 192.168.0.8/DB"
	* "No response from statistics server. Statistics not sent"
- Error
	* "Internal error. Cannot process request #12345 Error:...."
	* "Cannot perform login: credentials DB not responding"
- Critical
	* "Critical panic received: .... Shutting down"
	* "Fatal error: ... App is shutting down to prevent data corruption or loss"
	
### Example

	func internalCalculationFunc(x, y int) (result int, err error) {
		beego.Debug("calculating z. x:",x," y:",y)
		z := y
		switch {
		case x == 3 :
			beego.Trace("x == 3")
			panic("Failure.")
		case y == 1 :
			beego.Trace("y == 1")
			return 0, errors.New("Error!")
		case y == 2 :
			beego.Trace("y == 2")
			z = x
		default :
			beego.Trace("default")
			z += x
		}
		retVal := z-3
		beego.Debug("Returning ", retVal)
				
		return retVal, nil
	}	
	
	func processInput(input inputData) {
		defer func() {
			if r := recover(); r != nil {
	            beego.Error("Unexpected error occurred: ", r)
				outputs <- outputData{result : 0, error : true}
	        }
		}()
		beego.Info("Received input signal. x:",input.x," y:", input.y)	
		
		res, err := internalCalculationFunc(input.x, input.y)	
		if err != nil {
			beego.Warn("Error in calculation:", err.Error())
		}
		
		beego.Info("Returning result: ",res," error: ",err)		
		outputs <- outputData{result : res, error : err != nil}	
	}
	
	func main() {
		inputs = make(chan inputData)
		outputs = make(chan outputData)
		criticalChan = make(chan int)
		beego.Info("App started.")
		
		go consumeResults(outputs)
		beego.Info("Started receiving results.")
		
		go generateInputs(inputs)
		beego.Info("Started sending signals.")
		
		for {
			select {
				case input := <- inputs:
					processInput(input)
				case <- criticalChan:
					beego.Critical("Caught value from criticalChan: Go shut down.")
					panic("Shut down due to critical fault.")
			}	
		}
	}

## 配置管理
beego支持解析ini文件, beego默认会解析当前应用下的`conf/app.conf`文件

通过这个文件你可以初始化很多beego的默认参数

	appname = beepkg
	httpaddr = "127.0.0.1"
	httpport = 9090
	runmode ="dev"
	autorender = false
	autorecover = false
	viewspath = "myview"
	
上面这些参数会替换beego默认的一些参数。

你可以在配置文件中配置应用需要用的一些配置信息，例如下面所示的数据库信息：

	mysqluser = "root"
	mysqlpass = "rootpass"
	mysqlurls = "127.0.0.1"
	mysqldb   = "beego"
	
那么你就可以通过如下的方式获取设置的配置信息:

	beego.AppConfig.String("mysqluser")
	beego.AppConfig.String("mysqlpass")
	beego.AppConfig.String("mysqlurls")
	beego.AppConfig.String("mysqldb")

AppConfig支持如下方法

- Bool(key string) (bool, error)
- Int(key string) (int, error)
- Int64(key string) (int64, error)
- Float(key string) (float64, error)
- String(key string) string

## 系统默认参数
beego中带有很多可配置的参数，我们来一一认识一下它们，这样有利于我们在接下来的beego开发中可以充分的发挥他们的作用：

* BeeApp

	beego默认启动的一个应用器入口，在应用import beego的时候，在init中已经初始化的。
	
* AppConfig

	beego的配置文件解析之后的对象，也是在init的时候初始化的，里面保存有解析`conf/app.conf`下面所有的参数数据
	
* HttpAddr

	应用监听地址，默认为空，监听所有的网卡IP
	
* HttpPort

	应用监听端口，默认为8080
	
* AppName

	应用名称，默认是beego
	
* RunMode 

	应用的模式，默认是dev，为开发模式，在开发模式下出错会提示友好的出错页面，如前面错误描述中所述。
	
* AutoRender

	是否模板自动渲染，默认值为true，对于API类型的应用，应用需要把该选项设置为false，不需要渲染模板。
	
* RecoverPanic

	是否异常恢复，默认值为true，即当应用出现异常的情况，通过recover恢复回来，而不会导致应用异常退出。
	
* PprofOn

	是否启用pprof，默认是false，当开启之后，用户可以通过如下地址查看相应的goroutine执行情况
	
		/debug/pprof
		/debug/pprof/cmdline
		/debug/pprof/profile
		/debug/pprof/symbol 
	关于pprof的信息，请参考官方的描述[pprof](http://golang.org/pkg/net/http/pprof/)	
	
* ViewsPath

	模板路径，默认值是views
	
* SessionOn

	session是否开启，默认是false
	
* SessionProvider

	session的引擎，默认是memory
	
* SessionName

	存在客户端的cookie名称，默认值是beegosessionID
	
* SessionGCMaxLifetime

	session过期时间，默认值是3600秒
	
* SessionSavePath

	session保存路径，默认是空
	
* UseFcgi

	是否启用fastcgi，默认是false
	
* MaxMemory

	文件上传默认内存缓存大小，默认值是`1 << 26`(64M)

## 第三方应用集成
beego支持第三方应用的集成，用户可以自定义`http.Handler`,用户可以通过如下方式进行注册路由：

	beego.RouterHandler("/chat/:info(.*)", sockjshandler)
	
sockjshandler实现了接口`http.Handler`。

目前在beego的example中有支持sockjs的chat例子，示例代码如下：

	package main
	
	import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/fzzy/sockjs-go/sockjs"
		"strings"
	)
	
	var users *sockjs.SessionPool = sockjs.NewSessionPool()
	
	func chatHandler(s sockjs.Session) {
		users.Add(s)
		defer users.Remove(s)
	
		for {
			m := s.Receive()
			if m == nil {
				break
			}
			fullAddr := s.Info().RemoteAddr
			addr := fullAddr[:strings.LastIndex(fullAddr, ":")]
			m = []byte(fmt.Sprintf("%s: %s", addr, m))
			users.Broadcast(m)
		}
	}
	
	type MainController struct {
		beego.Controller
	}
	
	func (m *MainController) Get() {
		m.TplNames = "index.html"
	}
	
	func main() {
		conf := sockjs.NewConfig()
		sockjshandler := sockjs.NewHandler("/chat", chatHandler, conf)
		beego.Router("/", &MainController{})
		beego.RouterHandler("/chat/:info(.*)", sockjshandler)
		beego.Run()
	}

通过上面的代码很简单的实现了一个多人的聊天室。上面这个只是一个sockjs的例子，我想通过大家自定义`http.Handler`，可以有很多种方式来进行扩展beego应用。

## 部署编译应用
Go语言的应用最后编译之后是一个二进制文件，你只需要copy这个应用到服务器上，运行起来就行。beego由于带有几个静态文件、配置文件、模板文件三个目录，所以用户部署的时候需要同时copy这三个目录到相应的部署应用之下，下面以我实际的应用部署为例：

	$ mkdir /opt/app/beepkg
	$ cp beepkg /opt/app/beepkg
	$ cp -fr views /opt/app/beepkg
	$ cp -fr static /opt/app/beepkg
	$ cp -fr conf /opt/app/beepkg
	
这样在`/opt/app/beepkg`目录下面就会显示如下的目录结构：

	.
	├── conf
	│   ├── app.conf
	├── static
	│   ├── css
	│   ├── img
	│   └── js
	└── views
	    └── index.tpl
	├── beepkg	

这样我们就已经把我们需要的应用搬到服务器了，那么接下来就可以开始部署了，我现在服务器端用两种方式来run，

- Supervisord 
	
	安装和配置见[Supervisord](Supervisord.md)

- nohup方式

	nohup ./beepkg &

个人比较推荐第一种方式，可以很好的管理起来应用

- [Introduction](README.md)
- [Step by step](Tutorial.md)