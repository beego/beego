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

package redis_cluster

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProvider_SessionInit(t *testing.T) {
	savePath := `
{ "save_path": "my save path", "idle_timeout": "3s"}
`
	cp := &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "my save path", cp.SavePath)
	assert.Equal(t, 3*time.Second, cp.idleTimeout)
	assert.Equal(t, int64(12), cp.maxlifetime)
}
