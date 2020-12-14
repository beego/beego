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
	"github.com/beego/beego/v2/client/httplib/testing"
)

var port = ""
var baseURL = "http://localhost:"

// TestHTTPRequest beego test request client
type TestHTTPRequest testing.TestHTTPRequest

// Get returns test client in GET method
func Get(path string) *TestHTTPRequest {
	return (*TestHTTPRequest)(testing.Get(path))
}

// Post returns test client in POST method
func Post(path string) *TestHTTPRequest {
	return (*TestHTTPRequest)(testing.Post(path))
}

// Put returns test client in PUT method
func Put(path string) *TestHTTPRequest {
	return (*TestHTTPRequest)(testing.Put(path))
}

// Delete returns test client in DELETE method
func Delete(path string) *TestHTTPRequest {
	return (*TestHTTPRequest)(testing.Delete(path))
}

// Head returns test client in HEAD method
func Head(path string) *TestHTTPRequest {
	return (*TestHTTPRequest)(testing.Head(path))
}
