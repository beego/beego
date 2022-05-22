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

// RegisterModel register models
func RegisterModel(models ...interface{}) {
	RegisterModelWithPrefix("", models...)
}

// RegisterModelWithPrefix register models with a prefix
func RegisterModelWithPrefix(prefix string, models ...interface{}) {
	if err := defaultModelCache.register(prefix, true, models...); err != nil {
		panic(err)
	}
}

// RegisterModelWithSuffix register models with a suffix
func RegisterModelWithSuffix(suffix string, models ...interface{}) {
	if err := defaultModelCache.register(suffix, false, models...); err != nil {
		panic(err)
	}
}

// BootStrap bootstrap models.
// make all model parsed and can not add more models
func BootStrap() {
	defaultModelCache.bootstrap()
}
