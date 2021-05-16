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

type slideSshowResponse struct {
	Resp  *http.Response
	bytes []byte

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

func (s *slideSshowResponse) SetHttpResponse(resp *http.Response) {
	s.Resp = resp
}

func (s *slideSshowResponse) SetBytes(bytes []byte) {
	s.bytes = bytes
}

func (s *slideSshowResponse) Bytes() []byte {
	return s.bytes
}

func (s *slideSshowResponse) String() string {
	return string(s.bytes)
}

func TestClient_Get(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	// json
	var s *slideSshowResponse
	err = client.Get(&s, "/json")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Sample Slide Show", s.Slideshow.Title)
	assert.Equal(t, 2, len(s.Slideshow.Slides))
	assert.Equal(t, "Overview", s.Slideshow.Slides[1].Title)

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
	s = nil
	err = client.Get(&s, "/base64/c2xpZGVzaG93OgogIGF1dGhvcjogWW91cnMgVHJ1bHkKICBkYXRlOiBkYXRlIG9mIHB1YmxpY2F0aW9uCiAgc2xpZGVzOgogIC0gdGl0bGU6IFdha2UgdXAgdG8gV29uZGVyV2lkZ2V0cyEKICAgIHR5cGU6IGFsbAogIC0gaXRlbXM6CiAgICAtIFdoeSA8ZW0+V29uZGVyV2lkZ2V0czwvZW0+IGFyZSBncmVhdAogICAgLSBXaG8gPGVtPmJ1eXM8L2VtPiBXb25kZXJXaWRnZXRzCiAgICB0aXRsZTogT3ZlcnZpZXcKICAgIHR5cGU6IGFsbAogIHRpdGxlOiBTYW1wbGUgU2xpZGUgU2hvdw==")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Sample Slide Show", s.Slideshow.Title)
	assert.Equal(t, 2, len(s.Slideshow.Slides))
	assert.Equal(t, "Overview", s.Slideshow.Slides[1].Title)

}

func TestClient_Post(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	var resp = &slideSshowResponse{}
	err = client.Get(resp, "/json")
	if err != nil {
		t.Fatal(err)
	}

	jsonStr := resp.String()
	err = client.Post(resp, "/post", jsonStr)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodPost, resp.Resp.Request.Method)
}

func TestClient_Put(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	var resp = &slideSshowResponse{}
	err = client.Get(resp, "/json")
	if err != nil {
		t.Fatal(err)
	}

	jsonStr := resp.String()
	err = client.Put(resp, "/put", jsonStr)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodPut, resp.Resp.Request.Method)
}

func TestClient_Delete(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	var resp = &slideSshowResponse{}
	err = client.Delete(resp, "/delete")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodDelete, resp.Resp.Request.Method)
}

func TestClient_Head(t *testing.T) {
	client, err := NewClient("test", "http://beego.me")
	if err != nil {
		t.Fatal(err)
	}

	var resp = &slideSshowResponse{}
	err = client.Head(resp, "")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodHead, resp.Resp.Request.Method)
}
