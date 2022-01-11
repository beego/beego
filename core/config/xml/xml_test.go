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

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/core/config"
)

func TestXML(t *testing.T) {
	var (
		// xml parse should include in <config></config> tags
		xmlcontext = `<?xml version="1.0" encoding="UTF-8"?>
<config>
<appname>beeapi</appname>
<httpport>8080</httpport>
<mysqlport>3600</mysqlport>
<PI>3.1415976</PI>
<runmode>dev</runmode>
<autorender>false</autorender>
<copyrequestbody>true</copyrequestbody>
<path1>${GOPATH}</path1>
<path2>${GOPATH||/home/go}</path2>
<mysection>
<id>1</id>
<name>MySection</name>
</mysection>
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
			"path1":           os.Getenv("GOPATH"),
			"path2":           os.Getenv("GOPATH"),
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

	var xmlsection map[string]string
	xmlsection, err = xmlconf.GetSection("mysection")
	if err != nil {
		t.Fatal(err)
	}

	if len(xmlsection) == 0 {
		t.Error("section should not be empty")
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
			value, err = xmlconf.Strings(k)
		case string:
			value, err = xmlconf.String(k)
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

	res, _ := xmlconf.String("name")
	if res != "astaxie" {
		t.Fatal("get name error")
	}

	sub, err := xmlconf.Sub("mysection")
	assert.Nil(t, err)
	assert.NotNil(t, sub)
	name, err := sub.String("name")
	assert.Nil(t, err)
	assert.Equal(t, "MySection", name)

	id, err := sub.Int("id")
	assert.Nil(t, err)
	assert.Equal(t, 1, id)

	sec := &Section{}

	err = sub.Unmarshaler("", sec)
	assert.Nil(t, err)
	assert.Equal(t, "MySection", sec.Name)

	sec = &Section{}

	err = xmlconf.Unmarshaler("mysection", sec)
	assert.Nil(t, err)
	assert.Equal(t, "MySection", sec.Name)
}

func TestXMLMissConfig(t *testing.T) {
	xmlcontext1 := `
	<?xml version="1.0" encoding="UTF-8"?>
	<appname>beeapi</appname>
	`

	c := &Config{}
	_, err := c.ParseData([]byte(xmlcontext1))
	assert.Equal(t, "xml parse should include in <config></config> tags", err.Error())

	xmlcontext2 := `
	<?xml version="1.0" encoding="UTF-8"?>
	<config></config>
	`

	_, err = c.ParseData([]byte(xmlcontext2))
	assert.Equal(t, "xml parse <config></config> tags should include sub tags", err.Error())
}

type Section struct {
	Name string `xml:"name"`
}
