// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

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
	key = strings.ToLower(key)
	return c.data[key]
}

func (c *fakeConfigContainer) Set(key, val string) error {
	key = strings.ToLower(key)
	c.data[key] = val
	return nil
}

func (c *fakeConfigContainer) String(key string) string {
	return c.getData(key)
}

func (c *fakeConfigContainer) Strings(key string) []string {
	return strings.Split(c.getData(key), ";")
}

func (c *fakeConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.getData(key))
}

func (c *fakeConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getData(key), 10, 64)
}

func (c *fakeConfigContainer) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.getData(key))
}

func (c *fakeConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getData(key), 64)
}

func (c *fakeConfigContainer) DIY(key string) (interface{}, error) {
	key = strings.ToLower(key)
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("key not find")
}

var _ ConfigContainer = new(fakeConfigContainer)

func NewFakeConfig() ConfigContainer {
	return &fakeConfigContainer{
		data: make(map[string]string),
	}
}
