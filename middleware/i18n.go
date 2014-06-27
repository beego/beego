// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

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
	i18n := make(map[string]map[string]string)
	file, err := os.Open(filepath)
	if err != nil {
		panic("open " + filepath + " err :" + err.Error())
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("read " + filepath + " err :" + err.Error())
	}
	err = json.Unmarshal(data, &i18n)
	if err != nil {
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
