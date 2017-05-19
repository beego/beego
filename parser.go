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
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/astaxie/beego/context/param"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils"
)

var globalRouterTemplate = `package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {
{{.globalinfo}}
}
`

var (
	lastupdateFilename = "lastupdate.tmp"
	commentFilename    string
	pkgLastupdate      map[string]int64
	genInfoList        map[string][]ControllerComments
)

const commentPrefix = "commentsRouter_"

func init() {
	pkgLastupdate = make(map[string]int64)
}

func parserPkg(pkgRealpath, pkgpath string) error {
	rep := strings.NewReplacer("\\", "_", "/", "_", ".", "_")
	commentFilename, _ = filepath.Rel(AppPath, pkgRealpath)
	commentFilename = commentPrefix + rep.Replace(commentFilename) + ".go"
	if !compareFile(pkgRealpath) {
		logs.Info(pkgRealpath + " no changed")
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
						exp, ok := specDecl.Recv.List[0].Type.(*ast.StarExpr) // Check that the type is correct first beforing throwing to parser
						if ok {
							parserComments(specDecl, fmt.Sprint(exp.X), pkgpath)
						}
					}
				}
			}
		}
	}
	genRouterCode(pkgRealpath)
	savetoFile(pkgRealpath)
	return nil
}

type parsedComment struct {
	routerPath string
	methods    []string
	params     map[string]parsedParam
}

type parsedParam struct {
	name     string
	datatype string
	location string
	defValue string
	required bool
}

func parserComments(f *ast.FuncDecl, controllerName, pkgpath string) error {
	if f.Doc != nil {
		parsedComment, err := parseComment(f.Doc.List)
		if err != nil {
			return err
		}
		if parsedComment.routerPath != "" {
			key := pkgpath + ":" + controllerName
			cc := ControllerComments{}
			cc.Method = f.Name.String()
			cc.Router = parsedComment.routerPath
			cc.AllowHTTPMethods = parsedComment.methods
			cc.MethodParams = buildMethodParams(f.Type.Params.List, parsedComment)
			genInfoList[key] = append(genInfoList[key], cc)
		}

	}
	return nil
}

func buildMethodParams(funcParams []*ast.Field, pc *parsedComment) []*param.MethodParam {
	result := make([]*param.MethodParam, 0, len(funcParams))
	for _, fparam := range funcParams {
		methodParam := buildMethodParam(fparam, pc)
		result = append(result, methodParam)
	}
	return result
}

func buildMethodParam(fparam *ast.Field, pc *parsedComment) *param.MethodParam {
	options := []param.MethodParamOption{}
	name := fparam.Names[0].Name
	if cparam, ok := pc.params[name]; ok {
		//Build param from comment info
		name = cparam.name
		if cparam.required {
			options = append(options, param.IsRequired)
		}
		switch cparam.location {
		case "body":
			options = append(options, param.InBody)
		case "header":
			options = append(options, param.InHeader)
		case "path":
			options = append(options, param.InPath)
		}
		if cparam.defValue != "" {
			options = append(options, param.Default(cparam.defValue))
		}
	} else {
		if paramInPath(name, pc.routerPath) {
			options = append(options, param.InPath)
		}
	}
	return param.New(name, options...)
}

func paramInPath(name, route string) bool {
	return strings.HasSuffix(route, ":"+name) ||
		strings.Contains(route, ":"+name+"/")
}

var routeRegex = regexp.MustCompile(`@router\s+(\S+)(?:\s+\[(\S+)\])?`)

func parseComment(lines []*ast.Comment) (pc *parsedComment, err error) {
	pc = &parsedComment{}
	for _, c := range lines {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(t, "@router") {
			matches := routeRegex.FindStringSubmatch(t)
			if len(matches) == 3 {
				pc.routerPath = matches[1]
				methods := matches[2]
				if methods == "" {
					pc.methods = []string{"get"}
					//pc.hasGet = true
				} else {
					pc.methods = strings.Split(methods, ",")
					//pc.hasGet = strings.Contains(methods, "get")
				}
			} else {
				return nil, errors.New("Router information is missing")
			}
		} else if strings.HasPrefix(t, "@Param") {
			pv := getparams(strings.TrimSpace(strings.TrimLeft(t, "@Param")))
			if len(pv) < 4 {
				logs.Error("Invalid @Param format. Needs at least 4 parameters")
			}
			p := parsedParam{}
			names := strings.SplitN(pv[0], "=>", 2)
			p.name = names[0]
			funcParamName := p.name
			if len(names) > 1 {
				funcParamName = names[1]
			}
			p.location = pv[1]
			p.datatype = pv[2]
			switch len(pv) {
			case 5:
				p.required, _ = strconv.ParseBool(pv[3])
			case 6:
				p.defValue = pv[3]
				p.required, _ = strconv.ParseBool(pv[4])
			}
			if pc.params == nil {
				pc.params = map[string]parsedParam{}
			}
			pc.params[funcParamName] = p
		}
	}
	return
}

// direct copy from bee\g_docs.go
// analisys params return []string
// @Param	query		form	 string	true		"The email for login"
// [query form string true "The email for login"]
func getparams(str string) []string {
	var s []rune
	var j int
	var start bool
	var r []string
	var quoted int8
	for _, c := range str {
		if unicode.IsSpace(c) && quoted == 0 {
			if !start {
				continue
			} else {
				start = false
				j++
				r = append(r, string(s))
				s = make([]rune, 0)
				continue
			}
		}

		start = true
		if c == '"' {
			quoted ^= 1
			continue
		}
		s = append(s, c)
	}
	if len(s) > 0 {
		r = append(r, string(s))
	}
	return r
}

func genRouterCode(pkgRealpath string) {
	os.Mkdir(getRouterDir(pkgRealpath), 0755)
	logs.Info("generate router from comments")
	var (
		globalinfo string
		sortKey    []string
	)
	for k := range genInfoList {
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
			methodParams := "param.Make("
			if len(c.MethodParams) > 0 {
				lines := make([]string, 0, len(c.MethodParams))
				for _, m := range c.MethodParams {
					lines = append(lines, fmt.Sprint(m))
				}
				methodParams += "\n				" +
					strings.Join(lines, ",\n				") +
					",\n			"
			}
			methodParams += ")"
			globalinfo = globalinfo + `
	beego.GlobalControllerRouter["` + k + `"] = append(beego.GlobalControllerRouter["` + k + `"],
		beego.ControllerComments{
			Method: "` + strings.TrimSpace(c.Method) + `",
			` + "Router: `" + c.Router + "`" + `,
			AllowHTTPMethods: ` + allmethod + `,
			MethodParams: ` + methodParams + `,
			Params: ` + params + `})
`
		}
	}
	if globalinfo != "" {
		f, err := os.Create(filepath.Join(getRouterDir(pkgRealpath), commentFilename))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(strings.Replace(globalRouterTemplate, "{{.globalinfo}}", globalinfo, -1))
	}
}

func compareFile(pkgRealpath string) bool {
	if !utils.FileExists(filepath.Join(getRouterDir(pkgRealpath), commentFilename)) {
		return true
	}
	if utils.FileExists(lastupdateFilename) {
		content, err := ioutil.ReadFile(lastupdateFilename)
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
	ioutil.WriteFile(lastupdateFilename, d, os.ModePerm)
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

func getRouterDir(pkgRealpath string) string {
	dir := filepath.Dir(pkgRealpath)
	for {
		d := filepath.Join(dir, "routers")
		if utils.FileExists(d) {
			return d
		}

		if r, _ := filepath.Rel(dir, AppPath); r == "." {
			return d
		}
		// Parent dir.
		dir = filepath.Dir(dir)
	}
}
