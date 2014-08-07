// Beego (http://beego.me/)
//
// @description beego is an open-source, high-performance web framework for the Go programming language.
//
// @link        http://github.com/astaxie/beego for the canonical source repository
//
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
//
// @authors     astaxie
package config

import (
	"errors"
	"strconv"
	"strings"
)

type fakeConfigContainer struct {
	data map[string]string
}

func (c *fakeConfigContainer) getData(key string) string {
	return c.data[strings.ToLower(key)]
}

func (c *fakeConfigContainer) Set(key, val string) error {
	c.data[strings.ToLower(key)] = val
	return nil
}

func (c *fakeConfigContainer) String(key string) string {
	return c.getData(key)
}

func (c *fakeConfigContainer) DefaultString(key string, defaultval string) string {
	if v := c.getData(key); v == "" {
		return defaultval
	} else {
		return v
	}
}

func (c *fakeConfigContainer) Strings(key string) []string {
	return strings.Split(c.getData(key), ";")
}

func (c *fakeConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	if v := c.Strings(key); len(v) == 0 {
		return defaultval
	} else {
		return v
	}
}

func (c *fakeConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.getData(key))
}

func (c *fakeConfigContainer) DefaultInt(key string, defaultval int) int {
	if v, err := c.Int(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *fakeConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getData(key), 10, 64)
}

func (c *fakeConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	if v, err := c.Int64(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *fakeConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.getData(key))
}

func (c *fakeConfigContainer) DefaultBool(key string, defaultval bool) bool {
	if v, err := c.Bool(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *fakeConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getData(key), 64)
}

func (c *fakeConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	if v, err := c.Float(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *fakeConfigContainer) DIY(key string) (interface{}, error) {
	if v, ok := c.data[strings.ToLower(key)]; ok {
		return v, nil
	}
	return nil, errors.New("key not find")
}

func (c *fakeConfigContainer) GetSection(section string) (map[string]string, error) {
	return nil, errors.New("not implement in the fakeConfigContainer")
}

func (c *fakeConfigContainer) SaveConfigFile(filename string) error {
	return errors.New("not implement in the fakeConfigContainer")
}

var _ ConfigContainer = new(fakeConfigContainer)

func NewFakeConfig() ConfigContainer {
	return &fakeConfigContainer{
		data: make(map[string]string),
	}
}
