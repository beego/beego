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

// Package Kubernetes ConfigMap for config provider.
//
// depend on:
// "k8s.io/apimachinery/pkg/apis/meta/v1"
// "k8s.io/client-go/kubernetes"
// "k8s.io/client-go/rest"
//
// Usage:
//  import(
//    _ "github.com/astaxie/beego/config/configmap"
//      "github.com/astaxie/beego/config"
//  )
//
//  cnf, err := config.NewConfig("configMap", "myConfAppName")
//
//More docs http://beego.me/docs/module/config.md
package configmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego/config"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getClient() (*kubernetes.Clientset, error) {
	conf, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(conf)
}

func makeMap(kv map[string]string) map[string]interface{} {

	data := make(map[string]interface{})

	for k, v := range kv {
		data[k] = make(map[string]interface{})

		vals := strings.Split(v, "\n")

		mp := make(map[string]interface{})
		for _, h := range vals {
			m, n := split(string(h), "=")
			mp[m] = n
		}

		data[k] = mp
	}

	return data
}

func split(s string, sp string) (k string, v string) {
	kv := strings.Split(s, sp)
	return kv[0], kv[1]
}

// Config is a json config parser and implements Config interface.
type Config struct {
	client    *kubernetes.Clientset
	cerr      error
	name      string
	namespace string
}

// Parse returns a ConfigContainer with parsed json config map.
func (k *Config) Parse(name string) (config.Configer, error) {
	if k.cerr != nil {
		return nil, k.cerr
	}

	cmp, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(name, v1.GetOptions{})

	if err != nil {
		return nil, err
	}

	data := makeMap(cmp.Data)

	b, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	return k.ParseData(b)
}

// ParseData json data
func (k *Config) ParseData(data []byte) (config.Configer, error) {
	x := &ConfigContainer{data: make(map[string]interface{})}

	cnf := map[string]interface{}{}
	_ = json.Unmarshal(data, &cnf)

	x.data = config.ExpandValueEnvForMap(cnf)

	return x, nil
}

// ConfigContainer A Config represents the json configuration.
type ConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

// Bool returns the boolean value for a given key.
func (c *ConfigContainer) Bool(key string) (bool, error) {
	if v := c.data[key]; v != nil {
		return config.ParseBool(v)
	}
	return false, fmt.Errorf("not exist key: %q", key)
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
	return strconv.Atoi(c.data[key].(string))
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
	return strconv.ParseInt(c.data[key].(string), 10, 64)
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
	return strconv.ParseFloat(c.data[key].(string), 64)
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
func (c *ConfigContainer) String(key string) string {
	if v, ok := c.data[key].(string); ok {
		return v
	}
	return ""
}

// DefaultString returns the string value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultString(key string, defaultVal string) string {
	v := c.String(key)
	if v == "" {
		return defaultVal
	}
	return v
}

// Strings returns the []string value for a given key.
func (c *ConfigContainer) Strings(key string) []string {
	v := c.String(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaultVal
func (c *ConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	v := c.Strings(key)
	if v == nil {
		return defaultVal
	}
	return v
}

// GetSection returns map for the given section
func (c *ConfigContainer) GetSection(section string) (map[string]string, error) {
	if v, ok := c.data[section].(map[string]interface{}); ok {
		mapstr := make(map[string]string)
		for k, val := range v {
			mapstr[k] = config.ToString(val)
		}
		return mapstr, nil
	}
	return nil, fmt.Errorf("section '%s' not found", section)
}

// SaveConfigFile save the config into file
func (c *ConfigContainer) SaveConfigFile(filename string) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(c.data, "  ", "    ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
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
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist key")
}

func init() {
	client, err := getClient()

	config.Register("configMap", &Config{
		name:      "beego",
		namespace: "default",
		cerr:      err,
		client:    client,
	})
}
