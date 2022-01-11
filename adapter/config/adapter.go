// Copyright 2020
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/pkg/errors"

	"github.com/beego/beego/v2/core/config"
)

type newToOldConfigerAdapter struct {
	delegate config.Configer
}

func (c *newToOldConfigerAdapter) Set(key, val string) error {
	return c.delegate.Set(key, val)
}

func (c *newToOldConfigerAdapter) String(key string) string {
	res, _ := c.delegate.String(key)
	return res
}

func (c *newToOldConfigerAdapter) Strings(key string) []string {
	res, _ := c.delegate.Strings(key)
	return res
}

func (c *newToOldConfigerAdapter) Int(key string) (int, error) {
	return c.delegate.Int(key)
}

func (c *newToOldConfigerAdapter) Int64(key string) (int64, error) {
	return c.delegate.Int64(key)
}

func (c *newToOldConfigerAdapter) Bool(key string) (bool, error) {
	return c.delegate.Bool(key)
}

func (c *newToOldConfigerAdapter) Float(key string) (float64, error) {
	return c.delegate.Float(key)
}

func (c *newToOldConfigerAdapter) DefaultString(key string, defaultVal string) string {
	return c.delegate.DefaultString(key, defaultVal)
}

func (c *newToOldConfigerAdapter) DefaultStrings(key string, defaultVal []string) []string {
	return c.delegate.DefaultStrings(key, defaultVal)
}

func (c *newToOldConfigerAdapter) DefaultInt(key string, defaultVal int) int {
	return c.delegate.DefaultInt(key, defaultVal)
}

func (c *newToOldConfigerAdapter) DefaultInt64(key string, defaultVal int64) int64 {
	return c.delegate.DefaultInt64(key, defaultVal)
}

func (c *newToOldConfigerAdapter) DefaultBool(key string, defaultVal bool) bool {
	return c.delegate.DefaultBool(key, defaultVal)
}

func (c *newToOldConfigerAdapter) DefaultFloat(key string, defaultVal float64) float64 {
	return c.delegate.DefaultFloat(key, defaultVal)
}

func (c *newToOldConfigerAdapter) DIY(key string) (interface{}, error) {
	return c.delegate.DIY(key)
}

func (c *newToOldConfigerAdapter) GetSection(section string) (map[string]string, error) {
	return c.delegate.GetSection(section)
}

func (c *newToOldConfigerAdapter) SaveConfigFile(filename string) error {
	return c.delegate.SaveConfigFile(filename)
}

type oldToNewConfigerAdapter struct {
	delegate Configer
}

func (o *oldToNewConfigerAdapter) Set(key, val string) error {
	return o.delegate.Set(key, val)
}

func (o *oldToNewConfigerAdapter) String(key string) (string, error) {
	return o.delegate.String(key), nil
}

func (o *oldToNewConfigerAdapter) Strings(key string) ([]string, error) {
	return o.delegate.Strings(key), nil
}

func (o *oldToNewConfigerAdapter) Int(key string) (int, error) {
	return o.delegate.Int(key)
}

func (o *oldToNewConfigerAdapter) Int64(key string) (int64, error) {
	return o.delegate.Int64(key)
}

func (o *oldToNewConfigerAdapter) Bool(key string) (bool, error) {
	return o.delegate.Bool(key)
}

func (o *oldToNewConfigerAdapter) Float(key string) (float64, error) {
	return o.delegate.Float(key)
}

func (o *oldToNewConfigerAdapter) DefaultString(key string, defaultVal string) string {
	return o.delegate.DefaultString(key, defaultVal)
}

func (o *oldToNewConfigerAdapter) DefaultStrings(key string, defaultVal []string) []string {
	return o.delegate.DefaultStrings(key, defaultVal)
}

func (o *oldToNewConfigerAdapter) DefaultInt(key string, defaultVal int) int {
	return o.delegate.DefaultInt(key, defaultVal)
}

func (o *oldToNewConfigerAdapter) DefaultInt64(key string, defaultVal int64) int64 {
	return o.delegate.DefaultInt64(key, defaultVal)
}

func (o *oldToNewConfigerAdapter) DefaultBool(key string, defaultVal bool) bool {
	return o.delegate.DefaultBool(key, defaultVal)
}

func (o *oldToNewConfigerAdapter) DefaultFloat(key string, defaultVal float64) float64 {
	return o.delegate.DefaultFloat(key, defaultVal)
}

func (o *oldToNewConfigerAdapter) DIY(key string) (interface{}, error) {
	return o.delegate.DIY(key)
}

func (o *oldToNewConfigerAdapter) GetSection(section string) (map[string]string, error) {
	return o.delegate.GetSection(section)
}

func (o *oldToNewConfigerAdapter) Unmarshaler(prefix string, obj interface{}, opt ...config.DecodeOption) error {
	return errors.New("unsupported operation, please use actual config.Configer")
}

func (o *oldToNewConfigerAdapter) Sub(key string) (config.Configer, error) {
	return nil, errors.New("unsupported operation, please use actual config.Configer")
}

func (o *oldToNewConfigerAdapter) OnChange(key string, fn func(value string)) {
	// do nothing
}

func (o *oldToNewConfigerAdapter) SaveConfigFile(filename string) error {
	return o.delegate.SaveConfigFile(filename)
}

type oldToNewConfigAdapter struct {
	delegate Config
}

func (o *oldToNewConfigAdapter) Parse(key string) (config.Configer, error) {
	old, err := o.delegate.Parse(key)
	if err != nil {
		return nil, err
	}
	return &oldToNewConfigerAdapter{delegate: old}, nil
}

func (o *oldToNewConfigAdapter) ParseData(data []byte) (config.Configer, error) {
	old, err := o.delegate.ParseData(data)
	if err != nil {
		return nil, err
	}
	return &oldToNewConfigerAdapter{delegate: old}, nil
}
