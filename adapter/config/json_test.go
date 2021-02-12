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

package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonStartsWithArray(t *testing.T) {

	const jsoncontextwitharray = `[
	{
		"url": "user",
		"serviceAPI": "http://www.test.com/user"
	},
	{
		"url": "employee",
		"serviceAPI": "http://www.test.com/employee"
	}
]`
	cfgFileName := "testjsonWithArray.conf"
	f, err := os.Create(cfgFileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(jsoncontextwitharray)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(cfgFileName)
	jsonconf, err := NewConfig("json", cfgFileName)
	if err != nil {
		t.Fatal(err)
	}
	rootArray, err := jsonconf.DIY("rootArray")
	if err != nil {
		t.Error("array does not exist as element")
	}
	rootArrayCasted := rootArray.([]interface{})
	if rootArrayCasted == nil {
		t.Error("array from root is nil")
	} else {
		elem := rootArrayCasted[0].(map[string]interface{})
		if elem["url"] != "user" || elem["serviceAPI"] != "http://www.test.com/user" {
			t.Error("array[0] values are not valid")
		}

		elem2 := rootArrayCasted[1].(map[string]interface{})
		if elem2["url"] != "employee" || elem2["serviceAPI"] != "http://www.test.com/employee" {
			t.Error("array[1] values are not valid")
		}
	}
}

func TestJson(t *testing.T) {

	var (
		jsoncontext = `{
"appname": "beeapi",
"testnames": "foo;bar",
"httpport": 8080,
"mysqlport": 3600,
"PI": 3.1415976, 
"runmode": "dev",
"autorender": false,
"copyrequestbody": true,
"session": "on",
"cookieon": "off",
"newreg": "OFF",
"needlogin": "ON",
"enableSession": "Y",
"enableCookie": "N",
"flag": 1,
"path1": "${GOPATH}",
"path2": "${GOPATH||/home/go}",
"database": {
        "host": "host",
        "port": "port",
        "database": "database",
        "username": "username",
        "password": "${GOPATH}",
		"conns":{
			"maxconnection":12,
			"autoconnect":true,
			"connectioninfo":"info",
			"root": "${GOPATH}"
		}
    }
}`
		keyValue = map[string]interface{}{
			"appname":                         "beeapi",
			"testnames":                       []string{"foo", "bar"},
			"httpport":                        8080,
			"mysqlport":                       int64(3600),
			"PI":                              3.1415976,
			"runmode":                         "dev",
			"autorender":                      false,
			"copyrequestbody":                 true,
			"session":                         true,
			"cookieon":                        false,
			"newreg":                          false,
			"needlogin":                       true,
			"enableSession":                   true,
			"enableCookie":                    false,
			"flag":                            true,
			"path1":                           os.Getenv("GOPATH"),
			"path2":                           os.Getenv("GOPATH"),
			"database::host":                  "host",
			"database::port":                  "port",
			"database::database":              "database",
			"database::password":              os.Getenv("GOPATH"),
			"database::conns::maxconnection":  12,
			"database::conns::autoconnect":    true,
			"database::conns::connectioninfo": "info",
			"database::conns::root":           os.Getenv("GOPATH"),
			"unknown":                         "",
		}
	)

	cfgFileName := "testjson.conf"
	f, err := os.Create(cfgFileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(jsoncontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(cfgFileName)
	jsonconf, err := NewConfig("json", cfgFileName)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range keyValue {
		var err error
		var value interface{}
		switch v.(type) {
		case int:
			value, err = jsonconf.Int(k)
		case int64:
			value, err = jsonconf.Int64(k)
		case float64:
			value, err = jsonconf.Float(k)
		case bool:
			value, err = jsonconf.Bool(k)
		case []string:
			value = jsonconf.Strings(k)
		case string:
			value = jsonconf.String(k)
		default:
			value, err = jsonconf.DIY(k)
		}

		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("%v", v), fmt.Sprintf("%v", value))
	}

	assert.Nil(t, jsonconf.Set("name", "astaxie"))

	assert.Equal(t, "astaxie", jsonconf.String("name"))

	db, err := jsonconf.DIY("database")
	assert.Nil(t, err)

	m, ok := db.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t,"host" , m["host"])

	_, err = jsonconf.Int("unknown")
	assert.NotNil(t, err)

	_, err = jsonconf.Int64("unknown")
	assert.NotNil(t, err)

	_, err = jsonconf.Float("unknown")
	assert.NotNil(t, err)

	_, err = jsonconf.DIY("unknown")
	assert.NotNil(t, err)

	val := jsonconf.String("unknown")
	assert.Equal(t, "", val)

	_, err = jsonconf.Bool("unknown")
	assert.NotNil(t, err)

	assert.True(t, jsonconf.DefaultBool("unknown", true))
}
