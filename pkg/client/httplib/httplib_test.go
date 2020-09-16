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
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
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

func TestPostFile(t *testing.T) {
	rtn := struct {
		Files map[string]string `json:"files"`
	}{}
	v1 := "smallfish"
	v2 := "smallfish"
	_ = ioutil.WriteFile("./test1.info", []byte(v1), 0600)
	_ = ioutil.WriteFile("./test2.info", []byte(v2), 0600)
	defer func() {
		_ = os.RemoveAll("./test1.info")
		_ = os.RemoveAll("./test2.info")
	}()
	req := Post("http://httpbin.org/post")
	req.Debug(true)
	req.SetFileChunkSize(4)
	req.PostFile("uploadfile", "./test1.info")
	req.PostFile("file2", "./test2.info")

	err := req.ToJSON(&rtn)
	if err != nil {
		t.Fatal(err)
	}
	if rtn.Files["uploadfile"] == v1 && rtn.Files["file2"] == v2 {
		return
	}
	t.Fatal(rtn)
}

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
	err := req.ToFile(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f)
	b, err := ioutil.ReadFile(f)
	if n := strings.Index(string(b), "origin"); n == -1 {
		t.Fatal(err)
	}
}

func TestToFileDir(t *testing.T) {
	f := "./files/beego_testfile"
	req := Get("http://httpbin.org/ip")
	listener := &ProgressListener{}
	err := req.ToFile(f, listener)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("./files")
	b, err := ioutil.ReadFile(f)
	if n := strings.Index(string(b), "origin"); n == -1 {
		t.Fatal(err)
	}
	if listener.Current == listener.Total || listener.Current == 0 {
		t.Fatal(nil)
	}
	fmt.Println(listener)
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
