// Copyright 2020 beego-dev
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

// KV is common structure to store key-value data.
// when you need something like Pair, you can use this
type KV struct {
	Key   interface{}
	Value interface{}
}

// KVs will store KV collection as map
type KVs struct {
	kvs map[interface{}]interface{}
}

// GetValueOr check whether this contains the key,
// if the key not found, the default value will be return
func (kvs *KVs) GetValueOr(key interface{}, defValue interface{}) interface{} {
	v, ok := kvs.kvs[key]
	if ok {
		return v
	}
	return defValue
}

// Contains will check whether contains the key
func (kvs *KVs) Contains(key interface{}) bool {
	_, ok := kvs.kvs[key]
	return ok
}

// IfContains is a functional API that if the key is in KVs, the action will be invoked
func (kvs *KVs) IfContains(key interface{}, action func(value interface{})) *KVs {
	v, ok := kvs.kvs[key]
	if ok {
		action(v)
	}
	return kvs
}

// Put store the value
func (kvs *KVs) Put(key interface{}, value interface{}) *KVs {
	kvs.kvs[key] = value
	return kvs
}

// NewKVs will create the *KVs instance
func NewKVs(kvs ...KV) *KVs {
	res := &KVs{
		kvs: make(map[interface{}]interface{}, len(kvs)),
	}
	for _, kv := range kvs {
		res.kvs[kv.Key] = kv.Value
	}
	return res
}
