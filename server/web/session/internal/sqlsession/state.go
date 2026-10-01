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
	"database/sql"
	"fmt"
	"sync"
)

// PoolState stores the shared database pool and its session configuration.
type PoolState struct {
	lock        sync.RWMutex
	db          *sql.DB
	savePath    string
	maxLifetime int64
}

// Init initializes the pool or updates the lifetime for the same save path.
func (s *PoolState) Init(driverName, providerName string, maxLifetime int64, savePath string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.db != nil {
		if s.savePath != savePath {
			return fmt.Errorf("%s session provider is already initialized with a different configuration", providerName)
		}
		s.maxLifetime = maxLifetime
		return nil
	}

	db, err := sql.Open(driverName, savePath)
	if err != nil {
		return err
	}
	s.db = db
	s.savePath = savePath
	s.maxLifetime = maxLifetime
	return nil
}

// DB returns the shared database pool.
func (s *PoolState) DB() *sql.DB {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.db
}

// DBAndMaxLifetime returns the pool and current session lifetime together.
func (s *PoolState) DBAndMaxLifetime() (*sql.DB, int64) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.db, s.maxLifetime
}
