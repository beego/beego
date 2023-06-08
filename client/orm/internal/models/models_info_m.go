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

package models

import (
	"fmt"
	"os"
	"reflect"
)

// ModelInfo single model info
type ModelInfo struct {
	Manual    bool
	IsThrough bool
	Pkg       string
	Name      string
	FullName  string
	Table     string
	Model     interface{}
	Fields    *Fields
	AddrField reflect.Value // store the original struct value
	Uniques   []string
}

// NewModelInfo new model info
func NewModelInfo(val reflect.Value) (mi *ModelInfo) {
	mi = &ModelInfo{}
	mi.Fields = NewFields()
	ind := reflect.Indirect(val)
	mi.AddrField = val
	mi.Name = ind.Type().Name()
	mi.FullName = GetFullName(ind.Type())
	AddModelFields(mi, ind, "", []int{})
	return
}

// AddModelFields index: FieldByIndex returns the nested field corresponding to index
func AddModelFields(mi *ModelInfo, ind reflect.Value, mName string, index []int) {
	var (
		err error
		fi  *FieldInfo
		sf  reflect.StructField
	)

	for i := 0; i < ind.NumField(); i++ {
		field := ind.Field(i)
		sf = ind.Type().Field(i)
		// if the field is unexported skip
		if sf.PkgPath != "" {
			continue
		}
		// add anonymous struct Fields
		if sf.Anonymous {
			AddModelFields(mi, field, mName+"."+sf.Name, append(index, i))
			continue
		}

		fi, err = NewFieldInfo(mi, field, sf, mName)
		if err == errSkipField {
			err = nil
			continue
		} else if err != nil {
			break
		}
		// record current field index
		fi.FieldIndex = append(fi.FieldIndex, index...)
		fi.FieldIndex = append(fi.FieldIndex, i)
		fi.Mi = mi
		fi.InModel = true
		if !mi.Fields.Add(fi) {
			err = fmt.Errorf("duplicate column name: %s", fi.Column)
			break
		}
		if fi.Pk {
			if mi.Fields.Pk != nil {
				err = fmt.Errorf("one model must have one pk field only")
				break
			} else {
				mi.Fields.Pk = fi
			}
		}
	}

	if err != nil {
		fmt.Println(fmt.Errorf("field: %s.%s, %s", ind.Type(), sf.Name, err))
		os.Exit(2)
	}
}

// NewM2MModelInfo combine related model info to new model info.
// prepare for relation models query.
func NewM2MModelInfo(m1, m2 *ModelInfo) (mi *ModelInfo) {
	mi = new(ModelInfo)
	mi.Fields = NewFields()
	mi.Table = m1.Table + "_" + m2.Table + "s"
	mi.Name = CamelString(mi.Table)
	mi.FullName = m1.Pkg + "." + mi.Name

	fa := new(FieldInfo) // pk
	f1 := new(FieldInfo) // m1 table RelForeignKey
	f2 := new(FieldInfo) // m2 table RelForeignKey
	fa.FieldType = TypeBigIntegerField
	fa.Auto = true
	fa.Pk = true
	fa.DBcol = true
	fa.Name = "Id"
	fa.Column = "id"
	fa.FullName = mi.FullName + "." + fa.Name

	f1.DBcol = true
	f2.DBcol = true
	f1.FieldType = RelForeignKey
	f2.FieldType = RelForeignKey
	f1.Name = CamelString(m1.Table)
	f2.Name = CamelString(m2.Table)
	f1.FullName = mi.FullName + "." + f1.Name
	f2.FullName = mi.FullName + "." + f2.Name
	f1.Column = m1.Table + "_id"
	f2.Column = m2.Table + "_id"
	f1.Rel = true
	f2.Rel = true
	f1.RelTable = m1.Table
	f2.RelTable = m2.Table
	f1.RelModelInfo = m1
	f2.RelModelInfo = m2
	f1.Mi = mi
	f2.Mi = mi

	mi.Fields.Add(fa)
	mi.Fields.Add(f1)
	mi.Fields.Add(f2)
	mi.Fields.Pk = fa

	mi.Uniques = []string{f1.Column, f2.Column}
	return
}
