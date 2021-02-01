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

package memcache

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/beego/beego/v2/adapter/cache"
)

const (
	initError = "init err"
	setError = "set Error"
	checkError = "check err"
	getError = "get err"
	getMultiError = "GetMulti Error"
)

func TestMemcacheCache(t *testing.T) {

	addr := os.Getenv("MEMCACHE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11211"
	}

	bm, err := cache.NewCache("memcache", fmt.Sprintf(`{"conn": "%s"}`, addr))
	if err != nil {
		t.Error(initError)
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("astaxie", "1", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie") {
		t.Error(checkError)
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist("astaxie") {
		t.Error(checkError)
	}
	if err = bm.Put("astaxie", "1", timeoutDuration); err != nil {
		t.Error(setError, err)
	}

	if v, err := strconv.Atoi(string(bm.Get("astaxie").([]byte))); err != nil || v != 1 {
		t.Error(getError)
	}

	if err = bm.Incr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get("astaxie").([]byte))); err != nil || v != 2 {
		t.Error(getError)
	}

	if err = bm.Decr("astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get("astaxie").([]byte))); err != nil || v != 1 {
		t.Error(getError)
	}
	bm.Delete("astaxie")
	if bm.IsExist("astaxie") {
		t.Error("delete err")
	}

	// test string
	if err = bm.Put("astaxie", "author", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie") {
		t.Error(checkError)
	}

	if v := bm.Get("astaxie").([]byte); string(v) != "author" {
		t.Error(getError)
	}

	// test GetMulti
	if err = bm.Put("astaxie1", "author1", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie1") {
		t.Error(checkError)
	}

	vv := bm.GetMulti([]string{"astaxie", "astaxie1"})
	if len(vv) != 2 {
		t.Error(getMultiError)
	}
	if string(vv[0].([]byte)) != "author" && string(vv[0].([]byte)) != "author1" {
		t.Error(getMultiError)
	}
	if string(vv[1].([]byte)) != "author1" && string(vv[1].([]byte)) != "author" {
		t.Error(getMultiError)
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}
