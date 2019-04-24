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

// Package Kubernetes ConfigMapMap for config provider.
//
// depend on:
// "k8s.io/apimachinery/pkg/apis/meta/v1"
// "k8s.io/client-go/kubernetes"
// "k8s.io/client-go/rest"
//
// Usage:
//  import(
//    _ "github.com/astaxie/beego/config/configmap"
//      "github.com/astaxie/beego/config"
//  )
//
//  cnf, err := NewConfigMap("configMap", "myConfAppName")
//
//More docs http://beego.me/docs/module/md
package configmap

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/encoder/json"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/config/base"
)

func getClient() (*kubernetes.Clientset, error) {
	conf, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(conf)
}

func makeMap(kv map[string]string) map[string]interface{} {

	data := make(map[string]interface{})

	for k, v := range kv {
		data[k] = make(map[string]interface{})

		vals := strings.Split(v, "\n")

		mp := make(map[string]interface{})
		for _, h := range vals {
			m, n := split(string(h), "=")
			mp[m] = n
		}

		data[k] = mp
	}

	return data
}

func split(s string, sp string) (k string, v string) {
	kv := strings.Split(s, sp)
	return kv[0], kv[1]
}

// ConfigMap is a json config parser and implements ConfigMap interface.
type ConfigMap struct {
	client    *kubernetes.Clientset
	cerr      error
	name      string
	namespace string
	opts      *config.Option
}

func (k *ConfigMap) SetOption(option config.Option) {
	k.opts = &option
}

// Parse returns a ConfigMapContainer with parsed json config map.
func (k *ConfigMap) Parse() (config.Configer, error) {
	if k.cerr != nil {
		return nil, k.cerr
	}

	cmp, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(k.opts.ConfigName, v1.GetOptions{})

	if err != nil {
		return nil, err
	}

	data := makeMap(cmp.Data)

	b, err := k.opts.Encoder.Encode(data)

	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	return k.ParseData(b)
}

// ParseData json data
func (k *ConfigMap) ParseData(data []byte) (config.Configer, error) {
	x := &ConfigMapContainer{ConfigBaseContainer: base.ConfigBaseContainer{
		Data:          make(map[string]interface{}),
		SeparatorKeys: k.opts.SeparatorKeys,
	}}

	cnf := map[string]interface{}{}
	_ = k.opts.Encoder.Decode(data, &cnf)

	x.Data = config.ExpandValueEnvForMap(cnf)

	return x, nil
}

// ConfigMapContainer A ConfigMap represents the json configuration.
type ConfigMapContainer struct {
	base.ConfigBaseContainer
}

func NewConfigMap(option config.Option) *ConfigMap {
	client, err := getClient()

	return &ConfigMap{
		opts:      &option,
		name:      option.ConfigName,
		namespace: "default",
		cerr:      err,
		client:    client,
	}
}

func init() {
	config.Register("configMap", NewConfigMap(config.Option{
		ConfigName:    "beego",
		Context:       context.Background(),
		Encoder:       json.NewEncoder(),
		SeparatorKeys: "::",
	}))
}
