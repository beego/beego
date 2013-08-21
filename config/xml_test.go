package config

import (
	"os"
	"testing"
)

//xml parse should incluce in <config></config> tags
var xmlcontext = `<?xml version="1.0" encoding="UTF-8"?>
<config>
<appname>beeapi</appname>
<httpport>8080</httpport>
<mysqlport>3600</mysqlport>
<PI>3.1415976</PI>
<runmode>dev</runmode>
<autorender>false</autorender>
<copyrequestbody>true</copyrequestbody>
</config>
`

func TestXML(t *testing.T) {
	f, err := os.Create("testxml.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(xmlcontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testxml.conf")
	xmlconf, err := NewConfig("xml", "testxml.conf")
	if err != nil {
		t.Fatal(err)
	}
	if xmlconf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}
	if port, err := xmlconf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}
	if port, err := xmlconf.Int64("mysqlport"); err != nil || port != 3600 {
		t.Error(port)
		t.Fatal(err)
	}
	if pi, err := xmlconf.Float("PI"); err != nil || pi != 3.1415976 {
		t.Error(pi)
		t.Fatal(err)
	}
	if xmlconf.String("runmode") != "dev" {
		t.Fatal("runmode not equal to dev")
	}
	if v, err := xmlconf.Bool("autorender"); err != nil || v != false {
		t.Error(v)
		t.Fatal(err)
	}
	if v, err := xmlconf.Bool("copyrequestbody"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}
	if err = xmlconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	if xmlconf.String("name") != "astaxie" {
		t.Fatal("get name error")
	}
}
