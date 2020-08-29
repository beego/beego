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

package json

import (
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/pkg/infrastructure/config"
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
	f, err := os.Create("testjsonWithArray.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(jsoncontextwitharray)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testjsonWithArray.conf")
	jsonconf, err := config.NewConfig("json", "testjsonWithArray.conf")
	if err != nil {
		t.Fatal(err)
	}
	rootArray, err := jsonconf.DIY(nil, "rootArray")
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

	f, err := os.Create("testjson.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(jsoncontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testjson.conf")
	jsonconf, err := config.NewConfig("json", "testjson.conf")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range keyValue {
		var err error
		var value interface{}
		switch v.(type) {
		case int:
			value, err = jsonconf.Int(nil, k)
		case int64:
			value, err = jsonconf.Int64(nil, k)
		case float64:
			value, err = jsonconf.Float(nil, k)
		case bool:
			value, err = jsonconf.Bool(nil, k)
		case []string:
			value, err = jsonconf.Strings(nil, k)
		case string:
			value, err = jsonconf.String(nil, k)
		default:
			value, err = jsonconf.DIY(nil, k)
		}
		if err != nil {
			t.Fatalf("get key %q value fatal,%v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Fatalf("get key %q value, want %v got %v .", k, v, value)
		}

	}
	if err = jsonconf.Set(nil, "name", "astaxie"); err != nil {
		t.Fatal(err)
	}

	res, _ := jsonconf.String(nil, "name")
	if res != "astaxie" {
		t.Fatal("get name error")
	}

	if db, err := jsonconf.DIY(nil, "database"); err != nil {
		t.Fatal(err)
	} else if m, ok := db.(map[string]interface{}); !ok {
		t.Log(db)
		t.Fatal("db not map[string]interface{}")
	} else {
		if m["host"].(string) != "host" {
			t.Fatal("get host err")
		}
	}

	if _, err := jsonconf.Int(nil, "unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int")
	}

	if _, err := jsonconf.Int64(nil, "unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int64")
	}

	if _, err := jsonconf.Float(nil, "unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Float")
	}

	if _, err := jsonconf.DIY(nil, "unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an interface{}")
	}

	if val, _ := jsonconf.String(nil, "unknown"); val != "" {
		t.Error("unknown keys should return an empty string when expecting a String")
	}

	if _, err := jsonconf.Bool(nil, "unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Bool")
	}

	if !jsonconf.DefaultBool(nil, "unknown", true) {
		t.Error("unknown keys with default value wrong")
	}
}
