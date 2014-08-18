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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// JsonConfig is a json config parser and implements Config interface.
type JsonConfig struct {
}

// Parse returns a ConfigContainer with parsed json config map.
func (js *JsonConfig) Parse(filename string) (ConfigContainer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	x := &JsonConfigContainer{
		data: make(map[string]interface{}),
	}
	err = json.Unmarshal(content, &x.data)
	if err != nil {
		var wrappingArray []interface{}
		err2 := json.Unmarshal(content, &wrappingArray)
		if err2 != nil {
			return nil, err
		}
		x.data["rootArray"] = wrappingArray
	}
	return x, nil
}

func (js *JsonConfig) ParseData(data []byte) (ConfigContainer, error) {
	// Save memory data to temporary file
	tmpName := path.Join(os.TempDir(), "beego", fmt.Sprintf("%d", time.Now().Nanosecond()))
	os.MkdirAll(path.Dir(tmpName), os.ModePerm)
	if err := ioutil.WriteFile(tmpName, data, 0655); err != nil {
		return nil, err
	}
	return js.Parse(tmpName)
}

// A Config represents the json configuration.
// Only when get value, support key as section:name type.
type JsonConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

// Bool returns the boolean value for a given key.
func (c *JsonConfigContainer) Bool(key string) (bool, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(bool); ok {
			return v, nil
		}
		return false, errors.New("not bool value")
	}
	return false, errors.New("not exist key:" + key)
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultval
func (c *JsonConfigContainer) DefaultBool(key string, defaultval bool) bool {
	if v, err := c.Bool(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// Int returns the integer value for a given key.
func (c *JsonConfigContainer) Int(key string) (int, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int(v), nil
		}
		return 0, errors.New("not int value")
	}
	return 0, errors.New("not exist key:" + key)
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaltval
func (c *JsonConfigContainer) DefaultInt(key string, defaultval int) int {
	if v, err := c.Int(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// Int64 returns the int64 value for a given key.
func (c *JsonConfigContainer) Int64(key string) (int64, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int64(v), nil
		}
		return 0, errors.New("not int64 value")
	}
	return 0, errors.New("not exist key:" + key)
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaltval
func (c *JsonConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	if v, err := c.Int64(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// Float returns the float value for a given key.
func (c *JsonConfigContainer) Float(key string) (float64, error) {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return v, nil
		}
		return 0.0, errors.New("not float64 value")
	}
	return 0.0, errors.New("not exist key:" + key)
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaltval
func (c *JsonConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	if v, err := c.Float(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// String returns the string value for a given key.
func (c *JsonConfigContainer) String(key string) string {
	val := c.getData(key)
	if val != nil {
		if v, ok := val.(string); ok {
			return v
		}
	}
	return ""
}

// DefaultString returns the string value for a given key.
// if err != nil return defaltval
func (c *JsonConfigContainer) DefaultString(key string, defaultval string) string {
	if v := c.String(key); v == "" {
		return defaultval
	} else {
		return v
	}
}

// Strings returns the []string value for a given key.
func (c *JsonConfigContainer) Strings(key string) []string {
	return strings.Split(c.String(key), ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaltval
func (c *JsonConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	if v := c.Strings(key); len(v) == 0 {
		return defaultval
	} else {
		return v
	}
}

// GetSection returns map for the given section
func (c *JsonConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section]; ok {
		return v.(map[string]string), nil
	} else {
		return nil, errors.New("not exist setction")
	}
}

// SaveConfigFile save the config into file
func (c *JsonConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

// WriteValue writes a new value for key.
func (c *JsonConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

// DIY returns the raw value by a given key.
func (c *JsonConfigContainer) DIY(key string) (v interface{}, err error) {
	val := c.getData(key)
	if val != nil {
		return val, nil
	}
	return nil, errors.New("not exist key")
}

// section.key or key
func (c *JsonConfigContainer) getData(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	if len(key) == 0 {
		return nil
	}
	sectionKey := strings.Split(key, "::")
	if len(sectionKey) >= 2 {
		curValue, ok := c.data[sectionKey[0]]
		if !ok {
			return nil
		}
		for _, key := range sectionKey[1:] {
			if v, ok := curValue.(map[string]interface{}); ok {
				if curValue, ok = v[key]; !ok {
					return nil
				}
			}
		}
		return curValue
	}
	if v, ok := c.data[key]; ok {
		return v
	}
	return nil
}

func init() {
	Register("json", &JsonConfig{})
}
