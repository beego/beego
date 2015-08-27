package jwt

import (
	"github.com/astaxie/beego"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()

	return recorder, request
}

func Test_IssueTokenAction(t *testing.T) {
	url := "/v1/jwt/issue-token"

	mux := beego.NewControllerRegister()

	mux.InsertFilter("*", beego.BeforeRouter, AuthRequest(&Options{
		PrivateKeyPath: "test/jwt.rsa",
		PublicKeyPath:  "test/jwt.rsa.pub",
		WhiteList:      []string{"/v1/jwt/issue-token", "/docs"},
	}))

	mux.Add("/v1/jwt/issue-token", &JwtController{}, "get:IssueToken")

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if rw.Code != http.StatusOK {
		t.Errorf("Shoud return 200")
	}
}

func (tc *JwtController) Foo() {
	tc.Ctx.Output.Body([]byte("ok"))
}

func Test_AuthRequestWithAuthorizationHeader(t *testing.T) {

	url := "/foo"

	mux := beego.NewControllerRegister()

	mux.InsertFilter("*", beego.BeforeRouter, AuthRequest(&Options{
		PrivateKeyPath: "test/jwt.rsa",
		PublicKeyPath:  "test/jwt.rsa.pub",
		WhiteList:      []string{"/v1/jwt/issue-token", "/docs"},
	}))

	mux.Add("/foo", &JwtController{}, "get:Foo")
	newToken := CreateToken()

	rw, r := testRequest("GET", url)
	r.Header.Add("Authorization", "Bearer "+newToken["token"])
	mux.ServeHTTP(rw, r)

	if rw.Code != http.StatusOK {
		t.Errorf("Shoud return 200")
	}
	if rw.Body.String() != "ok" {
		t.Errorf("Should output ok")
	}
}

func Test_AuthRequestWithoutAuthorizationHeader(t *testing.T) {
	url := "/foo"

	mux := beego.NewControllerRegister()

	mux.InsertFilter("*", beego.BeforeRouter, AuthRequest(&Options{
		PrivateKeyPath: "test/jwt.rsa",
		PublicKeyPath:  "test/jwt.rsa.pub",
		WhiteList:      []string{"/v1/jwt/issue-token", "/docs"},
	}))

	mux.Add("/foo", &JwtController{}, "get:Foo")

	rw, r := testRequest("GET", url)
	mux.ServeHTTP(rw, r)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("Shoud return 401")
	}
}
