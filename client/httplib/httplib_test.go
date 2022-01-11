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

package httplib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	req := Get("http://httpbin.org/get")
	resp, err := req.Response()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestDoRequest(t *testing.T) {
	req := Get("https://goolnk.com/33BD2j")
	retryAmount := 1
	req.Retries(1)
	req.RetryDelay(1400 * time.Millisecond)
	retryDelay := 1400 * time.Millisecond

	req.setting.CheckRedirect = func(redirectReq *http.Request, redirectVia []*http.Request) error {
		return errors.New("Redirect triggered")
	}

	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	_, err := req.Response()
	if err == nil {
		t.Fatal("Response should have yielded an error")
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	elapsedTime := endTime - startTime
	delayedTime := int64(retryAmount) * retryDelay.Milliseconds()

	if elapsedTime < delayedTime {
		t.Errorf("Not enough retries. Took %dms. Delay was meant to take %dms", elapsedTime, delayedTime)
	}
}

func TestGet(t *testing.T) {
	req := Get("http://httpbin.org/get")
	b, err := req.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)

	s, err := req.String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)

	if string(b) != s {
		t.Fatal("request data not match")
	}
}

func TestSimplePost(t *testing.T) {
	v := "smallfish"
	req := Post("http://httpbin.org/post")
	req.Param("username", v)

	str, err := req.String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)

	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in post")
	}
}

// func TestPostFile(t *testing.T) {
//	v := "smallfish"
//	req := Post("http://httpbin.org/post")
//	req.Debug(true)
//	req.Param("username", v)
//	req.PostFile("uploadfile", "httplib_test.go")

//	str, err := req.String()
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(str)

//	n := strings.Index(str, v)
//	if n == -1 {
//		t.Fatal(v + " not found in post")
//	}
// }

func TestSimplePut(t *testing.T) {
	str, err := Put("http://httpbin.org/put").String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
}

func TestSimpleDelete(t *testing.T) {
	str, err := Delete("http://httpbin.org/delete").String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
}

func TestSimpleDeleteParam(t *testing.T) {
	str, err := Delete("http://httpbin.org/delete").Param("key", "val").String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
}

func TestWithCookie(t *testing.T) {
	v := "smallfish"
	str, err := Get("http://httpbin.org/cookies/set?k1=" + v).SetEnableCookie(true).String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)

	str, err = Get("http://httpbin.org/cookies").SetEnableCookie(true).String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)

	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in cookie")
	}
}

func TestWithBasicAuth(t *testing.T) {
	str, err := Get("http://httpbin.org/basic-auth/user/passwd").SetBasicAuth("user", "passwd").String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
	n := strings.Index(str, "authenticated")
	if n == -1 {
		t.Fatal("authenticated not found in response")
	}
}

func TestWithUserAgent(t *testing.T) {
	v := "beego"
	str, err := Get("http://httpbin.org/headers").SetUserAgent(v).String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)

	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in user-agent")
	}
}

func TestWithSetting(t *testing.T) {
	v := "beego"
	var setting BeegoHTTPSettings
	setting.EnableCookie = true
	setting.UserAgent = v
	setting.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          50,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	setting.ReadWriteTimeout = 5 * time.Second
	SetDefaultSetting(setting)

	str, err := Get("http://httpbin.org/get").String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)

	n := strings.Index(str, v)
	if n == -1 {
		t.Fatal(v + " not found in user-agent")
	}
}

func TestToJson(t *testing.T) {
	req := Get("http://httpbin.org/ip")
	resp, err := req.Response()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)

	// httpbin will return http remote addr
	type IP struct {
		Origin string `json:"origin"`
	}
	var ip IP
	err = req.ToJSON(&ip)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip.Origin)
	ips := strings.Split(ip.Origin, ",")
	if len(ips) == 0 {
		t.Fatal("response is not valid ip")
	}
	for i := range ips {
		if net.ParseIP(strings.TrimSpace(ips[i])).To4() == nil {
			t.Fatal("response is not valid ip")
		}
	}
}

func TestToFile(t *testing.T) {
	f := "beego_testfile"
	req := Get("http://httpbin.org/ip")
	err := req.ToFile(f)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f)
	b, err := ioutil.ReadFile(f)
	if n := bytes.Index(b, []byte("origin")); n == -1 {
		t.Fatal(err)
	}
}

func TestToFileDir(t *testing.T) {
	f := "./files/beego_testfile"
	req := Get("http://httpbin.org/ip")
	err := req.ToFile(f)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("./files")
	b, err := ioutil.ReadFile(f)
	if n := bytes.Index(b, []byte("origin")); n == -1 {
		t.Fatal(err)
	}
}

func TestHeader(t *testing.T) {
	req := Get("http://httpbin.org/headers")
	req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36")
	str, err := req.String()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(str)
}

// TestAddFilter make sure that AddFilters only work for the specific request
func TestAddFilter(t *testing.T) {
	req := Get("http://beego.vip")
	req.AddFilters(func(next Filter) Filter {
		return func(ctx context.Context, req *BeegoHTTPRequest) (*http.Response, error) {
			return next(ctx, req)
		}
	})

	r := Get("http://beego.vip")
	assert.Equal(t, 1, len(req.setting.FilterChains)-len(r.setting.FilterChains))
}

func TestFilterChainOrder(t *testing.T) {
	req := Get("http://beego.vip")
	req.AddFilters(func(next Filter) Filter {
		return func(ctx context.Context, req *BeegoHTTPRequest) (*http.Response, error) {
			return NewHttpResponseWithJsonBody("first"), nil
		}
	})

	req.AddFilters(func(next Filter) Filter {
		return func(ctx context.Context, req *BeegoHTTPRequest) (*http.Response, error) {
			return NewHttpResponseWithJsonBody("second"), nil
		}
	})

	resp, err := req.DoRequestWithCtx(context.Background())
	assert.Nil(t, err)
	data := make([]byte, 5)
	_, _ = resp.Body.Read(data)
	assert.Equal(t, "first", string(data))
}

func TestHead(t *testing.T) {
	req := Head("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "HEAD", req.req.Method)
}

func TestDelete(t *testing.T) {
	req := Delete("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "DELETE", req.req.Method)
}

func TestPost(t *testing.T) {
	req := Post("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.req.Method)
}

func TestNewBeegoRequest(t *testing.T) {
	req := NewBeegoRequest("http://beego.vip", "GET")
	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.req.Method)

	// invalid case but still go request
	req = NewBeegoRequest("httpa\ta://beego.vip", "GET")
	assert.NotNil(t, req)
}

func TestBeegoHTTPRequestSetProtocolVersion(t *testing.T) {
	req := NewBeegoRequest("http://beego.vip", "GET")
	req.SetProtocolVersion("HTTP/3.10")
	assert.Equal(t, "HTTP/3.10", req.req.Proto)
	assert.Equal(t, 3, req.req.ProtoMajor)
	assert.Equal(t, 10, req.req.ProtoMinor)

	req.SetProtocolVersion("")
	assert.Equal(t, "HTTP/1.1", req.req.Proto)
	assert.Equal(t, 1, req.req.ProtoMajor)
	assert.Equal(t, 1, req.req.ProtoMinor)

	// invalid case
	req.SetProtocolVersion("HTTP/aaa1.1")
	assert.Equal(t, "HTTP/1.1", req.req.Proto)
	assert.Equal(t, 1, req.req.ProtoMajor)
	assert.Equal(t, 1, req.req.ProtoMinor)
}

func TestPut(t *testing.T) {
	req := Put("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "PUT", req.req.Method)
}

func TestBeegoHTTPRequestHeader(t *testing.T) {
	req := Post("http://beego.vip")
	key, value := "test-header", "test-header-value"
	req.Header(key, value)
	assert.Equal(t, value, req.req.Header.Get(key))
}

func TestBeegoHTTPRequestSetHost(t *testing.T) {
	req := Post("http://beego.vip")
	host := "test-hose"
	req.SetHost(host)
	assert.Equal(t, host, req.req.Host)
}

func TestBeegoHTTPRequestParam(t *testing.T) {
	req := Post("http://beego.vip")
	key, value := "test-param", "test-param-value"
	req.Param(key, value)
	assert.Equal(t, value, req.params[key][0])

	value1 := "test-param-value-1"
	req.Param(key, value1)
	assert.Equal(t, value1, req.params[key][1])
}

func TestBeegoHTTPRequestBody(t *testing.T) {
	req := Post("http://beego.vip")
	body := `hello, world`
	req.Body([]byte(body))
	assert.Equal(t, int64(len(body)), req.req.ContentLength)
	assert.NotNil(t, req.req.GetBody)
	assert.NotNil(t, req.req.Body)

	body = "hhhh, I am test"
	req.Body(body)
	assert.Equal(t, int64(len(body)), req.req.ContentLength)
	assert.NotNil(t, req.req.GetBody)
	assert.NotNil(t, req.req.Body)

	// invalid case
	req.Body(13)
}

type user struct {
	Name string `xml:"name"`
}

func TestBeegoHTTPRequestXMLBody(t *testing.T) {
	req := Post("http://beego.vip")
	body := &user{
		Name: "Tom",
	}
	_, err := req.XMLBody(body)
	assert.True(t, req.req.ContentLength > 0)
	assert.Nil(t, err)
	assert.NotNil(t, req.req.GetBody)
}

// TODO
func TestBeegoHTTPRequestResponseForValue(t *testing.T) {
}

func TestBeegoHTTPRequestJSONMarshal(t *testing.T) {
	req := Post("http://beego.vip")
	req.SetEscapeHTML(false)
	body := map[string]interface{}{
		"escape": "left&right",
	}
	b, _ := req.JSONMarshal(body)
	assert.Equal(t, fmt.Sprintf(`{"escape":"left&right"}%s`, "\n"), string(b))
}
