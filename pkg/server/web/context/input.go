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

package context

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego/pkg/infrastructure/session"
)

// Regexes for checking the accept headers
// TODO make sure these are correct
var (
	acceptsHTMLRegex = regexp.MustCompile(`(text/html|application/xhtml\+xml)(?:,|$)`)
	acceptsXMLRegex  = regexp.MustCompile(`(application/xml|text/xml)(?:,|$)`)
	acceptsJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
	acceptsYAMLRegex = regexp.MustCompile(`(application/x-yaml)(?:,|$)`)
	maxParam         = 50
)

// BeegoInput operates the http request header, data, cookie and body.
// Contains router params and current session.
type BeegoInput struct {
	Context       *Context
	CruSession    session.Store
	pnames        []string
	pvalues       []string
	data          map[interface{}]interface{} // store some values in this context when calling context in filter or controller.
	dataLock      sync.RWMutex
	RequestBody   []byte
	RunMethod     string
	RunController reflect.Type
}

// NewInput returns the BeegoInput generated by context.
func NewInput() *BeegoInput {
	return &BeegoInput{
		pnames:  make([]string, 0, maxParam),
		pvalues: make([]string, 0, maxParam),
		data:    make(map[interface{}]interface{}),
	}
}

// Reset initializes the BeegoInput
func (input *BeegoInput) Reset(ctx *Context) {
	input.Context = ctx
	input.CruSession = nil
	input.pnames = input.pnames[:0]
	input.pvalues = input.pvalues[:0]
	input.dataLock.Lock()
	input.data = nil
	input.dataLock.Unlock()
	input.RequestBody = []byte{}
}

// Protocol returns the request protocol name, such as HTTP/1.1 .
func (input *BeegoInput) Protocol() string {
	return input.Context.Request.Proto
}

// URI returns the full request url with query, string and fragment.
func (input *BeegoInput) URI() string {
	return input.Context.Request.RequestURI
}

// URL returns the request url path (without query, string and fragment).
func (input *BeegoInput) URL() string {
	return input.Context.Request.URL.Path
}

// Site returns the base site url as scheme://domain type.
func (input *BeegoInput) Site() string {
	return input.Scheme() + "://" + input.Domain()
}

// Scheme returns the request scheme as "http" or "https".
func (input *BeegoInput) Scheme() string {
	if scheme := input.Header("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if input.Context.Request.URL.Scheme != "" {
		return input.Context.Request.URL.Scheme
	}
	if input.Context.Request.TLS == nil {
		return "http"
	}
	return "https"
}

// Domain returns the host name (alias of host method)
func (input *BeegoInput) Domain() string {
	return input.Host()
}

// Host returns the host name.
// If no host info in request, return localhost.
func (input *BeegoInput) Host() string {
	if input.Context.Request.Host != "" {
		if hostPart, _, err := net.SplitHostPort(input.Context.Request.Host); err == nil {
			return hostPart
		}
		return input.Context.Request.Host
	}
	return "localhost"
}

// Method returns http request method.
func (input *BeegoInput) Method() string {
	return input.Context.Request.Method
}

// Is returns the boolean value of this request is on given method, such as Is("POST").
func (input *BeegoInput) Is(method string) bool {
	return input.Method() == method
}

// IsGet Is this a GET method request?
func (input *BeegoInput) IsGet() bool {
	return input.Is("GET")
}

// IsPost Is this a POST method request?
func (input *BeegoInput) IsPost() bool {
	return input.Is("POST")
}

// IsHead Is this a Head method request?
func (input *BeegoInput) IsHead() bool {
	return input.Is("HEAD")
}

// IsOptions Is this a OPTIONS method request?
func (input *BeegoInput) IsOptions() bool {
	return input.Is("OPTIONS")
}

// IsPut Is this a PUT method request?
func (input *BeegoInput) IsPut() bool {
	return input.Is("PUT")
}

// IsDelete Is this a DELETE method request?
func (input *BeegoInput) IsDelete() bool {
	return input.Is("DELETE")
}

// IsPatch Is this a PATCH method request?
func (input *BeegoInput) IsPatch() bool {
	return input.Is("PATCH")
}

// IsAjax returns boolean of is this request generated by ajax.
func (input *BeegoInput) IsAjax() bool {
	return input.Header("X-Requested-With") == "XMLHttpRequest"
}

// IsSecure returns boolean of this request is in https.
func (input *BeegoInput) IsSecure() bool {
	return input.Scheme() == "https"
}

// IsWebsocket returns boolean of this request is in webSocket.
func (input *BeegoInput) IsWebsocket() bool {
	return input.Header("Upgrade") == "websocket"
}

// IsUpload returns boolean of whether file uploads in this request or not..
func (input *BeegoInput) IsUpload() bool {
	return strings.Contains(input.Header("Content-Type"), "multipart/form-data")
}

// AcceptsHTML Checks if request accepts html response
func (input *BeegoInput) AcceptsHTML() bool {
	return acceptsHTMLRegex.MatchString(input.Header("Accept"))
}

// AcceptsXML Checks if request accepts xml response
func (input *BeegoInput) AcceptsXML() bool {
	return acceptsXMLRegex.MatchString(input.Header("Accept"))
}

// AcceptsJSON Checks if request accepts json response
func (input *BeegoInput) AcceptsJSON() bool {
	return acceptsJSONRegex.MatchString(input.Header("Accept"))
}

// AcceptsYAML Checks if request accepts json response
func (input *BeegoInput) AcceptsYAML() bool {
	return acceptsYAMLRegex.MatchString(input.Header("Accept"))
}

// IP returns request client ip.
// if in proxy, return first proxy id.
// if error, return RemoteAddr.
func (input *BeegoInput) IP() string {
	ips := input.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip, _, err := net.SplitHostPort(ips[0])
		if err != nil {
			rip = ips[0]
		}
		return rip
	}
	if ip, _, err := net.SplitHostPort(input.Context.Request.RemoteAddr); err == nil {
		return ip
	}
	return input.Context.Request.RemoteAddr
}

// Proxy returns proxy client ips slice.
func (input *BeegoInput) Proxy() []string {
	if ips := input.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

// Referer returns http referer header.
func (input *BeegoInput) Referer() string {
	return input.Header("Referer")
}

// Refer returns http referer header.
func (input *BeegoInput) Refer() string {
	return input.Referer()
}

// SubDomains returns sub domain string.
// if aa.bb.domain.com, returns aa.bb
func (input *BeegoInput) SubDomains() string {
	parts := strings.Split(input.Host(), ".")
	if len(parts) >= 3 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

// Port returns request client port.
// when error or empty, return 80.
func (input *BeegoInput) Port() int {
	if _, portPart, err := net.SplitHostPort(input.Context.Request.Host); err == nil {
		port, _ := strconv.Atoi(portPart)
		return port
	}
	return 80
}

// UserAgent returns request client user agent string.
func (input *BeegoInput) UserAgent() string {
	return input.Header("User-Agent")
}

// ParamsLen return the length of the params
func (input *BeegoInput) ParamsLen() int {
	return len(input.pnames)
}

// Param returns router param by a given key.
func (input *BeegoInput) Param(key string) string {
	for i, v := range input.pnames {
		if v == key && i <= len(input.pvalues) {
			// we cannot use url.PathEscape(input.pvalues[i])
			// for example, if the value is /a/b
			// after url.PathEscape(input.pvalues[i]), the value is %2Fa%2Fb
			// However, the value is used in ControllerRegister.ServeHTTP
			// and split by "/", so function crash...
			return input.pvalues[i]
		}
	}
	return ""
}

// Params returns the map[key]value.
func (input *BeegoInput) Params() map[string]string {
	m := make(map[string]string)
	for i, v := range input.pnames {
		if i <= len(input.pvalues) {
			m[v] = input.pvalues[i]
		}
	}
	return m
}

// SetParam sets the param with key and value
func (input *BeegoInput) SetParam(key, val string) {
	// check if already exists
	for i, v := range input.pnames {
		if v == key && i <= len(input.pvalues) {
			input.pvalues[i] = val
			return
		}
	}
	input.pvalues = append(input.pvalues, val)
	input.pnames = append(input.pnames, key)
}

// ResetParams clears any of the input's params
// Used to clear parameters so they may be reset between filter passes.
func (input *BeegoInput) ResetParams() {
	input.pnames = input.pnames[:0]
	input.pvalues = input.pvalues[:0]
}

// Query returns input data item string by a given string.
func (input *BeegoInput) Query(key string) string {
	if val := input.Param(key); val != "" {
		return val
	}
	if input.Context.Request.Form == nil {
		input.dataLock.Lock()
		if input.Context.Request.Form == nil {
			input.Context.Request.ParseForm()
		}
		input.dataLock.Unlock()
	}
	input.dataLock.RLock()
	defer input.dataLock.RUnlock()
	return input.Context.Request.Form.Get(key)
}

// Header returns request header item string by a given string.
// if non-existed, return empty string.
func (input *BeegoInput) Header(key string) string {
	return input.Context.Request.Header.Get(key)
}

// Cookie returns request cookie item string by a given key.
// if non-existed, return empty string.
func (input *BeegoInput) Cookie(key string) string {
	ck, err := input.Context.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

// Session returns current session item value by a given key.
// if non-existed, return nil.
func (input *BeegoInput) Session(key interface{}) interface{} {
	return input.CruSession.Get(nil, key)
}

// CopyBody returns the raw request body data as bytes.
func (input *BeegoInput) CopyBody(MaxMemory int64) []byte {
	if input.Context.Request.Body == nil {
		return []byte{}
	}

	var requestbody []byte
	safe := &io.LimitedReader{R: input.Context.Request.Body, N: MaxMemory}
	if input.Header("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(safe)
		if err != nil {
			return nil
		}
		requestbody, _ = ioutil.ReadAll(reader)
	} else {
		requestbody, _ = ioutil.ReadAll(safe)
	}

	input.Context.Request.Body.Close()
	bf := bytes.NewBuffer(requestbody)
	input.Context.Request.Body = http.MaxBytesReader(input.Context.ResponseWriter, ioutil.NopCloser(bf), MaxMemory)
	input.RequestBody = requestbody
	return requestbody
}

// Data returns the implicit data in the input
func (input *BeegoInput) Data() map[interface{}]interface{} {
	input.dataLock.Lock()
	defer input.dataLock.Unlock()
	if input.data == nil {
		input.data = make(map[interface{}]interface{})
	}
	return input.data
}

// GetData returns the stored data in this context.
func (input *BeegoInput) GetData(key interface{}) interface{} {
	input.dataLock.Lock()
	defer input.dataLock.Unlock()
	if v, ok := input.data[key]; ok {
		return v
	}
	return nil
}

// SetData stores data with given key in this context.
// This data is only available in this context.
func (input *BeegoInput) SetData(key, val interface{}) {
	input.dataLock.Lock()
	defer input.dataLock.Unlock()
	if input.data == nil {
		input.data = make(map[interface{}]interface{})
	}
	input.data[key] = val
}

// ParseFormOrMulitForm parseForm or parseMultiForm based on Content-type
func (input *BeegoInput) ParseFormOrMulitForm(maxMemory int64) error {
	// Parse the body depending on the content type.
	if strings.Contains(input.Header("Content-Type"), "multipart/form-data") {
		if err := input.Context.Request.ParseMultipartForm(maxMemory); err != nil {
			return errors.New("Error parsing request body:" + err.Error())
		}
	} else if err := input.Context.Request.ParseForm(); err != nil {
		return errors.New("Error parsing request body:" + err.Error())
	}
	return nil
}

// Bind data from request.Form[key] to dest
// like /?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=astaxie
// var id int  beegoInput.Bind(&id, "id")  id ==123
// var isok bool  beegoInput.Bind(&isok, "isok")  isok ==true
// var ft float64  beegoInput.Bind(&ft, "ft")  ft ==1.2
// ol := make([]int, 0, 2)  beegoInput.Bind(&ol, "ol")  ol ==[1 2]
// ul := make([]string, 0, 2)  beegoInput.Bind(&ul, "ul")  ul ==[str array]
// user struct{Name}  beegoInput.Bind(&user, "user")  user == {Name:"astaxie"}
func (input *BeegoInput) Bind(dest interface{}, key string) error {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return errors.New("beego: non-pointer passed to Bind: " + key)
	}
	value = value.Elem()
	if !value.CanSet() {
		return errors.New("beego: non-settable variable passed to Bind: " + key)
	}
	typ := value.Type()
	// Get real type if dest define with interface{}.
	// e.g  var dest interface{} dest=1.0
	if value.Kind() == reflect.Interface {
		typ = value.Elem().Type()
	}
	rv := input.bind(key, typ)
	if !rv.IsValid() {
		return errors.New("beego: reflect value is empty")
	}
	value.Set(rv)
	return nil
}

func (input *BeegoInput) bind(key string, typ reflect.Type) reflect.Value {
	if input.Context.Request.Form == nil {
		input.Context.Request.ParseForm()
	}
	rv := reflect.Zero(typ)
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := input.Query(key)
		if len(val) == 0 {
			return rv
		}
		rv = input.bindInt(val, typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := input.Query(key)
		if len(val) == 0 {
			return rv
		}
		rv = input.bindUint(val, typ)
	case reflect.Float32, reflect.Float64:
		val := input.Query(key)
		if len(val) == 0 {
			return rv
		}
		rv = input.bindFloat(val, typ)
	case reflect.String:
		val := input.Query(key)
		if len(val) == 0 {
			return rv
		}
		rv = input.bindString(val, typ)
	case reflect.Bool:
		val := input.Query(key)
		if len(val) == 0 {
			return rv
		}
		rv = input.bindBool(val, typ)
	case reflect.Slice:
		rv = input.bindSlice(&input.Context.Request.Form, key, typ)
	case reflect.Struct:
		rv = input.bindStruct(&input.Context.Request.Form, key, typ)
	case reflect.Ptr:
		rv = input.bindPoint(key, typ)
	case reflect.Map:
		rv = input.bindMap(&input.Context.Request.Form, key, typ)
	}
	return rv
}

func (input *BeegoInput) bindValue(val string, typ reflect.Type) reflect.Value {
	rv := reflect.Zero(typ)
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv = input.bindInt(val, typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv = input.bindUint(val, typ)
	case reflect.Float32, reflect.Float64:
		rv = input.bindFloat(val, typ)
	case reflect.String:
		rv = input.bindString(val, typ)
	case reflect.Bool:
		rv = input.bindBool(val, typ)
	case reflect.Slice:
		rv = input.bindSlice(&url.Values{"": {val}}, "", typ)
	case reflect.Struct:
		rv = input.bindStruct(&url.Values{"": {val}}, "", typ)
	case reflect.Ptr:
		rv = input.bindPoint(val, typ)
	case reflect.Map:
		rv = input.bindMap(&url.Values{"": {val}}, "", typ)
	}
	return rv
}

func (input *BeegoInput) bindInt(val string, typ reflect.Type) reflect.Value {
	intValue, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetInt(intValue)
	return pValue.Elem()
}

func (input *BeegoInput) bindUint(val string, typ reflect.Type) reflect.Value {
	uintValue, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetUint(uintValue)
	return pValue.Elem()
}

func (input *BeegoInput) bindFloat(val string, typ reflect.Type) reflect.Value {
	floatValue, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetFloat(floatValue)
	return pValue.Elem()
}

func (input *BeegoInput) bindString(val string, typ reflect.Type) reflect.Value {
	return reflect.ValueOf(val)
}

func (input *BeegoInput) bindBool(val string, typ reflect.Type) reflect.Value {
	val = strings.TrimSpace(strings.ToLower(val))
	switch val {
	case "true", "on", "1":
		return reflect.ValueOf(true)
	}
	return reflect.ValueOf(false)
}

type sliceValue struct {
	index int           // Index extracted from brackets.  If -1, no index was provided.
	value reflect.Value // the bound value for this slice element.
}

func (input *BeegoInput) bindSlice(params *url.Values, key string, typ reflect.Type) reflect.Value {
	maxIndex := -1
	numNoIndex := 0
	sliceValues := []sliceValue{}
	for reqKey, vals := range *params {
		if !strings.HasPrefix(reqKey, key+"[") {
			continue
		}
		// Extract the index, and the index where a sub-key starts. (e.g. field[0].subkey)
		index := -1
		leftBracket, rightBracket := len(key), strings.Index(reqKey[len(key):], "]")+len(key)
		if rightBracket > leftBracket+1 {
			index, _ = strconv.Atoi(reqKey[leftBracket+1 : rightBracket])
		}
		subKeyIndex := rightBracket + 1

		// Handle the indexed case.
		if index > -1 {
			if index > maxIndex {
				maxIndex = index
			}
			sliceValues = append(sliceValues, sliceValue{
				index: index,
				value: input.bind(reqKey[:subKeyIndex], typ.Elem()),
			})
			continue
		}

		// It's an un-indexed element.  (e.g. element[])
		numNoIndex += len(vals)
		for _, val := range vals {
			// Unindexed values can only be direct-bound.
			sliceValues = append(sliceValues, sliceValue{
				index: -1,
				value: input.bindValue(val, typ.Elem()),
			})
		}
	}
	resultArray := reflect.MakeSlice(typ, maxIndex+1, maxIndex+1+numNoIndex)
	for _, sv := range sliceValues {
		if sv.index != -1 {
			resultArray.Index(sv.index).Set(sv.value)
		} else {
			resultArray = reflect.Append(resultArray, sv.value)
		}
	}
	return resultArray
}

func (input *BeegoInput) bindStruct(params *url.Values, key string, typ reflect.Type) reflect.Value {
	result := reflect.New(typ).Elem()
	fieldValues := make(map[string]reflect.Value)
	for reqKey, val := range *params {
		var fieldName string
		if strings.HasPrefix(reqKey, key+".") {
			fieldName = reqKey[len(key)+1:]
		} else if strings.HasPrefix(reqKey, key+"[") && reqKey[len(reqKey)-1] == ']' {
			fieldName = reqKey[len(key)+1 : len(reqKey)-1]
		} else {
			continue
		}

		if _, ok := fieldValues[fieldName]; !ok {
			// Time to bind this field.  Get it and make sure we can set it.
			fieldValue := result.FieldByName(fieldName)
			if !fieldValue.IsValid() {
				continue
			}
			if !fieldValue.CanSet() {
				continue
			}
			boundVal := input.bindValue(val[0], fieldValue.Type())
			fieldValue.Set(boundVal)
			fieldValues[fieldName] = boundVal
		}
	}

	return result
}

func (input *BeegoInput) bindPoint(key string, typ reflect.Type) reflect.Value {
	return input.bind(key, typ.Elem()).Addr()
}

func (input *BeegoInput) bindMap(params *url.Values, key string, typ reflect.Type) reflect.Value {
	var (
		result    = reflect.MakeMap(typ)
		keyType   = typ.Key()
		valueType = typ.Elem()
	)
	for paramName, values := range *params {
		if !strings.HasPrefix(paramName, key+"[") || paramName[len(paramName)-1] != ']' {
			continue
		}

		key := paramName[len(key)+1 : len(paramName)-1]
		result.SetMapIndex(input.bindValue(key, keyType), input.bindValue(values[0], valueType))
	}
	return result
}
