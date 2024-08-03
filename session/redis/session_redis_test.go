package redis

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/session"
)

func TestRedis(t *testing.T) {
	globalSession, err := setupSessionManager(t)
	if err != nil {
		t.Fatal(err)
	}

	go globalSession.GC()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSession.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}
	defer sess.SessionRelease(w)

	// SET AND GET
	err = sess.Set("username", "astaxie")
	if err != nil {
		t.Fatal("set username failed:", err)
	}
	username := sess.Get("username")
	if username != "astaxie" {
		t.Fatal("get username failed")
	}

	// DELETE
	err = sess.Delete("username")
	if err != nil {
		t.Fatal("delete username failed:", err)
	}
	username = sess.Get("username")
	if username != nil {
		t.Fatal("delete username failed")
	}

	// FLUSH
	err = sess.Set("username", "astaxie")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	err = sess.Set("password", "1qaz2wsx")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	username = sess.Get("username")
	if username != "astaxie" {
		t.Fatal("get username failed")
	}
	password := sess.Get("password")
	if password != "1qaz2wsx" {
		t.Fatal("get password failed")
	}
	err = sess.Flush()
	if err != nil {
		t.Fatal("flush failed:", err)
	}
	username = sess.Get("username")
	if username != nil {
		t.Fatal("flush failed")
	}
	password = sess.Get("password")
	if password != nil {
		t.Fatal("flush failed")
	}

	sess.SessionRelease(w)
}

func TestProvider_SessionInit(t *testing.T) {
	savePath := `
{ "save_path": "my save path", "idle_timeout": "3s"}
`
	cp := &Provider{}
	cp.SessionInit(12, savePath)
	assert.Equal(t, int64(12), cp.maxlifetime)
}

func TestStoreSessionReleaseIfPresentAndSessionDestroy(t *testing.T) {
	globalSessions, err := setupSessionManager(t)
	if err != nil {
		t.Fatal(err)
	}
	// todo test if e==nil
	go globalSessions.GC()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}

	if err := globalSessions.GetProvider().SessionDestroy(sess.SessionID()); err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sess.SessionReleaseIfPresent(httptest.NewRecorder())
	}()
	wg.Wait()

	if globalSessions.GetProvider().SessionExist(sess.SessionID()) {
		t.Fatalf("session %s should exist", sess.SessionID())
	}
}

func setupSessionManager(t *testing.T) (*session.Manager, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	redisConfig := fmt.Sprintf("%s,100,,0,30", redisAddr)
	sessionConfig := &session.ManagerConfig{}
	sessionConfig.CookieName = "gosessionid"
	sessionConfig.EnableSetCookie = true
	sessionConfig.Gclifetime = 3600
	sessionConfig.Maxlifetime = 3600
	sessionConfig.Secure = false
	sessionConfig.CookieLifeTime = 3600
	sessionConfig.ProviderConfig = redisConfig
	globalSessions, err := session.NewManager("redis", sessionConfig)
	if err != nil {
		t.Log("could not create manager: ", err)
		return nil, err
	}
	return globalSessions, nil
}
