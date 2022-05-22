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
)

const (
	odCascade             = "cascade"
	odSetNULL             = "set_null"
	odSetDefault          = "set_default"
	odDoNothing           = "do_nothing"
	defaultStructTagName  = "orm"
	defaultStructTagDelim = ";"
)

var defaultModelCache = NewModelCacheHandler()

// model info collection
type modelCache struct {
	sync.RWMutex    // only used outsite for bootStrap
	orders          []string
	cache           map[string]*modelInfo
	cacheByFullName map[string]*modelInfo
	done            bool
}

// NewModelCacheHandler generator of modelCache
func NewModelCacheHandler() *modelCache {
	return &modelCache{
		cache:           make(map[string]*modelInfo),
		cacheByFullName: make(map[string]*modelInfo),
	}
}

// get all model info
func (mc *modelCache) all() map[string]*modelInfo {
	m := make(map[string]*modelInfo, len(mc.cache))
	for k, v := range mc.cache {
		m[k] = v
	}
	return m
}

// get ordered model info
func (mc *modelCache) allOrdered() []*modelInfo {
	m := make([]*modelInfo, 0, len(mc.orders))
	for _, table := range mc.orders {
		m = append(m, mc.cache[table])
	}
	return m
}

// get model info by table name
func (mc *modelCache) get(table string) (mi *modelInfo, ok bool) {
	mi, ok = mc.cache[table]
	return
}

// get model info by full name
func (mc *modelCache) getByFullName(name string) (mi *modelInfo, ok bool) {
	mi, ok = mc.cacheByFullName[name]
	return
}

func (mc *modelCache) getByMd(md interface{}) (*modelInfo, bool) {
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	name := getFullName(typ)
	return mc.getByFullName(name)
}

// set model info to collection
func (mc *modelCache) set(table string, mi *modelInfo) *modelInfo {
	mii := mc.cache[table]
	mc.cache[table] = mi
	mc.cacheByFullName[mi.fullName] = mi
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
	mc.cache = make(map[string]*modelInfo)
	mc.cacheByFullName = make(map[string]*modelInfo)
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
		models map[string]*modelInfo
	)
	if dataBaseCache.getDefault() == nil {
		err = fmt.Errorf("must have one register DataBase alias named `default`")
		goto end
	}

	// set rel and reverse model
	// RelManyToMany set the relTable
	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.fields.columns {
			if fi.rel || fi.reverse {
				elm := fi.addrValue.Type().Elem()
				if fi.fieldType == RelReverseMany || fi.fieldType == RelManyToMany {
					elm = elm.Elem()
				}
				// check the rel or reverse model already register
				name := getFullName(elm)
				mii, ok := mc.getByFullName(name)
				if !ok || mii.pkg != elm.PkgPath() {
					err = fmt.Errorf("can not find rel in field `%s`, `%s` may be miss register", fi.fullName, elm.String())
					goto end
				}
				fi.relModelInfo = mii

				switch fi.fieldType {
				case RelManyToMany:
					if fi.relThrough != "" {
						if i := strings.LastIndex(fi.relThrough, "."); i != -1 && len(fi.relThrough) > (i+1) {
							pn := fi.relThrough[:i]
							rmi, ok := mc.getByFullName(fi.relThrough)
							if !ok || pn != rmi.pkg {
								err = fmt.Errorf("field `%s` wrong rel_through value `%s` cannot find table", fi.fullName, fi.relThrough)
								goto end
							}
							fi.relThroughModelInfo = rmi
							fi.relTable = rmi.table
						} else {
							err = fmt.Errorf("field `%s` wrong rel_through value `%s`", fi.fullName, fi.relThrough)
							goto end
						}
					} else {
						i := newM2MModelInfo(mi, mii)
						if fi.relTable != "" {
							i.table = fi.relTable
						}
						if v := mc.set(i.table, i); v != nil {
							err = fmt.Errorf("the rel table name `%s` already registered, cannot be use, please change one", fi.relTable)
							goto end
						}
						fi.relTable = i.table
						fi.relThroughModelInfo = i
					}

					fi.relThroughModelInfo.isThrough = true
				}
			}
		}
	}

	// check the rel filed while the relModelInfo also has filed point to current model
	// if not exist, add a new field to the relModelInfo
	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.fields.fieldsRel {
			switch fi.fieldType {
			case RelForeignKey, RelOneToOne, RelManyToMany:
				inModel := false
				for _, ffi := range fi.relModelInfo.fields.fieldsReverse {
					if ffi.relModelInfo == mi {
						inModel = true
						break
					}
				}
				if !inModel {
					rmi := fi.relModelInfo
					ffi := new(fieldInfo)
					ffi.name = mi.name
					ffi.column = ffi.name
					ffi.fullName = rmi.fullName + "." + ffi.name
					ffi.reverse = true
					ffi.relModelInfo = mi
					ffi.mi = rmi
					if fi.fieldType == RelOneToOne {
						ffi.fieldType = RelReverseOne
					} else {
						ffi.fieldType = RelReverseMany
					}
					if !rmi.fields.Add(ffi) {
						added := false
						for cnt := 0; cnt < 5; cnt++ {
							ffi.name = fmt.Sprintf("%s%d", mi.name, cnt)
							ffi.column = ffi.name
							ffi.fullName = rmi.fullName + "." + ffi.name
							if added = rmi.fields.Add(ffi); added {
								break
							}
						}
						if !added {
							panic(fmt.Errorf("cannot generate auto reverse field info `%s` to `%s`", fi.fullName, ffi.fullName))
						}
					}
				}
			}
		}
	}

	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.fields.fieldsRel {
			switch fi.fieldType {
			case RelManyToMany:
				for _, ffi := range fi.relThroughModelInfo.fields.fieldsRel {
					switch ffi.fieldType {
					case RelOneToOne, RelForeignKey:
						if ffi.relModelInfo == fi.relModelInfo {
							fi.reverseFieldInfoTwo = ffi
						}
						if ffi.relModelInfo == mi {
							fi.reverseField = ffi.name
							fi.reverseFieldInfo = ffi
						}
					}
				}
				if fi.reverseFieldInfoTwo == nil {
					err = fmt.Errorf("can not find m2m field for m2m model `%s`, ensure your m2m model defined correct",
						fi.relThroughModelInfo.fullName)
					goto end
				}
			}
		}
	}

	models = mc.all()
	for _, mi := range models {
		for _, fi := range mi.fields.fieldsReverse {
			switch fi.fieldType {
			case RelReverseOne:
				found := false
			mForA:
				for _, ffi := range fi.relModelInfo.fields.fieldsByType[RelOneToOne] {
					if ffi.relModelInfo == mi {
						found = true
						fi.reverseField = ffi.name
						fi.reverseFieldInfo = ffi

						ffi.reverseField = fi.name
						ffi.reverseFieldInfo = fi
						break mForA
					}
				}
				if !found {
					err = fmt.Errorf("reverse field `%s` not found in model `%s`", fi.fullName, fi.relModelInfo.fullName)
					goto end
				}
			case RelReverseMany:
				found := false
			mForB:
				for _, ffi := range fi.relModelInfo.fields.fieldsByType[RelForeignKey] {
					if ffi.relModelInfo == mi {
						found = true
						fi.reverseField = ffi.name
						fi.reverseFieldInfo = ffi

						ffi.reverseField = fi.name
						ffi.reverseFieldInfo = fi

						break mForB
					}
				}
				if !found {
				mForC:
					for _, ffi := range fi.relModelInfo.fields.fieldsByType[RelManyToMany] {
						conditions := fi.relThrough != "" && fi.relThrough == ffi.relThrough ||
							fi.relTable != "" && fi.relTable == ffi.relTable ||
							fi.relThrough == "" && fi.relTable == ""
						if ffi.relModelInfo == mi && conditions {
							found = true

							fi.reverseField = ffi.reverseFieldInfoTwo.name
							fi.reverseFieldInfo = ffi.reverseFieldInfoTwo
							fi.relThroughModelInfo = ffi.relThroughModelInfo
							fi.reverseFieldInfoTwo = ffi.reverseFieldInfo
							fi.reverseFieldInfoM2M = ffi
							ffi.reverseFieldInfoM2M = fi

							break mForC
						}
					}
				}
				if !found {
					err = fmt.Errorf("reverse field for `%s` not found in model `%s`", fi.fullName, fi.relModelInfo.fullName)
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
			err = fmt.Errorf("<orm.RegisterModel> cannot use non-ptr model struct `%s`", getFullName(typ))
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
		table := getTableName(val)

		if prefixOrSuffixStr != "" {
			if prefixOrSuffix {
				table = prefixOrSuffixStr + table
			} else {
				table = table + prefixOrSuffixStr
			}
		}

		// models's fullname is pkgpath + struct name
		name := getFullName(typ)
		if _, ok := mc.getByFullName(name); ok {
			err = fmt.Errorf("<orm.RegisterModel> model `%s` repeat register, must be unique\n", name)
			return
		}

		if _, ok := mc.get(table); ok {
			return nil
		}

		mi := newModelInfo(val)
		if mi.fields.pk == nil {
		outFor:
			for _, fi := range mi.fields.fieldsDB {
				if strings.ToLower(fi.name) == "id" {
					switch fi.addrValue.Elem().Kind() {
					case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
						fi.auto = true
						fi.pk = true
						mi.fields.pk = fi
						break outFor
					}
				}
			}
		}

		mi.table = table
		mi.pkg = typ.PkgPath()
		mi.model = model
		mi.manual = true

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
		queries = append(queries, fmt.Sprintf(`DROP TABLE IF EXISTS %s%s%s`, Q, mi.table, Q))
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
		sql += fmt.Sprintf("--  Table Structure for `%s`\n", mi.fullName)
		sql += fmt.Sprintf("-- %s\n", strings.Repeat("-", 50))

		sql += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s%s%s (\n", Q, mi.table, Q)

		columns := make([]string, 0, len(mi.fields.fieldsDB))

		sqlIndexes := [][]string{}
		var commentIndexes []int // store comment indexes for postgres

		for i, fi := range mi.fields.fieldsDB {
			column := fmt.Sprintf("    %s%s%s ", Q, fi.column, Q)
			col := getColumnTyp(al, fi)

			if fi.auto {
				switch al.Driver {
				case DRSqlite, DRPostgres:
					column += T["auto"]
				default:
					column += col + " " + T["auto"]
				}
			} else if fi.pk {
				column += col + " " + T["pk"]
			} else {
				column += col

				if !fi.null {
					column += " " + "NOT NULL"
				}

				// if fi.initial.String() != "" {
				//	column += " DEFAULT " + fi.initial.String()
				// }

				// Append attribute DEFAULT
				column += getColumnDefault(fi)

				if fi.unique {
					column += " " + "UNIQUE"
				}

				if fi.index {
					sqlIndexes = append(sqlIndexes, []string{fi.column})
				}
			}

			if strings.Contains(column, "%COL%") {
				column = strings.Replace(column, "%COL%", fi.column, -1)
			}

			if fi.description != "" && al.Driver != DRSqlite {
				if al.Driver == DRPostgres {
					commentIndexes = append(commentIndexes, i)
				} else {
					column += " " + fmt.Sprintf("COMMENT '%s'", fi.description)
				}
			}

			columns = append(columns, column)
		}

		if mi.model != nil {
			allnames := getTableUnique(mi.addrField)
			if !mi.manual && len(mi.uniques) > 0 {
				allnames = append(allnames, mi.uniques)
			}
			for _, names := range allnames {
				cols := make([]string, 0, len(names))
				for _, name := range names {
					if fi, ok := mi.fields.GetByAny(name); ok && fi.dbcol {
						cols = append(cols, fi.column)
					} else {
						panic(fmt.Errorf("cannot found column `%s` when parse UNIQUE in `%s.TableUnique`", name, mi.fullName))
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
			if mi.model != nil {
				engine = getTableEngine(mi.addrField)
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
					mi.table,
					Q,
					Q,
					mi.fields.fieldsDB[index].column,
					Q,
					mi.fields.fieldsDB[index].description)
			}
		}
		queries = append(queries, sql)

		if mi.model != nil {
			for _, names := range getTableIndex(mi.addrField) {
				cols := make([]string, 0, len(names))
				for _, name := range names {
					if fi, ok := mi.fields.GetByAny(name); ok && fi.dbcol {
						cols = append(cols, fi.column)
					} else {
						panic(fmt.Errorf("cannot found column `%s` when parse INDEX in `%s.TableIndex`", name, mi.fullName))
					}
				}
				sqlIndexes = append(sqlIndexes, cols)
			}
		}

		for _, names := range sqlIndexes {
			name := mi.table + "_" + strings.Join(names, "_")
			cols := strings.Join(names, sep)
			sql := fmt.Sprintf("CREATE INDEX %s%s%s ON %s%s%s (%s%s%s);", Q, name, Q, Q, mi.table, Q, Q, cols, Q)

			index := dbIndex{}
			index.Table = mi.table
			index.Name = name
			index.SQL = sql

			tableIndexes[mi.table] = append(tableIndexes[mi.table], index)
		}

	}

	return
}

// ResetModelCache Clean model cache. Then you can re-RegisterModel.
// Common use this api for test case.
func ResetModelCache() {
	defaultModelCache.clean()
}
