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
	"strings"
	"sync"
)

// QueryComments stores SQL query comments and provides thread-safe access.
// Comments will be included in generated SQL queries for improved debugging and tracing.
type QueryComments struct {
	mu       sync.RWMutex
	comments []string
}

// NewQueryComments creates a new QueryComments instance.
// The returned instance is safe for concurrent use.
func NewQueryComments() *QueryComments {
	return &QueryComments{
		comments: make([]string, 0),
	}
}

// AddComment adds a comment to the query comments.
// Multiple comments will be joined with semicolons in the final SQL query.
// This method is safe for concurrent use.
func (qc *QueryComments) AddComment(comment string) {
	if comment == "" {
		return
	}
	qc.mu.Lock()
	qc.comments = append(qc.comments, comment)
	qc.mu.Unlock()
}

// ClearComments removes all comments.
// This method is safe for concurrent use.
func (qc *QueryComments) ClearComments() {
	qc.mu.Lock()
	qc.comments = qc.comments[:0]
	qc.mu.Unlock()
}

// String returns all comments formatted as a SQL comment string.
// Multiple comments are joined with semicolons.
// Returns an empty string if there are no comments.
// This method is safe for concurrent use.
func (qc *QueryComments) String() string {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	if len(qc.comments) == 0 {
		return ""
	}
	return "/* " + strings.Join(qc.comments, "; ") + " */ "
}

// DefaultQueryComments is the default QueryComments instance used by the ORM.
var DefaultQueryComments = NewQueryComments()

// AddQueryComment adds a comment that will be included in subsequent queries.
// This is a convenience function that uses DefaultQueryComments.
func AddQueryComment(comment string) {
	DefaultQueryComments.AddComment(comment)
}

// ClearQueryComments removes all query comments.
// This is a convenience function that uses DefaultQueryComments.
func ClearQueryComments() {
	DefaultQueryComments.ClearComments()
}

var _ QueryCommenter = (*QueryComments)(nil)
