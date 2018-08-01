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

package utils

import (
	"testing"
)

func TestInSlice(t *testing.T) {
	sl := []string{"A", "b"}
	if !InSlice("A", sl) {
		t.Error("should be true")
	}
	if InSlice("B", sl) {
		t.Error("should be false")
	}
}

func TestInSliceIface(t *testing.T) {
	sl := []string{"A", "b"}
	if !InSliceIface("A", sl) {
		t.Error("should be true")
	}
	if InSliceIface("B", sl) {
		t.Error("should be false")
	}

	newsl := []interface{}{"A", "b"}
	if !InSliceIface("A", newsl) {
		t.Error("should be true")
	}
	if InSliceIface("B", newsl) {
		t.Error("should be false")
	}

	if InSliceIface("B", "C") {
		t.Error("should be false")
	}

	type testStruct struct {
		name string
		age  int
	}
	xiaoming := &testStruct{name: "xiaoming", age: 10}
	lilei := &testStruct{name: "lilei", age: 11}
	hanmeimei := &testStruct{name: "hanmeimei", age: 11}

	testSlice := make([]*testStruct, 2, 2)
	testSlice[0] = hanmeimei
	testSlice[1] = lilei
	if !InSliceIface(lilei, testSlice) {
		t.Error("should be true")
	}

	if InSliceIface(xiaoming, testSlice) {
		t.Error("should be false")
	}
}

func TestSliceDiff(t *testing.T) {
	var sl1 = []interface{}{"A", "b"}
	var sl2 = []interface{}{"A", "c"}

	var sl3 = SliceDiff(sl1, sl2)
	for _, v := range sl3 {
		if v != "b" {
			t.Error("should be b")
		}
	}
}
