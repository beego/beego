// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package auth

// Example:
// func SecretAuth(username, password string) bool {
// 	if username == "astaxie" && password == "helloBeego" {
// 		return true
// 	}
// 	return false
// }
// authPlugin := auth.NewBasicAuthenticator(SecretAuth, "My Realm")
// beego.InsertFilter("*", beego.BeforeRouter,authPlugin)

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func NewBasicAuthenticator(secrets SecretProvider, Realm string) beego.FilterFunc {
	return func(ctx *context.Context) {
		a := &BasicAuth{Secrets: secrets, Realm: Realm}
		if username := a.CheckAuth(ctx.Request); username == "" {
			a.RequireAuth(ctx.ResponseWriter, ctx.Request)
		}
	}
}

type SecretProvider func(user, pass string) bool

type BasicAuth struct {
	Secrets SecretProvider
	Realm   string
}

/*
 Checks the username/password combination from the request. Returns
 either an empty string (authentication failed) or the name of the
 authenticated user.

 Supports MD5 and SHA1 password entries
*/
func (a *BasicAuth) CheckAuth(r *http.Request) string {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return ""
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return ""
	}

	if a.Secrets(pair[0], pair[1]) {
		return pair[0]
	}
	return ""
}

/*
 http.Handler for BasicAuth which initiates the authentication process
 (or requires reauthentication).
*/
func (a *BasicAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+a.Realm+`"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}
