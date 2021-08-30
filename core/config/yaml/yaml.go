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
// Usage:
//  import(
//   _ "github.com/beego/beego/v2/core/config/yaml"
//     "github.com/beego/beego/v2/core/config"
//  )
//
//  cnf, err := config.NewConfig("yaml", "config.yaml")
//
// More docs http://beego.me/docs/module/config.md
package yaml

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

// Config is a yaml config parser and implements Config interface.
type Config struct{}

// Parse returns a ConfigContainer with parsed yaml config map.
func (*Config) Parse(filename string) (y config.Configer, err error) {
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
func (*Config) ParseData(data []byte) (config.Configer, error) {
	cnf, err := parseYML(data)
	if err != nil {
		return nil, err
	}

	return &ConfigContainer{
		data: cnf,
	}, nil
}

// ReadYmlReader Read yaml file to map.
func ReadYmlReader(path string) (cnf map[string]interface{}, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	return parseYML(buf)
}

// parseYML parse yaml formatted []byte to map.
func parseYML(buf []byte) (map[string]interface{}, error) {
	cnf := make(map[string]interface{})
	err := yaml.Unmarshal(buf, cnf)
	if err != nil {
		return nil, err
	}
	cnf = config.ExpandValueEnvForMap(cnf)
	return cnf, err
}

// ConfigContainer is a config which represents the yaml configuration.
type ConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

// Unmarshaler is similar to Sub
func (c *ConfigContainer) Unmarshaler(prefix string, obj interface{}, _ ...config.DecodeOption) error {
	sub, err := c.subMap(prefix)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(sub)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, obj)
}

func (c *ConfigContainer) Sub(key string) (config.Configer, error) {
	sub, err := c.subMap(key)
	if err != nil {
		return nil, err
	}
	return &ConfigContainer{
		data: sub,
	}, nil
}

func (c *ConfigContainer) subMap(key string) (map[string]interface{}, error) {
	tmpData := c.data
	keys := strings.Split(key, ".")
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch val := v.(type) {
			case map[string]interface{}:
				tmpData = val
				if idx == len(keys)-1 {
					return tmpData, nil
				}
			default:
				return nil, fmt.Errorf("the key is invalid: %s", key)
			}
		}
	}

	return tmpData, nil
}

func (*ConfigContainer) OnChange(_ string, _ func(value string)) {
	// do nothing
	logs.Warn("Unsupported operation: OnChange")
}

// Bool returns the boolean value for a given key.
func (c *ConfigContainer) Bool(key string) (bool, error) {
	v, err := c.getData(key)
	if err != nil {
		return false, err
	}
	return config.ParseBool(v)
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultVal
func (c *ConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultVal
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
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultInt(key string, defaultVal int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultVal
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *ConfigContainer) Int64(key string) (int64, error) {
	v, err := c.getData(key)
	if err != nil {
		return 0, err
	}
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int64:
		return val, nil
	default:
		return 0, errors.New("not int or int64 value")
	}
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultVal
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
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultVal
	}
	return v
}

// String returns the string value for a given key.
func (c *ConfigContainer) String(key string) (string, error) {
	if v, err := c.getData(key); err == nil {
		if vv, ok := v.(string); ok {
			return vv, nil
		}
	}
	return "", nil
}

// DefaultString returns the string value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultString(key string, defaultVal string) string {
	v, err := c.String(key)
	if v == "" || err != nil {
		return defaultVal
	}
	return v
}

// Strings returns the []string value for a given key.
func (c *ConfigContainer) Strings(key string) ([]string, error) {
	v, err := c.String(key)
	if v == "" || err != nil {
		return nil, err
	}
	return strings.Split(v, ";"), nil
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	v, err := c.Strings(key)
	if v == nil || err != nil {
		return defaultVal
	}
	return v
}

// GetSection returns map for the given section
func (c *ConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section]; ok {
		switch val := v.(type) {
		case map[string]interface{}:
			res := make(map[string]string, len(val))
			for k2, v2 := range val {
				res[k2] = fmt.Sprintf("%v", v2)
			}
			return res, nil
		case map[string]string:
			return val, nil
		default:
			return nil, fmt.Errorf("unexpected type: %v", v)
		}
	}
	return nil, errors.New("not exist section")
}

// SaveConfigFile save the config into file
func (c *ConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	buf, err := yaml.Marshal(c.data)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
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
	if key == "" {
		return nil, errors.New("key is empty")
	}
	c.RLock()
	defer c.RUnlock()

	keys := strings.Split(key, ".")
	tmpData := c.data
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch val := v.(type) {
			case map[string]interface{}:
				tmpData = val
				if idx == len(keys)-1 {
					return tmpData, nil
				}
			default:
				return v, nil
			}
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}

func init() {
	config.Register("yaml", &Config{})
}
