package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func testRequest(t *testing.T, handler *web.ControllerRegister, requestIP, method, path string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	r.Header.Set("X-Real-Ip", requestIP)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", requestIP, method, path, w.Code, code)
	}
}

func TestLimiter(t *testing.T) {
	handler := web.NewControllerRegister()
	err := handler.InsertFilter("/foo/*", web.BeforeRouter, NewLimiter(WithRate(1*time.Millisecond), WithCapacity(1), WithSessionKey(RemoteIPSessionKey)))
	if err != nil {
		t.Error(err)
	}
	handler.Any("*", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
	})

	route := "/foo/1"
	ip := "127.0.0.1"
	testRequest(t, handler, ip, "GET", route, 200)
	testRequest(t, handler, ip, "GET", route, 429)
	testRequest(t, handler, "127.0.0.2", "GET", route, 200)
	time.Sleep(1 * time.Millisecond)
	testRequest(t, handler, ip, "GET", route, 200)
}

func BenchmarkWithoutLimiter(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := web.NewControllerRegister()
	web.BConfig.RunMode = web.PROD
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			handler.ServeHTTP(recorder, r)
		}
	})
}

func BenchmarkWithLimiter(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := web.NewControllerRegister()
	web.BConfig.RunMode = web.PROD
	err := handler.InsertFilter("*", web.BeforeRouter, NewLimiter(WithRate(1*time.Millisecond), WithCapacity(100)))
	if err != nil {
		b.Error(err)
	}
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			handler.ServeHTTP(recorder, r)
		}
	})
}
