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
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var DBARGS = struct {
	Driver string
	Source string
	Debug  string
}{
	os.Getenv("ORM_DRIVER"),
	os.Getenv("ORM_SOURCE"),
	os.Getenv("ORM_DEBUG"),
}

func TestRegisterDataBase(t *testing.T) {
	err := RegisterDataBase("test-adapt1", DBARGS.Driver, DBARGS.Source)
	assert.Nil(t, err)
	err = RegisterDataBase("test-adapt2", DBARGS.Driver, DBARGS.Source, 20)
	assert.Nil(t, err)
	err = RegisterDataBase("test-adapt3", DBARGS.Driver, DBARGS.Source, 20, 300)
	assert.Nil(t, err)
	err = RegisterDataBase("test-adapt4", DBARGS.Driver, DBARGS.Source, 20, 300, 60*1000)
	assert.Nil(t, err)
}
