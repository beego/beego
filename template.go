package beego

//@todo add template funcs

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
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
	beegoTplFuncMap["dateformat"] = DateFormat
	beegoTplFuncMap["date"] = Date
	beegoTplFuncMap["compare"] = Compare
	beegoTplFuncMap["substr"] = Substr
	beegoTplFuncMap["html2str"] = Html2str
	beegoTplFuncMap["str2html"] = Str2html
	beegoTplFuncMap["htmlquote"] = Htmlquote
	beegoTplFuncMap["htmlunquote"] = Htmlunquote
	beegoTplFuncMap["renderform"] = RenderForm

	// go1.2 added template funcs
	// Comparisons
	beegoTplFuncMap["eq"] = eq // ==
	beegoTplFuncMap["ge"] = ge // >=
	beegoTplFuncMap["gt"] = gt // >
	beegoTplFuncMap["le"] = le // <=
	beegoTplFuncMap["lt"] = lt // <
	beegoTplFuncMap["ne"] = ne // !=
}

// AddFuncMap let user to register a func in the template
func AddFuncMap(key string, funname interface{}) error {
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
	if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
		return nil
	}
	if !HasTemplateEXt(paths) {
		return nil
	}

	replace := strings.NewReplacer("\\", "/")
	a := []byte(paths)
	a = a[len([]byte(self.root)):]
	file := strings.TrimLeft(replace.Replace(string(a)), "/")
	subdir := filepath.Dir(file)
	if _, ok := self.files[subdir]; ok {
		self.files[subdir] = append(self.files[subdir], file)
	} else {
		m := make([]string, 1)
		m[0] = file
		self.files[subdir] = m
	}

	return nil
}

func HasTemplateEXt(paths string) bool {
	for _, v := range BeeTemplateExt {
		if strings.HasSuffix(paths, "."+v) {
			return true
		}
	}
	return false
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
			return nil
		} else {
			return errors.New("dir open err")
		}
	}
	self := &templatefile{
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
	for _, v := range self.files {
		for _, file := range v {
			t, err := getTemplate(self.root, file, v...)
			if err != nil {
				Trace("parse template err:", file, err)
			} else {
				BeeTemplates[file] = t
			}
		}
	}
	return nil
}

func getTplDeep(root, file string, t *template.Template) (*template.Template, [][]string, error) {
	fileabspath := filepath.Join(root, file)
	data, err := ioutil.ReadFile(fileabspath)
	if err != nil {
		return nil, [][]string{}, err
	}
	t, err = t.New(file).Parse(string(data))
	if err != nil {
		return nil, [][]string{}, err
	}
	reg := regexp.MustCompile("{{[ ]*template[ ]+\"([^\"]+)\"")
	allsub := reg.FindAllStringSubmatch(string(data), -1)
	for _, m := range allsub {
		if len(m) == 2 {
			tlook := t.Lookup(m[1])
			if tlook != nil {
				continue
			}
			if !HasTemplateEXt(m[1]) {
				continue
			}
			if e, _ := FileExists(filepath.Join(root, m[1])); e {
				t, _, err = getTplDeep(root, m[1], t)
				if err != nil {
					return nil, [][]string{}, err
				}
			} else {
				relativefile := filepath.Join(filepath.Dir(file), m[1])
				if e, _ := FileExists(relativefile); e {
					t, _, err = getTplDeep(root, m[1], t)
					if err != nil {
						return nil, [][]string{}, err
					}
				} else {
					panic("can't find template file" + m[1])
				}
			}

		}
	}
	return t, allsub, nil
}

func getTemplate(root, file string, others ...string) (t *template.Template, err error) {
	t = template.New(file).Delims(TemplateLeft, TemplateRight).Funcs(beegoTplFuncMap)
	var submods [][]string
	t, submods, err = getTplDeep(root, file, t)
	if err != nil {
		return nil, err
	}
	t, err = _getTemplate(t, root, submods, others...)

	if err != nil {
		return nil, err
	}
	return
}

func _getTemplate(t0 *template.Template, root string, submods [][]string, others ...string) (t *template.Template, err error) {
	t = t0
	for _, m := range submods {
		if len(m) == 2 {
			templ := t.Lookup(m[1])
			if templ != nil {
				continue
			}
			//first check filename
			for _, otherfile := range others {
				if otherfile == m[1] {
					var submods1 [][]string
					t, submods1, err = getTplDeep(root, otherfile, t)
					if err != nil {
						Trace("template parse file err:", err)
					} else if submods1 != nil && len(submods1) > 0 {
						t, err = _getTemplate(t, root, submods1, others...)
					}
					break
				}
			}
			//second check define
			for _, otherfile := range others {
				fileabspath := filepath.Join(root, otherfile)
				data, err := ioutil.ReadFile(fileabspath)
				if err != nil {
					continue
				}
				reg := regexp.MustCompile("{{[ ]*define[ ]+\"([^\"]+)\"")
				allsub := reg.FindAllStringSubmatch(string(data), -1)
				for _, sub := range allsub {
					if len(sub) == 2 && sub[1] == m[1] {
						var submods1 [][]string
						t, submods1, err = getTplDeep(root, otherfile, t)
						if err != nil {
							Trace("template parse file err:", err)
						} else if submods1 != nil && len(submods1) > 0 {
							t, err = _getTemplate(t, root, submods1, others...)
						}
						break
					}
				}
			}
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
	integerKind
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
		if k1 != k2 {
			return false, errBadComparison
		}
		truth := false
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
	if k1 != k2 {
		return false, errBadComparison
	}
	truth := false
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
		panic("invalid kind")
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

// go1.2 added template funcs. end
