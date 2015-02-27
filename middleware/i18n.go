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

// Usage:
//
// import "github.com/astaxie/beego/middleware"
//
// I18N = middleware.NewLocale("conf/i18n.conf", beego.AppConfig.String("language"))
//
// more docs: http://beego.me/docs/module/i18n.md
package middleware

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Translation struct {
	filepath     string
	CurrentLocal string
	Locales      map[string]map[string]string
}

func NewLocale(filepath string, defaultlocal string) *Translation {
	file, err := os.Open(filepath)
	if err != nil {
		panic("open " + filepath + " err :" + err.Error())
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("read " + filepath + " err :" + err.Error())
	}

	i18n := make(map[string]map[string]string)
	if err = json.Unmarshal(data, &i18n); err != nil {
		panic("json.Unmarshal " + filepath + " err :" + err.Error())
	}
	return &Translation{
		filepath:     filepath,
		CurrentLocal: defaultlocal,
		Locales:      i18n,
	}
}

func (t *Translation) SetLocale(local string) {
	t.CurrentLocal = local
}

func (t *Translation) Translate(key string, local string) string {
	if local == "" {
		local = t.CurrentLocal
	}
	if ct, ok := t.Locales[key]; ok {
		if v, o := ct[local]; o {
			return v
		}
	}
	return key
}
