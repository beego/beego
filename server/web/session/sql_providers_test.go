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

package session_test

import (
	"context"
	"sync"
	"testing"

	"github.com/beego/beego/v2/server/web/session/mysql"
	"github.com/beego/beego/v2/server/web/session/postgres"
)

type sqlProvider interface {
	SessionInit(context.Context, int64, string) error
	SessionGC(context.Context)
}

func TestSQLProvidersSharePoolConfigurationRules(t *testing.T) {
	tests := []struct {
		name                string
		savePath            string
		conflictingSavePath string
		wantError           string
		newProvider         func() sqlProvider
	}{
		{
			name:                "mysql",
			savePath:            "user:password@unix(/nonexistent/beego-session-test.sock)/database",
			conflictingSavePath: "other:password@unix(/nonexistent/beego-session-test.sock)/database",
			wantError:           "mysql session provider is already initialized with a different configuration",
			newProvider: func() sqlProvider {
				return &mysql.Provider{}
			},
		},
		{
			name:                "postgresql",
			savePath:            "host=/nonexistent/beego-session-test user=test dbname=test sslmode=disable",
			conflictingSavePath: "host=/nonexistent/beego-session-test user=other dbname=test sslmode=disable",
			wantError:           "postgresql session provider is already initialized with a different configuration",
			newProvider: func() sqlProvider {
				return &postgres.Provider{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const (
				iterations  = 25
				maxLifetime = int64(3600)
			)
			ctx := context.Background()
			provider := tt.newProvider()
			if err := provider.SessionInit(ctx, maxLifetime, tt.savePath); err != nil {
				t.Fatalf("first SessionInit returned an error: %v", err)
			}
			if err := provider.SessionInit(ctx, maxLifetime, tt.savePath); err != nil {
				t.Fatalf("second SessionInit returned an error: %v", err)
			}
			if err := provider.SessionInit(ctx, maxLifetime+1, tt.savePath); err != nil {
				t.Fatalf("SessionInit with a new maximum lifetime returned an error: %v", err)
			}

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					lifetime := maxLifetime + int64(i%2)
					if err := provider.SessionInit(ctx, lifetime, tt.savePath); err != nil {
						t.Errorf("SessionInit returned an error: %v", err)
						return
					}
				}
			}()
			go func() {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					provider.SessionGC(ctx)
				}
			}()
			wg.Wait()

			err := provider.SessionInit(ctx, maxLifetime, tt.conflictingSavePath)
			if err == nil {
				t.Fatal("SessionInit accepted a different save path")
			} else if err.Error() != tt.wantError {
				t.Errorf("SessionInit error = %q; want %q", err, tt.wantError)
			}
		})
	}
}
