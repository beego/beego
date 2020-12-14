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
	"net/http"

	"github.com/beego/beego/v2/server/web/session"
)

type NewToOldStoreAdapter struct {
	delegate session.Store
}

func CreateNewToOldStoreAdapter(s session.Store) Store {
	return &NewToOldStoreAdapter{
		delegate: s,
	}
}

func (n *NewToOldStoreAdapter) Set(key, value interface{}) error {
	return n.delegate.Set(context.Background(), key, value)
}

func (n *NewToOldStoreAdapter) Get(key interface{}) interface{} {
	return n.delegate.Get(context.Background(), key)
}

func (n *NewToOldStoreAdapter) Delete(key interface{}) error {
	return n.delegate.Delete(context.Background(), key)
}

func (n *NewToOldStoreAdapter) SessionID() string {
	return n.delegate.SessionID(context.Background())
}

func (n *NewToOldStoreAdapter) SessionRelease(w http.ResponseWriter) {
	n.delegate.SessionRelease(context.Background(), w)
}

func (n *NewToOldStoreAdapter) Flush() error {
	return n.delegate.Flush(context.Background())
}

type oldToNewStoreAdapter struct {
	delegate Store
}

func (o *oldToNewStoreAdapter) Set(ctx context.Context, key, value interface{}) error {
	return o.delegate.Set(key, value)
}

func (o *oldToNewStoreAdapter) Get(ctx context.Context, key interface{}) interface{} {
	return o.delegate.Get(key)
}

func (o *oldToNewStoreAdapter) Delete(ctx context.Context, key interface{}) error {
	return o.delegate.Delete(key)
}

func (o *oldToNewStoreAdapter) SessionID(ctx context.Context) string {
	return o.delegate.SessionID()
}

func (o *oldToNewStoreAdapter) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	o.delegate.SessionRelease(w)
}

func (o *oldToNewStoreAdapter) Flush(ctx context.Context) error {
	return o.delegate.Flush()
}
