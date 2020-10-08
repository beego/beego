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

package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpServerWithCfg(t *testing.T) {
	// we should make sure that update server's config won't change
	BConfig.AppName = "Before"
	svr := NewHttpServerWithCfg(*BConfig)
	svr.Cfg.AppName = "hello"
	assert.NotEqual(t, "hello", BConfig.AppName)
	assert.Equal(t, "Before", BConfig.AppName)

}
