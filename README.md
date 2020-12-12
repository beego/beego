# Beego [![Build Status](https://travis-ci.org/astaxie/beego.svg?branch=master)](https://travis-ci.org/astaxie/beego) [![GoDoc](http://godoc.org/github.com/astaxie/beego?status.svg)](http://godoc.org/github.com/astaxie/beego) [![Foundation](https://img.shields.io/badge/Golang-Foundation-green.svg)](http://golangfoundation.org) [![Go Report Card](https://goreportcard.com/badge/github.com/astaxie/beego)](https://goreportcard.com/report/github.com/astaxie/beego)


beego is used for rapid development of RESTful APIs, web apps and backend services in Go.
It is inspired by Tornado, Sinatra and Flask. beego has some Go-specific features such as interfaces and struct embedding.

###### More info at [beego.me](http://beego.me). 

> If you could not open this website, go to [beedoc](https://github.com/beego/beedoc)

## beego 1.x and 2.x

We recently release beego 2.0.0-beta, and its structure change a lot, so you may get some error

1. If you are working on beego v1.x please try `go get github.com/astaxie/beego@v1.12.3`
2. If you want to try beego 2.0.0, run `go get github.com/astaxie/beego@develop`

We are still working on fix bug and documentation of v2.x. And v2.x's doc will be released with v2.0.0.

## Next version

v1.12.4 will be released on Jan 2021 And v2.0.0 will be released next month.

## Quick Start

###### Please see [Documentation](http://beego.me/docs) for more.

###### [beego-example](https://github.com/beego-dev/beego-example)

### Web Application

#### Create `hello` directory, cd `hello` directory

    mkdir hello
    cd hello
 
#### Init module

    go mod init

#### Download and install

    go get -u github.com/astaxie/beego@develop

Now we are working on beego v2.0.0, so using `@develop`.

#### Create file `hello.go`
```go
package main

import "github.com/astaxie/beego/server/web"

func main(){
    web.Run()
}
```
#### Build and run

    go build hello.go
    ./hello

#### Go to [http://localhost:8080](http://localhost:8080)

Congratulations! You've just built your first **beego** app.

### Using ORM module

```go

package main

import (
	"github.com/astaxie/beego/client/orm"
	"github.com/astaxie/beego/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

// User -
type User struct {
	ID   int    `orm:"column(id)"`
	Name string `orm:"column(name)"`
}

func init() {
	// need to register models in init
	orm.RegisterModel(new(User))

	// need to register db driver
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// need to register default database
	orm.RegisterDataBase("default", "mysql", "beego:test@tcp(192.168.0.105:13306)/orm_test?charset=utf8")
}

func main() {
	// automatically build table
	orm.RunSyncdb("default", false, true)

	// create orm object, and it will use `default` database
	o := orm.NewOrm()

	// data
	user := new(User)
	user.Name = "mike"

	// insert data
	id, err := o.Insert(user)
	if err != nil {
		logs.Info(err)
	}
	
	// ...
}
```

### Using httplib as http client
```go
package main

import (
	"github.com/astaxie/beego/client/httplib"
	"github.com/astaxie/beego/core/logs"
)

func main() {
	// Get, more methods please read docs
	req := httplib.Get("http://beego.me/")
	str, err := req.String()
	if err != nil {
		logs.Error(err)
	}
	logs.Info(str)
}

```

### Using config module

```go
package main

import (
	"context"

	"github.com/astaxie/beego/core/config"
	"github.com/astaxie/beego/core/logs"
)

var (
	ConfigFile = "./app.conf"
)

func main() {
	cfg, err := config.NewConfig("ini", ConfigFile)
	if err != nil {
		logs.Critical("An error occurred:", err)
		panic(err)
	}
	res, _ := cfg.String(context.Background(), "name")
	logs.Info("load config name is", res)
}
```
### Using logs module
```go
package main

import (
	"github.com/astaxie/beego/core/logs"
)

func main() {
	err := logs.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	if err != nil {
		panic(err)
	}
	logs.Info("hello beego")
}
```
### Using timed task

```go
package main

import (
	"context"
	"time"

	"github.com/astaxie/beego/core/logs"
	"github.com/astaxie/beego/task"
)

func main() {
	// create a task
	tk1 := task.NewTask("tk1", "0/3 * * * * *", func(ctx context.Context) error { logs.Info("tk1"); return nil })

	// check task
	err := tk1.Run(context.Background())
	if err != nil {
		logs.Error(err)
	}

	// add task to global todolist
	task.AddTask("tk1", tk1)

	// start tasks
	task.StartTask()

	// wait 12 second
	time.Sleep(12 * time.Second)
	defer task.StopTask()
}
``` 

### Using cache module

```go
package main

import (
	"context"
	"time"

	"github.com/astaxie/beego/client/cache"

	// don't forget this
	_ "github.com/astaxie/beego/client/cache/redis"

	"github.com/astaxie/beego/core/logs"
)

func main() {
	// create cache
	bm, err := cache.NewCache("redis", `{"key":"default", "conn":":6379", "password":"123456", "dbNum":"0"}`)
	if err != nil {
		logs.Error(err)
	}

	// put
	isPut := bm.Put(context.Background(), "astaxie", 1, time.Second*10)
	logs.Info(isPut)

	isPut = bm.Put(context.Background(), "hello", "world", time.Second*10)
	logs.Info(isPut)

	// get
	result, _ := bm.Get(context.Background(),"astaxie")
	logs.Info(string(result.([]byte)))

	multiResult, _ := bm.GetMulti(context.Background(), []string{"astaxie", "hello"})
	for i := range multiResult {
		logs.Info(string(multiResult[i].([]byte)))
	}

	// isExist
	isExist, _ := bm.IsExist(context.Background(), "astaxie")
	logs.Info(isExist)

	// delete
	isDelete := bm.Delete(context.Background(), "astaxie")
	logs.Info(isDelete)
}
```


## Features

* RESTful support
* MVC architecture
* Modularity
* Auto API documents
* Annotation router
* Namespace
* Powerful development tools
* Full stack for Web & API

## Documentation

* [English](http://beego.me/docs/intro/)
* [中文文档](http://beego.me/docs/intro/)
* [Русский](http://beego.me/docs/intro/)

## Community

* [http://beego.me/community](http://beego.me/community)
* Welcome to join us in Slack: [https://beego.slack.com](https://beego.slack.com), you can get invited from [here](https://github.com/beego/beedoc/issues/232)
* QQ Group Group ID:523992905

## License

beego source code is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
