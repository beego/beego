// Copyright 2026 beego Author. All Rights Reserved.
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

// Package sqlsessiontest provides shared tests for SQL session pool lifetimes.
package sqlsessiontest

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/beego/beego/v2/server/web/session"
)

const operationSessionID = "operation-session"

// ProviderFactory initializes a provider with the supplied sqlmock DSN.
type ProviderFactory func(dsn string) (session.Provider, *sql.DB, error)

// Queries contains the SQL expected from a session provider.
type Queries struct {
	Select, Update, Destroy, All, GC string
	GCArgs                           []driver.Value
}

// MySQLQueries returns the expected MySQL session queries.
func MySQLQueries() Queries {
	return Queries{
		Select:  "select session_data from session where session_key=?",
		Update:  "UPDATE session set `session_data`=?, `session_expiry`=? where session_key=?",
		Destroy: "DELETE FROM session where session_key=?",
		All:     "SELECT count(*) as num from session",
		GC:      "DELETE from session where session_expiry < ?",
		GCArgs:  []driver.Value{sqlmock.AnyArg()},
	}
}

// PostgreSQLQueries returns the expected PostgreSQL session queries.
func PostgreSQLQueries() Queries {
	return Queries{
		Select:  "select session_data from session where session_key=$1",
		Update:  "UPDATE session set session_data=$1, session_expiry=$2 where session_key=$3",
		Destroy: "DELETE FROM session where session_key=$1",
		All:     "SELECT count(*) as num from session",
		GC:      "DELETE from session where EXTRACT(EPOCH FROM (current_timestamp - session_expiry)) > $1",
		GCArgs:  []driver.Value{int64(3600)},
	}
}

type fixture struct {
	provider session.Provider
	mock     sqlmock.Sqlmock
	queries  Queries
}

// RunStores checks that stores borrow the provider's pool in either release order.
func RunStores(t *testing.T, factory ProviderFactory, queries Queries) {
	t.Helper()
	for _, tc := range []struct {
		name           string
		firstIfPresent bool
	}{
		{"SessionRelease", false},
		{"SessionReleaseIfPresent", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			checkStores(ctx, t, factory, queries, tc.name, tc.firstIfPresent)
		})
	}
}

func checkStores(ctx context.Context, t *testing.T, factory ProviderFactory, queries Queries, operation string, firstIfPresent bool) {
	t.Helper()
	f := newFixture(t, factory, queries)
	first := readStore(ctx, t, f, "first-store")
	second := readStore(ctx, t, f, "second-store")
	releaseStore(ctx, f, first, "first-store", firstIfPresent)
	releaseStore(ctx, f, second, "second-store", !firstIfPresent)
	assertPoolReusable(ctx, t, f, operation, "after-release")
	assertExpectations(t, f.mock)
}

func readStore(ctx context.Context, t *testing.T, f *fixture, sid string) session.Store {
	t.Helper()
	expectSessionRow(f, sid)
	store, err := f.provider.SessionRead(ctx, sid)
	if err != nil {
		t.Fatalf("SessionRead(%q) returned an error: %v", sid, err)
	}
	return store
}

func releaseStore(ctx context.Context, f *fixture, store session.Store, sid string, ifPresent bool) {
	f.mock.ExpectExec(f.queries.Update).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sid).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if ifPresent {
		store.SessionReleaseIfPresent(ctx, nil)
	} else {
		store.SessionRelease(ctx, nil)
	}
}

// RunOperations checks that each provider operation leaves its pool usable.
func RunOperations(t *testing.T, factory ProviderFactory, queries Queries) {
	t.Helper()
	for _, tc := range []struct {
		name  string
		check func(context.Context, *testing.T, *fixture)
	}{
		{"SessionExist", checkExist},
		{"SessionDestroy", checkDestroy},
		{"SessionAll", checkAll},
		{"SessionGC", checkGC},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			f := newFixture(t, factory, queries)
			tc.check(ctx, t, f)
			assertPoolReusable(ctx, t, f, tc.name, "after-operation")
			assertExpectations(t, f.mock)
		})
	}
}

func checkExist(ctx context.Context, t *testing.T, f *fixture) {
	t.Helper()
	expectSessionRow(f, operationSessionID)
	exists, err := f.provider.SessionExist(ctx, operationSessionID)
	if err != nil || !exists {
		t.Fatalf("SessionExist = (%v, %v); want (true, nil)", exists, err)
	}
}

func checkDestroy(ctx context.Context, t *testing.T, f *fixture) {
	t.Helper()
	f.mock.ExpectExec(f.queries.Destroy).WithArgs(operationSessionID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := f.provider.SessionDestroy(ctx, operationSessionID); err != nil {
		t.Fatalf("SessionDestroy returned an error: %v", err)
	}
}

func checkAll(ctx context.Context, t *testing.T, f *fixture) {
	t.Helper()
	f.mock.ExpectQuery(f.queries.All).
		WillReturnRows(sqlmock.NewRows([]string{"num"}).AddRow(2)).RowsWillBeClosed()
	if got := f.provider.SessionAll(ctx); got != 2 {
		t.Fatalf("SessionAll = %d; want 2", got)
	}
}

func checkGC(ctx context.Context, t *testing.T, f *fixture) {
	t.Helper()
	f.mock.ExpectExec(f.queries.GC).WithArgs(f.queries.GCArgs...).
		WillReturnResult(sqlmock.NewResult(0, 1))
	f.provider.SessionGC(ctx)
}

func expectSessionRow(f *fixture, sid string) {
	f.mock.ExpectQuery(f.queries.Select).WithArgs(sid).
		WillReturnRows(sqlmock.NewRows([]string{"session_data"}).AddRow([]byte{})).RowsWillBeClosed()
}

func assertPoolReusable(ctx context.Context, t *testing.T, f *fixture, operation, sid string) {
	t.Helper()
	expectSessionRow(f, sid)
	exists, err := f.provider.SessionExist(ctx, sid)
	if err != nil {
		t.Fatalf("SessionExist after %s returned an error: %v", operation, err)
	}
	if !exists {
		t.Fatalf("SessionExist after %s = false; want true", operation)
	}
}

func assertExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("session operation did not complete expected SQL: %v", err)
	}
}

func newFixture(t *testing.T, factory ProviderFactory, queries Queries) *fixture {
	t.Helper()
	dsn := t.Name()
	initialDB, mock, err := sqlmock.NewWithDSN(dsn, sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("creating SQL mock: %v", err)
	}
	t.Cleanup(func() { closeMockDB(t, mock, initialDB, "initial mock") })
	provider, pool, err := factory(dsn)
	if err != nil {
		t.Fatalf("initializing SQL mock provider: %v", err)
	}
	t.Cleanup(func() { closeMockDB(t, mock, pool, "provider") })
	if err := pool.Ping(); err != nil {
		t.Fatalf("opening provider database: %v", err)
	}
	return &fixture{provider: provider, mock: mock, queries: queries}
}

func closeMockDB(t *testing.T, mock sqlmock.Sqlmock, db *sql.DB, name string) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations before closing %s database: %v", name, err)
	}
	mock.ExpectClose()
	if err := db.Close(); err != nil {
		t.Errorf("closing %s database: %v", name, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations after closing %s database: %v", name, err)
	}
}
