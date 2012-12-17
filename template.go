package beego

//@todo add template funcs

import (
	//"fmt"
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"strings"
	"time"
)

var beegoTplFuncMap template.FuncMap

func init() {
	beegoTplFuncMap = make(template.FuncMap)
	beegoTplFuncMap["markdown"] = MarkDown
	beegoTplFuncMap["dateformat"] = DateFormat
	beegoTplFuncMap["compare"] = Compare
}

// MarkDown parses a string in MarkDown format and returns HTML. Used by the template parser as "markdown"
func MarkDown(raw string) (output template.HTML) {
	input := []byte(raw)
	bOutput := blackfriday.MarkdownBasic(input)
	output = template.HTML(string(bOutput))
	return
}

// DateFormat takes a time and a layout string and returns a string with the formatted date. Used by the template parser as "dateformat"
func DateFormat(t time.Time, layout string) (datestring string) {
	datestring = t.Format(layout)
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

// AddFuncMap let user to register a func in the template
func AddFuncMap(key string, funname interface{}) error {
	if _, ok := beegoTplFuncMap[key]; ok {
		beegoTplFuncMap[key] = funname
		return nil
	}
	return errors.New("funcmap already has the key")
}
