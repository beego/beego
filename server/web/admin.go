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

package web

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

// BeeAdminApp is the default adminApp used by admin module.
var beeAdminApp *adminApp

// FilterMonitorFunc is default monitor filter when admin module is enable.
// if this func returns, admin module records qps for this request by condition of this function logic.
// usage:
// 	func MyFilterMonitor(method, requestPath string, t time.Duration, pattern string, statusCode int) bool {
//	 	if method == "POST" {
//			return false
//	 	}
//	 	if t.Nanoseconds() < 100 {
//			return false
//	 	}
//	 	if strings.HasPrefix(requestPath, "/astaxie") {
//			return false
//	 	}
//	 	return true
// 	}
// 	beego.FilterMonitorFunc = MyFilterMonitor.
var FilterMonitorFunc func(string, string, time.Duration, string, int) bool

func init() {
	FilterMonitorFunc = func(string, string, time.Duration, string, int) bool { return true }
}

func list(root string, p interface{}, m M) {
	pt := reflect.TypeOf(p)
	pv := reflect.ValueOf(p)
	if pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
		pv = pv.Elem()
	}
	for i := 0; i < pv.NumField(); i++ {
		var key string
		if root == "" {
			key = pt.Field(i).Name
		} else {
			key = root + "." + pt.Field(i).Name
		}
		if pv.Field(i).Kind() == reflect.Struct {
			list(key, pv.Field(i).Interface(), m)
		} else {
			m[key] = pv.Field(i).Interface()
		}
	}
}

func writeJSON(rw http.ResponseWriter, jsonData []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonData)
}

// adminApp is an http.HandlerFunc map used as beeAdminApp.
type adminApp struct {
	*HttpServer
}

// Run start Beego admin
func (admin *adminApp) Run() {
	logs.Debug("now we don't start tasks here, if you use task module," +
		" please invoke task.StartTask, or task will not be executed")
	addr := BConfig.Listen.AdminAddr
	if BConfig.Listen.AdminPort != 0 {
		addr = net.JoinHostPort(BConfig.Listen.AdminAddr, fmt.Sprintf("%d", BConfig.Listen.AdminPort))
	}
	logs.Info("Admin server Running on %s", addr)
	admin.HttpServer.Run(addr)
}

func registerAdmin() error {
	if BConfig.Listen.EnableAdmin {

		c := &adminController{
			servers: make([]*HttpServer, 0, 2),
		}

		// copy config to avoid conflict
		adminCfg := *BConfig
		adminCfg.Listen.EnableHTTPS = false
		adminCfg.Listen.EnableMutualHTTPS = false
		beeAdminApp = &adminApp{
			HttpServer: NewHttpServerWithCfg(&adminCfg),
		}
		// keep in mind that all data should be html escaped to avoid XSS attack
		beeAdminApp.Router("/", c, "get:AdminIndex")
		beeAdminApp.Router("/qps", c, "get:QpsIndex")
		beeAdminApp.Router("/prof", c, "get:ProfIndex")
		beeAdminApp.Router("/healthcheck", c, "get:Healthcheck")
		beeAdminApp.Router("/task", c, "get:TaskStatus")
		beeAdminApp.Router("/listconf", c, "get:ListConf")
		beeAdminApp.Router("/metrics", c, "get:PrometheusMetrics")

		go beeAdminApp.Run()
	}
	return nil
}
