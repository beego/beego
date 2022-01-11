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

package web

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

// Substr returns the substr from start to length.
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

// HTML2str returns escaping text convert from html.
func HTML2str(html string) string {
	re := regexp.MustCompile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllStringFunc(html, strings.ToLower)

	// remove STYLE
	re = regexp.MustCompile(`\<style[\S\s]+?\</style\>`)
	html = re.ReplaceAllString(html, "")

	// remove SCRIPT
	re = regexp.MustCompile(`\<script[\S\s]+?\</script\>`)
	html = re.ReplaceAllString(html, "")

	re = regexp.MustCompile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllString(html, "\n")

	re = regexp.MustCompile(`\s{2,}`)
	html = re.ReplaceAllString(html, "\n")

	return strings.TrimSpace(html)
}

// DateFormat takes a time and a layout string and returns a string with the formatted date. Used by the template parser as "dateformat"
func DateFormat(t time.Time, layout string) (datestring string) {
	datestring = t.Format(layout)
	return
}

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", // A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// DateParse Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

// Date takes a PHP like date func to Go's time format.
func Date(t time.Time, format string) string {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

// Compare is a quick and dirty comparison function. It will convert whatever you give it to strings and see if the two values are equal.
// Whitespace is trimmed. Used by the template parser as "eq".
func Compare(a, b interface{}) (equal bool) {
	equal = false
	if strings.TrimSpace(fmt.Sprintf("%v", a)) == strings.TrimSpace(fmt.Sprintf("%v", b)) {
		equal = true
	}
	return
}

// CompareNot !Compare
func CompareNot(a, b interface{}) (equal bool) {
	return !Compare(a, b)
}

// NotNil the same as CompareNot
func NotNil(a interface{}) (isNil bool) {
	return CompareNot(a, nil)
}

// GetConfig get the Appconfig
func GetConfig(returnType, key string, defaultVal interface{}) (value interface{}, err error) {
	switch returnType {
	case "String":
		value, err = AppConfig.String(key)
	case "Bool":
		value, err = AppConfig.Bool(key)
	case "Int":
		value, err = AppConfig.Int(key)
	case "Int64":
		value, err = AppConfig.Int64(key)
	case "Float":
		value, err = AppConfig.Float(key)
	case "DIY":
		value, err = AppConfig.DIY(key)
	default:
		err = errors.New("config keys must be of type String, Bool, Int, Int64, Float, or DIY")
	}

	if err != nil {
		if reflect.TypeOf(returnType) != reflect.TypeOf(defaultVal) {
			err = errors.New("defaultVal type does not match returnType")
		} else {
			value, err = defaultVal, nil
		}
	} else if reflect.TypeOf(value).Kind() == reflect.String {
		if value == "" {
			if reflect.TypeOf(defaultVal).Kind() != reflect.String {
				err = errors.New("defaultVal type must be a String if the returnType is a String")
			} else {
				value = defaultVal.(string)
			}
		}
	}

	return
}

// Str2html Convert string to template.HTML type.
func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

// Htmlquote returns quoted html string.
func Htmlquote(text string) string {
	// HTML编码为实体符号
	/*
	   Encodes `text` for raw use in HTML.
	       >>> htmlquote("<'&\\">")
	       '&lt;&#39;&amp;&quot;&gt;'
	*/

	text = html.EscapeString(text)
	text = strings.NewReplacer(
		`“`, "&ldquo;",
		`”`, "&rdquo;",
		` `, "&nbsp;",
	).Replace(text)

	return strings.TrimSpace(text)
}

// Htmlunquote returns unquoted html string.
func Htmlunquote(text string) string {
	// 实体符号解释为HTML
	/*
	   Decodes `text` that's HTML quoted.
	       >>> htmlunquote('&lt;&#39;&amp;&quot;&gt;')
	       '<\\'&">'
	*/

	text = html.UnescapeString(text)

	return strings.TrimSpace(text)
}

// URLFor returns url string with another registered controller handler with params.
//	usage:
//
//	URLFor(".index")
//	print URLFor("index")
//  router /login
//	print URLFor("login")
//	print URLFor("login", "next","/"")
//  router /profile/:username
//	print UrlFor("profile", ":username","John Doe")
//	result:
//	/
//	/login
//	/login?next=/
//	/user/John%20Doe
//
//  more detail http://beego.vip/docs/mvc/controller/urlbuilding.md
func URLFor(endpoint string, values ...interface{}) string {
	return BeeApp.Handlers.URLFor(endpoint, values...)
}

// AssetsJs returns script tag with src string.
func AssetsJs(text string) template.HTML {
	text = "<script src=\"" + text + "\"></script>"

	return template.HTML(text)
}

// AssetsCSS returns stylesheet link tag with src string.
func AssetsCSS(text string) template.HTML {
	text = "<link href=\"" + text + "\" rel=\"stylesheet\" />"

	return template.HTML(text)
}

// ParseForm will parse form values to struct via tag.
func ParseForm(form url.Values, obj interface{}) error {
	return context.ParseForm(form, obj)
}

var unKind = map[reflect.Kind]bool{
	reflect.Uintptr:       true,
	reflect.Complex64:     true,
	reflect.Complex128:    true,
	reflect.Array:         true,
	reflect.Chan:          true,
	reflect.Func:          true,
	reflect.Map:           true,
	reflect.Ptr:           true,
	reflect.Slice:         true,
	reflect.Struct:        true,
	reflect.UnsafePointer: true,
}

// RenderForm will render object to form html.
// obj must be a struct pointer.
// nolint
func RenderForm(obj interface{}) template.HTML {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if objT.Kind() != reflect.Ptr || objT.Elem().Kind() != reflect.Struct {
		return template.HTML("")
	}
	objT = objT.Elem()
	objV = objV.Elem()

	var raw []string
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() || unKind[fieldV.Kind()] {
			continue
		}

		fieldT := objT.Field(i)

		label, name, fType, id, class, ignored, required := parseFormTag(fieldT)
		if ignored {
			continue
		}

		raw = append(raw, renderFormField(label, name, fType, fieldV.Interface(), id, class, required))
	}
	return template.HTML(strings.Join(raw, "</br>"))
}

// renderFormField returns a string containing HTML of a single form field.
func renderFormField(label, name, fType string, value interface{}, id string, class string, required bool) string {
	if id != "" {
		id = " id=\"" + id + "\""
	}

	if class != "" {
		class = " class=\"" + class + "\""
	}

	requiredString := ""
	if required {
		requiredString = " required"
	}

	if isValidForInput(fType) {
		return fmt.Sprintf(`%v<input%v%v name="%v" type="%v" value="%v"%v>`, label, id, class, name, fType, value, requiredString)
	}

	return fmt.Sprintf(`%v<%v%v%v name="%v"%v>%v</%v>`, label, fType, id, class, name, requiredString, value, fType)
}

// isValidForInput checks if fType is a valid value for the `type` property of an HTML input element.
func isValidForInput(fType string) bool {
	validInputTypes := strings.Fields("text password checkbox radio submit reset hidden image file button search email url tel number range date month week time datetime datetime-local color")
	for _, validType := range validInputTypes {
		if fType == validType {
			return true
		}
	}
	return false
}

// parseFormTag takes the stuct-tag of a StructField and parses the `form` value.
// returned are the form label, name-property, type and whether the field should be ignored.
func parseFormTag(fieldT reflect.StructField) (label, name, fType string, id string, class string, ignored bool, required bool) {
	tags := strings.Split(fieldT.Tag.Get("form"), ",")
	label = fieldT.Name + ": "
	name = fieldT.Name
	fType = "text"
	ignored = false
	id = fieldT.Tag.Get("id")
	class = fieldT.Tag.Get("class")

	required = false
	requiredField := fieldT.Tag.Get("required")
	if requiredField != "-" && requiredField != "" {
		required, _ = strconv.ParseBool(requiredField)
	}

	switch len(tags) {
	case 1:
		if tags[0] == "-" {
			ignored = true
		}
		if len(tags[0]) > 0 {
			name = tags[0]
		}
	case 2:
		if len(tags[0]) > 0 {
			name = tags[0]
		}
		if len(tags[1]) > 0 {
			fType = tags[1]
		}
	case 3:
		if len(tags[0]) > 0 {
			name = tags[0]
		}
		if len(tags[1]) > 0 {
			fType = tags[1]
		}
		if len(tags[2]) > 0 {
			label = tags[2]
		}
	}

	return
}

// go1.2 added template funcs. begin
var (
	errBadComparisonType = errors.New("invalid type for comparison")
	errBadComparison     = errors.New("incompatible types for comparison")
	errNoComparison      = errors.New("missing argument for comparison")
)

type kind int

const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	stringKind
	uintKind
)

func basicKind(v reflect.Value) (kind, error) {
	switch v.Kind() {
	case reflect.Bool:
		return boolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind, nil
	case reflect.Float32, reflect.Float64:
		return floatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return complexKind, nil
	case reflect.String:
		return stringKind, nil
	}
	return invalidKind, errBadComparisonType
}

// eq evaluates the comparison a == b || a == c || ...
func eq(arg1 interface{}, arg2 ...interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := basicKind(v1)
	if err != nil {
		return false, err
	}
	if len(arg2) == 0 {
		return false, errNoComparison
	}
	for _, arg := range arg2 {
		v2 := reflect.ValueOf(arg)
		k2, err := basicKind(v2)
		if err != nil {
			return false, err
		}
		truth := false
		if k1 != k2 {
			// Special case: Can compare integer values regardless of type's sign.
			switch {
			case k1 == intKind && k2 == uintKind:
				truth = v1.Int() >= 0 && uint64(v1.Int()) == v2.Uint()
			case k1 == uintKind && k2 == intKind:
				truth = v2.Int() >= 0 && v1.Uint() == uint64(v2.Int())
			default:
				return false, errBadComparison
			}
			if truth {
				return true, nil
			} else {
				return false, nil
			}
		}
		switch k1 {
		case boolKind:
			truth = v1.Bool() == v2.Bool()
		case complexKind:
			truth = v1.Complex() == v2.Complex()
		case floatKind:
			truth = v1.Float() == v2.Float()
		case intKind:
			truth = v1.Int() == v2.Int()
		case stringKind:
			truth = v1.String() == v2.String()
		case uintKind:
			truth = v1.Uint() == v2.Uint()
		default:
			panic("invalid kind")
		}
		if truth {
			return true, nil
		}
	}
	return false, nil
}

// ne evaluates the comparison a != b.
func ne(arg1, arg2 interface{}) (bool, error) {
	// != is the inverse of ==.
	equal, err := eq(arg1, arg2)
	return !equal, err
}

// lt evaluates the comparison a < b.
func lt(arg1, arg2 interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := basicKind(v1)
	if err != nil {
		return false, err
	}
	v2 := reflect.ValueOf(arg2)
	k2, err := basicKind(v2)
	if err != nil {
		return false, err
	}
	truth := false
	if k1 != k2 {
		// Special case: Can compare integer values regardless of type's sign.
		switch {
		case k1 == intKind && k2 == uintKind:
			truth = v1.Int() < 0 || uint64(v1.Int()) < v2.Uint()
		case k1 == uintKind && k2 == intKind:
			truth = v2.Int() >= 0 && v1.Uint() < uint64(v2.Int())
		default:
			return false, errBadComparison
		}
	} else {
		switch k1 {
		case boolKind, complexKind:
			return false, errBadComparisonType
		case floatKind:
			truth = v1.Float() < v2.Float()
		case intKind:
			truth = v1.Int() < v2.Int()
		case stringKind:
			truth = v1.String() < v2.String()
		case uintKind:
			truth = v1.Uint() < v2.Uint()
		default:
			return false, errBadComparisonType
		}
	}
	return truth, nil
}

// le evaluates the comparison <= b.
func le(arg1, arg2 interface{}) (bool, error) {
	// <= is < or ==.
	lessThan, err := lt(arg1, arg2)
	if lessThan || err != nil {
		return lessThan, err
	}
	return eq(arg1, arg2)
}

// gt evaluates the comparison a > b.
func gt(arg1, arg2 interface{}) (bool, error) {
	// > is the inverse of <=.
	lessOrEqual, err := le(arg1, arg2)
	if err != nil {
		return false, err
	}
	return !lessOrEqual, nil
}

// ge evaluates the comparison a >= b.
func ge(arg1, arg2 interface{}) (bool, error) {
	// >= is the inverse of <.
	lessThan, err := lt(arg1, arg2)
	if err != nil {
		return false, err
	}
	return !lessThan, nil
}

// MapGet getting value from map by keys
// usage:
// Data["m"] = M{
//     "a": 1,
//     "1": map[string]float64{
//         "c": 4,
//     },
// }
//
// {{ map_get m "a" }} // return 1
// {{ map_get m 1 "c" }} // return 4
func MapGet(arg1 interface{}, arg2 ...interface{}) (interface{}, error) {
	arg1Type := reflect.TypeOf(arg1)
	arg1Val := reflect.ValueOf(arg1)

	if arg1Type.Kind() == reflect.Map && len(arg2) > 0 {
		// check whether arg2[0] type equals to arg1 key type
		// if they are different, make conversion
		arg2Val := reflect.ValueOf(arg2[0])
		arg2Type := reflect.TypeOf(arg2[0])
		if arg2Type.Kind() != arg1Type.Key().Kind() {
			// convert arg2Value to string
			var arg2ConvertedVal interface{}
			arg2String := fmt.Sprintf("%v", arg2[0])

			// convert string representation to any other type
			switch arg1Type.Key().Kind() {
			case reflect.Bool:
				arg2ConvertedVal, _ = strconv.ParseBool(arg2String)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				arg2ConvertedVal, _ = strconv.ParseInt(arg2String, 0, 64)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				arg2ConvertedVal, _ = strconv.ParseUint(arg2String, 0, 64)
			case reflect.Float32, reflect.Float64:
				arg2ConvertedVal, _ = strconv.ParseFloat(arg2String, 64)
			case reflect.String:
				arg2ConvertedVal = arg2String
			default:
				arg2ConvertedVal = arg2Val.Interface()
			}
			arg2Val = reflect.ValueOf(arg2ConvertedVal)
		}

		storedVal := arg1Val.MapIndex(arg2Val)

		if storedVal.IsValid() {
			var result interface{}

			switch arg1Type.Elem().Kind() {
			case reflect.Bool:
				result = storedVal.Bool()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result = storedVal.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				result = storedVal.Uint()
			case reflect.Float32, reflect.Float64:
				result = storedVal.Float()
			case reflect.String:
				result = storedVal.String()
			default:
				result = storedVal.Interface()
			}

			// if there is more keys, handle this recursively
			if len(arg2) > 1 {
				return MapGet(result, arg2[1:]...)
			}
			return result, nil
		}
		return nil, nil

	}
	return nil, nil
}
