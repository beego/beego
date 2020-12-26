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

package cache

import (
	"github.com/beego/beego/v2/client/cache"
)

// GetString convert interface to string.
func GetString(v interface{}) string {
	return cache.GetString(v)
}

// GetInt convert interface to int.
func GetInt(v interface{}) int {
	return cache.GetInt(v)
}

// GetInt64 convert interface to int64.
func GetInt64(v interface{}) int64 {
	return cache.GetInt64(v)
}

// GetFloat64 convert interface to float64.
func GetFloat64(v interface{}) float64 {
	return cache.GetFloat64(v)
}

// GetBool convert interface to bool.
func GetBool(v interface{}) bool {
	return cache.GetBool(v)
}
