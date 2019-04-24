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

// Package Kubernetes ConfigConsulMap for config provider.
//
// depend on:
// "github.com/hashicorp/consul/api"
//
// Usage:
//  import(
//    _ "github.com/astaxie/beego/config/consul"
//      "github.com/astaxie/beego/config"
//  )
//
//  cnf, err := NewConfigConsul("consul", "myConfAppName")
//
//More docs http://beego.me/docs/module/md
package configmap

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/config/base"
	"github.com/astaxie/beego/encoder"
	"github.com/hashicorp/consul/api"
)

const DefaultPrefix = "/beego/config/"

type addressKey struct{}
type prefixKey struct{}
type stripPrefixKey struct{}
type dcKey struct{}
type tokenKey struct{}

// ConfigConsul is a json config parser and implements ConfigConsul interface.
type ConfigConsul struct {
	prefix      string
	stripPrefix string
	addr        string
	client      *api.Client
	opts        *config.Option
}

func (k *ConfigConsul) SetOption(option config.Option) {
	k.opts = &option

	conf, prefix, sp, client := getClient(option)

	k.prefix = prefix
	k.client = client
	k.stripPrefix = sp
	k.addr = conf.Address
}

// Parse returns a ConfigConsulContainer with parsed json config map.
func (k *ConfigConsul) Parse() (config.Configer, error) {
	kv, _, err := k.client.KV().List(k.prefix, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil || len(kv) == 0 {
		return nil, fmt.Errorf("source not found: %s", k.prefix)
	}

	data, err := makeMap(k.opts.Encoder, kv, k.stripPrefix)

	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	b, err := k.opts.Encoder.Encode(data)

	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	return k.ParseData(b)
}

// ParseData json data
func (k *ConfigConsul) ParseData(data []byte) (config.Configer, error) {
	x := &ConfigConsulContainer{ConfigBaseContainer: base.ConfigBaseContainer{
		Data:          make(map[string]interface{}),
		SeparatorKeys: k.opts.SeparatorKeys,
	}}

	cnf := map[string]interface{}{}
	_ = k.opts.Encoder.Decode(data, &cnf)

	x.Data = config.ExpandValueEnvForMap(cnf)

	return x, nil
}

// ConfigConsulContainer A ConfigConsul represents the json configuration.
type ConfigConsulContainer struct {
	base.ConfigBaseContainer
}

func NewConfigConsul(option config.Option) *ConfigConsul {
	conf, prefix, sp, client := getClient(option)

	return &ConfigConsul{
		prefix:      prefix,
		stripPrefix: sp,
		opts:        &option,
		addr:        conf.Address,
		client:      client,
	}
}

func init() {
	config.Register("consul", NewConfigConsul(config.Option{
		ConfigName:    DefaultPrefix,
		SeparatorKeys: "::",
		Context:       context.Background(),
	}))
}

func makeMap(e encoder.Encoder, kv api.KVPairs, stripPrefix string) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// consul guarantees lexicographic order, so no need to sort
	for _, v := range kv {
		pathString := strings.TrimPrefix(strings.TrimPrefix(v.Key, stripPrefix), "/")
		var val map[string]interface{}

		// ensure a valid value is stored at this location
		if len(v.Value) > 0 {
			if err := e.Decode(v.Value, &val); err != nil {
				return nil, fmt.Errorf("faild decode value. path: %s, error: %s", pathString, err)
			}
		}

		// set target at the root
		target := data

		// then descend to the target location, creating as we go, if need be
		if pathString != "" {
			path := strings.Split(pathString, "/")
			// find (or create) the location we want to put this value at
			for _, dir := range path {
				if _, ok := target[dir]; !ok {
					target[dir] = make(map[string]interface{})
				}
				target = target[dir].(map[string]interface{})
			}

		}

		// copy over the keys from the value
		for k := range val {
			target[k] = val[k]
		}
	}

	return data, nil
}

func getClient(option config.Option) (*api.Config, string, string, *api.Client) {
	// use default config
	config := api.DefaultConfig()
	// check if there are any addrs
	a, ok := option.Context.Value(addressKey{}).(string)
	if ok {
		addr, port, err := net.SplitHostPort(a)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "8500"
			addr = a
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		} else if err == nil {
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		}
	}
	dc, ok := option.Context.Value(dcKey{}).(string)
	if ok {
		config.Datacenter = dc
	}
	token, ok := option.Context.Value(tokenKey{}).(string)
	if ok {
		config.Token = token
	}
	prefix := option.ConfigName
	sp := ""
	f, ok := option.Context.Value(prefixKey{}).(string)
	if ok {
		prefix = f
	}
	if b, ok := option.Context.Value(stripPrefixKey{}).(bool); ok && b {
		sp = prefix
	}
	// create the client
	client, _ := api.NewClient(config)
	return config, prefix, sp, client
}
