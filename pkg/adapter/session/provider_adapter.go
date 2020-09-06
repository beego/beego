// Copyright 2020
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

package session

import (
	"context"

	"github.com/astaxie/beego/pkg/infrastructure/session"
)

type oldToNewProviderAdapter struct {
	delegate Provider
}

func (o *oldToNewProviderAdapter) SessionInit(ctx context.Context, gclifetime int64, config string) error {
	return o.delegate.SessionInit(gclifetime, config)
}

func (o *oldToNewProviderAdapter) SessionRead(ctx context.Context, sid string) (session.Store, error) {
	store, err := o.delegate.SessionRead(sid)
	return &oldToNewStoreAdapter{
		delegate: store,
	}, err
}

func (o *oldToNewProviderAdapter) SessionExist(ctx context.Context, sid string) (bool, error) {
	return o.delegate.SessionExist(sid), nil
}

func (o *oldToNewProviderAdapter) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	s, err := o.delegate.SessionRegenerate(oldsid, sid)
	return &oldToNewStoreAdapter{
		delegate: s,
	}, err
}

func (o *oldToNewProviderAdapter) SessionDestroy(ctx context.Context, sid string) error {
	return o.delegate.SessionDestroy(sid)
}

func (o *oldToNewProviderAdapter) SessionAll(ctx context.Context) int {
	return o.delegate.SessionAll()
}

func (o *oldToNewProviderAdapter) SessionGC(ctx context.Context) {
	o.delegate.SessionGC()
}

type newToOldProviderAdapter struct {
	delegate session.Provider
}

func (n *newToOldProviderAdapter) SessionInit(gclifetime int64, config string) error {
	return n.delegate.SessionInit(context.Background(), gclifetime, config)
}

func (n *newToOldProviderAdapter) SessionRead(sid string) (Store, error) {
	s, err := n.delegate.SessionRead(context.Background(), sid)
	if adt, ok := s.(*oldToNewStoreAdapter); err == nil && ok {
		return adt.delegate, err
	}
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

func (n *newToOldProviderAdapter) SessionExist(sid string) bool {
	res, _ := n.delegate.SessionExist(context.Background(), sid)
	return res
}

func (n *newToOldProviderAdapter) SessionRegenerate(oldsid, sid string) (Store, error) {
	s, err := n.delegate.SessionRegenerate(context.Background(), oldsid, sid)
	if adt, ok := s.(*oldToNewStoreAdapter); err == nil && ok {
		return adt.delegate, err
	}
	return &NewToOldStoreAdapter{
		delegate: s,
	}, err
}

func (n *newToOldProviderAdapter) SessionDestroy(sid string) error {
	return n.delegate.SessionDestroy(context.Background(), sid)
}

func (n *newToOldProviderAdapter) SessionAll() int {
	return n.delegate.SessionAll(context.Background())
}

func (n *newToOldProviderAdapter) SessionGC() {
	n.delegate.SessionGC(context.Background())
}
