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
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) Set(key, val string) error {
	c.data[strings.ToLower(key)] = val
	return nil
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) String(key string) string {
	return c.getData(key)
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) DefaultString(key string, defaultval string) string {
	v := c.String(key)
	if v == "" {
		return defaultval
	}
	return v
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) Strings(key string) []string {
	v := c.String(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) DefaultStrings(key string, defaultval []string) []string {
	v := c.Strings(key)
	if v == nil {
		return defaultval
	}
	return v
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) Int(key string) (int, error) {
	return strconv.Atoi(c.getData(key))
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) DefaultInt(key string, defaultval int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultval
	}
	return v
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getData(key), 10, 64)
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultval
	}
	return v
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) Bool(key string) (bool, error) {
	return ParseBool(c.getData(key))
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) DefaultBool(key string, defaultval bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultval
	}
	return v
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getData(key), 64)
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) DIY(key string) (interface{}, error) {
	if v, ok := c.data[strings.ToLower(key)]; ok {
		return v, nil
	}
	return nil, errors.New("key not find")
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) GetSection(section string) (map[string]string, error) {
	return nil, errors.New("not implement in the fakeConfigContainer")
}
// Deprecated: using pkg/config, we will delete this in v2.1.0
func (c *fakeConfigContainer) SaveConfigFile(filename string) error {
	return errors.New("not implement in the fakeConfigContainer")
}

var _ Configer = new(fakeConfigContainer)

// NewFakeConfig return a fake Configer
// Deprecated: using pkg/config, we will delete this in v2.1.0
func NewFakeConfig() Configer {
	return &fakeConfigContainer{
		data: make(map[string]string),
	}
}
