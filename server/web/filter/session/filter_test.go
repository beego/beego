package session

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/beego/beego/v2/server/web"
	webContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

var errProviderReconfigured = errors.New("test session provider is already initialized with a different configuration")

type reinitializationProvider struct {
	session.Provider
	savePath string
}

func (p *reinitializationProvider) SessionInit(_ context.Context, _ int64, savePath string) error {
	if p.savePath == "" {
		p.savePath = savePath
		return nil
	}
	if p.savePath != savePath {
		return errProviderReconfigured
	}
	return nil
}

func testRequest(t *testing.T, handler *web.ControllerRegister, path string, method string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != code {
		t.Errorf("%s, %s: %d, supposed to be %d", path, method, w.Code, code)
	}
}

func TestSession(t *testing.T) {
	storeKey := uuid.New().String()
	handler := web.NewControllerRegister()
	handler.InsertFilterChain(
		"*",
		Session(
			session.ProviderMemory,
			session.CfgCookieName(`go_session_id`),
			session.CfgSetCookie(true),
			session.CfgGcLifeTime(3600),
			session.CfgMaxLifeTime(3600),
			session.CfgSecure(false),
			session.CfgCookieLifeTime(3600),
		),
	)
	handler.InsertFilterChain(
		"*",
		func(next web.FilterFunc) web.FilterFunc {
			return func(ctx *webContext.Context) {
				if store := ctx.Input.GetData(storeKey); store == nil {
					t.Error(`store should not be nil`)
				}
				next(ctx)
			}
		},
	)
	handler.Any("*", func(ctx *webContext.Context) {
		ctx.Output.SetStatus(200)
	})

	testRequest(t, handler, "/dataset1/resource1", "GET", 200)
}

func TestSession1(t *testing.T) {
	handler := web.NewControllerRegister()
	handler.InsertFilterChain(
		"*",
		Session(
			session.ProviderMemory,
			session.CfgCookieName(`go_session_id`),
			session.CfgSetCookie(true),
			session.CfgGcLifeTime(3600),
			session.CfgMaxLifeTime(3600),
			session.CfgSecure(false),
			session.CfgCookieLifeTime(3600),
		),
	)
	handler.InsertFilterChain(
		"*",
		func(next web.FilterFunc) web.FilterFunc {
			return func(ctx *webContext.Context) {
				if store, err := ctx.Session(); store == nil || err != nil {
					t.Error(`store should not be nil`)
				}
				next(ctx)
			}
		},
	)
	handler.Any("*", func(ctx *webContext.Context) {
		ctx.Output.SetStatus(200)
	})

	testRequest(t, handler, "/dataset1/resource1", "GET", 200)
}

func TestSessionPanicsSynchronouslyWhenProviderInitializationFails(t *testing.T) {
	providerType := session.ProviderType("reinitialization-test-" + uuid.New().String())
	provider := &reinitializationProvider{}
	const initialSavePath = "user=test password=initial-secret dbname=first"
	if err := provider.SessionInit(context.Background(), 3600, initialSavePath); err != nil {
		t.Fatalf("initialize test provider: %v", err)
	}
	session.Register(string(providerType), provider)

	const conflictingSavePath = "user=test password=conflicting-secret dbname=second"
	returned := false
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		Session(
			providerType,
			session.CfgGcLifeTime(3600),
			session.CfgMaxLifeTime(3600),
			session.CfgProviderConfig(conflictingSavePath),
		)
		returned = true
	}()

	if returned {
		t.Fatal("Session returned after provider initialization failed")
	}
	if recovered != errProviderReconfigured {
		t.Fatalf("panic value = %v; want provider initialization error", recovered)
	}
	message := errProviderReconfigured.Error()
	if strings.Contains(message, conflictingSavePath) || strings.Contains(message, "conflicting-secret") {
		t.Fatalf("panic exposed provider credentials: %q", message)
	}
}
