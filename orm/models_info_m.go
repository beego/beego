package orm

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

type modelInfo struct {
	pkg       string
	name      string
	fullName  string
	table     string
	model     interface{}
	fields    *fields
	manual    bool
	addrField reflect.Value
	uniques   []string
	isThrough bool
}

func newModelInfo(val reflect.Value) (info *modelInfo) {
	var (
		err error
		fi  *fieldInfo
		sf  reflect.StructField
	)

	info = &modelInfo{}
	info.fields = newFields()

	ind := reflect.Indirect(val)
	typ := ind.Type()

	info.addrField = val

	info.name = typ.Name()
	info.fullName = getFullName(typ)

	for i := 0; i < ind.NumField(); i++ {
		field := ind.Field(i)
		sf = ind.Type().Field(i)
		fi, err = newFieldInfo(info, field, sf)

		if err != nil {
			if err == errSkipField {
				err = nil
				continue
			}
			break
		}

		added := info.fields.Add(fi)
		if added == false {
			err = errors.New(fmt.Sprintf("duplicate column name: %s", fi.column))
			break
		}

		if fi.pk {
			if info.fields.pk != nil {
				err = errors.New(fmt.Sprintf("one model must have one pk field only"))
				break
			} else {
				info.fields.pk = fi
			}
		}

		fi.fieldIndex = i
		fi.mi = info
		fi.inModel = true
	}

	if err != nil {
		fmt.Println(fmt.Errorf("field: %s.%s, %s", ind.Type(), sf.Name, err))
		os.Exit(2)
	}

	return
}

func newM2MModelInfo(m1, m2 *modelInfo) (info *modelInfo) {
	info = new(modelInfo)
	info.fields = newFields()
	info.table = m1.table + "_" + m2.table + "s"
	info.name = camelString(info.table)
	info.fullName = m1.pkg + "." + info.name

	fa := new(fieldInfo)
	f1 := new(fieldInfo)
	f2 := new(fieldInfo)
	fa.fieldType = TypeBigIntegerField
	fa.auto = true
	fa.pk = true
	fa.dbcol = true
	fa.name = "Id"
	fa.column = "id"
	fa.fullName = info.fullName + "." + fa.name

	f1.dbcol = true
	f2.dbcol = true
	f1.fieldType = RelForeignKey
	f2.fieldType = RelForeignKey
	f1.name = camelString(m1.table)
	f2.name = camelString(m2.table)
	f1.fullName = info.fullName + "." + f1.name
	f2.fullName = info.fullName + "." + f2.name
	f1.column = m1.table + "_id"
	f2.column = m2.table + "_id"
	f1.rel = true
	f2.rel = true
	f1.relTable = m1.table
	f2.relTable = m2.table
	f1.relModelInfo = m1
	f2.relModelInfo = m2
	f1.mi = info
	f2.mi = info

	info.fields.Add(fa)
	info.fields.Add(f1)
	info.fields.Add(f2)
	info.fields.pk = fa

	info.uniques = []string{f1.column, f2.column}
	return
}
