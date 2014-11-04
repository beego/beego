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

package beego

import "github.com/astaxie/beego/context"

// FilterFunc defines filter function type.
type FilterFunc func(*context.Context)

// FilterRouter defines filter operation before controller handler execution.
// it can match patterned url and do filter function when action arrives.
type FilterRouter struct {
	filterFunc     FilterFunc
	tree           *Tree
	pattern        string
	returnOnOutput bool
}

// ValidRouter check current request is valid for this filter.
// if matched, returns parsed params in this request by defined filter router pattern.
func (f *FilterRouter) ValidRouter(router string) (bool, map[string]string) {
	isok, params := f.tree.Match(router)
	if isok == nil {
		return false, nil
	}
	if isok, ok := isok.(bool); ok {
		return isok, params
	} else {
		return false, nil
	}
}
