package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type JsonConfig struct {
}

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

type JsonConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

func (c *JsonConfigContainer) Bool(key string) (bool, error) {
	if v, ok := c.data[key].(bool); ok {
		return v, nil
	}
	return false, errors.New("not bool value")
}

func (c *JsonConfigContainer) Int(key string) (int, error) {
	if v, ok := c.data[key].(float64); ok {
		return int(v), nil
	}
	return 0, errors.New("not int value")
}

func (c *JsonConfigContainer) Int64(key string) (int64, error) {
	if v, ok := c.data[key].(float64); ok {
		return int64(v), nil
	}
	return 0, errors.New("not bool value")
}

func (c *JsonConfigContainer) Float(key string) (float64, error) {
	if v, ok := c.data[key].(float64); ok {
		return v, nil
	}
	return 0.0, errors.New("not float64 value")
}

func (c *JsonConfigContainer) String(key string) string {
	if v, ok := c.data[key].(string); ok {
		return v
	}
	return ""
}

func (c *JsonConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

func (c *JsonConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist key")
}

func init() {
	Register("json", &JsonConfig{})
}
