package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web"
	webContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/google/uuid"
)

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
