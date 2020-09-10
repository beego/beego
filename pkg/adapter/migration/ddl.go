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

package migration

import (
	"github.com/astaxie/beego/pkg/client/orm/migration"
)

// Index struct defines the structure of Index Columns
type Index migration.Index

// Unique struct defines a single unique key combination
type Unique migration.Unique

// Column struct defines a single column of a table
type Column migration.Column

// Foreign struct defines a single foreign relationship
type Foreign migration.Foreign

// RenameColumn struct allows renaming of columns
type RenameColumn migration.RenameColumn

// CreateTable creates the table on system
func (m *Migration) CreateTable(tablename, engine, charset string, p ...func()) {
	(*migration.Migration)(m).CreateTable(tablename, engine, charset, p...)
}

// AlterTable set the ModifyType to alter
func (m *Migration) AlterTable(tablename string) {
	(*migration.Migration)(m).AlterTable(tablename)
}

// NewCol creates a new standard column and attaches it to m struct
func (m *Migration) NewCol(name string) *Column {
	return (*Column)((*migration.Migration)(m).NewCol(name))
}

// PriCol creates a new primary column and attaches it to m struct
func (m *Migration) PriCol(name string) *Column {
	return (*Column)((*migration.Migration)(m).PriCol(name))
}

// UniCol creates / appends columns to specified unique key and attaches it to m struct
func (m *Migration) UniCol(uni, name string) *Column {
	return (*Column)((*migration.Migration)(m).UniCol(uni, name))
}

// ForeignCol creates a new foreign column and returns the instance of column
func (m *Migration) ForeignCol(colname, foreigncol, foreigntable string) (foreign *Foreign) {
	return (*Foreign)((*migration.Migration)(m).ForeignCol(colname, foreigncol, foreigntable))
}

// SetOnDelete sets the on delete of foreign
func (foreign *Foreign) SetOnDelete(del string) *Foreign {
	(*migration.Foreign)(foreign).SetOnDelete(del)
	return foreign
}

// SetOnUpdate sets the on update of foreign
func (foreign *Foreign) SetOnUpdate(update string) *Foreign {
	(*migration.Foreign)(foreign).SetOnUpdate(update)
	return foreign
}

// Remove marks the columns to be removed.
// it allows reverse m to create the column.
func (c *Column) Remove() {
	(*migration.Column)(c).Remove()
}

// SetAuto enables auto_increment of column (can be used once)
func (c *Column) SetAuto(inc bool) *Column {
	(*migration.Column)(c).SetAuto(inc)
	return c
}

// SetNullable sets the column to be null
func (c *Column) SetNullable(null bool) *Column {
	(*migration.Column)(c).SetNullable(null)
	return c
}

// SetDefault sets the default value, prepend with "DEFAULT "
func (c *Column) SetDefault(def string) *Column {
	(*migration.Column)(c).SetDefault(def)
	return c
}

// SetUnsigned sets the column to be unsigned int
func (c *Column) SetUnsigned(unsign bool) *Column {
	(*migration.Column)(c).SetUnsigned(unsign)
	return c
}

// SetDataType sets the dataType of the column
func (c *Column) SetDataType(dataType string) *Column {
	(*migration.Column)(c).SetDataType(dataType)
	return c
}

// SetOldNullable allows reverting to previous nullable on reverse ms
func (c *RenameColumn) SetOldNullable(null bool) *RenameColumn {
	(*migration.RenameColumn)(c).SetOldNullable(null)
	return c
}

// SetOldDefault allows reverting to previous default on reverse ms
func (c *RenameColumn) SetOldDefault(def string) *RenameColumn {
	(*migration.RenameColumn)(c).SetOldDefault(def)
	return c
}

// SetOldUnsigned allows reverting to previous unsgined on reverse ms
func (c *RenameColumn) SetOldUnsigned(unsign bool) *RenameColumn {
	(*migration.RenameColumn)(c).SetOldUnsigned(unsign)
	return c
}

// SetOldDataType allows reverting to previous datatype on reverse ms
func (c *RenameColumn) SetOldDataType(dataType string) *RenameColumn {
	(*migration.RenameColumn)(c).SetOldDataType(dataType)
	return c
}

// SetPrimary adds the columns to the primary key (can only be used any number of times in only one m)
func (c *Column) SetPrimary(m *Migration) *Column {
	(*migration.Column)(c).SetPrimary((*migration.Migration)(m))
	return c
}

// AddColumnsToUnique adds the columns to Unique Struct
func (unique *Unique) AddColumnsToUnique(columns ...*Column) *Unique {
	cls := toNewColumnsArray(columns)
	(*migration.Unique)(unique).AddColumnsToUnique(cls...)
	return unique
}

// AddColumns adds columns to m struct
func (m *Migration) AddColumns(columns ...*Column) *Migration {
	cls := toNewColumnsArray(columns)
	(*migration.Migration)(m).AddColumns(cls...)
	return m
}

func toNewColumnsArray(columns []*Column) []*migration.Column {
	cls := make([]*migration.Column, 0, len(columns))
	for _, c := range columns {
		cls = append(cls, (*migration.Column)(c))
	}
	return cls
}

// AddPrimary adds the column to primary in m struct
func (m *Migration) AddPrimary(primary *Column) *Migration {
	(*migration.Migration)(m).AddPrimary((*migration.Column)(primary))
	return m
}

// AddUnique adds the column to unique in m struct
func (m *Migration) AddUnique(unique *Unique) *Migration {
	(*migration.Migration)(m).AddUnique((*migration.Unique)(unique))
	return m
}

// AddForeign adds the column to foreign in m struct
func (m *Migration) AddForeign(foreign *Foreign) *Migration {
	(*migration.Migration)(m).AddForeign((*migration.Foreign)(foreign))
	return m
}

// AddIndex adds the column to index in m struct
func (m *Migration) AddIndex(index *Index) *Migration {
	(*migration.Migration)(m).AddIndex((*migration.Index)(index))
	return m
}

// RenameColumn allows renaming of columns
func (m *Migration) RenameColumn(from, to string) *RenameColumn {
	return (*RenameColumn)((*migration.Migration)(m).RenameColumn(from, to))
}

// GetSQL returns the generated sql depending on ModifyType
func (m *Migration) GetSQL() (sql string) {
	return (*migration.Migration)(m).GetSQL()
}
