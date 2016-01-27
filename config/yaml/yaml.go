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

// Package yaml for config provider
//
// depend on github.com/beego/goyaml2
//
// go install github.com/beego/goyaml2
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/config/yaml"
//   "github.com/astaxie/beego/config"
// )
//
//  cnf, err := config.NewConfig("yaml", "config.yaml")
//
//  more docs http://beego.me/docs/module/config.md
package yaml

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/beego/goyaml2"
)

// Config is a yaml config parser and implements Config interface.
type Config struct{}

// Parse returns a ConfigContainer with parsed yaml config map.
func (yaml *Config) Parse(filename string) (y config.Configer, err error) {
	cnf, err := ReadYmlReader(filename)
	if err != nil {
		return
	}
	y = &ConfigContainer{
		data: cnf,
	}
	return
}

// ParseData parse yaml data
func (yaml *Config) ParseData(data []byte) (config.Configer, error) {
	// Save memory data to temporary file
	tmpName := path.Join(os.TempDir(), "beego", fmt.Sprintf("%d", time.Now().Nanosecond()))
	os.MkdirAll(path.Dir(tmpName), os.ModePerm)
	if err := ioutil.WriteFile(tmpName, data, 0655); err != nil {
		return nil, err
	}
	return yaml.Parse(tmpName)
}

// ReadYmlReader Read yaml file to map.
// if json like, use json package, unless goyaml2 package.
func ReadYmlReader(path string) (cnf map[string]interface{}, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil || len(buf) < 3 {
		return
	}

	if string(buf[0:1]) == "{" {
		log.Println("Look like a Json, try json umarshal")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			log.Println("It is Json Map")
			return
		}
	}

	data, err := goyaml2.Read(bytes.NewBuffer(buf))
	if err != nil {
		log.Println("Goyaml2 ERR>", string(buf), err)
		return
	}

	if data == nil {
		log.Println("Goyaml2 output nil? Pls report bug\n" + string(buf))
		return
	}
	cnf, ok := data.(map[string]interface{})
	if !ok {
		log.Println("Not a Map? >> ", string(buf), data)
		cnf = nil
	}
	return
}

// ConfigContainer A Config represents the yaml configuration.
type ConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

// Bool returns the boolean value for a given key.
func (c *ConfigContainer) Bool(key string) (bool, error) {

	if v, err := c.getData(key); err != nil {
		return false, err
	} else {
		return config.ParseBool(v)
	}
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultval
func (c *ConfigContainer) DefaultBool(key string, defaultval bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int returns the integer value for a given key.
func (c *ConfigContainer) Int(key string) (int, error) {
	if v, err := c.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int); ok {
		return vv, nil
	} else if vv, ok := v.(int64); ok {
		return int(vv), nil
	}
	return 0, errors.New("not int value")
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaltval
func (c *ConfigContainer) DefaultInt(key string, defaultval int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *ConfigContainer) Int64(key string) (int64, error) {
	if v, err := c.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int64); ok {
		return vv, nil
	}
	return 0, errors.New("not bool value")
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaltval
func (c *ConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Float returns the float value for a given key.
func (c *ConfigContainer) Float(key string) (float64, error) {
	if v, err := c.getData(key); err != nil {
		return 0.0, err
	} else if vv, ok := v.(float64); ok {
		return vv, nil
	} else if vv, ok := v.(int); ok {
		return float64(vv), nil
	} else if vv, ok := v.(int64); ok {
		return float64(vv), nil
	}
	return 0.0, errors.New("not float64 value")
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaltval
func (c *ConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}

// String returns the string value for a given key.
func (c *ConfigContainer) String(key string) string {
	if v, err := c.getData(key); err == nil {
		if vv, ok := v.(string); ok {
			return vv
		}
	}
	return ""
}

// DefaultString returns the string value for a given key.
// if err != nil return defaltval
func (c *ConfigContainer) DefaultString(key string, defaultval string) string {
	v := c.String(key)
	if v == "" {
		return defaultval
	}
	return v
}

// Strings returns the []string value for a given key.
func (c *ConfigContainer) Strings(key string) []string {
	return strings.Split(c.String(key), ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaltval
func (c *ConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	v := c.Strings(key)
	if len(v) == 0 {
		return defaultval
	}
	return v
}

// GetSection returns map for the given section
func (c *ConfigContainer) GetSection(section string) (map[string]string, error) {
	v, ok := c.data[section]
	if ok {
		var interfaceToStr func(map[string]interface{}) map[string]string

		interfaceToStr = func(values map[string]interface{}) map[string]string {
			strValues := make(map[string]string, len(values))
			for k, vv := range values {
				if vv == nil {
					strValues[k] = ""
				} else if env, ok := config.Getenv(vv); ok {
					strValues[k] = env
				} else if str, ok := vv.(string); ok {
					strValues[k] = str
				} else if m, ok := vv.(map[string]interface{}); ok {
					strValues[k] = fmt.Sprintf("%v", interfaceToStr(m))
				} else {
					// TODO: no better.
					strValues[k] = fmt.Sprintf("%v", vv)
				}
			}
			return strValues
		}
		return interfaceToStr(v.(map[string]interface{})), nil
	}
	return nil, errors.New("not exist setction")
}

// SaveConfigFile save the config into file
func (c *ConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = goyaml2.Write(f, c.data)
	return err
}

// Set writes a new value for key.
func (c *ConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

// DIY returns the raw value by a given key.
func (c *ConfigContainer) DIY(key string) (v interface{}, err error) {
	return c.getData(key)
}

func (c *ConfigContainer) getData(key string) (interface{}, error) {

	if len(key) == 0 {
		return nil, errors.New("key is emtpy")
	}

	if v, ok := c.data[key]; ok {
		if env, ok := config.Getenv(v); ok {
			return env, nil
		} else {
			return v, nil
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}

func init() {
	config.Register("yaml", &Config{})
}
