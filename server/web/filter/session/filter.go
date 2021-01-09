package session

import (
	"context"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	webContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

//Session maintain session for web service
//Session new a session storage and store it into webContext.Context
//
//params:
//ctx: pointer of beego web context
//storeKey: set the storage key in ctx.Input
//
//if you want to get session storage, just see GetStore
func Session(providerType session.ProviderType, options ...session.ManagerConfigOpt) web.FilterChain {
	sessionConfig := session.NewManagerConfig(options...)
	sessionManager, _ := session.NewManager(string(providerType), sessionConfig)
	go sessionManager.GC()

	return func(next web.FilterFunc) web.FilterFunc {
		return func(ctx *webContext.Context) {
			if ctx.Input.CruSession != nil {
				return
			}

			if sess, err := sessionManager.SessionStart(ctx.ResponseWriter, ctx.Request); err != nil {
				logs.Warning(`init session error:%s`, err.Error())
			} else {
				//release session at the end of request
				defer sess.SessionRelease(context.Background(), ctx.ResponseWriter)
				ctx.Input.CruSession = sess
			}

			next(ctx)

		}
	}
}

//GetStore get session storage in beego web context
func GetStore(ctx *webContext.Context) (store session.Store, err error) {
	if ctx == nil {
		err = errors.New(`ctx is nil`)
		return
	}

	if s := ctx.Input.CruSession; s != nil {
		store = s
		return
	} else {
		err = errors.New(`can not get a valid session store`)
		return
	}
}
