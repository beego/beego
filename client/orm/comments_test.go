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
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryComments(t *testing.T) {
	t.Run("Basic QueryComments Operations", func(t *testing.T) {
		qc := NewQueryComments()

		// Test empty comments
		assert.Equal(t, "", qc.String())

		// Test single comment
		qc.AddComment("test comment")
		assert.Equal(t, "/* test comment */ ", qc.String())

		// Test multiple comments
		qc.AddComment("another comment")
		assert.Equal(t, "/* test comment; another comment */ ", qc.String())

		// Test clear
		qc.ClearComments()
		assert.Equal(t, "", qc.String())
	})

	t.Run("Empty Comment Handling", func(t *testing.T) {
		qc := NewQueryComments()
		qc.AddComment("")
		assert.Equal(t, "", qc.String(), "Empty comment should be ignored")

		qc.AddComment("valid comment")
		qc.AddComment("")
		assert.Equal(t, "/* valid comment */ ", qc.String(), "Empty comment should be ignored among valid ones")
	})

	t.Run("Thread Safety", func(t *testing.T) {
		qc := NewQueryComments()
		done := make(chan bool)

		// Test concurrent access
		go func() {
			for i := 0; i < 100; i++ {
				qc.AddComment("comment1")
				qc.ClearComments()
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				qc.AddComment("comment2")
				qc.String()
			}
			done <- true
		}()

		<-done
		<-done
	})
}

func TestPrependComments(t *testing.T) {
	t.Run("Basic Comment Prepending", func(t *testing.T) {
		query := "SELECT * FROM users"
		qc := NewQueryComments()
		qc.AddComment("trace_id:123")

		mockDB := &mockDBQuerier{comments: qc}
		result := prependComments(mockDB, query)
		assert.Equal(t, "/* trace_id:123 */ SELECT * FROM users", result)
	})

	t.Run("Nil Handling", func(t *testing.T) {
		query := "SELECT * FROM users"
		result := prependComments(nil, query)
		assert.Equal(t, query, result, "Should handle nil querier gracefully")
	})

	t.Run("Empty Comments", func(t *testing.T) {
		query := "SELECT * FROM users"
		mockDB := &mockDBQuerier{comments: NewQueryComments()}
		result := prependComments(mockDB, query)
		assert.Equal(t, query, result, "Should handle empty comments gracefully")
	})

	t.Run("Multiple Comments", func(t *testing.T) {
		query := "SELECT * FROM users"
		comments := NewQueryComments()
		comments.AddComment("trace_id:123")
		comments.AddComment("user_id:456")
		mockDB := &mockDBQuerier{comments: comments}

		result := prependComments(mockDB, query)
		expected := "/* trace_id:123; user_id:456 */ SELECT * FROM users"
		assert.Equal(t, expected, result)
	})
}

// mockDBQuerier implements dbQuerier interface for testing
type mockDBQuerier struct {
	comments *QueryComments
}

func (m *mockDBQuerier) GetQueryComments() *QueryComments {
	return m.comments
}

func (m *mockDBQuerier) Prepare(query string) (*sql.Stmt, error) {
	return new(sql.Stmt), nil
}

func (m *mockDBQuerier) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return new(sql.Stmt), nil
}

func (m *mockDBQuerier) Exec(query string, args ...interface{}) (sql.Result, error) {
	return &mockResult{}, nil // Use composite literal instead of new
}

func (m *mockDBQuerier) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return &mockResult{}, nil // Use composite literal instead of new
}

func (m *mockDBQuerier) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return new(sql.Rows), nil
}

func (m *mockDBQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return new(sql.Rows), nil
}

func (m *mockDBQuerier) QueryRow(query string, args ...interface{}) *sql.Row {
	return new(sql.Row)
}

func (m *mockDBQuerier) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return new(sql.Row)
}

// mockResult implements sql.Result for testing
type mockResult struct{}

func (r *mockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r *mockResult) RowsAffected() (int64, error) {
	return 0, nil
}
