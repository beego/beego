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

package orm

import (
	"fmt"
	"runtime/debug"

	imodels "github.com/beego/beego/v2/client/orm/internal/models"
)

var defaultModelCache = imodels.NewModelCacheHandler()

// RegisterModel Register models
func RegisterModel(models ...interface{}) {
	RegisterModelWithPrefix("", models...)
}

// RegisterModelWithPrefix Register models with a prefix
func RegisterModelWithPrefix(prefix string, models ...interface{}) {
	if err := defaultModelCache.Register(prefix, true, models...); err != nil {
		panic(err)
	}
}

// RegisterModelWithSuffix Register models with a suffix
func RegisterModelWithSuffix(suffix string, models ...interface{}) {
	if err := defaultModelCache.Register(suffix, false, models...); err != nil {
		panic(err)
	}
}

// BootStrap Bootstrap models.
// make All model parsed and can not add more models
func BootStrap() {
	BootStrapWithAlias("default")
}

// BootStrap with alias
func BootStrapWithAlias(alias string) {
	if _, ok := dataBaseCache.get(alias); !ok {
		fmt.Printf("must have one Register DataBase alias named %q\n", alias)
		debug.PrintStack()
		return
	}
	defaultModelCache.Bootstrap()
}

// ResetModelCache Clean model cache. Then you can re-RegisterModel.
// Common use this api for test case.
func ResetModelCache() {
	defaultModelCache.Clean()
}
