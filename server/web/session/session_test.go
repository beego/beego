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

package session

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

type contextTestKey string

const (
	testCtxKey contextTestKey = "requestContextKey"
	testCtxVal string         = "requestContextVal"
)

type mockStore struct {
	sid string
}

func (*mockStore) Set(_ context.Context, _, _ interface{}) error {
	return nil
}

func (*mockStore) Get(_ context.Context, _ interface{}) interface{} {
	return nil
}

func (*mockStore) Delete(_ context.Context, _ interface{}) error {
	return nil
}

func (m *mockStore) SessionID(_ context.Context) string {
	return m.sid
}

func (*mockStore) SessionReleaseIfPresent(_ context.Context, _ http.ResponseWriter) {}

func (*mockStore) SessionRelease(_ context.Context, _ http.ResponseWriter) {}

func (*mockStore) Flush(_ context.Context) error {
	return nil
}

type mockContextAwareProvider struct {
	lock         sync.Mutex
	onRead       func(ctx context.Context, sid string)
	onExist      func(ctx context.Context, sid string)
	onRegenerate func(ctx context.Context, oldsid, sid string)
	onDestroy    func(ctx context.Context, sid string)
	onAll        func(ctx context.Context)
	onGC         func(ctx context.Context)
	existResult  bool
}

func (*mockContextAwareProvider) SessionInit(_ context.Context, _ int64, _ string) error {
	return nil
}

func (p *mockContextAwareProvider) SessionRead(ctx context.Context, sid string) (Store, error) {
	if ctx == nil {
		panic("SessionRead received nil context")
	}
	p.lock.Lock()
	fn := p.onRead
	p.lock.Unlock()
	if fn != nil {
		fn(ctx, sid)
	}
	return &mockStore{sid: sid}, nil
}

func (p *mockContextAwareProvider) SessionExist(ctx context.Context, sid string) (bool, error) {
	if ctx == nil {
		panic("SessionExist received nil context")
	}
	p.lock.Lock()
	fn := p.onExist
	res := p.existResult
	p.lock.Unlock()
	if fn != nil {
		fn(ctx, sid)
	}
	return res, nil
}

func (p *mockContextAwareProvider) SessionRegenerate(ctx context.Context, oldsid, sid string) (Store, error) {
	if ctx == nil {
		panic("SessionRegenerate received nil context")
	}
	p.lock.Lock()
	fn := p.onRegenerate
	p.lock.Unlock()
	if fn != nil {
		fn(ctx, oldsid, sid)
	}
	return &mockStore{sid: sid}, nil
}

func (p *mockContextAwareProvider) SessionDestroy(ctx context.Context, sid string) error {
	if ctx == nil {
		panic("SessionDestroy received nil context")
	}
	p.lock.Lock()
	fn := p.onDestroy
	p.lock.Unlock()
	if fn != nil {
		fn(ctx, sid)
	}
	return nil
}

func (p *mockContextAwareProvider) SessionAll(ctx context.Context) int {
	if ctx == nil {
		panic("SessionAll received nil context")
	}
	p.lock.Lock()
	fn := p.onAll
	p.lock.Unlock()
	if fn != nil {
		fn(ctx)
	}
	return 42
}

func (p *mockContextAwareProvider) SessionGC(ctx context.Context) {
	if ctx == nil {
		panic("SessionGC received nil context")
	}
	p.lock.Lock()
	fn := p.onGC
	p.lock.Unlock()
	if fn != nil {
		fn(ctx)
	}
}

func setupTestManager(providerName string, provider *mockContextAwareProvider) (*Manager, error) {
	Register(providerName, provider)
	cf := &ManagerConfig{
		CookieName:      "gosessionid",
		Gclifetime:      3600,
		Maxlifetime:     3600,
		EnableSetCookie: true,
	}
	return NewManager(providerName, cf)
}

func TestSessionStart_PropagatesRequestContext_NewSession(t *testing.T) {
	var mu sync.Mutex
	var receivedCtx context.Context

	provider := &mockContextAwareProvider{
		onRead: func(ctx context.Context, _ string) {
			mu.Lock()
			receivedCtx = ctx
			mu.Unlock()
		},
	}
	manager, err := setupTestManager("mock_new_sess", provider)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), testCtxKey, testCtxVal))
	defer cancel()

	req := httptest.NewRequest("GET", "http://example.com/", http.NoBody).WithContext(ctx)
	w := httptest.NewRecorder()

	store, err := manager.SessionStart(w, req)
	if err != nil {
		t.Fatalf("SessionStart error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	mu.Lock()
	captured := receivedCtx
	mu.Unlock()

	if captured == nil {
		t.Fatal("expected non-nil context in SessionRead")
	}
	if val := captured.Value(testCtxKey); val != testCtxVal {
		t.Fatalf("expected context value %q, got %v", testCtxVal, val)
	}

	cancel()
	select {
	case <-captured.Done():
		if captured.Err() != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", captured.Err())
		}
	default:
		t.Fatal("expected captured context to be canceled")
	}
}

func TestSessionStart_PropagatesRequestContext_ExistingSession(t *testing.T) {
	var mu sync.Mutex
	var existCtx, readCtx context.Context

	provider := &mockContextAwareProvider{
		existResult: true,
		onExist: func(ctx context.Context, _ string) {
			mu.Lock()
			existCtx = ctx
			mu.Unlock()
		},
		onRead: func(ctx context.Context, _ string) {
			mu.Lock()
			readCtx = ctx
			mu.Unlock()
		},
	}
	manager, err := setupTestManager("mock_exist_sess", provider)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), testCtxKey, testCtxVal))
	defer cancel()

	existingSid := "existing_sid_12345"
	req := httptest.NewRequest("GET", "http://example.com/", http.NoBody).WithContext(ctx)
	req.AddCookie(&http.Cookie{
		Name:  manager.config.CookieName,
		Value: url.QueryEscape(existingSid),
	})
	w := httptest.NewRecorder()

	store, err := manager.SessionStart(w, req)
	if err != nil {
		t.Fatalf("SessionStart error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	mu.Lock()
	capturedExist := existCtx
	capturedRead := readCtx
	mu.Unlock()

	if capturedExist == nil {
		t.Fatal("expected non-nil context in SessionExist")
	}
	if val := capturedExist.Value(testCtxKey); val != testCtxVal {
		t.Fatalf("SessionExist expected context value %q, got %v", testCtxVal, val)
	}

	if capturedRead == nil {
		t.Fatal("expected non-nil context in SessionRead")
	}
	if val := capturedRead.Value(testCtxKey); val != testCtxVal {
		t.Fatalf("SessionRead expected context value %q, got %v", testCtxVal, val)
	}
}

func TestSessionStart_PropagatesRequestContext_SessionNotFoundGeneratesNew(t *testing.T) {
	var mu sync.Mutex
	var existCtx, readCtx context.Context

	provider := &mockContextAwareProvider{
		existResult: false,
		onExist: func(ctx context.Context, _ string) {
			mu.Lock()
			existCtx = ctx
			mu.Unlock()
		},
		onRead: func(ctx context.Context, _ string) {
			mu.Lock()
			readCtx = ctx
			mu.Unlock()
		},
	}
	manager, err := setupTestManager("mock_notfound_sess", provider)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	ctx := context.WithValue(context.Background(), testCtxKey, testCtxVal)
	req := httptest.NewRequest("GET", "http://example.com/", http.NoBody).WithContext(ctx)
	req.AddCookie(&http.Cookie{
		Name:  manager.config.CookieName,
		Value: url.QueryEscape("stale_sid"),
	})
	w := httptest.NewRecorder()

	store, err := manager.SessionStart(w, req)
	if err != nil {
		t.Fatalf("SessionStart error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	mu.Lock()
	capturedExist := existCtx
	capturedRead := readCtx
	mu.Unlock()

	if capturedExist == nil || capturedExist.Value(testCtxKey) != testCtxVal {
		t.Fatalf("SessionExist context missing value %v", capturedExist)
	}
	if capturedRead == nil || capturedRead.Value(testCtxKey) != testCtxVal {
		t.Fatalf("SessionRead context missing value %v", capturedRead)
	}
}

func TestSessionDestroy_PropagatesRequestContext(t *testing.T) {
	var mu sync.Mutex
	var destroyCtx context.Context

	provider := &mockContextAwareProvider{
		onDestroy: func(ctx context.Context, _ string) {
			mu.Lock()
			destroyCtx = ctx
			mu.Unlock()
		},
	}
	manager, err := setupTestManager("mock_destroy_sess", provider)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), testCtxKey, testCtxVal))
	defer cancel()

	sid := "sid_to_destroy"
	req := httptest.NewRequest("POST", "http://example.com/logout", http.NoBody).WithContext(ctx)
	req.AddCookie(&http.Cookie{
		Name:  manager.config.CookieName,
		Value: url.QueryEscape(sid),
	})
	w := httptest.NewRecorder()

	manager.SessionDestroy(w, req)

	mu.Lock()
	captured := destroyCtx
	mu.Unlock()

	if captured == nil {
		t.Fatal("expected non-nil context in SessionDestroy")
	}
	if val := captured.Value(testCtxKey); val != testCtxVal {
		t.Fatalf("expected context value %q, got %v", testCtxVal, val)
	}

	cancel()
	select {
	case <-captured.Done():
		if captured.Err() != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", captured.Err())
		}
	default:
		t.Fatal("expected destroyCtx to be canceled")
	}
}

func TestSessionRegenerateID_PropagatesRequestContext_ExistingCookie(t *testing.T) {
	var mu sync.Mutex
	var regenCtx context.Context

	provider := &mockContextAwareProvider{
		onRegenerate: func(ctx context.Context, _, _ string) {
			mu.Lock()
			regenCtx = ctx
			mu.Unlock()
		},
	}
	manager, err := setupTestManager("mock_regen_exist", provider)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), testCtxKey, testCtxVal))
	defer cancel()

	oldSid := "old_sid_to_rotate"
	req := httptest.NewRequest("POST", "http://example.com/login", http.NoBody).WithContext(ctx)
	req.AddCookie(&http.Cookie{
		Name:  manager.config.CookieName,
		Value: url.QueryEscape(oldSid),
	})
	w := httptest.NewRecorder()

	store, err := manager.SessionRegenerateID(w, req)
	if err != nil {
		t.Fatalf("SessionRegenerateID error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	mu.Lock()
	captured := regenCtx
	mu.Unlock()

	if captured == nil {
		t.Fatal("expected non-nil context in SessionRegenerate")
	}
	if val := captured.Value(testCtxKey); val != testCtxVal {
		t.Fatalf("expected context value %q, got %v", testCtxVal, val)
	}

	cancel()
	select {
	case <-captured.Done():
		if captured.Err() != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", captured.Err())
		}
	default:
		t.Fatal("expected regenCtx to be canceled")
	}
}

func TestSessionRegenerateID_PropagatesRequestContext_NoCookie(t *testing.T) {
	var mu sync.Mutex
	var readCtx context.Context

	provider := &mockContextAwareProvider{
		onRead: func(ctx context.Context, _ string) {
			mu.Lock()
			readCtx = ctx
			mu.Unlock()
		},
	}
	manager, err := setupTestManager("mock_regen_nocookie", provider)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	ctx := context.WithValue(context.Background(), testCtxKey, testCtxVal)
	req := httptest.NewRequest("POST", "http://example.com/login", http.NoBody).WithContext(ctx)
	w := httptest.NewRecorder()

	store, err := manager.SessionRegenerateID(w, req)
	if err != nil {
		t.Fatalf("SessionRegenerateID error: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	mu.Lock()
	captured := readCtx
	mu.Unlock()

	if captured == nil {
		t.Fatal("expected non-nil context in SessionRead")
	}
	if val := captured.Value(testCtxKey); val != testCtxVal {
		t.Fatalf("expected context value %q, got %v", testCtxVal, val)
	}
}

func TestGC_SuppliesNonNilContext(t *testing.T) {
	var mu sync.Mutex
	var gcCtx context.Context

	provider := &mockContextAwareProvider{
		onGC: func(ctx context.Context) {
			mu.Lock()
			gcCtx = ctx
			mu.Unlock()
		},
	}
	manager := &Manager{
		provider: provider,
		config: &ManagerConfig{
			Gclifetime: 360000,
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("GC() triggered panic with nil context: %v", r)
		}
	}()

	manager.GC()

	mu.Lock()
	captured := gcCtx
	mu.Unlock()

	if captured == nil {
		t.Fatal("GC() did not supply context to SessionGC")
	}
	if captured != context.Background() {
		t.Fatalf("expected context.Background(), got %v", captured)
	}
}

func TestGetActiveSession_SuppliesNonNilContext(t *testing.T) {
	var mu sync.Mutex
	var allCtx context.Context

	provider := &mockContextAwareProvider{
		onAll: func(ctx context.Context) {
			mu.Lock()
			allCtx = ctx
			mu.Unlock()
		},
	}
	manager := &Manager{
		provider: provider,
		config: &ManagerConfig{
			Gclifetime: 3600,
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("GetActiveSession() triggered panic with nil context: %v", r)
		}
	}()

	count := manager.GetActiveSession()
	if count != 42 {
		t.Fatalf("expected 42, got %d", count)
	}

	mu.Lock()
	captured := allCtx
	mu.Unlock()

	if captured == nil {
		t.Fatal("GetActiveSession() did not supply context to SessionAll")
	}
	if captured != context.Background() {
		t.Fatalf("expected context.Background(), got %v", captured)
	}
}

func TestGetSessionStoreWithCtx(t *testing.T) {
	var mu sync.Mutex
	var readCtx context.Context

	provider := &mockContextAwareProvider{
		onRead: func(ctx context.Context, _ string) {
			mu.Lock()
			readCtx = ctx
			mu.Unlock()
		},
	}
	manager := &Manager{
		provider: provider,
		config:   &ManagerConfig{},
	}

	testSid := "sid_get_store"

	// 1. With non-nil context carrying values
	ctx := context.WithValue(context.Background(), testCtxKey, testCtxVal)
	store, err := manager.GetSessionStoreWithCtx(ctx, testSid)
	if err != nil {
		t.Fatalf("GetSessionStoreWithCtx error: %v", err)
	}
	if store == nil || store.SessionID(ctx) != testSid {
		t.Fatalf("unexpected store result: %v", store)
	}

	mu.Lock()
	captured := readCtx
	mu.Unlock()

	if captured == nil || captured.Value(testCtxKey) != testCtxVal {
		t.Fatalf("expected context with value %q, got %v", testCtxVal, captured)
	}

	// 2. With nil context variable (should fall back to context.Background())
	var nilCtx context.Context
	store2, err := manager.GetSessionStoreWithCtx(nilCtx, testSid)
	if err != nil {
		t.Fatalf("GetSessionStoreWithCtx(nilCtx) error: %v", err)
	}
	if store2 == nil {
		t.Fatal("expected non-nil store")
	}

	mu.Lock()
	captured2 := readCtx
	mu.Unlock()

	if captured2 == nil || captured2 != context.Background() {
		t.Fatalf("expected context.Background(), got %v", captured2)
	}
}

func TestGetSessionStore_BackwardsCompatibility(t *testing.T) {
	var mu sync.Mutex
	var readCtx context.Context

	provider := &mockContextAwareProvider{
		onRead: func(ctx context.Context, _ string) {
			mu.Lock()
			readCtx = ctx
			mu.Unlock()
		},
	}
	manager := &Manager{
		provider: provider,
		config:   &ManagerConfig{},
	}

	testSid := "sid_compat_store"
	store, err := manager.GetSessionStore(testSid)
	if err != nil {
		t.Fatalf("GetSessionStore error: %v", err)
	}
	if store == nil || store.SessionID(context.Background()) != testSid {
		t.Fatalf("unexpected store result: %v", store)
	}

	mu.Lock()
	captured := readCtx
	mu.Unlock()

	if captured == nil || captured != context.Background() {
		t.Fatalf("expected context.Background(), got %v", captured)
	}
}

func TestReqContext_Fallback(t *testing.T) {
	// nil request variable
	var nilReq *http.Request
	ctx := reqContext(nilReq)
	if ctx != context.Background() {
		t.Fatalf("expected context.Background() for nil request, got %v", ctx)
	}

	// request with non-nil context
	customCtx := context.WithValue(context.Background(), testCtxKey, testCtxVal)
	req := httptest.NewRequest("GET", "/", http.NoBody).WithContext(customCtx)
	ctx2 := reqContext(req)
	if ctx2.Value(testCtxKey) != testCtxVal {
		t.Fatalf("expected custom context value, got %v", ctx2.Value(testCtxKey))
	}
}
