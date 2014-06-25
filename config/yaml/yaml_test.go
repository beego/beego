// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package config

import (
	"os"
	"testing"

	"github.com/astaxie/beego/config"
)

var yamlcontext = `
"appname": beeapi
"httpport": 8080
"mysqlport": 3600
"PI": 3.1415976
"runmode": dev
"autorender": false
"copyrequestbody": true
`

func TestYaml(t *testing.T) {
	f, err := os.Create("testyaml.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(yamlcontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testyaml.conf")
	yamlconf, err := config.NewConfig("yaml", "testyaml.conf")
	if err != nil {
		t.Fatal(err)
	}
	if yamlconf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}
	if port, err := yamlconf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}
	if port, err := yamlconf.Int64("mysqlport"); err != nil || port != 3600 {
		t.Error(port)
		t.Fatal(err)
	}
	if pi, err := yamlconf.Float("PI"); err != nil || pi != 3.1415976 {
		t.Error(pi)
		t.Fatal(err)
	}
	if yamlconf.String("runmode") != "dev" {
		t.Fatal("runmode not equal to dev")
	}
	if v, err := yamlconf.Bool("autorender"); err != nil || v != false {
		t.Error(v)
		t.Fatal(err)
	}
	if v, err := yamlconf.Bool("copyrequestbody"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}
	if err = yamlconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	if yamlconf.String("name") != "astaxie" {
		t.Fatal("get name error")
	}
}
