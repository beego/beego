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
	"io"
	"io/ioutil"
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

type slideShowResponse struct {
	Resp       *http.Response
	bytes      []byte
	StatusCode int
	Body       io.ReadCloser
	Header     map[string][]string

	Slideshow slideshow `json:"slideshow" yaml:"slideshow"`
}

func (r *slideShowResponse) SetHTTPResponse(resp *http.Response) {
	r.Resp = resp
}

func (r *slideShowResponse) SetBytes(bytes []byte) {
	r.bytes = bytes
}

func (r *slideShowResponse) SetReader(reader io.ReadCloser) {
	r.Body = reader
}

func (r *slideShowResponse) SetStatusCode(status int) {
	r.StatusCode = status
}

func (r *slideShowResponse) SetHeader(header map[string][]string) {
	r.Header = header
}

func (r *slideShowResponse) String() string {
	return string(r.bytes)
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

func TestClientHandleCarrier(t *testing.T) {
	v := "beego"
	client, err := NewClient("test", "http://httpbin.org/",
		WithUserAgent(v))
	if err != nil {
		t.Fatal(err)
	}

	s := &slideShowResponse{}
	err = client.Get(s, "/json")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Body.Close()

	assert.NotNil(t, s.Resp)
	assert.NotNil(t, s.Body)
	assert.Equal(t, "429", s.Header["Content-Length"][0])
	assert.Equal(t, 200, s.StatusCode)

	b, err := ioutil.ReadAll(s.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 429, len(b))
	assert.Equal(t, s.String(), string(b))
}

func TestClientGet(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org/")
	if err != nil {
		t.Fatal(err)
	}

	// json
	var s *slideShowResponse
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

func TestClientPost(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	resp := &slideShowResponse{}
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

func TestClientPut(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	resp := &slideShowResponse{}
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

func TestClientDelete(t *testing.T) {
	client, err := NewClient("test", "http://httpbin.org")
	if err != nil {
		t.Fatal(err)
	}

	resp := &slideShowResponse{}
	err = client.Delete(resp, "/delete")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Resp.Body.Close()
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodDelete, resp.Resp.Request.Method)
}

func TestClientHead(t *testing.T) {
	client, err := NewClient("test", "http://beego.me")
	if err != nil {
		t.Fatal(err)
	}

	resp := &slideShowResponse{}
	err = client.Head(resp, "")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Resp.Body.Close()
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodHead, resp.Resp.Request.Method)
}
