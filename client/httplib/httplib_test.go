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
	json "encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HttplibTestSuite struct {
	suite.Suite
	l net.Listener
}

func (h *HttplibTestSuite) SetupSuite() {
	listener, err := net.Listen("tcp", ":8080")
	require.NoError(h.T(), err)
	h.l = listener

	handler := http.NewServeMux()

	handler.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		agent := request.Header.Get("User-Agent")
		_, _ = writer.Write([]byte("hello, " + agent))
	})

	handler.HandleFunc("/put", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello, put"))
	})

	handler.HandleFunc("/post", func(writer http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		_, _ = writer.Write(body)
	})

	handler.HandleFunc("/delete", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello, delete"))
	})

	handler.HandleFunc("/cookies/set", func(writer http.ResponseWriter, request *http.Request) {
		k1 := request.URL.Query().Get("k1")
		http.SetCookie(writer, &http.Cookie{
			Name:  "k1",
			Value: k1,
		})
		_, _ = writer.Write([]byte("hello, set cookie"))
	})

	handler.HandleFunc("/cookies", func(writer http.ResponseWriter, request *http.Request) {
		body := request.Cookies()[0].String()
		_, _ = writer.Write([]byte(body))
	})

	handler.HandleFunc("/basic-auth/user/passwd", func(writer http.ResponseWriter, request *http.Request) {
		_, _, ok := request.BasicAuth()
		if ok {
			_, _ = writer.Write([]byte("authenticated"))
		} else {
			_, _ = writer.Write([]byte("no auth"))
		}
	})

	handler.HandleFunc("/headers", func(writer http.ResponseWriter, request *http.Request) {
		agent := request.Header.Get("User-Agent")
		_, _ = writer.Write([]byte(agent))
	})

	handler.HandleFunc("/ip", func(writer http.ResponseWriter, request *http.Request) {
		data := map[string]string{"origin": "127.0.0.1"}
		jsonBytes, _ := json.Marshal(data)
		_, _ = writer.Write(jsonBytes)
	})

	handler.HandleFunc("/redirect", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "redirect_dst", http.StatusTemporaryRedirect)
	})
	handler.HandleFunc("redirect_dst", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello"))
	})
	go func() {
		_ = http.Serve(listener, handler)
	}()
}

func (h *HttplibTestSuite) TearDownSuite() {
	_ = h.l.Close()
}

func TestHttplib(t *testing.T) {
	suite.Run(t, &HttplibTestSuite{})
}

func (h *HttplibTestSuite) TestResponse() {
	req := Get("http://localhost:8080/get")
	_, err := req.Response()
	require.NoError(h.T(), err)
}

func (h *HttplibTestSuite) TestDoRequest() {
	t := h.T()
	req := Get("http://localhost:8080/redirect")
	retryAmount := 1
	req.Retries(1)
	req.RetryDelay(1400 * time.Millisecond)
	retryDelay := 1400 * time.Millisecond

	req.setting.CheckRedirect = func(redirectReq *http.Request, redirectVia []*http.Request) error {
		return errors.New("Redirect triggered")
	}

	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	_, err := req.Response()
	require.Error(t, err)

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	elapsedTime := endTime - startTime
	delayedTime := int64(retryAmount) * retryDelay.Milliseconds()

	if elapsedTime < delayedTime {
		t.Errorf("Not enough retries. Took %dms. Delay was meant to take %dms", elapsedTime, delayedTime)
	}
}

func (h *HttplibTestSuite) TestGet() {
	t := h.T()
	req := Get("http://localhost:8080/get")
	b, err := req.Bytes()
	require.NoError(t, err)

	s, err := req.String()
	require.NoError(t, err)
	require.Equal(t, string(b), s)
}

func (h *HttplibTestSuite) TestSimplePost() {
	t := h.T()
	v := "smallfish"
	req := Post("http://localhost:8080/post")
	req.Param("username", v)

	str, err := req.String()
	require.NoError(t, err)
	n := strings.Index(str, v)
	require.NotEqual(t, -1, n)
}

func (h *HttplibTestSuite) TestSimplePut() {
	t := h.T()
	_, err := Put("http://localhost:8080/put").String()
	require.NoError(t, err)
}

func (h *HttplibTestSuite) TestSimpleDelete() {
	t := h.T()
	_, err := Delete("http://localhost:8080/delete").String()
	require.NoError(t, err)
}

func (h *HttplibTestSuite) TestSimpleDeleteParam() {
	t := h.T()
	_, err := Delete("http://localhost:8080/delete").Param("key", "val").String()
	require.NoError(t, err)
}

func (h *HttplibTestSuite) TestWithCookie() {
	t := h.T()
	v := "smallfish"
	_, err := Get("http://localhost:8080/cookies/set?k1=" + v).SetEnableCookie(true).String()
	require.NoError(t, err)

	str, err := Get("http://localhost:8080/cookies").SetEnableCookie(true).String()
	require.NoError(t, err)

	n := strings.Index(str, v)
	require.NotEqual(t, -1, n)
}

func (h *HttplibTestSuite) TestWithBasicAuth() {
	t := h.T()
	str, err := Get("http://localhost:8080/basic-auth/user/passwd").SetBasicAuth("user", "passwd").String()
	require.NoError(t, err)
	n := strings.Index(str, "authenticated")
	require.NotEqual(t, -1, n)
}

func (h *HttplibTestSuite) TestWithUserAgent() {
	t := h.T()
	v := "beego"
	str, err := Get("http://localhost:8080/headers").SetUserAgent(v).String()
	require.NoError(t, err)
	n := strings.Index(str, v)
	require.NotEqual(t, -1, n)
}

func (h *HttplibTestSuite) TestWithSetting() {
	t := h.T()
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

	str, err := Get("http://localhost:8080/get").String()
	require.NoError(t, err)
	n := strings.Index(str, v)
	require.NotEqual(t, -1, n)
}

func (h *HttplibTestSuite) TestToJson() {
	t := h.T()
	req := Get("http://localhost:8080/ip")
	resp, err := req.Response()
	require.NoError(t, err)
	t.Log(resp)

	type IP struct {
		Origin string `json:"origin"`
	}
	var ip IP
	err = req.ToJSON(&ip)
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", ip.Origin)

	ips := strings.Split(ip.Origin, ",")
	require.NotEmpty(t, ips)
}

func (h *HttplibTestSuite) TestToFile() {
	t := h.T()
	f := "beego_testfile"
	req := Get("http://localhost:8080/ip")
	err := req.ToFile(f)
	require.NoError(t, err)
	defer os.Remove(f)

	b, err := os.ReadFile(f)
	n := bytes.Index(b, []byte("origin"))
	require.NotEqual(t, -1, n)
}

func (h *HttplibTestSuite) TestToFileDir() {
	t := h.T()
	f := "./files/beego_testfile"
	req := Get("http://localhost:8080/ip")
	err := req.ToFile(f)
	require.NoError(t, err)
	defer os.RemoveAll("./files")
	b, err := os.ReadFile(f)
	require.NoError(t, err)
	n := bytes.Index(b, []byte("origin"))
	require.NotEqual(t, -1, n)
}

func (h *HttplibTestSuite) TestHeader() {
	t := h.T()
	req := Get("http://localhost:8080/headers")
	req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36")
	_, err := req.String()
	require.NoError(t, err)
}

// TestAddFilter make sure that AddFilters only work for the specific request
func (h *HttplibTestSuite) TestAddFilter() {
	t := h.T()
	req := Get("http://beego.vip")
	req.AddFilters(func(next Filter) Filter {
		return func(ctx context.Context, req *BeegoHTTPRequest) (*http.Response, error) {
			return next(ctx, req)
		}
	})

	r := Get("http://beego.vip")
	assert.Equal(t, 1, len(req.setting.FilterChains)-len(r.setting.FilterChains))
}

func (h *HttplibTestSuite) TestFilterChainOrder() {
	t := h.T()
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

func (h *HttplibTestSuite) TestHead() {
	t := h.T()
	req := Head("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "HEAD", req.req.Method)
}

func (h *HttplibTestSuite) TestDelete() {
	t := h.T()
	req := Delete("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "DELETE", req.req.Method)
}

func (h *HttplibTestSuite) TestPost() {
	t := h.T()
	req := Post("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.req.Method)
}

func (h *HttplibTestSuite) TestPut() {
	t := h.T()
	req := Put("http://beego.vip")
	assert.NotNil(t, req)
	assert.Equal(t, "PUT", req.req.Method)
}

func TestNewBeegoRequest(t *testing.T) {
	req := NewBeegoRequest("http://beego.vip", "GET")
	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.req.Method)

	// invalid case but still go request
	req = NewBeegoRequest("httpa\ta://beego.vip", "GET")
	assert.NotNil(t, req)
}

func TestNewBeegoRequestWithCtx(t *testing.T) {
	req := NewBeegoRequestWithCtx(context.Background(), "http://beego.vip", "GET")
	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.req.Method)

	// bad url but still get request
	req = NewBeegoRequestWithCtx(context.Background(), "httpa\ta://beego.vip", "GET")
	assert.NotNil(t, req)

	// bad method but still get request
	req = NewBeegoRequestWithCtx(context.Background(), "http://beego.vip", "G\tET")
	assert.NotNil(t, req)
}

func TestBeegoHTTPRequestSetProtocolVersion(t *testing.T) {
	req := NewBeegoRequest("http://beego.vip", "GET")
	assert.Equal(t, 1, req.req.ProtoMajor)
	assert.Equal(t, 1, req.req.ProtoMinor)

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
