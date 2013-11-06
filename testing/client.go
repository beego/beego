package testing

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/httplib"
)

var port = ""
var baseUrl = "http://localhost:"

type TestHttpRequest struct {
	httplib.BeegoHttpRequest
}

func getPort() string {
	if port == "" {
		config, err := config.NewConfig("ini", "../conf/app.conf")
		if err != nil {
			return "8080"
		}
		port = config.String("httpport")
		return port
	}
	return port
}

func Get(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Get(baseUrl + getPort() + path)}
}

func Post(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Post(baseUrl + getPort() + path)}
}

func Put(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Put(baseUrl + getPort() + path)}
}

func Delete(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Delete(baseUrl + getPort() + path)}
}

func Head(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Head(baseUrl + getPort() + path)}
}
