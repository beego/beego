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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/astaxie/beego/pkg/core/logs"
)

func TestDefaultIndexNaming_IndexName(t *testing.T) {
	tm := time.Date(2020, 9, 12, 1, 34, 45, 234, time.UTC)
	lm := &logs.LogMsg{
		When: tm,
	}

	res := (&defaultIndexNaming{}).IndexName(lm)
	assert.Equal(t, "2020.09.12", res)
}
