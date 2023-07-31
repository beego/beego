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
//
//	import "github.com/beego/beego/v2/core/config"
//
// Examples.
//
//	cnf, err := config.NewConfig("ini", "config.conf")
//
//	cnf APIS:
//
//	cnf.Set(key, val string) error
//	cnf.String(key string) string
//	cnf.Strings(key string) []string
//	cnf.Int(key string) (int, error)
//	cnf.Int64(key string) (int64, error)
//	cnf.Bool(key string) (bool, error)
//	cnf.Float(key string) (float64, error)
//	cnf.DefaultString(key string, defaultVal string) string
//	cnf.DefaultStrings(key string, defaultVal []string) []string
//	cnf.DefaultInt(key string, defaultVal int) int
//	cnf.DefaultInt64(key string, defaultVal int64) int64
//	cnf.DefaultBool(key string, defaultVal bool) bool
//	cnf.DefaultFloat(key string, defaultVal float64) float64
//	cnf.DIY(key string) (interface{}, error)
//	cnf.GetSection(section string) (map[string]string, error)
//	cnf.SaveConfigFile(filename string) error
package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Configer defines how to get and set value from configuration raw data.
type Configer interface {
	// Set support section::key type in given key when using ini type.
	Set(key, val string) error

	// String support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	String(key string) (string, error)
	// Strings get string slice
	Strings(key string) ([]string, error)
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	// DefaultString support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultString(key string, defaultVal string) string
	// DefaultStrings get string slice
	DefaultStrings(key string, defaultVal []string) []string
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64

	// DIY return the original value
	DIY(key string) (interface{}, error)

	GetSection(section string) (map[string]string, error)

	Unmarshaler(prefix string, obj interface{}, opt ...DecodeOption) error
	Sub(key string) (Configer, error)
	OnChange(key string, fn func(value string))
	SaveConfigFile(filename string) error
}

type BaseConfiger struct {
	// The reader should support key like "a.b.c"
	reader func(ctx context.Context, key string) (string, error)
}

func NewBaseConfiger(reader func(ctx context.Context, key string) (string, error)) BaseConfiger {
	return BaseConfiger{
		reader: reader,
	}
}

func (c *BaseConfiger) Int(key string) (int, error) {
	res, err := c.reader(context.TODO(), key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (c *BaseConfiger) Int64(key string) (int64, error) {
	res, err := c.reader(context.TODO(), key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(res, 10, 64)
}

func (c *BaseConfiger) Bool(key string) (bool, error) {
	res, err := c.reader(context.TODO(), key)
	if err != nil {
		return false, err
	}
	return ParseBool(res)
}

func (c *BaseConfiger) Float(key string) (float64, error) {
	res, err := c.reader(context.TODO(), key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(res, 64)
}

// DefaultString returns the string value for a given key.
// if err != nil or value is empty return defaultval
func (c *BaseConfiger) DefaultString(key string, defaultVal string) string {
	if res, err := c.String(key); res != "" && err == nil {
		return res
	}
	return defaultVal
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaultval
func (c *BaseConfiger) DefaultStrings(key string, defaultVal []string) []string {
	if res, err := c.Strings(key); len(res) > 0 && err == nil {
		return res
	}
	return defaultVal
}

func (c *BaseConfiger) DefaultInt(key string, defaultVal int) int {
	if res, err := c.Int(key); err == nil {
		return res
	}
	return defaultVal
}

func (c *BaseConfiger) DefaultInt64(key string, defaultVal int64) int64 {
	if res, err := c.Int64(key); err == nil {
		return res
	}
	return defaultVal
}

func (c *BaseConfiger) DefaultBool(key string, defaultVal bool) bool {
	if res, err := c.Bool(key); err == nil {
		return res
	}
	return defaultVal
}

func (c *BaseConfiger) DefaultFloat(key string, defaultVal float64) float64 {
	if res, err := c.Float(key); err == nil {
		return res
	}
	return defaultVal
}

func (c *BaseConfiger) String(key string) (string, error) {
	return c.reader(context.TODO(), key)
}

// Strings returns the []string value for a given key.
// Return nil if config value does not exist or is empty.
func (c *BaseConfiger) Strings(key string) ([]string, error) {
	res, err := c.String(key)
	if err != nil || res == "" {
		return nil, err
	}
	return strings.Split(res, ";"), nil
}

func (*BaseConfiger) Sub(string) (Configer, error) {
	return nil, errors.New("unsupported operation")
}

func (*BaseConfiger) OnChange(_ string, _ func(value string)) {
	// do nothing
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
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("config: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewConfig adapterName is ini/json/xml/yaml.
// filename is the config file path.
func NewConfig(adapterName, filename string) (Configer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	return adapter.Parse(filename)
}

// NewConfigData adapterName is ini/json/xml/yaml.
// data is the config data.
func NewConfigData(adapterName string, data []byte) (Configer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q (forgotten import?)", adapterName)
	}
	return adapter.ParseData(data)
}

// ExpandValueEnvForMap convert all string value with environment variable.
func ExpandValueEnvForMap(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch value := v.(type) {
		case string:
			m[k] = ExpandValueEnv(value)
		case map[string]interface{}:
			m[k] = ExpandValueEnvForMap(value)
		case map[string]string:
			for k2, v2 := range value {
				value[k2] = ExpandValueEnv(v2)
			}
			m[k] = value
		case map[interface{}]interface{}:
			tmp := make(map[string]interface{}, len(value))
			for k2, v2 := range value {
				tmp[k2.(string)] = v2
			}
			m[k] = ExpandValueEnvForMap(tmp)
		}
	}
	return m
}

// ExpandValueEnv returns value of convert with environment variable.
//
// Return environment variable if value start with "${" and end with "}".
// Return default value if environment variable is empty or not exist.
//
// It accept value formats "${env}" , "${env||}}" , "${env||defaultValue}" , "defaultvalue".
// Examples:
//
//	v1 := config.ExpandValueEnv("${GOPATH}")			// return the GOPATH environment variable.
//	v2 := config.ExpandValueEnv("${GOAsta||/usr/local/go}")	// return the default value "/usr/local/go/".
//	v3 := config.ExpandValueEnv("Astaxie")				// return the value "Astaxie".
func ExpandValueEnv(value string) (realValue string) {
	realValue = value

	vLen := len(value)
	// 3 = ${}
	if vLen < 3 {
		return
	}
	// Need start with "${" and end with "}", then return.
	if value[0] != '$' || value[1] != '{' || value[vLen-1] != '}' {
		return
	}

	key := ""
	defaultV := ""
	// value start with "${"
	for i := 2; i < vLen; i++ {
		if value[i] == '|' && (i+1 < vLen && value[i+1] == '|') {
			key = value[2:i]
			defaultV = value[i+2 : vLen-1] // other string is default value.
			break
		} else if value[i] == '}' {
			key = value[2:i]
			break
		}
	}

	realValue = os.Getenv(key)
	if realValue == "" {
		realValue = defaultV
	}

	return
}

// ParseBool returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func ParseBool(val interface{}) (value bool, err error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%d", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1.0 {
				return true, nil
			} else if v == 0.0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

// ToString converts values of any type to string.
func ToString(x interface{}) string {
	switch y := x.(type) {

	// Handle dates with special logic
	// This needs to come above the fmt.Stringer
	// test since time.Time's have a .String()
	// method
	case time.Time:
		return y.Format("A Monday")

	// Handle type string
	case string:
		return y

	// Handle type with .String() method
	case fmt.Stringer:
		return y.String()

	// Handle type with .Error() method
	case error:
		return y.Error()

	}

	// Handle named string type
	if v := reflect.ValueOf(x); v.Kind() == reflect.String {
		return v.String()
	}

	// Fallback to fmt package for anything else like numeric types
	return fmt.Sprint(x)
}

type DecodeOption func(options decodeOptions)

type decodeOptions struct{}
