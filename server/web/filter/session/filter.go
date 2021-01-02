package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	webContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/google/uuid"
	"sync"
)

var (
	sessionKey     string
	sessionKeyOnce sync.Once
)

func getSessionKey() string {

	sessionKeyOnce.Do(func() {
		//generate an unique session store key
		sessionKey = fmt.Sprintf(`sess_store:%d`, uuid.New().ID())
	})

	return sessionKey
}

func Session(providerType session.ProviderType, options ...session.ManagerConfigOpt) web.FilterChain {
	sessionConfig := session.NewManagerConfig(options...)
	sessionManager, _ := session.NewManager(string(providerType), sessionConfig)
	go sessionManager.GC()

	return func(next web.FilterFunc) web.FilterFunc {
		return func(ctx *webContext.Context) {

			if sess, err := sessionManager.SessionStart(ctx.ResponseWriter, ctx.Request); err != nil {
				logs.Warning(`init session error:%s`, err.Error())
			} else {
				//release session at the end of request
				defer sess.SessionRelease(context.Background(), ctx.ResponseWriter)
				ctx.Input.SetData(getSessionKey(), sess)
			}

			next(ctx)

		}
	}
}

func GetStore(ctx *webContext.Context) (store session.Store, err error) {
	if ctx == nil {
		err = errors.New(`ctx is nil`)
		return
	}

	if s, ok := ctx.Input.GetData(getSessionKey()).(session.Store); ok {
		store = s
		return
	} else {
		err = errors.New(`can not get a valid session store`)
		return
	}
}
