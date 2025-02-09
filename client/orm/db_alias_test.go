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

	"github.com/stretchr/testify/assert"
)

func TestRegisterDataBase(t *testing.T) {
	err := RegisterDataBase("test-params", DBARGS.Driver, DBARGS.Source,
		MaxIdleConnections(20),
		MaxOpenConnections(300),
		ConnMaxLifetime(time.Minute),
		ConnMaxIdletime(time.Minute))
	assert.Nil(t, err)

	al := getDB("test-params")
	assert.NotNil(t, al)
	assert.Equal(t, al.MaxIdleConns, 20)
	assert.Equal(t, al.MaxOpenConns, 300)
	assert.Equal(t, al.ConnMaxLifetime, time.Minute)
	assert.Equal(t, al.ConnMaxIdletime, time.Minute)
}

func TestRegisterDataBaseMaxStmtCacheSizeNegative1(t *testing.T) {
	dbName := "TestRegisterDataBase_MaxStmtCacheSizeNegative1"
	err := RegisterDataBase(dbName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(-1))
	assert.Nil(t, err)

	al := getDB(dbName)
	assert.NotNil(t, al)
	assert.Equal(t, al.stmtDecoratorsLimit, 0)
}

func TestRegisterDataBaseMaxStmtCacheSize0(t *testing.T) {
	dbName := "TestRegisterDataBase_MaxStmtCacheSize0"
	err := RegisterDataBase(dbName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(0))
	assert.Nil(t, err)

	al := getDB(dbName)
	assert.NotNil(t, al)
	assert.Equal(t, al.stmtDecoratorsLimit, 0)
}

func TestRegisterDataBaseMaxStmtCacheSize1(t *testing.T) {
	dbName := "TestRegisterDataBase_MaxStmtCacheSize1"
	err := RegisterDataBase(dbName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(1))
	assert.Nil(t, err)

	al := getDB(dbName)
	assert.NotNil(t, al)
	assert.Equal(t, al.stmtDecoratorsLimit, 1)
}

func TestRegisterDataBaseMaxStmtCacheSize841(t *testing.T) {
	dbName := "TestRegisterDataBase_MaxStmtCacheSize841"
	err := RegisterDataBase(dbName, DBARGS.Driver, DBARGS.Source, MaxStmtCacheSize(841))
	assert.Nil(t, err)

	al := getDB(dbName)
	assert.NotNil(t, al)
	assert.Equal(t, al.stmtDecoratorsLimit, 841)
}

func TestDBCache(t *testing.T) {
	dataBaseCache.add("test1", &DB{})
	dataBaseCache.add("default", &DB{})
	al := dataBaseCache.getDefault()
	assert.NotNil(t, al)
	al, ok := dataBaseCache.get("test1")
	assert.NotNil(t, al)
	assert.True(t, ok)
}
