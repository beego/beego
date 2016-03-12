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

package xml

import (
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/config"
)

func TestXML(t *testing.T) {

	var (
		//xml parse should incluce in <config></config> tags
		xmlcontext = `<?xml version="1.0" encoding="UTF-8"?>
<config>
<appname>beeapi</appname>
<httpport>8080</httpport>
<mysqlport>3600</mysqlport>
<PI>3.1415976</PI>
<runmode>dev</runmode>
<autorender>false</autorender>
<copyrequestbody>true</copyrequestbody>
<path1>$$GOROOT</path1>
<path2>$$GOROOT||/home/go</path2>
<path3>$$GOROOT$$GOPATH2||/home/go</path3>
<token1>$$TOKEN</token1>
<token2>$$TOKEN||</token2>
<token3>$$TOKEN||astaxie</token3>
<token4>token$$TOKEN</token4>
<token5>$$TOKEN$$TOKEN||TOKEN</token5>
</config>
`
		keyValue = map[string]interface{}{
			"appname":         "beeapi",
			"httpport":        8080,
			"mysqlport":       int64(3600),
			"PI":              3.1415976,
			"runmode":         "dev",
			"autorender":      false,
			"copyrequestbody": true,
			"path1":           os.Getenv("GOROOT"),
			"path2":           os.Getenv("GOROOT"),
			"path3":           "/home/go",
			"token1":          "",
			"token2":          "",
			"token3":          "astaxie",
			"token4":          "token$$TOKEN",
			"token5":          "TOKEN",
			"error":           "",
			"emptystrings":    []string{},
		}
	)

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
	xmlconf, err := config.NewConfig("xml", "testxml.conf")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range keyValue {

		var (
			value interface{}
			err   error
		)

		switch v.(type) {
		case int:
			value, err = xmlconf.Int(k)
		case int64:
			value, err = xmlconf.Int64(k)
		case float64:
			value, err = xmlconf.Float(k)
		case bool:
			value, err = xmlconf.Bool(k)
		case []string:
			value = xmlconf.Strings(k)
		case string:
			value = xmlconf.String(k)
		default:
			value, err = xmlconf.DIY(k)
		}
		if err != nil {
			t.Errorf("get key %q value fatal,%v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Errorf("get key %q value, want %v got %v .", k, v, value)
		}

	}

	if err = xmlconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	if xmlconf.String("name") != "astaxie" {
		t.Fatal("get name error")
	}
}
