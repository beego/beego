// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/httplib"
)

var port = ""
var baseUrl = "http://localhost:"

// beego test request client
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

// returns test client in GET method
func Get(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Get(baseUrl + getPort() + path)}
}

// returns test client in POST method
func Post(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Post(baseUrl + getPort() + path)}
}

// returns test client in PUT method
func Put(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Put(baseUrl + getPort() + path)}
}

// returns test client in DELETE method
func Delete(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Delete(baseUrl + getPort() + path)}
}

// returns test client in HEAD method
func Head(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Head(baseUrl + getPort() + path)}
}
