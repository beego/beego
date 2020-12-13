// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package authz provides handlers to enable ACL, RBAC, ABAC authorization support.
// Simple Usage:
//	import(
//		"github.com/beego/beego"
//		"github.com/beego/beego/plugins/authz"
//		"github.com/casbin/casbin"
//	)
//
//	func main(){
//		// mediate the access for every request
//		beego.InsertFilter("*", beego.BeforeRouter, authz.NewAuthorizer(casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")))
//		beego.Run()
//	}
//
//
// Advanced Usage:
//
//	func main(){
//		e := casbin.NewEnforcer("authz_model.conf", "")
//		e.AddRoleForUser("alice", "admin")
//		e.AddPolicy(...)
//
//		beego.InsertFilter("*", beego.BeforeRouter, authz.NewAuthorizer(e))
//		beego.Run()
//	}
package authz

import (
	"net/http"

	"github.com/casbin/casbin"

	beego "github.com/beego/beego/adapter"
	"github.com/beego/beego/adapter/context"
	beecontext "github.com/beego/beego/server/web/context"
	"github.com/beego/beego/server/web/filter/authz"
)

// NewAuthorizer returns the authorizer.
// Use a casbin enforcer as input
func NewAuthorizer(e *casbin.Enforcer) beego.FilterFunc {
	f := authz.NewAuthorizer(e)
	return func(context *context.Context) {
		f((*beecontext.Context)(context))
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer authz.BasicAuthorizer

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func (a *BasicAuthorizer) GetUserName(r *http.Request) string {
	return (*authz.BasicAuthorizer)(a).GetUserName(r)
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request) bool {
	return (*authz.BasicAuthorizer)(a).CheckPermission(r)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(w http.ResponseWriter) {
	(*authz.BasicAuthorizer)(a).RequirePermission(w)
}
