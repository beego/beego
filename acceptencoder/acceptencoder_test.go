// Copyright 2015 beego Author. All Rights Reserved.
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

package acceptencoder

import "testing"

func Test_ExtractEncoding(t *testing.T) {
	if parseEncoding("gzip,deflate") != "gzip" {
		t.Fail()
	}
	if parseEncoding("deflate,gzip") != "deflate" {
		t.Fail()
	}

	if parseEncoding("gzip;q=.5,deflate") != "deflate" {
		t.Fail()
	}

	if parseEncoding("gzip;q=0,deflate") != "deflate" {
		t.Fail()
	}

	if parseEncoding("deflate;q=0.5,gzip;q=0.5,identity") != "identity" {
		t.Fail()
	}

	if parseEncoding("*") != "gzip" {
		t.Fail()
	}
}
