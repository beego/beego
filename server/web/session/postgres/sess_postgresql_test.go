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

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/beego/beego/v2/server/web/session"
	"github.com/beego/beego/v2/server/web/session/internal/sqlsessiontest"
)

func TestProviderUsesPostgreSQLDriver(t *testing.T) {
	provider := new(Provider)
	if err := provider.SessionInit(context.Background(), 3600,
		"host=/nonexistent/beego-session-test user=test dbname=test sslmode=disable"); err != nil {
		t.Fatalf("SessionInit returned an error: %v", err)
	}
	pool := provider.connectInit()
	t.Cleanup(func() { _ = pool.Close() })

	driverType := fmt.Sprintf("%T", pool.Driver())
	if driverType != "*pq.Driver" {
		t.Fatalf("database driver = %q; want %q", driverType, "*pq.Driver")
	}
	provider.SessionGC(context.Background())
}

func TestSessionStoresKeepPoolOpen(t *testing.T) {
	sqlsessiontest.RunStores(t, newMockProvider, sqlsessiontest.PostgreSQLQueries())
}

func TestProviderOperationsKeepPoolOpen(t *testing.T) {
	sqlsessiontest.RunOperations(t, newMockProvider, sqlsessiontest.PostgreSQLQueries())
}

func newMockProvider(dsn string) (session.Provider, *sql.DB, error) {
	provider := &Provider{}
	if err := provider.state.Init("sqlmock", "postgresql", 3600, dsn); err != nil {
		return nil, nil, err
	}
	return provider, provider.connectInit(), nil
}
