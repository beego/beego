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

package session

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/orm/internal/models"
)

func Test_getColumnTyp(t *testing.T) {
	testCases := []struct {
		name string
		fi   *models.FieldInfo
		al   *DB

		wantCol string
	}{
		{
			// https://github.com/beego/beego/issues/5254
			name: "issue 5254",
			fi: &models.FieldInfo{
				FieldType: TypePositiveIntegerField,
				Column:    "my_col",
			},
			al: &DB{
				dbBaser: newdbBasePostgres(),
			},
			wantCol: `bigint CHECK("my_col" >= 0)`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			col := getColumnTyp(tc.al, tc.fi)
			assert.Equal(t, tc.wantCol, col)
		})
	}
}
