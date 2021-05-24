package bean

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// the mock object must be pointer of struct
// the element in mock object can be slices, structures, basic data types, pointers and interface
func Mock(v interface{}) (err error) {
	pv := reflect.ValueOf(v)
	// the input must be pointer of struct
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		err = fmt.Errorf("not a pointer of struct")
		return
	}
	err = mock(pv)
	return
}

func mock(pv reflect.Value) (err error) {
	pt := pv.Type()
	for i := 0; i < pt.Elem().NumField(); i++ {
		ptt := pt.Elem().Field(i)
		pvv := pv.Elem().FieldByName(ptt.Name)
		if !pvv.CanSet() {
			continue
		}
		kt := ptt.Type.Kind()
		tagValue := ptt.Tag.Get("mock")
		switch kt {
		case reflect.Map:
			continue
		case reflect.Interface:
			if pvv.IsNil() { // when interface is nil,can not sure the type
				continue
			}
			pvv.Set(reflect.New(pvv.Elem().Type().Elem()))
			err = mock(pvv.Elem())
		case reflect.Ptr:
			err = mockPtr(pvv, ptt.Type.Elem())
		case reflect.Struct:
			err = mock(pvv.Addr())
		case reflect.Array, reflect.Slice:
			err = mockSlice(tagValue, pvv)
		case reflect.String:
			pvv.SetString(tagValue)
		case reflect.Bool:
			err = mockBool(tagValue, pvv)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value, e := strconv.ParseInt(tagValue, 10, 64)
			if e != nil || pvv.OverflowInt(value) {
				err = fmt.Errorf("the value:%s is invalid", tagValue)
			}
			pvv.SetInt(value)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value, e := strconv.ParseUint(tagValue, 10, 64)
			if e != nil || pvv.OverflowUint(value) {
				err = fmt.Errorf("the value:%s is invalid", tagValue)
			}
			pvv.SetUint(value)
		case reflect.Float32, reflect.Float64:
			value, e := strconv.ParseFloat(tagValue, pvv.Type().Bits())
			if e != nil || pvv.OverflowFloat(value) {
				err = fmt.Errorf("the value:%s is invalid", tagValue)
			}
			pvv.SetFloat(value)
		default:
		}
		if err != nil {
			return
		}
	}
	return
}

// mock slice value
func mockSlice(tagValue string, pvv reflect.Value) (err error) {
	if len(tagValue) == 0 {
		return
	}
	sliceMetas := strings.Split(tagValue, ":")
	if len(sliceMetas) != 2 || sliceMetas[0] != "length" {
		err = fmt.Errorf("the value:%s is invalid", tagValue)
		return
	}
	length, e := strconv.Atoi(sliceMetas[1])
	if e != nil {
		return e
	}

	sliceType := reflect.SliceOf(pvv.Type().Elem()) // get slice type
	itemType := sliceType.Elem()                    // get the type of item in slice
	value := reflect.MakeSlice(sliceType, 0, length)
	newSliceValue := make([]reflect.Value, 0, length)
	for k := 0; k < length; k++ {
		itemValue := reflect.New(itemType).Elem()
		// if item in slice is struct or pointer,must set zero value
		switch itemType.Kind() {
		case reflect.Struct:
			err = mock(itemValue.Addr())
		case reflect.Ptr:
			if itemValue.IsNil() {
				itemValue.Set(reflect.New(itemType.Elem()))
				if e := mock(itemValue); e != nil {
					return e
				}
			}
		}
		newSliceValue = append(newSliceValue, itemValue)
		if err != nil {
			return
		}
	}
	value = reflect.Append(value, newSliceValue...)
	pvv.Set(value)
	return
}

// mock bool value
func mockBool(tagValue string, pvv reflect.Value) (err error) {
	switch tagValue {
	case "true":
		pvv.SetBool(true)
	case "false":
		pvv.SetBool(false)
	default:
		err = fmt.Errorf("the value:%s is invalid", tagValue)
	}
	return
}

// mock pointer
func mockPtr(pvv reflect.Value, ptt reflect.Type) (err error) {
	if pvv.IsNil() {
		pvv.Set(reflect.New(ptt)) // must set nil value to zero value
	}
	err = mock(pvv)
	return
}
