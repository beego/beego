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

// Package apiauth provides handlers to enable apiauth support.
//
// Simple Usage:
//	import(
//		"github.com/astaxie/beego"
//		"github.com/astaxie/beego/plugins/apiauth"
//	)
//
//	func main(){
//		// apiauth every request
//		beego.InsertFilter("*", beego.BeforeRouter,apiauth.APIBaiscAuth("appid","appkey"))
//		beego.Run()
//	}
//
// Advanced Usage:
//
//	func getAppSecret(appid string) string {
//		// get appsecret by appid
//		// maybe store in configure, maybe in database
//	}
//
//	beego.InsertFilter("*", beego.BeforeRouter,apiauth.APISecretAuth(getAppSecret, 360))
//
// Information:
//
// In the request user should include these params in the query
//
// 1. appid
//
//		 appid is assigned to the application
//
// 2. signature
//
//	get the signature use apiauth.Signature()
//
//	when you send to server remember use url.QueryEscape()
//
// 3. timestamp:
//
//       send the request time, the format is yyyy-mm-dd HH:ii:ss
//
package apiauth

import (
	"net/url"

	beego "github.com/astaxie/beego/adapter"
	"github.com/astaxie/beego/adapter/context"
	beecontext "github.com/astaxie/beego/server/web/context"
	"github.com/astaxie/beego/server/web/filter/apiauth"
)

// AppIDToAppSecret is used to get appsecret throw appid
type AppIDToAppSecret apiauth.AppIDToAppSecret

// APIBasicAuth use the basic appid/appkey as the AppIdToAppSecret
func APIBasicAuth(appid, appkey string) beego.FilterFunc {
	f := apiauth.APIBasicAuth(appid, appkey)
	return func(c *context.Context) {
		f((*beecontext.Context)(c))
	}
}

// APIBaiscAuth calls APIBasicAuth for previous callers
func APIBaiscAuth(appid, appkey string) beego.FilterFunc {
	return APIBasicAuth(appid, appkey)
}

// APISecretAuth use AppIdToAppSecret verify and
func APISecretAuth(f AppIDToAppSecret, timeout int) beego.FilterFunc {
	ft := apiauth.APISecretAuth(apiauth.AppIDToAppSecret(f), timeout)
	return func(ctx *context.Context) {
		ft((*beecontext.Context)(ctx))
	}
}

// Signature used to generate signature with the appsecret/method/params/RequestURI
func Signature(appsecret, method string, params url.Values, requestURL string) string {
	return apiauth.Signature(appsecret, method, params, requestURL)
}
