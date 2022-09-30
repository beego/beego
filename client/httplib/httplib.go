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

// Package httplib is used as http.Client
// Usage:
//
// import "github.com/beego/beego/v2/client/httplib"
//
//	b := httplib.Post("http://beego.vip/")
//	b.Param("username","astaxie")
//	b.Param("password","123456")
//	b.PostFile("uploadfile1", "httplib.pdf")
//	b.PostFile("uploadfile2", "httplib.txt")
//	str, err := b.String()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(str)
//
//  more docs http://beego.vip/docs/module/httplib.md
package httplib

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/beego/beego/v2/core/berror"
	"github.com/beego/beego/v2/core/logs"
)

const (
	contentTypeKey         = "Content-Type"
	ApplicationJSON        = "application/json"
	ApplicationAtomXML     = "application/atom+xml"
	ApplicationEcmascript  = "application/ecmascript"
	ApplicationJavaScript  = "application/javascript"
	ApplicationVndJSON     = "application/vnd.api+json"
	ApplicationOctetStream = "application/octet-stream"
	ApplicationOgg         = "application/ogg"
	ApplicationPdf         = "application/pdf"
	ApplicationPostscript  = "application/postscript"
	ApplicationRdfXML      = "application/rdf+xml"
	ApplicationRssXML      = "application/rss+xml"
	ApplicationXML         = "application/xml"
	ApplicationUrlencoded  = "application/x-www-form-urlencoded"
	ApplicationFontWoff    = "application/font-woff"
	ApplicationSoap        = "application/rdf+soap"
	ApplicationXYaml       = "application/x-yaml"
	ApplicationXHTMlXML    = "application/xhtml+xml"
	ApplicationXMLDtd      = "application/xml-dtd"
	ApplicationXMLXop      = "application/xml-xop"
	ApplicationZip         = "application/zip"
	ApplicationGZip        = "application/gzip"
	ApplicationGraphql     = "application/graphql"
)

// it will be the last filter and execute request.Do
var doRequestFilter = func(ctx context.Context, req *BeegoHTTPRequest) (*http.Response, error) {
	return req.doRequest(ctx)
}

var doRequestFilterWithMediaType = func(ctx context.Context, mediaType string, req *BeegoHTTPRequest) (*http.Response,
	error) {
	return req.doRequestWithMediaType(ctx, mediaType)
}

// NewBeegoRequest returns *BeegoHttpRequest with specific method
// TODO add error as return value
// I think if we don't return error
// users are hard to check whether we create Beego request successfully
func NewBeegoRequest(rawurl, method string) *BeegoHTTPRequest {
	return NewBeegoRequestWithCtx(context.Background(), rawurl, method)
}

// NewBeegoRequestWithCtx returns a new BeegoHTTPRequest given a method, URL
func NewBeegoRequestWithCtx(ctx context.Context, rawurl, method string) *BeegoHTTPRequest {
	req, err := http.NewRequestWithContext(ctx, method, rawurl, nil)
	if err != nil {
		logs.Error("%+v", berror.Wrapf(err, InvalidURLOrMethod, "invalid raw url or method: %s %s", rawurl, method))
	}

	return &BeegoHTTPRequest{
		url:     rawurl,
		req:     req,
		params:  map[string][]string{},
		files:   map[string]string{},
		setting: defaultSetting,
		resp:    &http.Response{},
	}
}

// Get returns *BeegoHttpRequest with GET method.
func Get(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "GET")
}

// Post returns *BeegoHttpRequest with POST method.
func Post(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "POST")
}

// Put returns *BeegoHttpRequest with PUT method.
func Put(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "PUT")
}

// Delete returns *BeegoHttpRequest DELETE method.
func Delete(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "DELETE")
}

// Head returns *BeegoHttpRequest with HEAD method.
func Head(url string) *BeegoHTTPRequest {
	return NewBeegoRequest(url, "HEAD")
}

// BeegoHTTPRequest provides more useful methods than http.Request for requesting a url.
type BeegoHTTPRequest struct {
	url     string
	req     *http.Request
	params  map[string][]string
	files   map[string]string
	setting BeegoHTTPSettings
	resp    *http.Response
	body    []byte
}

// GetRequest returns the request object
func (b *BeegoHTTPRequest) GetRequest() *http.Request {
	return b.req
}

// Setting changes request settings
func (b *BeegoHTTPRequest) Setting(setting BeegoHTTPSettings) *BeegoHTTPRequest {
	b.setting = setting
	return b
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (b *BeegoHTTPRequest) SetBasicAuth(username, password string) *BeegoHTTPRequest {
	b.req.SetBasicAuth(username, password)
	return b
}

// SetEnableCookie sets enable/disable cookiejar
func (b *BeegoHTTPRequest) SetEnableCookie(enable bool) *BeegoHTTPRequest {
	b.setting.EnableCookie = enable
	return b
}

// SetUserAgent sets User-Agent header field
func (b *BeegoHTTPRequest) SetUserAgent(useragent string) *BeegoHTTPRequest {
	b.setting.UserAgent = useragent
	return b
}

// Retries sets Retries times.
// default is 0 (never retry)
// -1 retry indefinitely (forever)
// Other numbers specify the exact retry amount
func (b *BeegoHTTPRequest) Retries(times int) *BeegoHTTPRequest {
	b.setting.Retries = times
	return b
}

// RetryDelay sets the time to sleep between reconnection attempts
func (b *BeegoHTTPRequest) RetryDelay(delay time.Duration) *BeegoHTTPRequest {
	b.setting.RetryDelay = delay
	return b
}

// SetTimeout sets connect time out and read-write time out for BeegoRequest.
func (b *BeegoHTTPRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *BeegoHTTPRequest {
	b.setting.ConnectTimeout = connectTimeout
	b.setting.ReadWriteTimeout = readWriteTimeout
	return b
}

// SetTLSClientConfig sets TLS connection configuration if visiting HTTPS url.
func (b *BeegoHTTPRequest) SetTLSClientConfig(config *tls.Config) *BeegoHTTPRequest {
	b.setting.TLSClientConfig = config
	return b
}

// Header adds header item string in request.
func (b *BeegoHTTPRequest) Header(key, value string) *BeegoHTTPRequest {
	b.req.Header.Set(key, value)
	return b
}

// SetHost set the request host
func (b *BeegoHTTPRequest) SetHost(host string) *BeegoHTTPRequest {
	b.req.Host = host
	return b
}

// SetProtocolVersion sets the protocol version for incoming requests.
// Client requests always use HTTP/1.1
func (b *BeegoHTTPRequest) SetProtocolVersion(vers string) *BeegoHTTPRequest {
	if vers == "" {
		vers = "HTTP/1.1"
	}

	major, minor, ok := http.ParseHTTPVersion(vers)
	if ok {
		b.req.Proto = vers
		b.req.ProtoMajor = major
		b.req.ProtoMinor = minor
		return b
	}
	logs.Error("%+v", berror.Errorf(InvalidUrlProtocolVersion, "invalid protocol: %s", vers))
	return b
}

// SetCookie adds a cookie to the request.
func (b *BeegoHTTPRequest) SetCookie(cookie *http.Cookie) *BeegoHTTPRequest {
	b.req.Header.Add("Cookie", cookie.String())
	return b
}

// SetTransport sets the transport field
func (b *BeegoHTTPRequest) SetTransport(transport http.RoundTripper) *BeegoHTTPRequest {
	b.setting.Transport = transport
	return b
}

// SetProxy sets the HTTP proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (b *BeegoHTTPRequest) SetProxy(proxy func(*http.Request) (*url.URL, error)) *BeegoHTTPRequest {
	b.setting.Proxy = proxy
	return b
}

// SetCheckRedirect specifies the policy for handling redirects.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func (b *BeegoHTTPRequest) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *BeegoHTTPRequest {
	b.setting.CheckRedirect = redirect
	return b
}

// SetFilters will use the filter as the invocation filters
func (b *BeegoHTTPRequest) SetFilters(fcs ...FilterChain) *BeegoHTTPRequest {
	b.setting.FilterChains = fcs
	return b
}

// AddFilters adds filter
func (b *BeegoHTTPRequest) AddFilters(fcs ...FilterChain) *BeegoHTTPRequest {
	b.setting.FilterChains = append(b.setting.FilterChains, fcs...)
	return b
}

// SetEscapeHTML is used to set the flag whether escape HTML special characters during processing
func (b *BeegoHTTPRequest) SetEscapeHTML(isEscape bool) *BeegoHTTPRequest {
	b.setting.EscapeHTML = isEscape
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *BeegoHTTPRequest) Param(key, value string) *BeegoHTTPRequest {
	if param, ok := b.params[key]; ok {
		b.params[key] = append(param, value)
	} else {
		b.params[key] = []string{value}
	}
	return b
}

// PostFile adds a post file to the request
func (b *BeegoHTTPRequest) PostFile(formname, filename string) *BeegoHTTPRequest {
	b.files[formname] = filename
	return b
}

// Body adds request raw body.
// Supports string and []byte.
// TODO return error if data is invalid
func (b *BeegoHTTPRequest) Body(data interface{}) *BeegoHTTPRequest {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bf), nil
		}
		b.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bf), nil
		}
		b.req.ContentLength = int64(len(t))
	default:
		logs.Error("%+v", berror.Errorf(UnsupportedBodyType, "unsupported body data type: %s", t))
	}
	return b
}

// XMLBody adds the request raw body encoded in XML.
func (b *BeegoHTTPRequest) XMLBody(obj interface{}) (*BeegoHTTPRequest, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := xml.Marshal(obj)
		if err != nil {
			return b, berror.Wrap(err, InvalidXMLBody, "obj could not be converted to XML data")
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(byts)), nil
		}
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set(contentTypeKey, "application/xml")
	}
	return b, nil
}

// YAMLBody adds the request raw body encoded in YAML.
func (b *BeegoHTTPRequest) YAMLBody(obj interface{}) (*BeegoHTTPRequest, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := yaml.Marshal(obj)
		if err != nil {
			return b, berror.Wrap(err, InvalidYAMLBody, "obj could not be converted to YAML data")
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set(contentTypeKey, "application/x+yaml")
	}
	return b, nil
}

// JSONBody adds the request raw body encoded in JSON.
func (b *BeegoHTTPRequest) JSONBody(obj interface{}) (*BeegoHTTPRequest, error) {
	if b.req.Body == nil && obj != nil {
		byts, err := b.JSONMarshal(obj)
		if err != nil {
			return b, berror.Wrap(err, InvalidJSONBody, "obj could not be converted to JSON body")
		}
		b.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		b.req.ContentLength = int64(len(byts))
		b.req.Header.Set(contentTypeKey, "application/json")
	}
	return b, nil
}

func (b *BeegoHTTPRequest) JSONMarshal(obj interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(b.setting.EscapeHTML)
	err := jsonEncoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (b *BeegoHTTPRequest) buildURL(paramBody string) {
	// build GET url with query string
	if b.req.Method == "GET" && len(paramBody) > 0 {
		if strings.Contains(b.url, "?") {
			b.url += "&" + paramBody
		} else {
			b.url = b.url + "?" + paramBody
		}
		return
	}

	// build POST/PUT/PATCH url and body
	if (b.req.Method == "POST" || b.req.Method == "PUT" || b.req.Method == "PATCH" || b.req.Method == "DELETE") && b.req.Body == nil {
		// with files
		if len(b.files) > 0 {
			b.handleFiles()
			return
		}

		// with params
		if len(paramBody) > 0 {
			b.Header(contentTypeKey, "application/x-www-form-urlencoded")
			b.Body(paramBody)
		}
	}
}

func (b *BeegoHTTPRequest) handleFiles() {
	pr, pw := io.Pipe()
	bodyWriter := multipart.NewWriter(pw)
	go func() {
		for formname, filename := range b.files {
			b.handleFileToBody(bodyWriter, formname, filename)
		}
		for k, v := range b.params {
			for _, vv := range v {
				_ = bodyWriter.WriteField(k, vv)
			}
		}
		_ = bodyWriter.Close()
		_ = pw.Close()
	}()
	b.Header(contentTypeKey, bodyWriter.FormDataContentType())
	b.req.Body = ioutil.NopCloser(pr)
	b.Header("Transfer-Encoding", "chunked")
}

func (*BeegoHTTPRequest) handleFileToBody(bodyWriter *multipart.Writer, formname string, filename string) {
	fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
	const errFmt = "Httplib: %+v"
	if err != nil {
		logs.Error(errFmt, berror.Wrapf(err, CreateFormFileFailed,
			"could not create form file, formname: %s, filename: %s", formname, filename))
	}
	fh, err := os.Open(filename)
	if err != nil {
		logs.Error(errFmt, berror.Wrapf(err, ReadFileFailed, "could not open this file %s", filename))
	}
	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		logs.Error(errFmt, berror.Wrapf(err, CopyFileFailed, "could not copy this file %s", filename))
	}
	err = fh.Close()
	if err != nil {
		logs.Error(errFmt, berror.Wrapf(err, CloseFileFailed, "could not close this file %s", filename))
	}
}

func (b *BeegoHTTPRequest) getResponse() (*http.Response, error) {
	if b.resp.StatusCode != 0 {
		return b.resp, nil
	}
	resp, err := b.DoRequest()
	if err != nil {
		return nil, err
	}
	b.resp = resp
	return resp, nil
}

func (b *BeegoHTTPRequest) getResponseWithMediaType(mediaType string) (*http.Response, error) {
	if b.resp.StatusCode != 0 {
		return b.resp, nil
	}
	resp, err := b.DoRequestWithMediaType(mediaType)
	if err != nil {
		return nil, err
	}
	b.resp = resp
	return resp, nil
}

// DoRequest executes client.Do
func (b *BeegoHTTPRequest) DoRequest() (resp *http.Response, err error) {
	root := doRequestFilter
	if len(b.setting.FilterChains) > 0 {
		for i := len(b.setting.FilterChains) - 1; i >= 0; i-- {
			root = b.setting.FilterChains[i](root)
		}
	}
	return root(b.req.Context(), b)
}

// Deprecated: please use NewBeegoRequestWithContext
func (b *BeegoHTTPRequest) DoRequestWithCtx(ctx context.Context) (resp *http.Response, err error) {
	root := doRequestFilter
	if len(b.setting.FilterChains) > 0 {
		for i := len(b.setting.FilterChains) - 1; i >= 0; i-- {
			root = b.setting.FilterChains[i](root)
		}
	}
	return root(ctx, b)
}

// DoRequest executes client.DoWithMediaType
func (b *BeegoHTTPRequest) DoRequestWithMediaType(mediaType string) (resp *http.Response, err error) {
	root := doRequestFilterWithMediaType
	if len(b.setting.FilterChains) > 0 {
		for i := len(b.setting.FilterChains) - 1; i >= 0; i-- {
			root = b.setting.FilterWithMediaChains[i](root)
		}
	}
	return root(b.req.Context(), mediaType, b)
}

func (b *BeegoHTTPRequest) doRequest(_ context.Context) (*http.Response, error) {
	paramBody := b.buildParamBody()

	b.buildURL(paramBody)
	urlParsed, err := url.Parse(b.url)
	if err != nil {
		return nil, berror.Wrapf(err, InvalidUrl, "parse url failed, the url is %s", b.url)
	}

	b.req.URL = urlParsed

	trans := b.buildTrans()

	jar := b.buildCookieJar()

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if b.setting.UserAgent != "" && b.req.Header.Get("User-Agent") == "" {
		b.req.Header.Set("User-Agent", b.setting.UserAgent)
	}

	if b.setting.CheckRedirect != nil {
		client.CheckRedirect = b.setting.CheckRedirect
	}

	return b.sendRequest(client)
}

func (b *BeegoHTTPRequest) doRequestWithMediaType(_ context.Context, mediaType string) (*http.Response, error) {
	paramBody := b.buildParamBody()

	b.buildURL(paramBody)
	urlParsed, err := url.Parse(b.url)
	if err != nil {
		return nil, berror.Wrapf(err, InvalidUrl, "parse url failed, the url is %s", b.url)
	}

	b.req.URL = urlParsed

	trans := b.buildTrans()

	jar := b.buildCookieJar()

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if b.setting.UserAgent != "" && b.req.Header.Get("User-Agent") == "" {
		b.req.Header.Set("User-Agent", b.setting.UserAgent)
	}

	if b.setting.CheckRedirect != nil {
		client.CheckRedirect = b.setting.CheckRedirect
	}

	return b.sendRequestWithMediaType(client, mediaType)
}

func (b *BeegoHTTPRequest) sendRequest(client *http.Client) (resp *http.Response, err error) {
	// retries default value is 0, it will run once.
	// retries equal to -1, it will run forever until success
	// retries is setted, it will retries fixed times.
	// Sleeps for a 400ms between calls to reduce spam
	for i := 0; b.setting.Retries == -1 || i <= b.setting.Retries; i++ {
		resp, err = client.Do(b.req)
		if err == nil {
			return
		}
		time.Sleep(b.setting.RetryDelay)
	}
	return nil, berror.Wrap(err, SendRequestFailed, "sending request fail")
}

func (b *BeegoHTTPRequest) sendRequestWithMediaType(client *http.Client, mediaType string) (resp *http.Response, err error) {
	// retries default value is 0, it will run once.
	// retries equal to -1, it will run forever until success
	// retries is setted, it will retries fixed times.
	// Sleeps for a 400ms between calls to reduce spam
	for i := 0; b.setting.Retries == -1 || i <= b.setting.Retries; i++ {
		resp, err = client.DoWithMediaType(b.req, mediaType)
		if err == nil {
			return
		}
		time.Sleep(b.setting.RetryDelay)
	}
	return nil, berror.Wrap(err, SendRequestFailed, "sending request fail")
}

func (b *BeegoHTTPRequest) buildCookieJar() http.CookieJar {
	var jar http.CookieJar
	if b.setting.EnableCookie {
		if defaultCookieJar == nil {
			createDefaultCookie()
		}
		jar = defaultCookieJar
	}
	return jar
}

func (b *BeegoHTTPRequest) buildTrans() http.RoundTripper {
	trans := b.setting.Transport

	if trans == nil {
		// create default transport
		trans = &http.Transport{
			TLSClientConfig:     b.setting.TLSClientConfig,
			Proxy:               b.setting.Proxy,
			DialContext:         TimeoutDialerCtx(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout),
			MaxIdleConnsPerHost: 100,
		}
	} else if t, ok := trans.(*http.Transport); ok {
		// if b.transport is *http.Transport then set the settings.
		if t.TLSClientConfig == nil {
			t.TLSClientConfig = b.setting.TLSClientConfig
		}
		if t.Proxy == nil {
			t.Proxy = b.setting.Proxy
		}
		if t.DialContext == nil {
			t.DialContext = TimeoutDialerCtx(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout)
		}
	}
	return trans
}

func (b *BeegoHTTPRequest) buildParamBody() string {
	var paramBody string
	if len(b.params) > 0 {
		var buf bytes.Buffer
		for k, v := range b.params {
			for _, vv := range v {
				buf.WriteString(url.QueryEscape(k))
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(vv))
				buf.WriteByte('&')
			}
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}
	return paramBody
}

// String returns the body string in response.
// Calls Response inner.
func (b *BeegoHTTPRequest) String() (string, error) {
	data, err := b.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// Calls Response inner.
func (b *BeegoHTTPRequest) Bytes() ([]byte, error) {
	if b.body != nil {
		return b.body, nil
	}
	resp, err := b.getResponse()
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	if b.setting.Gzip && resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, berror.Wrap(err, ReadGzipBodyFailed, "building gzip reader failed")
		}
		b.body, err = ioutil.ReadAll(reader)
		return b.body, berror.Wrap(err, ReadGzipBodyFailed, "reading gzip data failed")
	}
	b.body, err = ioutil.ReadAll(resp.Body)
	return b.body, err
}

// ToFile saves the body data in response to one file.
// Calls Response inner.
func (b *BeegoHTTPRequest) ToFile(filename string) error {
	resp, err := b.getResponse()
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	err = pathExistAndMkdir(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// Check if the file directory exists. If it doesn't then it's created
func pathExistAndMkdir(filename string) (err error) {
	filename = path.Dir(filename)
	_, err = os.Stat(filename)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(filename, os.ModePerm)
		if err == nil {
			return nil
		}
	}
	return berror.Wrapf(err, CreateFileIfNotExistFailed, "try to create(if not exist) failed: %s", filename)
}

// ToJSON returns the map that marshals from the body bytes as json in response.
// Calls Response inner.
func (b *BeegoHTTPRequest) ToJSON(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return berror.Wrap(json.Unmarshal(data, v),
		UnmarshalJSONResponseToObjectFailed, "unmarshal json body to object failed.")
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// Calls Response inner.
func (b *BeegoHTTPRequest) ToXML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return berror.Wrap(xml.Unmarshal(data, v),
		UnmarshalXMLResponseToObjectFailed, "unmarshal xml body to object failed.")
}

// ToYAML returns the map that marshals from the body bytes as yaml in response .
// Calls Response inner.
func (b *BeegoHTTPRequest) ToYAML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return berror.Wrap(yaml.Unmarshal(data, v),
		UnmarshalYAMLResponseToObjectFailed, "unmarshal yaml body to object failed.")
}

// ToValue attempts to resolve the response body to value using an existing method.
// Calls Response inner.
// If response header contain Content-Type, func will call ToJSON\ToXML\ToYAML.
// Else it will try to parse body as json\yaml\xml, If all attempts fail, an error will be returned
func (b *BeegoHTTPRequest) ToValue(value interface{}) error {
	if value == nil {
		return nil
	}

	contentType := strings.Split(b.resp.Header.Get(contentTypeKey), ";")[0]
	// try to parse it as content type
	switch contentType {
	case "application/json":
		return b.ToJSON(value)
	case "text/xml", "application/xml":
		return b.ToXML(value)
	case "text/yaml", "application/x-yaml", "application/x+yaml":
		return b.ToYAML(value)
	}

	// try to parse it anyway
	if err := b.ToJSON(value); err == nil {
		return nil
	}
	if err := b.ToYAML(value); err == nil {
		return nil
	}
	if err := b.ToXML(value); err == nil {
		return nil
	}

	return berror.Error(UnmarshalResponseToObjectFailed, "unmarshal body to object failed.")
}

// Response executes request client gets response manually.
func (b *BeegoHTTPRequest) Response() (*http.Response, error) {
	return b.getResponse()
}

// ResponseWithMediaType Response executes request client gets response with media type manually.
func (b *BeegoHTTPRequest) ResponseWithMediaType(mediaType string) (*http.Response, error) {
	return b.getResponseWithMediaType(mediaType)
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
// Deprecated
// we will move this at the end of 2021
// please use TimeoutDialerCtx
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		return TimeoutDialerCtx(cTimeout, rwTimeout)(context.Background(), netw, addr)
	}
}

func TimeoutDialerCtx(cTimeout time.Duration,
	rwTimeout time.Duration) func(ctx context.Context, net, addr string) (c net.Conn, err error) {
	return func(ctx context.Context, netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, err
	}
}
