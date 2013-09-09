package middleware

import (
	"github.com/astaxie/beego/session"
)

var (
	GlobalSessions *session.Manager
)

func StartSession(provideName, cookieName string, maxlifetime int64, savePath string) {
	GlobalSessions, _ = session.NewManager(provideName, cookieName, maxlifetime, savePath)
	go GlobalSessions.GC()
}
