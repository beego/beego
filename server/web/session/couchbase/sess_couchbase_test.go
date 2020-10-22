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

package couchbase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_SessionInit(t *testing.T) {
	// using old style
	savePath := `http://host:port/,Pool,Bucket`
	cp := &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "http://host:port/", cp.SavePath)
	assert.Equal(t, "Pool", cp.Pool)
	assert.Equal(t, "Bucket", cp.Bucket)
	assert.Equal(t, int64(12), cp.maxlifetime)

	savePath = `
{ "save_path": "my save path", "pool": "mypool", "bucket": "mybucket"}
`
	cp = &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "my save path", cp.SavePath)
	assert.Equal(t, "mypool", cp.Pool)
	assert.Equal(t, "mybucket", cp.Bucket)
	assert.Equal(t, int64(12), cp.maxlifetime)
}
