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

// Package migration is used for migration
//
// The table structure is as follow:
//
//	CREATE TABLE `migrations` (
//		`id_migration` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate key',
//		`name` varchar(255) DEFAULT NULL COMMENT 'migration name, unique',
//		`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date migrated or rolled back',
//		`statements` longtext COMMENT 'SQL statements for this migration',
//		`rollback_statements` longtext,
//		`status` enum('update','rollback') DEFAULT NULL COMMENT 'update indicates it is a normal migration while rollback means this migration is rolled back',
//		PRIMARY KEY (`id_migration`)
//	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
package migration

import (
	"github.com/beego/beego/client/orm/migration"
)

// const the data format for the bee generate migration datatype
const (
	DateFormat   = "20060102_150405"
	DBDateFormat = "2006-01-02 15:04:05"
)

// Migrationer is an interface for all Migration struct
type Migrationer interface {
	Up()
	Down()
	Reset()
	Exec(name, status string) error
	GetCreated() int64
}

// Migration defines the migrations by either SQL or DDL
type Migration migration.Migration

// Up implement in the Inheritance struct for upgrade
func (m *Migration) Up() {
	(*migration.Migration)(m).Up()
}

// Down implement in the Inheritance struct for down
func (m *Migration) Down() {
	(*migration.Migration)(m).Down()
}

// Migrate adds the SQL to the execution list
func (m *Migration) Migrate(migrationType string) {
	(*migration.Migration)(m).Migrate(migrationType)
}

// SQL add sql want to execute
func (m *Migration) SQL(sql string) {
	(*migration.Migration)(m).SQL(sql)
}

// Reset the sqls
func (m *Migration) Reset() {
	(*migration.Migration)(m).Reset()
}

// Exec execute the sql already add in the sql
func (m *Migration) Exec(name, status string) error {
	return (*migration.Migration)(m).Exec(name, status)
}

// GetCreated get the unixtime from the Created
func (m *Migration) GetCreated() int64 {
	return (*migration.Migration)(m).GetCreated()
}

// Register register the Migration in the map
func Register(name string, m Migrationer) error {
	return migration.Register(name, m)
}

// Upgrade upgrade the migration from lasttime
func Upgrade(lasttime int64) error {
	return migration.Upgrade(lasttime)
}

// Rollback rollback the migration by the name
func Rollback(name string) error {
	return migration.Rollback(name)
}

// Reset reset all migration
// run all migration's down function
func Reset() error {
	return migration.Reset()
}

// Refresh first Reset, then Upgrade
func Refresh() error {
	return migration.Refresh()
}
