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
	"github.com/beego/beego/core/utils"
)

// BeeMap is a map with lock
type BeeMap utils.BeeMap

// NewBeeMap return new safemap
func NewBeeMap() *BeeMap {
	return (*BeeMap)(utils.NewBeeMap())
}

// Get from maps return the k's value
func (m *BeeMap) Get(k interface{}) interface{} {
	return (*utils.BeeMap)(m).Get(k)
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BeeMap) Set(k interface{}, v interface{}) bool {
	return (*utils.BeeMap)(m).Set(k, v)
}

// Check Returns true if k is exist in the map.
func (m *BeeMap) Check(k interface{}) bool {
	return (*utils.BeeMap)(m).Check(k)
}

// Delete the given key and value.
func (m *BeeMap) Delete(k interface{}) {
	(*utils.BeeMap)(m).Delete(k)
}

// Items returns all items in safemap.
func (m *BeeMap) Items() map[interface{}]interface{} {
	return (*utils.BeeMap)(m).Items()
}

// Count returns the number of items within the map.
func (m *BeeMap) Count() int {
	return (*utils.BeeMap)(m).Count()
}
