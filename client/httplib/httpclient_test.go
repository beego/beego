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
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("test1", "http://beego.me", WithEnableCookie(true))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, true, client.Setting.EnableCookie)
}

func TestClient_Response(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	var resp *http.Response
	err = client.Response(&resp).Get(nil, "status/203")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 203, resp.StatusCode)
}

func TestClient_StatusCode(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	var statusCode *int
	err = client.StatusCode(&statusCode).Get(nil, "status/203")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 203, *statusCode)
}

func TestClient_Headers(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	var header *http.Header
	err = client.Headers(&header).Get(nil, "bytes/123")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "123", header.Get("Content-Length"))
}

func TestClient_HeaderValue(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	var val *string
	err = client.Headers(nil).HeaderValue("Content-Length", &val).Get(nil, "bytes/123")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "123", *val)
}

func TestClient_ContentType(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	var contentType *string
	err = client.ContentType(&contentType).Get(nil, "bytes/123")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "application/octet-stream", *contentType)
}

func TestClient_ContentLength(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	var contentLength *int64
	err = client.ContentLength(&contentLength).Get(nil, "bytes/123")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(123), *contentLength)
}

type total struct {
	Slideshow slideshow `json:"slideshow" yaml:"slideshow"`
}

type slideshow struct {
	XMLName xml.Name `xml:"slideshow"`

	Title  string  `json:"title" yaml:"title" xml:"title,attr"`
	Author string  `json:"author" yaml:"author" xml:"author,attr"`
	Date   string  `json:"date" yaml:"date" xml:"date,attr"`
	Slides []slide `json:"slides" yaml:"slides" xml:"slide"`
}

type slide struct {
	XMLName xml.Name `xml:"slide"`

	Title string `json:"title" yaml:"title" xml:"title"`
}

func TestClient_Get(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	// basic type
	var s *string
	err = client.Get(&s, "/base64/SFRUUEJJTiBpcyBhd2Vzb21l")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "HTTPBIN is awesome", *s)

	var bytes *[]byte
	err = client.Get(&bytes, "/base64/SFRUUEJJTiBpcyBhd2Vzb21l")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []byte("HTTPBIN is awesome"), *bytes)

	// json
	var tp *total
	err = client.Get(&tp, "/json")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Sample Slide Show", tp.Slideshow.Title)
	assert.Equal(t, 2, len(tp.Slideshow.Slides))
	assert.Equal(t, "Overview", tp.Slideshow.Slides[1].Title)

	// xml
	var ssp *slideshow
	err = client.Get(&ssp, "/base64/PD94bWwgPz48c2xpZGVzaG93CnRpdGxlPSJTYW1wbGUgU2xpZGUgU2hvdyIKZGF0ZT0iRGF0ZSBvZiBwdWJsaWNhdGlvbiIKYXV0aG9yPSJZb3VycyBUcnVseSI+PHNsaWRlIHR5cGU9ImFsbCI+PHRpdGxlPldha2UgdXAgdG8gV29uZGVyV2lkZ2V0cyE8L3RpdGxlPjwvc2xpZGU+PHNsaWRlIHR5cGU9ImFsbCI+PHRpdGxlPk92ZXJ2aWV3PC90aXRsZT48aXRlbT5XaHkgPGVtPldvbmRlcldpZGdldHM8L2VtPiBhcmUgZ3JlYXQ8L2l0ZW0+PGl0ZW0vPjxpdGVtPldobyA8ZW0+YnV5czwvZW0+IFdvbmRlcldpZGdldHM8L2l0ZW0+PC9zbGlkZT48L3NsaWRlc2hvdz4=")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Sample Slide Show", ssp.Title)
	assert.Equal(t, 2, len(ssp.Slides))
	assert.Equal(t, "Overview", ssp.Slides[1].Title)

	// yaml
	tp = nil
	err = client.Get(&tp, "/base64/c2xpZGVzaG93OgogIGF1dGhvcjogWW91cnMgVHJ1bHkKICBkYXRlOiBkYXRlIG9mIHB1YmxpY2F0aW9uCiAgc2xpZGVzOgogIC0gdGl0bGU6IFdha2UgdXAgdG8gV29uZGVyV2lkZ2V0cyEKICAgIHR5cGU6IGFsbAogIC0gaXRlbXM6CiAgICAtIFdoeSA8ZW0+V29uZGVyV2lkZ2V0czwvZW0+IGFyZSBncmVhdAogICAgLSBXaG8gPGVtPmJ1eXM8L2VtPiBXb25kZXJXaWRnZXRzCiAgICB0aXRsZTogT3ZlcnZpZXcKICAgIHR5cGU6IGFsbAogIHRpdGxlOiBTYW1wbGUgU2xpZGUgU2hvdw==")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Sample Slide Show", tp.Slideshow.Title)
	assert.Equal(t, 2, len(tp.Slideshow.Slides))
	assert.Equal(t, "Overview", tp.Slideshow.Slides[1].Title)

}

func TestClient_Post(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	var s *string
	err = client.Get(&s, "/json")
	if err != nil {
		t.Fatal(err)
	}

	var resp *http.Response
	err = client.Response(&resp).Post(&s, "/post", *s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodPost, resp.Request.Method)
}

func TestClient_Put(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	var s *string
	err = client.Get(&s, "/json")
	if err != nil {
		t.Fatal(err)
	}

	var resp *http.Response
	err = client.Response(&resp).Put(&s, "/put", *s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodPut, resp.Request.Method)
}

func TestClient_Delete(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	var resp *http.Response
	err = client.Response(&resp).Delete(nil, "/delete")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodDelete, resp.Request.Method)
}

func TestClient_Head(t *testing.T) {
	client, err := NewClient("test", "http://beego.me")
	if err != nil {
		t.Fatal(err)
	}

	var resp *http.Response
	err = client.Response(&resp).Head(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodHead, resp.Request.Method)
}
