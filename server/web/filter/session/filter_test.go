package session

import (
	"context"
	"github.com/beego/beego/v2/server/web"
	webContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
	"net/http"
	"net/http/httptest"
	"testing"
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
				if store := ctx.Input.GetData(getSessionKey()); store == nil {
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

func TestGetStore(t *testing.T) {
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
				var (
					checkKey   = `asodiuasdk1j)AS(87`
					checkValue = `ASsd-09812-3`

					store session.Store
					err   error

					c = context.Background()
				)

				if store, err = GetStore(ctx); err == nil {
					if store == nil {
						t.Error(`store should not be nil`)
					} else {
						_ = store.Set(c, checkKey, checkValue)
					}
				} else {
					t.Error(err)
				}

				next(ctx)

				if store != nil {
					if v := store.Get(c, checkKey); v != checkValue {
						t.Error(v, `is not equals to`, checkValue)
					}
				}else{
					t.Error(`store should not be nil`)
				}

			}
		},
	)
	handler.Any("*", func(ctx *webContext.Context) {
		ctx.Output.SetStatus(200)
	})

	testRequest(t, handler, "/dataset1/resource1", "GET", 200)
}
