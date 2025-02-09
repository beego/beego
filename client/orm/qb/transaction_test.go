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
	"context"
	"database/sql"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTx_Commit(t *testing.T) {
	ormer := orm.NewMockOrmer(gomock.NewController(t))
	db, err := OpenDB("mysql", ormer)
	assert.Nil(t, err)
	ormer.EXPECT().BeginWithCtxAndOpts(gomock.Any(), &sql.TxOptions{}).Return(nil, errors.New("begin failed"))
	// Begin 失败
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	assert.Equal(t, errors.New("begin failed"), err)
	assert.Nil(t, tx)

	ormer.EXPECT().BeginWithCtxAndOpts(gomock.Any(), &sql.TxOptions{}).Return(orm.NewMockTxOrmer(gomock.NewController(t)), nil)
	tx, err = db.BeginTx(context.Background(), &sql.TxOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, tx)
}

func TestTx_Rollback(t *testing.T) {
	ormer := orm.NewMockOrmer(gomock.NewController(t))
	db, err := OpenDB("mysql", ormer)
	assert.Nil(t, err)
	txmer := orm.NewMockTxOrmer(gomock.NewController(t))
	ormer.EXPECT().BeginWithCtxAndOpts(gomock.Any(), &sql.TxOptions{}).Return(txmer, nil)
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	ormer.EXPECT()
	txmer.EXPECT().Rollback().Return(nil)
	err = tx.Rollback()
	assert.Nil(t, err)
}
