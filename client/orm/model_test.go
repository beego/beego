// Copyright 2022
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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ModelWithComments struct {
	Id       int    `orm:"description(user id)"`
	UserName string `orm:"size(30);unique;description(user name)"`
	Email    string `orm:"size(100)"`
	Password string `orm:"size(100)"`
}

func TestGetDbCreateSQLWithComment(t *testing.T) {
	al := getDbAlias("default")
	modelCache.clean()
	RegisterModel(&ModelWithComments{})
	queries, _, _ := modelCache.getDbCreateSQL(al)
	header := `-- --------------------------------------------------
--  Table Structure for ` + fmt.Sprintf("`%s`", modelCache.allOrdered()[0].fullName) + `
-- --------------------------------------------------`
	if al.Driver == DRPostgres {
		assert.Equal(t, queries, []string{header + `
CREATE TABLE IF NOT EXISTS "model_with_comments" (
    "id" serial NOT NULL PRIMARY KEY,
    "user_name" varchar(30) NOT NULL DEFAULT ''  UNIQUE,
    "email" varchar(100) NOT NULL DEFAULT '' ,
    "password" varchar(100) NOT NULL DEFAULT '' 
);
COMMENT ON COLUMN "model_with_comments"."id" is 'user id';
COMMENT ON COLUMN "model_with_comments"."user_name" is 'user name';`,
		})
	} else if al.Driver == DRMySQL {
		assert.Equal(t, queries, []string{header + `
CREATE TABLE IF NOT EXISTS ` + "`model_with_comments`" + ` (
    ` + "`id`" + ` integer AUTO_INCREMENT NOT NULL PRIMARY KEY COMMENT 'user id',
    ` + "`user_name`" + ` varchar(30) NOT NULL DEFAULT ''  UNIQUE COMMENT 'user name',
    ` + "`email`" + ` varchar(100) NOT NULL DEFAULT '' ,
    ` + "`password`" + ` varchar(100) NOT NULL DEFAULT '' 
) ENGINE=INNODB;`,
		})
	} else if al.Driver == DRSqlite {
		assert.Equal(t, queries, []string{header + `
CREATE TABLE IF NOT EXISTS ` + "`model_with_comments`" + ` (
    ` + "`id`" + ` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    ` + "`user_name`" + ` varchar(30) NOT NULL DEFAULT ''  UNIQUE,
    ` + "`email`" + ` varchar(100) NOT NULL DEFAULT '' ,
    ` + "`password`" + ` varchar(100) NOT NULL DEFAULT '' 
);`,
		})
	}
	modelCache.clean()
}
