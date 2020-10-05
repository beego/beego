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

package bean

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeTypeAdapter_DefaultValue(t *testing.T) {
	typeAdapter := &TimeTypeAdapter{Layout: "2006-01-02 15:04:05"}
	tm, err := typeAdapter.DefaultValue(context.Background(), "2018-02-03 12:34:11")
	assert.Nil(t, err)
	assert.NotNil(t, tm)
}
