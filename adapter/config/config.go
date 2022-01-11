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

// Package config is used to parse config.
// Usage:
//  import "github.com/beego/beego/v2/core/config"
// Examples.
//
//  cnf, err := config.NewConfig("ini", "config.conf")
//
//  cnf APIS:
//
//  cnf.Set(key, val string) error
//  cnf.String(key string) string
//  cnf.Strings(key string) []string
//  cnf.Int(key string) (int, error)
//  cnf.Int64(key string) (int64, error)
//  cnf.Bool(key string) (bool, error)
//  cnf.Float(key string) (float64, error)
//  cnf.DefaultString(key string, defaultVal string) string
//  cnf.DefaultStrings(key string, defaultVal []string) []string
//  cnf.DefaultInt(key string, defaultVal int) int
//  cnf.DefaultInt64(key string, defaultVal int64) int64
//  cnf.DefaultBool(key string, defaultVal bool) bool
//  cnf.DefaultFloat(key string, defaultVal float64) float64
//  cnf.DIY(key string) (interface{}, error)
//  cnf.GetSection(section string) (map[string]string, error)
//  cnf.SaveConfigFile(filename string) error
// More docs http://beego.vip/docs/module/config.md
package config

import (
	"github.com/beego/beego/v2/core/config"
)

// Configer defines how to get and set value from configuration raw data.
type Configer interface {
	Set(key, val string) error   // support section::key type in given key when using ini type.
	String(key string) string    // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	Strings(key string) []string // get string slice
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultVal []string) []string // get string slice
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
	DIY(key string) (interface{}, error)
	GetSection(section string) (map[string]string, error)
	SaveConfigFile(filename string) error
}

// Config is the adapter interface for parsing config file to get raw data to Configer.
type Config interface {
	Parse(key string) (Configer, error)
	ParseData(data []byte) (Configer, error)
}

var adapters = make(map[string]Config)

// Register makes a config adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Config) {
	config.Register(name, &oldToNewConfigAdapter{delegate: adapter})
}

// NewConfig adapterName is ini/json/xml/yaml.
// filename is the config file path.
func NewConfig(adapterName, filename string) (Configer, error) {
	cfg, err := config.NewConfig(adapterName, filename)
	if err != nil {
		return nil, err
	}

	// it was registered by using Register method
	res, ok := cfg.(*oldToNewConfigerAdapter)
	if ok {
		return res.delegate, nil
	}

	return &newToOldConfigerAdapter{
		delegate: cfg,
	}, nil
}

// NewConfigData adapterName is ini/json/xml/yaml.
// data is the config data.
func NewConfigData(adapterName string, data []byte) (Configer, error) {
	cfg, err := config.NewConfigData(adapterName, data)
	if err != nil {
		return nil, err
	}

	// it was registered by using Register method
	res, ok := cfg.(*oldToNewConfigerAdapter)
	if ok {
		return res.delegate, nil
	}

	return &newToOldConfigerAdapter{
		delegate: cfg,
	}, nil
}

// ExpandValueEnvForMap convert all string value with environment variable.
func ExpandValueEnvForMap(m map[string]interface{}) map[string]interface{} {
	return config.ExpandValueEnvForMap(m)
}

// ExpandValueEnv returns value of convert with environment variable.
//
// Return environment variable if value start with "${" and end with "}".
// Return default value if environment variable is empty or not exist.
//
// It accept value formats "${env}" , "${env||}}" , "${env||defaultValue}" , "defaultvalue".
// Examples:
//	v1 := config.ExpandValueEnv("${GOPATH}")			// return the GOPATH environment variable.
//	v2 := config.ExpandValueEnv("${GOAsta||/usr/local/go}")	// return the default value "/usr/local/go/".
//	v3 := config.ExpandValueEnv("Astaxie")				// return the value "Astaxie".
func ExpandValueEnv(value string) string {
	return config.ExpandValueEnv(value)
}

// ParseBool returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func ParseBool(val interface{}) (value bool, err error) {
	return config.ParseBool(val)
}

// ToString converts values of any type to string.
func ToString(x interface{}) string {
	return config.ToString(x)
}
