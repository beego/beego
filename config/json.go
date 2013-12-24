package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
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
	x := &JsonConfigContainer{
		data: make(map[string]interface{}),
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &x.data)
	if err != nil {
		return nil, err
	}
	return x, nil
}

// A Config represents the json configuration.
// Only when get value, support key as section:name type.
type JsonConfigContainer struct {
	data map[string]interface{}
	sync.RWMutex
}

// Bool returns the boolean value for a given key.
func (c *JsonConfigContainer) Bool(key string) (bool, error) {
	val := c.getdata(key)
	if val != nil {
		if v, ok := val.(bool); ok {
			return v, nil
		} else {
			return false, errors.New("not bool value")
		}
	} else {
		return false, errors.New("not exist key:" + key)
	}
	return false, nil
}

// Int returns the integer value for a given key.
func (c *JsonConfigContainer) Int(key string) (int, error) {
	val := c.getdata(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int(v), nil
		} else {
			return 0, errors.New("not int value")
		}
	} else {
		return 0, errors.New("not exist key:" + key)
	}
	return 0, nil
}

// Int64 returns the int64 value for a given key.
func (c *JsonConfigContainer) Int64(key string) (int64, error) {
	val := c.getdata(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int64(v), nil
		} else {
			return 0, errors.New("not int64 value")
		}
	} else {
		return 0, errors.New("not exist key:" + key)
	}
	return 0, nil
}

// Float returns the float value for a given key.
func (c *JsonConfigContainer) Float(key string) (float64, error) {
	val := c.getdata(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return v, nil
		} else {
			return 0.0, errors.New("not float64 value")
		}
	} else {
		return 0.0, errors.New("not exist key:" + key)
	}
	return 0.0, nil
}

// String returns the string value for a given key.
func (c *JsonConfigContainer) String(key string) string {
	val := c.getdata(key)
	if val != nil {
		if v, ok := val.(string); ok {
			return v
		} else {
			return ""
		}
	} else {
		return ""
	}
	return ""
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
	val := c.getdata(key)
	if val != nil {
		return val, nil
	} else {
		return nil, errors.New("not exist key")
	}
	return nil, nil
}

// section.key or key
func (c *JsonConfigContainer) getdata(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	if len(key) == 0 {
		return nil
	}
	sectionkey := strings.Split(key, "::")
	if len(sectionkey) >= 2 {
		cruval, ok := c.data[sectionkey[0]]
		if !ok {
			return nil
		}
		for _, key := range sectionkey[1:] {
			if v, ok := cruval.(map[string]interface{}); !ok {
				return nil
			} else if cruval, ok = v[key]; !ok {
				return nil
			}
		}
		return cruval
	} else {
		if v, ok := c.data[key]; ok {
			return v
		}
	}
	return nil
}

func init() {
	Register("json", &JsonConfig{})
}
