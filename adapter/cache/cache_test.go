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
	"os"
	"sync"
	"testing"
	"time"
)

const (
	initError = "init err"
	setError = "set Error"
	checkError = "check err"
	getError = "get err"
	getMultiError = "GetMulti Error"
)

func TestCacheIncr(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error(initError)
	}
	// timeoutDuration := 10 * time.Second

	bm.Put("edwardhey", 0, time.Second*20)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			bm.Incr("edwardhey")
		}()
	}
	wg.Wait()
	if bm.Get("edwardhey").(int) != 10 {
		t.Error("Incr err")
	}
}

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error(initError)
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie") {
		t.Error(checkError)
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error(getError)
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("astaxie") {
		t.Error(checkError)
	}

	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error(setError, err)
	}

	if err = bm.Incr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 2 {
		t.Error(getError)
	}

	if err = bm.Decr("astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error(getError)
	}
	bm.Delete("astaxie")
	if bm.IsExist("astaxie") {
		t.Error("delete err")
	}

	// test GetMulti
	if err = bm.Put("astaxie", "author", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie") {
		t.Error(checkError)
	}
	if v := bm.Get("astaxie"); v.(string) != "author" {
		t.Error(getError)
	}

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
	if vv[0].(string) != "author" {
		t.Error(getMultiError)
	}
	if vv[1].(string) != "author1" {
		t.Error(getMultiError)
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)
	if err != nil {
		t.Error(initError)
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("astaxie", 1, timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !bm.IsExist("astaxie") {
		t.Error(checkError)
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error(getError)
	}

	if err = bm.Incr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 2 {
		t.Error(getError)
	}

	if err = bm.Decr("astaxie"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
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
	if v := bm.Get("astaxie"); v.(string) != "author" {
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
	if vv[0].(string) != "author" {
		t.Error(getMultiError)
	}
	if vv[1].(string) != "author1" {
		t.Error(getMultiError)
	}

	os.RemoveAll("cache")
}
