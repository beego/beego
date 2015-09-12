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

// Package jwt provides JWT (Json Web Token) authentication
//
// Usage
// In file main.go
//
//	import (
// 		"github.com/astaxie/beego"
//		"github.com/astaxie/beego/plugins/jwt"
// )
//
//	func main() {
//		// JWT for Url matching /v1/*
//		// PrivateKeyPath: The path for the private RSA key used by JWT
//		// PublicKeyPath: The path for the public RSA key used by JWT
//		// The list of Urls should be excluded from the JWT Auth
//		beego.InsertFilter("/v1/*", beego.BeforeRouter, jwt.AuthRequest(&jwt.Options{
//			PrivateKeyPath: "conf/beeblog.rsa",
//			PublicKeyPath:  "conf/beeblog.rsa.pub",
//			WhiteList:      []string{"/v1/jwt/issue-token", "/docs"},
//		}))
//		beego.Run()
//	}
//
// In file routers/router.go
//
//	import (
// 		"github.com/astaxie/beego"
//		"github.com/astaxie/beego/plugins/jwt"
// )
//	func init() {
//		ns := beego.NSNamespace("/jwt",
//			beego.NSInclude(
//				&jwt.JwtController{},
//			),
//		)
//		beego.AddNamespace(ns)
//	}
//

package jwt

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	goJwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"time"
)

// Options for the JWT Auth
type Options struct {
	PrivateKeyPath string
	PublicKeyPath  string
	WhiteList      []string
}

var RSAKeys struct {
	PrivateKey []byte
	PublicKey  []byte
}

func AuthRequest(o *Options) beego.FilterFunc {
	RSAKeys.PrivateKey, _ = ioutil.ReadFile(o.PrivateKeyPath)
	RSAKeys.PublicKey, _ = ioutil.ReadFile(o.PublicKeyPath)

	return func(ctx *context.Context) {
		// :TODO the url patterns should be considered here.
		// Shouldn't only use the string equal
		for _, method := range o.WhiteList {
			if method == ctx.Request.URL.Path {
				return
			}
		}

		parsedToken, err := goJwt.ParseFromRequest(ctx.Request, func(t *goJwt.Token) (interface{}, error) {
			return RSAKeys.PublicKey, nil
		})

		if err == nil && parsedToken.Valid {
			ctx.Output.SetStatus(http.StatusOK)
		} else {
			ctx.Output.SetStatus(http.StatusUnauthorized)
		}

	}
}

// oprations for Jwt
type JwtController struct {
	beego.Controller
}

func (this *JwtController) URLMapping() {
	this.Mapping("IssueToken", this.IssueToken)
}

// @Title IssueToken
// @Description Issue a Json Web Token
// @Success 200 string
// @Failure 403 no privilege to access
// @Failure 500 server inner error
// @router /issue-token [get]
func (this *JwtController) IssueToken() {
	this.Data["json"] = CreateToken()
	this.ServeJson()
}

func CreateToken() map[string]string {
	log := logs.NewLogger(10000)
	log.SetLogger("console", "")

	token := goJwt.New(goJwt.GetSigningMethod("RS256")) // Create a Token that will be signed with RSA 256.
	token.Claims["ID"] = "This is my super fake ID"
	token.Claims["exp"] = time.Now().Unix() + 36000
	// The claims object allows you to store information in the actual token.
	tokenString, _ := token.SignedString(RSAKeys.PrivateKey)
	// tokenString Contains the actual token you should share with your client.
	return map[string]string{"token": tokenString}
}
