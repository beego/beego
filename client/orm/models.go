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
	"reflect"
	"runtime/debug"
	"strings"
	"sync"

	imodels "github.com/beego/beego/v2/client/orm/internal/models"
)

var defaultModelCache = NewModelCacheHandler()

// model info collection
type modelCache struct {
	sync.RWMutex    // only used outsite for bootStrap
	orders          []string
	cache           map[string]*imodels.ModelInfo
	cacheByFullName map[string]*imodels.ModelInfo
	done            bool
}

// NewModelCacheHandler generator of modelCache
func NewModelCacheHandler() *modelCache {
	return &modelCache{
		cache:           make(map[string]*imodels.ModelInfo),
		cacheByFullName: make(map[string]*imodels.ModelInfo),
	}
}

// get all model info
func (mc *modelCache) all() map[string]*imodels.ModelInfo {
	m := make(map[string]*imodels.ModelInfo, len(mc.cache))
	for k, v := range mc.cache {
		m[k] = v
	}
	return m
}

// get ordered model info
func (mc *modelCache) allOrdered() []*imodels.ModelInfo {
	m := make([]*imodels.ModelInfo, 0, len(mc.orders))
	for _, table := range mc.orders {
		m = append(m, mc.cache[table])
	}
	return m
}

// get model info by table name
func (mc *modelCache) get(table string) (mi *imodels.ModelInfo, ok bool) {
	mi, ok = mc.cache[table]
	return
}

// get model info by full name
func (mc *modelCache) getByFullName(name string) (mi *imodels.ModelInfo, ok bool) {
	mi, ok = mc.cacheByFullName[name]
	return
}

func (mc *modelCache) getByMd(md interface{}) (*imodels.ModelInfo, bool) {
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	name := imodels.GetFullName(typ)
	return mc.getByFullName(name)
}

// set model info to collection
func (mc *modelCache) set(table string, mi *imodels.ModelInfo) *imodels.ModelInfo {
	mii := mc.cache[table]
	mc.cache[table] = mi
	mc.cacheByFullName[mi.FullName] = mi
	if mii == nil {
		mc.orders = append(mc.orders, table)
	}
	return mii
}

// clean all model info.
func (mc *modelCache) clean() {
	mc.Lock()
	defer mc.Unlock()

	mc.orders = make([]string, 0)
	mc.cache = make(map[string]*imodels.ModelInfo)
	mc.cacheByFullName = make(map[string]*imodels.ModelInfo)
	mc.done = false
}

// bootstrap bootstrap for models
func (mc *modelCache) bootstrap() {
	mc.Lock()
	defer mc.Unlock()
	if mc.done {
		return
	}
	var (
		err    error
		models map[string]*imodels.ModelInfo
	)
	if dataBaseCache.getDefault() == nil {
		err = fmt.Errorf("must have one register DataBase alias named `default`")
		goto end
	}

	// set rel and reverse model
	// RelManyToMany set the relTable
	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.Fields.Columns {
			if fi.Rel || fi.Reverse {
				elm := fi.AddrValue.Type().Elem()
				if fi.FieldType == RelReverseMany || fi.FieldType == RelManyToMany {
					elm = elm.Elem()
				}
				// check the rel or reverse model already register
				name := imodels.GetFullName(elm)
				mii, ok := mc.getByFullName(name)
				if !ok || mii.Pkg != elm.PkgPath() {
					err = fmt.Errorf("can not find rel in field `%s`, `%s` may be miss register", fi.FullName, elm.String())
					goto end
				}
				fi.RelModelInfo = mii

				switch fi.FieldType {
				case RelManyToMany:
					if fi.RelThrough != "" {
						if i := strings.LastIndex(fi.RelThrough, "."); i != -1 && len(fi.RelThrough) > (i+1) {
							pn := fi.RelThrough[:i]
							rmi, ok := mc.getByFullName(fi.RelThrough)
							if !ok || pn != rmi.Pkg {
								err = fmt.Errorf("field `%s` wrong rel_through value `%s` cannot find table", fi.FullName, fi.RelThrough)
								goto end
							}
							fi.RelThroughModelInfo = rmi
							fi.RelTable = rmi.Table
						} else {
							err = fmt.Errorf("field `%s` wrong rel_through value `%s`", fi.FullName, fi.RelThrough)
							goto end
						}
					} else {
						i := imodels.NewM2MModelInfo(mi, mii)
						if fi.RelTable != "" {
							i.Table = fi.RelTable
						}
						if v := mc.set(i.Table, i); v != nil {
							err = fmt.Errorf("the rel table name `%s` already registered, cannot be use, please change one", fi.RelTable)
							goto end
						}
						fi.RelTable = i.Table
						fi.RelThroughModelInfo = i
					}

					fi.RelThroughModelInfo.IsThrough = true
				}
			}
		}
	}

	// check the rel filed while the relModelInfo also has filed point to current model
	// if not exist, add a new field to the relModelInfo
	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.Fields.FieldsRel {
			switch fi.FieldType {
			case RelForeignKey, RelOneToOne, RelManyToMany:
				inModel := false
				for _, ffi := range fi.RelModelInfo.Fields.FieldsReverse {
					if ffi.RelModelInfo == mi {
						inModel = true
						break
					}
				}
				if !inModel {
					rmi := fi.RelModelInfo
					ffi := new(imodels.FieldInfo)
					ffi.Name = mi.Name
					ffi.Column = ffi.Name
					ffi.FullName = rmi.FullName + "." + ffi.Name
					ffi.Reverse = true
					ffi.RelModelInfo = mi
					ffi.Mi = rmi
					if fi.FieldType == RelOneToOne {
						ffi.FieldType = RelReverseOne
					} else {
						ffi.FieldType = RelReverseMany
					}
					if !rmi.Fields.Add(ffi) {
						added := false
						for cnt := 0; cnt < 5; cnt++ {
							ffi.Name = fmt.Sprintf("%s%d", mi.Name, cnt)
							ffi.Column = ffi.Name
							ffi.FullName = rmi.FullName + "." + ffi.Name
							if added = rmi.Fields.Add(ffi); added {
								break
							}
						}
						if !added {
							panic(fmt.Errorf("cannot generate auto reverse field info `%s` to `%s`", fi.FullName, ffi.FullName))
						}
					}
				}
			}
		}
	}

	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.Fields.FieldsRel {
			switch fi.FieldType {
			case RelManyToMany:
				for _, ffi := range fi.RelThroughModelInfo.Fields.FieldsRel {
					switch ffi.FieldType {
					case RelOneToOne, RelForeignKey:
						if ffi.RelModelInfo == fi.RelModelInfo {
							fi.ReverseFieldInfoTwo = ffi
						}
						if ffi.RelModelInfo == mi {
							fi.ReverseField = ffi.Name
							fi.ReverseFieldInfo = ffi
						}
					}
				}
				if fi.ReverseFieldInfoTwo == nil {
					err = fmt.Errorf("can not find m2m field for m2m model `%s`, ensure your m2m model defined correct",
						fi.RelThroughModelInfo.FullName)
					goto end
				}
			}
		}
	}

	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.Fields.FieldsReverse {
			switch fi.FieldType {
			case RelReverseOne:
				found := false
			mForA:
				for _, ffi := range fi.RelModelInfo.Fields.FieldsByType[RelOneToOne] {
					if ffi.RelModelInfo == mi {
						found = true
						fi.ReverseField = ffi.Name
						fi.ReverseFieldInfo = ffi

						ffi.ReverseField = fi.Name
						ffi.ReverseFieldInfo = fi
						break mForA
					}
				}
				if !found {
					err = fmt.Errorf("reverse field `%s` not found in model `%s`", fi.FullName, fi.RelModelInfo.FullName)
					goto end
				}
			case RelReverseMany:
				found := false
			mForB:
				for _, ffi := range fi.RelModelInfo.Fields.FieldsByType[RelForeignKey] {
					if ffi.RelModelInfo == mi {
						found = true
						fi.ReverseField = ffi.Name
						fi.ReverseFieldInfo = ffi

						ffi.ReverseField = fi.Name
						ffi.ReverseFieldInfo = fi

						break mForB
					}
				}
				if !found {
				mForC:
					for _, ffi := range fi.RelModelInfo.Fields.FieldsByType[RelManyToMany] {
						conditions := fi.RelThrough != "" && fi.RelThrough == ffi.RelThrough ||
							fi.RelTable != "" && fi.RelTable == ffi.RelTable ||
							fi.RelThrough == "" && fi.RelTable == ""
						if ffi.RelModelInfo == mi && conditions {
							found = true

							fi.ReverseField = ffi.ReverseFieldInfoTwo.Name
							fi.ReverseFieldInfo = ffi.ReverseFieldInfoTwo
							fi.RelThroughModelInfo = ffi.RelThroughModelInfo
							fi.ReverseFieldInfoTwo = ffi.ReverseFieldInfo
							fi.ReverseFieldInfoM2M = ffi
							ffi.ReverseFieldInfoM2M = fi

							break mForC
						}
					}
				}
				if !found {
					err = fmt.Errorf("reverse field for `%s` not found in model `%s`", fi.FullName, fi.RelModelInfo.FullName)
					goto end
				}
			}
		}
	}

end:
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
	}
	mc.done = true
}

// register register models to model cache
func (mc *modelCache) register(prefixOrSuffixStr string, prefixOrSuffix bool, models ...interface{}) (err error) {
	for _, model := range models {
		val := reflect.ValueOf(model)
		typ := reflect.Indirect(val).Type()

		if val.Kind() != reflect.Ptr {
			err = fmt.Errorf("<orm.RegisterModel> cannot use non-ptr model struct `%s`", imodels.GetFullName(typ))
			return
		}
		// For this case:
		// u := &User{}
		// registerModel(&u)
		if typ.Kind() == reflect.Ptr {
			err = fmt.Errorf("<orm.RegisterModel> only allow ptr model struct, it looks you use two reference to the struct `%s`", typ)
			return
		}
		if val.Elem().Kind() == reflect.Slice {
			val = reflect.New(val.Elem().Type().Elem())
		}
		table := imodels.GetTableName(val)

		if prefixOrSuffixStr != "" {
			if prefixOrSuffix {
				table = prefixOrSuffixStr + table
			} else {
				table = table + prefixOrSuffixStr
			}
		}

		// models's fullname is pkgpath + struct name
		name := imodels.GetFullName(typ)
		if _, ok := mc.getByFullName(name); ok {
			err = fmt.Errorf("<orm.RegisterModel> model `%s` repeat register, must be unique\n", name)
			return
		}

		if _, ok := mc.get(table); ok {
			return nil
		}

		mi := imodels.NewModelInfo(val)
		if mi.Fields.Pk == nil {
		outFor:
			for _, fi := range mi.Fields.FieldsDB {
				if strings.ToLower(fi.Name) == "id" {
					switch fi.AddrValue.Elem().Kind() {
					case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
						fi.Auto = true
						fi.Pk = true
						mi.Fields.Pk = fi
						break outFor
					}
				}
			}
		}

		mi.Table = table
		mi.Pkg = typ.PkgPath()
		mi.Model = model
		mi.Manual = true

		mc.set(table, mi)
	}
	return
}

// getDbDropSQL get database scheme drop sql queries
func (mc *modelCache) getDbDropSQL(al *alias) (queries []string, err error) {
	if len(mc.cache) == 0 {
		err = errors.New("no Model found, need register your model")
		return
	}

	Q := al.DbBaser.TableQuote()

	for _, mi := range mc.allOrdered() {
		queries = append(queries, fmt.Sprintf(`DROP TABLE IF EXISTS %s%s%s`, Q, mi.Table, Q))
	}
	return queries, nil
}

// getDbCreateSQL get database scheme creation sql queries
func (mc *modelCache) getDbCreateSQL(al *alias) (queries []string, tableIndexes map[string][]dbIndex, err error) {
	if len(mc.cache) == 0 {
		err = errors.New("no Model found, need register your model")
		return
	}

	Q := al.DbBaser.TableQuote()
	T := al.DbBaser.DbTypes()
	sep := fmt.Sprintf("%s, %s", Q, Q)

	tableIndexes = make(map[string][]dbIndex)

	for _, mi := range mc.allOrdered() {
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

			if fi.Auto {
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

// ResetModelCache Clean model cache. Then you can re-RegisterModel.
// Common use this api for test case.
func ResetModelCache() {
	defaultModelCache.clean()
}
