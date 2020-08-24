// Copyright 2020 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"context"
	"time"
)

// Invocation represents an "Orm" invocation
type Invocation struct {
	Method string
	// Md may be nil in some cases. It depends on method
	Md interface{}
	// the args are all arguments except context.Context
	Args []interface{}

	mi *modelInfo
	// f is the Orm operation
	f func(ctx context.Context) []interface{}

	// insideTx indicates whether this is inside a transaction
	InsideTx    bool
	TxStartTime time.Time
	TxName      string
}

func (inv *Invocation) GetTableName() string {
	if inv.mi != nil {
		return inv.mi.table
	}
	return ""
}

func (inv *Invocation) execute(ctx context.Context) []interface{} {
	return inv.f(ctx)
}

// GetPkFieldName return the primary key of this table
// if not found, "" is returned
func (inv *Invocation) GetPkFieldName() string {
	if inv.mi.fields.pk != nil {
		return inv.mi.fields.pk.name
	}
	return ""
}
