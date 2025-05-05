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

// Package orm provides SQL query comment support to enhance query tracing and debugging.
//
// Query comments are SQL comments (/* comment */) that are automatically prepended to
// SQL queries. They can be used for tracing, debugging, and monitoring purposes.
// Comments are thread-safe and automatically cleared after each query execution.
//
// Example usage:
//
//	o := orm.NewOrm()
//	o.AddQueryComment("trace_id:123")
//	o.AddQueryComment("user_id:456")
//
//	// Generated SQL: /* trace_id:123; user_id:456 */ SELECT * FROM users
//	var users []User
//	o.QueryTable("users").All(&users)
//
//	// Comments are automatically cleared after query execution
//	// Next query will have no comments unless explicitly added
//
//	// Comments can be manually cleared if needed
//	o.ClearQueryComments()
//
// For a complete example, see examples/query_comments.go.
package orm

import (
	"strings"
	"sync"
)

// QueryComments stores SQL query comments and provides thread-safe access.
// Each comment is wrapped in /* */ SQL comment syntax and multiple comments
// are joined with semicolons.
//
// Comments are automatically prepended to generated SQL queries for improved
// debugging, tracing, and monitoring. They are cleared after each query execution
// to prevent comments from leaking between unrelated operations.
//
// Thread-safety is guaranteed by using a sync.RWMutex to protect access to
// the comments slice.
//
// Example:
//
//	qc := orm.NewQueryComments()
//	qc.AddComment("trace_id:123")
//	qc.AddComment("user_id:456")
//	fmt.Println(qc.String()) // Outputs: /* trace_id:123; user_id:456 */
//
//	// Comments cleared for next operation
//	qc.ClearComments()
//	fmt.Println(qc.String()) // Outputs: ""
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
// If multiple comments are added, they will be joined by "; " within the final SQL comment block.
// Empty comment strings are ignored to prevent empty comments in the SQL.
//
// The method is safe for concurrent use via sync.RWMutex protection.
//
// Example:
//
//	qc := orm.NewQueryComments()
//	qc.AddComment("trace_id:abc123")
//	qc.AddComment("user_id:456")
//	// Result: /* trace_id:abc123; user_id:456 */ SELECT ...
//
//	// Empty comments are ignored automatically
//	qc.AddComment("")
//	// Result unchanged: /* trace_id:abc123; user_id:456 */ SELECT ...
func (qc *QueryComments) AddComment(comment string) {
	if comment == "" {
		return
	}
	qc.mu.Lock()
	defer qc.mu.Unlock() // Use defer for unlock
	qc.comments = append(qc.comments, comment)
}

// ClearComments removes all previously added comments for the current query context.
// This provides a way to explicitly clear comments when needed, though comments
// are automatically cleared after query execution.
//
// The method is safe for concurrent use via sync.RWMutex protection.
//
// Example:
//
//	qc := orm.NewQueryComments()
//	qc.AddComment("trace_id:123")
//
//	// Before clear: /* trace_id:123 */ SELECT ...
//	qc.ClearComments()
//	// After clear: SELECT ...
//
//	// Adding new comments starts fresh
//	qc.AddComment("new_trace:789")
//	// In SQL: /* new_trace:789 */ SELECT ...
func (qc *QueryComments) ClearComments() {
	qc.mu.Lock()
	defer qc.mu.Unlock()          // Use defer for unlock
	qc.comments = qc.comments[:0] // Reset slice length to 0, keeping allocated capacity
}

// String formats the collected comments into a single SQL comment string surrounded
// by /* */ and adding spaces for readability. If no comments have been added,
// it returns an empty string.
//
// Multiple comments are joined with semicolons. Each comment is preserved exactly
// as added (no escaping or modification of comment text).
//
// The method is safe for concurrent use via sync.RWMutex protection.
//
// Example:
//
//	qc := orm.NewQueryComments()
//
//	// No comments
//	fmt.Println(qc.String()) // Outputs: ""
//
//	// Single comment
//	qc.AddComment("trace_id:123")
//	fmt.Println(qc.String()) // Outputs: "/* trace_id:123 */ "
//
//	// Multiple comments
//	qc.AddComment("user_id:456")
//	fmt.Println(qc.String()) // Outputs: "/* trace_id:123; user_id:456 */ "
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

// ClearQueryComments removes all global query comments.
// This is a convenience function that uses DefaultQueryComments.
//
// Example workflow:
//
//	orm.AddQueryComment("debug=true")
//	// ... run some queries ...
//	orm.ClearQueryComments() // Clean up when done debugging
func ClearQueryComments() {
	DefaultQueryComments.ClearComments()
}

var _ QueryCommenter = (*QueryComments)(nil)
