// Copyright 2023 beego. All Rights Reserved.
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

package models

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
)

// ModelCache info collection
type ModelCache struct {
	sync.RWMutex    // only used outsite for bootStrap
	orders          []string
	cache           map[string]*ModelInfo
	cacheByFullName map[string]*ModelInfo
	done            bool
}

// NewModelCacheHandler generator of ModelCache
func NewModelCacheHandler() *ModelCache {
	return &ModelCache{
		cache:           make(map[string]*ModelInfo),
		cacheByFullName: make(map[string]*ModelInfo),
	}
}

// All return all model info
func (mc *ModelCache) All() map[string]*ModelInfo {
	m := make(map[string]*ModelInfo, len(mc.cache))
	for k, v := range mc.cache {
		m[k] = v
	}
	return m
}

func (mc *ModelCache) Empty() bool {
	return len(mc.cache) == 0
}

func (mc *ModelCache) AllOrdered() []*ModelInfo {
	m := make([]*ModelInfo, 0, len(mc.orders))
	for _, table := range mc.orders {
		m = append(m, mc.cache[table])
	}
	return m
}

// Get model info by table name
func (mc *ModelCache) Get(table string) (mi *ModelInfo, ok bool) {
	mi, ok = mc.cache[table]
	return
}

// GetByFullName model info by full name
func (mc *ModelCache) GetByFullName(name string) (mi *ModelInfo, ok bool) {
	mi, ok = mc.cacheByFullName[name]
	return
}

func (mc *ModelCache) GetByMd(md interface{}) (*ModelInfo, bool) {
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	name := GetFullName(typ)
	return mc.GetByFullName(name)
}

// Set model info to collection
func (mc *ModelCache) Set(table string, mi *ModelInfo) *ModelInfo {
	mii := mc.cache[table]
	mc.cache[table] = mi
	mc.cacheByFullName[mi.FullName] = mi
	if mii == nil {
		mc.orders = append(mc.orders, table)
	}
	return mii
}

// Clean All model info.
func (mc *ModelCache) Clean() {
	mc.Lock()
	defer mc.Unlock()

	mc.orders = make([]string, 0)
	mc.cache = make(map[string]*ModelInfo)
	mc.cacheByFullName = make(map[string]*ModelInfo)
	mc.done = false
}

// Bootstrap Bootstrap for models
func (mc *ModelCache) Bootstrap() {
	mc.Lock()
	defer mc.Unlock()
	if mc.done {
		return
	}
	var (
		err    error
		models map[string]*ModelInfo
	)
	// Set rel and reverse model
	// RelManyToMany Set the relTable
	models = mc.All()
	for _, mi := range models {
		for _, fi := range mi.Fields.Columns {
			if fi.Rel || fi.Reverse {
				elm := fi.AddrValue.Type().Elem()
				if fi.FieldType == RelReverseMany || fi.FieldType == RelManyToMany {
					elm = elm.Elem()
				}
				// check the rel or reverse model already Register
				name := GetFullName(elm)
				mii, ok := mc.GetByFullName(name)
				if !ok || mii.Pkg != elm.PkgPath() {
					err = fmt.Errorf("can not find rel in field `%s`, `%s` may be miss Register", fi.FullName, elm.String())
					goto end
				}
				fi.RelModelInfo = mii

				switch fi.FieldType {
				case RelManyToMany:
					if fi.RelThrough != "" {
						if i := strings.LastIndex(fi.RelThrough, "."); i != -1 && len(fi.RelThrough) > (i+1) {
							pn := fi.RelThrough[:i]
							rmi, ok := mc.GetByFullName(fi.RelThrough)
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
						i := NewM2MModelInfo(mi, mii)
						if fi.RelTable != "" {
							i.Table = fi.RelTable
						}
						if v := mc.Set(i.Table, i); v != nil {
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
	models = mc.All()
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
					ffi := new(FieldInfo)
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

	models = mc.All()
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

	models = mc.All()
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

// Register Register models to model cache
func (mc *ModelCache) Register(prefixOrSuffixStr string, prefixOrSuffix bool, models ...interface{}) (err error) {
	for _, model := range models {
		val := reflect.ValueOf(model)
		typ := reflect.Indirect(val).Type()

		if val.Kind() != reflect.Ptr {
			err = fmt.Errorf("<orm.RegisterModel> cannot use non-ptr model struct `%s`", GetFullName(typ))
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
		table := GetTableName(val)

		if prefixOrSuffixStr != "" {
			if prefixOrSuffix {
				table = prefixOrSuffixStr + table
			} else {
				table = table + prefixOrSuffixStr
			}
		}

		// models's fullname is pkgpath + struct name
		name := GetFullName(typ)
		if _, ok := mc.GetByFullName(name); ok {
			err = fmt.Errorf("<orm.RegisterModel> model `%s` repeat Register, must be unique\n", name)
			return
		}

		if _, ok := mc.Get(table); ok {
			return nil
		}

		mi := NewModelInfo(val)
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

		mc.Set(table, mi)
	}
	return
}
