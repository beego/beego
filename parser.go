// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie
package beego

import (
	"os"
	"path/filepath"
)

var globalControllerRouter = `package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	{{.globalinfo}}
}
`

func parserPkg(pkgpath string) error {
	err := filepath.Walk(pkgpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Error("error scan app Controller source:", err)
			return err
		}
		//if is normal file or name is temp skip
		//directory is needed
		if !info.IsDir() || info.Name() == "tmp" {
			return nil
		}

		//fileSet := token.NewFileSet()
		//astPkgs, err := parser.ParseDir(fileSet, path, func(info os.FileInfo) bool {
		//	name := info.Name()
		//	return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
		//}, parser.ParseComments)

		return nil
	})

	return err
}
