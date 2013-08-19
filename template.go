package beego

//@todo add template funcs

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
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
	if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
		return nil
	}
	if !HasTemplateEXt(paths) {
		return nil
	}

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
		BeeTemplates[k] = template.Must(template.New("beegoTemplate"+k).Funcs(beegoTplFuncMap).ParseFiles(v...)).Delims(TemplatLeft, TemplatRight)
	}
	return nil
}
