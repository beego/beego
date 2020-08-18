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

package orm

import (
	"testing"
	"time"

	"github.com/astaxie/beego/pkg/orm/hints"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDataBase(t *testing.T) {
	err := RegisterDataBase("test-params", DBARGS.Driver, DBARGS.Source,
		hints.MaxIdleConnections(20),
		hints.MaxOpenConnections(300),
		hints.ConnMaxLifetime(time.Minute))
	assert.Nil(t, err)

	al := getDbAlias("test-params")
	assert.NotNil(t, al)
	assert.Equal(t, al.MaxIdleConns, 20)
	assert.Equal(t, al.MaxOpenConns, 300)
	assert.Equal(t, al.ConnMaxLifetime, time.Minute)
}

func TestRegisterDataBase_MaxStmtCacheSizeNegative1(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSizeNegative1"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, hints.MaxStmtCacheSize(-1))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 0)
}

func TestRegisterDataBase_MaxStmtCacheSize0(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSize0"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, hints.MaxStmtCacheSize(0))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 0)
}

func TestRegisterDataBase_MaxStmtCacheSize1(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSize1"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, hints.MaxStmtCacheSize(1))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 1)
}

func TestRegisterDataBase_MaxStmtCacheSize841(t *testing.T) {
	aliasName := "TestRegisterDataBase_MaxStmtCacheSize841"
	err := RegisterDataBase(aliasName, DBARGS.Driver, DBARGS.Source, hints.MaxStmtCacheSize(841))
	assert.Nil(t, err)

	al := getDbAlias(aliasName)
	assert.NotNil(t, al)
	assert.Equal(t, al.DB.stmtDecoratorsLimit, 841)
}

func TestDBCache(t *testing.T) {
	dataBaseCache.add("test1", &alias{})
	dataBaseCache.add("default", &alias{})
	al := dataBaseCache.getDefault()
	assert.NotNil(t, al)
	al, ok := dataBaseCache.get("test1")
	assert.NotNil(t, al)
	assert.True(t, ok)
}
