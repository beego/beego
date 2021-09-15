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

package mock

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/orm"
)

const mockErrorMsg = "mock error"

func init() {
	orm.RegisterModel(&User{})
}

func TestMockDBStats(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	stats := &sql.DBStats{}
	s.Mock(MockDBStats(stats))

	o := orm.NewOrm()

	res := o.DBStats()

	assert.Equal(t, stats, res)
}

func TestMockDeleteWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	s.Mock(MockDeleteWithCtx((&User{}).TableName(), 12, nil))
	o := orm.NewOrm()
	rows, err := o.Delete(&User{})
	assert.Equal(t, int64(12), rows)
	assert.Nil(t, err)
}

func TestMockInsertOrUpdateWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	s.Mock(MockInsertOrUpdateWithCtx((&User{}).TableName(), 12, nil))
	o := orm.NewOrm()
	id, err := o.InsertOrUpdate(&User{})
	assert.Equal(t, int64(12), id)
	assert.Nil(t, err)
}

func TestMockRead(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	err := errors.New(mockErrorMsg)
	s.Mock(MockRead((&User{}).TableName(), func(data interface{}) {
		u := data.(*User)
		u.Name = "Tom"
	}, err))
	o := orm.NewOrm()
	u := &User{}
	e := o.Read(u)
	assert.Equal(t, err, e)
	assert.Equal(t, "Tom", u.Name)
}

func TestMockQueryM2MWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := &DoNothingQueryM2Mer{}
	s.Mock(MockQueryM2MWithCtx((&User{}).TableName(), "Tags", mock))
	o := orm.NewOrm()
	res := o.QueryM2M(&User{}, "Tags")
	assert.Equal(t, mock, res)
}

func TestMockQueryTableWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := &DoNothingQuerySetter{}
	s.Mock(MockQueryTableWithCtx((&User{}).TableName(), mock))
	o := orm.NewOrm()
	res := o.QueryTable(&User{})
	assert.Equal(t, mock, res)
}

func TestMockTable(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockTable((&User{}).TableName(), mock))
	o := orm.NewOrm()
	res := o.Read(&User{})
	assert.Equal(t, mock, res)
}

func TestMockInsertMultiWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockInsertMultiWithCtx((&User{}).TableName(), 12, mock))
	o := orm.NewOrm()
	res, err := o.InsertMulti(11, []interface{}{&User{}})
	assert.Equal(t, int64(12), res)
	assert.Equal(t, mock, err)
}

func TestMockInsertWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockInsertWithCtx((&User{}).TableName(), 13, mock))
	o := orm.NewOrm()
	res, err := o.Insert(&User{})
	assert.Equal(t, int64(13), res)
	assert.Equal(t, mock, err)
}

func TestMockUpdateWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockUpdateWithCtx((&User{}).TableName(), 12, mock))
	o := orm.NewOrm()
	res, err := o.Update(&User{})
	assert.Equal(t, int64(12), res)
	assert.Equal(t, mock, err)
}

func TestMockLoadRelatedWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockLoadRelatedWithCtx((&User{}).TableName(), "T", 12, mock))
	o := orm.NewOrm()
	res, err := o.LoadRelated(&User{}, "T")
	assert.Equal(t, int64(12), res)
	assert.Equal(t, mock, err)
}

func TestMockMethod(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockMethod("ReadWithCtx", mock))
	o := orm.NewOrm()
	err := o.Read(&User{})
	assert.Equal(t, mock, err)
}

func TestMockReadForUpdateWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockReadForUpdateWithCtx((&User{}).TableName(), func(data interface{}) {
		u := data.(*User)
		u.Name = "Tom"
	}, mock))
	o := orm.NewOrm()
	u := &User{}
	err := o.ReadForUpdate(u)
	assert.Equal(t, mock, err)
	assert.Equal(t, "Tom", u.Name)
}

func TestMockRawWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := &DoNothingRawSetter{}
	s.Mock(MockRawWithCtx(mock))
	o := orm.NewOrm()
	res := o.Raw("")
	assert.Equal(t, mock, res)
}

func TestMockReadOrCreateWithCtx(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockReadOrCreateWithCtx((&User{}).TableName(), func(data interface{}) {
		u := data.(*User)
		u.Name = "Tom"
	}, false, 12, mock))
	o := orm.NewOrm()
	u := &User{}
	inserted, id, err := o.ReadOrCreate(u, "")
	assert.Equal(t, mock, err)
	assert.Equal(t, int64(12), id)
	assert.False(t, inserted)
	assert.Equal(t, "Tom", u.Name)
}

func TestTransactionClosure(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockRead((&User{}).TableName(), func(data interface{}) {
		u := data.(*User)
		u.Name = "Tom"
	}, mock))
	u, err := originalTxUsingClosure()
	assert.Equal(t, mock, err)
	assert.Equal(t, "Tom", u.Name)
}

func TestTransactionManually(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockRead((&User{}).TableName(), func(data interface{}) {
		u := data.(*User)
		u.Name = "Tom"
	}, mock))
	u, err := originalTxManually()
	assert.Equal(t, mock, err)
	assert.Equal(t, "Tom", u.Name)
}

func TestTransactionRollback(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockRead((&User{}).TableName(), nil, errors.New("read error")))
	s.Mock(MockRollback(mock))
	_, err := originalTx()
	assert.Equal(t, mock, err)
}

func TestTransactionRollbackUnlessCommit(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockRollbackUnlessCommit(mock))

	// u := &User{}
	o := orm.NewOrm()
	txOrm, _ := o.Begin()
	err := txOrm.RollbackUnlessCommit()
	assert.Equal(t, mock, err)
}

func TestTransactionCommit(t *testing.T) {
	s := StartMock()
	defer s.Clear()
	mock := errors.New(mockErrorMsg)
	s.Mock(MockRead((&User{}).TableName(), func(data interface{}) {
		u := data.(*User)
		u.Name = "Tom"
	}, nil))
	s.Mock(MockCommit(mock))
	u, err := originalTx()
	assert.Equal(t, mock, err)
	assert.Equal(t, "Tom", u.Name)
}

func originalTx() (*User, error) {
	u := &User{}
	o := orm.NewOrm()
	txOrm, _ := o.Begin()
	err := txOrm.Read(u)
	if err == nil {
		err = txOrm.Commit()
		return u, err
	} else {
		err = txOrm.Rollback()
		return nil, err
	}
}

func originalTxManually() (*User, error) {
	u := &User{}
	o := orm.NewOrm()
	txOrm, _ := o.Begin()
	err := txOrm.Read(u)
	_ = txOrm.Commit()
	return u, err
}

func originalTxUsingClosure() (*User, error) {
	u := &User{}
	var err error
	o := orm.NewOrm()
	_ = o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		err = txOrm.Read(u)
		return nil
	})
	return u, err
}

type User struct {
	Id   int
	Name string
}

func (u *User) TableName() string {
	return "user"
}
