# Beego [![Test](https://github.com/beego/beego/actions/workflows/test.yml/badge.svg?branch=develop)](https://github.com/beego/beego/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/beego/beego)](https://goreportcard.com/report/github.com/beego/beego) [![Go Reference](https://pkg.go.dev/badge/github.com/beego/beego/v2.svg)](https://pkg.go.dev/github.com/beego/beego/v2)

Beego is used for rapid development of enterprise application in Go, including RESTful APIs, web apps and backend services.

It is inspired by Tornado, Sinatra and Flask. beego has some Go-specific features such as interfaces and struct embedding.

## Quick Start
- [New Doc Website - unavailable](https://beego.gocn.vip)
- [New Doc Website Backup @flycash](https://doc.meoying.com/en-US/beego/developing/)
- [New Doc Website source code](https://github.com/beego/beego-doc)
- [Old Doc - github](https://github.com/beego/beedoc)
- [Example](https://github.com/beego/beego-example)

> Kindly remind that sometimes the HTTPS certificate is expired, you may get some NOT SECURE warning

### Web Application

#### Create `hello` directory, cd `hello` directory

    mkdir hello
    cd hello

#### Init module

    go mod init

#### Download and install

    go get github.com/beego/beego/v2@latest

#### Create file `hello.go`

```go
package main

import "github.com/beego/beego/v2/server/web"

func main() {
	web.Run()
}
```

#### Download required dependencies

    go mod tidy

#### Build and run

    go build hello.go
    ./hello

#### Go to [http://localhost:8080](http://localhost:8080)

Congratulations! You've just built your first **beego** app.

## Features

* RESTful support
* [MVC architecture](https://github.com/beego/beedoc/tree/master/en-US/mvc)
* Modularity
* [Auto API documents](https://github.com/beego/beedoc/blob/master/en-US/advantage/docs.md)
* [Annotation router](https://github.com/beego/beedoc/blob/master/en-US/mvc/controller/router.md)
* [Namespace](https://github.com/beego/beedoc/blob/master/en-US/mvc/controller/router.md#namespace)
* [Powerful development tools](https://github.com/beego/bee)
* Full stack for Web & API

## Modules

* [orm](https://github.com/beego/beedoc/tree/master/en-US/mvc/model)
* [session](https://github.com/beego/beedoc/blob/master/en-US/module/session.md)
* [logs](https://github.com/beego/beedoc/blob/master/en-US/module/logs.md)
* [config](https://github.com/beego/beedoc/blob/master/en-US/module/config.md)
* [cache](https://github.com/beego/beedoc/blob/master/en-US/module/cache.md)
* [context](https://github.com/beego/beedoc/blob/master/en-US/module/context.md)
* [admin](https://github.com/beego/beedoc/blob/master/en-US/module/admin.md)
* [httplib](https://github.com/beego/beedoc/blob/master/en-US/module/httplib.md)
* [task](https://github.com/beego/beedoc/blob/master/en-US/module/task.md)
* [i18n](https://github.com/beego/beedoc/blob/master/en-US/module/i18n.md)

## Community

* Welcome to join us in Slack: [https://beego.slack.com invite](https://join.slack.com/t/beego/shared_invite/zt-fqlfjaxs-_CRmiITCSbEqQG9NeBqXKA),
* QQ Group ID:523992905
* [Contribution Guide](https://github.com/beego/beedoc/blob/master/en-US/intro/contributing.md).

## License

beego source code is licensed under the Apache Licence, Version 2.0
([https://www.apache.org/licenses/LICENSE-2.0.html](https://www.apache.org/licenses/LICENSE-2.0.html)).
