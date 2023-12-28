// Copyright 2023 beego. All Rights Reserved.
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

package qb

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestResult_RowsAffected(t *testing.T) {
	testCases := []struct {
		name         string
		res          Result
		wantAffected int64
		wantErr      error
	}{
		{
			name:    "err",
			wantErr: errors.New("exec err"),
			res:     Result{err: errors.New("exec err")},
		},
		{
			name:    "unknown error",
			wantErr: errors.New("unknown error"),
			res:     Result{res: sqlmock.NewErrorResult(errors.New("unknown error"))},
		},
		{
			name:         "no err",
			wantAffected: int64(234),
			res:          Result{res: sqlmock.NewResult(123, 234)},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			affected, err := tc.res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantAffected, affected)
		})
	}
}

func TestResult_LastInsertId(t *testing.T) {
	testCases := []struct {
		name       string
		res        Result
		wantLastId int64
		wantErr    error
	}{
		{
			name:    "err",
			wantErr: errors.New("exec err"),
			res:     Result{err: errors.New("exec err")},
		},
		{
			name:    "res err",
			wantErr: errors.New("exec err"),
			res:     Result{res: sqlmock.NewErrorResult(errors.New("exec err"))},
		},
		{
			name:       "no err",
			wantLastId: int64(123),
			res:        Result{res: sqlmock.NewResult(123, 234)},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := tc.res.LastInsertId()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLastId, id)
		})
	}
}
