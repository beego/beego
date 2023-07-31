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
	"encoding/json"
	"encoding/xml"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("test1", "http://beego.vip", WithEnableCookie(true))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, true, client.Setting.EnableCookie)
}

type slideShowResponse struct {
	Resp       *http.Response      `json:"resp,omitempty"`
	bytes      []byte              `json:"bytes,omitempty"`
	StatusCode int                 `json:"status_code,omitempty"`
	Body       io.ReadCloser       `json:"body,omitempty"`
	Header     map[string][]string `json:"header,omitempty"`

	Slideshow slideshow `json:"slideshow,omitempty" yaml:"slideshow" xml:"slideshow"`
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
	//XMLName xml.Name `xml:"slideshow"`

	Title  string  `json:"title" yaml:"title" xml:"title,attr"`
	Author string  `json:"author" yaml:"author" xml:"author,attr"`
	Date   string  `json:"date" yaml:"date" xml:"date,attr"`
	Slides []slide `json:"slides" yaml:"slides" xml:"slide"`
}

type slide struct {
	XMLName xml.Name `xml:"slide"`

	Title string `json:"title" yaml:"title" xml:"title"`
}

type ClientTestSuite struct {
	suite.Suite
	l net.Listener
}

func (c *ClientTestSuite) SetupSuite() {
	listener, err := net.Listen("tcp", ":8080")
	require.NoError(c.T(), err)
	c.l = listener

	handler := http.NewServeMux()
	handler.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		data, _ := json.Marshal(slideshow{})
		_, _ = writer.Write(data)
	})

	ssr := slideShowResponse{
		Slideshow: slideshow{
			Title: "Sample Slide Show",
			Slides: []slide{
				{
					Title: "Content",
				},
				{
					Title: "Overview",
				},
			},
		},
	}

	handler.HandleFunc("/req2resp", func(writer http.ResponseWriter, request *http.Request) {
		data, _ := io.ReadAll(request.Body)
		_, _ = writer.Write(data)
	})

	handler.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		data, _ := json.Marshal(ssr)
		_, _ = writer.Write(data)
	})

	handler.HandleFunc("/get/xml", func(writer http.ResponseWriter, request *http.Request) {
		data, err := xml.Marshal(ssr.Slideshow)
		require.NoError(c.T(), err)
		_, _ = writer.Write(data)
	})

	handler.HandleFunc("/get/yaml", func(writer http.ResponseWriter, request *http.Request) {
		data, _ := yaml.Marshal(ssr)
		_, _ = writer.Write(data)
	})

	go func() {
		_ = http.Serve(listener, handler)
	}()
}

func (c *ClientTestSuite) TearDownSuite() {
	_ = c.l.Close()
}

func TestClient(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}

func (c *ClientTestSuite) TestClientHandleCarrier() {
	t := c.T()
	v := "beego"
	client, err := NewClient("test", "http://localhost:8080/",
		WithUserAgent(v))
	require.NoError(t, err)
	resp := &slideShowResponse{}
	err = client.Get(resp, "/json")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.NotNil(t, resp.Resp)
	assert.NotNil(t, resp.Body)
	assert.Equal(t, "48", resp.Header["Content-Length"][0])
	assert.Equal(t, 200, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 48, len(b))
	assert.Equal(t, resp.String(), string(b))
}

func (c *ClientTestSuite) TestClientGet() {
	t := c.T()
	client, err := NewClient("test", "http://localhost:8080/")
	if err != nil {
		t.Fatal(err)
	}

	// json
	var s slideShowResponse
	err = client.Get(&s, "/get")
	require.NoError(t, err)
	assert.Equal(t, "Sample Slide Show", s.Slideshow.Title)
	assert.Equal(t, 2, len(s.Slideshow.Slides))
	assert.Equal(t, "Overview", s.Slideshow.Slides[1].Title)

	// xml
	var ss slideshow
	err = client.Get(&ss, "/get/xml")
	require.NoError(t, err)
	assert.Equal(t, "Sample Slide Show", ss.Title)
	assert.Equal(t, 2, len(ss.Slides))
	assert.Equal(t, "Overview", ss.Slides[1].Title)

	// yaml
	s = slideShowResponse{}
	err = client.Get(&s, "/get/yaml")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Sample Slide Show", s.Slideshow.Title)
	assert.Equal(t, 2, len(s.Slideshow.Slides))
	assert.Equal(t, "Overview", s.Slideshow.Slides[1].Title)
}

func (c *ClientTestSuite) TestClientPost() {
	t := c.T()
	client, err := NewClient("test", "http://localhost:8080")
	require.NoError(t, err)

	input := slideShowResponse{
		Slideshow: slideshow{
			Title: "Sample Slide Show",
			Slides: []slide{
				{
					Title: "Content",
				},
				{
					Title: "Overview",
				},
			},
		},
	}

	jsonStr, err := json.Marshal(input)
	require.NoError(t, err)
	resp := slideShowResponse{}
	err = client.Post(&resp, "/req2resp", jsonStr)
	require.NoError(t, err)
	assert.Equal(t, input.Slideshow, resp.Slideshow)
	assert.Equal(t, http.MethodPost, resp.Resp.Request.Method)
}

func (c *ClientTestSuite) TestClientPut() {
	t := c.T()
	client, err := NewClient("test", "http://localhost:8080")
	require.NoError(t, err)

	input := slideShowResponse{
		Slideshow: slideshow{
			Title: "Sample Slide Show",
			Slides: []slide{
				{
					Title: "Content",
				},
				{
					Title: "Overview",
				},
			},
		},
	}

	jsonStr, err := json.Marshal(input)
	require.NoError(t, err)
	resp := slideShowResponse{}
	err = client.Put(&resp, "/req2resp", jsonStr)
	require.NoError(t, err)
	assert.Equal(t, input.Slideshow, resp.Slideshow)
	assert.Equal(t, http.MethodPut, resp.Resp.Request.Method)
}

func (c *ClientTestSuite) TestClientDelete() {
	t := c.T()
	client, err := NewClient("test", "http://localhost:8080")
	require.NoError(t, err)

	resp := &slideShowResponse{}
	err = client.Delete(resp, "/req2resp")
	require.NoError(t, err)
	defer resp.Resp.Body.Close()

	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodDelete, resp.Resp.Request.Method)
}

func (c *ClientTestSuite) TestClientHead() {
	t := c.T()
	client, err := NewClient("test", "http://localhost:8080")
	require.NoError(t, err)
	resp := &slideShowResponse{}
	err = client.Head(resp, "/req2resp")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Resp.Body.Close()
	assert.NotNil(t, resp)
	assert.Equal(t, http.MethodHead, resp.Resp.Request.Method)
}
