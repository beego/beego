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

package utils

type KV interface {
	GetKey() interface{}
	GetValue() interface{}
}

// SimpleKV is common structure to store key-value pairs.
// When you need something like Pair, you can use this
type SimpleKV struct {
	Key   interface{}
	Value interface{}
}

var _ KV = new(SimpleKV)

func (s *SimpleKV) GetKey() interface{} {
	return s.Key
}

func (s *SimpleKV) GetValue() interface{} {
	return s.Value
}

// KVs interface
type KVs interface {
	GetValueOr(key interface{}, defValue interface{}) interface{}
	Contains(key interface{}) bool
	IfContains(key interface{}, action func(value interface{})) KVs
}

// SimpleKVs will store SimpleKV collection as map
type SimpleKVs struct {
	kvs map[interface{}]interface{}
}

var _ KVs = new(SimpleKVs)

// GetValueOr returns the value for a given key, if non-existent
// it returns defValue
func (kvs *SimpleKVs) GetValueOr(key interface{}, defValue interface{}) interface{} {
	v, ok := kvs.kvs[key]
	if ok {
		return v
	}
	return defValue
}

// Contains checks if a key exists
func (kvs *SimpleKVs) Contains(key interface{}) bool {
	_, ok := kvs.kvs[key]
	return ok
}

// IfContains invokes the action on a key if it exists
func (kvs *SimpleKVs) IfContains(key interface{}, action func(value interface{})) KVs {
	v, ok := kvs.kvs[key]
	if ok {
		action(v)
	}
	return kvs
}

// NewKVs creates the *KVs instance
func NewKVs(kvs ...KV) KVs {
	res := &SimpleKVs{
		kvs: make(map[interface{}]interface{}, len(kvs)),
	}
	for _, kv := range kvs {
		res.kvs[kv.GetKey()] = kv.GetValue()
	}
	return res
}
