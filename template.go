package beego

import (
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	beegoTplFuncMap template.FuncMap
	BeeTemplates    map[string]*template.Template
	BeeTemplateExt  []string
)

func init() {
	BeeTemplates = make(map[string]*template.Template)
	beegoTplFuncMap = make(template.FuncMap)
	BeeTemplateExt = make([]string, 0)
	BeeTemplateExt = append(BeeTemplateExt, "tpl", "html")
	beegoTplFuncMap["markdown"] = MarkDown
	beegoTplFuncMap["dateformat"] = DateFormat
	beegoTplFuncMap["date"] = Date
	beegoTplFuncMap["compare"] = Compare
	beegoTplFuncMap["substr"] = Substr
	beegoTplFuncMap["html2str"] = Html2str
	beegoTplFuncMap["str2html"] = Str2html
	beegoTplFuncMap["htmlquote"] = Htmlquote
	beegoTplFuncMap["htmlunquote"] = Htmlunquote
}

// MarkDown parses a string in MarkDown format and returns HTML. Used by the template parser as "markdown"
func MarkDown(raw string) (output template.HTML) {
	input := []byte(raw)
	bOutput := blackfriday.MarkdownBasic(input)
	output = template.HTML(string(bOutput))
	return
}

func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt) - 1
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

// AddFuncMap let user to register a func in the template
func AddFuncMap(key string, funname interface{}) error {
	if _, ok := beegoTplFuncMap[key]; ok {
		return errors.New("funcmap already has the key")
	}
	beegoTplFuncMap[key] = funname
	return nil
}

type templatefile struct {
	root  string
	files map[string][]string
}

func (self *templatefile) visit(paths string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	} else if (f.Mode() & os.ModeSymlink) > 0 {
		return nil
	} else {
		hasExt := false
		for _, v := range BeeTemplateExt {
			if strings.HasSuffix(paths, v) {
				hasExt = true
				break
			}
		}
		if hasExt {
			replace := strings.NewReplacer("\\", "/")
			a := []byte(paths)
			a = a[len([]byte(self.root)):]
			subdir := path.Dir(strings.TrimLeft(replace.Replace(string(a)), "/"))
			if _, ok := self.files[subdir]; ok {
				self.files[subdir] = append(self.files[subdir], paths)
			} else {
				m := make([]string, 1)
				m[0] = paths
				self.files[subdir] = m
			}

		}
	}
	return nil
}

func AddTemplateExt(ext string) {
	for _, v := range BeeTemplateExt {
		if v == ext {
			return
		}
	}
	BeeTemplateExt = append(BeeTemplateExt, ext)
}

func BuildTemplate(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return err
		} else {
			return errors.New("dir open err")
		}
	}
	self := templatefile{
		root:  dir,
		files: make(map[string][]string),
	}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		return self.visit(path, f, err)
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return err
	}
	for k, v := range self.files {
		BeeTemplates[k] = template.Must(template.New("beegoTemplate" + k).Funcs(beegoTplFuncMap).ParseFiles(v...))
	}
	return nil
}
