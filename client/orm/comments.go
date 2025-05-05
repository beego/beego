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
// Use the AddComment and ClearComments methods on an Ormer instance to manage comments.
type QueryComments struct {
	mu       sync.RWMutex // mu protects the comments slice
	comments []string     // comments stores the list of comments added
}

// NewQueryComments creates a new, empty QueryComments instance.
// The returned instance is safe for concurrent use.
// This is typically called internally when creating a new Ormer.
func NewQueryComments() *QueryComments {
	return &QueryComments{
		comments: make([]string, 0),
	}
}

// AddComment appends a new comment string to the list of comments for the current query context.
// If multiple comments are added, they will be joined by "; " within the final SQL comment block (e.g., /* comment1; comment2 */).
// An empty comment string is ignored.
// This method is safe for concurrent use.
func (qc *QueryComments) AddComment(comment string) {
	if comment == "" {
		return
	}
	qc.mu.Lock()
	defer qc.mu.Unlock() // Use defer for unlock
	qc.comments = append(qc.comments, comment)
}

// ClearComments removes all previously added comments for the current query context.
// This method is safe for concurrent use.
func (qc *QueryComments) ClearComments() {
	qc.mu.Lock()
	defer qc.mu.Unlock()          // Use defer for unlock
	qc.comments = qc.comments[:0] // Reset slice length to 0, keeping allocated capacity
}

// String formats the collected comments into a single SQL comment string (e.g., "/* comment1; comment2 */ ").
// If no comments have been added, it returns an empty string.
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
