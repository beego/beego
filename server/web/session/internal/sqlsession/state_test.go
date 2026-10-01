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

package sqlsession

import (
	"strings"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql" // Register the MySQL driver used by these initialization tests.
)

func TestPoolStateReusesDatabaseAndUpdatesLifetime(t *testing.T) {
	const savePath = "user@tcp(127.0.0.1:3306)/database"
	state := &PoolState{}
	if err := state.Init("mysql", "test", 3600, savePath); err != nil {
		t.Fatalf("first Init returned an error: %v", err)
	}
	first := state.DB()
	t.Cleanup(func() {
		_ = first.Close()
	})
	if got := first.Stats().OpenConnections; got != 0 {
		t.Fatalf("Init opened %d physical connections; want 0", got)
	}

	if err := state.Init("mysql", "test", 7200, savePath); err != nil {
		t.Fatalf("second Init returned an error: %v", err)
	}
	second, maxLifetime := state.DBAndMaxLifetime()
	if second != first {
		t.Fatal("Init replaced the database pool for the same configuration")
	}
	if maxLifetime != 7200 {
		t.Fatalf("maximum lifetime = %d; want 7200", maxLifetime)
	}
}

func TestPoolStateRejectsSavePathChangesWithoutExposingConfiguration(t *testing.T) {
	const (
		initialSavePath     = "initial-user@tcp(127.0.0.1:3306)/database"
		conflictingSavePath = "conflicting-user@tcp(127.0.0.1:3306)/database"
	)
	state := &PoolState{}
	if err := state.Init("mysql", "test", 3600, initialSavePath); err != nil {
		t.Fatalf("first Init returned an error: %v", err)
	}
	first := state.DB()
	t.Cleanup(func() {
		_ = first.Close()
	})

	err := state.Init("mysql", "test", 3600, conflictingSavePath)
	if err == nil {
		t.Fatal("Init accepted a different save path")
	}
	if strings.Contains(err.Error(), initialSavePath) || strings.Contains(err.Error(), conflictingSavePath) {
		t.Fatalf("Init error exposed a database configuration: %q", err)
	}
	if state.DB() != first {
		t.Fatal("Init replaced the database pool after rejecting a different save path")
	}
}

func TestPoolStateConcurrentInitAndRead(t *testing.T) {
	const (
		iterations = 25
		savePath   = "user@tcp(127.0.0.1:3306)/database"
	)
	state := &PoolState{}
	if err := state.Init("mysql", "test", 3600, savePath); err != nil {
		t.Fatalf("Init returned an error: %v", err)
	}
	db := state.DB()
	t.Cleanup(func() {
		_ = db.Close()
	})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if err := state.Init("mysql", "test", 3600+int64(i%2), savePath); err != nil {
				t.Errorf("Init returned an error: %v", err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			current, _ := state.DBAndMaxLifetime()
			if current != db {
				t.Error("Init replaced the database pool")
				return
			}
		}
	}()
	wg.Wait()
}
