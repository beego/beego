package redis_sentinel

import (
	"github.com/astaxie/beego/session"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedisSentinel(t *testing.T) {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "119.23.132.234:26379,100,,0,master",
	}
	globalSessions, _ := session.NewManager("redis_sentinel", sessionConfig)
	go globalSessions.GC()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
	assert.Nil(t, err)
	defer sess.SessionRelease(w)

	// SET AND GET
	err = sess.Set("username", "astaxie")
	assert.Nil(t, err)
	username := sess.Get("username")
	assert.Equal(t, "astaxie", username)

	// DELETE
	err = sess.Delete("username")
	assert.Nil(t, err)
	username = sess.Get("username")
	assert.Equal(t, nil, username)

	// FLUSH
	err = sess.Set("username", "astaxie")
	assert.Nil(t, err)
	err = sess.Set("password", "1qaz2wsx")
	assert.Nil(t, err)
	username = sess.Get("username")
	assert.Equal(t, "astaxie", username)
	password := sess.Get("password")
	assert.Equal(t, "1qaz2wsx", password)
	err = sess.Flush()
	assert.Nil(t, err)
	username = sess.Get("username")
	assert.Equal(t, nil, username)
	password = sess.Get("password")
	assert.Equal(t, nil, password)

	sess.SessionRelease(w)

}
