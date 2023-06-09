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

package orm

import (
	"testing"

	"github.com/beego/beego/v2/client/orm/internal/models"

	"github.com/stretchr/testify/assert"
)

type Interface struct {
	Id   int
	Name string

	Index1 string
	Index2 string

	Unique1 string
	Unique2 string
}

func (i *Interface) TableIndex() [][]string {
	return [][]string{{"index1"}, {"index2"}}
}

func (i *Interface) TableUnique() [][]string {
	return [][]string{{"unique1"}, {"unique2"}}
}

func (i *Interface) TableName() string {
	return "INTERFACE_"
}

func (i *Interface) TableEngine() string {
	return "innodb"
}

func TestDbBase_GetTables(t *testing.T) {
	RegisterModel(&Interface{})
	mi, ok := defaultModelCache.get("INTERFACE_")
	assert.True(t, ok)
	assert.NotNil(t, mi)

	engine := models.GetTableEngine(mi.AddrField)
	assert.Equal(t, "innodb", engine)
	uniques := models.GetTableUnique(mi.AddrField)
	assert.Equal(t, [][]string{{"unique1"}, {"unique2"}}, uniques)
	indexes := models.GetTableIndex(mi.AddrField)
	assert.Equal(t, [][]string{{"index1"}, {"index2"}}, indexes)
}
