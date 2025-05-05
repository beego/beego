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

package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryComments(t *testing.T) {
	qc := NewQueryComments()

	// Test empty comments
	assert.Equal(t, "", qc.String())

	// Test single comment
	qc.AddComment("test comment") // Renamed from Add
	assert.Equal(t, "/* test comment */ ", qc.String())

	// Test multiple comments
	qc.AddComment("another comment") // Renamed from Add
	assert.Equal(t, "/* test comment; another comment */ ", qc.String())

	// Test clear
	qc.ClearComments() // Renamed from Clear
	assert.Equal(t, "", qc.String())
}

// TODO: Add a test using a real Ormer instance to verify comments appear in logs/SQL
