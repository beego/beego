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
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/astaxie/beego/utils"
)

var globalRouterTemplate = `package routers

import (
	"github.com/astaxie/beego"
)

func init() {
{{.globalinfo}}
}
`

var (
	lastupdateFilename string = "lastupdate.tmp"
	commentFilename    string
	pkgLastupdate      map[string]int64
	genInfoList        map[string][]ControllerComments
)

const COMMENTFL = "commentsRouter_"

func init() {
	pkgLastupdate = make(map[string]int64)
}

func parserPkg(pkgRealpath, pkgpath string) error {
	rep := strings.NewReplacer("/", "_", ".", "_")
	commentFilename = COMMENTFL + rep.Replace(pkgpath) + ".go"
	if !compareFile(pkgRealpath) {
		Info(pkgRealpath + " has not changed, not reloading")
		return nil
	}
	genInfoList = make(map[string][]ControllerComments)
	fileSet := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fileSet, pkgRealpath, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)

	if err != nil {
		return err
	}
	for _, pkg := range astPkgs {
		for _, fl := range pkg.Files {
			for _, d := range fl.Decls {
				switch specDecl := d.(type) {
				case *ast.FuncDecl:
					if specDecl.Recv != nil {
						parserComments(specDecl.Doc, specDecl.Name.String(), fmt.Sprint(specDecl.Recv.List[0].Type.(*ast.StarExpr).X), pkgpath)
					}
				}
			}
		}
	}
	genRouterCode()
	savetoFile(pkgRealpath)
	return nil
}

func parserComments(comments *ast.CommentGroup, funcName, controllerName, pkgpath string) error {
	if comments != nil && comments.List != nil {
		for _, c := range comments.List {
			t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
			if strings.HasPrefix(t, "@router") {
				elements := strings.TrimLeft(t, "@router ")
				e1 := strings.SplitN(elements, " ", 2)
				if len(e1) < 1 {
					return errors.New("you should has router infomation")
				}
				key := pkgpath + ":" + controllerName
				cc := ControllerComments{}
				cc.Method = funcName
				cc.Router = e1[0]
				if len(e1) == 2 && e1[1] != "" {
					e1 = strings.SplitN(e1[1], " ", 2)
					if len(e1) >= 1 {
						cc.AllowHTTPMethods = strings.Split(strings.Trim(e1[0], "[]"), ",")
					} else {
						cc.AllowHTTPMethods = append(cc.AllowHTTPMethods, "get")
					}
				} else {
					cc.AllowHTTPMethods = append(cc.AllowHTTPMethods, "get")
				}
				if len(e1) == 2 && e1[1] != "" {
					keyval := strings.Split(strings.Trim(e1[1], "[]"), " ")
					for _, kv := range keyval {
						kk := strings.Split(kv, ":")
						cc.Params = append(cc.Params, map[string]string{strings.Join(kk[:len(kk)-1], ":"): kk[len(kk)-1]})
					}
				}
				genInfoList[key] = append(genInfoList[key], cc)
			}
		}
	}
	return nil
}

func genRouterCode() {
	os.Mkdir(path.Join(workPath, "routers"), 0755)
	Info("generate router from comments")
	var globalinfo string
	sortKey := make([]string, 0)
	for k, _ := range genInfoList {
		sortKey = append(sortKey, k)
	}
	sort.Strings(sortKey)
	for _, k := range sortKey {
		cList := genInfoList[k]
		for _, c := range cList {
			allmethod := "nil"
			if len(c.AllowHTTPMethods) > 0 {
				allmethod = "[]string{"
				for _, m := range c.AllowHTTPMethods {
					allmethod += `"` + m + `",`
				}
				allmethod = strings.TrimRight(allmethod, ",") + "}"
			}
			params := "nil"
			if len(c.Params) > 0 {
				params = "[]map[string]string{"
				for _, p := range c.Params {
					for k, v := range p {
						params = params + `map[string]string{` + k + `:"` + v + `"},`
					}
				}
				params = strings.TrimRight(params, ",") + "}"
			}
			globalinfo = globalinfo + `
	beego.GlobalControllerRouter["` + k + `"] = append(beego.GlobalControllerRouter["` + k + `"],
		beego.ControllerComments{
			"` + strings.TrimSpace(c.Method) + `",
			` + "`" + c.Router + "`" + `,
			` + allmethod + `,
			` + params + `})
`
		}
	}
	if globalinfo != "" {
		f, err := os.Create(path.Join(workPath, "routers", commentFilename))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(strings.Replace(globalRouterTemplate, "{{.globalinfo}}", globalinfo, -1))
	}
}

func compareFile(pkgRealpath string) bool {
	if !utils.FileExists(path.Join(workPath, "routers", commentFilename)) {
		return true
	}
	if utils.FileExists(path.Join(workPath, lastupdateFilename)) {
		content, err := ioutil.ReadFile(path.Join(workPath, lastupdateFilename))
		if err != nil {
			return true
		}
		json.Unmarshal(content, &pkgLastupdate)
		lastupdate, err := getpathTime(pkgRealpath)
		if err != nil {
			return true
		}
		if v, ok := pkgLastupdate[pkgRealpath]; ok {
			if lastupdate <= v {
				return false
			}
		}
	}
	return true
}

func savetoFile(pkgRealpath string) {
	lastupdate, err := getpathTime(pkgRealpath)
	if err != nil {
		return
	}
	pkgLastupdate[pkgRealpath] = lastupdate
	d, err := json.Marshal(pkgLastupdate)
	if err != nil {
		return
	}
	ioutil.WriteFile(path.Join(workPath, lastupdateFilename), d, os.ModePerm)
}

func getpathTime(pkgRealpath string) (lastupdate int64, err error) {
	fl, err := ioutil.ReadDir(pkgRealpath)
	if err != nil {
		return lastupdate, err
	}
	for _, f := range fl {
		if lastupdate < f.ModTime().UnixNano() {
			lastupdate = f.ModTime().UnixNano()
		}
	}
	return lastupdate, nil
}
