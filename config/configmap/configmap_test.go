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

package configmap

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/config"
)

// Is needed add to kubernetes ex: kubectl get configmap beego --namespace default {"appname": "beeapi", ...}
var jsoncontext = `{
		"appname": "beeapi",
		"httpport": "8080",
		"mysqlport": "3600",
		"PI": "3.1415976",
		"runmode": "dev",
		"autorender": "false",
		"copyrequestbody": "true",
		"path1": "${GOPATH}",
		"path2": "${GOPATH||/home/go}",
		"mysection": {
			"id": "1",
			"name": "MySection"
		}
}`

func TestConfigMap(t *testing.T) {
	var (
		err      error
		jsonconf config.Configer
	)
	if tr := os.Getenv("TRAVIS"); len(tr) > 0 {
		j, _ := json.Marshal(jsoncontext)
		jsonconf, err = config.NewConfigData("configMap", j)
	} else {
		jsonconf, err = config.NewConfig("configMap", "beego")
	}

	var keyValue = map[string]interface{}{
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

	if err != nil {
		t.Fatal(err)
	}

	var jsonsection map[string]string
	jsonsection, err = jsonconf.GetSection("mysection")
	if err != nil {
		t.Fatal(err)
	}

	if len(jsonsection) == 0 {
		t.Error("section should not be empty")
	}

	for k, v := range keyValue {

		var (
			value interface{}
			err   error
		)

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
		if err != nil {
			t.Errorf("get key %q value fatal,%v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Errorf("get key %q value, want %v got %v .", k, v, value)
		}
	}

	if err = jsonconf.Set("name", "astaxie"); err != nil {
		t.Fatal(err)
	}
	if jsonconf.String("name") != "astaxie" {
		t.Fatal("get name error")
	}
}
