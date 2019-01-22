package redis_sentinel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego/session"
)

func TestRedisSentinel(t *testing.T) {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "127.0.0.1:6379,100,,0,master",
	}
	globalSessions, e := session.NewManager("redis_sentinel", sessionConfig)
	if e != nil {
		t.Log(e)
		return
	}
	//todo test if e==nil
	go globalSessions.GC()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
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
