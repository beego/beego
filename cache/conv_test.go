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
	"testing"
)

func TestGetString(t *testing.T) {
	var t1 = "test1"
	if GetString(t1) != "test1" {
		t.Error("get string from string error")
	}
	var t2 = []byte("test2")
	if GetString(t2) != "test2" {
		t.Error("get string from byte array error")
	}
	var t3 = 1
	if GetString(t3) != "1" {
		t.Error("get string from int error")
	}
	var t4 int64 = 1
	if GetString(t4) != "1" {
		t.Error("get string from int64 error")
	}
	var t5 = 1.1
	if GetString(t5) != "1.1" {
		t.Error("get string from float64 error")
	}

	if GetString(nil) != "" {
		t.Error("get string from nil error")
	}
}

func TestGetInt(t *testing.T) {
	var t1 = 1
	if GetInt(t1) != 1 {
		t.Error("get int from int error")
	}
	var t2 int32 = 32
	if GetInt(t2) != 32 {
		t.Error("get int from int32 error")
	}
	var t3 int64 = 64
	if GetInt(t3) != 64 {
		t.Error("get int from int64 error")
	}
	var t4 = "128"
	if GetInt(t4) != 128 {
		t.Error("get int from num string error")
	}
	if GetInt(nil) != 0 {
		t.Error("get int from nil error")
	}
}

func TestGetInt64(t *testing.T) {
	var i int64 = 1
	var t1 = 1
	if GetInt64(t1) != i {
		t.Error("get int64 from int error")
	}
	var t2 int32 = 1
	if GetInt64(t2) != i {
		t.Error("get int64 from int32 error")
	}
	var t3 int64 = 1
	if GetInt64(t3) != i {
		t.Error("get int64 from int64 error")
	}
	var t4 = "1"
	if GetInt64(t4) != i {
		t.Error("get int64 from num string error")
	}
	if GetInt64(nil) != 0 {
		t.Error("get int64 from nil")
	}
}

func TestGetFloat64(t *testing.T) {
	var f = 1.11
	var t1 float32 = 1.11
	if GetFloat64(t1) != f {
		t.Error("get float64 from float32 error")
	}
	var t2 = 1.11
	if GetFloat64(t2) != f {
		t.Error("get float64 from float64 error")
	}
	var t3 = "1.11"
	if GetFloat64(t3) != f {
		t.Error("get float64 from string error")
	}

	var f2 float64 = 1
	var t4 = 1
	if GetFloat64(t4) != f2 {
		t.Error("get float64 from int error")
	}

	if GetFloat64(nil) != 0 {
		t.Error("get float64 from nil error")
	}
}

func TestGetBool(t *testing.T) {
	var t1 = true
	if !GetBool(t1) {
		t.Error("get bool from bool error")
	}
	var t2 = "true"
	if !GetBool(t2) {
		t.Error("get bool from string error")
	}
	if GetBool(nil) {
		t.Error("get bool from nil error")
	}
}

func byteArrayEquals(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
