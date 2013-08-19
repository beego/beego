package orm

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func registerModel(model interface{}, prefix string) {
	val := reflect.ValueOf(model)
	ind := reflect.Indirect(val)
	typ := ind.Type()

	if val.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("<orm.RegisterModel> cannot use non-ptr model struct `%s`", getFullName(typ)))
	}

	table := getTableName(val)

	if prefix != "" {
		table = prefix + table
	}

	name := getFullName(typ)
	if _, ok := modelCache.getByFN(name); ok {
		fmt.Printf("<orm.RegisterModel> model `%s` repeat register, must be unique\n", name)
		os.Exit(2)
	}

	if _, ok := modelCache.get(table); ok {
		fmt.Printf("<orm.RegisterModel> table name `%s` repeat register, must be unique\n", table)
		os.Exit(2)
	}

	info := newModelInfo(val)

	if info.fields.pk == nil {
	outFor:
		for _, fi := range info.fields.fieldsDB {
			if fi.name == "Id" {
				if fi.sf.Tag.Get(defaultStructTagName) == "" {
					switch fi.addrValue.Elem().Kind() {
					case reflect.Int, reflect.Int64, reflect.Int32:
						fi.auto = true
						fi.pk = true
						info.fields.pk = fi
						break outFor
					}
				}
			}
		}

		if info.fields.pk == nil {
			fmt.Printf("<orm.RegisterModel> `%s` need a primary key field\n", name)
			os.Exit(2)
		}

	}

	info.table = table
	info.pkg = typ.PkgPath()
	info.model = model
	info.manual = true

	modelCache.set(table, info)
}

func bootStrap() {
	if modelCache.done {
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

	models = modelCache.all()
	for _, mi := range models {
		for _, fi := range mi.fields.columns {
			if fi.rel || fi.reverse {
				elm := fi.addrValue.Type().Elem()
				switch fi.fieldType {
				case RelReverseMany, RelManyToMany:
					elm = elm.Elem()
				}

				name := getFullName(elm)
				mii, ok := modelCache.getByFN(name)
				if ok == false || mii.pkg != elm.PkgPath() {
					err = fmt.Errorf("can not found rel in field `%s`, `%s` may be miss register", fi.fullName, elm.String())
					goto end
				}
				fi.relModelInfo = mii

				switch fi.fieldType {
				case RelManyToMany:
					if fi.relThrough != "" {
						msg := fmt.Sprintf("field `%s` wrong rel_through value `%s`", fi.fullName, fi.relThrough)
						if i := strings.LastIndex(fi.relThrough, "."); i != -1 && len(fi.relThrough) > (i+1) {
							pn := fi.relThrough[:i]
							mn := fi.relThrough[i+1:]
							tn := snakeString(mn)
							rmi, ok := modelCache.get(tn)
							if ok == false || pn != rmi.pkg {
								err = errors.New(msg + " cannot find table")
								goto end
							}

							fi.relThroughModelInfo = rmi
							fi.relTable = rmi.table

						} else {
							err = errors.New(msg)
							goto end
						}
						err = nil
					} else {
						i := newM2MModelInfo(mi, mii)
						if fi.relTable != "" {
							i.table = fi.relTable
						}

						if v := modelCache.set(i.table, i); v != nil {
							err = fmt.Errorf("the rel table name `%s` already registered, cannot be use, please change one", fi.relTable)
							goto end
						}
						fi.relTable = i.table
						fi.relThroughModelInfo = i
					}
				}
			}
		}
	}

	models = modelCache.all()
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
				if inModel == false {
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
					if rmi.fields.Add(ffi) == false {
						added := false
						for cnt := 0; cnt < 5; cnt++ {
							ffi.name = fmt.Sprintf("%s%d", mi.name, cnt)
							ffi.column = ffi.name
							ffi.fullName = rmi.fullName + "." + ffi.name
							if added = rmi.fields.Add(ffi); added {
								break
							}
						}
						if added == false {
							panic(fmt.Sprintf("cannot generate auto reverse field info `%s` to `%s`", fi.fullName, ffi.fullName))
						}
					}
				}
			}
		}
	}

	for _, mi := range models {
		if fields, ok := mi.fields.fieldsByType[RelReverseOne]; ok {
			for _, fi := range fields {
				found := false
			mForA:
				for _, ffi := range fi.relModelInfo.fields.fieldsByType[RelOneToOne] {
					if ffi.relModelInfo == mi {
						found = true
						fi.reverseField = ffi.name
						fi.reverseFieldInfo = ffi
						break mForA
					}
				}
				if found == false {
					err = fmt.Errorf("reverse field `%s` not found in model `%s`", fi.fullName, fi.relModelInfo.fullName)
					goto end
				}
			}
		}
		if fields, ok := mi.fields.fieldsByType[RelReverseMany]; ok {
			for _, fi := range fields {
				found := false
			mForB:
				for _, ffi := range fi.relModelInfo.fields.fieldsByType[RelForeignKey] {
					if ffi.relModelInfo == mi {
						found = true
						fi.reverseField = ffi.name
						fi.reverseFieldInfo = ffi
						break mForB
					}
				}
				if found == false {
				mForC:
					for _, ffi := range fi.relModelInfo.fields.fieldsByType[RelManyToMany] {
						if ffi.relModelInfo == mi {
							found = true
							fi.reverseField = ffi.name
							fi.reverseFieldInfo = ffi
							break mForC
						}
					}
				}
				if found == false {
					err = fmt.Errorf("reverse field `%s` not found in model `%s`", fi.fullName, fi.relModelInfo.fullName)
					goto end
				}
			}
		}
	}

end:
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func RegisterModel(models ...interface{}) {
	if modelCache.done {
		panic(fmt.Errorf("RegisterModel must be run before BootStrap"))
	}

	for _, model := range models {
		registerModel(model, "")
	}
}

// register model with a prefix
func RegisterModelWithPrefix(prefix string, models ...interface{}) {
	if modelCache.done {
		panic(fmt.Errorf("RegisterModel must be run before BootStrap"))
	}

	for _, model := range models {
		registerModel(model, prefix)
	}
}

func BootStrap() {
	if modelCache.done {
		return
	}

	modelCache.Lock()
	defer modelCache.Unlock()
	bootStrap()
	modelCache.done = true
}
