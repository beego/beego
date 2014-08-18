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
	"os"
	"testing"
)

var jsoncontext = `{
"appname": "beeapi",
"httpport": 8080,
"mysqlport": 3600,
"PI": 3.1415976,
"runmode": "dev",
"autorender": false,
"copyrequestbody": true,
"database": {
        "host": "host",                 
        "port": "port",                 
        "database": "database",
        "username": "username",
        "password": "password",
		"conns":{
			"maxconnection":12,
			"autoconnect":true,
			"connectioninfo":"info"
		}
    }
}`

var jsoncontextwitharray = `[
	{
		"url": "user",
		"serviceAPI": "http://www.test.com/user"
	},
	{
		"url": "employee",
		"serviceAPI": "http://www.test.com/employee"
	}
]`

func TestJsonStartsWithArray(t *testing.T) {
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
	jsonconf, err := NewConfig("json", "testjsonWithArray.conf")
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
	jsonconf, err := NewConfig("json", "testjson.conf")
	if err != nil {
		t.Fatal(err)
	}
	if jsonconf.String("appname") != "beeapi" {
		t.Fatal("appname not equal to beeapi")
	}
	if port, err := jsonconf.Int("httpport"); err != nil || port != 8080 {
		t.Error(port)
		t.Fatal(err)
	}
	if port, err := jsonconf.Int64("mysqlport"); err != nil || port != 3600 {
		t.Error(port)
		t.Fatal(err)
	}
	if pi, err := jsonconf.Float("PI"); err != nil || pi != 3.1415976 {
		t.Error(pi)
		t.Fatal(err)
	}
	if jsonconf.String("runmode") != "dev" {
		t.Fatal("runmode not equal to dev")
	}
	if v, err := jsonconf.Bool("autorender"); err != nil || v != false {
		t.Error(v)
		t.Fatal(err)
	}
	if v, err := jsonconf.Bool("copyrequestbody"); err != nil || v != true {
		t.Error(v)
		t.Fatal(err)
	}
	if err = jsonconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	if jsonconf.String("name") != "astaxie" {
		t.Fatal("get name error")
	}
	if jsonconf.String("database::host") != "host" {
		t.Fatal("get database::host error")
	}
	if jsonconf.String("database::conns::connectioninfo") != "info" {
		t.Fatal("get database::conns::connectioninfo error")
	}
	if maxconnection, err := jsonconf.Int("database::conns::maxconnection"); err != nil || maxconnection != 12 {
		t.Fatal("get database::conns::maxconnection error")
	}
	if db, err := jsonconf.DIY("database"); err != nil {
		t.Fatal(err)
	} else if m, ok := db.(map[string]interface{}); !ok {
		t.Log(db)
		t.Fatal("db not map[string]interface{}")
	} else {
		if m["host"].(string) != "host" {
			t.Fatal("get host err")
		}
	}

	if _, err := jsonconf.Int("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int")
	}

	if _, err := jsonconf.Int64("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an Int64")
	}

	if _, err := jsonconf.Float("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Float")
	}

	if _, err := jsonconf.DIY("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting an interface{}")
	}

	if val := jsonconf.String("unknown"); val != "" {
		t.Error("unknown keys should return an empty string when expecting a String")
	}

	if _, err := jsonconf.Bool("unknown"); err == nil {
		t.Error("unknown keys should return an error when expecting a Bool")
	}
}
