package config

import (
	"os"
	"testing"
)

var inicontext = `
;comment one
#comment two
appname = beeapi
httpport = 8080
mysqlport = 3600
PI = 3.1415976
runmode = "dev"
autorender = false
copyrequestbody = true
[demo]
key1="asta"
key2 = "xie"
CaseInsensitive = true
`

func TestIni(t *testing.T) {
	f, err := os.Create("testini.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(inicontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testini.conf")
	iniconf, err := NewConfig("ini", "testini.conf")
	if err != nil {
		t.Fatal(err)
	}
	if iniconf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}
	if port, err := iniconf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}
	if port, err := iniconf.Int64("mysqlport"); err != nil || port != 3600 {
		t.Error(port)
		t.Fatal(err)
	}
	if pi, err := iniconf.Float("PI"); err != nil || pi != 3.1415976 {
		t.Error(pi)
		t.Fatal(err)
	}
	if iniconf.String("runmode") != "dev" {
		t.Fatal("runmode not equal to dev")
	}
	if v, err := iniconf.Bool("autorender"); err != nil || v != false {
		t.Error(v)
		t.Fatal(err)
	}
	if v, err := iniconf.Bool("copyrequestbody"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}
	if err = iniconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	if iniconf.String("name") != "astaxie" {
		t.Fatal("get name error")
	}
	if iniconf.String("demo::key1") != "asta" {
		t.Fatal("get demo.key1 error")
	}
	if iniconf.String("demo::key2") != "xie" {
		t.Fatal("get demo.key2 error")
	}
	if v, err := iniconf.Bool("demo::caseinsensitive"); err != nil || v != true {
		t.Fatal("get demo.caseinsensitive error")
	}
}
