package base

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/encoder"
)

// ConfigBaseContainer A ConfigFile represents the json configuration.
type ConfigBaseContainer struct {
	Data          map[string]interface{}
	SeparatorKeys string
	sync.RWMutex
}

func (c *ConfigBaseContainer) keyJoin(keys []string) string {
	if len(keys) == 0 {
		panic("must provide some key")
	}
	return strings.Join(keys, c.SeparatorKeys)
}

// Bool returns the boolean value for a given key.
func (c *ConfigBaseContainer) Bool(keys ...string) (bool, error) {
	v, err := c.getData(c.keyJoin(keys))
	if err != nil {
		return false, err
	}
	return config.ParseBool(v)
}

// DefaultBool return the bool value if has no error
// otherwise return the defaultVal
func (c *ConfigBaseContainer) DefaultBool(defaultVal bool, keys ...string) bool {
	v, err := c.Bool(c.keyJoin(keys))
	if err != nil {
		return defaultVal
	}
	return v
}

// Int returns the integer value for a given key.
func (c *ConfigBaseContainer) Int(keys ...string) (int, error) {
	if v, err := c.getData(c.keyJoin(keys)); err != nil {
		return 0, err
	} else {
		switch v.(type) {
		case int:
			return v.(int), nil
		case int64:
			return int(v.(int64)), nil
		case string:
			if vv, ok := v.(int); ok {
				return vv, nil
			}
		}
	}
	return 0, errors.New("not int value")
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaultVal
func (c *ConfigBaseContainer) DefaultInt(defaultVal int, keys ...string) int {
	v, err := c.Int(c.keyJoin(keys))
	if err != nil {
		return defaultVal
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *ConfigBaseContainer) Int64(keys ...string) (int64, error) {
	if v, err := c.getData(c.keyJoin(keys)); err != nil {
		return 0, err
	} else {
		switch v.(type) {
		case interface{}:
			sVal := fmt.Sprintf("%v", v)
			i, _ := strconv.Atoi(sVal)
			return int64(i), nil
		case int:
			return int64(v.(int)), nil
		case int64:
			return v.(int64), nil
		case string:
			if vv, ok := v.(int64); ok {
				return vv, nil
			}
		}
	}
	return 0, errors.New("not int64 value")
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaultVal
func (c *ConfigBaseContainer) DefaultInt64(defaultVal int64, keys ...string) int64 {
	v, err := c.Int64(c.keyJoin(keys))
	if err != nil {
		return defaultVal
	}
	return v
}

// Float returns the float value for a given key.
func (c *ConfigBaseContainer) Float(keys ...string) (float64, error) {
	if v, err := c.getData(c.keyJoin(keys)); err != nil {
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
func (c *ConfigBaseContainer) DefaultFloat(defaultVal float64, keys ...string) float64 {
	v, err := c.Float(c.keyJoin(keys))
	if err != nil {
		return defaultVal
	}
	return v
}

// String returns the string value for a given key.
func (c *ConfigBaseContainer) String(keys ...string) string {
	if v, err := c.getData(c.keyJoin(keys)); err == nil {
		if vv, ok := v.(string); ok {
			return vv
		}
	}
	return ""
}

// DefaultString returns the string value for a given key.
// if err != nil return defaultVal
func (c *ConfigBaseContainer) DefaultString(defaultVal string, keys ...string) string {
	v := c.String(c.keyJoin(keys))
	if v == "" {
		return defaultVal
	}
	return v
}

// Strings returns the []string value for a given key.
func (c *ConfigBaseContainer) Strings(keys ...string) []string {
	v := c.String(c.keyJoin(keys))
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaultVal
func (c *ConfigBaseContainer) DefaultStrings(defaultVal []string, keys ...string) []string {
	v := c.Strings(c.keyJoin(keys))
	if v == nil {
		return defaultVal
	}
	return v
}

// GetSection returns map for the given section
func (c *ConfigBaseContainer) GetSection(section string) (map[string]string, error) {

	if v, ok := c.Data[section]; ok {
		return v.(map[string]string), nil
	}
	return nil, errors.New("not exist section")
}

// SaveConfigFile save the config into file
func (c *ConfigBaseContainer) SaveConfigFile(filename string, encode encoder.Encoder) (err error) {
	// Write configuration file by filename.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := encode.Encode(c.Data)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

// Set writes a new value for key.
func (c *ConfigBaseContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.Data[key] = val
	return nil
}

// DIY returns the raw value by a given key.
func (c *ConfigBaseContainer) DIY(keys ...string) (v interface{}, err error) {
	return c.getData(c.keyJoin(keys))
}

func (c *ConfigBaseContainer) getRealKeys(key string) []string {
	keys := strings.Split(key, c.SeparatorKeys)

	var realKeys []string

	for _, myKey := range keys {
		isPresent := false
		for realKey := range c.Data {
			if strings.ToLower(realKey) == strings.ToLower(myKey) {
				realKeys = append(realKeys, realKey)
				isPresent = true
				break
			}
		}

		if !isPresent {
			realKeys = append(realKeys, myKey)
		}
	}

	return realKeys
}

func (c *ConfigBaseContainer) getData(key string) (interface{}, error) {
	if len(key) == 0 {
		return nil, errors.New("key is empty")
	}
	c.RLock()
	defer c.RUnlock()

	keys := c.getRealKeys(key)
	tmpData := c.Data
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {

			switch v.(type) {
			case map[interface{}]interface{}:
				{
					tmpData = config.ToStringKeyMap(v.(map[interface{}]interface{}))
					if idx == len(keys)-1 {
						return tmpData, nil
					}
				}
			case map[string]interface{}:
				{
					tmpData = v.(map[string]interface{})
					if idx == len(keys)-1 {
						return tmpData, nil
					}
				}
			default:
				{
					return v, nil
				}

			}
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}

var _ config.Configer = NewBaseConfig()

// NewBaseonfig return a base Configer
func NewBaseConfig() *ConfigBaseContainer {
	return &ConfigBaseContainer{
		Data:          make(map[string]interface{}),
		SeparatorKeys: "::",
	}
}
