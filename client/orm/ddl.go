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

package orm

import (
	"errors"
	"fmt"
	"strings"

	imodels "github.com/beego/beego/v2/client/orm/internal/models"
)

// getDbDropSQL Get database scheme drop sql queries
func getDbDropSQL(mc *imodels.ModelCache, al *alias) (queries []string, err error) {
	if mc.Empty() {
		err = errors.New("no Model found, need Register your model")
		return
	}

	Q := al.DbBaser.TableQuote()

	for _, mi := range mc.AllOrdered() {
		queries = append(queries, fmt.Sprintf(`DROP TABLE IF EXISTS %s%s%s`, Q, mi.Table, Q))
	}
	return queries, nil
}

// getDbCreateSQL Get database scheme creation sql queries
func getDbCreateSQL(mc *imodels.ModelCache, al *alias) (queries []string, tableIndexes map[string][]dbIndex, err error) {
	if mc.Empty() {
		err = errors.New("no Model found, need Register your model")
		return
	}

	Q := al.DbBaser.TableQuote()
	T := al.DbBaser.DbTypes()
	sep := fmt.Sprintf("%s, %s", Q, Q)

	tableIndexes = make(map[string][]dbIndex)

	for _, mi := range mc.AllOrdered() {
		sql := fmt.Sprintf("-- %s\n", strings.Repeat("-", 50))
		sql += fmt.Sprintf("--  Table Structure for `%s`\n", mi.FullName)
		sql += fmt.Sprintf("-- %s\n", strings.Repeat("-", 50))

		sql += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s%s%s (\n", Q, mi.Table, Q)

		columns := make([]string, 0, len(mi.Fields.FieldsDB))

		sqlIndexes := [][]string{}
		var commentIndexes []int // store comment indexes for postgres

		for i, fi := range mi.Fields.FieldsDB {
			column := fmt.Sprintf("    %s%s%s ", Q, fi.Column, Q)
			col := getColumnTyp(al, fi)
			if fi.DBType != "" {
				column += fi.DBType
			} else if fi.Auto {
				switch al.Driver {
				case DRSqlite, DRPostgres:
					column += T["auto"]
				default:
					column += col + " " + T["auto"]
				}
			} else if fi.Pk {
				column += col + " " + T["pk"]
			} else {
				column += col

				if !fi.Null {
					column += " " + "NOT NULL"
				}

				// if fi.initial.String() != "" {
				//	column += " DEFAULT " + fi.initial.String()
				// }

				// Append attribute DEFAULT
				column += getColumnDefault(fi)

				if fi.Unique {
					column += " " + "UNIQUE"
				}

				if fi.Index {
					sqlIndexes = append(sqlIndexes, []string{fi.Column})
				}
			}

			if strings.Contains(column, "%COL%") {
				column = strings.Replace(column, "%COL%", fi.Column, -1)
			}

			if fi.Description != "" && al.Driver != DRSqlite {
				if al.Driver == DRPostgres {
					commentIndexes = append(commentIndexes, i)
				} else {
					column += " " + fmt.Sprintf("COMMENT '%s'", fi.Description)
				}
			}

			columns = append(columns, column)
		}

		if mi.Model != nil {
			allnames := imodels.GetTableUnique(mi.AddrField)
			if !mi.Manual && len(mi.Uniques) > 0 {
				allnames = append(allnames, mi.Uniques)
			}
			for _, names := range allnames {
				cols := make([]string, 0, len(names))
				for _, name := range names {
					if fi, ok := mi.Fields.GetByAny(name); ok && fi.DBcol {
						cols = append(cols, fi.Column)
					} else {
						panic(fmt.Errorf("cannot found column `%s` when parse UNIQUE in `%s.TableUnique`", name, mi.FullName))
					}
				}
				column := fmt.Sprintf("    UNIQUE (%s%s%s)", Q, strings.Join(cols, sep), Q)
				columns = append(columns, column)
			}
		}

		sql += strings.Join(columns, ",\n")
		sql += "\n)"

		if al.Driver == DRMySQL {
			var engine string
			if mi.Model != nil {
				engine = imodels.GetTableEngine(mi.AddrField)
			}
			if engine == "" {
				engine = al.Engine
			}
			sql += " ENGINE=" + engine
		}

		sql += ";"
		if al.Driver == DRPostgres && len(commentIndexes) > 0 {
			// append comments for postgres only
			for _, index := range commentIndexes {
				sql += fmt.Sprintf("\nCOMMENT ON COLUMN %s%s%s.%s%s%s is '%s';",
					Q,
					mi.Table,
					Q,
					Q,
					mi.Fields.FieldsDB[index].Column,
					Q,
					mi.Fields.FieldsDB[index].Description)
			}
		}
		queries = append(queries, sql)

		if mi.Model != nil {
			for _, names := range imodels.GetTableIndex(mi.AddrField) {
				cols := make([]string, 0, len(names))
				for _, name := range names {
					if fi, ok := mi.Fields.GetByAny(name); ok && fi.DBcol {
						cols = append(cols, fi.Column)
					} else {
						panic(fmt.Errorf("cannot found column `%s` when parse INDEX in `%s.TableIndex`", name, mi.FullName))
					}
				}
				sqlIndexes = append(sqlIndexes, cols)
			}
		}

		for _, names := range sqlIndexes {
			name := mi.Table + "_" + strings.Join(names, "_")
			cols := strings.Join(names, sep)
			sql := fmt.Sprintf("CREATE INDEX %s%s%s ON %s%s%s (%s%s%s);", Q, name, Q, Q, mi.Table, Q, Q, cols, Q)

			index := dbIndex{}
			index.Table = mi.Table
			index.Name = name
			index.SQL = sql

			tableIndexes[mi.Table] = append(tableIndexes[mi.Table], index)
		}

	}

	return
}
