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

package etcd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdConfigerProvider_Parse(t *testing.T) {
	provider := &EtcdConfigerProvider{}
	cfger, err := provider.Parse(readEtcdConfig())
	assert.Nil(t, err)
	assert.NotNil(t, cfger)
}

func TestEtcdConfiger(t *testing.T) {

	provider := &EtcdConfigerProvider{}
	cfger, _ := provider.Parse(readEtcdConfig())

	subCfger, err := cfger.Sub("sub.")
	assert.Nil(t, err)
	assert.NotNil(t, subCfger)

	subSubCfger, err := subCfger.Sub("sub.")
	assert.NotNil(t, subSubCfger)
	assert.Nil(t, err)

	str, err := subSubCfger.String("key1")
	assert.Nil(t, err)
	assert.Equal(t, "sub.sub.key", str)

	// we cannot test it
	subSubCfger.OnChange("watch", func(value string) {
		// do nothing
	})

	defStr := cfger.DefaultString("not_exit", "default value")
	assert.Equal(t, "default value", defStr)

	defInt64 := cfger.DefaultInt64("not_exit", -1)
	assert.Equal(t, int64(-1), defInt64)

	defInt := cfger.DefaultInt("not_exit", -2)
	assert.Equal(t, -2, defInt)

	defFlt := cfger.DefaultFloat("not_exit", 12.3)
	assert.Equal(t, 12.3, defFlt)

	defBl := cfger.DefaultBool("not_exit", true)
	assert.True(t, defBl)

	defStrs := cfger.DefaultStrings("not_exit", []string{"hello"})
	assert.Equal(t, []string{"hello"}, defStrs)

	fl, err := cfger.Float("current.float")
	assert.Nil(t, err)
	assert.Equal(t, 1.23, fl)

	bl, err := cfger.Bool("current.bool")
	assert.Nil(t, err)
	assert.True(t, bl)

	it, err := cfger.Int("current.int")
	assert.Nil(t, err)
	assert.Equal(t, 11, it)

	str, err = cfger.String("current.string")
	assert.Nil(t, err)
	assert.Equal(t, "hello", str)

	tn := &TestEntity{}
	err = cfger.Unmarshaler("current.serialize.", tn)
	assert.Nil(t, err)
	assert.Equal(t, "test", tn.Name)
}

type TestEntity struct {
	Name string    `yaml:"name"`
	Sub  SubEntity `yaml:"sub"`
}

type SubEntity struct {
	SubName string `yaml:"subName"`
}

func readEtcdConfig() string {
	addr := os.Getenv("ETCD_ADDR")
	if addr == "" {
		addr = "localhost:2379"
	}

	obj := clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 3 * time.Second,
	}
	cfg, _ := json.Marshal(obj)
	return string(cfg)
}
