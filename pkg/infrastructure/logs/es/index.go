// Copyright 2020
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

package es

import (
	"fmt"

	"github.com/astaxie/beego/pkg/infrastructure/logs"
)

// IndexNaming generate the index name
type IndexNaming interface {
	IndexName(lm *logs.LogMsg) string
}

var indexNaming IndexNaming = &defaultIndexNaming{}

// SetIndexNaming will register global IndexNaming
func SetIndexNaming(i IndexNaming) {
	indexNaming = i
}

type defaultIndexNaming struct{}

func (d *defaultIndexNaming) IndexName(lm *logs.LogMsg) string {
	return fmt.Sprintf("%04d.%02d.%02d", lm.When.Year(), lm.When.Month(), lm.When.Day())
}
