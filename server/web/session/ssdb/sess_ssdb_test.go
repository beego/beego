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

package ssdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_SessionInit(t *testing.T) {
	// using old style
	savePath := `localhost:8080`
	cp := &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "localhost", cp.Host)
	assert.Equal(t, 8080, cp.Port)
	assert.Equal(t, int64(12), cp.maxLifetime)

	savePath = `
{ "host": "localhost", "port": 8080}
`
	cp = &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "localhost", cp.Host)
	assert.Equal(t, 8080, cp.Port)
	assert.Equal(t, int64(12), cp.maxLifetime)
}
