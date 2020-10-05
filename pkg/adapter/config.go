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

package adapter

import (
	context2 "context"

	"github.com/astaxie/beego/pkg/adapter/session"
	newCfg "github.com/astaxie/beego/pkg/core/config"
	"github.com/astaxie/beego/pkg/server/web"
)

// Config is the main struct for BConfig
type Config web.Config

// Listen holds for http and https related config
type Listen web.Listen

// WebConfig holds web related config
type WebConfig web.WebConfig

// SessionConfig holds session related config
type SessionConfig web.SessionConfig

// LogConfig holds Log related config
type LogConfig web.LogConfig

var (
	// BConfig is the default config for Application
	BConfig *Config
	// AppConfig is the instance of Config, store the config information from file
	AppConfig *beegoAppConfig
	// AppPath is the absolute path to the app
	AppPath string
	// GlobalSessions is the instance for the session manager
	GlobalSessions *session.Manager

	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "ini"
	// WorkPath is the absolute path to project root directory
	WorkPath string
)

func init() {
	BConfig = (*Config)(web.BConfig)
	AppPath = web.AppPath

	WorkPath = web.WorkPath

	AppConfig = &beegoAppConfig{innerConfig: (newCfg.Configer)(web.AppConfig)}
}

// LoadAppConfig allow developer to apply a config file
func LoadAppConfig(adapterName, configPath string) error {
	return web.LoadAppConfig(adapterName, configPath)
}

type beegoAppConfig struct {
	innerConfig newCfg.Configer
}

func (b *beegoAppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(context2.Background(), BConfig.RunMode+"::"+key, val); err != nil {
		return b.innerConfig.Set(context2.Background(), key, val)
	}
	return nil
}

func (b *beegoAppConfig) String(key string) string {
	if v, err := b.innerConfig.String(context2.Background(), BConfig.RunMode+"::"+key); v != "" && err != nil {
		return v
	}
	res, _ := b.innerConfig.String(context2.Background(), key)
	return res
}

func (b *beegoAppConfig) Strings(key string) []string {
	if v, err := b.innerConfig.Strings(context2.Background(), BConfig.RunMode+"::"+key); len(v) > 0 && err != nil {
		return v
	}
	res, _ := b.innerConfig.Strings(context2.Background(), key)
	return res
}

func (b *beegoAppConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(context2.Background(), BConfig.RunMode+"::"+key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(context2.Background(), key)
}

func (b *beegoAppConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(context2.Background(), BConfig.RunMode+"::"+key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(context2.Background(), key)
}

func (b *beegoAppConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(context2.Background(), BConfig.RunMode+"::"+key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(context2.Background(), key)
}

func (b *beegoAppConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(context2.Background(), BConfig.RunMode+"::"+key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(context2.Background(), key)
}

func (b *beegoAppConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *beegoAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(context2.Background(), key)
}

func (b *beegoAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(context2.Background(), section)
}

func (b *beegoAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(context2.Background(), filename)
}
