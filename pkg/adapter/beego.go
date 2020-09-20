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

package adapter

import (
	"github.com/astaxie/beego/pkg"
	"github.com/astaxie/beego/pkg/server/web"
)

const (

	// VERSION represent beego web framework version.
	VERSION = pkg.VERSION

	// DEV is for develop
	DEV = web.DEV
	// PROD is for production
	PROD = web.PROD
)

// M is Map shortcut
type M web.M

// Hook function to run
type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) // hook function slice to store the hookfunc
)

// AddAPPStartHook is used to register the hookfunc
// The hookfuncs will run in beego.Run()
// such as initiating session , starting middleware , building template, starting admin control and so on.
func AddAPPStartHook(hf ...hookfunc) {
	for _, f := range hf {
		web.AddAPPStartHook(func() error {
			return f()
		})
	}
}

// Run beego application.
// beego.Run() default run on HttpPort
// beego.Run("localhost")
// beego.Run(":8089")
// beego.Run("127.0.0.1:8089")
func Run(params ...string) {
	web.Run(params...)
}

// RunWithMiddleWares Run beego application with middlewares.
func RunWithMiddleWares(addr string, mws ...MiddleWare) {
	newMws := oldMiddlewareToNew(mws)
	web.RunWithMiddleWares(addr, newMws...)
}

// TestBeegoInit is for test package init
func TestBeegoInit(ap string) {
	web.TestBeegoInit(ap)
}

// InitBeegoBeforeTest is for test package init
func InitBeegoBeforeTest(appConfigPath string) {
	web.InitBeegoBeforeTest(appConfigPath)
}
