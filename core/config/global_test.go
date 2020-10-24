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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalInstance(t *testing.T) {
	cfgStr := `
appname = beeapi
httpport = 8080
mysqlport = 3600
PI = 3.1415926
runmode = "dev"
autorender = false
copyrequestbody = true
session= on
cookieon= off
newreg = OFF
needlogin = ON
enableSession = Y
enableCookie = N
developer="tom;jerry"
flag = 1
path1 = ${GOPATH}
path2 = ${GOPATH||/home/go}
[demo]
key1="asta"
key2 = "xie"
CaseInsensitive = true
peers = one;two;three
password = ${GOPATH}
`
	path := os.TempDir() + string(os.PathSeparator) + "test_global_instance.ini"
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(cfgStr)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(path)

	err = InitGlobalInstance("ini", path)
	assert.Nil(t, err)

	val, err := String("appname")
	assert.Nil(t, err)
	assert.Equal(t, "beeapi", val)

	val = DefaultString("appname__", "404")
	assert.Equal(t, "404", val)

	vi, err := Int("httpport")
	assert.Nil(t, err)
	assert.Equal(t, 8080, vi)
	vi = DefaultInt("httpport__", 404)
	assert.Equal(t, 404, vi)

	vi64, err := Int64("mysqlport")
	assert.Nil(t, err)
	assert.Equal(t, int64(3600), vi64)
	vi64 = DefaultInt64("mysqlport__", 404)
	assert.Equal(t, int64(404), vi64)

	vf, err := Float("PI")
	assert.Nil(t, err)
	assert.Equal(t, 3.1415926, vf)
	vf = DefaultFloat("PI__", 4.04)
	assert.Equal(t, 4.04, vf)

	vb, err := Bool("copyrequestbody")
	assert.Nil(t, err)
	assert.True(t, vb)

	vb = DefaultBool("copyrequestbody__", false)
	assert.False(t, vb)

	vss := DefaultStrings("developer__", []string{"tom", ""})
	assert.Equal(t, []string{"tom", ""}, vss)

	vss, err = Strings("developer")
	assert.Nil(t, err)
	assert.Equal(t, []string{"tom", "jerry"}, vss)
}
