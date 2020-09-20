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
	"context"
	"errors"
	"strconv"
	"strings"
)

type fakeConfigContainer struct {
	BaseConfiger
	data map[string]string
}

func (c *fakeConfigContainer) getData(key string) string {
	return c.data[strings.ToLower(key)]
}

func (c *fakeConfigContainer) Set(ctx context.Context, key, val string) error {
	c.data[strings.ToLower(key)] = val
	return nil
}

func (c *fakeConfigContainer) Int(ctx context.Context, key string) (int, error) {
	return strconv.Atoi(c.getData(key))
}

func (c *fakeConfigContainer) DefaultInt(ctx context.Context, key string, defaultVal int) int {
	v, err := c.Int(ctx, key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) Int64(ctx context.Context, key string) (int64, error) {
	return strconv.ParseInt(c.getData(key), 10, 64)
}

func (c *fakeConfigContainer) DefaultInt64(ctx context.Context, key string, defaultVal int64) int64 {
	v, err := c.Int64(ctx, key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) Bool(ctx context.Context, key string) (bool, error) {
	return ParseBool(c.getData(key))
}

func (c *fakeConfigContainer) DefaultBool(ctx context.Context, key string, defaultVal bool) bool {
	v, err := c.Bool(ctx, key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) Float(ctx context.Context, key string) (float64, error) {
	return strconv.ParseFloat(c.getData(key), 64)
}

func (c *fakeConfigContainer) DefaultFloat(ctx context.Context, key string, defaultVal float64) float64 {
	v, err := c.Float(ctx, key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *fakeConfigContainer) DIY(ctx context.Context, key string) (interface{}, error) {
	if v, ok := c.data[strings.ToLower(key)]; ok {
		return v, nil
	}
	return nil, errors.New("key not find")
}

func (c *fakeConfigContainer) GetSection(ctx context.Context, section string) (map[string]string, error) {
	return nil, errors.New("not implement in the fakeConfigContainer")
}

func (c *fakeConfigContainer) SaveConfigFile(ctx context.Context, filename string) error {
	return errors.New("not implement in the fakeConfigContainer")
}

var _ Configer = new(fakeConfigContainer)

// NewFakeConfig return a fake Configer
func NewFakeConfig() Configer {
	res := &fakeConfigContainer{
		data: make(map[string]string),
	}
	res.BaseConfiger = NewBaseConfiger(func(ctx context.Context, key string) (string, error) {
		return res.getData(key), nil
	})
	return res
}
