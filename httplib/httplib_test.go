// Beego (http://beego.me)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package httplib

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGetUrl(t *testing.T) {
	resp, err := Get("http://beego.me").Debug(true).Response()
	if err != nil {
		t.Fatal(err)
	}
	if resp.Body == nil {
		t.Fatal("body is nil")
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("data is no")
	}

	str, err := Get("http://beego.me").String()
	if err != nil {
		t.Fatal(err)
	}
	if len(str) == 0 {
		t.Fatal("has no info")
	}
}

func ExamplePost(t *testing.T) {
	b := Post("http://beego.me/").Debug(true)
	b.Param("username", "astaxie")
	b.Param("password", "hello")
	b.PostFile("uploadfile", "httplib_test.go")
	str, err := b.String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}

func TestSimpleGetString(t *testing.T) {
	fmt.Println("TestSimpleGetString==========================================")
	html, err := Get("http://httpbin.org/headers").SetAgent("beegoooooo").String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(html)
	fmt.Println("TestSimpleGetString==========================================")
}

func TestSimpleGetStringWithDefaultCookie(t *testing.T) {
	fmt.Println("TestSimpleGetStringWithDefaultCookie==========================================")
	html, err := Get("http://httpbin.org/cookies/set?k1=v1").SetEnableCookie(true).String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(html)
	html, err = Get("http://httpbin.org/cookies").SetEnableCookie(true).String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(html)
	fmt.Println("TestSimpleGetStringWithDefaultCookie==========================================")
}

func TestDefaultSetting(t *testing.T) {
	fmt.Println("TestDefaultSetting==========================================")
	var def BeegoHttpSettings
	def.EnableCookie = true
	//def.ShowDebug = true
	def.UserAgent = "UserAgent"
	//def.ConnectTimeout = 60*time.Second
	//def.ReadWriteTimeout = 60*time.Second
	def.Transport = nil //http.DefaultTransport
	SetDefaultSetting(def)

	html, err := Get("http://httpbin.org/headers").String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(html)
	html, err = Get("http://httpbin.org/headers").String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(html)
	fmt.Println("TestDefaultSetting==========================================")
}
