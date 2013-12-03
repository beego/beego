//xml parse should incluce in <config></config> tags

package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/beego/x2j"
)

type XMLConfig struct {
}

func (xmls *XMLConfig) Parse(filename string) (ConfigContainer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	x := &XMLConfigContainer{
		data: make(map[string]interface{}),
	}
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

type XMLConfigContainer struct {
	data map[string]interface{}
	sync.Mutex
}

func (c *XMLConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.data[key].(string))
}

func (c *XMLConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key].(string))
}

func (c *XMLConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.data[key].(string), 10, 64)
}

func (c *XMLConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key].(string), 64)
}

func (c *XMLConfigContainer) String(key string) string {
	if v, ok := c.data[key].(string); ok {
		return v
	}
	return ""
}

func (c *XMLConfigContainer) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	return nil
}

func (c *XMLConfigContainer) DIY(key string) (v interface{}, err error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("not exist key")
}

func init() {
	Register("xml", &XMLConfig{})
}
