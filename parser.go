// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie
package beego

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

var globalRouterTemplate = `package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	{{.globalinfo}}
}
`

var genInfoList map[string][]ControllerComments

func init() {
	genInfoList = make(map[string][]ControllerComments)
}

func parserPkg(pkgRealpath, pkgpath string) error {
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
					parserComments(specDecl.Doc, specDecl.Name.String(), fmt.Sprint(specDecl.Recv.List[0].Type.(*ast.StarExpr).X), pkgpath)
				}
			}
		}
	}
	genRouterCode()
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
	os.Mkdir(path.Join(AppPath, "routers"), 0755)
	Info("generate router from comments")
	var globalinfo string
	for k, cList := range genInfoList {
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
			globalinfo = globalinfo + fmt.Sprintln(`beego.GlobalControllerRouter["`+k+`"] = append(beego.GlobalControllerRouter["`+k+`"], beego.ControllerComments{"`+
				strings.TrimSpace(c.Method)+`", "`+c.Router+`", `+allmethod+", "+params+"})")
		}
	}
	if globalinfo != "" {
		f, err := os.Create(path.Join(AppPath, "routers", "commentsRouter.go"))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(strings.Replace(globalRouterTemplate, "{{.globalinfo}}", globalinfo, -1))
	}
}
