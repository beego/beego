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

import "strings"

// QueryComments stores the comments to be added to queries
type QueryComments struct {
	comments []string
}

// NewQueryComments creates a new QueryComments instance
func NewQueryComments() *QueryComments {
	return &QueryComments{
		comments: make([]string, 0),
	}
}

// Add appends a new comment
func (qc *QueryComments) Add(comment string) {
	qc.comments = append(qc.comments, comment)
}

// Clear removes all comments
func (qc *QueryComments) Clear() {
	qc.comments = qc.comments[:0]
}

// String returns all comments formatted as a SQL comment string
func (qc *QueryComments) String() string {
	if len(qc.comments) == 0 {
		return ""
	}
	return "/* " + strings.Join(qc.comments, "; ") + " */ "
}

var (
	// DefaultQueryComments is the default QueryComments instance used by the ORM
	DefaultQueryComments = NewQueryComments()
)

// AddQueryComment adds a comment that will be included in subsequent queries
func AddQueryComment(comment string) {
	DefaultQueryComments.Add(comment)
}

// ClearQueryComments removes all query comments
func ClearQueryComments() {
	DefaultQueryComments.Clear()
}