package beego

import (
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func webTime(t time.Time) string {
	ftime := t.Format(time.RFC1123)
	if strings.HasSuffix(ftime, "UTC") {
		ftime = ftime[0:len(ftime)-3] + "GMT"
	}
	return ftime
}

func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

// Html2str() returns escaping text convert from html
func Html2str(html string) string {
	src := string(html)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

// DateFormat takes a time and a layout string and returns a string with the formatted date. Used by the template parser as "dateformat"
func DateFormat(t time.Time, layout string) (datestring string) {
	datestring = t.Format(layout)
	return
}

// Date takes a PHP like date func to Go's time fomate
func Date(t time.Time, format string) (datestring string) {
	patterns := []string{
		// year
		"Y", "2006", // A full numeric representation of a year, 4 digits	Examples: 1999 or 2003
		"y", "06", //A two digit representation of a year	Examples: 99 or 03

		// month
		"m", "01", // Numeric representation of a month, with leading zeros	01 through 12
		"n", "1", // Numeric representation of a month, without leading zeros	1 through 12
		"M", "Jan", // A short textual representation of a month, three letters	Jan through Dec
		"F", "January", // A full textual representation of a month, such as January or March	January through December

		// day
		"d", "02", // Day of the month, 2 digits with leading zeros	01 to 31
		"j", "2", // Day of the month without leading zeros	1 to 31

		// week
		"D", "Mon", // A textual representation of a day, three letters	Mon through Sun
		"l", "Monday", // A full textual representation of the day of the week	Sunday through Saturday

		// time
		"g", "3", // 12-hour format of an hour without leading zeros	1 through 12
		"G", "15", // 24-hour format of an hour without leading zeros	0 through 23
		"h", "03", // 12-hour format of an hour with leading zeros	01 through 12
		"H", "15", // 24-hour format of an hour with leading zeros	00 through 23

		"a", "pm", // Lowercase Ante meridiem and Post meridiem	am or pm
		"A", "PM", // Uppercase Ante meridiem and Post meridiem	AM or PM

		"i", "04", // Minutes with leading zeros	00 to 59
		"s", "05", // Seconds, with leading zeros	00 through 59
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)
	datestring = t.Format(format)
	return
}

// Compare is a quick and dirty comparison function. It will convert whatever you give it to strings and see if the two values are equal.
// Whitespace is trimmed. Used by the template parser as "eq"
func Compare(a, b interface{}) (equal bool) {
	equal = false
	if strings.TrimSpace(fmt.Sprintf("%v", a)) == strings.TrimSpace(fmt.Sprintf("%v", b)) {
		equal = true
	}
	return
}

func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

func Htmlquote(src string) string {
	//HTML编码为实体符号
	/*
	   Encodes `text` for raw use in HTML.
	       >>> htmlquote("<'&\\">")
	       '&lt;&#39;&amp;&quot;&gt;'
	*/

	text := string(src)

	text = strings.Replace(text, "&", "&amp;", -1) // Must be done first!
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	text = strings.Replace(text, "'", "&#39;", -1)
	text = strings.Replace(text, "\"", "&quot;", -1)
	text = strings.Replace(text, "“", "&ldquo;", -1)
	text = strings.Replace(text, "”", "&rdquo;", -1)
	text = strings.Replace(text, " ", "&nbsp;", -1)

	return strings.TrimSpace(text)
}

func Htmlunquote(src string) string {
	//实体符号解释为HTML
	/*
	   Decodes `text` that's HTML quoted.
	       >>> htmlunquote('&lt;&#39;&amp;&quot;&gt;')
	       '<\\'&">'
	*/

	// strings.Replace(s, old, new, n)
	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换

	text := string(src)
	text = strings.Replace(text, "&nbsp;", " ", -1)
	text = strings.Replace(text, "&rdquo;", "”", -1)
	text = strings.Replace(text, "&ldquo;", "“", -1)
	text = strings.Replace(text, "&quot;", "\"", -1)
	text = strings.Replace(text, "&#39;", "'", -1)
	text = strings.Replace(text, "&gt;", ">", -1)
	text = strings.Replace(text, "&lt;", "<", -1)
	text = strings.Replace(text, "&amp;", "&", -1) // Must be done last!

	return strings.TrimSpace(text)
}

func inSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// parse form values to struct via tag
func ParseForm(form url.Values, obj interface{}) error {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if !(objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("%v must be  a struct pointer", obj)
	}
	objT = objT.Elem()
	objV = objV.Elem()

	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		fieldT := objT.Field(i)

		tags := strings.Split(fieldT.Tag.Get("form"), ",")
		var tag string
		if len(tags) == 0 || len(tags[0]) == 0 {
			tag = fieldT.Name
		} else {
			tag = tags[0]
		}

		value := form.Get(tag)
		if len(value) == 0 {
			continue
		}

		switch fieldT.Type.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldV.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			fieldV.SetFloat(x)
		case reflect.Interface:
			fieldV.Set(reflect.ValueOf(value))
		case reflect.String:
			fieldV.SetString(value)
		}
	}
	return nil
}

func stringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
}
