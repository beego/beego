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

package beego

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// beego web framework version.
const VERSION = "1.5.0"

//hook function to run
type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) //hook function slice to store the hookfunc
)

// AddAPPStartHook is used to register the hookfunc
// The hookfuncs will run in beego.Run()
// such as sessionInit, middlerware start, buildtemplate, admin start
func AddAPPStartHook(hf hookfunc) {
	hooks = append(hooks, hf)
}

// Run beego application.
// beego.Run() default run on HttpPort
// beego.Run(":8089")
// beego.Run("127.0.0.1:8089")
func Run(params ...string) {
	initBeforeHTTPRun()

	params = append(params, fmt.Sprintf("%s:%d", HTTPAddr, HTTPPort))
	addr := strings.Split(params[0], ":")
	if len(addr) == 2 {
		HTTPAddr = addr[0]
		HTTPPort, _ = strconv.Atoi(addr[1])
	}

	BeeApp.Run()
}

func initBeforeHTTPRun() {
	// if AppConfigPath not In the conf/app.conf reParse config
	if AppConfigPath != filepath.Join(AppPath, "conf", "app.conf") {
		err := ParseConfig()
		if err != nil && AppConfigPath != filepath.Join(workPath, "conf", "app.conf") {
			// configuration is critical to app, panic here if parse failed
			panic(err)
		}
	}

	//init hooks
	AddAPPStartHook(registerMime)
	AddAPPStartHook(registerDefaultErrorHandler)
	AddAPPStartHook(registerSession)
	AddAPPStartHook(registerDocs)
	AddAPPStartHook(registerTemplate)
	AddAPPStartHook(registerAdmin)

	for _, hk := range hooks {
		if err := hk(); err != nil {
			panic(err)
		}
	}
}

// TestBeegoInit is for test package init
func TestBeegoInit(apppath string) {
	AppPath = apppath
	os.Setenv("BEEGO_RUNMODE", "test")
	AppConfigPath = filepath.Join(AppPath, "conf", "app.conf")
	err := ParseConfig()
	if err != nil && !os.IsNotExist(err) {
		// for init if doesn't have app.conf will not panic
		Info(err)
	}
	os.Chdir(AppPath)
	initBeforeHTTPRun()
}
