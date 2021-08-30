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

package yaml

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/core/config"
)

func TestYaml(t *testing.T) {
	var (
		yamlcontext = `
"appname": beeapi
"httpport": 8080
"mysqlport": 3600
"PI": 3.1415976
"runmode": dev
"autorender": false
"copyrequestbody": true
"PATH": GOPATH
"path1": ${GOPATH}
"path2": ${GOPATH||/home/go}
"empty": "" 
"user":
  "name": "tom"
  "age": 13
`

		keyValue = map[string]interface{}{
			"appname":         "beeapi",
			"httpport":        8080,
			"mysqlport":       int64(3600),
			"PI":              3.1415976,
			"runmode":         "dev",
			"autorender":      false,
			"copyrequestbody": true,
			"PATH":            "GOPATH",
			"path1":           os.Getenv("GOPATH"),
			"path2":           os.Getenv("GOPATH"),
			"error":           "",
			"emptystrings":    []string{},
		}
	)

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

	m, err := ReadYmlReader("testyaml.conf")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, m, yamlconf.(*ConfigContainer).data)

	shadow, err := (&Config{}).ParseData([]byte(yamlcontext))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, shadow, yamlconf)

	yamlconf.OnChange("abc", func(value string) {
		fmt.Printf("on change, value is %s \n", value)
	})

	res, _ := yamlconf.String("appname")
	if res != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}

	for k, v := range keyValue {

		var (
			value interface{}
			err   error
		)

		switch v.(type) {
		case int:
			value, err = yamlconf.Int(k)
		case int64:
			value, err = yamlconf.Int64(k)
		case float64:
			value, err = yamlconf.Float(k)
		case bool:
			value, err = yamlconf.Bool(k)
		case []string:
			value, err = yamlconf.Strings(k)
		case string:
			value, err = yamlconf.String(k)
		default:
			value, err = yamlconf.DIY(k)
		}
		if err != nil {
			t.Errorf("get key %q value fatal, %v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Errorf("get key %q value, want %v got %v .", k, v, value)
		}

	}

	if err = yamlconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	res, _ = yamlconf.String("name")
	if res != "astaxie" {
		t.Fatal("get name error")
	}

	sub, err := yamlconf.Sub("user")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, sub)
	name, err := sub.String("name")
	assert.Nil(t, err)
	assert.Equal(t, "tom", name)

	age, err := sub.Int("age")
	assert.Nil(t, err)
	assert.Equal(t, 13, age)

	user := &User{}

	err = sub.Unmarshaler("", user)
	assert.Nil(t, err)
	assert.Equal(t, "tom", user.Name)
	assert.Equal(t, 13, user.Age)

	user = &User{}

	err = yamlconf.Unmarshaler("user", user)
	assert.Nil(t, err)
	assert.Equal(t, "tom", user.Name)
	assert.Equal(t, 13, user.Age)

	// default value
	assert.Equal(t, "beeapi", yamlconf.DefaultString("appname", "invalid"))
	assert.Equal(t, "invalid", yamlconf.DefaultString("i-appname", "invalid"))
	assert.Equal(t, 8080, yamlconf.DefaultInt("httpport", 8090))
	assert.Equal(t, 8090, yamlconf.DefaultInt("i-httpport", 8090))
	assert.Equal(t, 3.1415976, yamlconf.DefaultFloat("PI", 3.14))
	assert.Equal(t, 3.14, yamlconf.DefaultFloat("1-PI", 3.14))
	assert.True(t, yamlconf.DefaultBool("copyrequestbody", false))
	assert.True(t, yamlconf.DefaultBool("i-copyrequestbody", true))
	assert.Equal(t, int64(8080), yamlconf.DefaultInt64("httpport", 8090))
	assert.Equal(t, int64(8090), yamlconf.DefaultInt64("i-httpport", 8090))
	assert.Equal(t, "tom", yamlconf.DefaultString("user.name", "invalid"))
	assert.Equal(t, "invalid", yamlconf.DefaultString("user.1-name", "invalid"))
	assert.Equal(t, []string{"tom"}, yamlconf.DefaultStrings("strings", []string{"tom"}))

	appName, err := yamlconf.DIY("appname")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "beeapi", appName)

	err = yamlconf.SaveConfigFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	section, err := yamlconf.GetSection("user")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "tom", section["name"])
}

type User struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}
