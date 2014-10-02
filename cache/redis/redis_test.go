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

package redis

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/astaxie/beego/cache"
)

func TestRedisCache(t *testing.T) {
	bm, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379"}`)
	if err != nil {
		t.Error("init err")
	}

	// clear all previous cache
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}

	if err = bm.Put("astaxie", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	if !bm.IsExist("astaxie") {
		t.Error("check err")
	}

	time.Sleep(10 * time.Second)

	if bm.IsExist("astaxie") {
		t.Error("check err")
	}

	if err = bm.Put("astaxie", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	v, err := redis.Int(bm.Get("astaxie"), nil)
	if err != nil {
		t.Error(err)
	} else if v != 1 {
		t.Error("value not as expected", v)
	}

	counter, err := bm.Incr("astaxie")
	if err != nil {
		t.Error("Incr Error", err)
	} else if counter != 2 {
		t.Error("value not as expected", counter)
	}

	counter, err = bm.Decr("astaxie")
	if err != nil {
		t.Error("Decr Error", err)
	} else if counter != 1 {
		t.Error("value not as expected", counter)
	}

	bm.Delete("astaxie")
	if bm.IsExist("astaxie") {
		t.Error("delete err")
	}
	//test string
	if err = bm.Put("astaxie", "author", 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie") {
		t.Error("check err")
	}

	if v, _ := redis.String(bm.Get("astaxie"), err); v != "author" {
		t.Error("get err")
	}
	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}

	// test HashCache interface
	hashCache := bm.(cache.HashCache)
	key := "test"
	field := "field"
	value := "value"

	// HPut
	if err = hashCache.HPut(key, field, value); err != nil {
		t.Error(err)
	}

	// HGet
	if v, _ := redis.String(hashCache.HGet(key, field), nil); v != value {
		t.Error("value get '%s' is not expected: '%s'", v, value)
	}

	// HDelete
	if err = hashCache.HDelete(key, field); err != nil {
		t.Error(err)
	}

	// HGet
	if v := hashCache.HGet(key, field); v != nil {
		t.Error("value is not deleted properly: '%s'", v)
	}

	// HGetAll
	fieldAndValues := [4]string{"field0", "value0", "field1", "value1"}
	for index := 0; index < len(fieldAndValues); index += 2 {
		if err = hashCache.HPut(key, fieldAndValues[index], fieldAndValues[index+1]); err != nil {
			t.Error(err)
		}
	}

	fieldAndValuesReturned, err := hashCache.HGetAll(key)
	if err != nil {
		t.Error(err)
	}
	if len(fieldAndValuesReturned) != len(fieldAndValues) {
		t.Error("incorrect number of fields and values returned")
	}

	for index, elem := range fieldAndValuesReturned {
		str, _ := redis.String(elem, nil)
		if str != fieldAndValues[index] {
			t.Error("unexpected returned data '%s', expecting '%s'", str, fieldAndValues[index])
		}
	}

	// HIncrBy
	initValue := uint64(100)
	delta := uint64(10)
	field = "counter"
	if err = hashCache.HPut(key, field, initValue); err != nil {
		t.Error(err)
	}

	counter, err = hashCache.HIncrBy(key, field, delta)
	if err != nil {
		t.Error(err)
	} else if counter != initValue+delta {
		t.Error("value %d is not expected: %d", counter, initValue+delta)
	}

	// HDecrBy
	if err = hashCache.HPut(key, field, initValue); err != nil {
		t.Error(err)
	}

	counter, err = hashCache.HDecrBy(key, field, delta)
	if err != nil {
		t.Error(err)
	} else if counter != initValue-delta {
		t.Error("value %d is not expected: %d", counter, initValue-delta)
	}
}
