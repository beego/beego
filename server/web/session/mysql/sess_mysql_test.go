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

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/beego/beego/v2/server/web/session"
	"github.com/beego/beego/v2/server/web/session/internal/sqlsessiontest"
)

func TestProviderUsesMySQLDriver(t *testing.T) {
	provider := &Provider{}
	ctx := context.Background()
	if err := provider.SessionInit(ctx, 3600, "user:password@unix(/nonexistent/beego-session-test.sock)/database"); err != nil {
		t.Fatalf("SessionInit returned an error: %v", err)
	}
	db := provider.connectInit()
	t.Cleanup(func() {
		_ = db.Close()
	})

	if got, want := fmt.Sprintf("%T", db.Driver()), "*mysql.MySQLDriver"; got != want {
		t.Errorf("database driver = %q; want %q", got, want)
	}
	provider.SessionGC(ctx)
}

func TestSessionStoresKeepPoolOpen(t *testing.T) {
	sqlsessiontest.RunStores(t, newMockProvider, sqlsessiontest.MySQLQueries())
}

func TestProviderOperationsKeepPoolOpen(t *testing.T) {
	sqlsessiontest.RunOperations(t, newMockProvider, sqlsessiontest.MySQLQueries())
}

func newMockProvider(dsn string) (session.Provider, *sql.DB, error) {
	provider := &Provider{}
	if err := provider.state.Init("sqlmock", "mysql", 3600, dsn); err != nil {
		return nil, nil, err
	}
	return provider, provider.connectInit(), nil
}
