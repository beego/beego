package redis_sentinel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/server/web/mock/common"
	"github.com/beego/beego/v2/server/web/session"
)

func TestRedisSentinel(t *testing.T) {
	sessionConfig := session.NewManagerConfig(
		session.CfgCookieName(`gosessionid`),
		session.CfgSetCookie(true),
		session.CfgGcLifeTime(3600),
		session.CfgMaxLifeTime(3600),
		session.CfgSecure(false),
		session.CfgCookieLifeTime(3600),
		session.CfgProviderConfig("127.0.0.1:6379,100,,0,master"),
	)
	globalSessions, e := session.NewManager("redis_sentinel", sessionConfig)
	if e != nil {
		t.Log(e)
		return
	}
	// todo test if e==nil
	go globalSessions.GC()

	ctx := context.Background()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}
	defer sess.SessionRelease(ctx, w)

	// SET AND GET
	err = sess.Set(ctx, "username", "astaxie")
	if err != nil {
		t.Fatal("set username failed:", err)
	}
	username := sess.Get(ctx, "username")
	if username != "astaxie" {
		t.Fatal("get username failed")
	}

	// DELETE
	err = sess.Delete(ctx, "username")
	if err != nil {
		t.Fatal("delete username failed:", err)
	}
	username = sess.Get(ctx, "username")
	if username != nil {
		t.Fatal("delete username failed")
	}

	// FLUSH
	err = sess.Set(ctx, "username", "astaxie")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	err = sess.Set(ctx, "password", "1qaz2wsx")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	username = sess.Get(ctx, "username")
	if username != "astaxie" {
		t.Fatal("get username failed")
	}
	password := sess.Get(ctx, "password")
	if password != "1qaz2wsx" {
		t.Fatal("get password failed")
	}
	err = sess.Flush(ctx)
	if err != nil {
		t.Fatal("flush failed:", err)
	}
	username = sess.Get(ctx, "username")
	if username != nil {
		t.Fatal("flush failed")
	}
	password = sess.Get(ctx, "password")
	if password != nil {
		t.Fatal("flush failed")
	}

	sess.SessionRelease(ctx, w)
}

func TestProvider_SessionInit(t *testing.T) {
	savePath := `
{ "save_path": "my save path", "idle_timeout": "3s"}
`
	cp := &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "my save path", cp.SavePath)
	assert.Equal(t, 3*time.Second, cp.idleTimeout)
	assert.Equal(t, int64(12), cp.maxlifetime)
}

func TestStoreSessionReleaseIfPresentAndSessionDestroy(t *testing.T) {
	sessionConfig := session.NewManagerConfig(
		session.CfgCookieName(`gosessionid`),
		session.CfgSetCookie(true),
		session.CfgGcLifeTime(3600),
		session.CfgMaxLifeTime(3600),
		session.CfgSecure(false),
		session.CfgCookieLifeTime(3600),
		session.CfgProviderConfig("127.0.0.1:6379,100,,0,master"),
	)
	globalSessions, e := session.NewManager("redis_sentinel", sessionConfig)
	if e != nil {
		t.Log(e)
		return
	}
	// todo test if e==nil
	go globalSessions.GC()

	ctx := context.Background()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}

	if err := globalSessions.GetProvider().SessionDestroy(ctx, sess.SessionID(ctx)); err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sess.SessionReleaseIfPresent(context.Background(), common.NewMockHttpResponse())
	}()
	wg.Wait()
	exist, err := globalSessions.GetProvider().SessionExist(ctx, sess.SessionID(ctx))
	if err != nil {
		t.Error(err)
	}
	if exist {
		t.Fatalf("session %s should exist", sess.SessionID(ctx))
	}
}
