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
        "password": "password"
    }
}`

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
	if db, err := jsonconf.DIY("database"); err != nil {
		t.Fatal(err)
	} else if m, ok := db.(map[string]interface{}); !ok {
		t.Fatal("db not map[string]interface{}")
	} else {
		if m["host"].(string) != "host" {
			t.Fatal("get host err")
		}
	}
}
