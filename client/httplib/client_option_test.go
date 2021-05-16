// Copyright 2020 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httplib

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type respCarrier struct {
	Resp  *http.Response
	bytes []byte
}

func (r *respCarrier) SetHttpResponse(resp *http.Response) {
	r.Resp = resp
}

func (r *respCarrier) SetBytes(bytes []byte) {
	r.bytes = bytes
}

func (r *respCarrier) Bytes() []byte {
	return r.bytes
}

func (r *respCarrier) String() string {
	return string(r.bytes)
}

func TestOption_WithEnableCookie(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/",
		WithEnableCookie(true))
	if err != nil {
		t.Fatal(err)
	}

	v := "smallfish"
	var resp = &respCarrier{}
	err = client.Get(resp, "/cookies/set?k1="+v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	err = client.Get(resp, "/cookies")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	n := strings.Index(resp.String(), v)
	if n == -1 {
		t.Fatal(v + " not found in cookie")
	}
}

func TestOption_WithUserAgent(t *testing.T) {
	v := "beego"
	client, err := NewClient("test", "http://httpbin.org/",
		WithUserAgent(v))
	if err != nil {
		t.Fatal(err)
	}

	var resp = &respCarrier{}
	err = client.Get(resp, "/headers")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	n := strings.Index(resp.String(), v)
	if n == -1 {
		t.Fatal(v + " not found in user-agent")
	}
}

func TestOption_WithCheckRedirect(t *testing.T) {
	client, err := NewClient("test", "https://goolnk.com/33BD2j",
		WithCheckRedirect(func(redirectReq *http.Request, redirectVia []*http.Request) error {
			return errors.New("Redirect triggered")
		}))
	if err != nil {
		t.Fatal(err)
	}
	err = client.Get(nil, "")
	assert.NotNil(t, err)
}

func TestOption_WithHTTPSetting(t *testing.T) {
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

	client, err := NewClient("test", "http://httpbin.org/",
		WithHTTPSetting(setting))
	if err != nil {
		t.Fatal(err)
	}

	var resp = &respCarrier{}
	err = client.Get(resp, "/get")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	n := strings.Index(resp.String(), v)
	if n == -1 {
		t.Fatal(v + " not found in user-agent")
	}
}

func TestOption_WithHeader(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}
	client.CommonOpts = append(client.CommonOpts, WithHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36"))

	var resp = &respCarrier{}
	err = client.Get(resp, "/headers")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	n := strings.Index(resp.String(), "Mozilla/5.0")
	if n == -1 {
		t.Fatal("Mozilla/5.0 not found in user-agent")
	}
}

func TestOption_WithTokenFactory(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}
	client.CommonOpts = append(client.CommonOpts,
		WithTokenFactory(func() string {
			return "testauth"
		}))

	var resp = &respCarrier{}
	err = client.Get(resp, "/headers")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	n := strings.Index(resp.String(), "testauth")
	if n == -1 {
		t.Fatal("Auth is not set in request")
	}
}

func TestOption_WithBasicAuth(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	var resp = &respCarrier{}
	err = client.Get(resp, "/basic-auth/user/passwd",
		WithBasicAuth(func() (string, string) {
			return "user", "passwd"
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())
	n := strings.Index(resp.String(), "authenticated")
	if n == -1 {
		t.Fatal("authenticated not found in response")
	}
}

func TestOption_WithContentType(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	v := "application/json"
	var resp = &respCarrier{}
	err = client.Get(resp, "/headers", WithContentType(v))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	n := strings.Index(resp.String(), v)
	if n == -1 {
		t.Fatal(v + " not found in header")
	}
}

func TestOption_WithParam(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	v := "smallfish"
	var resp = &respCarrier{}
	err = client.Get(resp, "/get", WithParam("username", v))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	n := strings.Index(resp.String(), v)
	if n == -1 {
		t.Fatal(v + " not found in header")
	}
}

func TestOption_WithRetry(t *testing.T) {
	client, err := NewClient("test", "https://goolnk.com/33BD2j",
		WithCheckRedirect(func(redirectReq *http.Request, redirectVia []*http.Request) error {
			return errors.New("Redirect triggered")
		}))
	if err != nil {
		t.Fatal(err)
	}

	retryAmount := 1
	retryDelay := 1400 * time.Millisecond
	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	_ = client.Get(nil, "", WithRetry(retryAmount, retryDelay))

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	elapsedTime := endTime - startTime
	delayedTime := int64(retryAmount) * retryDelay.Milliseconds()
	if elapsedTime < delayedTime {
		t.Errorf("Not enough retries. Took %dms. Delay was meant to take %dms", elapsedTime, delayedTime)
	}
}
