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

package beego

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/astaxie/beego/utils"
)

var (
	beegoTplFuncMap template.FuncMap
	// beego template caching map and supported template file extensions.
	BeeTemplates   map[string]*template.Template
	BeeTemplateExt []string
)

func init() {
	BeeTemplates = make(map[string]*template.Template)
	beegoTplFuncMap = make(template.FuncMap)
	BeeTemplateExt = make([]string, 0)
	BeeTemplateExt = append(BeeTemplateExt, "tpl", "html")
	beegoTplFuncMap["dateformat"] = DateFormat
	beegoTplFuncMap["date"] = Date
	beegoTplFuncMap["compare"] = Compare
	beegoTplFuncMap["compare_not"] = CompareNot
	beegoTplFuncMap["not_nil"] = NotNil
	beegoTplFuncMap["not_null"] = NotNil
	beegoTplFuncMap["substr"] = Substr
	beegoTplFuncMap["html2str"] = Html2str
	beegoTplFuncMap["str2html"] = Str2html
	beegoTplFuncMap["htmlquote"] = Htmlquote
	beegoTplFuncMap["htmlunquote"] = Htmlunquote
	beegoTplFuncMap["renderform"] = RenderForm
	beegoTplFuncMap["assets_js"] = AssetsJs
	beegoTplFuncMap["assets_css"] = AssetsCss
	beegoTplFuncMap["config"] = Config

	// go1.2 added template funcs
	// Comparisons
	beegoTplFuncMap["eq"] = eq // ==
	beegoTplFuncMap["ge"] = ge // >=
	beegoTplFuncMap["gt"] = gt // >
	beegoTplFuncMap["le"] = le // <=
	beegoTplFuncMap["lt"] = lt // <
	beegoTplFuncMap["ne"] = ne // !=

	beegoTplFuncMap["urlfor"] = UrlFor // !=
}

// AddFuncMap let user to register a func in the template.
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
	if !HasTemplateExt(paths) {
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

// return this path contains supported template extension of beego or not.
func HasTemplateExt(paths string) bool {
	for _, v := range BeeTemplateExt {
		if strings.HasSuffix(paths, "."+v) {
			return true
		}
	}
	return false
}

// add new extension for template.
func AddTemplateExt(ext string) {
	for _, v := range BeeTemplateExt {
		if v == ext {
			return
		}
	}
	BeeTemplateExt = append(BeeTemplateExt, ext)
}

// build all template files in a directory.
// it makes beego can render any template file in view directory.
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

func getTplDeep(root, file, parent string, t *template.Template) (*template.Template, [][]string, error) {
	var fileabspath string
	if filepath.HasPrefix(file, "../") {
		fileabspath = filepath.Join(root, filepath.Dir(parent), file)
	} else {
		fileabspath = filepath.Join(root, file)
	}
	if e := utils.FileExists(fileabspath); !e {
		panic("can't find template file:" + file)
	}
	data, err := ioutil.ReadFile(fileabspath)
	if err != nil {
		return nil, [][]string{}, err
	}
	t, err = t.New(file).Parse(string(data))
	if err != nil {
		return nil, [][]string{}, err
	}
	reg := regexp.MustCompile(TemplateLeft + "[ ]*template[ ]+\"([^\"]+)\"")
	allsub := reg.FindAllStringSubmatch(string(data), -1)
	for _, m := range allsub {
		if len(m) == 2 {
			tlook := t.Lookup(m[1])
			if tlook != nil {
				continue
			}
			if !HasTemplateExt(m[1]) {
				continue
			}
			t, _, err = getTplDeep(root, m[1], file, t)
			if err != nil {
				return nil, [][]string{}, err
			}
		}
	}
	return t, allsub, nil
}

func getTemplate(root, file string, others ...string) (t *template.Template, err error) {
	t = template.New(file).Delims(TemplateLeft, TemplateRight).Funcs(beegoTplFuncMap)
	var submods [][]string
	t, submods, err = getTplDeep(root, file, "", t)
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
					t, submods1, err = getTplDeep(root, otherfile, "", t)
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
				reg := regexp.MustCompile(TemplateLeft + "[ ]*define[ ]+\"([^\"]+)\"")
				allsub := reg.FindAllStringSubmatch(string(data), -1)
				for _, sub := range allsub {
					if len(sub) == 2 && sub[1] == m[1] {
						var submods1 [][]string
						t, submods1, err = getTplDeep(root, otherfile, "", t)
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
