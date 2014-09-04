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

// package xml for config provider
//
// depend on github.com/beego/x2j
//
// go install github.com/beego/x2j
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/config/xml"
//   "github.com/astaxie/beego/config"
// )
//
//  cnf, err := config.NewConfig("xml", "config.xml")
//
//  more docs http://beego.me/docs/module/config.md
package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/beego/x2j"
)

// XmlConfig is a xml config parser and implements Config interface.
// xml configurations should be included in <config></config> tag.
// only support key/value pair as <key>value</key> as each item.
type XMLConfig struct{}

// Parse returns a ConfigContainer with parsed xml config map.
func (xc *XMLConfig) Parse(filename string) (config.ConfigContainer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	x := &XMLConfigContainer{data: make(map[string]interface{})}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	d, err := x2j.DocToMap(string(content))
	if err != nil {
		return nil, err
	}

	x.data = d["config"].(map[string]interface{})
	return x, nil
}

func (x *XMLConfig) ParseData(data []byte) (config.ConfigContainer, error) {
	// Save memory data to temporary file
	tmpName := path.Join(os.TempDir(), "beego", fmt.Sprintf("%d", time.Now().Nanosecond()))
	os.MkdirAll(path.Dir(tmpName), os.ModePerm)
	if err := ioutil.WriteFile(tmpName, data, 0655); err != nil {
		return nil, err
	}
	return x.Parse(tmpName)
}

// A Config represents the xml configuration.
type XMLConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

// Bool returns the boolean value for a given key.
func (c *XMLConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.data[key].(string))
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultval
func (c *XMLConfigContainer) DefaultBool(key string, defaultval bool) bool {
	if v, err := c.Bool(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// Int returns the integer value for a given key.
func (c *XMLConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key].(string))
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaltval
func (c *XMLConfigContainer) DefaultInt(key string, defaultval int) int {
	if v, err := c.Int(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// Int64 returns the int64 value for a given key.
func (c *XMLConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.data[key].(string), 10, 64)
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaltval
func (c *XMLConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	if v, err := c.Int64(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// Float returns the float value for a given key.
func (c *XMLConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key].(string), 64)
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaltval
func (c *XMLConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	if v, err := c.Float(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

// String returns the string value for a given key.
func (c *XMLConfigContainer) String(key string) string {
	if v, ok := c.data[key].(string); ok {
		return v
	}
	return ""
}

// DefaultString returns the string value for a given key.
// if err != nil return defaltval
func (c *XMLConfigContainer) DefaultString(key string, defaultval string) string {
	if v := c.String(key); v == "" {
		return defaultval
	} else {
		return v
	}
}

// Strings returns the []string value for a given key.
func (c *XMLConfigContainer) Strings(key string) []string {
	return strings.Split(c.String(key), ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaltval
func (c *XMLConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	if v := c.Strings(key); len(v) == 0 {
		return defaultval
	} else {
		return v
	}
}

// GetSection returns map for the given section
func (c *XMLConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section]; ok {
		return v.(map[string]string), nil
	} else {
		return nil, errors.New("not exist setction")
	}
}

// SaveConfigFile save the config into file
func (c *XMLConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := xml.MarshalIndent(c.data, "  ", "    ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

// WriteValue writes a new value for key.
func (c *XMLConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

// DIY returns the raw value by a given key.
func (c *XMLConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist key")
}

func init() {
	config.Register("xml", &XMLConfig{})
}
