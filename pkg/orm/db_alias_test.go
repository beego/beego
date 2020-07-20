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

	"github.com/astaxie/beego/pkg/common"
)

func TestRegisterDataBase(t *testing.T) {
	err := RegisterDataBase("test-params", DBARGS.Driver, DBARGS.Source, common.KV{
		Key:   MaxIdleConnsKey,
		Value: 20,
	}, common.KV{
		Key:   MaxOpenConnsKey,
		Value: 300,
	}, common.KV{
		Key:   ConnMaxLifetimeKey,
		Value: time.Minute,
	})
	assert.Nil(t, err)

	al := getDbAlias("test-params")
	assert.NotNil(t, al)
	assert.Equal(t, al.MaxIdleConns, 20)
	assert.Equal(t, al.MaxOpenConns, 300)
	assert.Equal(t, al.ConnMaxLifetime, time.Minute)
}
